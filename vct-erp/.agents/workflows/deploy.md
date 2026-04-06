---
description: Quy trình deployment — Build, Stage, Test, Production
---

# /deploy — Deployment Workflow

## BƯỚC 1: PRE-DEPLOY CHECK
1. All tests pass (unit + integration)
2. No critical security vulnerabilities
3. Migration tested on staging
4. Changelog documented

## BƯỚC 2: BUILD
// turbo
1. Backend: `go build -o finance-api ./cmd/api/`
2. Frontend: `npm run build`
3. Docker: `docker build -t vct-erp-api .`
4. Verify build artifacts

## BƯỚC 3: STAGING DEPLOY
1. Apply database migrations
2. Deploy backend containers
3. Deploy frontend
4. Run smoke tests
5. Verify health check endpoints

## BƯỚC 4: PRODUCTION DEPLOY (Cấp 3 — Chairman approve)
1. Blue/Green or Canary deployment
2. Apply migrations (with backup)
3. Deploy with health checks
4. Monitor error rates (< 1%)
5. Rollback if issues detected

## BƯỚC 5: POST-DEPLOY
// turbo
1. Verify health check: GET /healthz
2. Check monitoring dashboards
3. Confirm key user flows work
4. Update deployment log

// turbo-all

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
