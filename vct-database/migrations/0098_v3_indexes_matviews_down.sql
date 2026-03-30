-- Down migration: 0098_v3_indexes_matviews

BEGIN;

-- Cancel cron job if exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_cron') THEN
        PERFORM cron.unschedule('refresh_vct_dashboard');
    END IF;
EXCEPTION WHEN OTHERS THEN
    NULL;
END $$;

DROP FUNCTION IF EXISTS core.refresh_all_matviews();
DROP MATERIALIZED VIEW IF EXISTS core.mv_athlete_rankings;
DROP MATERIALIZED VIEW IF EXISTS core.mv_national_dashboard;

DROP INDEX IF EXISTS core.idx_core_belt_latest;
DROP INDEX IF EXISTS core.idx_core_events_entity_time;
DROP INDEX IF EXISTS core.idx_core_athletes_province_belt;

COMMIT;
