BEGIN;

CREATE TABLE IF NOT EXISTS finance_report_access_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_code    VARCHAR(32) NOT NULL,
    report_code     VARCHAR(64) NOT NULL,
    actor_id        VARCHAR(128) NOT NULL,
    actor_role      VARCHAR(64) NOT NULL,
    ip_address      INET,
    user_agent      TEXT,
    filters         JSONB NOT NULL DEFAULT '{}'::JSONB,
    accessed_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_finance_report_access_logs_report_time
    ON finance_report_access_logs (company_code, report_code, accessed_at DESC);

CREATE INDEX IF NOT EXISTS idx_finance_report_access_logs_actor_time
    ON finance_report_access_logs (actor_id, accessed_at DESC);

CREATE OR REPLACE VIEW v_finance_reconciliation AS
WITH account_catalog AS (
    SELECT
        a.company_code,
        a.id AS account_id,
        a.code AS account_code,
        a.name AS account_name,
        a.account_type,
        a.normal_side,
        'VND'::CHAR(3) AS default_currency_code
    FROM accounts AS a
    WHERE a.is_active
),
movement_totals AS (
    SELECT
        je.company_code,
        ji.account_id,
        ji.currency_code,
        COALESCE(SUM(CASE WHEN je.status = 'posted' AND ji.side = 'debit' THEN ji.amount ELSE 0 END), 0)::NUMERIC(20, 4) AS posted_debit,
        COALESCE(SUM(CASE WHEN je.status = 'posted' AND ji.side = 'credit' THEN ji.amount ELSE 0 END), 0)::NUMERIC(20, 4) AS posted_credit,
        COALESCE(SUM(CASE WHEN je.status = 'reversed' AND ji.side = 'debit' THEN ji.amount ELSE 0 END), 0)::NUMERIC(20, 4) AS reversed_debit,
        COALESCE(SUM(CASE WHEN je.status = 'reversed' AND ji.side = 'credit' THEN ji.amount ELSE 0 END), 0)::NUMERIC(20, 4) AS reversed_credit,
        COALESCE(SUM(CASE WHEN je.status = 'posted' THEN ji.amount_signed ELSE 0 END), 0)::NUMERIC(20, 4) AS posted_net_balance,
        COALESCE(SUM(CASE WHEN je.status = 'reversed' THEN ji.amount_signed ELSE 0 END), 0)::NUMERIC(20, 4) AS reversed_net_balance,
        MAX(je.posting_date) FILTER (WHERE je.status IN ('posted', 'reversed')) AS last_posting_date,
        MAX(ji.updated_at) AS last_movement_at
    FROM journal_entries AS je
    INNER JOIN journal_items AS ji ON ji.journal_entry_id = je.id
    WHERE je.status IN ('posted', 'reversed')
    GROUP BY je.company_code, ji.account_id, ji.currency_code
)
SELECT
    catalog.company_code,
    catalog.account_id,
    catalog.account_code,
    catalog.account_name,
    catalog.account_type,
    catalog.normal_side,
    COALESCE(movement.currency_code, catalog.default_currency_code) AS currency_code,
    COALESCE(movement.posted_debit, 0)::NUMERIC(20, 4) AS posted_debit,
    COALESCE(movement.posted_credit, 0)::NUMERIC(20, 4) AS posted_credit,
    COALESCE(movement.reversed_debit, 0)::NUMERIC(20, 4) AS reversed_debit,
    COALESCE(movement.reversed_credit, 0)::NUMERIC(20, 4) AS reversed_credit,
    (COALESCE(movement.posted_debit, 0) + COALESCE(movement.reversed_debit, 0))::NUMERIC(20, 4) AS historical_debit,
    (COALESCE(movement.posted_credit, 0) + COALESCE(movement.reversed_credit, 0))::NUMERIC(20, 4) AS historical_credit,
    COALESCE(movement.posted_net_balance, 0)::NUMERIC(20, 4) AS posted_net_balance,
    COALESCE(movement.reversed_net_balance, 0)::NUMERIC(20, 4) AS reversed_net_balance,
    (COALESCE(movement.posted_net_balance, 0) + COALESCE(movement.reversed_net_balance, 0))::NUMERIC(20, 4) AS effective_net_balance,
    CASE
        WHEN COALESCE(movement.posted_net_balance, 0) + COALESCE(movement.reversed_net_balance, 0) >= 0
            THEN (COALESCE(movement.posted_net_balance, 0) + COALESCE(movement.reversed_net_balance, 0))::NUMERIC(20, 4)
        ELSE 0::NUMERIC(20, 4)
    END AS effective_debit_balance,
    CASE
        WHEN COALESCE(movement.posted_net_balance, 0) + COALESCE(movement.reversed_net_balance, 0) < 0
            THEN ABS(COALESCE(movement.posted_net_balance, 0) + COALESCE(movement.reversed_net_balance, 0))::NUMERIC(20, 4)
        ELSE 0::NUMERIC(20, 4)
    END AS effective_credit_balance,
    movement.last_posting_date,
    movement.last_movement_at
