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
      AND je.status = 'posted'
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
    gross_revenue,
    revenue_deductions,
    financial_income,
    other_income,
    gross_revenue - revenue_deductions + financial_income + other_income AS net_revenue
FROM cost_center_window
ORDER BY cost_center;
