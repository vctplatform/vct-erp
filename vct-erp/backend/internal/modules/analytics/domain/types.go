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

// Repository exposes the analytics read model.
type Repository interface {
	RevenueStream(ctx context.Context, companyCode string, from time.Time, to time.Time) ([]RevenueStreamPoint, error)
	CashRunway(ctx context.Context, companyCode string, asOf time.Time, months int) (CashRunway, error)
}