FROM account_catalog AS catalog
LEFT JOIN movement_totals AS movement
    ON movement.company_code = catalog.company_code
   AND movement.account_id = catalog.account_id;

CREATE OR REPLACE VIEW v_gross_profit_by_segment AS
WITH company_catalog AS (
    SELECT DISTINCT company_code
    FROM accounts
),
segment_catalog AS (
    SELECT *
    FROM (
        VALUES
            ('saas', 1, 'SaaS'),
            ('dojo', 2, 'Dojo'),
            ('retail', 3, 'Retail'),
            ('rental', 4, 'Rental')
    ) AS segment(segment_key, segment_order, segment_label)
),
segment_movements AS (
    SELECT
        je.company_code,
        COALESCE(je.metadata->>'cost_center', je.metadata->>'business_line', 'unassigned') AS segment_key,
        a.code AS account_code,
        ji.side,
        ji.amount
    FROM journal_entries AS je
    INNER JOIN journal_items AS ji ON ji.journal_entry_id = je.id
    INNER JOIN accounts AS a ON a.id = ji.account_id
    WHERE je.status IN ('posted', 'reversed')
      AND (
          a.code LIKE '511%' OR
          a.code LIKE '521%' OR
          a.code LIKE '632%' OR
          a.code LIKE '711%'
      )
),
aggregated AS (
    SELECT
        company_code,
        segment_key,
        COALESCE(SUM(CASE WHEN account_code LIKE '511%' THEN CASE WHEN side = 'credit' THEN amount ELSE amount * -1 END ELSE 0 END), 0)::NUMERIC(20, 4) AS gross_revenue,
        COALESCE(SUM(CASE WHEN account_code LIKE '521%' THEN CASE WHEN side = 'debit' THEN amount ELSE amount * -1 END ELSE 0 END), 0)::NUMERIC(20, 4) AS revenue_deductions,
        COALESCE(SUM(CASE WHEN account_code LIKE '711%' THEN CASE WHEN side = 'credit' THEN amount ELSE amount * -1 END ELSE 0 END), 0)::NUMERIC(20, 4) AS other_income,
        COALESCE(SUM(CASE WHEN account_code LIKE '632%' THEN CASE WHEN side = 'debit' THEN amount ELSE amount * -1 END ELSE 0 END), 0)::NUMERIC(20, 4) AS cost_of_goods_sold
    FROM segment_movements
    GROUP BY company_code, segment_key
)
SELECT
    company.company_code,
    segment.segment_key,
    segment.segment_label,
    segment.segment_order,
    COALESCE(aggregated.gross_revenue, 0)::NUMERIC(20, 4) AS gross_revenue,
    COALESCE(aggregated.revenue_deductions, 0)::NUMERIC(20, 4) AS revenue_deductions,
    (COALESCE(aggregated.gross_revenue, 0) - COALESCE(aggregated.revenue_deductions, 0) + COALESCE(aggregated.other_income, 0))::NUMERIC(20, 4) AS net_revenue,
    COALESCE(aggregated.other_income, 0)::NUMERIC(20, 4) AS other_income,
    COALESCE(aggregated.cost_of_goods_sold, 0)::NUMERIC(20, 4) AS cost_of_goods_sold,
    (
        COALESCE(aggregated.gross_revenue, 0)
        - COALESCE(aggregated.revenue_deductions, 0)
        + COALESCE(aggregated.other_income, 0)
        - COALESCE(aggregated.cost_of_goods_sold, 0)
    )::NUMERIC(20, 4) AS gross_profit,
    CASE
        WHEN (COALESCE(aggregated.gross_revenue, 0) - COALESCE(aggregated.revenue_deductions, 0) + COALESCE(aggregated.other_income, 0)) = 0
            THEN 0::NUMERIC(12, 6)
        ELSE ROUND((
            (
                COALESCE(aggregated.gross_revenue, 0)
                - COALESCE(aggregated.revenue_deductions, 0)
                + COALESCE(aggregated.other_income, 0)
                - COALESCE(aggregated.cost_of_goods_sold, 0)
            ) / NULLIF(
                COALESCE(aggregated.gross_revenue, 0)
                - COALESCE(aggregated.revenue_deductions, 0)
                + COALESCE(aggregated.other_income, 0),
                0
            )
        ), 6)::NUMERIC(12, 6)
    END AS gross_margin_ratio
