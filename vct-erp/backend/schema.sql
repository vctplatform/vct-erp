BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE account_category AS ENUM (
    'asset',
    'liability',
    'equity',
    'revenue',
    'expense',
    'contra_asset',
    'contra_liability',
    'contra_equity',
    'contra_revenue',
    'contra_expense'
);

CREATE TYPE normal_side AS ENUM ('debit', 'credit');
CREATE TYPE journal_entry_status AS ENUM ('draft', 'posted', 'reversed');
CREATE TYPE voucher_type AS ENUM ('PT', 'PC', 'PK');
CREATE TYPE outbox_status AS ENUM ('pending', 'processing', 'published', 'failed');
CREATE TYPE idempotency_status AS ENUM ('processing', 'completed', 'failed');
CREATE TYPE revenue_schedule_status AS ENUM ('scheduled', 'recognized', 'failed');
CREATE TYPE receivable_status AS ENUM ('open', 'settled', 'voided');
CREATE TYPE deposit_status AS ENUM ('held', 'released', 'applied', 'forfeited');
CREATE TYPE statement_line_status AS ENUM ('open', 'matched', 'manual', 'ignored');

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION set_account_depth()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    parent_depth SMALLINT;
    parent_company_code VARCHAR(32);
    cycle_found BOOLEAN;
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.depth := 0;
        RETURN NEW;
    END IF;

    IF NEW.parent_id = NEW.id THEN
        RAISE EXCEPTION 'account cannot be its own parent';
    END IF;

    SELECT depth, company_code
    INTO parent_depth, parent_company_code
    FROM accounts
    WHERE id = NEW.parent_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'parent account % not found', NEW.parent_id;
    END IF;

    IF parent_company_code <> NEW.company_code THEN
        RAISE EXCEPTION 'parent account company mismatch for %', NEW.code;
    END IF;

    WITH RECURSIVE parent_chain AS (
        SELECT id, parent_id
        FROM accounts
        WHERE id = NEW.parent_id

        UNION ALL

        SELECT a.id, a.parent_id
        FROM accounts a
        INNER JOIN parent_chain pc ON a.id = pc.parent_id
    )
    SELECT EXISTS (
        SELECT 1
        FROM parent_chain
        WHERE id = NEW.id
    )
    INTO cycle_found;

    IF cycle_found THEN
        RAISE EXCEPTION 'account hierarchy cycle detected for %', NEW.code;
    END IF;

    NEW.depth := parent_depth + 1;
    RETURN NEW;
END;
$$;

CREATE TABLE accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_code    VARCHAR(32) NOT NULL DEFAULT 'VCT_GROUP',
    code            VARCHAR(32) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    parent_id       UUID REFERENCES accounts(id) ON DELETE RESTRICT,
    depth           SMALLINT NOT NULL DEFAULT 0,
    account_type    account_category NOT NULL,
    normal_side     normal_side NOT NULL,
    is_postable     BOOLEAN NOT NULL DEFAULT TRUE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    description     TEXT,
    metadata        JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_accounts_company_code UNIQUE (company_code, code),
    CONSTRAINT ck_accounts_parent_not_self CHECK (id IS DISTINCT FROM parent_id)
);

CREATE INDEX idx_accounts_company_parent ON accounts (company_code, parent_id);
CREATE INDEX idx_accounts_company_active ON accounts (company_code, is_active) WHERE is_active;

CREATE TRIGGER tr_accounts_set_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER tr_accounts_set_depth
    BEFORE INSERT OR UPDATE OF parent_id, company_code
    ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION set_account_depth();

CREATE TABLE voucher_sequences (
    company_code    VARCHAR(32) NOT NULL,
    voucher_type    voucher_type NOT NULL,
    period_key      CHAR(7) NOT NULL,
    last_value      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (company_code, voucher_type, period_key),
    CONSTRAINT ck_voucher_sequences_last_value_positive CHECK (last_value >= 0)
);

