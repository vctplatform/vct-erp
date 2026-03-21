WITH journal_lines AS (
    SELECT
        je.posting_date,
        je.entry_no,
        je.voucher_type,
        je.source_module,
        je.external_ref,
        je.description AS entry_description,
        ji.line_no,
        a.code AS account_code,
        a.name AS account_name,
        ji.side,
        ji.amount,
        ji.amount_signed
    FROM journal_entries AS je
    INNER JOIN journal_items AS ji ON ji.journal_entry_id = je.id
    INNER JOIN accounts AS a ON a.id = ji.account_id
    WHERE je.company_code = $1
      AND je.posting_date BETWEEN $2 AND $3
      AND je.status IN ('posted', 'reversed')
),
windowed AS (
    SELECT
        *,
        SUM(CASE WHEN side = 'debit' THEN amount ELSE 0 END)
            OVER (ORDER BY posting_date, entry_no, line_no) AS running_debit,
        SUM(CASE WHEN side = 'credit' THEN amount ELSE 0 END)
            OVER (ORDER BY posting_date, entry_no, line_no) AS running_credit
    FROM journal_lines
)
SELECT
    posting_date,
    entry_no,
    voucher_type,
    source_module,
    external_ref,
    entry_description,
    line_no,
    account_code,
    account_name,
    CASE WHEN side = 'debit' THEN amount ELSE 0 END AS debit_amount,
    CASE WHEN side = 'credit' THEN amount ELSE 0 END AS credit_amount,
    running_debit,
    running_credit
FROM windowed
ORDER BY posting_date, entry_no, line_no;
