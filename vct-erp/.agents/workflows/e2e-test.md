---
description: Run end-to-end tests for VCT Platform (Playwright for web, Detox/Maestro for mobile)
---

# /e2e-test - End-to-End Testing Workflow

// turbo-all

## Architecture
```
E2E Tests
├── Web (Playwright)
│   ├── Auth flows (Supabase login/register)
│   ├── CRUD flows (athletes, clubs, tournaments)
│   ├── Permission boundaries (RLS)
│   └── Responsive (desktop + mobile viewport)
└── Mobile (Detox / Maestro)
    ├── Auth flows (Supabase + deep links)
    ├── Tab navigation
    ├── Offline / network error handling
    └── Platform-specific (iOS + Android)
```

## Prerequisites
- Playwright installed: `npx playwright install`
- Detox CLI: `npm install -g detox-cli`
- Maestro CLI: `curl -Ls https://get.maestro.mobile.dev | bash`
- Neon test branch created (isolated data)
- Supabase test project or local Supabase (`supabase start`)

## Steps

### Step 1: Create Neon branch for E2E tests
```bash
neonctl branches create \
  --project-id ${NEON_PROJECT_ID} \
  --name e2e-test-$(date +%Y%m%d) \
  --parent main
```

### Step 2: Seed test data
```bash
# Apply migrations to test branch
migrate -path supabase/migrations \
  -database "$(neonctl connection-string e2e-test-$(date +%Y%m%d) --project-id ${NEON_PROJECT_ID})" \
  up

# Run seed script
psql "$(neonctl connection-string e2e-test-$(date +%Y%m%d) --project-id ${NEON_PROJECT_ID})" \
  -f tests/e2e/seed.sql
```

### Step 3: Run Playwright tests (Web)
```bash
cd apps/web

# Run all E2E tests
npx playwright test

# Run specific test file
npx playwright test tests/e2e/athlete.spec.ts

# Run with UI mode (interactive debugging)
npx playwright test --ui

# Run in specific browser
npx playwright test --project=chromium
npx playwright test --project=mobile-chrome
```

### Step 4: View Playwright report
```bash
npx playwright show-report
```

### Step 5: Run Detox tests (Mobile - React Native)
```bash
cd apps/mobile

# Build for E2E (iOS)
detox build --configuration ios.sim.release

# Run tests (iOS)
detox test --configuration ios.sim.release

# Build for E2E (Android)
detox build --configuration android.emu.release

# Run tests (Android)
detox test --configuration android.emu.release
```

### Step 6: Run Maestro tests (Mobile - Alternative)
```bash
cd apps/mobile

# Run single flow
maestro test tests/e2e/flows/login.yaml

# Run all flows
maestro test tests/e2e/flows/

# Record flow (interactive)
maestro record
```

### Step 7: Cleanup Neon test branch
```bash
neonctl branches delete e2e-test-$(date +%Y%m%d) --project-id ${NEON_PROJECT_ID}
```

## Test Categories

| Category | Web (Playwright) | Mobile (Detox/Maestro) |
|----------|-----------------|----------------------|
| Auth | ✅ Login, register, MFA | ✅ Biometric, deep links |
| Navigation | ✅ All routes | ✅ Tab bar, stack nav |
| CRUD | ✅ All modules | ✅ Core modules |
| Permissions | ✅ Role-based access | ✅ Role-based access |
| Responsive | ✅ 375px-1440px | N/A (native) |
| Offline | N/A | ✅ Network error states |
| i18n | ✅ VI/EN switch | ✅ VI/EN switch |
| Accessibility | ✅ Keyboard nav, ARIA | ✅ VoiceOver, TalkBack |

## CI Integration (GitHub Actions)
```yaml
e2e-web:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-node@v4
    - run: npx playwright install --with-deps
    - run: npx playwright test
    - uses: actions/upload-artifact@v4
      if: failure()
      with:
        name: playwright-report
        path: playwright-report/
```

## Troubleshooting

| Issue | Solution |
|-------|---------|
| Playwright timeout | Increase `timeout` in `playwright.config.ts` |
| Detox build fail | `detox clean-framework-cache && detox build` |
| Maestro not finding element | Use `maestro studio` for element inspection |
| Supabase auth flaky | Use service key for test user creation |
| Neon branch stale | Delete and recreate branch from main |
