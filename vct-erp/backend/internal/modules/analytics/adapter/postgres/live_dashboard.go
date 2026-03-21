package postgres

import (
	"context"
	"fmt"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

// ReconciliationBalances reads the effective account balances from the live dashboard view.
func (r *Repository) ReconciliationBalances(ctx context.Context, companyCode string) ([]analyticsdomain.ReconciliationBalance, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT
    company_code,
    account_code,
    account_name,
    effective_debit_balance,
    effective_credit_balance,
    effective_net_balance,
    COALESCE(last_posting_date, CURRENT_DATE)::date AS last_posting_date
FROM v_finance_reconciliation
WHERE company_code = $1
ORDER BY account_code`,
		companyCode,
	)
	if err != nil {
		return nil, fmt.Errorf("query finance reconciliation: %w", err)
	}
	defer rows.Close()

	balances := make([]analyticsdomain.ReconciliationBalance, 0, 64)
	for rows.Next() {
		var (
			item            analyticsdomain.ReconciliationBalance
			debitRaw        string
			creditRaw       string
			netRaw          string
			lastPostingDate time.Time
		)

		if err := rows.Scan(
			&item.CompanyCode,
			&item.AccountCode,
			&item.AccountName,
			&debitRaw,
			&creditRaw,
			&netRaw,
			&lastPostingDate,
		); err != nil {
			return nil, fmt.Errorf("scan finance reconciliation row: %w", err)
		}

		var err error
		if item.EffectiveDebitBalance, err = ledgerdomain.ParseMoney(debitRaw); err != nil {
			return nil, fmt.Errorf("parse reconciliation debit balance: %w", err)
		}
		if item.EffectiveCreditBalance, err = ledgerdomain.ParseMoney(creditRaw); err != nil {
			return nil, fmt.Errorf("parse reconciliation credit balance: %w", err)
		}
		if item.EffectiveNetBalance, err = ledgerdomain.ParseMoney(netRaw); err != nil {
			return nil, fmt.Errorf("parse reconciliation net balance: %w", err)
		}
		item.LastPostingDate = lastPostingDate.UTC()
		balances = append(balances, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate finance reconciliation rows: %w", err)
	}

	return balances, nil
}

// FinanceSeries returns monthly revenue, expense, profit, and ending cash points for dashboard charts.
func (r *Repository) FinanceSeries(ctx context.Context, companyCode string, asOf time.Time, months int) ([]analyticsdomain.FinanceSeriesPoint, error) {
	if months <= 0 {
		months = 6
	}

	rows, err := r.db.QueryContext(ctx, `
WITH bounds AS (
    SELECT
        date_trunc('month', $2::date)::date AS anchor_month,
        (date_trunc('month', $2::date)::date - make_interval(months => $3 - 1))::date AS window_start
),
month_series AS (
    SELECT generate_series(bounds.window_start, bounds.anchor_month, INTERVAL '1 month')::date AS month_start
    FROM bounds
),
opening_cash AS (
    SELECT
        COALESCE(SUM(ji.amount_signed), 0)::NUMERIC(20, 4) AS opening_cash
    FROM journal_entries AS je
    INNER JOIN journal_items AS ji ON ji.journal_entry_id = je.id
    INNER JOIN accounts AS a ON a.id = ji.account_id
    CROSS JOIN bounds
    WHERE je.company_code = $1
      AND je.status IN ('posted', 'reversed')
      AND (a.code LIKE '111%' OR a.code LIKE '112%')
      AND je.posting_date < bounds.window_start
),
monthly_aggregates AS (
    SELECT
        date_trunc('month', je.posting_date)::date AS month_start,
        COALESCE(SUM(CASE
            WHEN a.code LIKE '111%' OR a.code LIKE '112%' THEN ji.amount_signed
            ELSE 0
        END), 0)::NUMERIC(20, 4) AS cash_delta,
        COALESCE(SUM(CASE
            WHEN a.code LIKE '511%' THEN CASE WHEN ji.side = 'credit' THEN ji.amount ELSE ji.amount * -1 END
            WHEN a.code LIKE '521%' THEN CASE WHEN ji.side = 'debit' THEN ji.amount * -1 ELSE ji.amount END
            WHEN a.code LIKE '711%' THEN CASE WHEN ji.side = 'credit' THEN ji.amount ELSE ji.amount * -1 END
            ELSE 0
        END), 0)::NUMERIC(20, 4) AS revenue_amount,
        COALESCE(SUM(CASE
            WHEN a.code LIKE '632%' OR
                 a.code LIKE '635%' OR
                 a.code LIKE '641%' OR
                 a.code LIKE '642%' OR
                 a.code LIKE '811%' OR
                 a.code LIKE '821%'
                THEN CASE WHEN ji.side = 'debit' THEN ji.amount ELSE ji.amount * -1 END
            ELSE 0
        END), 0)::NUMERIC(20, 4) AS expense_amount
    FROM journal_entries AS je
    INNER JOIN journal_items AS ji ON ji.journal_entry_id = je.id
    INNER JOIN accounts AS a ON a.id = ji.account_id
    CROSS JOIN bounds
    WHERE je.company_code = $1
      AND je.status IN ('posted', 'reversed')
      AND je.posting_date >= bounds.window_start
      AND je.posting_date < (bounds.anchor_month + INTERVAL '1 month')
    GROUP BY date_trunc('month', je.posting_date)::date
),
series AS (
    SELECT
        month_series.month_start,
        COALESCE(monthly_aggregates.revenue_amount, 0)::NUMERIC(20, 4) AS revenue_amount,
        COALESCE(monthly_aggregates.expense_amount, 0)::NUMERIC(20, 4) AS expense_amount,
        (
            COALESCE((SELECT opening_cash FROM opening_cash), 0)
            + SUM(COALESCE(monthly_aggregates.cash_delta, 0)) OVER (
                ORDER BY month_series.month_start
                ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
            )
        )::NUMERIC(20, 4) AS ending_cash
    FROM month_series
    LEFT JOIN monthly_aggregates
        ON monthly_aggregates.month_start = month_series.month_start
)
SELECT
    month_start,
    revenue_amount,
    expense_amount,
    (revenue_amount - expense_amount)::NUMERIC(20, 4) AS profit_amount,
    ending_cash
FROM series
ORDER BY month_start`,
		companyCode,
		asOf.Format("2006-01-02"),
		months,
	)
	if err != nil {
		return nil, fmt.Errorf("query finance series: %w", err)
	}
	defer rows.Close()

	points := make([]analyticsdomain.FinanceSeriesPoint, 0, months)
	for rows.Next() {
		var (
			point         analyticsdomain.FinanceSeriesPoint
			revenueRaw    string
			expenseRaw    string
			profitRaw     string
			cashEndingRaw string
		)

		if err := rows.Scan(
			&point.PeriodStart,
			&revenueRaw,
			&expenseRaw,
			&profitRaw,
			&cashEndingRaw,
		); err != nil {
			return nil, fmt.Errorf("scan finance series row: %w", err)
		}

		var err error
		if point.Revenue, err = ledgerdomain.ParseMoney(revenueRaw); err != nil {
			return nil, fmt.Errorf("parse finance series revenue: %w", err)
		}
		if point.Expense, err = ledgerdomain.ParseMoney(expenseRaw); err != nil {
			return nil, fmt.Errorf("parse finance series expense: %w", err)
		}
		if point.Profit, err = ledgerdomain.ParseMoney(profitRaw); err != nil {
			return nil, fmt.Errorf("parse finance series profit: %w", err)
		}
		if point.CashEnding, err = ledgerdomain.ParseMoney(cashEndingRaw); err != nil {
			return nil, fmt.Errorf("parse finance series ending cash: %w", err)
		}

		point.PeriodStart = point.PeriodStart.UTC()
		point.PeriodLabel = point.PeriodStart.Format("2006-01")
		points = append(points, point)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate finance series rows: %w", err)
	}

	return points, nil
}
