package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"vct-platform/backend/internal/modules/ledger/domain"
	"vct-platform/backend/internal/shared/repository"
	"vct-platform/backend/internal/shared/sqltx"
)

type queryable interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Store implements the ledger repositories on top of PostgreSQL.
type Store struct {
	db *sql.DB
}

// NewStore builds the PostgreSQL store.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// WithinTransaction runs the callback inside a SQL transaction and propagates it through context.
func (s *Store) WithinTransaction(ctx context.Context, opts repository.TxOptions, fn func(ctx context.Context) error) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("postgres store is not configured")
	}
	if _, ok := sqltx.FromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: toSQLIsolationLevel(opts.Isolation),
		ReadOnly:  opts.ReadOnly,
	})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txCtx := sqltx.WithTx(ctx, tx)
	if err := fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("rollback after %v: %w", err, rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// GetByIDs loads chart-of-accounts nodes in a single query.
func (s *Store) GetByIDs(ctx context.Context, ids []string) (map[string]domain.Account, error) {
	result := make(map[string]domain.Account, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	query := fmt.Sprintf(`
SELECT
    id,
    company_code,
    code,
    name,
    parent_id,
    depth,
    account_type,
    normal_side,
    is_postable,
    is_active,
    created_at,
    updated_at
FROM accounts
WHERE id IN (%s)`, dollarPlaceholders(1, len(ids)))

	rows, err := s.current(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query accounts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			account     domain.Account
			parentID    sql.NullString
			accountType string
			normalSide  string
		)
		if err := rows.Scan(
			&account.ID,
			&account.CompanyCode,
			&account.Code,
			&account.Name,
			&parentID,
			&account.Depth,
			&accountType,
			&normalSide,
			&account.IsPostable,
			&account.IsActive,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}

		if parentID.Valid {
			account.ParentID = &parentID.String
		}
		account.Type = domain.AccountType(accountType)
		account.NormalSide = domain.Side(normalSide)
		result[account.ID] = account
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate accounts: %w", err)
	}

	return result, nil
}

// NextVoucherNumber allocates the next voucher sequence for the month and voucher type.
func (s *Store) NextVoucherNumber(ctx context.Context, companyCode string, voucherType domain.VoucherType, postingDate time.Time) (string, error) {
	periodKey := postingDate.UTC().Format("2006-01")
	now := time.Now().UTC()

	var nextValue int
	if err := s.current(ctx).QueryRowContext(ctx, `
INSERT INTO voucher_sequences (
    company_code,
    voucher_type,
    period_key,
    last_value,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, 1, $4, $4)
ON CONFLICT (company_code, voucher_type, period_key)
DO UPDATE SET
    last_value = voucher_sequences.last_value + 1,
    updated_at = EXCLUDED.updated_at
RETURNING last_value`,
		companyCode,
		string(voucherType),
		periodKey,
		now,
	).Scan(&nextValue); err != nil {
		return "", fmt.Errorf("allocate voucher number: %w", err)
	}

	return fmt.Sprintf("%s-%04d/%s", voucherType, nextValue, postingDate.UTC().Format("01-06")), nil
}

// CreateEntry inserts the journal header.
func (s *Store) CreateEntry(ctx context.Context, entry *domain.JournalEntry) error {
	if entry == nil {
		return fmt.Errorf("journal entry is required")
	}

	metadata, err := marshalJSON(entry.Metadata)
	if err != nil {
		return fmt.Errorf("marshal entry metadata: %w", err)
	}

	_, err = s.current(ctx).ExecContext(ctx, `
INSERT INTO journal_entries (
    id,
    entry_no,
    voucher_type,
    company_code,
    source_module,
    external_ref,
    description,
    currency_code,
    posting_date,
    status,
    posted_at,
    metadata,
    reversal_of_entry_id,
    void_reason,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CAST($12 AS JSONB), NULLIF($13, '')::uuid, NULLIF($14, ''), $15, $16)`,
		entry.ID,
		entry.ReferenceNo,
		string(entry.VoucherType),
		entry.CompanyCode,
		entry.SourceModule,
		entry.ExternalRef,
		entry.Description,
		entry.CurrencyCode,
		entry.PostingDate.Format("2006-01-02"),
		string(entry.Status),
		entry.PostedAt,
		metadata,
		entry.ReversalOfEntryID,
		entry.VoidReason,
		entry.CreatedAt,
		entry.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert journal entry %s: %w", entry.ID, err)
	}

	return nil
}

// GetEntry loads a journal entry and its detail lines.
func (s *Store) GetEntry(ctx context.Context, entryID string) (domain.JournalEntry, error) {
	return s.getEntry(ctx, entryID, false)
}

// GetEntryForUpdate loads and locks a journal entry for a privileged mutation.
func (s *Store) GetEntryForUpdate(ctx context.Context, entryID string) (domain.JournalEntry, error) {
	return s.getEntry(ctx, entryID, true)
}

// CreateItems inserts the journal lines and lets PostgreSQL route them to the right partition by created_at.
func (s *Store) CreateItems(ctx context.Context, entryID string, items []domain.JournalItem, createdAt time.Time, companyCode string, currencyCode string) error {
	for _, item := range items {
		item.JournalEntryID = entryID
		_, err := s.current(ctx).ExecContext(ctx, `
INSERT INTO journal_items (
    journal_entry_id,
    line_no,
    company_code,
    account_id,
    side,
    amount,
    amount_signed,
    currency_code,
    description,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			entryID,
			item.LineNo,
			companyCode,
			item.AccountID,
			string(item.Side),
			item.Amount.String(),
			item.SignedAmount().String(),
			currencyCode,
			item.Description,
			createdAt,
			createdAt,
		)
		if err != nil {
			return fmt.Errorf("insert journal item %d for %s: %w", item.LineNo, entryID, err)
		}
	}

	return nil
}

// MarkReversed ties the original entry to its reversal after the reversal posting succeeds.
func (s *Store) MarkReversed(ctx context.Context, originalEntryID string, reversalEntryID string, reversedAt time.Time, reason string) error {
	result, err := s.current(ctx).ExecContext(ctx, `
UPDATE journal_entries
SET
    status = 'reversed',
    reversal_entry_id = $2,
    reversed_at = $3,
    void_reason = NULLIF($4, ''),
    updated_at = $3
WHERE id = $1
  AND status = 'posted'
  AND reversal_entry_id IS NULL`,
		originalEntryID,
		reversalEntryID,
		reversedAt,
		reason,
	)
	if err != nil {
		return fmt.Errorf("mark journal entry %s reversed: %w", originalEntryID, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrEntryAlreadyReversed
	}

	return nil
}

// ApplyDeltas updates the real-time balance table inside the same transaction as the journal post.
func (s *Store) ApplyDeltas(ctx context.Context, deltas []domain.AccountBalanceDelta, entryID string, updatedAt time.Time) error {
	for _, delta := range deltas {
		_, err := s.current(ctx).ExecContext(ctx, `
INSERT INTO account_balances (
    company_code,
    account_id,
    currency_code,
    debit_balance,
    credit_balance,
    net_balance,
    last_journal_entry_id,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (company_code, account_id, currency_code)
DO UPDATE SET
    debit_balance = account_balances.debit_balance + EXCLUDED.debit_balance,
    credit_balance = account_balances.credit_balance + EXCLUDED.credit_balance,
    net_balance = account_balances.net_balance + EXCLUDED.net_balance,
    last_journal_entry_id = EXCLUDED.last_journal_entry_id,
    updated_at = EXCLUDED.updated_at`,
			delta.CompanyCode,
			delta.AccountID,
			delta.CurrencyCode,
			delta.DebitDelta.String(),
			delta.CreditDelta.String(),
			delta.NetDelta.String(),
			entryID,
			updatedAt,
		)
		if err != nil {
			return fmt.Errorf("upsert account balance %s: %w", delta.AccountID, err)
		}
	}

	return nil
}

// Enqueue stores an outbox event in the same transaction as the ledger write.
func (s *Store) Enqueue(ctx context.Context, event domain.OutboxEvent) error {
	availableAt := event.AvailableAt
	if availableAt.IsZero() {
		availableAt = event.CreatedAt
	}

	_, err := s.current(ctx).ExecContext(ctx, `
INSERT INTO outbox_events (
    id,
    aggregate_type,
    aggregate_id,
    event_type,
    stream_key,
    status,
    payload,
    attempt_count,
    available_at,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, CAST($7 AS JSONB), $8, $9, $10, $11)`,
		event.ID,
		event.AggregateType,
		event.AggregateID,
		event.EventType,
		event.StreamKey,
		string(event.Status),
		string(event.Payload),
		event.AttemptCount,
		availableAt,
		event.CreatedAt,
		event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert outbox event %s: %w", event.ID, err)
	}

	return nil
}

// MarkPublished updates the outbox status after a successful Redis Stream append.
func (s *Store) MarkPublished(ctx context.Context, eventID string, publishedAt time.Time) error {
	_, err := s.current(ctx).ExecContext(ctx, `
UPDATE outbox_events
SET
    status = 'published',
    locked_at = NULL,
    locked_by = NULL,
    last_error = NULL,
    published_at = $2,
    updated_at = $2
WHERE id = $1`, eventID, publishedAt)
	if err != nil {
		return fmt.Errorf("mark outbox event %s published: %w", eventID, err)
	}

	return nil
}

// ClaimPending leases a batch of outbox rows for relay processing.
func (s *Store) ClaimPending(ctx context.Context, workerID string, batchSize int, now time.Time, staleClaimThreshold time.Duration) ([]domain.OutboxEvent, error) {
	rows, err := s.current(ctx).QueryContext(ctx, `
WITH candidates AS (
    SELECT id
    FROM outbox_events
    WHERE
        (
            status IN ('pending', 'failed')
            AND available_at <= $2
        )
        OR
        (
            status = 'processing'
            AND locked_at IS NOT NULL
            AND locked_at <= $3
        )
    ORDER BY available_at, created_at
    LIMIT $1
    FOR UPDATE SKIP LOCKED
)
UPDATE outbox_events AS oe
SET
    status = 'processing',
    attempt_count = oe.attempt_count + 1,
    locked_at = $2,
    locked_by = $4,
    updated_at = $2
FROM candidates
WHERE oe.id = candidates.id
RETURNING
    oe.id,
    oe.aggregate_type,
    oe.aggregate_id,
    oe.event_type,
    oe.stream_key,
    oe.status,
    oe.payload,
    oe.attempt_count,
    oe.available_at,
    oe.last_error,
    oe.created_at,
    oe.locked_at,
    oe.locked_by,
    oe.published_at`,
		batchSize,
		now,
		now.Add(-staleClaimThreshold),
		workerID,
	)
	if err != nil {
		return nil, fmt.Errorf("claim outbox rows: %w", err)
	}
	defer rows.Close()

	events := make([]domain.OutboxEvent, 0, batchSize)
	for rows.Next() {
		event, err := scanOutboxEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate claimed outbox rows: %w", err)
	}

	return events, nil
}

// MarkFailed reschedules an outbox row after a publish failure.
func (s *Store) MarkFailed(ctx context.Context, eventID string, failedAt time.Time, lastError string, retryAt time.Time) error {
	_, err := s.current(ctx).ExecContext(ctx, `
UPDATE outbox_events
SET
    status = 'failed',
    locked_at = NULL,
    locked_by = NULL,
    last_error = $3,
    available_at = $4,
    updated_at = $2
WHERE id = $1`,
		eventID,
		failedAt,
		lastError,
		retryAt,
	)
	if err != nil {
		return fmt.Errorf("mark outbox event %s failed: %w", eventID, err)
	}

	return nil
}

func (s *Store) current(ctx context.Context) queryable {
	if tx, ok := sqltx.FromContext(ctx); ok {
		return tx
	}
	return s.db
}

func (s *Store) getEntry(ctx context.Context, entryID string, forUpdate bool) (domain.JournalEntry, error) {
	var (
		entry             domain.JournalEntry
		status            string
		voucherType       string
		metadataRaw       []byte
		reversalOfEntryID sql.NullString
		reversalEntryID   sql.NullString
		reversedAt        sql.NullTime
		voidReason        sql.NullString
		lockClause        string
	)

	if forUpdate {
		lockClause = " FOR UPDATE"
	}

	err := s.current(ctx).QueryRowContext(ctx, `
SELECT
    id,
    entry_no,
    voucher_type,
    company_code,
    source_module,
    external_ref,
    description,
    currency_code,
    posting_date,
    status,
    posted_at,
    COALESCE(metadata::text, '{}'),
    reversal_of_entry_id,
    reversal_entry_id,
    reversed_at,
    void_reason,
    created_at,
    updated_at
FROM journal_entries
WHERE id = $1`+lockClause,
		entryID,
	).Scan(
		&entry.ID,
		&entry.ReferenceNo,
		&voucherType,
		&entry.CompanyCode,
		&entry.SourceModule,
		&entry.ExternalRef,
		&entry.Description,
		&entry.CurrencyCode,
		&entry.PostingDate,
		&status,
		&entry.PostedAt,
		&metadataRaw,
		&reversalOfEntryID,
		&reversalEntryID,
		&reversedAt,
		&voidReason,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.JournalEntry{}, domain.ErrJournalEntryNotFound
		}
		return domain.JournalEntry{}, fmt.Errorf("query journal entry %s: %w", entryID, err)
	}

	entry.Status = domain.EntryStatus(status)
	entry.VoucherType = domain.VoucherType(voucherType)
	if len(metadataRaw) > 0 && string(metadataRaw) != "null" {
		if err := json.Unmarshal(metadataRaw, &entry.Metadata); err != nil {
			return domain.JournalEntry{}, fmt.Errorf("decode entry metadata %s: %w", entryID, err)
		}
	}
	if reversalOfEntryID.Valid {
		entry.ReversalOfEntryID = reversalOfEntryID.String
	}
	if reversalEntryID.Valid {
		entry.ReversalEntryID = reversalEntryID.String
	}
	if reversedAt.Valid {
		value := reversedAt.Time
		entry.ReversedAt = &value
	}
	if voidReason.Valid {
		entry.VoidReason = voidReason.String
	}

	rows, err := s.current(ctx).QueryContext(ctx, `
SELECT
    journal_entry_id,
    line_no,
    account_id,
    side,
    amount,
    description
FROM journal_items
WHERE journal_entry_id = $1
ORDER BY line_no, created_at`,
		entryID,
	)
	if err != nil {
		return domain.JournalEntry{}, fmt.Errorf("query journal items for %s: %w", entryID, err)
	}
	defer rows.Close()

	entry.Items = make([]domain.JournalItem, 0, 4)
	for rows.Next() {
		var (
			item      domain.JournalItem
			side      string
			amountRaw string
		)
		if err := rows.Scan(
			&item.JournalEntryID,
			&item.LineNo,
			&item.AccountID,
			&side,
			&amountRaw,
			&item.Description,
		); err != nil {
			return domain.JournalEntry{}, fmt.Errorf("scan journal item for %s: %w", entryID, err)
		}

		amount, err := domain.ParseMoney(amountRaw)
		if err != nil {
			return domain.JournalEntry{}, fmt.Errorf("parse journal item amount for %s: %w", entryID, err)
		}
		item.Side = domain.Side(side)
		item.Amount = amount
		entry.Items = append(entry.Items, item)
	}

	if err := rows.Err(); err != nil {
		return domain.JournalEntry{}, fmt.Errorf("iterate journal items for %s: %w", entryID, err)
	}

	return entry, nil
}

func dollarPlaceholders(start int, count int) string {
	values := make([]string, count)
	for index := 0; index < count; index++ {
		values[index] = fmt.Sprintf("$%d", start+index)
	}
	return strings.Join(values, ", ")
}

func toSQLIsolationLevel(level repository.IsolationLevel) sql.IsolationLevel {
	switch level {
	case repository.IsolationRepeatableRead:
		return sql.LevelRepeatableRead
	case repository.IsolationSerializable:
		return sql.LevelSerializable
	default:
		return sql.LevelReadCommitted
	}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanOutboxEvent(row rowScanner) (domain.OutboxEvent, error) {
	var (
		event       domain.OutboxEvent
		status      string
		payload     []byte
		lastError   sql.NullString
		lockedAt    sql.NullTime
		lockedBy    sql.NullString
		publishedAt sql.NullTime
	)

	if err := row.Scan(
		&event.ID,
		&event.AggregateType,
		&event.AggregateID,
		&event.EventType,
		&event.StreamKey,
		&status,
		&payload,
		&event.AttemptCount,
		&event.AvailableAt,
		&lastError,
		&event.CreatedAt,
		&lockedAt,
		&lockedBy,
		&publishedAt,
	); err != nil {
		return domain.OutboxEvent{}, fmt.Errorf("scan outbox event: %w", err)
	}

	event.Status = domain.OutboxStatus(status)
	event.Payload = payload
	if lastError.Valid {
		event.LastError = lastError.String
	}
	if lockedAt.Valid {
		value := lockedAt.Time
		event.LockedAt = &value
	}
	if lockedBy.Valid {
		value := lockedBy.String
		event.LockedBy = &value
	}
	if publishedAt.Valid {
		value := publishedAt.Time
		event.PublishedAt = &value
	}

	return event, nil
}

func marshalJSON(payload any) (string, error) {
	if payload == nil {
		return "{}", nil
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