CREATE TRIGGER tr_voucher_sequences_set_updated_at
    BEFORE UPDATE ON voucher_sequences
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE journal_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_no        VARCHAR(50) NOT NULL,
    voucher_type    voucher_type NOT NULL DEFAULT 'PK',
    company_code    VARCHAR(32) NOT NULL DEFAULT 'VCT_GROUP',
    source_module   VARCHAR(64) NOT NULL,
    external_ref    VARCHAR(100),
    description     TEXT NOT NULL,
    currency_code   CHAR(3) NOT NULL DEFAULT 'VND',
    posting_date    DATE NOT NULL,
    status          journal_entry_status NOT NULL DEFAULT 'posted',
    posted_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata        JSONB NOT NULL DEFAULT '{}'::JSONB,
    reversal_of_entry_id UUID REFERENCES journal_entries(id) ON DELETE RESTRICT,
    reversal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    reversed_at     TIMESTAMPTZ,
    void_reason     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_journal_entries_company_entry_no UNIQUE (company_code, entry_no),
    CONSTRAINT ck_journal_entries_reversal_self CHECK (
        id IS DISTINCT FROM reversal_of_entry_id
        AND id IS DISTINCT FROM reversal_entry_id
    )
);

CREATE INDEX idx_journal_entries_company_posting_date
    ON journal_entries (company_code, posting_date DESC);

CREATE INDEX idx_journal_entries_source_module
    ON journal_entries (source_module, external_ref);

CREATE UNIQUE INDEX uq_journal_entries_reversal_of_entry
    ON journal_entries (reversal_of_entry_id)
    WHERE reversal_of_entry_id IS NOT NULL;

CREATE INDEX idx_journal_entries_reversal_entry
    ON journal_entries (reversal_entry_id)
    WHERE reversal_entry_id IS NOT NULL;

