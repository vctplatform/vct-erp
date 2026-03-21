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
