BEGIN;

-- Legacy VAS_TT200 profile for VCT.
-- Note: from 2026-01-01, Thong tu 99/2025/TT-BTC replaced TT200.
-- This migration keeps a TT200-compatible chart because the current rollout explicitly targets it.

DO $$
BEGIN
    CREATE TYPE voucher_type AS ENUM ('PT', 'PC', 'PK');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END;
$$;

DO $$
BEGIN
    CREATE TYPE statement_line_status AS ENUM ('open', 'matched', 'manual', 'ignored');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END;
$$;

CREATE TABLE IF NOT EXISTS voucher_sequences (
    company_code    VARCHAR(32) NOT NULL,
    voucher_type    voucher_type NOT NULL,
    period_key      CHAR(7) NOT NULL,
    last_value      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (company_code, voucher_type, period_key),
    CONSTRAINT ck_voucher_sequences_last_value_positive CHECK (last_value >= 0)
);

ALTER TABLE journal_entries
    ADD COLUMN IF NOT EXISTS voucher_type voucher_type NOT NULL DEFAULT 'PK',
    ADD COLUMN IF NOT EXISTS reversal_of_entry_id UUID REFERENCES journal_entries(id) ON DELETE RESTRICT,
    ADD COLUMN IF NOT EXISTS reversal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS reversed_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS void_reason TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS uq_journal_entries_reversal_of_entry
    ON journal_entries (reversal_of_entry_id)
    WHERE reversal_of_entry_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_journal_entries_reversal_entry
    ON journal_entries (reversal_entry_id)
    WHERE reversal_entry_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS bank_statement_lines (
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

CREATE INDEX IF NOT EXISTS idx_bank_statement_lines_status_booking
    ON bank_statement_lines (company_code, bank_account_no, status, booking_date);

CREATE INDEX IF NOT EXISTS idx_bank_statement_lines_matched_entry
    ON bank_statement_lines (matched_entry_id)
    WHERE matched_entry_id IS NOT NULL;

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
