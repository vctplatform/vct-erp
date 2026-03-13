---
description: Set up developer environment with all required tools for VCT Platform development (Go, Node, Docker, Neon, Supabase, Redis)
---

# /setup-environment - Developer Environment Setup

// turbo-all

## Step 1: Check system prerequisites
```powershell
Write-Host "=== System Check ===" -ForegroundColor Cyan
Write-Host "OS: $([System.Environment]::OSVersion.VersionString)"
Write-Host "Architecture: $([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture)"
```

## Step 2: Install Go 1.26+ (if not installed)
```powershell
# Check Go version
go version
# If not installed or outdated:
# winget install GoLang.Go
```

## Step 3: Install Node.js 25+ (if not installed)
```powershell
# Check Node version
node --version
npm --version
# If not installed:
# winget install OpenJS.NodeJS
```

## Step 4: Install Docker Desktop (if not installed)
```powershell
# Check Docker 27+
docker --version
docker compose version
# If not installed:
# winget install Docker.DockerDesktop
```

## Step 5: Install Git (if not installed)
```powershell
git --version
# If not installed:
# winget install Git.Git
```

## Step 6: Install Neon CLI
```bash
npm install -g neonctl
neonctl --version

# Auth with Neon
neonctl auth
```

## Step 7: Install Supabase CLI
```bash
npm install -g supabase
supabase --version

# Login to Supabase
supabase login
```

## Step 8: Install Expo CLI (for mobile development)
```bash
npm install -g expo-cli eas-cli
expo --version
eas --version

# Login to Expo
eas login
```

## Step 9: Install Go tools
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/cosmtrek/air@latest  # Hot reload
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/vektra/mockery/v2@latest
```

## Step 10: Install global npm packages
```bash
npm install -g typescript eslint prettier
```

## Step 11: Configure Git
```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
git config --global init.defaultBranch main
git config --global pull.rebase true
git config --global core.autocrlf true  # Windows
```

## Step 12: Set up environment variables
```bash
# Create .env file from template
cp .env.example .env

# Required variables:
# NEON_DATABASE_URL=postgres://user:pass@ep-xxx.neon.tech/vct_platform?sslmode=require
# SUPABASE_URL=https://your-project.supabase.co
# SUPABASE_ANON_KEY=eyJ...
# SUPABASE_SERVICE_KEY=eyJ...
# SUPABASE_JWT_SECRET=your-jwt-secret
# REDIS_URL=redis://localhost:6379
```

## Step 13: Start local services (Redis only, DB is Neon)
```bash
docker compose up -d redis
docker compose exec redis redis-cli ping
```

## Step 14: Start Supabase local dev (optional)
```bash
supabase start
# Supabase Studio: http://localhost:54323
# Supabase API:    http://localhost:54321
```

## Step 15: Install backend dependencies
```bash
cd backend
go mod download
go mod verify
```

## Step 16: Install frontend dependencies
```bash
cd frontend
npm ci
```

## Step 17: Run database migrations (on Neon dev branch)
```bash
# Apply migrations to Neon dev branch
migrate -path migrations -database "${NEON_DATABASE_URL}" up

# Or using Supabase CLI
supabase db push
```

## Step 18: Start development servers
```bash
# Terminal 1: Backend with hot reload
cd backend && air

# Terminal 2: Frontend (Web) with Vite
cd frontend && npm run dev

# Terminal 3: Mobile (Expo) - optional
cd mobile && npx expo start
```

## Step 19: Verify development environment
```bash
echo "=== Environment Ready ==="
echo "Backend API:     http://localhost:8080"
echo "Frontend (Web):  http://localhost:3000"
echo "Mobile (Expo):   exp://localhost:8081"
echo "API Docs:        http://localhost:8080/swagger/index.html"
echo "Neon Dashboard:  https://console.neon.tech"
echo "Supabase Studio: https://supabase.com/dashboard"
```

## Troubleshooting

| Issue | Solution |
|-------|---------|
| Neon connection timeout | Check IP allowlist, verify SSL mode |
| Supabase JWT invalid | Regenerate keys in Supabase dashboard |
| Port 3000 in use | Kill process: `npx kill-port 3000` |
| Docker not starting | Restart Docker Desktop |
| Go module errors | `go clean -modcache && go mod download` |
| npm install fails | `rm -rf node_modules package-lock.json && npm install` |
| Expo build error | `npx expo install --fix` |
