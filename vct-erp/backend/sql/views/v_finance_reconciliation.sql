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
