package postgres

import (
	"context"
	"database/sql"
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

// CreateEntry inserts the journal header.
func (s *Store) CreateEntry(ctx context.Context, entry *domain.JournalEntry) error {
	if entry == nil {
		return fmt.Errorf("journal entry is required")
	}

	_, err := s.current(ctx).ExecContext(ctx, `
INSERT INTO journal_entries (
    id,
    entry_no,
    company_code,
    source_module,
    external_ref,
    description,
    currency_code,
    posting_date,
    status,
    posted_at,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		entry.ID,
		entry.ReferenceNo,
		entry.CompanyCode,
		entry.SourceModule,
		entry.ExternalRef,
		entry.Description,
		entry.CurrencyCode,
		entry.PostingDate.Format("2006-01-02"),
		string(entry.Status),
		entry.PostedAt,
		entry.CreatedAt,
		entry.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert journal entry %s: %w", entry.ID, err)
	}

	return nil
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