CREATE TRIGGER tr_journal_entries_set_updated_at
    BEFORE UPDATE ON journal_entries
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE journal_items (
    journal_entry_id    UUID NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    line_no             SMALLINT NOT NULL,
    company_code        VARCHAR(32) NOT NULL,
    account_id          UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    side                normal_side NOT NULL,
    amount              NUMERIC(20, 4) NOT NULL,
    amount_signed       NUMERIC(20, 4) NOT NULL,
    currency_code       CHAR(3) NOT NULL DEFAULT 'VND',
    description         TEXT,
    metadata            JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (journal_entry_id, line_no, created_at),
    CONSTRAINT ck_journal_items_amount_positive CHECK (amount > 0),
    CONSTRAINT ck_journal_items_signed_consistency CHECK (
        (side = 'debit' AND amount_signed = amount)
        OR
        (side = 'credit' AND amount_signed = amount * -1)
    )
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_journal_items_company_account_created_at
    ON journal_items (company_code, account_id, created_at DESC);

CREATE INDEX idx_journal_items_company_entry_created_at
    ON journal_items (company_code, journal_entry_id, created_at DESC);

CREATE INDEX brin_journal_items_created_at
    ON journal_items USING BRIN (created_at);

CREATE OR REPLACE FUNCTION validate_journal_item_dimension()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    account_company_code VARCHAR(32);
    account_is_postable BOOLEAN;
    account_is_active BOOLEAN;
    entry_company_code VARCHAR(32);
BEGIN
    SELECT company_code, is_postable, is_active
    INTO account_company_code, account_is_postable, account_is_active
    FROM accounts
    WHERE id = NEW.account_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'account % not found', NEW.account_id;
    END IF;

    IF NOT account_is_active THEN
        RAISE EXCEPTION 'account % is inactive', NEW.account_id;
    END IF;

    IF NOT account_is_postable THEN
        RAISE EXCEPTION 'account % does not accept postings', NEW.account_id;
    END IF;

    SELECT company_code
    INTO entry_company_code
    FROM journal_entries
    WHERE id = NEW.journal_entry_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'journal entry % not found', NEW.journal_entry_id;
    END IF;

    IF entry_company_code <> NEW.company_code THEN
        RAISE EXCEPTION 'journal item company mismatch for entry %', NEW.journal_entry_id;
    END IF;

    IF account_company_code <> NEW.company_code THEN
        RAISE EXCEPTION 'journal item company mismatch for account %', NEW.account_id;
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER tr_journal_items_validate_dimension
    BEFORE INSERT OR UPDATE ON journal_items
    FOR EACH ROW
    EXECUTE FUNCTION validate_journal_item_dimension();

CREATE OR REPLACE FUNCTION validate_posted_journal_entry_balance()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    debit_total NUMERIC(20, 4);
    credit_total NUMERIC(20, 4);
    line_count BIGINT;
BEGIN
    IF NEW.status <> 'posted' THEN
        RETURN NULL;
    END IF;

    SELECT
        COUNT(*),
        COALESCE(SUM(CASE WHEN side = 'debit' THEN amount ELSE 0 END), 0),
        COALESCE(SUM(CASE WHEN side = 'credit' THEN amount ELSE 0 END), 0)
    INTO line_count, debit_total, credit_total
    FROM journal_items
    WHERE journal_entry_id = NEW.id;

    IF line_count < 2 THEN
        RAISE EXCEPTION 'journal entry % must contain at least two journal lines', NEW.id;
    END IF;

    IF debit_total <> credit_total THEN
        RAISE EXCEPTION
            'journal entry % is unbalanced (debit %, credit %)',
            NEW.id,
            debit_total,
            credit_total;
    END IF;

    RETURN NULL;
END;
$$;

CREATE CONSTRAINT TRIGGER tr_journal_entries_validate_balance
    AFTER INSERT OR UPDATE OF status, posted_at
    ON journal_entries
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW
    EXECUTE FUNCTION validate_posted_journal_entry_balance();

CREATE OR REPLACE FUNCTION create_journal_items_quarter_partition(p_year INT, p_quarter INT)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    partition_start DATE;
    partition_end DATE;
    partition_name TEXT;
    account_index_name TEXT;
    entry_index_name TEXT;
    brin_index_name TEXT;
BEGIN
    IF p_quarter < 1 OR p_quarter > 4 THEN
        RAISE EXCEPTION 'quarter must be between 1 and 4';
    END IF;

    partition_start := make_date(p_year, ((p_quarter - 1) * 3) + 1, 1);
    partition_end := (partition_start + INTERVAL '3 months')::DATE;
    partition_name := format('journal_items_%s_q%s', p_year, p_quarter);
    account_index_name := format('idx_%s_company_account_created_at', partition_name);
    entry_index_name := format('idx_%s_company_entry_created_at', partition_name);
    brin_index_name := format('brin_%s_created_at', partition_name);

    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF journal_items FOR VALUES FROM (%L) TO (%L)',
        partition_name,
        partition_start,
        partition_end
    );

    EXECUTE format(
        'CREATE INDEX IF NOT EXISTS %I ON %I (company_code, account_id, created_at DESC)',
        account_index_name,
        partition_name
    );

    EXECUTE format(
        'CREATE INDEX IF NOT EXISTS %I ON %I (company_code, journal_entry_id, created_at DESC)',
        entry_index_name,
        partition_name
    );

    EXECUTE format(
        'CREATE INDEX IF NOT EXISTS %I ON %I USING BRIN (created_at)',
        brin_index_name,
        partition_name
    );
END;
$$;

DO $$
DECLARE
    year_no INT;
    quarter_no INT;
    current_year INT := EXTRACT(YEAR FROM CURRENT_DATE);
BEGIN
    FOR year_no IN current_year - 1 .. current_year + 2 LOOP
        FOR quarter_no IN 1 .. 4 LOOP
            PERFORM create_journal_items_quarter_partition(year_no, quarter_no);
        END LOOP;
    END LOOP;
END;
$$;

CREATE TABLE IF NOT EXISTS journal_items_default
    PARTITION OF journal_items DEFAULT;

