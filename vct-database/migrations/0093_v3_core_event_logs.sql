-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0093: Core Event Logs (v3.0)
-- Global Event Sourcing — Append-Only Event Log
-- Trái tim của hệ thống Event-Sourcing toàn cầu
-- ══════════════════════════════════════════════════════════════════

BEGIN;

-- ═══════════════════════════════════════════════
-- 1. CORE.EVENT_LOGS — Nhật ký Sự kiện Toàn cầu
-- ⚠️ APPEND-ONLY: KHÔNG UPDATE, KHÔNG DELETE
-- ═══════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS core.event_logs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type    TEXT NOT NULL,
    -- Các event_type quan trọng:
    -- 'MEDAL_WON'                — VĐV đoạt huy chương
    -- 'BELT_PROMOTED'            — Thăng đai
    -- 'ATHLETE_SUSPENDED'        — Đình chỉ thi đấu
    -- 'CLUB_TRANSFER'            — Chuyển CLB
    -- 'TOURNAMENT_CLOSED'        — Bế mạc giải đấu
    -- 'TOURNAMENT_CREATED'       — Mở giải đấu
    -- 'CRISIS_ALERT'             — Lệnh cấm khẩn cấp
    -- 'REGISTRATION_APPROVED'    — Duyệt đăng ký
    -- 'PAYMENT_RECEIVED'         — Nhận thanh toán

    entity_type   TEXT NOT NULL,
    -- entity_type: 'athlete', 'tournament', 'club', 'federation', 'user'

    entity_id     UUID NOT NULL,
    payload       JSONB NOT NULL DEFAULT '{}',
    source_schema TEXT,
    actor_id      UUID REFERENCES core.users(id),
    correlation_id UUID,
    ip_address    INET,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
    -- ⚠️ APPEND-ONLY: Không UPDATE, Không DELETE — Bất biến
);

COMMENT ON TABLE core.event_logs IS 'Nhật ký Sự kiện Toàn cầu — APPEND-ONLY. Trái tim Event-Sourcing. Mọi sự kiện quan trọng trong hệ thống đều được ghi nhận tại đây.';
COMMENT ON COLUMN core.event_logs.event_type IS 'Loại sự kiện: MEDAL_WON, BELT_PROMOTED, CLUB_TRANSFER, TOURNAMENT_CLOSED, CRISIS_ALERT...';
COMMENT ON COLUMN core.event_logs.payload IS 'Dữ liệu JSON chi tiết của sự kiện — VD: {"medal": "GOLD", "category": "Nam_50kg"}';
COMMENT ON COLUMN core.event_logs.source_schema IS 'Schema gốc phát sự kiện — VD: t_2026_quocgia_001';
COMMENT ON COLUMN core.event_logs.correlation_id IS 'ID nhóm các sự kiện liên quan — VD: tất cả events của 1 giải đấu';

-- ═══════════════════════════════════════════════
-- 2. INDEXES
-- ═══════════════════════════════════════════════

CREATE INDEX IF NOT EXISTS idx_core_events_type
    ON core.event_logs (event_type, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_core_events_entity
    ON core.event_logs (entity_type, entity_id);

CREATE INDEX IF NOT EXISTS idx_core_events_actor
    ON core.event_logs (actor_id) WHERE actor_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_core_events_created
    ON core.event_logs (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_core_events_correlation
    ON core.event_logs (correlation_id) WHERE correlation_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_core_events_source
    ON core.event_logs (source_schema) WHERE source_schema IS NOT NULL;

-- GIN index cho payload search
CREATE INDEX IF NOT EXISTS idx_core_events_payload
    ON core.event_logs USING gin (payload jsonb_path_ops);

-- ═══════════════════════════════════════════════
-- 3. HELPER FUNCTION — Ghi sự kiện
-- ═══════════════════════════════════════════════

CREATE OR REPLACE FUNCTION core.log_event(
    p_event_type    TEXT,
    p_entity_type   TEXT,
    p_entity_id     UUID,
    p_payload       JSONB DEFAULT '{}',
    p_source_schema TEXT DEFAULT NULL,
    p_actor_id      UUID DEFAULT NULL,
    p_correlation_id UUID DEFAULT NULL
)
RETURNS UUID AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO core.event_logs (
        event_type, entity_type, entity_id,
        payload, source_schema, actor_id, correlation_id
    ) VALUES (
        p_event_type, p_entity_type, p_entity_id,
        p_payload, p_source_schema, p_actor_id, p_correlation_id
    )
    RETURNING id INTO v_id;

    -- Notify listeners for real-time processing
    PERFORM pg_notify('vct_events', json_build_object(
        'id', v_id,
        'event_type', p_event_type,
        'entity_type', p_entity_type,
        'entity_id', p_entity_id
    )::TEXT);

    RETURN v_id;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION core.log_event IS 'Ghi sự kiện vào Event Log + notify real-time listeners';

COMMIT;
