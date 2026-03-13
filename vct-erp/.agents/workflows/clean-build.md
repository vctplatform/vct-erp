---
description: Clean build and dependency refresh for VCT Platform
---

# /clean-build - Clean Build & Dependency Reset

// turbo-all

## When to Use
- Build errors that don't make sense
- Dependency version conflicts
- Docker image issues
- Cache corruption
- After major dependency updates

## Steps

### Step 1: Clean Go build cache
```bash
cd backend
go clean -cache
go clean -testcache
go clean -modcache
```

### Step 2: Re-download Go dependencies
```bash
cd backend
go mod tidy
go mod download
go mod verify
```

### Step 3: Clean Node modules
```bash
cd frontend
rm -rf node_modules
rm -f package-lock.json
```

### Step 4: Fresh npm install
```bash
cd frontend
npm install
```

### Step 5: Clean Docker resources
```bash
# Stop all containers
docker compose down

# Remove unused images
docker image prune -f

# Remove unused volumes (⚠️ deletes data!)
# docker volume prune -f

# Remove build cache
docker builder prune -f
```

### Step 6: Rebuild Docker images
```bash
docker compose build --no-cache
```

### Step 7: Rebuild and test backend
```bash
cd backend
go build -v ./...
go vet ./...
go test -count=1 ./...
```

### Step 8: Rebuild and test frontend
```bash
cd frontend
npm run lint
npx tsc --noEmit
npm run build
npm run test -- --run
```

### Step 9: Start fresh
```bash
docker compose up -d
echo "✅ Clean build complete"
```

## Nuclear Option (full reset)
```bash
# ⚠️ WARNING: This removes ALL Docker data including database volumes!
docker compose down -v --remove-orphans
docker system prune -af --volumes
# Then re-run /setup-environment and /database-migration
```
