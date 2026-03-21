package domain

import "time"

// EntryStatus captures the lifecycle of a journal entry.
type EntryStatus string

const (
	EntryStatusDraft    EntryStatus = "draft"
	EntryStatusPosted   EntryStatus = "posted"
	EntryStatusReversed EntryStatus = "reversed"
)

// VoucherType captures the source document class used for Vietnamese voucher numbering.
type VoucherType string

const (
	VoucherTypeReceipt VoucherType = "PT"
	VoucherTypePayment VoucherType = "PC"
	VoucherTypeGeneral VoucherType = "PK"
)

// JournalEntry stores the journal header and its detail lines.
type JournalEntry struct {
	ID                string
	ReferenceNo       string
	VoucherType       VoucherType
	CompanyCode       string
	SourceModule      string
	ExternalRef       string
	Description       string
	CurrencyCode      string
	PostingDate       time.Time
	Status            EntryStatus
	Metadata          map[string]any
	ReversalOfEntryID string
	ReversalEntryID   string
	ReversedAt        *time.Time
	VoidReason        string
	PostedAt          time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Items             []JournalItem
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

// NormalizeVoucherType sanitizes a voucher type while keeping VAS-facing numbering conventions explicit.
func NormalizeVoucherType(value string) (VoucherType, bool) {
	switch VoucherType(value) {
	case VoucherTypeReceipt:
		return VoucherTypeReceipt, true
	case VoucherTypePayment:
		return VoucherTypePayment, true
	case VoucherTypeGeneral:
		return VoucherTypeGeneral, true
	default:
		return "", false
	}
}

// ReversalVoucherType returns the voucher type used when voiding an entry.
func ReversalVoucherType(original VoucherType) VoucherType {
	switch original {
	case VoucherTypeReceipt:
		return VoucherTypePayment
	case VoucherTypePayment:
		return VoucherTypeReceipt
	default:
		return VoucherTypeGeneral
	}
}