CREATE TABLE account_balances (
    company_code            VARCHAR(32) NOT NULL,
    account_id              UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    currency_code           CHAR(3) NOT NULL DEFAULT 'VND',
    debit_balance           NUMERIC(20, 4) NOT NULL DEFAULT 0,
    credit_balance          NUMERIC(20, 4) NOT NULL DEFAULT 0,
    net_balance             NUMERIC(20, 4) NOT NULL DEFAULT 0,
    last_journal_entry_id   UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (company_code, account_id, currency_code),
    CONSTRAINT ck_account_balances_non_negative CHECK (
        debit_balance >= 0
        AND credit_balance >= 0
    )
);

CREATE INDEX idx_account_balances_updated_at
    ON account_balances (updated_at DESC);

CREATE TABLE outbox_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type  VARCHAR(64) NOT NULL,
    aggregate_id    UUID NOT NULL,
    event_type      VARCHAR(128) NOT NULL,
    stream_key      VARCHAR(128) NOT NULL,
    status          outbox_status NOT NULL DEFAULT 'pending',
    payload         JSONB NOT NULL,
    attempt_count   INTEGER NOT NULL DEFAULT 0,
    available_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    locked_at       TIMESTAMPTZ,
    locked_by       VARCHAR(128),
    published_at    TIMESTAMPTZ,
    last_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_outbox_events_status_available
    ON outbox_events (status, available_at, created_at);

CREATE INDEX idx_outbox_events_processing_lock
    ON outbox_events (status, locked_at)
    WHERE status = 'processing';

CREATE INDEX idx_outbox_events_aggregate
    ON outbox_events (aggregate_type, aggregate_id);

CREATE TRIGGER tr_outbox_events_set_updated_at
    BEFORE UPDATE ON outbox_events
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE idempotency_keys (
    scope               VARCHAR(100) NOT NULL,
    idempotency_key     VARCHAR(128) NOT NULL,
    request_hash        CHAR(64) NOT NULL,
    status              idempotency_status NOT NULL DEFAULT 'processing',
    resource_id         VARCHAR(64),
    response_payload    JSONB,
    locked_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at        TIMESTAMPTZ,
    last_error          TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (scope, idempotency_key)
);

CREATE INDEX idx_idempotency_keys_status_locked_at
    ON idempotency_keys (status, locked_at DESC);

CREATE TRIGGER tr_idempotency_keys_set_updated_at
    BEFORE UPDATE ON idempotency_keys
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE saas_contracts (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_code                    VARCHAR(32) NOT NULL DEFAULT 'VCT_GROUP',
    contract_no                     VARCHAR(64) NOT NULL,
    customer_ref                    VARCHAR(64) NOT NULL,
    cash_account_id                 UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    deferred_revenue_account_id     UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    recognized_revenue_account_id   UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    currency_code                   CHAR(3) NOT NULL DEFAULT 'VND',
    start_date                      DATE NOT NULL,
    end_date                        DATE NOT NULL,
    term_months                     SMALLINT NOT NULL,
    total_amount                    NUMERIC(20, 4) NOT NULL,
    source_ref                      VARCHAR(100),
    initial_journal_entry_id        UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_saas_contracts_company_contract UNIQUE (company_code, contract_no),
    CONSTRAINT ck_saas_contracts_term_months CHECK (term_months > 0),
    CONSTRAINT ck_saas_contracts_total_amount_positive CHECK (total_amount > 0),
    CONSTRAINT ck_saas_contracts_date_range CHECK (end_date >= start_date)
);

CREATE INDEX idx_saas_contracts_customer_ref
    ON saas_contracts (company_code, customer_ref);

CREATE TRIGGER tr_saas_contracts_set_updated_at
    BEFORE UPDATE ON saas_contracts
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE saas_revenue_schedules (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id                 UUID NOT NULL REFERENCES saas_contracts(id) ON DELETE CASCADE,
    sequence_no                 SMALLINT NOT NULL,
    service_month               DATE NOT NULL,
    amount                      NUMERIC(20, 4) NOT NULL,
    status                      revenue_schedule_status NOT NULL DEFAULT 'scheduled',
    recognized_journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    recognized_at               TIMESTAMPTZ,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_saas_revenue_schedules_contract_sequence UNIQUE (contract_id, sequence_no),
    CONSTRAINT ck_saas_revenue_schedules_amount_positive CHECK (amount > 0)
);

