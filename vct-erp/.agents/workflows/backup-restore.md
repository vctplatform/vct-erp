---
description: Database backup and restore procedures for PostgreSQL
---

# /backup-restore - Database Backup & Restore

## Backup Strategy (Neon + Supabase)

| Method | Provider | Frequency | Retention |
|--------|----------|-----------|-----------|
| Point-in-Time Recovery | Neon | Continuous (WAL) | 7-30 days (plan dependent) |
| Neon Branching | Neon | On-demand | Until manually deleted |
| Supabase Daily Backup | Supabase | Daily | 7 days (Pro: 30 days) |
| pg_dump (manual) | Self | Weekly | 90 days (S3) |

## Neon PITR (Primary Method)

### Step 1: Restore to a point in time via branching
```bash
# Create a branch from a specific point in time
neonctl branches create \
  --project-id ${NEON_PROJECT_ID} \
  --name recovery-$(date +%Y%m%d) \
  --parent main \
  --set-as-default=false

# Verify data on recovery branch
psql "$(neonctl connection-string recovery-$(date +%Y%m%d) --project-id ${NEON_PROJECT_ID})" \
  -c "SELECT count(*) FROM athletes;"
```

### Step 2: Promote recovery branch (if needed)
```bash
# ⚠️ This makes the recovery branch the new main
neonctl branches set-primary \
  --project-id ${NEON_PROJECT_ID} \
  --branch recovery-$(date +%Y%m%d)
```

## Manual Backup (pg_dump from Neon)

### Step 3: Manual logical backup
```bash
# Full database backup from Neon
pg_dump -Fc \
  "${NEON_DATABASE_URL}" \
  -f "vct_platform_$(date +%Y%m%d_%H%M%S).dump"
```

### Step 4: Backup specific tables
```bash
pg_dump -Fc -t athletes "${NEON_DATABASE_URL}" -f athletes_backup.dump
pg_dump -s "${NEON_DATABASE_URL}" -f schema_only.sql
```

### Step 5: Verify backup integrity
```bash
pg_restore -l vct_platform_backup.dump | head -20
```

### Step 6: Upload to remote storage
```bash
aws s3 cp vct_platform_backup.dump \
  s3://vct-backups/postgres/$(date +%Y/%m/%d)/ \
  --storage-class STANDARD_IA
```

## Restore Steps

### Step 7: Restore to a new Neon branch
```bash
# Create fresh branch
neonctl branches create \
  --project-id ${NEON_PROJECT_ID} \
  --name restore-test

# Restore from pg_dump
pg_restore -d "$(neonctl connection-string restore-test --project-id ${NEON_PROJECT_ID})" \
  -j 4 vct_platform_backup.dump
```

## Supabase Backups

### Step 8: Check Supabase backup status
```bash
# Supabase Pro plans include daily backups
# Access via: Supabase Dashboard → Project → Database → Backups
# Point-in-Time Recovery available on Pro plans
```

## Emergency Checklist
- [ ] Identify the issue (data corruption? accidental delete?)
- [ ] Determine recovery target time
- [ ] Notify team of downtime
- [ ] Create Neon branch from target point in time
- [ ] Verify data integrity on recovery branch
- [ ] If correct, promote recovery branch
- [ ] Verify application health
- [ ] Write incident report
