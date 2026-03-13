---
name: qa-engineer
description: QA Engineer role - Testing strategy, test automation, unit/integration/E2E testing, performance testing, and quality metrics for VCT Platform.
---

# QA Engineer - VCT Platform

## Role Overview
Ensures software quality through comprehensive testing strategies, test automation, and quality metrics. Covers unit, integration, E2E, and performance testing across web (React 20), mobile (Expo React Native), and backend (Go 1.26).

## Technology Stack
- **Backend Unit Tests**: Go 1.26 `testing` + testify
- **Frontend Unit Tests (Web)**: Vitest + React Testing Library
- **Frontend Unit Tests (Mobile)**: Jest + @testing-library/react-native
- **Integration Tests**: Go test with Neon branch isolation
- **E2E Tests (Web)**: Playwright
- **E2E Tests (Mobile)**: Detox / Maestro
- **Performance Tests**: k6 / Artillery
- **API Testing**: Hurl / Postman / httptest
- **Coverage**: go test -cover / Vitest coverage (istanbul)
- **Database**: Neon branching (isolated test environments per CI run)

## Core Patterns

### 1. Testing Pyramid

```
         ╱╲
        ╱  ╲         E2E Tests (10%)
       ╱    ╲        - Critical user flows (web + mobile)
      ╱──────╲       - Smoke tests
     ╱        ╲
    ╱          ╲     Integration Tests (20%)
   ╱            ╲    - API endpoints
  ╱──────────────╲   - Database queries (Neon branch)
 ╱                ╲
╱                  ╲  Unit Tests (70%)
╱────────────────────╲ - Domain logic
                       - Use cases
                       - Utilities
```

### 2. Go Unit Test Patterns

#### Table-Driven Tests (Mandatory)
```go
func TestCalculateAge(t *testing.T) {
    tests := []struct {
        name     string
        dob      time.Time
        wantAge  int
        wantErr  bool
    }{
        {
            name:    "adult athlete",
            dob:     time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC),
            wantAge: 26,
        },
        {
            name:    "minor athlete",
            dob:     time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC),
            wantAge: 10,
        },
        {
            name:    "future date",
            dob:     time.Now().AddDate(1, 0, 0),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            age, err := CalculateAge(tt.dob)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.wantAge, age)
        })
    }
}
```

#### Mock Pattern
```go
// Generate mocks with mockery
//go:generate mockery --name=AthleteRepository --output=mocks

// Or manual mock
type MockAthleteRepository struct {
    mock.Mock
}

func (m *MockAthleteRepository) Create(ctx context.Context, a *domain.Athlete) error {
    args := m.Called(ctx, a)
    return args.Error(0)
}
```

### 3. Neon Branch Isolation for Integration Tests

```go
// Test setup with Neon branch
func TestIntegration_AthleteRepository(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Use Neon branch URL from CI environment
    dbURL := os.Getenv("NEON_TEST_DATABASE_URL")
    if dbURL == "" {
        t.Skip("NEON_TEST_DATABASE_URL not set")
    }

    pool, err := pgxpool.New(context.Background(), dbURL)
    require.NoError(t, err)
    defer pool.Close()

    repo := postgres.NewAthleteRepository(pool)

    t.Run("Create and Get", func(t *testing.T) {
        athlete := &domain.Athlete{
            FirstName: "Test",
            LastName:  "Athlete",
            Email:     fmt.Sprintf("test-%d@example.com", time.Now().UnixNano()),
        }
        err := repo.Create(context.Background(), athlete)
        assert.NoError(t, err)
        assert.NotEqual(t, uuid.Nil, athlete.ID)

        fetched, err := repo.GetByID(context.Background(), athlete.ID)
        assert.NoError(t, err)
        assert.Equal(t, athlete.FirstName, fetched.FirstName)
    })
}
```

### 4. API Integration Test Pattern

```go
func TestAthleteAPI(t *testing.T) {
    srv := setupTestServer(t)
    defer srv.Close()

    t.Run("POST /api/v1/athletes - success", func(t *testing.T) {
        body := `{"first_name":"Nguyen","last_name":"Van A","email":"a@test.com","gender":"male","date_of_birth":"2000-01-01"}`
        resp, err := http.Post(srv.URL+"/api/v1/athletes", "application/json", strings.NewReader(body))
        require.NoError(t, err)
        assert.Equal(t, http.StatusCreated, resp.StatusCode)

        var result map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&result)
        assert.True(t, result["success"].(bool))
    })

    t.Run("POST /api/v1/athletes - validation error", func(t *testing.T) {
        body := `{"first_name":"","email":"invalid"}`
        resp, err := http.Post(srv.URL+"/api/v1/athletes", "application/json", strings.NewReader(body))
        require.NoError(t, err)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("GET /api/v1/athletes - with pagination", func(t *testing.T) {
        resp, err := http.Get(srv.URL + "/api/v1/athletes?page=1&per_page=10")
        require.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)
    })
}
```