CREATE INDEX idx_saas_revenue_schedules_due
    ON saas_revenue_schedules (status, service_month);

CREATE TRIGGER tr_saas_revenue_schedules_set_updated_at
    BEFORE UPDATE ON saas_revenue_schedules
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE dojo_receivables (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_code            VARCHAR(32) NOT NULL DEFAULT 'VCT_GROUP',
    student_ref             VARCHAR(64) NOT NULL,
    billing_month           DATE NOT NULL,
    due_date                DATE NOT NULL,
    currency_code           CHAR(3) NOT NULL DEFAULT 'VND',
    receivable_account_id   UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    revenue_account_id      UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    amount_due              NUMERIC(20, 4) NOT NULL,
    amount_paid             NUMERIC(20, 4) NOT NULL DEFAULT 0,
    status                  receivable_status NOT NULL DEFAULT 'open',
    source_ref              VARCHAR(100),
    assessment_entry_id     UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    settlement_entry_id     UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_dojo_receivables_student_month UNIQUE (company_code, student_ref, billing_month),
    CONSTRAINT ck_dojo_receivables_amount_due_positive CHECK (amount_due > 0),
    CONSTRAINT ck_dojo_receivables_amount_paid_valid CHECK (amount_paid >= 0 AND amount_paid <= amount_due)
);

CREATE INDEX idx_dojo_receivables_status_due
    ON dojo_receivables (status, due_date, billing_month);

CREATE TRIGGER tr_dojo_receivables_set_updated_at
    BEFORE UPDATE ON dojo_receivables
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE rental_deposits (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_code            VARCHAR(32) NOT NULL DEFAULT 'VCT_GROUP',
    rental_order_id         VARCHAR(64) NOT NULL,
    customer_ref            VARCHAR(64) NOT NULL,
    cash_account_id         UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    holding_account_id      UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    currency_code           CHAR(3) NOT NULL DEFAULT 'VND',
    amount                  NUMERIC(20, 4) NOT NULL,
    status                  deposit_status NOT NULL DEFAULT 'held',
    held_entry_id           UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    released_entry_id       UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    held_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    released_at             TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_rental_deposits_order UNIQUE (company_code, rental_order_id),
    CONSTRAINT ck_rental_deposits_amount_positive CHECK (amount > 0)
);

CREATE INDEX idx_rental_deposits_status_held_at
    ON rental_deposits (status, held_at DESC);

CREATE TRIGGER tr_rental_deposits_set_updated_at
    BEFORE UPDATE ON rental_deposits
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE bank_statement_lines (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_code        VARCHAR(32) NOT NULL DEFAULT 'VCT_GROUP',
    bank_account_no     VARCHAR(64) NOT NULL,
    external_line_id    VARCHAR(128) NOT NULL,
    booking_date        DATE NOT NULL,
    value_date          DATE,
    reference_no        VARCHAR(128),
    description         TEXT,
    currency_code       CHAR(3) NOT NULL DEFAULT 'VND',
    amount_signed       NUMERIC(20, 4) NOT NULL,
    status              statement_line_status NOT NULL DEFAULT 'open',
    matched_entry_id    UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    matching_rule       VARCHAR(64),
    matched_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_bank_statement_lines_external UNIQUE (company_code, bank_account_no, external_line_id)
);

CREATE INDEX idx_bank_statement_lines_status_booking
    ON bank_statement_lines (company_code, bank_account_no, status, booking_date);

CREATE INDEX idx_bank_statement_lines_matched_entry
    ON bank_statement_lines (matched_entry_id)
    WHERE matched_entry_id IS NOT NULL;

CREATE TRIGGER tr_bank_statement_lines_set_updated_at
    BEFORE UPDATE ON bank_statement_lines
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

