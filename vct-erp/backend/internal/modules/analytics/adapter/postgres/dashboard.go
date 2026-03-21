package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

// FinanceSummary reads the executive KPI card values from the dashboard views.
func (r *Repository) FinanceSummary(ctx context.Context, companyCode string) (analyticsdomain.FinanceSummary, error) {
	var (
		totalRevenueRaw string
		grossProfitRaw  string
		netCashRaw      string
		margin          float64
	)

	err := r.db.QueryRowContext(ctx, `
SELECT
    COALESCE(SUM(segment.net_revenue), 0)::NUMERIC(20, 4) AS total_revenue,
    COALESCE(SUM(segment.gross_profit), 0)::NUMERIC(20, 4) AS gross_profit,
    COALESCE(
        CASE
            WHEN SUM(segment.net_revenue) = 0 THEN 0
            ELSE ROUND((SUM(segment.gross_profit) / NULLIF(SUM(segment.net_revenue), 0)) * 100, 4)
        END,
        0
    ) AS gross_profit_margin,
    COALESCE(MAX(runway.current_cash), 0)::NUMERIC(20, 4) AS net_cash
FROM v_gross_profit_by_segment AS segment
LEFT JOIN v_cash_runway_metrics AS runway
    ON runway.company_code = segment.company_code
WHERE segment.company_code = $1`,
		companyCode,
	).Scan(&totalRevenueRaw, &grossProfitRaw, &margin, &netCashRaw)
	if err != nil {
		return analyticsdomain.FinanceSummary{}, fmt.Errorf("query finance summary: %w", err)
	}

	totalRevenue, err := ledgerdomain.ParseMoney(totalRevenueRaw)
	if err != nil {
		return analyticsdomain.FinanceSummary{}, fmt.Errorf("parse total revenue: %w", err)
	}
	grossProfit, err := ledgerdomain.ParseMoney(grossProfitRaw)
	if err != nil {
		return analyticsdomain.FinanceSummary{}, fmt.Errorf("parse gross profit: %w", err)
	}
	netCash, err := ledgerdomain.ParseMoney(netCashRaw)
	if err != nil {
		return analyticsdomain.FinanceSummary{}, fmt.Errorf("parse net cash: %w", err)
	}

	return analyticsdomain.FinanceSummary{
		CompanyCode:       companyCode,
		TotalRevenue:      totalRevenue,
		GrossProfit:       grossProfit,
		GrossProfitMargin: margin,
		NetCash:           netCash,
		CurrencyCode:      "VND",
	}, nil
}

// Segments reads the gross profit split from the dashboard view and computes each slice share.
func (r *Repository) Segments(ctx context.Context, companyCode string) ([]analyticsdomain.SegmentGrossProfit, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT
    company_code,
    segment_key,
    segment_label,
    gross_revenue,
    revenue_deductions,
    net_revenue,
    other_income,
    cost_of_goods_sold,
    gross_profit,
    gross_margin_ratio,
    COALESCE(
        CASE
            WHEN SUM(gross_profit) OVER () = 0 THEN 0
            ELSE ROUND((gross_profit / NULLIF(SUM(gross_profit) OVER (), 0)) * 100, 4)
        END,
        0
    ) AS gross_profit_share
FROM v_gross_profit_by_segment
WHERE company_code = $1
ORDER BY segment_order`,
		companyCode,
	)
	if err != nil {
		return nil, fmt.Errorf("query finance segments: %w", err)
	}
	defer rows.Close()

	segments := make([]analyticsdomain.SegmentGrossProfit, 0, 4)
	for rows.Next() {
		var (
			item                           analyticsdomain.SegmentGrossProfit
			grossRevenueRaw, deductionsRaw string
			netRevenueRaw, otherIncomeRaw  string
			costRaw, grossProfitRaw        string
		)
		if err := rows.Scan(
			&item.CompanyCode,
			&item.SegmentKey,
			&item.SegmentLabel,
			&grossRevenueRaw,
			&deductionsRaw,
			&netRevenueRaw,
			&otherIncomeRaw,
			&costRaw,
			&grossProfitRaw,
			&item.GrossMarginRatio,
			&item.GrossProfitShare,
		); err != nil {
			return nil, fmt.Errorf("scan finance segment: %w", err)
		}

		if item.GrossRevenue, err = ledgerdomain.ParseMoney(grossRevenueRaw); err != nil {
			return nil, fmt.Errorf("parse segment gross revenue: %w", err)
		}
		if item.RevenueDeductions, err = ledgerdomain.ParseMoney(deductionsRaw); err != nil {
			return nil, fmt.Errorf("parse segment deductions: %w", err)
		}
		if item.NetRevenue, err = ledgerdomain.ParseMoney(netRevenueRaw); err != nil {
			return nil, fmt.Errorf("parse segment net revenue: %w", err)
		}
		if item.OtherIncome, err = ledgerdomain.ParseMoney(otherIncomeRaw); err != nil {
			return nil, fmt.Errorf("parse segment other income: %w", err)
		}
		if item.CostOfGoodsSold, err = ledgerdomain.ParseMoney(costRaw); err != nil {
			return nil, fmt.Errorf("parse segment cogs: %w", err)
		}
		if item.GrossProfit, err = ledgerdomain.ParseMoney(grossProfitRaw); err != nil {
			return nil, fmt.Errorf("parse segment gross profit: %w", err)
		}

		segments = append(segments, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate finance segments: %w", err)
	}
	return segments, nil
}

// RecordReportAccess writes a mandatory audit trail for executive report access.
func (r *Repository) RecordReportAccess(ctx context.Context, access analyticsdomain.ReportAccessLog) error {
	filtersRaw, err := json.Marshal(access.Filters)
	if err != nil {
		return fmt.Errorf("marshal finance report audit filters: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
INSERT INTO finance_report_access_logs (
    company_code,
    report_code,
    actor_id,
    actor_role,
    ip_address,
    user_agent,
    filters,
    accessed_at
)
VALUES ($1, $2, $3, $4, NULLIF($5, '')::INET, NULLIF($6, ''), CAST($7 AS JSONB), $8)`,
		access.CompanyCode,
		access.ReportCode,
		access.ActorID,
		access.ActorRole,
		access.IPAddress,
		access.UserAgent,
		string(filtersRaw),
		access.AccessedAt,
	)
	if err != nil {
		return fmt.Errorf("insert finance report audit log: %w", err)
	}
	return nil
}
