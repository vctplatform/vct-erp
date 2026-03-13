---
name: dba
description: Database Administrator role - PostgreSQL 18+ management, Neon serverless, Supabase database, schema design, migrations, indexing, performance tuning, backup, and disaster recovery for VCT Platform.
---

# Database Administrator (DBA) - VCT Platform

## Role Overview
Manages all database aspects including schema design, migrations, performance tuning, backup/restore, and monitoring for PostgreSQL 18+ running on Neon (serverless) and Supabase.

## Technology Stack
- **Primary DB**: PostgreSQL 18+ (async I/O, JSON_TABLE, incremental backup, virtual columns)
- **Serverless Provider**: Neon (autoscaling, branching, instant provisioning, point-in-time restore)
- **BaaS Provider**: Supabase (RLS policies, pg_notify, Supabase migrations, real-time)
- **Driver**: pgx v5 (Go)
- **Migration Tool**: golang-migrate / goose / Supabase CLI (`supabase db push`)
- **Connection Pool**: pgxpool (Go), Neon connection pooler (PgBouncer built-in)
- **Monitoring**: pg_stat_statements, Neon dashboard, Supabase dashboard
- **Backup**: Neon (automatic PITR), Supabase (daily backups), pg_dump (manual)

## Core Patterns

### 1. Schema Design Standards

#### Naming Conventions
| Object | Convention | Example |
|--------|-----------|---------|
| Table | snake_case, plural | `athletes`, `club_members` |
| Column | snake_case | `first_name`, `created_at` |
| Primary Key | `id` | `id UUID DEFAULT gen_random_uuid()` |
| Foreign Key | `{table_singular}_id` | `club_id`, `athlete_id` |
| Index | `idx_{table}_{columns}` | `idx_athletes_email` |
| Unique | `uq_{table}_{columns}` | `uq_athletes_email` |
| Check | `ck_{table}_{column}` | `ck_athletes_gender` |
| Enum Type | singular | `gender_type`, `status_type` |

#### Standard Columns (ALL tables must include)
```sql
CREATE TABLE athletes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- ... business columns ...
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,  -- soft delete
    created_by  UUID REFERENCES auth.users(id),  -- Supabase auth user
    updated_by  UUID REFERENCES auth.users(id)
);

-- Auto-update updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_athletes_updated_at
    BEFORE UPDATE ON athletes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
```

### 2. PostgreSQL 18+ Features

#### JSON_TABLE (query JSON as relational)
```sql
-- Parse JSON tournament brackets directly
SELECT jt.*
FROM tournaments t,
  JSON_TABLE(t.bracket_data, '$.rounds[*]'
    COLUMNS (
      round_number INT PATH '$.round',
      match_count INT PATH '$.matches.size()',
      status TEXT PATH '$.status'
    )
  ) AS jt
WHERE t.id = $1;
```

#### Virtual Generated Columns
```sql
-- Full name as virtual computed column (no storage cost)
ALTER TABLE athletes
  ADD COLUMN full_name TEXT
  GENERATED ALWAYS AS (first_name || ' ' || last_name) VIRTUAL;
```

#### Incremental Backup
```sql
-- PG 18+ incremental backup (Neon handles this automatically)
-- Manual setup for self-hosted:
-- pg_basebackup --incremental=/path/to/last/backup -D /new/backup
```

### 3. Neon Serverless Configuration

#### Branching Strategy
```
main (production)
├── staging       → Auto-reset nightly from main
├── dev           → Development branch
├── pr-123        → Per-PR preview branch (auto-created by CI)
├── pr-456        → Per-PR preview branch (auto-deleted on merge)
└── migration-test → Test migrations before applying to main
```

#### Neon Connection String
```
# Production (pooled - for application)
postgres://user:pass@ep-cool-name-123456.us-east-2.aws.neon.tech/vct_platform?sslmode=require

# Direct (for migrations)
postgres://user:pass@ep-cool-name-123456.us-east-2.aws.neon.tech/vct_platform?sslmode=require&options=endpoint%3Dep-cool-name-123456
```

#### Neon Autoscaling Config
| Setting | Development | Staging | Production |
|---------|------------|---------|------------|
| Min CU | 0.25 | 0.5 | 1 |
| Max CU | 1 | 2 | 8 |
| Suspend after | 5 min | 10 min | Never |
| Autoscaling | On | On | On |

### 4. Supabase RLS (Row Level Security)

```sql
-- Enable RLS on all tables
ALTER TABLE athletes ENABLE ROW LEVEL SECURITY;

-- Policy: Users can read athletes in their organization
CREATE POLICY "org_read_athletes" ON athletes
    FOR SELECT
    USING (
        organization_id IN (
            SELECT organization_id FROM user_organizations
            WHERE user_id = auth.uid()
        )
    );

-- Policy: Club admins can insert athletes
CREATE POLICY "club_admin_insert_athletes" ON athletes
    FOR INSERT
    WITH CHECK (
        club_id IN (
            SELECT club_id FROM club_admins
            WHERE user_id = auth.uid()
        )
    );

-- Policy: Athletes can update own profile
CREATE POLICY "athlete_update_own" ON athletes
    FOR UPDATE
    USING (auth.uid()::text = user_id::text)
    WITH CHECK (auth.uid()::text = user_id::text);
```

