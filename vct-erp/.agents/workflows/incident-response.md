---
description: Production incident response playbook - triage, resolve, and postmortem
---

# /incident-response - Production Incident Playbook

## Severity Levels

| Level | Description | Response Time | Example |
|-------|------------|---------------|---------|
| **SEV1** (Critical) | Service down, data loss | < 15 min | API 500 for all users, DB corruption |
| **SEV2** (High) | Major feature broken | < 1 hour | Auth not working, scoring broken during tournament |
| **SEV3** (Medium) | Degraded performance | < 4 hours | Slow queries, intermittent errors |
| **SEV4** (Low) | Minor issue | Next business day | UI glitch, non-critical feature broken |

## Incident Response Steps

### Phase 1: Detect & Triage (0-15 min)

#### Step 1: Confirm the incident
```bash
# Check API health
curl -f https://api.vct-platform.com/health
curl -f https://api.vct-platform.com/ready

# Check error rates (Grafana/Prometheus)
echo "Grafana: http://monitoring.vct-platform.com:3001"

# Check Neon status
echo "Neon Status: https://neonstatus.com"

# Check Supabase status
echo "Supabase Status: https://status.supabase.com"
```

#### Step 2: Classify severity
```markdown
- Who is affected? (all users / subset / single user)
- What functionality is broken? (core / peripheral)
- Is data at risk? (data loss / corruption / read-only)
- Is it tournament day? (peak traffic = auto-escalate)
```

#### Step 3: Notify team
```bash
# SEV1/SEV2: Immediate notification
echo "🚨 INCIDENT: [Brief description]"
echo "Severity: SEV[X]"
echo "Impact: [Who and what is affected]"
echo "Incident Commander: [Name]"
echo "Channel: #vct-incidents"
```

### Phase 2: Investigate & Mitigate (15-60 min)

#### Step 4: Check application logs
```bash
# Kubernetes logs
kubectl logs -n vct-production -l app=vct-api --tail=100 --since=10m

# Docker logs
docker compose logs api --tail=100 --since=10m
```

#### Step 5: Check database health (Neon)
```sql
-- Active connections
SELECT count(*) FROM pg_stat_activity WHERE state = 'active';

-- Long-running queries
SELECT pid, now() - query_start AS duration, query
FROM pg_stat_activity
WHERE state != 'idle' AND query_start < now() - interval '30 seconds'
ORDER BY duration DESC;

-- Locks
SELECT pid, relation::regclass, mode, granted
FROM pg_locks WHERE NOT granted;
```

#### Step 6: Check Neon compute status
```bash
# Check compute scaling
echo "Neon Dashboard → Project → Monitoring"
echo "  - Check: compute scaling events"
echo "  - Check: connection count"
echo "  - Check: storage usage"
```

#### Step 7: Quick mitigation
```bash
# Option A: Rollback application
kubectl rollout undo deployment/vct-api --namespace=vct-production

# Option B: Rollback database (Neon PITR)
neonctl branches create \
  --project-id ${NEON_PROJECT_ID} \
  --name recovery-$(date +%Y%m%d-%H%M) \
  --parent main

# Option C: Scale up Neon compute
echo "Neon Dashboard → Branch Settings → Increase Max CU"

# Option D: Kill problematic queries
psql "${NEON_DATABASE_URL}" -c "SELECT pg_terminate_backend(PID);"
```

### Phase 3: Resolve (30 min - 4 hours)

#### Step 8: Apply fix
```bash
# Create hotfix branch
git checkout -b hotfix/incident-$(date +%Y%m%d) main

# Apply fix, test, merge (see /hotfix workflow)
```

#### Step 9: Verify resolution
```bash
# Health check
curl -f https://api.vct-platform.com/health

# Smoke test critical endpoints
curl -f https://api.vct-platform.com/api/v1/athletes?page=1&per_page=1
curl -f https://vct-platform.com

# Check error rates returning to normal
echo "Monitor Grafana for 15 minutes"
```

### Phase 4: Postmortem (within 48 hours)

#### Step 10: Create postmortem document
```markdown
## Incident Postmortem: [Title]

### Summary
- **Date/Time**: [Start] - [End] (Duration: [X] minutes)
- **Severity**: SEV[X]
- **Impact**: [Who/what was affected, estimated user impact]
- **Incident Commander**: [Name]

### Timeline
| Time | Event |
|------|-------|
| HH:MM | Alert triggered / User report received |
| HH:MM | Incident confirmed, team notified |
| HH:MM | Root cause identified |
| HH:MM | Mitigation applied |
| HH:MM | Service restored |

### Root Cause
[Detailed technical explanation]

### Resolution
[What was done to fix it]

### Action Items
| # | Action | Owner | Due Date | Status |
|---|--------|-------|----------|--------|
| 1 | [Preventive measure] | @dev | [date] | TODO |
| 2 | [Monitoring improvement] | @devops | [date] | TODO |
| 3 | [Process change] | @pm | [date] | TODO |

### Lessons Learned
- What went well?
- What could be improved?
- Where did we get lucky?
```

## Quick Reference

| Scenario | First Action |
|----------|-------------|
| API returning 500s | Check logs → rollback if needed |
| Database connection errors | Check Neon dashboard → compute scaling |
| Auth not working | Check Supabase status → verify JWT config |
| Slow responses | Check long-running queries → kill / optimize |
| Data inconsistency | Create Neon recovery branch → investigate |
| DDoS suspected | Enable rate limiting → check Cloudflare |
