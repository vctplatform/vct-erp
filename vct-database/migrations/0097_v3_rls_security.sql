-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0097: RLS & Security (v3.0)
-- Row-Level Security cho tất cả Core tables
-- Phân quyền: Federation → Club → Athlete
-- ══════════════════════════════════════════════════════════════════

BEGIN;

-- ═══════════════════════════════════════════════
-- 1. ENABLE RLS trên tất cả Core tables
-- ═══════════════════════════════════════════════

ALTER TABLE core.users ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.global_athletes ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.belt_history ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.federations ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.event_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.sessions ENABLE ROW LEVEL SECURITY;
-- audit_logs đã bật RLS ở migration 0094

-- ═══════════════════════════════════════════════
-- 2. USERS POLICIES
-- ═══════════════════════════════════════════════

-- Admin xem tất cả users
CREATE POLICY users_admin_all ON core.users
    FOR ALL
    USING (
        current_setting('app.current_role', true) = 'admin'
    );

-- User xem chính mình
CREATE POLICY users_self_select ON core.users
    FOR SELECT
    USING (
        id::TEXT = current_setting('app.current_user', true)
    );

-- User update chính mình
CREATE POLICY users_self_update ON core.users
    FOR UPDATE
    USING (
        id::TEXT = current_setting('app.current_user', true)
    );

-- ═══════════════════════════════════════════════
-- 3. GLOBAL ATHLETES POLICIES
-- ═══════════════════════════════════════════════

-- Admin xem tất cả athletes
CREATE POLICY athletes_admin_all ON core.global_athletes
    FOR ALL
    USING (
        current_setting('app.current_role', true) = 'admin'
    );

-- Federation admin chỉ thấy VĐV thuộc tỉnh mình
CREATE POLICY federation_athletes_read ON core.global_athletes
    FOR SELECT
    USING (
        province = (
            SELECT province_code FROM core.federations
            WHERE admin_id::TEXT = current_setting('app.current_user', true)
            LIMIT 1
        )
        OR current_setting('app.current_role', true) = 'admin'
    );

-- Athlete xem chính mình
CREATE POLICY athletes_self_select ON core.global_athletes
    FOR SELECT
    USING (
        user_id::TEXT = current_setting('app.current_user', true)
    );

-- ═══════════════════════════════════════════════
-- 4. BELT HISTORY POLICIES
-- ═══════════════════════════════════════════════

-- Ai cũng có thể đọc belt history (thông tin công khai)
CREATE POLICY belt_history_read_all ON core.belt_history
    FOR SELECT
    USING (true);

-- Chỉ admin và federation_admin được INSERT
CREATE POLICY belt_history_insert ON core.belt_history
    FOR INSERT
    WITH CHECK (
        current_setting('app.current_role', true) IN ('admin', 'federation_admin')
    );

-- ═══════════════════════════════════════════════
-- 5. FEDERATIONS POLICIES
-- ═══════════════════════════════════════════════

-- Ai cũng đọc được danh sách liên đoàn
CREATE POLICY federations_read_all ON core.federations
    FOR SELECT
    USING (true);

-- Chỉ admin được sửa/xóa
CREATE POLICY federations_admin_write ON core.federations
    FOR ALL
    USING (
        current_setting('app.current_role', true) = 'admin'
    );

-- ═══════════════════════════════════════════════
-- 6. EVENT LOGS POLICIES
-- ═══════════════════════════════════════════════

-- Admin + federation_admin đọc event logs
CREATE POLICY events_admin_read ON core.event_logs
    FOR SELECT
    USING (
        current_setting('app.current_role', true) IN ('admin', 'federation_admin')
    );

-- Mọi function đều có thể INSERT (qua core.log_event)
CREATE POLICY events_insert ON core.event_logs
    FOR INSERT
    WITH CHECK (true);

-- CẤM UPDATE + DELETE event logs
CREATE POLICY events_no_update ON core.event_logs
    FOR UPDATE USING (false);
CREATE POLICY events_no_delete ON core.event_logs
    FOR DELETE USING (false);

-- ═══════════════════════════════════════════════
-- 7. SESSIONS POLICIES
-- ═══════════════════════════════════════════════

-- User chỉ thấy sessions của mình
CREATE POLICY sessions_self ON core.sessions
    FOR ALL
    USING (
        user_id::TEXT = current_setting('app.current_user', true)
        OR current_setting('app.current_role', true) = 'admin'
    );

-- ═══════════════════════════════════════════════
-- 8. BYPASS RLS cho Service Role
-- Backend service account vẫn cần truy cập đầy đủ
-- ═══════════════════════════════════════════════

-- Tạo role cho backend service (nếu chưa có)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'vct_service') THEN
        CREATE ROLE vct_service;
    END IF;
END $$;

-- Cấp quyền usage trên core schema
GRANT USAGE ON SCHEMA core TO vct_service;
GRANT ALL ON ALL TABLES IN SCHEMA core TO vct_service;
GRANT ALL ON ALL SEQUENCES IN SCHEMA core TO vct_service;

COMMENT ON POLICY users_admin_all ON core.users IS 'Admin có full access users';
COMMENT ON POLICY federation_athletes_read ON core.global_athletes IS 'Liên đoàn chỉ thấy VĐV tỉnh mình quản lý';

COMMIT;
