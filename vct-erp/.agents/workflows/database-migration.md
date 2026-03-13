---
description: Create, review, and apply database migrations for VCT Platform (Neon + Supabase)
---

# /database-migration - Database Migration Workflow

## When to Use
- Adding/modifying database tables, columns, or indexes
- Changing constraints, triggers, or functions
- Adding/updating RLS policies (Supabase)
- Data transformations or seed data

## Steps

// turbo-all

### Step 1: Create a Neon branch for testing
```bash
# Create isolated branch to test migrations
neonctl branches create \
  --project-id ${NEON_PROJECT_ID} \
  --name migration-test \
  --parent main

# Get branch connection string
neonctl connection-string migration-test --project-id ${NEON_PROJECT_ID}
```

### Step 2: Create migration file (Supabase CLI)
```bash
supabase migration new {migration_name}
# Example: supabase migration new create_athletes
# Creates: supabase/migrations/{timestamp}_{name}.sql
```

### Step 3: Write migration SQL
```sql
-- supabase/migrations/{timestamp}_create_athletes.sql

CREATE TYPE gender_type AS ENUM ('male', 'female');
CREATE TYPE athlete_status_type AS ENUM ('active', 'inactive', 'suspended');

CREATE TABLE IF NOT EXISTS athletes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    full_name       TEXT GENERATED ALWAYS AS (first_name || ' ' || last_name) VIRTUAL,
    date_of_birth   DATE NOT NULL,
    gender          gender_type NOT NULL,
    email           VARCHAR(255) NOT NULL,
    phone           VARCHAR(20),
    club_id         UUID REFERENCES clubs(id),
    status          athlete_status_type NOT NULL DEFAULT 'active',
    user_id         UUID REFERENCES auth.users(id),
    organization_id UUID REFERENCES organizations(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    created_by      UUID REFERENCES auth.users(id),
    updated_by      UUID REFERENCES auth.users(id)
);

-- Indexes
CREATE UNIQUE INDEX IF NOT EXISTS uq_athletes_email ON athletes(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_athletes_club_id ON athletes(club_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_athletes_org ON athletes(organization_id) WHERE deleted_at IS NULL;

-- Enable RLS
ALTER TABLE athletes ENABLE ROW LEVEL SECURITY;

-- RLS policies
CREATE POLICY "org_read_athletes" ON athletes
    FOR SELECT USING (
        organization_id IN (
            SELECT organization_id FROM user_organizations
            WHERE user_id = auth.uid()
        )
    );

-- Updated_at trigger
CREATE TRIGGER tr_athletes_updated_at
    BEFORE UPDATE ON athletes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
```

### Step 4: Test migration on Neon branch
```bash
# Apply to test branch
supabase db push --db-url "${NEON_MIGRATION_BRANCH_URL}"

# Or with golang-migrate:
migrate -path supabase/migrations -database "${NEON_MIGRATION_BRANCH_URL}" up

# Verify
psql "${NEON_MIGRATION_BRANCH_URL}" -c "\d athletes"
```

### Step 5: Test rollback (if using golang-migrate)
```bash
migrate -path supabase/migrations -database "${NEON_MIGRATION_BRANCH_URL}" down 1
migrate -path supabase/migrations -database "${NEON_MIGRATION_BRANCH_URL}" up
```

### Step 6: Apply to Neon main branch (after PR merge)
```bash
# Apply to production/main branch
supabase db push --db-url "${NEON_DATABASE_URL}"

# Or sync with Supabase cloud project
supabase db push --linked
```

### Step 7: Cleanup test branch
```bash
neonctl branches delete migration-test --project-id ${NEON_PROJECT_ID}
```

## Migration Rules

| Rule | Description |
|------|------------|
| **Idempotent** | Use `IF NOT EXISTS` / `IF EXISTS` |
| **Transactional** | Neon supports transactional DDL |
| **RLS Included** | Always include RLS policies for new tables |
| **Test on Branch** | Always test on Neon branch before main |
| **Non-destructive** | Never `DROP COLUMN` without data backup |
| **Small** | One logical change per migration |
| **Named** | Descriptive names: `create_athletes`, `add_club_id_to_athletes` |

## Common Patterns

```sql
-- Add column
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS phone VARCHAR(20);

-- Add NOT NULL column with default
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active';

-- Create enum type (safe)
DO $$ BEGIN
    CREATE TYPE gender_type AS ENUM ('male', 'female');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Add RLS policy
ALTER TABLE {table} ENABLE ROW LEVEL SECURITY;
CREATE POLICY "{name}" ON {table} FOR SELECT
    USING (organization_id IN (
        SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
    ));
```
