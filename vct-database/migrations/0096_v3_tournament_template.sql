-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0096: Tournament Schema Template (v3.0)
-- Arena/Dã chiến — Schema tách biệt cho mỗi giải đấu
-- Schema pattern: t_{year}_{scope}_{seq}
-- Archive pattern: arch_t_{year}_{scope}_{seq}
-- ══════════════════════════════════════════════════════════════════

BEGIN;

-- ═══════════════════════════════════════════════
-- 1. FUNCTION: Tạo Tournament Schema
-- Gọi: SELECT core.create_tournament_schema('2026', 'quocgia', '001');
-- ═══════════════════════════════════════════════

CREATE OR REPLACE FUNCTION core.create_tournament_schema(
    p_year  TEXT,
    p_scope TEXT,
    p_seq   TEXT
)
RETURNS TEXT AS $$
DECLARE
    v_schema TEXT;
BEGIN
    v_schema := 't_' || p_year || '_' || lower(p_scope) || '_' || p_seq;

    -- Tạo schema
    EXECUTE format('CREATE SCHEMA IF NOT EXISTS %I', v_schema);
    EXECUTE format(
        'COMMENT ON SCHEMA %I IS %L',
        v_schema,
        'VCT Tournament Arena — Giải ' || p_scope || ' ' || p_year || ' #' || p_seq
    );

    -- ── Bảng Đội hình (Snapshot từ core — IMMUTABLE sau khi khóa) ──
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.roster (
            id              UUID PRIMARY KEY,
            athlete_name    TEXT NOT NULL,
            cccd            TEXT,
            dob             DATE,
            gender          TEXT,
            belt_level_name TEXT NOT NULL,
            belt_rank_order INT,
            weight_kg       NUMERIC(5,2),
            height_cm       NUMERIC(5,2),
            club_name       TEXT,
            federation_name TEXT,
            province        TEXT,
            face_image_url  TEXT,
            is_locked       BOOLEAN NOT NULL DEFAULT false,
            is_suspended    BOOLEAN NOT NULL DEFAULT false,
            suspend_reason  TEXT,
            snapshot_at     TIMESTAMPTZ NOT NULL DEFAULT now()
        )', v_schema);

    -- ── Bảng Sơ đồ Nhánh đấu ──
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.brackets (
            id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            round_number    INT NOT NULL,
            match_order     INT NOT NULL,
            category        TEXT NOT NULL,
            red_athlete_id  UUID,
            blue_athlete_id UUID,
            winner_id       UUID,
            status          TEXT DEFAULT ''pending''
                              CHECK (status IN (''pending'', ''in_progress'', ''completed'', ''cancelled'')),
            scheduled_at    TIMESTAMPTZ,
            started_at      TIMESTAMPTZ,
            finished_at     TIMESTAMPTZ,
            arena_name      TEXT,
            notes           TEXT,
            created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
        )', v_schema);

    -- FK cho brackets → roster
    EXECUTE format('
        ALTER TABLE %I.brackets
        ADD CONSTRAINT fk_brackets_red FOREIGN KEY (red_athlete_id) REFERENCES %I.roster(id),
        ADD CONSTRAINT fk_brackets_blue FOREIGN KEY (blue_athlete_id) REFERENCES %I.roster(id)
    ', v_schema, v_schema, v_schema);

    -- ── Bảng Điểm Thi đấu (Redis → Flush mỗi cuối vòng) ──
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.match_scores (
            id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            bracket_id       UUID NOT NULL,
            referee_id       UUID NOT NULL,
            referee_name     TEXT,
            referee_position TEXT DEFAULT ''corner''
                               CHECK (referee_position IN (''corner'', ''center'', ''head'')),
            red_points       INT NOT NULL DEFAULT 0,
            blue_points      INT NOT NULL DEFAULT 0,
            penalties_red    INT DEFAULT 0,
            penalties_blue   INT DEFAULT 0,
            round_number     INT NOT NULL DEFAULT 1,
            notes            TEXT,
            scored_at        TIMESTAMPTZ NOT NULL DEFAULT now()
        )', v_schema);

    -- FK cho match_scores → brackets
    EXECUTE format('
        ALTER TABLE %I.match_scores
        ADD CONSTRAINT fk_scores_bracket FOREIGN KEY (bracket_id) REFERENCES %I.brackets(id)
    ', v_schema, v_schema);

    -- ── Bảng Kết quả Giải ──
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.results (
            id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            category      TEXT NOT NULL,
            gold_id       UUID,
            gold_name     TEXT,
            silver_id     UUID,
            silver_name   TEXT,
            bronze1_id    UUID,
            bronze1_name  TEXT,
            bronze2_id    UUID,
            bronze2_name  TEXT,
            is_finalized  BOOLEAN NOT NULL DEFAULT false,
            finalized_by  UUID,
            finalized_at  TIMESTAMPTZ,
            created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
        )', v_schema);

    -- ── Indexes ──
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_roster_locked ON %I.roster (is_locked)', v_schema, v_schema);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_roster_suspended ON %I.roster (is_suspended) WHERE is_suspended = true', v_schema, v_schema);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_brackets_cat ON %I.brackets (category, round_number)', v_schema, v_schema);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_brackets_status ON %I.brackets (status)', v_schema, v_schema);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%s_scores_bracket ON %I.match_scores (bracket_id)', v_schema, v_schema);

    -- ── Ghi event log ──
    PERFORM core.log_event(
        'TOURNAMENT_SCHEMA_CREATED',
        'tournament',
        gen_random_uuid(),
        jsonb_build_object(
            'year', p_year,
            'scope', p_scope,
            'seq', p_seq,
            'schema_name', v_schema
        )
    );

    RAISE NOTICE 'Tournament schema % created successfully.', v_schema;
    RETURN v_schema;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION core.create_tournament_schema IS 'Tạo schema mới cho giải đấu. VD: SELECT core.create_tournament_schema(''2026'', ''quocgia'', ''001'');';

