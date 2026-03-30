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
