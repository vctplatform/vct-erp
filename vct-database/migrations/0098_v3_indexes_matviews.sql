-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0098: Indexes & Materialized Views (v3.0)
-- Performance optimization: Dashboard CEO, National statistics
-- ══════════════════════════════════════════════════════════════════

BEGIN;

-- ═══════════════════════════════════════════════
-- 1. MATERIALIZED VIEW — Dashboard CEO Quốc gia
-- Cập nhật 15 phút/lần bằng pg_cron
-- ═══════════════════════════════════════════════

CREATE MATERIALIZED VIEW IF NOT EXISTS core.mv_national_dashboard AS
SELECT
    f.id                   AS federation_id,
    f.name                 AS federation_name,
    f.province_code,
    f.region,
    (
        SELECT COUNT(*)
        FROM core.global_athletes ga
        WHERE ga.province = f.province_code
          AND ga.deleted_at IS NULL
          AND ga.status = 'active'
    ) AS total_active_athletes,
    (
        SELECT COUNT(DISTINCT bh.id)
        FROM core.belt_history bh
        JOIN core.global_athletes ga ON bh.athlete_id = ga.id
        WHERE ga.province = f.province_code
    ) AS total_belt_certifications,
    (
        SELECT COUNT(*)
        FROM core.event_logs el
        WHERE el.event_type = 'MEDAL_WON'
          AND el.payload->>'province' = f.province_code
    ) AS total_medals_won,
    now() AS refreshed_at
FROM core.federations f
WHERE f.deleted_at IS NULL
  AND f.is_active = true;

COMMENT ON MATERIALIZED VIEW core.mv_national_dashboard IS 'Dashboard CEO — Tổng hợp thống kê toàn quốc theo tỉnh. Refresh mỗi 15 phút.';

-- Index cho dashboard
CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_dashboard_fed
    ON core.mv_national_dashboard (federation_id);

CREATE INDEX IF NOT EXISTS idx_mv_dashboard_region
    ON core.mv_national_dashboard (region);

-- ═══════════════════════════════════════════════
-- 2. MATERIALIZED VIEW — Bảng xếp hạng VĐV
-- ═══════════════════════════════════════════════

CREATE MATERIALIZED VIEW IF NOT EXISTS core.mv_athlete_rankings AS
SELECT
    ga.id              AS athlete_id,
    ga.full_name,
    ga.province,
    ga.elo_rating,
    ga.total_medals,
    ga.gender,
    bl.name            AS current_belt,
    bl.rank_order      AS belt_rank,
    bh.exam_date       AS last_belt_date,
    RANK() OVER (ORDER BY ga.elo_rating DESC) AS national_rank,
    RANK() OVER (
        PARTITION BY ga.province
        ORDER BY ga.elo_rating DESC
    ) AS province_rank,
    now() AS refreshed_at
FROM core.global_athletes ga
LEFT JOIN LATERAL (
    SELECT belt_level_id, exam_date
    FROM core.belt_history
    WHERE athlete_id = ga.id
    ORDER BY exam_date DESC, created_at DESC
    LIMIT 1
) bh ON true
LEFT JOIN core.belt_levels bl ON bl.id = bh.belt_level_id
WHERE ga.deleted_at IS NULL
  AND ga.status = 'active';

COMMENT ON MATERIALIZED VIEW core.mv_athlete_rankings IS 'Bảng xếp hạng VĐV toàn quốc + theo tỉnh. Refresh mỗi 15 phút.';

CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_rankings_athlete
    ON core.mv_athlete_rankings (athlete_id);

CREATE INDEX IF NOT EXISTS idx_mv_rankings_national
    ON core.mv_athlete_rankings (national_rank);

CREATE INDEX IF NOT EXISTS idx_mv_rankings_province
    ON core.mv_athlete_rankings (province, province_rank);

CREATE INDEX IF NOT EXISTS idx_mv_rankings_belt
    ON core.mv_athlete_rankings (belt_rank DESC);

-- ═══════════════════════════════════════════════
-- 3. FUNCTION: Refresh tất cả Materialized Views
-- ═══════════════════════════════════════════════

CREATE OR REPLACE FUNCTION core.refresh_all_matviews()
RETURNS VOID AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY core.mv_national_dashboard;
    REFRESH MATERIALIZED VIEW CONCURRENTLY core.mv_athlete_rankings;
    RAISE NOTICE 'All materialized views refreshed at %', now();
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION core.refresh_all_matviews IS 'Refresh tất cả materialized views. Gọi bằng pg_cron mỗi 15 phút.';

-- ═══════════════════════════════════════════════
-- 4. PG_CRON SCHEDULE (nếu Extension khả dụng)
-- Trên Supabase: pg_cron đã được cài sẵn
-- ═══════════════════════════════════════════════

DO $$
BEGIN
    -- Kiểm tra pg_cron có khả dụng không
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_cron') THEN
        -- Schedule refresh mỗi 15 phút
        PERFORM cron.schedule(
            'refresh_vct_dashboard',
            '*/15 * * * *',
            'SELECT core.refresh_all_matviews();'
        );
        RAISE NOTICE 'pg_cron scheduled: refresh_vct_dashboard every 15 minutes';
    ELSE
        RAISE NOTICE 'pg_cron not available. Manual refresh required: SELECT core.refresh_all_matviews();';
    END IF;
END $$;

-- ═══════════════════════════════════════════════
-- 5. COMPOSITE INDEXES cho Common Queries
-- ═══════════════════════════════════════════════

-- Tìm kiếm VĐV theo tỉnh + đai
CREATE INDEX IF NOT EXISTS idx_core_athletes_province_belt
    ON core.global_athletes (province, status)
    WHERE deleted_at IS NULL;

-- Event logs: Truy vấn sự kiện theo entity + time range
CREATE INDEX IF NOT EXISTS idx_core_events_entity_time
    ON core.event_logs (entity_type, entity_id, created_at DESC);

-- Belt history: Lấy đai mới nhất nhanh
CREATE INDEX IF NOT EXISTS idx_core_belt_latest
    ON core.belt_history (athlete_id, exam_date DESC, created_at DESC);

COMMIT;
