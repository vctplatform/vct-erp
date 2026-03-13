---
name: devops-engineer
description: DevOps Engineer role - CI/CD pipelines, Docker, Kubernetes, infrastructure as code, monitoring, and deployment automation for VCT Platform.
---

# DevOps Engineer - VCT Platform

## Role Overview
Manages infrastructure, CI/CD pipelines, containerization, deployments, and monitoring. Ensures the platform is reliable, scalable, and deployable with minimal downtime. Integrates with Neon (serverless PostgreSQL) and Supabase (BaaS) for managed services.

## Technology Stack
- **Containers**: Docker 27+ / Docker Compose
- **Orchestration**: Kubernetes (EKS/GKE/AKS) for production
- **CI/CD**: GitHub Actions
- **Mobile CI**: Expo EAS Build (cloud builds for iOS/Android)
- **IaC**: Terraform / Docker Compose
- **Registry**: GitHub Container Registry (ghcr.io)
- **Database**: Neon (serverless PostgreSQL 18+, branching per PR)
- **BaaS**: Supabase (Auth, Realtime, Storage, Edge Functions)
- **Monitoring**: Prometheus + Grafana + Neon Dashboard + Supabase Dashboard
- **Logging**: ELK Stack / Loki
- **Secrets**: GitHub Secrets / HashiCorp Vault
- **DNS/CDN**: Cloudflare
- **SSL**: Let's Encrypt (auto-renewal)

## Core Patterns

### 1. Docker Configuration

#### Backend Dockerfile (Multi-stage)
```dockerfile
# Build stage
FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# Runtime stage
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
COPY --from=builder /server /server
COPY --from=builder /app/migrations /migrations
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
```

#### Frontend Dockerfile
```dockerfile
FROM node:25 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:latest
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
```

#### Docker Compose (Development)
```yaml
services:
  api:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=${NEON_DATABASE_URL}          # Neon serverless PostgreSQL
      - SUPABASE_URL=${SUPABASE_URL}               # Supabase project URL
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}     # Supabase anon key
      - SUPABASE_SERVICE_KEY=${SUPABASE_SERVICE_KEY}
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${SUPABASE_JWT_SECRET}
    volumes:
      - ./backend:/app
    depends_on:
      redis:
        condition: service_started

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

  redis:
    image: redis:7
    ports:
      - "6379:6379"

  # Supabase local development
  supabase:
    image: supabase/supabase-dev
    ports:
      - "54321:54321"  # Supabase Studio
      - "54322:54322"  # Supabase API
    environment:
      POSTGRES_PASSWORD: postgres

volumes:
  postgres_data:
```

### 2. GitHub Actions CI/CD

#### CI Pipeline (`.github/workflows/ci.yml`)
```yaml
name: CI
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

env:
  NEON_PROJECT_ID: ${{ secrets.NEON_PROJECT_ID }}
  NEON_API_KEY: ${{ secrets.NEON_API_KEY }}

jobs:
  backend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      # Create Neon branch for this PR
      - name: Create Neon Branch
        id: neon-branch
        uses: neondatabase/create-branch-action@v5
        with:
          project_id: ${{ env.NEON_PROJECT_ID }}
          api_key: ${{ env.NEON_API_KEY }}
          branch_name: ci/${{ github.sha }}
          parent: main

      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - run: go mod download
        working-directory: backend
      - run: golangci-lint run ./...
        working-directory: backend
      - run: go test -race -coverprofile=coverage.out ./...
        working-directory: backend
        env:
          DATABASE_URL: ${{ steps.neon-branch.outputs.db_url_with_pooler }}

      # Cleanup Neon branch
      - name: Delete Neon Branch
        if: always()
        uses: neondatabase/delete-branch-action@v3
        with:
          project_id: ${{ env.NEON_PROJECT_ID }}
          api_key: ${{ env.NEON_API_KEY }}
          branch: ci/${{ github.sha }}

  frontend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '25'
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json
      - run: npm ci
        working-directory: frontend
      - run: npm run lint
        working-directory: frontend
      - run: npm run type-check
        working-directory: frontend
      - run: npm run test -- --run
        working-directory: frontend
      - run: npm run build
        working-directory: frontend

  mobile-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '25'
      - uses: expo/expo-github-action@v8
        with:
          eas-version: latest
          token: ${{ secrets.EXPO_TOKEN }}
      - run: npm ci
        working-directory: mobile
      - run: eas build --platform all --non-interactive --no-wait
        working-directory: mobile

  docker-build:
    needs: [backend-test, frontend-test]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v6
        with:
          context: ./backend
          push: ${{ github.ref == 'refs/heads/main' }}
          tags: ghcr.io/${{ github.repository }}/api:${{ github.sha }}
```

