---
description: Docker development environment setup and management for VCT Platform
---

# /docker-dev - Docker Development Environment

// turbo-all

## Architecture
```
┌──────────────────────────────────────────────────────┐
│  Docker Network: vct-network                          │
│                                                       │
│  ┌─────────┐  ┌─────────┐                            │
│  │   API   │  │   Web   │   External Services:       │
│  │  :8080  │  │  :3000  │   ┌──────────┐             │
│  └────┬────┘  └────┬────┘   │  Neon DB │ (cloud)     │
│       │             │        │  :5432   │             │
│  ┌────▼────┐                 └──────────┘             │
│  │  Redis  │                 ┌──────────┐             │
│  │  :6379  │                 │ Supabase │ (cloud)     │
│  └─────────┘                 │  Auth/RT │             │
│                              └──────────┘             │
└──────────────────────────────────────────────────────┘
```

> **Note**: Database (PostgreSQL) is hosted on Neon (serverless). Auth/Realtime/Storage on Supabase. Only Redis runs locally.

## Steps

### Step 1: Create docker-compose.yml (if not exists)
```yaml
services:
  api:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=${NEON_DATABASE_URL}               # Neon serverless PostgreSQL
      - SUPABASE_URL=${SUPABASE_URL}                    # Supabase project URL
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
      - SUPABASE_SERVICE_KEY=${SUPABASE_SERVICE_KEY}
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${SUPABASE_JWT_SECRET}
      - ENV=development
    volumes:
      - ./backend:/app
    depends_on:
      redis:
        condition: service_started
    networks:
      - vct-network

  web:
    build:
      context: ./frontend
      dockerfile: Dockerfile.dev
    ports:
      - "3000:3000"
    environment:
      - VITE_API_URL=http://localhost:8080
      - VITE_SUPABASE_URL=${SUPABASE_URL}
      - VITE_SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
    volumes:
      - ./frontend/src:/app/src
      - ./frontend/public:/app/public
    networks:
      - vct-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - vct-network

networks:
  vct-network:
    driver: bridge

volumes:
  redis_data:
```

### Step 2: Create .env file
```bash
# Copy from template
cp .env.example .env

# Fill in values from Neon & Supabase dashboards:
# NEON_DATABASE_URL=postgres://user:pass@ep-xxx.neon.tech/vct_platform?sslmode=require
# SUPABASE_URL=https://your-project.supabase.co
# SUPABASE_ANON_KEY=eyJ...
# SUPABASE_SERVICE_KEY=eyJ...
# SUPABASE_JWT_SECRET=your-jwt-secret
```

### Step 3: Start development stack
```bash
docker compose up -d
```

### Step 4: Check all services are running
```bash
docker compose ps
```

### Step 5: View logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f api
docker compose logs -f web
```

### Step 6: Access services
```bash
echo "API:             http://localhost:8080"
echo "Frontend:        http://localhost:3000"
echo "Neon Dashboard:  https://console.neon.tech"
echo "Supabase Studio: https://supabase.com/dashboard"
echo "Redis:           localhost:6379"
```

### Step 7: Run commands inside containers
```bash
# Go shell
docker compose exec api sh

# Redis CLI
docker compose exec redis redis-cli

# Connect to Neon DB directly
psql "${NEON_DATABASE_URL}"
```

### Step 8: Rebuild after Dockerfile changes
```bash
docker compose up -d --build
```

### Step 9: Stop all services
```bash
docker compose down
```

## Optional: Supabase Local Development
```bash
# Start Supabase locally (instead of cloud)
supabase start

# Use local Supabase URLs in .env:
# SUPABASE_URL=http://localhost:54321
# SUPABASE_ANON_KEY=<local-anon-key>
# SUPABASE_SERVICE_KEY=<local-service-key>
```

## Troubleshooting

| Issue | Solution |
|-------|---------|
| Neon connection refused | Check IP allowlist in Neon dashboard |
| Supabase auth errors | Verify JWT secret matches Supabase config |
| Port already in use | `docker compose down` then retry |
| Container won't start | `docker compose logs {service}` |
| Slow on Windows | Use WSL2 backend for Docker Desktop |
| Volume sync slow | Exclude `node_modules` from volume mounts |
