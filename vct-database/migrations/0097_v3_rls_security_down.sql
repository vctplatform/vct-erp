-- Down migration: 0097_v3_rls_security

BEGIN;

-- Drop all policies
DROP POLICY IF EXISTS sessions_self ON core.sessions;
DROP POLICY IF EXISTS events_no_delete ON core.event_logs;
DROP POLICY IF EXISTS events_no_update ON core.event_logs;
DROP POLICY IF EXISTS events_insert ON core.event_logs;
DROP POLICY IF EXISTS events_admin_read ON core.event_logs;
DROP POLICY IF EXISTS federations_admin_write ON core.federations;
DROP POLICY IF EXISTS federations_read_all ON core.federations;
DROP POLICY IF EXISTS belt_history_insert ON core.belt_history;
DROP POLICY IF EXISTS belt_history_read_all ON core.belt_history;
DROP POLICY IF EXISTS athletes_self_select ON core.global_athletes;
DROP POLICY IF EXISTS federation_athletes_read ON core.global_athletes;
DROP POLICY IF EXISTS athletes_admin_all ON core.global_athletes;
DROP POLICY IF EXISTS users_self_update ON core.users;
DROP POLICY IF EXISTS users_self_select ON core.users;
DROP POLICY IF EXISTS users_admin_all ON core.users;

-- Disable RLS
ALTER TABLE core.users DISABLE ROW LEVEL SECURITY;
ALTER TABLE core.global_athletes DISABLE ROW LEVEL SECURITY;
ALTER TABLE core.belt_history DISABLE ROW LEVEL SECURITY;
ALTER TABLE core.federations DISABLE ROW LEVEL SECURITY;
ALTER TABLE core.event_logs DISABLE ROW LEVEL SECURITY;
ALTER TABLE core.sessions DISABLE ROW LEVEL SECURITY;

-- Revoke grants
REVOKE ALL ON ALL TABLES IN SCHEMA core FROM vct_service;
REVOKE USAGE ON SCHEMA core FROM vct_service;

COMMIT;
