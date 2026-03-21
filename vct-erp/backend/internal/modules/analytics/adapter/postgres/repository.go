package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

// Repository implements dashboard reads on PostgreSQL.
type Repository struct {
	db *sql.DB
}

// NewRepository constructs the Postgres analytics adapter.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// RevenueStream returns net revenue by cost center for the selected period.
func (r *Repository) RevenueStream(ctx context.Context, companyCode string, from time.Time, to time.Time) ([]analyticsdomain.RevenueStreamPoint, error) {
	rows, err := r.db.QueryContext(ctx, `
WITH revenue_source AS (
    SELECT
        COALESCE(je.metadata->>'cost_center', je.metadata->>'business_line', 'unassigned') AS cost_center,
        a.code AS account_code,
        CASE
            WHEN a.code LIKE '521%' THEN ji.amount
            WHEN ji.side = 'credit' THEN ji.amount
            ELSE ji.amount * -1
        END AS signed_amount
    FROM journal_entries AS je
    INNER JOIN journal_items AS ji ON ji.journal_entry_id = je.id
    INNER JOIN accounts AS a ON a.id = ji.account_id
    WHERE je.company_code = $1
      AND je.status IN ('posted', 'reversed')
      AND je.posting_date BETWEEN $2 AND $3
      AND (
          a.code LIKE '511%' OR
          a.code LIKE '521%' OR
          a.code LIKE '515%' OR
          a.code LIKE '711%'
      )
),
cost_center_window AS (
    SELECT DISTINCT
        cost_center,
        COALESCE(SUM(CASE WHEN account_code LIKE '511%' THEN signed_amount ELSE 0 END) OVER (PARTITION BY cost_center), 0) AS gross_revenue,
        COALESCE(SUM(CASE WHEN account_code LIKE '521%' THEN signed_amount ELSE 0 END) OVER (PARTITION BY cost_center), 0) AS revenue_deductions,
        COALESCE(SUM(CASE WHEN account_code LIKE '515%' THEN signed_amount ELSE 0 END) OVER (PARTITION BY cost_center), 0) AS financial_income,
        COALESCE(SUM(CASE WHEN account_code LIKE '711%' THEN signed_amount ELSE 0 END) OVER (PARTITION BY cost_center), 0) AS other_income
    FROM revenue_source
)
SELECT
    cost_center,
    gross_revenue - revenue_deductions + financial_income + other_income AS net_revenue
FROM cost_center_window
ORDER BY cost_center`,
		companyCode,
		from.Format("2006-01-02"),
		to.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("query revenue stream: %w", err)
	}
	defer rows.Close()

	points := make([]analyticsdomain.RevenueStreamPoint, 0, 8)
	for rows.Next() {
		var (
			point     analyticsdomain.RevenueStreamPoint
			amountRaw string
		)
		if err := rows.Scan(&point.CostCenter, &amountRaw); err != nil {
			return nil, fmt.Errorf("scan revenue stream row: %w", err)
		}
		amount, err := ledgerdomain.ParseMoney(amountRaw)
		if err != nil {
			return nil, fmt.Errorf("parse revenue stream amount: %w", err)
		}
		point.Amount = amount
		points = append(points, point)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate revenue stream rows: %w", err)
	}
	return points, nil
}

// CashRunway returns a 3-month projected runway using current cash, recent burn, and contracted SaaS schedules.
func (r *Repository) CashRunway(ctx context.Context, companyCode string, asOf time.Time, months int) (analyticsdomain.CashRunway, error) {
	var currentCashRaw string
	if err := r.db.QueryRowContext(ctx, `
SELECT current_cash
FROM v_cash_runway_metrics
WHERE company_code = $1`,
		companyCode,
	).Scan(&currentCashRaw); err != nil {
		return analyticsdomain.CashRunway{}, fmt.Errorf("query current cash: %w", err)
	}
	currentCash, err := ledgerdomain.ParseMoney(currentCashRaw)
	if err != nil {
		return analyticsdomain.CashRunway{}, fmt.Errorf("parse current cash: %w", err)
	}

	startMonth := time.Date(asOf.Year(), asOf.Month(), 1, 0, 0, 0, 0, time.UTC)
	var burnRaw string
	if err := r.db.QueryRowContext(ctx, `
SELECT average_monthly_burn
FROM v_cash_runway_metrics
WHERE company_code = $1`,
		companyCode,
	).Scan(&burnRaw); err != nil {
		return analyticsdomain.CashRunway{}, fmt.Errorf("query average burn: %w", err)
	}

	averageBurn := ledgerdomain.MustParseMoney("0.0000")
	averageBurn, err = ledgerdomain.ParseMoney(burnRaw)
	if err != nil {
		return analyticsdomain.CashRunway{}, fmt.Errorf("parse average burn: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT
    date_trunc('month', schedule.service_month)::date AS month_start,
    COALESCE(SUM(schedule.amount), 0) AS contracted_amount
FROM saas_revenue_schedules AS schedule
INNER JOIN saas_contracts AS contract ON contract.id = schedule.contract_id
WHERE contract.company_code = $1
  AND schedule.status = 'scheduled'
  AND schedule.service_month >= $2::date
  AND schedule.service_month < ($2::date + make_interval(months => $3))
GROUP BY date_trunc('month', schedule.service_month)::date
ORDER BY month_start`,
		companyCode,
		startMonth.Format("2006-01-02"),
		months,
	)
	if err != nil {
		return analyticsdomain.CashRunway{}, fmt.Errorf("query contracted inflow: %w", err)
	}
	defer rows.Close()

	contractedByMonth := make(map[string]ledgerdomain.Money, months)
	for rows.Next() {
		var (
			monthStart time.Time
			amountRaw  string
		)
		if err := rows.Scan(&monthStart, &amountRaw); err != nil {
			return analyticsdomain.CashRunway{}, fmt.Errorf("scan contracted inflow row: %w", err)
		}
		amount, err := ledgerdomain.ParseMoney(amountRaw)
		if err != nil {
			return analyticsdomain.CashRunway{}, fmt.Errorf("parse contracted inflow amount: %w", err)
		}
		contractedByMonth[monthStart.Format("2006-01")] = amount
	}

	projection := analyticsdomain.CashRunway{
		AsOf:               asOf,
		CurrentCash:        currentCash,
		AverageMonthlyBurn: averageBurn,
		Months:             make([]analyticsdomain.CashRunwayMonth, 0, months),
	}

	opening := currentCash
	for index := 0; index < months; index++ {
		monthStart := startMonth.AddDate(0, index, 0)
		contracted := contractedByMonth[monthStart.Format("2006-01")]
		projectedEnding := opening.Add(contracted).Sub(averageBurn)
		projection.Months = append(projection.Months, analyticsdomain.CashRunwayMonth{
			MonthLabel:       monthStart.Format("2006-01"),
			OpeningCash:      opening,
			ContractedInflow: contracted,
			ProjectedBurn:    averageBurn,
			ProjectedEnding:  projectedEnding,
		})
		opening = projectedEnding
	}

	return projection, nil
}
