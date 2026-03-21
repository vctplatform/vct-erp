package domain

import (
	"encoding/json"
	"time"

	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

// BusinessLine identifies which VCT business vertical emitted the financial event.
type BusinessLine string

const (
	BusinessLineSaaS   BusinessLine = "saas"
	BusinessLineDojo   BusinessLine = "dojo"
	BusinessLineRetail BusinessLine = "retail"
	BusinessLineRental BusinessLine = "rental"
)

// CaptureOperation identifies the business action that should be transformed into ledger postings.
type CaptureOperation string

const (
	OperationSaaSCaptureAnnualContract CaptureOperation = "saas.capture_annual_contract"
	OperationSaaSRecognizeDueRevenue   CaptureOperation = "saas.recognize_due_revenue"
	OperationDojoAssessMonthlyTuition  CaptureOperation = "dojo.assess_monthly_tuition"
	OperationDojoCapturePayment        CaptureOperation = "dojo.capture_payment"
	OperationRetailCaptureSale         CaptureOperation = "retail.capture_sale"
	OperationRetailCaptureRefund       CaptureOperation = "retail.capture_refund"
	OperationRentalCaptureDeposit      CaptureOperation = "rental.capture_deposit"
	OperationRentalReleaseDeposit      CaptureOperation = "rental.release_deposit"
)

// IdempotencyReservationStatus reports how the idempotency repository handled the incoming key.
type IdempotencyReservationStatus string

const (
	IdempotencyStatusAcquired   IdempotencyReservationStatus = "acquired"
	IdempotencyStatusReplay     IdempotencyReservationStatus = "replay"
	IdempotencyStatusConflict   IdempotencyReservationStatus = "conflict"
	IdempotencyStatusInProgress IdempotencyReservationStatus = "in_progress"
)

// IdempotencyReservation contains the state returned when a request reserves an idempotency key.
type IdempotencyReservation struct {
	Status          IdempotencyReservationStatus
	RequestHash     string
	ResourceID      string
	ResponsePayload []byte
}

// SaaSContract tracks a prepaid SaaS agreement and its deferred revenue parameters.
type SaaSContract struct {
	ID                         string
	CompanyCode                string
	ContractNo                 string
	CustomerRef                string
	CashAccountID              string
	DeferredRevenueAccountID   string
	RecognizedRevenueAccountID string
	CurrencyCode               string
	StartDate                  time.Time
	EndDate                    time.Time
	TermMonths                 int
	TotalAmount                ledgerdomain.Money
	SourceRef                  string
	InitialJournalEntryID      string
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

// RevenueSchedule tracks a single SaaS recognition slice.
type RevenueSchedule struct {
	ID                       string
	ContractID               string
	SequenceNo               int
	ServiceMonth             time.Time
	Amount                   ledgerdomain.Money
	Status                   string
	RecognizedJournalEntryID string
	RecognizedAt             *time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

// DueRevenueSchedule contains the data needed to recognize a due SaaS slice.
type DueRevenueSchedule struct {
	ScheduleID                 string
	ContractID                 string
	ContractNo                 string
	CompanyCode                string
	CustomerRef                string
	CurrencyCode               string
	ServiceMonth               time.Time
	Amount                     ledgerdomain.Money
	DeferredRevenueAccountID   string
	RecognizedRevenueAccountID string
}

// DojoReceivable tracks a monthly tuition receivable for a student.
type DojoReceivable struct {
	ID                  string
	CompanyCode         string
	StudentRef          string
	BillingMonth        time.Time
	DueDate             time.Time
	CurrencyCode        string
	ReceivableAccountID string
	RevenueAccountID    string
	AmountDue           ledgerdomain.Money
	AmountPaid          ledgerdomain.Money
	Status              string
	SourceRef           string
	AssessmentEntryID   string
	SettlementEntryID   string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// RentalDeposit tracks a deposit held against a rental order.
type RentalDeposit struct {
	ID               string
	CompanyCode      string
	RentalOrderID    string
	CustomerRef      string
	CashAccountID    string
	HoldingAccountID string
	CurrencyCode     string
	Amount           ledgerdomain.Money
	Status           string
	HeldEntryID      string
	ReleasedEntryID  string
	HeldAt           time.Time
	ReleasedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CaptureRequest is the envelope accepted by the finance capture API.
type CaptureRequest struct {
	IdempotencyKey string
	BusinessLine   BusinessLine     `json:"business_line"`
	Operation      CaptureOperation `json:"operation"`
	Payload        json.RawMessage  `json:"payload"`
}

// CaptureResult returns the materialized response for an idempotent capture request.
type CaptureResult struct {
	BusinessLine BusinessLine     `json:"business_line"`
	Operation    CaptureOperation `json:"operation"`
	ResourceID   string           `json:"resource_id,omitempty"`
	Replay       bool             `json:"replay"`
	Payload      json.RawMessage  `json:"payload"`
}

// BankStatementLine stores imported bank statement movements for reconciliation with ledger account 112.
type BankStatementLine struct {
	ID             string
	CompanyCode    string
	BankAccountNo  string
	ExternalLineID string
	ReferenceNo    string
	Description    string
	CurrencyCode   string
	BookingDate    time.Time
	ValueDate      time.Time
	Amount         ledgerdomain.Money
	Status         string
	MatchedEntryID string
	MatchingRule   string
	MatchedAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// LedgerBankMovement represents a journal movement on a bank account that can be matched to a statement line.
type LedgerBankMovement struct {
	JournalEntryID string
	EntryNo        string
	AccountID      string
	AccountCode    string
	ExternalRef    string
	Description    string
	CurrencyCode   string
	PostingDate    time.Time
	Amount         ledgerdomain.Money
	Side           ledgerdomain.Side
}

// ReconcileBankRequest scopes a reconciliation batch for one bank ledger account.
type ReconcileBankRequest struct {
	CompanyCode      string
	LedgerAccountID  string
	BankAccountNo    string
	DateFrom         time.Time
	DateTo           time.Time
	MaxDateDriftDays int
	Limit            int
}

// BankMatchResult captures a single successful statement-to-ledger pairing.
type BankMatchResult struct {
	StatementLineID string `json:"statement_line_id"`
	JournalEntryID  string `json:"journal_entry_id"`
	EntryNo         string `json:"entry_no"`
	MatchingRule    string `json:"matching_rule"`
}

// ReconcileBankResult summarizes a reconciliation batch.
type ReconcileBankResult struct {
	MatchedCount            int               `json:"matched_count"`
	UnmatchedStatementCount int               `json:"unmatched_statement_count"`
	UnmatchedLedgerCount    int               `json:"unmatched_ledger_count"`
	Matches                 []BankMatchResult `json:"matches"`
}

// VoidAuditLog stores the before/after snapshots that must be preserved for inspections.
type VoidAuditLog struct {
	ID              string
	CompanyCode     string
	OriginalEntryID string
	ReversalEntryID string
	ActorID         string
	Reason          string
	VoidedAt        time.Time
	Before          any
	After           any
}

// VoidEntryResult exposes the reversal identifiers to privileged clients.
type VoidEntryResult struct {
	OriginalEntryID string `json:"original_entry_id"`
	ReversalEntryID string `json:"reversal_entry_id"`
	VoucherNo       string `json:"voucher_no"`
	Status          string `json:"status"`
}
