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
