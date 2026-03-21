package domain

import (
	"context"
	"time"

	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

// RevenueStreamPoint represents one business-line slice in the analytics dashboard.
type RevenueStreamPoint struct {
	CostCenter string             `json:"cost_center"`
	Amount     ledgerdomain.Money `json:"amount"`
}

// FinanceSummary returns the headline KPI cards for the executive dashboard.
type FinanceSummary struct {
	CompanyCode       string             `json:"company_code"`
	TotalRevenue      ledgerdomain.Money `json:"total_revenue"`
	GrossProfit       ledgerdomain.Money `json:"gross_profit"`
	GrossProfitMargin float64            `json:"gross_profit_margin"`
	NetCash           ledgerdomain.Money `json:"net_cash"`
	CurrencyCode      string             `json:"currency_code"`
}

// SegmentGrossProfit captures one operating segment slice for pie and bar charts.
type SegmentGrossProfit struct {
	CompanyCode       string             `json:"company_code"`
	SegmentKey        string             `json:"segment_key"`
	SegmentLabel      string             `json:"segment_label"`
	GrossRevenue      ledgerdomain.Money `json:"gross_revenue"`
	RevenueDeductions ledgerdomain.Money `json:"revenue_deductions"`
	NetRevenue        ledgerdomain.Money `json:"net_revenue"`
	OtherIncome       ledgerdomain.Money `json:"other_income"`
	CostOfGoodsSold   ledgerdomain.Money `json:"cost_of_goods_sold"`
	GrossProfit       ledgerdomain.Money `json:"gross_profit"`
	GrossMarginRatio  float64            `json:"gross_margin_ratio"`
	GrossProfitShare  float64            `json:"gross_profit_share"`
}

// CashRunwayMonth models one month in the projected runway horizon.
type CashRunwayMonth struct {
	MonthLabel       string             `json:"month_label"`
	OpeningCash      ledgerdomain.Money `json:"opening_cash"`
	ContractedInflow ledgerdomain.Money `json:"contracted_inflow"`
	ProjectedBurn    ledgerdomain.Money `json:"projected_burn"`
	ProjectedEnding  ledgerdomain.Money `json:"projected_ending_cash"`
}

// CashRunway aggregates the runway summary delivered to dashboards.
type CashRunway struct {
	AsOf               time.Time          `json:"as_of"`
	CurrentCash        ledgerdomain.Money `json:"current_cash"`
	AverageMonthlyBurn ledgerdomain.Money `json:"average_monthly_burn"`
	Months             []CashRunwayMonth  `json:"months"`
}

// ReconciliationBalance captures the live effective balance for one ledger account.
type ReconciliationBalance struct {
	CompanyCode            string
	AccountCode            string
	AccountName            string
	EffectiveDebitBalance  ledgerdomain.Money
	EffectiveCreditBalance ledgerdomain.Money
	EffectiveNetBalance    ledgerdomain.Money
	LastPostingDate        time.Time
}

// FinanceSeriesPoint captures one monthly point for dashboard charts.
type FinanceSeriesPoint struct {
	PeriodStart time.Time
	PeriodLabel string
	CashEnding  ledgerdomain.Money
	Revenue     ledgerdomain.Money
	Expense     ledgerdomain.Money
	Profit      ledgerdomain.Money
}

// DashboardEvent is the realtime notification delivered to executive clients.
type DashboardEvent struct {
	Event        string    `json:"event"`
	CompanyCode  string    `json:"company_code"`
	EntryID      string    `json:"entry_id"`
	ReferenceNo  string    `json:"reference_no,omitempty"`
	Amount       float64   `json:"amount"`
	Segment      string    `json:"segment"`
	SourceModule string    `json:"source_module,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// ReportAccessLog preserves who opened which finance dashboard report for auditability.
type ReportAccessLog struct {
	CompanyCode string
	ReportCode  string
	ActorID     string
	ActorRole   string
	IPAddress   string
	UserAgent   string
	Filters     map[string]string
	AccessedAt  time.Time
}

// Repository exposes the analytics read model.
type Repository interface {
	FinanceSummary(ctx context.Context, companyCode string) (FinanceSummary, error)
	Segments(ctx context.Context, companyCode string) ([]SegmentGrossProfit, error)
	RevenueStream(ctx context.Context, companyCode string, from time.Time, to time.Time) ([]RevenueStreamPoint, error)
	CashRunway(ctx context.Context, companyCode string, asOf time.Time, months int) (CashRunway, error)
	ReconciliationBalances(ctx context.Context, companyCode string) ([]ReconciliationBalance, error)
	FinanceSeries(ctx context.Context, companyCode string, asOf time.Time, months int) ([]FinanceSeriesPoint, error)
	RecordReportAccess(ctx context.Context, access ReportAccessLog) error
}
