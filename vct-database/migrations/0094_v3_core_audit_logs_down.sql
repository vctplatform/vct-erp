-- Down migration: 0094_v3_core_audit_logs

BEGIN;

DROP TRIGGER IF EXISTS trg_audit_federations ON core.federations;
DROP TRIGGER IF EXISTS trg_audit_belts ON core.belt_history;
DROP TRIGGER IF EXISTS trg_audit_athletes ON core.global_athletes;
DROP TRIGGER IF EXISTS trg_audit_users ON core.users;

DROP FUNCTION IF EXISTS core.fn_audit_trigger();
DROP TABLE IF EXISTS core.audit_logs CASCADE;

COMMIT;
