package usecase

import (
	"context"
	"time"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
	ledgerusecase "vct-platform/backend/internal/modules/ledger/usecase"
)

// LedgerPoster bridges finance business services to the core ledger posting use case.
type LedgerPoster interface {
	PostEntry(ctx context.Context, req ledgerusecase.PostEntryRequest) (*ledgerusecase.PostEntryResult, error)
}

// SaaSAccountingService exposes the SaaS-specific accounting operations used by the capture API.
type SaaSAccountingService interface {
	CaptureAnnualContract(ctx context.Context, req CaptureAnnualContractRequest) (*CaptureAnnualContractResult, error)
	RecognizeDueRevenue(ctx context.Context, req RecognizeDueRevenueRequest) (*RecognizeDueRevenueResult, error)
}

// DojoAccountingService exposes the dojo-specific accounting operations used by the capture API.
type DojoAccountingService interface {
	AssessMonthlyTuition(ctx context.Context, req AssessMonthlyTuitionRequest) (*AssessMonthlyTuitionResult, error)
	CapturePayment(ctx context.Context, req CaptureDojoPaymentRequest) (*CaptureDojoPaymentResult, error)
}

// RentalAccountingService exposes the rental-specific accounting operations used by the capture API.
type RentalAccountingService interface {
	CaptureDeposit(ctx context.Context, req CaptureRentalDepositRequest) (*CaptureRentalDepositResult, error)
	ReleaseDeposit(ctx context.Context, req ReleaseRentalDepositRequest) (*ReleaseRentalDepositResult, error)
}

// CaptureAnnualContractRequest describes a prepaid SaaS contract that should be deferred then recognized monthly.
type CaptureAnnualContractRequest struct {
	CompanyCode                string             `json:"company_code"`
	ContractNo                 string             `json:"contract_no"`
	CustomerRef                string             `json:"customer_ref"`
	CashAccountID              string             `json:"cash_account_id"`
	DeferredRevenueAccountID   string             `json:"deferred_revenue_account_id"`
	RecognizedRevenueAccountID string             `json:"recognized_revenue_account_id"`
	CurrencyCode               string             `json:"currency_code"`
	ServiceStartDate           time.Time          `json:"service_start_date"`
	TermMonths                 int                `json:"term_months"`
	TotalAmount                ledgerdomain.Money `json:"total_amount"`
	SourceRef                  string             `json:"source_ref,omitempty"`
}

// CaptureAnnualContractResult reports the created contract and the seed ledger entry.
type CaptureAnnualContractResult struct {
	ContractID            string   `json:"contract_id"`
	InitialJournalEntryID string   `json:"initial_journal_entry_id"`
	ScheduleCount         int      `json:"schedule_count"`
	MonthlyAmounts        []string `json:"monthly_amounts"`
}

// RecognizeDueRevenueRequest triggers recognition of all SaaS schedules due up to the provided month.
type RecognizeDueRevenueRequest struct {
	CompanyCode string    `json:"company_code"`
	UpTo        time.Time `json:"up_to"`
	Limit       int       `json:"limit,omitempty"`
}

// RecognizeDueRevenueResult reports how many schedules were recognized.
type RecognizeDueRevenueResult struct {
	RecognizedCount int      `json:"recognized_count"`
	JournalEntryIDs []string `json:"journal_entry_ids"`
}

// AssessMonthlyTuitionRequest creates a dojo receivable for a student.
type AssessMonthlyTuitionRequest struct {
	CompanyCode         string             `json:"company_code"`
	StudentRef          string             `json:"student_ref"`
	BillingMonth        time.Time          `json:"billing_month"`
	DueDate             time.Time          `json:"due_date"`
	ReceivableAccountID string             `json:"receivable_account_id"`
	RevenueAccountID    string             `json:"revenue_account_id"`
	CurrencyCode        string             `json:"currency_code"`
	Amount              ledgerdomain.Money `json:"amount"`
	SourceRef           string             `json:"source_ref,omitempty"`
}

// AssessMonthlyTuitionResult reports the created dojo receivable.
type AssessMonthlyTuitionResult struct {
	ReceivableID     string `json:"receivable_id"`
	JournalEntryID   string `json:"journal_entry_id"`
	ReceivableStatus string `json:"receivable_status"`
}

// CaptureDojoPaymentRequest records payment for a dojo receivable.
type CaptureDojoPaymentRequest struct {
	CompanyCode   string             `json:"company_code"`
	StudentRef    string             `json:"student_ref"`
	BillingMonth  time.Time          `json:"billing_month"`
	CashAccountID string             `json:"cash_account_id"`
	CurrencyCode  string             `json:"currency_code"`
	PaymentAmount ledgerdomain.Money `json:"payment_amount"`
	SourceRef     string             `json:"source_ref,omitempty"`
}

// CaptureDojoPaymentResult reports the receivable settlement result.
type CaptureDojoPaymentResult struct {
	ReceivableID   string `json:"receivable_id"`
	JournalEntryID string `json:"journal_entry_id"`
	Status         string `json:"status"`
}

// CaptureRentalDepositRequest holds the data needed to lock a rental deposit.
type CaptureRentalDepositRequest struct {
	CompanyCode      string             `json:"company_code"`
	RentalOrderID    string             `json:"rental_order_id"`
	CustomerRef      string             `json:"customer_ref"`
	CashAccountID    string             `json:"cash_account_id"`
	HoldingAccountID string             `json:"holding_account_id"`
	CurrencyCode     string             `json:"currency_code"`
	Amount           ledgerdomain.Money `json:"amount"`
	SourceRef        string             `json:"source_ref,omitempty"`
}

// CaptureRentalDepositResult reports the deposit hold state.
type CaptureRentalDepositResult struct {
	DepositID      string `json:"deposit_id"`
	JournalEntryID string `json:"journal_entry_id"`
	DepositStatus  string `json:"deposit_status"`
}

// ReleaseRentalDepositRequest releases a held rental deposit.
type ReleaseRentalDepositRequest struct {
	CompanyCode   string `json:"company_code"`
	RentalOrderID string `json:"rental_order_id"`
	SourceRef     string `json:"source_ref,omitempty"`
}

// ReleaseRentalDepositResult reports the deposit release state.
type ReleaseRentalDepositResult struct {
	DepositID      string `json:"deposit_id"`
	JournalEntryID string `json:"journal_entry_id"`
	DepositStatus  string `json:"deposit_status"`
}

var _ financedomain.BusinessLine
