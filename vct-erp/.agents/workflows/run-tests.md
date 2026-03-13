---
description: Run all test suites for VCT Platform (unit, integration, E2E, performance)
---

# /run-tests - Execute Test Suites

// turbo-all

## Quick Test (before commit)

### Step 1: Run backend unit tests
```bash
cd backend
go test -race -count=1 -short ./...
```

### Step 2: Run frontend unit tests
```bash
cd frontend
npx vitest run
```

---

## Full Test Suite (before merge/deploy)

### Step 3: Backend - Full test suite with coverage
```bash
cd backend
go test -race -count=1 -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out | tail -1
# Target: ≥ 80% total coverage
```

### Step 4: Backend - Generate coverage HTML report
```bash
cd backend
go tool cover -html=coverage.out -o coverage.html
echo "Coverage report: backend/coverage.html"
```

### Step 5: Frontend - Full test suite with coverage
```bash
cd frontend
npx vitest run --coverage
```

### Step 6: Backend - Lint check
```bash
cd backend
golangci-lint run --timeout 5m ./...
```

### Step 7: Frontend - Lint and type check
```bash
cd frontend
npm run lint
npx tsc --noEmit
```

### Step 8: Frontend - Build check
```bash
cd frontend
npm run build
```

---

## Integration Tests

### Step 9: Start test infrastructure
```bash
docker compose -f docker-compose.test.yml up -d db-test redis-test
sleep 5
```

### Step 10: Run integration tests
```bash
cd backend
DATABASE_URL="postgres://vct:vct_secret@localhost:5433/vct_test?sslmode=disable" \
go test -race -tags=integration -count=1 ./...
```

### Step 11: Cleanup test infrastructure
```bash
docker compose -f docker-compose.test.yml down -v
```

---

## E2E Tests

### Step 12: Start full application stack
```bash
docker compose up -d
sleep 10
```

### Step 13: Run E2E tests
```bash
cd frontend
npx playwright test --reporter=html
echo "E2E report: frontend/playwright-report/index.html"
```

### Step 14: Cleanup
```bash
docker compose down
```

---

## Performance Tests

### Step 15: Run load test (requires running API)
```bash
k6 run --vus 50 --duration 30s tests/performance/athlete-load.js
```

---

## Test Summary
```
┌──────────────────────────────────────────┐
│ Test Suite         │ Command             │
├──────────────────────────────────────────┤
│ Backend Unit       │ go test ./...       │
│ Backend Coverage   │ go test -cover ./...│
│ Backend Lint       │ golangci-lint run   │
│ Frontend Unit      │ vitest run          │
│ Frontend Lint      │ npm run lint        │
│ Frontend Types     │ tsc --noEmit        │
│ Frontend Build     │ npm run build       │
│ Integration        │ go test -tags=int..│
│ E2E                │ playwright test     │
│ Performance        │ k6 run             │
└──────────────────────────────────────────┘
```
