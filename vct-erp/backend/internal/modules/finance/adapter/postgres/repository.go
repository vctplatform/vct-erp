package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
	"vct-platform/backend/internal/shared/sqltx"
)

type queryable interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Store implements the finance repositories using database/sql and the shared transaction context.
type Store struct {
	db *sql.DB
}

// NewStore constructs the finance Postgres adapter.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Reserve reserves an idempotency key or replays the cached response when the request was already completed.
func (s *Store) Reserve(ctx context.Context, scope string, key string, requestHash string, now time.Time) (financedomain.IdempotencyReservation, error) {
	result, err := s.current(ctx).ExecContext(ctx, `
INSERT INTO idempotency_keys (
    scope,
    idempotency_key,
    request_hash,
    status,
    locked_at,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, 'processing', $4, $4, $4)
ON CONFLICT (scope, idempotency_key) DO NOTHING`,
		scope,
		key,
		requestHash,
		now,
	)
	if err != nil {
		return financedomain.IdempotencyReservation{}, fmt.Errorf("insert idempotency key: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
		return financedomain.IdempotencyReservation{Status: financedomain.IdempotencyStatusAcquired, RequestHash: requestHash}, nil
	}

	var (
		existingHash string
		status       string
		resourceID   sql.NullString
		responseRaw  []byte
	)

	if err := s.current(ctx).QueryRowContext(ctx, `
SELECT request_hash, status, resource_id, COALESCE(response_payload::text, '{}')
FROM idempotency_keys
WHERE scope = $1 AND idempotency_key = $2`,
		scope,
		key,
	).Scan(&existingHash, &status, &resourceID, &responseRaw); err != nil {
		return financedomain.IdempotencyReservation{}, fmt.Errorf("load idempotency key: %w", err)
	}

	if existingHash != requestHash {
		return financedomain.IdempotencyReservation{Status: financedomain.IdempotencyStatusConflict}, nil
	}

	switch status {
	case "completed":
		return financedomain.IdempotencyReservation{
			Status:          financedomain.IdempotencyStatusReplay,
			RequestHash:     existingHash,
			ResourceID:      resourceID.String,
			ResponsePayload: responseRaw,
		}, nil
	case "failed":
		if _, err := s.current(ctx).ExecContext(ctx, `
UPDATE idempotency_keys
SET
    status = 'processing',
    request_hash = $3,
    locked_at = $4,
    last_error = NULL,
    updated_at = $4
WHERE scope = $1 AND idempotency_key = $2`,
			scope,
			key,
			requestHash,
			now,
		); err != nil {
			return financedomain.IdempotencyReservation{}, fmt.Errorf("reopen failed idempotency key: %w", err)
		}
		return financedomain.IdempotencyReservation{Status: financedomain.IdempotencyStatusAcquired, RequestHash: requestHash}, nil
	default:
		return financedomain.IdempotencyReservation{Status: financedomain.IdempotencyStatusInProgress, RequestHash: existingHash}, nil
	}
}

// Complete stores the response payload for future idempotent replays.
func (s *Store) Complete(ctx context.Context, scope string, key string, responsePayload []byte, resourceID string, completedAt time.Time) error {
	_, err := s.current(ctx).ExecContext(ctx, `
UPDATE idempotency_keys
SET
    status = 'completed',
    resource_id = NULLIF($3, ''),
    response_payload = CAST($4 AS JSONB),
    completed_at = $5,
    locked_at = $5,
    last_error = NULL,
    updated_at = $5
WHERE scope = $1 AND idempotency_key = $2`,
		scope,
		key,
		resourceID,
		string(responsePayload),
		completedAt,
	)
	if err != nil {
		return fmt.Errorf("complete idempotency key %s/%s: %w", scope, key, err)
	}
	return nil
}

// Fail stores the latest failure for the idempotency key.
func (s *Store) Fail(ctx context.Context, scope string, key string, lastError string, failedAt time.Time) error {
	_, err := s.current(ctx).ExecContext(ctx, `
UPDATE idempotency_keys
SET
    status = 'failed',
    last_error = $3,
    locked_at = $4,
    updated_at = $4
WHERE scope = $1 AND idempotency_key = $2`,
		scope,
		key,
		lastError,
		failedAt,
	)
	if err != nil {
		return fmt.Errorf("mark idempotency key %s/%s failed: %w", scope, key, err)
	}
	return nil
}

// CreateContract inserts a SaaS contract definition.
func (s *Store) CreateContract(ctx context.Context, contract financedomain.SaaSContract) error {
	_, err := s.current(ctx).ExecContext(ctx, `
INSERT INTO saas_contracts (
    id,
    company_code,
    contract_no,
    customer_ref,
    cash_account_id,
    deferred_revenue_account_id,
    recognized_revenue_account_id,
    currency_code,
    start_date,
    end_date,
    term_months,
    total_amount,
    source_ref,
    initial_journal_entry_id,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NULLIF($14, '')::uuid, $15, $16)`,
		contract.ID,
		contract.CompanyCode,
		contract.ContractNo,
		contract.CustomerRef,
		contract.CashAccountID,
		contract.DeferredRevenueAccountID,
		contract.RecognizedRevenueAccountID,
		contract.CurrencyCode,
		contract.StartDate.Format("2006-01-02"),
		contract.EndDate.Format("2006-01-02"),
		contract.TermMonths,
		contract.TotalAmount.String(),
		contract.SourceRef,
		contract.InitialJournalEntryID,
		contract.CreatedAt,
		contract.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert saas contract %s: %w", contract.ID, err)
	}
	return nil
}

// CreateSchedules inserts SaaS revenue schedule rows.
func (s *Store) CreateSchedules(ctx context.Context, schedules []financedomain.RevenueSchedule) error {
	for _, schedule := range schedules {
		_, err := s.current(ctx).ExecContext(ctx, `
INSERT INTO saas_revenue_schedules (
    id,
    contract_id,
    sequence_no,
    service_month,
    amount,
    status,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			schedule.ID,
			schedule.ContractID,
			schedule.SequenceNo,
			schedule.ServiceMonth.Format("2006-01-02"),
			schedule.Amount.String(),
			schedule.Status,
			schedule.CreatedAt,
			schedule.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert saas schedule %s: %w", schedule.ID, err)
		}
	}
	return nil
}

// ListDueSchedules returns scheduled SaaS revenue rows that are ready for recognition.
func (s *Store) ListDueSchedules(ctx context.Context, companyCode string, upTo time.Time, limit int) ([]financedomain.DueRevenueSchedule, error) {
	rows, err := s.current(ctx).QueryContext(ctx, `
SELECT
    schedule.id,
    contract.id,
    contract.contract_no,
    contract.company_code,
    contract.customer_ref,
    contract.currency_code,
    schedule.service_month,
    schedule.amount,
    contract.deferred_revenue_account_id,
    contract.recognized_revenue_account_id
FROM saas_revenue_schedules AS schedule
INNER JOIN saas_contracts AS contract ON contract.id = schedule.contract_id
WHERE contract.company_code = $1
  AND schedule.status = 'scheduled'
  AND schedule.service_month <= $2
ORDER BY schedule.service_month, schedule.sequence_no
LIMIT $3`,
		companyCode,
		upTo.Format("2006-01-02"),
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query due saas schedules: %w", err)
	}
	defer rows.Close()

	schedules := make([]financedomain.DueRevenueSchedule, 0, limit)
	for rows.Next() {
		var (
			schedule  financedomain.DueRevenueSchedule
			amountRaw string
		)
		if err := rows.Scan(
			&schedule.ScheduleID,
			&schedule.ContractID,
			&schedule.ContractNo,
			&schedule.CompanyCode,
			&schedule.CustomerRef,
			&schedule.CurrencyCode,
			&schedule.ServiceMonth,
			&amountRaw,
			&schedule.DeferredRevenueAccountID,
			&schedule.RecognizedRevenueAccountID,
		); err != nil {
			return nil, fmt.Errorf("scan due saas schedule: %w", err)
		}

		amount, err := ledgerdomain.ParseMoney(amountRaw)
		if err != nil {
			return nil, fmt.Errorf("parse saas schedule amount: %w", err)
		}
		schedule.Amount = amount
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate due saas schedules: %w", err)
	}
	return schedules, nil
}

// MarkScheduleRecognized marks a SaaS schedule row as recognized.
func (s *Store) MarkScheduleRecognized(ctx context.Context, scheduleID string, journalEntryID string, recognizedAt time.Time) error {
	_, err := s.current(ctx).ExecContext(ctx, `
UPDATE saas_revenue_schedules
SET
    status = 'recognized',
    recognized_journal_entry_id = $2,
    recognized_at = $3,
    updated_at = $3
WHERE id = $1`,
		scheduleID,
		journalEntryID,
		recognizedAt,
	)
	if err != nil {
		return fmt.Errorf("mark saas schedule %s recognized: %w", scheduleID, err)
	}
	return nil
}

// CreateReceivable inserts a dojo receivable row.
func (s *Store) CreateReceivable(ctx context.Context, receivable financedomain.DojoReceivable) error {
	_, err := s.current(ctx).ExecContext(ctx, `
INSERT INTO dojo_receivables (
    id,
    company_code,
    student_ref,
    billing_month,
    due_date,
    currency_code,
    receivable_account_id,
    revenue_account_id,
    amount_due,
    amount_paid,
    status,
    source_ref,
    assessment_entry_id,
    settlement_entry_id,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NULLIF($13, '')::uuid, NULLIF($14, '')::uuid, $15, $16)`,
		receivable.ID,
		receivable.CompanyCode,
		receivable.StudentRef,
		receivable.BillingMonth.Format("2006-01-02"),
		receivable.DueDate.Format("2006-01-02"),
		receivable.CurrencyCode,
		receivable.ReceivableAccountID,
		receivable.RevenueAccountID,
		receivable.AmountDue.String(),
		receivable.AmountPaid.String(),
		receivable.Status,
		receivable.SourceRef,
		receivable.AssessmentEntryID,
		receivable.SettlementEntryID,
		receivable.CreatedAt,
		receivable.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert dojo receivable %s: %w", receivable.ID, err)
	}
	return nil
}

// GetReceivable loads a dojo receivable by student and billing month.
func (s *Store) GetReceivable(ctx context.Context, companyCode string, studentRef string, billingMonth time.Time) (financedomain.DojoReceivable, error) {
	var (
		receivable        financedomain.DojoReceivable
		amountDueRaw      string
		amountPaidRaw     string
		assessmentEntryID sql.NullString
		settlementEntryID sql.NullString
	)

	err := s.current(ctx).QueryRowContext(ctx, `
SELECT
    id,
    company_code,
    student_ref,
    billing_month,
    due_date,
    currency_code,
    receivable_account_id,
    revenue_account_id,
    amount_due,
    amount_paid,
    status,
    source_ref,
    assessment_entry_id,
    settlement_entry_id,
    created_at,
    updated_at
FROM dojo_receivables
WHERE company_code = $1 AND student_ref = $2 AND billing_month = $3`,
		companyCode,
		studentRef,
		billingMonth.Format("2006-01-02"),
	).Scan(
		&receivable.ID,
		&receivable.CompanyCode,
		&receivable.StudentRef,
		&receivable.BillingMonth,
		&receivable.DueDate,
		&receivable.CurrencyCode,
		&receivable.ReceivableAccountID,
		&receivable.RevenueAccountID,
		&amountDueRaw,
		&amountPaidRaw,
		&receivable.Status,
		&receivable.SourceRef,
		&assessmentEntryID,
		&settlementEntryID,
		&receivable.CreatedAt,
		&receivable.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return financedomain.DojoReceivable{}, nil
		}
		return financedomain.DojoReceivable{}, fmt.Errorf("query dojo receivable: %w", err)
	}

	amountDue, err := ledgerdomain.ParseMoney(amountDueRaw)
	if err != nil {
		return financedomain.DojoReceivable{}, fmt.Errorf("parse dojo amount due: %w", err)
	}
	amountPaid, err := ledgerdomain.ParseMoney(amountPaidRaw)
	if err != nil {
		return financedomain.DojoReceivable{}, fmt.Errorf("parse dojo amount paid: %w", err)
	}
	receivable.AmountDue = amountDue
	receivable.AmountPaid = amountPaid
	receivable.AssessmentEntryID = assessmentEntryID.String
	receivable.SettlementEntryID = settlementEntryID.String
	return receivable, nil
}

// MarkSettled increments the amount paid and marks the receivable settled when fully paid.
func (s *Store) MarkSettled(ctx context.Context, receivableID string, amountPaid ledgerdomain.Money, journalEntryID string, settledAt time.Time) error {
	_, err := s.current(ctx).ExecContext(ctx, `
UPDATE dojo_receivables
SET
    amount_paid = amount_paid + CAST($2 AS NUMERIC(20, 4)),
    status = CASE
        WHEN amount_paid + CAST($2 AS NUMERIC(20, 4)) >= amount_due THEN 'settled'
        ELSE 'open'
    END,
    settlement_entry_id = $3,
    updated_at = $4
WHERE id = $1`,
		receivableID,
		amountPaid.String(),
		journalEntryID,
		settledAt,
	)
	if err != nil {
		return fmt.Errorf("mark dojo receivable %s settled: %w", receivableID, err)
	}
	return nil
}

// CreateDeposit inserts a rental deposit record.
func (s *Store) CreateDeposit(ctx context.Context, deposit financedomain.RentalDeposit) error {
	_, err := s.current(ctx).ExecContext(ctx, `
INSERT INTO rental_deposits (
    id,
    company_code,
    rental_order_id,
    customer_ref,
    cash_account_id,
    holding_account_id,
    currency_code,
    amount,
    status,
    held_entry_id,
    released_entry_id,
    held_at,
    released_at,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NULLIF($10, '')::uuid, NULLIF($11, '')::uuid, $12, $13, $14, $15)`,
		deposit.ID,
		deposit.CompanyCode,
		deposit.RentalOrderID,
		deposit.CustomerRef,
		deposit.CashAccountID,
		deposit.HoldingAccountID,
		deposit.CurrencyCode,
		deposit.Amount.String(),
		deposit.Status,
		deposit.HeldEntryID,
		deposit.ReleasedEntryID,
		deposit.HeldAt,
		deposit.ReleasedAt,
		deposit.CreatedAt,
		deposit.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert rental deposit %s: %w", deposit.ID, err)
	}
	return nil
}

// GetByRentalOrder loads a rental deposit by rental order identifier.
func (s *Store) GetByRentalOrder(ctx context.Context, companyCode string, rentalOrderID string) (financedomain.RentalDeposit, error) {
	var (
		deposit         financedomain.RentalDeposit
		amountRaw       string
		releasedEntryID sql.NullString
		releasedAt      sql.NullTime
	)

	err := s.current(ctx).QueryRowContext(ctx, `
SELECT
    id,
    company_code,
    rental_order_id,
    customer_ref,
    cash_account_id,
    holding_account_id,
    currency_code,
    amount,
    status,
    held_entry_id,
    released_entry_id,
    held_at,
    released_at,
    created_at,
    updated_at
FROM rental_deposits
WHERE company_code = $1 AND rental_order_id = $2`,
		companyCode,
		rentalOrderID,
	).Scan(
		&deposit.ID,
		&deposit.CompanyCode,
		&deposit.RentalOrderID,
		&deposit.CustomerRef,
		&deposit.CashAccountID,
		&deposit.HoldingAccountID,
		&deposit.CurrencyCode,
		&amountRaw,
		&deposit.Status,
		&deposit.HeldEntryID,
		&releasedEntryID,
		&deposit.HeldAt,
		&releasedAt,
		&deposit.CreatedAt,
		&deposit.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return financedomain.RentalDeposit{}, nil
		}
		return financedomain.RentalDeposit{}, fmt.Errorf("query rental deposit: %w", err)
	}

	amount, err := ledgerdomain.ParseMoney(amountRaw)
	if err != nil {
		return financedomain.RentalDeposit{}, fmt.Errorf("parse rental deposit amount: %w", err)
	}
	deposit.Amount = amount
	deposit.ReleasedEntryID = releasedEntryID.String
	if releasedAt.Valid {
		value := releasedAt.Time
		deposit.ReleasedAt = &value
	}
	return deposit, nil
}

// MarkReleased marks a rental deposit as released.
func (s *Store) MarkReleased(ctx context.Context, depositID string, releasedJournalEntryID string, releasedAt time.Time) error {
	_, err := s.current(ctx).ExecContext(ctx, `
UPDATE rental_deposits
SET
    status = 'released',
    released_entry_id = $2,
    released_at = $3,
    updated_at = $3
WHERE id = $1`,
		depositID,
		releasedJournalEntryID,
		releasedAt,
	)
	if err != nil {
		return fmt.Errorf("mark rental deposit %s released: %w", depositID, err)
	}
	return nil
}

func (s *Store) current(ctx context.Context) queryable {
	if tx, ok := sqltx.FromContext(ctx); ok {
		return tx
	}
	return s.db
}
