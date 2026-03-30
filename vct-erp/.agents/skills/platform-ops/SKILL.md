---
name: platform-ops
description: >-
  Mega-Skill for DevOps, Security, and Infrastructure. CI/CD, Docker, Kubernetes,
  monitoring, OWASP security, audit compliance, and real-time systems.
metadata:
  author: VCT Platform
  version: "1.0.0"
  type: "Mega-Skill"
  locale: vi-VN
---

# PLATFORM-OPS — MEGA-SKILL

> Infrastructure, security, and reliability. Keep the ERP fortress running 24/7.

---

## 🔹 NĂNG LỰC: DEVOPS-ENGINEER

> *"Automate EVERYTHING you do more than twice."*

### CI/CD Pipeline
```
Best Practice Pipeline:
├── Commit → Lint + Format (30s)
├── Build → Compile + Docker build (2min)
├── Test → Unit + Integration (5min)
├── Security → SAST + Dependency scan (3min)
├── Staging → Auto-deploy + Smoke test (5min)
├── Approval → Manual gate for production
├── Production → Blue/Green or Canary
├── Verify → Health check + Synthetic monitor
└── Rollback → Automatic if error rate > 1%

Target: Commit → Production < 15 minutes
```

### Infrastructure as Code
| Tool | Use Case | VCT Stack |
|------|---------|-----------|
| Docker | Container packaging | ✅ All services |
| Docker Compose | Local dev environment | ✅ Dev |
| Kubernetes | Production orchestration | ✅ Production |
| Helm | K8s package management | ✅ Charts |
| GitHub Actions | CI/CD | ✅ Primary |

### Monitoring (3 Pillars)
```
1. METRICS (Prometheus/Grafana)
   ├── System: CPU, Memory, Disk, Network
   ├── Application: Request rate, Error rate, Duration (RED)
   ├── Business: Journal entries/day, Revenue, Balances
   └── SLO/SLI tracking

2. LOGS (Loki)
   ├── Structured JSON logs
   ├── Correlation IDs across services
   ├── Log levels: ERROR → WARN → INFO → DEBUG
   └── Retention: 30d hot, 90d warm, 1y cold

3. TRACES (Jaeger)
   ├── Distributed tracing
   ├── Latency breakdown per service
   └── Error propagation
```

### Incident Response
```
Severity:
├── SEV1: System down → ALL HANDS, < 15min response
├── SEV2: Major feature broken → On-call, < 30min
├── SEV3: Minor issue → On-call, < 2h
└── SEV4: Cosmetic → Next business day

Post-Mortem: Timeline → Impact → Root Cause (5 Whys) → Action Items
```

### Docker Compose (VCT ERP Dev)
```yaml
services:
  api:
    build: ./backend
    ports: ["8080:8080"]
    depends_on: [postgres, redis, mongo]
  postgres:
    image: postgres:18
    volumes: [pgdata:/var/lib/postgresql/data]
  redis:
    image: redis:7-alpine
  mongo:
    image: mongo:7
  frontend:
    build: ./frontend
    ports: ["3000:3000"]
```

---

## 🔹 NĂNG LỰC: SECURITY-ENGINEER

> *"Security is not a product, it's a process."*

### OWASP Top 10 for ERP
| Risk | Mitigation |
|------|-----------|
| Broken Access Control | RBAC, server-side checks, least privilege |
| Cryptographic Failures | TLS 1.3, AES-256, bcrypt for passwords |
| Injection | Parameterized queries (pgx), input validation |
| Insecure Design | Threat modeling, secure patterns |
| Misconfiguration | Hardened defaults, remove unused |
| Vulnerable Components | Dependabot, SCA scanning |
| Auth Failures | MFA, session management, rate limiting |
| Data Integrity | Code signing, SBOM, integrity checks |
| Logging Failures | Centralized logging, audit trail |
| SSRF | Allowlist URLs, disable protocols |

### ERP-Specific Security
```
Financial Data Protection:
├── Encrypt at rest (AES-256 for sensitive fields)
├── Encrypt in transit (TLS 1.3)
├── Row-Level Security (RLS) per company_code
├── Audit every financial operation (MongoDB)
├── No DELETE on financial records (soft delete / void)
├── Journal entries immutable after posting
├── Separation of duties (create vs approve)
└── API rate limiting per user/role

Authentication & Authorization:
├── JWT tokens (short-lived, 15min)
├── Refresh tokens (httpOnly cookies)
├── RBAC: admin, accountant, viewer, auditor
├── Permission: read, write, void, admin per module
└── Session management with Redis
```

### STRIDE Threat Model
```
Per feature, ask:
├── Spoofing: Can someone impersonate?
├── Tampering: Can data be modified?
├── Repudiation: Can actions be denied? (audit log!)
├── Information Disclosure: Can data leak?
├── Denial of Service: Can system be crashed?
└── Elevation of Privilege: Can viewer become admin?
```

### Compliance Readiness
| Standard | Focus |
|----------|-------|
| NĐ 13/2023 (VN PDPA) | Personal data protection |
| VAS TT200 | Vietnamese Accounting Standards |
| SOC 2 Type II | Security, Availability, Confidentiality |

---

## 🔹 NĂNG LỰC: REALTIME-ENGINEER

### WebSocket / SSE for ERP
```
Real-time Use Cases:
├── Balance updates after journal posting
├── Dashboard KPI refresh
├── Notification: approval requests
├── Multi-user conflict detection (editing same record)
└── System health status in control room

Implementation:
├── Redis Pub/Sub for broadcast
├── SSE for simple one-way (dashboard refresh)
├── WebSocket for bi-directional (collaboration)
└── Graceful degradation to polling
```

### Event-Driven Architecture
```
Outbox Pattern Flow:
├── 1. Business operation + outbox INSERT in 1 TX
├── 2. Relay worker picks up pending events
├── 3. Publish to Redis Stream
├── 4. Consumer processes event
├── 5. Mark outbox as published
└── 6. Failed? Retry with backoff

Benefits for ERP:
├── Loose coupling between modules
├── Audit trail built into flow
├── At-least-once delivery guarantee
└── Async processing for heavy operations
```

---

## Development Checklist

- [ ] Dockerfile multi-stage build (small images)
- [ ] Docker Compose for local dev
- [ ] CI pipeline: lint → test → build → scan
- [ ] Health check endpoint (/healthz)
- [ ] Structured JSON logging
- [ ] Secrets in env vars, never in code
- [ ] TLS for all external connections
- [ ] Rate limiting on APIs
- [ ] Audit logging for financial operations
- [ ] Backup strategy documented

## Trigger Patterns

- "deploy", "Docker", "Kubernetes", "CI/CD"
- "security", "bảo mật", "OWASP", "authentication"
- "monitoring", "logging", "alerts", "incident"
- "realtime", "WebSocket", "event", "stream"
