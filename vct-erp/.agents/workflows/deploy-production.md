---
description: Deploy VCT Platform to production environment with safety checks and rollback plan
---

# /deploy-production - Production Deployment

## ⚠️ CRITICAL: This workflow deploys to PRODUCTION. Every step must be verified.

## Prerequisites
- [ ] Staging deployment verified and tested
- [ ] All E2E tests pass on staging
- [ ] PM/CTO approval received
- [ ] Rollback plan reviewed
- [ ] Backup taken (see `/backup-restore`)
- [ ] Maintenance window communicated (if needed)

## Steps

### Step 1: Create release branch from develop
```bash
git checkout develop
git pull origin develop
git checkout -b release/v$(cat VERSION)
```

### Step 2: Final verification on release branch
```bash
# Backend tests
cd backend && go test -race -count=1 -v ./...

# Frontend tests and build
cd frontend && npm run test -- --run && npm run build

# Lint
cd backend && golangci-lint run ./...
cd frontend && npm run lint
```

### Step 3: Update version and changelog
```bash
# Update VERSION file
echo "1.2.0" > VERSION

# Update CHANGELOG.md
cat >> CHANGELOG.md << 'EOF'
## [1.2.0] - 2026-03-12
### Added
- Tournament live scoring system
- Provincial federation management
### Fixed
- Pagination offset bug in athlete list
### Changed
- Improved athlete search performance
EOF

git add VERSION CHANGELOG.md
git commit -m "chore(release): prepare v1.2.0"
```

### Step 4: Build production Docker images
```bash
VERSION=$(cat VERSION)
docker build -t ghcr.io/vct-platform/api:v${VERSION} -t ghcr.io/vct-platform/api:latest -f backend/Dockerfile backend/
docker build -t ghcr.io/vct-platform/web:v${VERSION} -t ghcr.io/vct-platform/web:latest -f frontend/Dockerfile frontend/
```

### Step 5: Push images to registry
```bash
docker push ghcr.io/vct-platform/api:v${VERSION}
docker push ghcr.io/vct-platform/api:latest
docker push ghcr.io/vct-platform/web:v${VERSION}
docker push ghcr.io/vct-platform/web:latest
```

### Step 6: Create Neon backup branch before deployment
```bash
# Create a branch as pre-deploy snapshot (instant, zero cost)
neonctl branches create \
  --project-id ${NEON_PROJECT_ID} \
  --name backup-pre-v${VERSION} \
  --parent main
```

### Step 7: Apply database migrations
```bash
# Apply migrations to Neon main branch
migrate -path supabase/migrations \
  -database "${NEON_DATABASE_URL}" \
  up

# Sync Supabase RLS policies and functions
supabase db push --linked
```

### Step 8: Deploy with rolling update
```bash
# Kubernetes
kubectl set image deployment/vct-api api=ghcr.io/vct-platform/api:v${VERSION} --namespace=vct-production
kubectl rollout status deployment/vct-api --namespace=vct-production --timeout=300s

kubectl set image deployment/vct-web web=ghcr.io/vct-platform/web:v${VERSION} --namespace=vct-production
kubectl rollout status deployment/vct-web --namespace=vct-production --timeout=300s
```

### Step 9: Post-deployment verification
```bash
# Health checks
curl -f https://api.vct-platform.com/health
curl -f https://api.vct-platform.com/ready

# Version check
curl -s https://api.vct-platform.com/version | jq '.version'

# Smoke test: key endpoints
curl -f https://api.vct-platform.com/api/v1/athletes?page=1&per_page=1
curl -f https://vct-platform.com
```

### Step 10: Merge release to main and tag
```bash
git checkout main
git merge release/v${VERSION} --no-ff -m "Release v${VERSION}"
git tag -a v${VERSION} -m "Release v${VERSION}"
git push origin main --tags

# Merge back to develop
git checkout develop
git merge main --no-ff
git push origin develop

# Cleanup
git branch -d release/v${VERSION}
```

### Step 11: Create GitHub release
```bash
gh release create v${VERSION} --title "v${VERSION}" --notes-file CHANGELOG.md
```

## Rollback Plan

### Immediate Rollback (< 5 min)
```bash
# Kubernetes
kubectl rollout undo deployment/vct-api --namespace=vct-production
kubectl rollout undo deployment/vct-web --namespace=vct-production

# Verify
kubectl rollout status deployment/vct-api --namespace=vct-production
```

### Database Rollback (Neon PITR)
```bash
# Option 1: Rollback last migration
migrate -path supabase/migrations -database "${NEON_DATABASE_URL}" down 1

# Option 2: Full restore via Neon branch promotion
neonctl branches set-primary \
  --project-id ${NEON_PROJECT_ID} \
  --branch backup-pre-v${VERSION}
```
