---
description: Quy trình testing — Unit, Integration, E2E Strategy
---

# /testing — Testing Workflow

## STRATEGY
```
Testing Pyramid:
├── Unit Tests (70%): Fast, isolated, mock dependencies
├── Integration Tests (20%): Real DB, real Redis
└── E2E Tests (10%): Full stack, critical paths only
```

## BƯỚC 1: UNIT TESTS (Skill: backend-engine)
// turbo
1. Table-driven tests for use cases
2. Mock repository interface
3. Test happy path + error paths
4. Test edge cases (empty, nil, boundary)
5. Run: `go test ./internal/modules/...`

## BƯỚC 2: INTEGRATION TESTS
// turbo
1. Use test database (Docker)
2. Test repository with real PostgreSQL
3. Test full handler flow (httptest)
4. Verify DB state after operations

## BƯỚC 3: FRONTEND TESTS
// turbo
1. Component tests (render, interaction)
2. API integration tests (mock server)
3. Visual regression (if applicable)

## BƯỚC 4: RUN & REPORT
// turbo
1. Run full suite: `go test ./... -cover`
2. Check coverage: target > 80%
3. Report failures with context

// turbo-all

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
