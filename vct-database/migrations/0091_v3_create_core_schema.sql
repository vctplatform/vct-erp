-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0091: Core Schema Foundation (v3.0 FINAL)
-- Creates the core schema with users and federations tables
-- Architecture: Hub-and-Spoke — core = Global Hub (Bất biến)
-- Approved by: Chairman Hoàng Bá Tùng — 30/03/2026
-- ══════════════════════════════════════════════════════════════════

BEGIN;

-- ═══════════════════════════════════════════════
-- 1. CREATE CORE SCHEMA
-- ═══════════════════════════════════════════════

CREATE SCHEMA IF NOT EXISTS core;

COMMENT ON SCHEMA core IS 'VCT Global Hub — Cục Căn cước Quốc gia của Võ Cổ Truyền. Dữ liệu bất biến: Users, Athletes, Belts, Events, Audit.';

-- ═══════════════════════════════════════════════
-- 2. EXTENSIONS
-- ═══════════════════════════════════════════════

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Immutable wrapper cho unaccent (Supabase yêu cầu IMMUTABLE cho GENERATED columns)
CREATE OR REPLACE FUNCTION core.immutable_unaccent(TEXT)
RETURNS TEXT AS $$
    SELECT public.unaccent($1);
$$ LANGUAGE sql IMMUTABLE PARALLEL SAFE;

COMMENT ON FUNCTION core.immutable_unaccent IS 'Immutable wrapper cho unaccent — cần cho GENERATED ALWAYS AS columns';

-- ═══════════════════════════════════════════════
-- 3. CORE.USERS — Bảng Người dùng (Tài khoản đăng nhập)
-- ═══════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS core.users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT UNIQUE NOT NULL,
    phone         TEXT UNIQUE,
    full_name     TEXT NOT NULL,
    search_name   TEXT GENERATED ALWAYS AS (
                    lower(core.immutable_unaccent(full_name))
                  ) STORED,
    password_hash TEXT,
    role          TEXT NOT NULL DEFAULT 'athlete'
                    CHECK (role IN (
                        'admin',
                        'federation_admin',
                        'club_owner',
                        'coach',
                        'athlete',
                        'referee',
                        'parent'
                    )),
    avatar_url    TEXT,
    is_active     BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

COMMENT ON TABLE core.users IS 'Tài khoản người dùng — Trung tâm xác thực toàn hệ thống VCT';
COMMENT ON COLUMN core.users.search_name IS 'Tên không dấu tự động sinh cho fuzzy search tiếng Việt';
COMMENT ON COLUMN core.users.deleted_at IS 'Soft-delete: NULL = hoạt động, NOT NULL = đã xóa';

-- ═══════════════════════════════════════════════
-- 4. CORE.FEDERATIONS — Bảng Liên đoàn
-- ═══════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS core.federations (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT NOT NULL,
    region        TEXT NOT NULL
                    CHECK (region IN ('quocgia', 'tinh', 'huyen')),
    province_code TEXT UNIQUE,
    admin_id      UUID REFERENCES core.users(id),
    address       TEXT,
    phone         TEXT,
    email         TEXT,
    website       TEXT,
    is_active     BOOLEAN NOT NULL DEFAULT true,
    founded_date  DATE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

COMMENT ON TABLE core.federations IS 'Liên đoàn VCT — Quốc gia, Tỉnh/Thành, Quận/Huyện';
COMMENT ON COLUMN core.federations.province_code IS 'Mã tỉnh: lamdong, binhdinh, ... — dùng tạo fed_{code} schema';

-- ═══════════════════════════════════════════════
-- 5. SESSIONS — Phiên đăng nhập
-- ═══════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS core.sessions (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID NOT NULL REFERENCES core.users(id) ON DELETE CASCADE,
    access_token_jti   TEXT UNIQUE NOT NULL,
    refresh_token_jti  TEXT UNIQUE NOT NULL,
    ip_address         INET,
    user_agent         TEXT,
    expires_at         TIMESTAMPTZ NOT NULL,
    refresh_expires_at TIMESTAMPTZ NOT NULL,
    revoked_at         TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON TABLE core.sessions IS 'JWT sessions — Access + Refresh token tracking';

-- ═══════════════════════════════════════════════
-- 6. UPDATED_AT TRIGGER
-- ═══════════════════════════════════════════════

CREATE OR REPLACE FUNCTION core.trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON core.users
    FOR EACH ROW EXECUTE FUNCTION core.trigger_set_updated_at();

CREATE TRIGGER trg_federations_updated_at
    BEFORE UPDATE ON core.federations
    FOR EACH ROW EXECUTE FUNCTION core.trigger_set_updated_at();

-- ═══════════════════════════════════════════════
-- 7. INDEXES
-- ═══════════════════════════════════════════════

CREATE INDEX IF NOT EXISTS idx_core_users_search
    ON core.users USING gin (search_name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_core_users_role
    ON core.users (role) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_core_users_email
    ON core.users (email) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_core_federations_province
    ON core.federations (province_code) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_core_federations_region
    ON core.federations (region) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_core_sessions_user
    ON core.sessions (user_id) WHERE revoked_at IS NULL;

COMMIT;
