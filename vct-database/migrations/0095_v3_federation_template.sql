-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0095: Federation Schema Template (v3.0)
-- Club Domain — Vệ tinh quay quanh Core Hub
-- Schema pattern: fed_{province_code}
-- ══════════════════════════════════════════════════════════════════

BEGIN;

-- ═══════════════════════════════════════════════
-- 1. FUNCTION: Tạo Federation Schema cho 1 tỉnh
-- Gọi: SELECT core.create_federation_schema('lamdong');
-- ═══════════════════════════════════════════════

CREATE OR REPLACE FUNCTION core.create_federation_schema(
    p_province_code TEXT
)
RETURNS VOID AS $$
DECLARE
    v_schema TEXT;
BEGIN
    v_schema := 'fed_' || lower(p_province_code);

    -- Tạo schema
    EXECUTE format('CREATE SCHEMA IF NOT EXISTS %I', v_schema);
    EXECUTE format(
        'COMMENT ON SCHEMA %I IS %L',
        v_schema,
        'VCT Federation Schema — Liên đoàn tỉnh ' || p_province_code
    );

    -- ── Bảng CLB/Võ Đường ──
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.clubs (
            id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            federation_id   UUID NOT NULL REFERENCES core.federations(id),
            name            TEXT NOT NULL,
            master_id       UUID REFERENCES core.users(id),
            master_name     TEXT,
            address         TEXT,
            phone           TEXT,
            email           TEXT,
            district        TEXT,
            is_active       BOOLEAN NOT NULL DEFAULT true,
            annual_fee_paid BOOLEAN NOT NULL DEFAULT false,
            total_members   INT NOT NULL DEFAULT 0,
            founded_date    DATE,
            created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
            deleted_at      TIMESTAMPTZ
        )', v_schema);

    -- ── Bảng Lịch sử Ghi danh (Timeline History) ──
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.club_memberships (
            id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            athlete_id    UUID NOT NULL REFERENCES core.global_athletes(id),
            club_id       UUID NOT NULL,
            start_date    DATE NOT NULL DEFAULT CURRENT_DATE,
            end_date      DATE,
            status        TEXT NOT NULL DEFAULT ''ACTIVE''
                            CHECK (status IN (''ACTIVE'', ''ARCHIVED'', ''SUSPENDED'')),
            transfer_note TEXT,
            approved_by   UUID REFERENCES core.users(id),
            created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
            deleted_at    TIMESTAMPTZ
        )', v_schema);

    -- FK cho club_id → clubs(id) trong cùng schema
    EXECUTE format('
        ALTER TABLE %I.club_memberships
        ADD CONSTRAINT fk_memberships_club
        FOREIGN KEY (club_id) REFERENCES %I.clubs(id)
    ', v_schema, v_schema);

    -- ── Bảng Điểm danh ──
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.attendance (
            id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            athlete_id    UUID NOT NULL,
            club_id       UUID NOT NULL,
            check_in_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
            session_type  TEXT DEFAULT ''regular''
                            CHECK (session_type IN (''regular'', ''extra'', ''exam'', ''event'')),
            notes         TEXT,
            recorded_by   UUID REFERENCES core.users(id)
        )', v_schema);

    -- ── Bảng Thu/Chi Nội bộ ──
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.local_finances (
            id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            club_id       UUID NOT NULL,
            athlete_id    UUID,
            type          TEXT NOT NULL CHECK (type IN (''income'', ''expense'')),
            category      TEXT NOT NULL,
            amount        NUMERIC(12,2) NOT NULL,
            description   TEXT,
            receipt_url   TEXT,
            recorded_by   UUID REFERENCES core.users(id),
            recorded_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
            deleted_at    TIMESTAMPTZ
        )', v_schema);

    -- ── Indexes ──
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_clubs_fed ON %I.clubs (federation_id) WHERE deleted_at IS NULL', p_province_code, v_schema);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_memberships_athlete ON %I.club_memberships (athlete_id)', p_province_code, v_schema);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_memberships_club ON %I.club_memberships (club_id, status)', p_province_code, v_schema);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_memberships_active ON %I.club_memberships (athlete_id) WHERE status = ''ACTIVE''', p_province_code, v_schema);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_attendance_club ON %I.attendance (club_id, check_in_at DESC)', p_province_code, v_schema);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_finances_club ON %I.local_finances (club_id, recorded_at DESC)', p_province_code, v_schema);

    -- ── Ghi event log ──
    PERFORM core.log_event(
        'FEDERATION_SCHEMA_CREATED',
        'federation',
        gen_random_uuid(),
        jsonb_build_object('province_code', p_province_code, 'schema_name', v_schema),
        NULL,
        NULLIF(current_setting('app.current_user', true), '')::UUID
    );

    RAISE NOTICE 'Federation schema % created successfully.', v_schema;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION core.create_federation_schema IS 'Tạo schema mới cho Liên đoàn tỉnh. VD: SELECT core.create_federation_schema(''lamdong'');';

-- ═══════════════════════════════════════════════
-- 2. TẠO SCHEMA MẪU: fed_lamdong (Lâm Đồng)
-- ═══════════════════════════════════════════════

SELECT core.create_federation_schema('lamdong');

-- ═══════════════════════════════════════════════
-- 3. SEED — Liên đoàn Lâm Đồng
-- ═══════════════════════════════════════════════

INSERT INTO core.federations (name, region, province_code, is_active)
VALUES ('Liên đoàn Võ Cổ Truyền Lâm Đồng', 'tinh', 'lamdong', true)
ON CONFLICT (province_code) DO NOTHING;

COMMIT;
