-- Down migration: 0092_v3_core_athletes_belts

BEGIN;

DROP TRIGGER IF EXISTS trg_athletes_updated_at ON core.global_athletes;

DROP TABLE IF EXISTS core.belt_history CASCADE;
DROP TABLE IF EXISTS core.belt_levels CASCADE;
DROP TABLE IF EXISTS core.global_athletes CASCADE;

COMMIT;