### 3. Neon Branching in CI/CD

```yaml
# Preview environment per PR with Neon branching
name: Preview Environment
on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  preview:
    runs-on: ubuntu-latest
    steps:
      - name: Create Neon preview branch
        uses: neondatabase/create-branch-action@v5
        id: create-branch
        with:
          project_id: ${{ secrets.NEON_PROJECT_ID }}
          api_key: ${{ secrets.NEON_API_KEY }}
          branch_name: preview/pr-${{ github.event.pull_request.number }}
          parent: main

      - name: Run migrations on preview branch
        run: |
          DATABASE_URL="${{ steps.create-branch.outputs.db_url }}" \
          go run cmd/migrate/main.go up

      - name: Comment PR with preview DB
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: '🗄️ Preview database branch created: `preview/pr-${{ github.event.pull_request.number }}`'
            })
```

### 4. Environment Strategy

| Environment | Purpose | Deployment | Database |
|-------------|---------|------------|----------|
| **Local** | Development | Docker Compose | Neon dev branch / Supabase local |
| **CI** | Testing | GitHub Actions | Neon ephemeral branch (per run) |
| **Preview** | PR review | Auto-deploy from PR | Neon branch (per PR) |
| **Staging** | Integration testing | Auto-deploy from `develop` | Neon staging branch |
| **Production** | Live system | Manual approve from `main` | Neon main branch |

### 5. Kubernetes Manifests (Production)

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vct-api
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  template:
    spec:
      containers:
        - name: api
          image: ghcr.io/org/vct-api:latest
          ports:
            - containerPort: 8080
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: neon-credentials
                  key: database-url
            - name: SUPABASE_URL
              valueFrom:
                secretKeyRef:
                  name: supabase-credentials
                  key: url
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 15
            periodSeconds: 20
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
```

### 6. Monitoring Stack

#### Health Check Endpoints
```
GET /health   → Basic health (up/down)
GET /ready    → Readiness (Neon DB connected, Redis connected, Supabase reachable)
GET /metrics  → Prometheus metrics
```

#### Key Metrics to Monitor
| Metric | Alert Threshold | Source |
|--------|----------------|--------|
| HTTP response time (P95) | > 500ms | Prometheus |
| Error rate (5xx) | > 1% | Prometheus |
| CPU usage | > 80% for 5 min | Prometheus |
| Memory usage | > 85% | Prometheus |
| Neon compute hours | > budget limit | Neon Dashboard |
| Neon storage | > 80% of plan | Neon Dashboard |
| Supabase API requests | > rate limit | Supabase Dashboard |
| Disk usage | > 90% | Prometheus |

### 7. Deployment Checklist
- [ ] All CI checks passing
- [ ] Docker image built and pushed
- [ ] Database migrations applied (Neon main branch)
- [ ] Supabase migrations synced (`supabase db push`)
- [ ] Environment variables updated
- [ ] Health checks responding
- [ ] Monitoring dashboards verified
- [ ] Rollback plan available (Neon PITR)
- [ ] Team notified of deployment
