---
description: Initialize a new VCT Platform project from scratch with full directory structure, config files, and boilerplate code
---

# /project-init - Initialize VCT Platform Project

// turbo-all

## Prerequisites
- Go 1.26+ installed
- Node.js 25+ installed
- Docker 27+ & Docker Compose installed
- Git installed
- Neon CLI (`npm install -g neonctl`)
- Supabase CLI (`npm install -g supabase`)
- Expo CLI (`npm install -g expo-cli eas-cli`)

## Steps

1. Create root project structure
```bash
mkdir -p backend/{cmd/server,internal/{modules,shared/{middleware,errors,pagination,config},pkg/{httputil,dbutil,testutil}},docs}
mkdir -p apps/web/src/{app/providers,modules,shared/{components/{ui,layout,feedback},hooks,utils,types},styles}
mkdir -p apps/mobile/app/{auth,"(tabs)"}
mkdir -p apps/mobile/components
mkdir -p packages/shared/{types,hooks,services,utils}
mkdir -p packages/ui
mkdir -p supabase/migrations
mkdir -p i18n/{vi,en}
mkdir -p deploy/{docker,k8s,scripts}
mkdir -p docs/{api,architecture,guides}
mkdir -p .github/workflows
```

2. Initialize Go module
```bash
cd backend
go mod init github.com/vct-platform/backend
```

3. Install Go dependencies
```bash
cd backend
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5
go get github.com/google/uuid
go get github.com/rs/zerolog
go get github.com/go-playground/validator/v10
go get github.com/golang-jwt/jwt/v5
go get github.com/golang-migrate/migrate/v4
go get github.com/swaggo/swag/cmd/swag
go get github.com/stretchr/testify
go get github.com/redis/go-redis/v9
```

4. Initialize web frontend
```bash
cd apps/web
npx -y create-vite@latest ./ -- --template react-ts
npm install @supabase/supabase-js @tanstack/react-query react-router-dom react-i18next i18next react-hook-form @hookform/resolvers zod recharts lucide-react
npm install -D @testing-library/react @testing-library/jest-dom vitest jsdom tailwindcss@latest
```

5. Initialize mobile app
```bash
cd apps/mobile
npx -y create-expo-app@latest ./ --template tabs
npm install @supabase/supabase-js @tanstack/react-query nativewind
```

6. Initialize Supabase project
```bash
supabase init
# Creates supabase/ directory with config.toml
```

7. Create environment files
```bash
cat > .env.example << 'EOF'
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENV=development

# Neon (Serverless PostgreSQL)
NEON_DATABASE_URL=postgres://user:pass@ep-xxx.neon.tech/vct_platform?sslmode=require
NEON_PROJECT_ID=your-neon-project-id
NEON_API_KEY=your-neon-api-key

# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=eyJ...
SUPABASE_SERVICE_KEY=eyJ...
SUPABASE_JWT_SECRET=your-jwt-secret

# Redis
REDIS_URL=redis://localhost:6379

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000
EOF

cp .env.example .env
```

8. Create .gitignore
```bash
cat > .gitignore << 'EOF'
# Go
backend/server
*.exe
*.test
*.out
vendor/

# Node
node_modules/
dist/
.env.local

# Mobile
apps/mobile/.expo/
apps/mobile/ios/
apps/mobile/android/

# IDE
.idea/
.vscode/
*.swp

# Environment
.env
.env.local
*.log

# OS
.DS_Store
Thumbs.db

# Docker
docker-compose.override.yml

# Supabase
supabase/.temp/
EOF
```

9. Initialize Git repository
```bash
git init
git add .
git commit -m "chore: initial project scaffold for VCT Platform"
```

10. Verify project structure
```bash
echo "=== Backend ==="
ls -la backend/
echo "=== Web App ==="
ls -la apps/web/
echo "=== Mobile App ==="
ls -la apps/mobile/
echo "=== Supabase ==="
ls -la supabase/
echo "=== Structure ==="
find . -type d -not -path '*/node_modules/*' -not -path '*/.git/*' -not -path '*/.expo/*' | head -60
```