INSERT INTO accounts (
    company_code,
    code,
    name,
    parent_id,
    account_type,
    normal_side,
    is_postable,
    description,
    metadata
)
SELECT
    seed.company_code,
    seed.code,
    seed.name,
    parent.id,
    seed.account_type::account_category,
    seed.normal_side::normal_side,
    seed.is_postable,
    seed.description,
    seed.metadata::jsonb
FROM (
    VALUES
        ('VCT_GROUP', '111', 'Tien mat', NULL::VARCHAR(32), 'asset', 'debit', FALSE, 'Tai khoan tong hop tien mat', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '112', 'Tien gui ngan hang', NULL, 'asset', 'debit', FALSE, 'Tai khoan tong hop tien gui ngan hang', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '131', 'Phai thu cua khach hang', NULL, 'asset', 'debit', TRUE, 'Cong no phai thu khach hang', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '133', 'Thue GTGT duoc khau tru', NULL, 'asset', 'debit', FALSE, 'Thue GTGT dau vao duoc khau tru', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '141', 'Tam ung', NULL, 'asset', 'debit', TRUE, 'Tam ung nhan vien va doi tac', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '242', 'Chi phi tra truoc', NULL, 'asset', 'debit', TRUE, 'Chi phi tra truoc ngan va dai han', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '331', 'Phai tra cho nguoi ban', NULL, 'liability', 'credit', TRUE, 'Cong no phai tra nha cung cap', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '333', 'Thue va cac khoan phai nop Nha nuoc', NULL, 'liability', 'credit', FALSE, 'Tai khoan tong hop nghia vu thue', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '334', 'Phai tra nguoi lao dong', NULL, 'liability', 'credit', TRUE, 'Cong no luong va thu nhap phai tra', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '338', 'Phai tra, phai nop khac', NULL, 'liability', 'credit', FALSE, 'Tai khoan tong hop cac khoan phai nop khac', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '411', 'Von dau tu cua chu so huu', NULL, 'equity', 'credit', TRUE, 'Von gop cua chu so huu', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '421', 'Loi nhuan sau thue chua phan phoi', NULL, 'equity', 'credit', TRUE, 'Loi nhuan sau thue chua phan phoi', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '511', 'Doanh thu ban hang va cung cap dich vu', NULL, 'revenue', 'credit', FALSE, 'Tai khoan tong hop doanh thu chinh', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '515', 'Doanh thu hoat dong tai chinh', NULL, 'revenue', 'credit', TRUE, 'Doanh thu tai chinh', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '521', 'Cac khoan giam tru doanh thu', NULL, 'contra_revenue', 'debit', FALSE, 'Giam tru doanh thu', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '632', 'Gia von hang ban', NULL, 'expense', 'debit', TRUE, 'Gia von hang hoa va dich vu', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '635', 'Chi phi tai chinh', NULL, 'expense', 'debit', TRUE, 'Chi phi hoat dong tai chinh', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '641', 'Chi phi ban hang', NULL, 'expense', 'debit', TRUE, 'Chi phi ban hang', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '642', 'Chi phi quan ly doanh nghiep', NULL, 'expense', 'debit', FALSE, 'Tai khoan tong hop chi phi quan ly doanh nghiep', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '711', 'Thu nhap khac', NULL, 'revenue', 'credit', TRUE, 'Thu nhap khac', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '811', 'Chi phi khac', NULL, 'expense', 'debit', TRUE, 'Chi phi khac', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '821', 'Chi phi thue thu nhap doanh nghiep', NULL, 'expense', 'debit', FALSE, 'Chi phi thue TNDN hien hanh va hoan lai', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '911', 'Xac dinh ket qua kinh doanh', NULL, 'equity', 'credit', TRUE, 'Tai khoan ket chuyen xac dinh ket qua kinh doanh', '{"profile":"VAS_TT200"}')
) AS seed(company_code, code, name, parent_code, account_type, normal_side, is_postable, description, metadata)
LEFT JOIN accounts AS parent
    ON parent.company_code = seed.company_code
   AND parent.code = seed.parent_code
