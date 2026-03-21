package domain

import "time"

// EntryStatus captures the lifecycle of a journal entry.
type EntryStatus string

const (
	EntryStatusDraft    EntryStatus = "draft"
	EntryStatusPosted   EntryStatus = "posted"
	EntryStatusReversed EntryStatus = "reversed"
)

// JournalEntry stores the journal header and its detail lines.
type JournalEntry struct {
	ID           string
	ReferenceNo  string
	CompanyCode  string
	SourceModule string
	ExternalRef  string
	Description  string
	CurrencyCode string
	PostingDate  time.Time
	Status       EntryStatus
	PostedAt     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Items        []JournalItem
}

// JournalItem stores a debit or credit line.
type JournalItem struct {
	JournalEntryID string
	LineNo         int
	AccountID      string
	Side           Side
	Amount         Money
	Description    string
}

// SignedAmount turns credits into negative values for balance deltas and analytics.
func (i JournalItem) SignedAmount() Money {
	if i.Side == SideDebit {
		return i.Amount
	}

	return Money{}.Sub(i.Amount)
}

// AccountBalanceDelta captures the net effect of a posted journal entry on a balance row.
type AccountBalanceDelta struct {
	CompanyCode  string
	AccountID    string
	CurrencyCode string
	DebitDelta   Money
	CreditDelta  Money
	NetDelta     Money
}

// OutboxStatus tracks the dispatch state of an integration event.
type OutboxStatus string

const (
	OutboxStatusPending    OutboxStatus = "pending"
	OutboxStatusProcessing OutboxStatus = "processing"
	OutboxStatusPublished  OutboxStatus = "published"
	OutboxStatusFailed     OutboxStatus = "failed"
)

// OutboxEvent is persisted transactionally with the journal entry and later published.
type OutboxEvent struct {
	ID            string
	AggregateType string
	AggregateID   string
	EventType     string
	StreamKey     string
	Payload       []byte
	Status        OutboxStatus
	AttemptCount  int
	AvailableAt   time.Time
	LastError     string
	CreatedAt     time.Time
	LockedAt      *time.Time
	LockedBy      *string
	PublishedAt   *time.Time
}
