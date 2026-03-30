-- Down migration: 0091_v3_create_core_schema

BEGIN;

DROP TRIGGER IF EXISTS trg_federations_updated_at ON core.federations;
DROP TRIGGER IF EXISTS trg_users_updated_at ON core.users;
DROP FUNCTION IF EXISTS core.trigger_set_updated_at();

DROP TABLE IF EXISTS core.sessions CASCADE;
DROP TABLE IF EXISTS core.federations CASCADE;
DROP TABLE IF EXISTS core.users CASCADE;

DROP SCHEMA IF EXISTS core CASCADE;

COMMIT;
