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

// RetailAccountingService exposes the retail POS accounting operations used by the capture API.
type RetailAccountingService interface {
	CaptureSale(ctx context.Context, req CaptureRetailSaleRequest) (*CaptureRetailSaleResult, error)
	CaptureRefund(ctx context.Context, req CaptureRetailRefundRequest) (*CaptureRetailRefundResult, error)
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
	CapturedAt                 time.Time          `json:"captured_at,omitempty"`
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
	PaidAt        time.Time          `json:"paid_at,omitempty"`
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

// CaptureRetailSaleRequest records a POS sale.
type CaptureRetailSaleRequest struct {
	CompanyCode      string             `json:"company_code"`
	OrderNo          string             `json:"order_no"`
	CustomerRef      string             `json:"customer_ref,omitempty"`
	OccurredAt       time.Time          `json:"occurred_at,omitempty"`
	CashAccountID    string             `json:"cash_account_id"`
	RevenueAccountID string             `json:"revenue_account_id"`
	CurrencyCode     string             `json:"currency_code"`
	Amount           ledgerdomain.Money `json:"amount"`
	SourceRef        string             `json:"source_ref,omitempty"`
}

// CaptureRetailSaleResult reports the posted retail sale.
type CaptureRetailSaleResult struct {
	OrderNo        string `json:"order_no"`
	JournalEntryID string `json:"journal_entry_id"`
}

// CaptureRetailRefundRequest records a POS refund as contra revenue.
type CaptureRetailRefundRequest struct {
	CompanyCode            string             `json:"company_code"`
	OrderNo                string             `json:"order_no"`
	CustomerRef            string             `json:"customer_ref,omitempty"`
	RefundedAt             time.Time          `json:"refunded_at,omitempty"`
	CashAccountID          string             `json:"cash_account_id"`
	ContraRevenueAccountID string             `json:"contra_revenue_account_id"`
	CurrencyCode           string             `json:"currency_code"`
	Amount                 ledgerdomain.Money `json:"amount"`
	SourceRef              string             `json:"source_ref,omitempty"`
}

// CaptureRetailRefundResult reports the posted retail refund.
type CaptureRetailRefundResult struct {
	OrderNo        string `json:"order_no"`
	JournalEntryID string `json:"journal_entry_id"`
}

// CaptureRentalDepositRequest holds the data needed to lock a rental deposit.
type CaptureRentalDepositRequest struct {
	CompanyCode      string             `json:"company_code"`
	RentalOrderID    string             `json:"rental_order_id"`
	CustomerRef      string             `json:"customer_ref"`
	HeldAt           time.Time          `json:"held_at,omitempty"`
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
	CompanyCode            string             `json:"company_code"`
	RentalOrderID          string             `json:"rental_order_id"`
	ReleasedAt             time.Time          `json:"released_at,omitempty"`
	CashAccountID          string             `json:"cash_account_id,omitempty"`
	DamageAmount           ledgerdomain.Money `json:"damage_amount,omitempty"`
	DamageRevenueAccountID string             `json:"damage_revenue_account_id,omitempty"`
	SourceRef              string             `json:"source_ref,omitempty"`
}

// ReleaseRentalDepositResult reports the deposit release state.
type ReleaseRentalDepositResult struct {
	DepositID      string `json:"deposit_id"`
	JournalEntryID string `json:"journal_entry_id"`
	DepositStatus  string `json:"deposit_status"`
}

var _ financedomain.BusinessLine
