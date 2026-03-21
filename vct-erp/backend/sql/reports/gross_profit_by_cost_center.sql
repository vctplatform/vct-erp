WITH pnl_source AS (
    SELECT
        COALESCE(je.metadata->>'cost_center', je.metadata->>'business_line', 'unassigned') AS cost_center,
        a.code AS account_code,
        ji.side,
        ji.amount
    FROM journal_entries AS je
    INNER JOIN journal_items AS ji ON ji.journal_entry_id = je.id
    INNER JOIN accounts AS a ON a.id = ji.account_id
    WHERE je.company_code = $1
      AND je.status = 'posted'
      AND je.posting_date BETWEEN $2 AND $3
      AND (
          a.code LIKE '511%' OR
          a.code LIKE '521%' OR
          a.code LIKE '515%' OR
          a.code LIKE '711%' OR
          a.code LIKE '632%'
      )
),
normalized AS (
    SELECT
        cost_center,
        account_code,
        CASE
            WHEN account_code LIKE '521%' THEN amount
            WHEN account_code LIKE '632%' THEN amount
            WHEN side = 'credit' THEN amount
            ELSE amount * -1
        END AS signed_amount
    FROM pnl_source
),
windowed AS (
    SELECT DISTINCT
        cost_center,
        COALESCE(SUM(CASE WHEN account_code LIKE '511%' THEN signed_amount ELSE 0 END) OVER (PARTITION BY cost_center), 0) AS gross_revenue,
        COALESCE(SUM(CASE WHEN account_code LIKE '521%' THEN signed_amount ELSE 0 END) OVER (PARTITION BY cost_center), 0) AS revenue_deductions,
        COALESCE(SUM(CASE WHEN account_code LIKE '515%' THEN signed_amount ELSE 0 END) OVER (PARTITION BY cost_center), 0) AS financial_income,
        COALESCE(SUM(CASE WHEN account_code LIKE '711%' THEN signed_amount ELSE 0 END) OVER (PARTITION BY cost_center), 0) AS other_income,
        COALESCE(SUM(CASE WHEN account_code LIKE '632%' THEN signed_amount ELSE 0 END) OVER (PARTITION BY cost_center), 0) AS cost_of_goods_sold
    FROM normalized
)
SELECT
    cost_center,
    gross_revenue,
    revenue_deductions,
    financial_income,
    other_income,
    cost_of_goods_sold,
    gross_revenue - revenue_deductions + financial_income + other_income - cost_of_goods_sold AS gross_profit
FROM windowed
ORDER BY cost_center;
