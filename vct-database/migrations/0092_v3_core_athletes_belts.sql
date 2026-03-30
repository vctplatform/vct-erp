-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0092: Core Athletes & Belt System (v3.0)
-- Global Athletes master records + Belt hierarchy + Belt history
-- Append-Only belt_history — Event Sourcing pattern
-- ══════════════════════════════════════════════════════════════════

BEGIN;

-- ═══════════════════════════════════════════════
-- 1. CORE.GLOBAL_ATHLETES — Master Record VĐV
-- "Cục Căn cước Quốc gia Võ sinh"
-- ═══════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS core.global_athletes (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID UNIQUE REFERENCES core.users(id),
    cccd           TEXT UNIQUE,
    full_name      TEXT NOT NULL,
    search_name    TEXT GENERATED ALWAYS AS (
                     lower(core.immutable_unaccent(full_name))
                   ) STORED,
    dob            DATE NOT NULL,
    gender         TEXT CHECK (gender IN ('male', 'female', 'other')),
    province       TEXT,
    address        TEXT,
    phone          TEXT,
    email          TEXT,
    nationality    TEXT NOT NULL DEFAULT 'Việt Nam',
    face_image_url TEXT,
    cccd_scan_url  TEXT,
    id_number      TEXT,
    weight         NUMERIC(5,2),
    height         NUMERIC(5,2),
    elo_rating     INT NOT NULL DEFAULT 1200,
    total_medals   INT NOT NULL DEFAULT 0,
    status         TEXT NOT NULL DEFAULT 'active'
                     CHECK (status IN ('active', 'suspended', 'retired', 'draft')),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ
);

COMMENT ON TABLE core.global_athletes IS 'Hồ sơ Võ sinh Toàn cầu — Master Record. Tách biệt hoàn toàn khỏi tournament/club.';
COMMENT ON COLUMN core.global_athletes.cccd IS 'Căn cước Công dân — Mã định danh duy nhất';
COMMENT ON COLUMN core.global_athletes.face_image_url IS 'Ảnh thẻ — Cloudflare R2 Storage';
COMMENT ON COLUMN core.global_athletes.cccd_scan_url IS 'Ảnh quét CCCD — Cloudflare R2 Storage';

-- ═══════════════════════════════════════════════
-- 2. CORE.BELT_LEVELS — Hệ thống Cấp Đai
-- ═══════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS core.belt_levels (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    rank_order  INT NOT NULL UNIQUE,
    color_hex   TEXT,
    description TEXT,
    branch      TEXT NOT NULL DEFAULT 'default',
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON TABLE core.belt_levels IS 'Bảng tra cứu Hệ thống Đai — Lam Đai, Lục Đai, Hoàng Đai...';
COMMENT ON COLUMN core.belt_levels.rank_order IS 'Thứ tự tăng dần: 1 → Sơ cấp, 10 → Cao cấp nhất';
COMMENT ON COLUMN core.belt_levels.branch IS 'Môn phái/Chi nhánh: default, vocotruyen, taekwondo...';

-- Seed dữ liệu Đai Võ Cổ Truyền
INSERT INTO core.belt_levels (name, rank_order, color_hex, description, branch) VALUES
    ('Bạch Đai (Trắng)',    1,  '#FFFFFF', 'Cấp nhập môn',                  'vocotruyen'),
    ('Lam Đai (Xanh Dương)',2,  '#0000FF', 'Cấp cơ bản',                    'vocotruyen'),
    ('Lục Đai (Xanh Lá)',   3,  '#00FF00', 'Cấp trung bình',                'vocotruyen'),
    ('Hoàng Đai (Vàng)',    4,  '#FFD700', 'Cấp khá',                        'vocotruyen'),
    ('Hồng Đai (Hồng)',     5,  '#FF69B4', 'Cấp giỏi',                      'vocotruyen'),
    ('Hồng Đai I',          6,  '#FF1493', 'Huyền Đai Đệ nhất đẳng',        'vocotruyen'),
    ('Hồng Đai II',         7,  '#DC143C', 'Huyền Đai Đệ nhị đẳng',         'vocotruyen'),
    ('Hồng Đai III',        8,  '#B22222', 'Huyền Đai Đệ tam đẳng',         'vocotruyen'),
    ('Huyền Đai (Đen)',     9,  '#000000', 'Huyền Đai — Bậc Thầy',          'vocotruyen'),
    ('Bạch Kim Đai',        10, '#E5E4E2', 'Đại sư — Cấp cao nhất',         'vocotruyen')
ON CONFLICT (rank_order) DO NOTHING;

-- ═══════════════════════════════════════════════
-- 3. CORE.BELT_HISTORY — Lịch sử Thăng Đai
-- ⚠️ APPEND-ONLY: KHÔNG CÓ updated_at, deleted_at
-- ═══════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS core.belt_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    athlete_id      UUID NOT NULL REFERENCES core.global_athletes(id),
    belt_level_id   UUID NOT NULL REFERENCES core.belt_levels(id),
    exam_date       DATE NOT NULL,
    examiner_id     UUID REFERENCES core.users(id),
    federation_id   UUID REFERENCES core.federations(id),
    certificate_url TEXT,
    qr_code_data    TEXT,
    source_event    TEXT NOT NULL DEFAULT 'manual'
                      CHECK (source_event IN ('manual', 'tournament_sync', 'import')),
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
    -- ⚠️ KHÔNG CÓ updated_at, deleted_at → APPEND-ONLY, BẤT BIẾN
);

COMMENT ON TABLE core.belt_history IS 'Lịch sử Thăng Đai — APPEND-ONLY. Mỗi record là 1 sự kiện thăng đai. KHÔNG BAO GIỜ UPDATE hay DELETE.';
COMMENT ON COLUMN core.belt_history.source_event IS 'Nguồn: manual (thi đai), tournament_sync (từ giải đấu), import (nhập liệu)';

-- ═══════════════════════════════════════════════
-- 4. TRIGGERS
-- ═══════════════════════════════════════════════

CREATE TRIGGER trg_athletes_updated_at
    BEFORE UPDATE ON core.global_athletes
    FOR EACH ROW EXECUTE FUNCTION core.trigger_set_updated_at();

-- ═══════════════════════════════════════════════
-- 5. INDEXES
-- ═══════════════════════════════════════════════

CREATE INDEX IF NOT EXISTS idx_core_athletes_search
    ON core.global_athletes USING gin (search_name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_core_athletes_user
    ON core.global_athletes (user_id) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_core_athletes_cccd
    ON core.global_athletes (cccd) WHERE cccd IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_core_athletes_province
    ON core.global_athletes (province) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_core_athletes_status
    ON core.global_athletes (status) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_core_belt_history_athlete
    ON core.belt_history (athlete_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_core_belt_history_federation
    ON core.belt_history (federation_id) WHERE federation_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_core_belt_levels_branch
    ON core.belt_levels (branch, rank_order);

COMMIT;
