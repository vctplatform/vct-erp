---
description: Performance profiling workflow - Go pprof, Lighthouse, Neon query analysis, and k6 load testing
---

# /performance-profile - Performance Profiling Workflow

// turbo-all

## When to Use
- Before release: baseline performance check
- After performance regression detected
- Tournament day preparation (peak traffic)
- After database schema changes

## Steps

### Step 1: Backend profiling (Go pprof)
```bash
cd backend

# Enable pprof endpoint (add to main.go if not present):
# import _ "net/http/pprof"
# go func() { http.ListenAndServe(":6060", nil) }()

# CPU profile (30 seconds)
go tool pprof -http=:8090 http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profile
go tool pprof -http=:8090 http://localhost:6060/debug/pprof/heap

# Goroutine analysis
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Trace (5 seconds)
curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5
go tool trace trace.out
```

### Step 2: Go benchmarks
```bash
cd backend

# Run all benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkAthleteList -benchmem -count=5 ./internal/modules/athlete/...

# Compare benchmarks (before/after)
go test -bench=. -benchmem ./... > bench_before.txt
# ... make changes ...
go test -bench=. -benchmem ./... > bench_after.txt
go install golang.org/x/perf/cmd/benchstat@latest
benchstat bench_before.txt bench_after.txt
```

### Step 3: Database query analysis (Neon)
```sql
-- Enable pg_stat_statements (usually enabled on Neon)
-- Top 10 slowest queries
SELECT
    round(mean_exec_time::numeric, 2) AS avg_ms,
    round(total_exec_time::numeric, 2) AS total_ms,
    calls,
    round((100 * total_exec_time / sum(total_exec_time) OVER ())::numeric, 2) AS pct,
    query
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Top 10 most called queries
SELECT calls, mean_exec_time, query
FROM pg_stat_statements
ORDER BY calls DESC
LIMIT 10;

-- Missing indexes (sequential scans on large tables)
SELECT relname, seq_scan, seq_tup_read, idx_scan, idx_tup_fetch
FROM pg_stat_user_tables
WHERE seq_scan > 100
ORDER BY seq_tup_read DESC
LIMIT 10;

-- Table bloat check
SELECT
    schemaname, tablename,
    pg_size_pretty(pg_total_relation_size(schemaname || '.' || tablename)) AS total_size,
    pg_size_pretty(pg_relation_size(schemaname || '.' || tablename)) AS table_size,
    pg_size_pretty(pg_indexes_size(schemaname || '.' || tablename::regclass)) AS index_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname || '.' || tablename) DESC;

-- Cache hit ratio (should be > 99%)
SELECT
    sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)) AS ratio
FROM pg_statio_user_tables;
```

### Step 4: EXPLAIN ANALYZE for slow queries
```sql
-- Analyze specific query
EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)
SELECT a.*, c.name AS club_name
FROM athletes a
LEFT JOIN clubs c ON c.id = a.club_id
WHERE a.organization_id = '...'
AND a.deleted_at IS NULL
ORDER BY a.last_name
LIMIT 20 OFFSET 0;
```

### Step 5: Frontend performance (Lighthouse)
```bash
cd apps/web

# Run Lighthouse CI
npx lighthouse https://vct-platform.com \
  --output=json --output=html \
  --output-path=./lighthouse-report

# Key metrics to check:
# - First Contentful Paint (FCP) < 1.8s
# - Largest Contentful Paint (LCP) < 2.5s
# - Cumulative Layout Shift (CLS) < 0.1
# - Total Blocking Time (TBT) < 200ms
# - Time to Interactive (TTI) < 3.8s
```

### Step 6: Bundle size analysis
```bash
cd apps/web

# Vite bundle analysis
npx vite-bundle-visualizer

# Check for large dependencies
du -sh node_modules/* | sort -rh | head -20
```

### Step 7: Load testing (k6)
```bash
# Install k6
# winget install k6

# Run load test
k6 run tests/performance/athlete-load.js \
  --env API_URL=https://staging-api.vct-platform.com

# Tournament day simulation (high concurrency)
k6 run tests/performance/tournament-peak.js \
  --env API_URL=https://staging-api.vct-platform.com \
  --vus 500 --duration 5m
```

### Step 8: Neon compute monitoring
```bash
echo "Check in Neon Dashboard:"
echo "  1. Compute hours usage"
echo "  2. Autoscaling events (did CU scale up?)"
echo "  3. Connection count during peak"
echo "  4. Storage growth rate"
echo ""
echo "Neon Console: https://console.neon.tech"
```

### Step 9: Generate performance report
```bash
echo "=== Performance Report ==="
echo "Date: $(date)"
echo ""
echo "Backend:"
echo "  API P95 latency: [from Grafana]"
echo "  API P99 latency: [from Grafana]"
echo "  Error rate: [from Grafana]"
echo ""
echo "Database:"
echo "  Slowest query avg: [from pg_stat_statements]"
echo "  Cache hit ratio: [from query]"
echo "  Neon CU usage: [from dashboard]"
echo ""
echo "Frontend:"
echo "  LCP: [from Lighthouse]"
echo "  FCP: [from Lighthouse]"
echo "  Bundle size: [from vite-bundle-visualizer]"
echo ""
echo "Load Test:"
echo "  Max concurrent users: [from k6]"
echo "  P95 under load: [from k6]"
echo "  Error rate under load: [from k6]"
```

## Performance Targets

| Metric | Target | Alert Threshold |
|--------|--------|----------------|
| API P95 | < 200ms | > 500ms |
| API P99 | < 500ms | > 1s |
| DB query avg | < 50ms | > 200ms |
| FCP (web) | < 1.8s | > 3s |
| LCP (web) | < 2.5s | > 4s |
| Bundle size (web) | < 500KB gzip | > 1MB |
| App launch (mobile) | < 2s | > 4s |
| Neon cache hit | > 99% | < 95% |
