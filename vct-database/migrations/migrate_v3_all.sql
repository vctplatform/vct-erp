-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Master Migration Script: Database v3.0 FINAL
-- Chạy tất cả migrations 0091-0098 theo thứ tự
-- 
-- HƯỚNG DẪN CHẠY:
--   psql "postgresql://postgres.xxx:password@host:5432/postgres" -f migrate_v3_all.sql
--
-- Hoặc chạy từng file:
--   psql $DATABASE_URL -f 0091_v3_create_core_schema.sql
--   psql $DATABASE_URL -f 0092_v3_core_athletes_belts.sql
--   ... (theo thứ tự)
--
-- Approved by: Chairman Hoàng Bá Tùng — 30/03/2026
-- ══════════════════════════════════════════════════════════════════

-- Phase 1: Core Foundation
\echo '━━━ Phase 1: Core Schema Foundation ━━━'
\i 0091_v3_create_core_schema.sql
\echo '✅ 0091 — Core schema + users + federations + sessions'

\i 0092_v3_core_athletes_belts.sql
\echo '✅ 0092 — Global athletes + belt levels + belt history'

\i 0093_v3_core_event_logs.sql
\echo '✅ 0093 — Event logs (Append-Only)'

\i 0094_v3_core_audit_logs.sql
\echo '✅ 0094 — Audit logs (Sổ Vàng) + triggers'

-- Phase 2: Federation
\echo '━━━ Phase 2: Federation Schema ━━━'
\i 0095_v3_federation_template.sql
\echo '✅ 0095 — Federation template + fed_lamdong example'

-- Phase 3: Tournament
\echo '━━━ Phase 3: Tournament Schema ━━━'
\i 0096_v3_tournament_template.sql
\echo '✅ 0096 — Tournament template + helper functions'

-- Phase 4: Security
\echo '━━━ Phase 4: RLS & Security ━━━'
\i 0097_v3_rls_security.sql
\echo '✅ 0097 — RLS policies for all core tables'

-- Phase 5: Performance
\echo '━━━ Phase 5: Indexes & Materialized Views ━━━'
\i 0098_v3_indexes_matviews.sql
\echo '✅ 0098 — Dashboard CEO + Athlete Rankings + pg_cron'

-- Phase 6: PostgreSQL 18.3 Features
\echo '━━━ Phase 6: PostgreSQL 18.3 Features ━━━'
\i 0099_v3_pg18_features.sql
\echo '✅ 0099 — PostgreSQL 18.3 features initialized'

-- Phase 7: Hierarchy & Seed
\echo '━━━ Phase 7: Hierarchical Federations ━━━'
\i 0100_v3_federations_hierarchy_and_seed.sql
\echo '✅ 0100 — 35 Federations seeded with parent_id'

\echo ''
\echo '══════════════════════════════════════════════'
\echo '  ✅ VCT DATABASE v3.0 MIGRATION COMPLETE!'
\echo '══════════════════════════════════════════════'
\echo ''
\echo 'Schemas created:'
\echo '  • core          — Global Hub (Users, Athletes, Belts, Events, Audit)'
\echo '  • fed_lamdong   — Federation example (Lâm Đồng)'
\echo ''
\echo 'Functions available:'
\echo '  • core.create_federation_schema(province_code)'
\echo '  • core.create_tournament_schema(year, scope, seq)'
\echo '  • core.export_tournament_roster(schema, athlete_ids)'
\echo '  • core.lock_tournament_roster(schema)'
\echo '  • core.archive_tournament_schema(schema)'
\echo '  • core.log_event(type, entity_type, entity_id, payload)'
\echo '  • core.refresh_all_matviews()'
