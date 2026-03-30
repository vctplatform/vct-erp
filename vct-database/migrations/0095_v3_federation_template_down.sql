-- Down migration: 0095_v3_federation_template

BEGIN;

DROP SCHEMA IF EXISTS fed_lamdong CASCADE;
DROP FUNCTION IF EXISTS core.create_federation_schema(TEXT);

DELETE FROM core.federations WHERE province_code = 'lamdong';

COMMIT;
