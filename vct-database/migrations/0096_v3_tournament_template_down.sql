-- Down migration: 0096_v3_tournament_template

BEGIN;

DROP FUNCTION IF EXISTS core.lock_tournament_roster(TEXT);
DROP FUNCTION IF EXISTS core.archive_tournament_schema(TEXT);
DROP FUNCTION IF EXISTS core.export_tournament_roster(TEXT, UUID[]);
DROP FUNCTION IF EXISTS core.create_tournament_schema(TEXT, TEXT, TEXT);

COMMIT;
