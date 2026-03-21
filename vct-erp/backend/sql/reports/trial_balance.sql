WITH account_ledger AS (
    SELECT
        a.id AS account_id,
        a.code AS account_code,
        a.name AS account_name,
        a.normal_side,
        je.posting_date,
        ji.side,
        ji.amount,
        CASE
            WHEN a.normal_side = 'debit' THEN ji.amount_signed
            ELSE ji.amount_signed * -1
        END AS natural_delta
    FROM accounts AS a
    LEFT JOIN journal_items AS ji
        ON ji.account_id = a.id
       AND ji.company_code = a.company_code
    LEFT JOIN journal_entries AS je
        ON je.id = ji.journal_entry_id
       AND je.company_code = a.company_code
    WHERE a.company_code = $1
      AND (je.status = 'posted' OR je.status IS NULL)
),
windowed AS (
    SELECT DISTINCT
        account_id,
        account_code,
        account_name,
        normal_side,
        COALESCE(SUM(CASE WHEN posting_date < $2 THEN natural_delta ELSE 0 END) OVER (PARTITION BY account_id), 0) AS opening_balance,
        COALESCE(SUM(CASE WHEN posting_date BETWEEN $2 AND $3 AND side = 'debit' THEN amount ELSE 0 END) OVER (PARTITION BY account_id), 0) AS period_debit,
        COALESCE(SUM(CASE WHEN posting_date BETWEEN $2 AND $3 AND side = 'credit' THEN amount ELSE 0 END) OVER (PARTITION BY account_id), 0) AS period_credit
    FROM account_ledger
)
SELECT
    account_code,
    account_name,
    normal_side,
    opening_balance,
    period_debit,
    period_credit,
    CASE
        WHEN normal_side = 'debit' THEN opening_balance + period_debit - period_credit
        ELSE opening_balance - period_debit + period_credit
    END AS closing_balance
FROM windowed
WHERE opening_balance <> 0
   OR period_debit <> 0
   OR period_credit <> 0
ORDER BY account_code;
