-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0094: Core Audit Logs (v3.0)
-- Sổ Vàng Kiểm Toán — CẤM DELETE, CẤM UPDATE
-- Audit Trail trigger tự động cho mọi bảng trọng yếu
-- ══════════════════════════════════════════════════════════════════

BEGIN;

-- ═══════════════════════════════════════════════
-- 1. CORE.AUDIT_LOGS — Sổ Vàng Kiểm Toán
-- ⚠️ KHÔNG AI ĐƯỢC QUYỀN DELETE/TRUNCATE — KỂ CẢ CTO
-- ═══════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS core.audit_logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_schema TEXT NOT NULL,
    table_name   TEXT NOT NULL,
    operation    TEXT NOT NULL CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    old_data     JSONB,
    new_data     JSONB,
    changed_by   UUID,
    ip_address   INET,
    user_agent   TEXT,
    reason       TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
    -- ⚠️ APPEND-ONLY: KHÔNG CÓ updated_at, deleted_at
);

COMMENT ON TABLE core.audit_logs IS 'SỔ VÀNG KIỂM TOÁN — Ghi lại MỌI thay đổi trên các bảng trọng yếu. CẤM DELETE/UPDATE bảng này.';
COMMENT ON COLUMN core.audit_logs.reason IS 'Lý do thay đổi — BẮT BUỘC cung cấp khi UPDATE dữ liệu nhạy cảm';

-- ═══════════════════════════════════════════════
-- 2. AUDIT TRIGGER FUNCTION
-- Tự động ghi log mọi INSERT/UPDATE/DELETE
-- ═══════════════════════════════════════════════

CREATE OR REPLACE FUNCTION core.fn_audit_trigger()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO core.audit_logs (
        table_schema, table_name, operation,
        old_data, new_data, changed_by
    ) VALUES (
        TG_TABLE_SCHEMA,
        TG_TABLE_NAME,
        TG_OP,
        CASE WHEN TG_OP IN ('UPDATE', 'DELETE')
             THEN row_to_json(OLD)::JSONB END,
        CASE WHEN TG_OP IN ('INSERT', 'UPDATE')
             THEN row_to_json(NEW)::JSONB END,
        NULLIF(current_setting('app.current_user', true), '')::UUID
    );
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION core.fn_audit_trigger IS 'Trigger function ghi audit log tự động. Dùng SECURITY DEFINER để bypass RLS khi ghi audit.';

-- ═══════════════════════════════════════════════
-- 3. GẮN TRIGGER VÀO CÁC BẢNG TRỌNG YẾU
-- ═══════════════════════════════════════════════

-- Users
CREATE TRIGGER trg_audit_users
    AFTER INSERT OR UPDATE OR DELETE ON core.users
    FOR EACH ROW EXECUTE FUNCTION core.fn_audit_trigger();

-- Global Athletes
CREATE TRIGGER trg_audit_athletes
    AFTER INSERT OR UPDATE OR DELETE ON core.global_athletes
    FOR EACH ROW EXECUTE FUNCTION core.fn_audit_trigger();

-- Belt History
CREATE TRIGGER trg_audit_belts
    AFTER INSERT ON core.belt_history
    FOR EACH ROW EXECUTE FUNCTION core.fn_audit_trigger();

-- Federations
CREATE TRIGGER trg_audit_federations
    AFTER INSERT OR UPDATE OR DELETE ON core.federations
    FOR EACH ROW EXECUTE FUNCTION core.fn_audit_trigger();

-- ═══════════════════════════════════════════════
-- 4. RLS — PROTECTION POLICIES CHO AUDIT LOGS
-- ═══════════════════════════════════════════════

ALTER TABLE core.audit_logs ENABLE ROW LEVEL SECURITY;

-- Policy: CẤM DELETE — Không ai được xóa audit logs
CREATE POLICY audit_no_delete ON core.audit_logs
    FOR DELETE USING (false);

-- Policy: CẤM UPDATE — Không ai được sửa audit logs
CREATE POLICY audit_no_update ON core.audit_logs
    FOR UPDATE USING (false);

-- Policy: Chỉ admin và federation_admin được đọc
CREATE POLICY audit_select_admin ON core.audit_logs
    FOR SELECT USING (true);
    -- Trong production: USING (current_setting('app.current_role', true) IN ('admin', 'federation_admin'))

-- Policy: Cho phép trigger INSERT (SECURITY DEFINER)
CREATE POLICY audit_insert_trigger ON core.audit_logs
    FOR INSERT WITH CHECK (true);

-- ═══════════════════════════════════════════════
-- 5. INDEXES
-- ═══════════════════════════════════════════════

CREATE INDEX IF NOT EXISTS idx_core_audit_table
    ON core.audit_logs (table_schema, table_name, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_core_audit_operation
    ON core.audit_logs (operation, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_core_audit_changed_by
    ON core.audit_logs (changed_by) WHERE changed_by IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_core_audit_created
    ON core.audit_logs (created_at DESC);

COMMIT;
