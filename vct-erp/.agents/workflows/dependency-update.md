---
description: Audit and upgrade project dependencies with security vulnerability checks
---

# /dependency-update - Dependency Update Workflow

// turbo-all

## When to Use
- Monthly scheduled dependency audit
- Security vulnerability alert (Dependabot/Snyk)
- Major version upgrade planning
- Before a release cycle

## Steps

### Step 1: Check for Go dependency updates
```bash
cd backend

# List outdated dependencies
go list -m -u all

# Check for vulnerabilities
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### Step 2: Update Go dependencies
```bash
cd backend

# Update all minor/patch versions
go get -u ./...
go mod tidy

# Update specific dependency (major version)
# go get github.com/jackc/pgx/v5@latest

# Verify
go build ./...
go test -race -count=1 ./...
```

### Step 3: Check for npm dependency updates (Web)
```bash
cd apps/web

# Check outdated packages
npm outdated

# Audit for vulnerabilities
npm audit

# Interactive upgrade
npx npm-check-updates -i
```

### Step 4: Update npm dependencies (Web)
```bash
cd apps/web

# Update minor/patch
npm update

# Apply security fixes
npm audit fix

# For major version upgrades
npx npm-check-updates -u
npm install

# Verify
npm run build
npm run test -- --run
```

### Step 5: Update Expo / React Native dependencies (Mobile)
```bash
cd apps/mobile

# Expo-managed update (recommended)
npx expo install --fix

# Check compatibility
npx expo-doctor

# Update Expo SDK
npx expo install expo@latest

# Verify
npx expo start --clear
```

### Step 6: Update Supabase CLI & types
```bash
# Update Supabase CLI
npm install -g supabase@latest

# Regenerate TypeScript types from Supabase schema
supabase gen types typescript --linked > packages/shared/types/database.ts
```

### Step 7: Update Neon CLI
```bash
npm install -g neonctl@latest
neonctl --version
```

### Step 8: Update Go linting tools
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/cosmtrek/air@latest
go install github.com/vektra/mockery/v2@latest
```

### Step 9: Run full test suite
```bash
# Backend
cd backend && go test -race -count=1 ./...

# Frontend (Web)
cd apps/web && npm run test -- --run && npm run build

# Mobile
cd apps/mobile && npx expo-doctor
```

### Step 10: Commit and create PR
```bash
git checkout -b chore/dependency-update-$(date +%Y%m%d)
git add -A
git commit -m "chore(deps): update all dependencies $(date +%Y-%m-%d)"
git push origin HEAD
gh pr create --title "chore(deps): dependency update $(date +%Y-%m-%d)" \
  --body "## Changes\n- Updated Go dependencies\n- Updated npm packages (web + mobile)\n- Updated CLI tools (Supabase, Neon, Expo)\n- All tests passing"
```

## Automation (GitHub Actions)

```yaml
# .github/workflows/dependency-update.yml
name: Dependency Audit
on:
  schedule:
    - cron: '0 9 1 * *'  # Monthly, 1st day at 9am
  workflow_dispatch:
jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.26' }
      - uses: actions/setup-node@v4
        with: { node-version: '25' }
      - run: govulncheck ./...
        working-directory: backend
      - run: npm audit --audit-level=high
        working-directory: apps/web
```

## Version Pinning Policy

| Category | Policy | Example |
|----------|--------|---------|
| Go stdlib | Follow Go releases | Go 1.26.x |
| Go deps (critical) | Pin major, float minor | `pgx v5.x` |
| React/Expo | Pin major | `react@20.x`, `expo@53.x` |
| Tailwind CSS | Pin major | `tailwindcss@4.x` |
| CLI tools | Always latest | `supabase@latest` |
| Dev deps | Float freely | `vitest@latest` |