FROM company_catalog AS company
CROSS JOIN segment_catalog AS segment
LEFT JOIN aggregated
    ON aggregated.company_code = company.company_code
   AND aggregated.segment_key = segment.segment_key;

CREATE OR REPLACE VIEW v_cash_runway_metrics AS
WITH company_catalog AS (
    SELECT DISTINCT company_code
    FROM accounts
),
as_of_anchor AS (
    SELECT CURRENT_DATE::DATE AS as_of_date
),
cash_position AS (
    SELECT
        je.company_code,
        COALESCE(SUM(ji.amount_signed), 0)::NUMERIC(20, 4) AS current_cash
    FROM journal_entries AS je
    INNER JOIN journal_items AS ji ON ji.journal_entry_id = je.id
    INNER JOIN accounts AS a ON a.id = ji.account_id
    WHERE je.status IN ('posted', 'reversed')
      AND (a.code LIKE '111%' OR a.code LIKE '112%')
    GROUP BY je.company_code
),
recent_monthly_burn AS (
    SELECT
        je.company_code,
        date_trunc('month', je.posting_date)::DATE AS month_start,
        COALESCE(SUM(CASE WHEN ji.side = 'debit' THEN ji.amount ELSE ji.amount * -1 END), 0)::NUMERIC(20, 4) AS burn_amount
    FROM journal_entries AS je
    INNER JOIN journal_items AS ji ON ji.journal_entry_id = je.id
    INNER JOIN accounts AS a ON a.id = ji.account_id
    CROSS JOIN as_of_anchor
    WHERE je.status IN ('posted', 'reversed')
      AND je.posting_date >= (date_trunc('month', as_of_anchor.as_of_date)::DATE - INTERVAL '3 months')
      AND je.posting_date < date_trunc('month', as_of_anchor.as_of_date)::DATE
      AND (
          a.code LIKE '632%' OR
          a.code LIKE '635%' OR
          a.code LIKE '641%' OR
          a.code LIKE '642%' OR
          a.code LIKE '811%' OR
          a.code LIKE '821%'
      )
    GROUP BY je.company_code, date_trunc('month', je.posting_date)::DATE
),
burn_metrics AS (
    SELECT
        company_code,
        COALESCE(CAST(AVG(burn_amount) AS NUMERIC(20, 4)), 0::NUMERIC(20, 4)) AS average_monthly_burn,
        MIN(month_start) AS burn_window_start,
        MAX(month_start) AS burn_window_end
    FROM recent_monthly_burn
    GROUP BY company_code
)
SELECT
    company.company_code,
    as_of_anchor.as_of_date,
    COALESCE(cash_position.current_cash, 0)::NUMERIC(20, 4) AS current_cash,
    COALESCE(burn_metrics.average_monthly_burn, 0)::NUMERIC(20, 4) AS average_monthly_burn,
    COALESCE(burn_metrics.burn_window_start, (date_trunc('month', as_of_anchor.as_of_date)::DATE - INTERVAL '3 months')::DATE) AS burn_window_start,
    COALESCE(burn_metrics.burn_window_end, (date_trunc('month', as_of_anchor.as_of_date)::DATE - INTERVAL '1 month')::DATE) AS burn_window_end
FROM company_catalog AS company
CROSS JOIN as_of_anchor
LEFT JOIN cash_position
    ON cash_position.company_code = company.company_code
LEFT JOIN burn_metrics
    ON burn_metrics.company_code = company.company_code;

COMMIT;