-- ═══════════════════════════════════════════════
-- 2. FUNCTION: Export Roster (Clone VĐV → Tournament)
-- ═══════════════════════════════════════════════

CREATE OR REPLACE FUNCTION core.export_tournament_roster(
    p_tournament_schema TEXT,
    p_athlete_ids       UUID[]
)
RETURNS INT AS $$
DECLARE
    v_count INT := 0;
BEGIN
    EXECUTE format('
        INSERT INTO %I.roster (
            id, athlete_name, cccd, dob, gender,
            belt_level_name, belt_rank_order,
            weight_kg, height_cm, face_image_url
        )
        SELECT
            ga.id,
            ga.full_name,
            ga.cccd,
            ga.dob,
            ga.gender,
            COALESCE(bl.name, ''Chưa xếp hạng''),
            COALESCE(bl.rank_order, 0),
            ga.weight,
            ga.height,
            ga.face_image_url
        FROM core.global_athletes ga
        LEFT JOIN LATERAL (
            SELECT bh.belt_level_id
            FROM core.belt_history bh
            WHERE bh.athlete_id = ga.id
            ORDER BY bh.exam_date DESC, bh.created_at DESC
            LIMIT 1
        ) latest_belt ON true
        LEFT JOIN core.belt_levels bl ON bl.id = latest_belt.belt_level_id
        WHERE ga.id = ANY($1)
          AND ga.deleted_at IS NULL
        ON CONFLICT (id) DO NOTHING
    ', p_tournament_schema)
    USING p_athlete_ids;

    GET DIAGNOSTICS v_count = ROW_COUNT;
    RETURN v_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION core.export_tournament_roster IS 'Clone VĐV từ core vào roster giải đấu (Snapshot bất biến). Chạy trước ngày mở giải.';

-- ═══════════════════════════════════════════════
-- 3. FUNCTION: Archive Tournament Schema
-- Bế mạc → chuyển sang arch_t_{...}
-- ═══════════════════════════════════════════════

CREATE OR REPLACE FUNCTION core.archive_tournament_schema(
    p_tournament_schema TEXT
)
RETURNS TEXT AS $$
DECLARE
    v_archive TEXT;
BEGIN
    v_archive := 'arch_' || p_tournament_schema;

    -- Rename schema
    EXECUTE format('ALTER SCHEMA %I RENAME TO %I', p_tournament_schema, v_archive);

    -- Ghi event log
    PERFORM core.log_event(
        'TOURNAMENT_ARCHIVED',
        'tournament',
        gen_random_uuid(),
        jsonb_build_object(
            'original_schema', p_tournament_schema,
            'archive_schema', v_archive
        )
    );

    RAISE NOTICE 'Tournament schema % archived to %.', p_tournament_schema, v_archive;
    RETURN v_archive;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION core.archive_tournament_schema IS 'Bế mạc giải → Chuyển schema sang archive (read-only vĩnh viễn). VD: SELECT core.archive_tournament_schema(''t_2026_quocgia_001'');';

-- ═══════════════════════════════════════════════
-- 4. FUNCTION: Lock Roster (Khóa data trước vòng 1)
-- ═══════════════════════════════════════════════

CREATE OR REPLACE FUNCTION core.lock_tournament_roster(
    p_tournament_schema TEXT
)
RETURNS INT AS $$
DECLARE
    v_count INT;
BEGIN
    EXECUTE format('
        UPDATE %I.roster SET is_locked = true WHERE is_locked = false
    ', p_tournament_schema);

    GET DIAGNOSTICS v_count = ROW_COUNT;

    PERFORM core.log_event(
        'ROSTER_LOCKED',
        'tournament',
        gen_random_uuid(),
        jsonb_build_object(
            'schema', p_tournament_schema,
            'locked_count', v_count
        )
    );

    RETURN v_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION core.lock_tournament_roster IS 'Khóa roster trước vòng 1 — Sau khi khóa, thay đổi trên core KHÔNG ảnh hưởng roster.';

COMMIT;