ON CONFLICT (company_code, code) DO UPDATE
SET
    name = EXCLUDED.name,
    account_type = EXCLUDED.account_type,
    normal_side = EXCLUDED.normal_side,
    is_postable = EXCLUDED.is_postable,
    description = EXCLUDED.description,
    metadata = EXCLUDED.metadata,
    updated_at = NOW();

INSERT INTO accounts (
    company_code,
    code,
    name,
    parent_id,
    account_type,
    normal_side,
    is_postable,
    description,
    metadata
)
SELECT
    seed.company_code,
    seed.code,
    seed.name,
    parent.id,
    seed.account_type::account_category,
    seed.normal_side::normal_side,
    seed.is_postable,
    seed.description,
    seed.metadata::jsonb
FROM (
    VALUES
        ('VCT_GROUP', '1111', 'Tien mat Viet Nam', '111', 'asset', 'debit', TRUE, 'Tien mat VND', '{"profile":"VAS_TT200","system_key":"cash_vnd"}'),
        ('VCT_GROUP', '1121', 'Tien gui ngan hang Viet Nam dong', '112', 'asset', 'debit', TRUE, 'Tien gui ngan hang VND', '{"profile":"VAS_TT200","system_key":"bank_vnd"}'),
        ('VCT_GROUP', '1331', 'Thue GTGT duoc khau tru cua hang hoa dich vu', '133', 'asset', 'debit', TRUE, 'Thue GTGT dau vao duoc khau tru', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '33311', 'Thue GTGT dau ra phai nop', '333', 'liability', 'credit', TRUE, 'Thue GTGT dau ra', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '3334', 'Thue thu nhap doanh nghiep', '333', 'liability', 'credit', TRUE, 'Thue TNDN phai nop', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '3387', 'Doanh thu chua thuc hien', '338', 'liability', 'credit', TRUE, 'Doanh thu nhan truoc chua ghi nhan', '{"profile":"VAS_TT200","system_key":"deferred_revenue"}'),
        ('VCT_GROUP', '3388', 'Phai tra, phai nop khac', '338', 'liability', 'credit', TRUE, 'Cac khoan phai nop khac', '{"profile":"VAS_TT200","system_key":"deposit_holding"}'),
        ('VCT_GROUP', '5111', 'Doanh thu ban hang hoa', '511', 'revenue', 'credit', TRUE, 'Doanh thu ban hang hoa', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '5113', 'Doanh thu cung cap dich vu', '511', 'revenue', 'credit', TRUE, 'Doanh thu cung cap dich vu SaaS va Vo duong', '{"profile":"VAS_TT200","system_key":"service_revenue"}'),
        ('VCT_GROUP', '5211', 'Chiet khau thuong mai', '521', 'contra_revenue', 'debit', TRUE, 'Chiet khau thuong mai giam tru doanh thu', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '6421', 'Chi phi nhan vien quan ly', '642', 'expense', 'debit', TRUE, 'Chi phi luong va phuc loi bo phan quan ly', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '6422', 'Chi phi vat lieu quan ly', '642', 'expense', 'debit', TRUE, 'Chi phi vat lieu, cong cu dung cu quan ly', '{"profile":"VAS_TT200"}'),
        ('VCT_GROUP', '8211', 'Chi phi thue thu nhap doanh nghiep hien hanh', '821', 'expense', 'debit', TRUE, 'Chi phi thue TNDN hien hanh', '{"profile":"VAS_TT200"}')
) AS seed(company_code, code, name, parent_code, account_type, normal_side, is_postable, description, metadata)
INNER JOIN accounts AS parent
    ON parent.company_code = seed.company_code
   AND parent.code = seed.parent_code
ON CONFLICT (company_code, code) DO UPDATE
SET
    name = EXCLUDED.name,
    parent_id = EXCLUDED.parent_id,
    account_type = EXCLUDED.account_type,
    normal_side = EXCLUDED.normal_side,
    is_postable = EXCLUDED.is_postable,
    description = EXCLUDED.description,
    metadata = EXCLUDED.metadata,
    updated_at = NOW();

COMMIT;
