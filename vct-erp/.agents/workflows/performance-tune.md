---
description: Quy trình tối ưu hiệu năng — Profile, Identify, Optimize
---

# /performance-tune — Performance Optimization Workflow

## BƯỚC 1: PROFILE
// turbo
1. Identify bottleneck (slow query? CPU? memory?)
2. Backend: `go tool pprof`, query EXPLAIN ANALYZE
3. Frontend: Lighthouse, React DevTools Profiler
4. Database: pg_stat_statements, slow query log

## BƯỚC 2: IDENTIFY
// turbo
1. Top 5 slowest queries
2. Top 5 slowest API endpoints
3. Top 5 largest frontend bundles
4. Memory leaks or goroutine leaks

## BƯỚC 3: OPTIMIZE
// turbo
1. Database: Add indexes, rewrite queries, use CTEs
2. Backend: Connection pooling, caching, pagination
3. Frontend: Code splitting, lazy loading, memo
4. Infra: CDN, compression, HTTP/2

## BƯỚC 4: VERIFY
// turbo
1. Benchmark before vs after
2. Verify no regression
3. Document optimization and rationale

// turbo-all