### 5. Migration Standards

#### File Naming (Supabase CLI)
```
supabase/migrations/
├── 20260101000000_create_users.sql
├── 20260101000001_create_athletes.sql
├── 20260101000002_create_clubs.sql
├── 20260102000000_add_rls_policies.sql
└── 20260103000000_add_tournament_tables.sql
```

#### Migration Template
```sql
-- 20260101000001_create_athletes.sql

CREATE TYPE gender_type AS ENUM ('male', 'female');
CREATE TYPE athlete_status_type AS ENUM ('active', 'inactive', 'suspended');

CREATE TABLE athletes (
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
    user_id         UUID REFERENCES auth.users(id),  -- Supabase auth
    organization_id UUID REFERENCES organizations(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    created_by      UUID REFERENCES auth.users(id),
    updated_by      UUID REFERENCES auth.users(id)
);

CREATE UNIQUE INDEX uq_athletes_email ON athletes(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_athletes_club_id ON athletes(club_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_athletes_status ON athletes(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_athletes_name ON athletes(last_name, first_name);
CREATE INDEX idx_athletes_org ON athletes(organization_id) WHERE deleted_at IS NULL;

-- Enable RLS
ALTER TABLE athletes ENABLE ROW LEVEL SECURITY;

CREATE TRIGGER tr_athletes_updated_at
    BEFORE UPDATE ON athletes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
```

### 6. Indexing Strategy

| Index Type | When to Use | Example |
|-----------|------------|---------|
| B-tree (default) | Equality, range, sorting | `CREATE INDEX idx_athletes_dob ON athletes(date_of_birth)` |
| Hash | Equality only | `CREATE INDEX idx_users_email_hash ON users USING hash(email)` |
| GIN | Full-text search, JSONB | `CREATE INDEX idx_athletes_search ON athletes USING gin(search_vector)` |
| GiST | Geometric, range types | `CREATE INDEX idx_events_daterange ON events USING gist(date_range)` |
| Partial | Filtered queries | `CREATE INDEX idx_active_athletes ON athletes(status) WHERE deleted_at IS NULL` |
| Covering | Include extra columns | `CREATE INDEX idx_athletes_club INCLUDE (first_name, last_name) ON athletes(club_id)` |

### 7. Performance Tuning

#### Neon-Specific Tuning
```
- Neon manages most PostgreSQL settings automatically
- Focus on query optimization, indexing, and connection management
- Use Neon's built-in connection pooler (PgBouncer)
- Monitor via Neon dashboard: query performance, storage, compute hours
```

#### Slow Query Investigation Checklist
1. `EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)` on the query
2. Check for sequential scans on large tables
3. Verify indexes exist for WHERE/JOIN columns
4. Check Neon dashboard for compute scaling events
5. Verify table statistics are up-to-date: `ANALYZE table_name`
6. Check connection pool utilization (Neon pooler)
7. Consider partitioning for tables > 10M rows

### 8. Backup Strategy

| Method | Provider | Frequency | Retention |
|--------|----------|-----------|-----------|
| PITR (Point-in-Time) | Neon | Continuous (WAL) | 7-30 days (plan dependent) |
| Neon Branching | Neon | On-demand | Unlimited (manual cleanup) |
| Daily Backup | Supabase | Daily | 7 days (Pro: 30 days) |
| pg_dump (manual) | Self | Weekly | 90 days (S3) |

#### Recovery with Neon Branching
```bash
# Create a branch from any point in time
neonctl branches create \
  --project-id <project_id> \
  --name recovery-branch \
  --parent main \
  --set-as-default=false

# Restore: promote recovery branch to main
neonctl branches set-primary \
  --project-id <project_id> \
  --branch recovery-branch
```

### 9. Monitoring Queries

```sql
-- Active connections
SELECT count(*) FROM pg_stat_activity WHERE state = 'active';

-- Slow queries (real-time)
SELECT pid, now() - pg_stat_activity.query_start AS duration, query
FROM pg_stat_activity
WHERE state != 'idle' AND query_start < now() - interval '5 seconds'
ORDER BY duration DESC;

-- Table sizes
SELECT relname, pg_size_pretty(pg_total_relation_size(relid))
FROM pg_catalog.pg_statio_user_tables
ORDER BY pg_total_relation_size(relid) DESC;

-- Index usage
SELECT indexrelname, idx_scan, idx_tup_read, idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- Unused indexes
SELECT indexrelname FROM pg_stat_user_indexes WHERE idx_scan = 0;

-- Cache hit ratio (should be > 99%)
SELECT sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)) AS ratio
FROM pg_statio_user_tables;
```

### 10. DBA Checklist
- [ ] All tables have `updated_at` trigger
- [ ] Soft delete pattern used (`deleted_at` column)
- [ ] Partial indexes exclude soft-deleted rows
- [ ] Foreign keys have corresponding indexes
- [ ] RLS policies enabled on all tables (Supabase)
- [ ] Migrations tested on Neon branch before applying to main
- [ ] Connection pooling via Neon built-in pooler
- [ ] PITR enabled on Neon (automatic)
- [ ] Supabase RLS policies reviewed for security
- [ ] `pg_stat_statements` enabled for query analysis
- [ ] Neon compute autoscaling configured per environment