### 5. Frontend Test Patterns

#### Web (Vitest + RTL)
```tsx
// AthleteCard.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { AthleteCard } from './AthleteCard';

const mockAthlete = {
  id: '1',
  firstName: 'Nguyễn',
  lastName: 'Văn A',
  status: 'active',
  email: 'a@test.com',
};

describe('AthleteCard', () => {
  it('renders athlete name', () => {
    render(<AthleteCard athlete={mockAthlete} />);
    expect(screen.getByText('Nguyễn Văn A')).toBeInTheDocument();
  });

  it('calls onEdit when edit button clicked', () => {
    const onEdit = vi.fn();
    render(<AthleteCard athlete={mockAthlete} onEdit={onEdit} />);
    fireEvent.click(screen.getByRole('button', { name: /edit/i }));
    expect(onEdit).toHaveBeenCalledWith('1');
  });
});
```

#### Mobile (Jest + RNTL)
```tsx
// AthleteCard.test.tsx (React Native)
import { render, fireEvent } from '@testing-library/react-native';
import { AthleteCard } from './AthleteCard';

describe('AthleteCard (Native)', () => {
  it('renders athlete name', () => {
    const { getByText } = render(<AthleteCard athlete={mockAthlete} />);
    expect(getByText('Nguyễn Văn A')).toBeTruthy();
  });

  it('handles press event', () => {
    const onPress = jest.fn();
    const { getByTestId } = render(
      <AthleteCard athlete={mockAthlete} onPress={onPress} />
    );
    fireEvent.press(getByTestId('athlete-card'));
    expect(onPress).toHaveBeenCalledWith('1');
  });
});
```

### 6. E2E Test Pattern (Playwright - Web)

```typescript
// tests/e2e/athlete.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Athlete Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
    await page.fill('[name="email"]', 'admin@vct.com');
    await page.fill('[name="password"]', 'password');
    await page.click('button[type="submit"]');
    await page.waitForURL('/dashboard');
  });

  test('should create new athlete', async ({ page }) => {
    await page.goto('/athletes/new');
    await page.fill('[name="firstName"]', 'Nguyễn');
    await page.fill('[name="lastName"]', 'Văn B');
    await page.fill('[name="email"]', 'b@test.com');
    await page.selectOption('[name="gender"]', 'male');
    await page.fill('[name="dateOfBirth"]', '2000-01-01');
    await page.click('button[type="submit"]');
    await expect(page.getByText('Thêm VĐV thành công')).toBeVisible();
  });
});
```

### 7. Performance Test (k6)

```javascript
// tests/performance/athlete-load.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 100 },   // Ramp up
    { duration: '3m', target: 100 },   // Sustained
    { duration: '1m', target: 500 },   // Peak (tournament day)
    { duration: '2m', target: 500 },   // Sustained peak
    { duration: '1m', target: 0 },     // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],   // 95% under 500ms
    http_req_failed: ['rate<0.01'],     // Error rate < 1%
  },
};

export default function () {
  const res = http.get(`${__ENV.API_URL}/api/v1/athletes?page=1&per_page=20`);
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
  sleep(1);
}
```

### 8. Quality Metrics

| Metric | Target | Tool |
|--------|--------|------|
| Unit test coverage | ≥ 80% | go test -cover / vitest |
| Integration test coverage | ≥ 60% | go test (Neon branch) |
| E2E critical paths | 100% | Playwright / Detox |
| Bug escape rate | < 5% | Manual tracking |
| Regression rate | 0% | CI/CD |
| Test execution time | < 10 min | CI pipeline |

### 9. QA Checklist per Feature
- [ ] Unit tests for all domain logic (Go 1.26)
- [ ] Unit tests for use cases with mocks
- [ ] API integration tests (happy + error paths)
- [ ] Frontend component tests (web: Vitest, mobile: Jest)
- [ ] E2E test for critical user flow (Playwright + Detox)
- [ ] Vietnamese text rendering verified
- [ ] Mobile responsiveness tested (web + native)
- [ ] Edge cases covered (empty, max, special chars)
- [ ] Supabase RLS permission boundaries tested
- [ ] Performance under expected load verified
- [ ] Neon branch used for test isolation in CI
