---
description: System health check workflow for VCT Platform
---

# /health-check - System Health Check

// turbo-all

## Quick Health Check

### Step 1: Check API health
```bash
curl -sf http://localhost:8080/health && echo "✅ API: Healthy" || echo "❌ API: Down"
```

### Step 2: Check API readiness
```bash
curl -sf http://localhost:8080/ready && echo "✅ Ready: All dependencies connected" || echo "❌ Ready: Some dependencies failing"
```

### Step 3: Check PostgreSQL
```bash
docker compose exec db pg_isready -U vct && echo "✅ PostgreSQL: Ready" || echo "❌ PostgreSQL: Down"
```

### Step 4: Check Redis
```bash
docker compose exec redis redis-cli ping | grep -q PONG && echo "✅ Redis: Ready" || echo "❌ Redis: Down"
```

### Step 5: Check frontend
```bash
curl -sf http://localhost:3000 > /dev/null && echo "✅ Frontend: Running" || echo "❌ Frontend: Down"
```

## Detailed Health Check

### Step 6: Check Docker containers
```bash
docker compose ps
```

### Step 7: Check resource usage
```bash
docker stats --no-stream
```

### Step 8: Check database connections
```bash
docker compose exec db psql -U vct -d vct_platform -c \
  "SELECT count(*) as active_connections FROM pg_stat_activity WHERE state = 'active';"
```

### Step 9: Check database size
```bash
docker compose exec db psql -U vct -d vct_platform -c \
  "SELECT pg_size_pretty(pg_database_size('vct_platform')) as db_size;"
```

### Step 10: Check disk usage
```bash
docker system df
```

### Step 11: Check recent error logs
```bash
docker compose logs --since 1h api 2>&1 | grep -i error | tail -20
```

### Step 12: Check API response times
```bash
# Quick performance check
for i in {1..5}; do
  time curl -sf http://localhost:8080/api/v1/athletes?page=1&per_page=1 > /dev/null
done
```

## Health Status Summary
```
┌──────────────────────────────────────┐
│ VCT Platform Health Check            │
├──────────────────────────────────────┤
│ API Server:     ✅ / ❌              │
│ PostgreSQL:     ✅ / ❌              │
│ Redis:          ✅ / ❌              │
│ Frontend:       ✅ / ❌              │
│ Monitoring:     ✅ / ❌              │
│ Disk Space:     XX% used             │
│ Memory:         XX% used             │
│ Active Conns:   XX / 200             │
└──────────────────────────────────────┘
```
