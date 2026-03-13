---
description: Build, test, and deploy VCT Platform to staging environment
---

# /deploy-staging - Staging Deployment

// turbo-all

## Prerequisites
- All CI checks passing on `develop` branch
- Docker images built successfully
- Staging server/cluster accessible

## Steps

### Step 1: Ensure develop branch is up to date
```bash
git checkout develop
git pull origin develop
```

### Step 2: Run full test suite
```bash
# Backend
cd backend && go test -race -count=1 ./...

# Frontend
cd frontend && npm run test -- --run && npm run build
```

### Step 3: Build Docker images for staging
```bash
# Build with staging tag
docker build -t ghcr.io/vct-platform/api:staging -f backend/Dockerfile backend/
docker build -t ghcr.io/vct-platform/web:staging -f frontend/Dockerfile frontend/
```

### Step 4: Push images to registry
```bash
docker push ghcr.io/vct-platform/api:staging
docker push ghcr.io/vct-platform/web:staging
```

### Step 5: Apply database migrations on Neon staging branch
```bash
# Apply migrations to Neon staging branch
migrate -path supabase/migrations \
  -database "${NEON_STAGING_DATABASE_URL}" \
  up

# Sync Supabase RLS policies
supabase db push --linked
```

### Step 6: Deploy to staging
```bash
# Docker Compose (simple staging)
docker compose -f deploy/docker/docker-compose.staging.yml pull
docker compose -f deploy/docker/docker-compose.staging.yml up -d

# OR Kubernetes
kubectl apply -f deploy/k8s/staging/ --namespace=vct-staging
kubectl rollout status deployment/vct-api --namespace=vct-staging
```

### Step 7: Verify staging deployment
```bash
# Health check
curl -f https://staging-api.vct-platform.com/health

# Readiness check
curl -f https://staging-api.vct-platform.com/ready

# Smoke test
curl -f https://staging-api.vct-platform.com/api/v1/health
```

### Step 8: Notify team
```bash
echo "✅ Staging deployed successfully"
echo "API: https://staging-api.vct-platform.com"
echo "Web: https://staging.vct-platform.com"
echo "Commit: $(git rev-parse --short HEAD)"
```

## Rollback
```bash
# Revert to previous image
docker compose -f deploy/docker/docker-compose.staging.yml down
docker tag ghcr.io/vct-platform/api:previous ghcr.io/vct-platform/api:staging
docker compose -f deploy/docker/docker-compose.staging.yml up -d

# Rollback migrations
# Rollback Neon staging branch
migrate -path supabase/migrations -database "${NEON_STAGING_DATABASE_URL}" down 1

# Or reset staging branch from main
neonctl branches reset staging --project-id ${NEON_PROJECT_ID}
```
