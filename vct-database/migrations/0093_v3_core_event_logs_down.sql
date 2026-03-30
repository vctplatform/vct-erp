-- Down migration: 0093_v3_core_event_logs

BEGIN;

DROP FUNCTION IF EXISTS core.log_event;
DROP TABLE IF EXISTS core.event_logs CASCADE;

COMMIT;
