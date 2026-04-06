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

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
