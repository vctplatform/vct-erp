package domain

import (
	"context"
	"time"

	"vct-platform/backend/internal/shared/repository"
)

// AccountCatalogRepository provides chart-of-accounts lookups.
type AccountCatalogRepository interface {
	repository.BulkReader[Account, string]
}

// JournalEntryRepository persists journal headers and items.
type JournalEntryRepository interface {
	CreateEntry(ctx context.Context, entry *JournalEntry) error
	CreateItems(ctx context.Context, entryID string, items []JournalItem, createdAt time.Time, companyCode string, currencyCode string) error
}

// AccountBalanceRepository keeps the real-time balance table synchronized.
type AccountBalanceRepository interface {
	ApplyDeltas(ctx context.Context, deltas []AccountBalanceDelta, entryID string, updatedAt time.Time) error
}

// OutboxRepository persists and updates integration events.
type OutboxRepository interface {
	Enqueue(ctx context.Context, event OutboxEvent) error
	MarkPublished(ctx context.Context, eventID string, publishedAt time.Time) error
}

// OutboxRelayRepository supports the asynchronous relay worker.
type OutboxRelayRepository interface {
	ClaimPending(ctx context.Context, workerID string, batchSize int, now time.Time, staleClaimThreshold time.Duration) ([]OutboxEvent, error)
	MarkPublished(ctx context.Context, eventID string, publishedAt time.Time) error
	MarkFailed(ctx context.Context, eventID string, failedAt time.Time, lastError string, retryAt time.Time) error
}

// ChartOfAccountsCache caches the chart-of-accounts lookup in Redis.
type ChartOfAccountsCache interface {
	GetMany(ctx context.Context, ids []string) (map[string]Account, error)
	SetMany(ctx context.Context, accounts []Account, ttl time.Duration) error
}

// EventPublisher pushes committed outbox events to downstream infrastructure.
type EventPublisher interface {
	Publish(ctx context.Context, event OutboxEvent) error
}

// TransactionManager executes the posting flow inside a single ACID transaction.
type TransactionManager interface {
	repository.TxManager
}
