---
name: cto-tech-lead
description: CTO / Tech Lead role - Architecture decisions, technology stack selection, team leadership, technical roadmap, and strategic planning for VCT Platform.
---

# CTO / Tech Lead - VCT Platform

## Role Overview
The CTO/Tech Lead is responsible for all high-level technology decisions, architecture oversight, team leadership, and aligning technical strategy with business goals for the VCT Platform (Vietnam Cycling & Triathlon management system).

## Core Responsibilities

### 1. Technology Stack Decisions
**VCT Platform Approved Stack:**

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| **Backend** | Go 1.26 | High performance, range-over-func, enhanced generics, improved stdlib |
| **Frontend (Web)** | React 20 / TypeScript | React Compiler, Server Components, Actions, use() hook |
| **Frontend (Mobile)** | Expo React Native | Shared codebase with web, cross-platform iOS/Android |
| **Styling** | Tailwind CSS v4 | CSS-first config, `@theme` directive, zero-JS config, lightning-fast |
| **Database** | PostgreSQL 18+ | Async I/O, JSON_TABLE, virtual generated columns, incremental backup |
| **Serverless DB** | Neon | Serverless PostgreSQL, branching, autoscaling, instant provisioning |
| **BaaS** | Supabase | Auth, real-time subscriptions, storage, edge functions, RLS |
| **Runtime** | Node.js 25.6.1 | Native fetch, test runner, permission model, ESM by default |
| **Cache** | Redis 7+ | In-memory caching, pub/sub for real-time |
| **API** | REST + WebSocket | RESTful for CRUD, WebSocket for live scoring |
| **Auth** | Supabase Auth + RBAC | Built-in OAuth, magic links, MFA, row-level security |
| **Container** | Docker + Docker Compose | Consistent dev/staging/prod environments |
| **Orchestration** | Kubernetes (production) | Auto-scaling, rolling deployments |
| **CI/CD** | GitHub Actions | Integrated with GitHub, free for open source |
| **Mobile CI** | Expo EAS Build | Cloud builds for iOS/Android, OTA updates |
| **Monitoring** | Prometheus + Grafana | Metrics collection + visualization |
| **Logging** | Structured JSON (zerolog) | High-performance structured logging |
| **API Docs** | Swagger/OpenAPI 3.1 | Auto-generated, interactive docs |

### 2. Architecture Decisions

#### Clean Architecture (Mandatory)
```
┌─────────────────────────────────────────────┐
│                  Handlers                    │  ← HTTP/gRPC endpoints
├─────────────────────────────────────────────┤
│                  Use Cases                   │  ← Business logic orchestration
├─────────────────────────────────────────────┤
│                  Domain                      │  ← Entities, Value Objects, Interfaces
├─────────────────────────────────────────────┤
│               Infrastructure                 │  ← DB, Cache, External APIs
└─────────────────────────────────────────────┘
```

**Rules:**
- Dependencies point INWARD only (outer layers depend on inner layers)
- Domain layer has ZERO external dependencies
- Use Cases define interfaces, Infrastructure implements them
- Handlers only call Use Cases, never Infrastructure directly

#### Platform Strategy (Web + Mobile)
```
┌──────────────────────────────────────────────────┐
│                SHARED LAYER                       │
│  ┌──────────┐  ┌──────────┐  ┌────────────────┐ │
│  │ Types/   │  │ API      │  │ Business Logic │ │
│  │ Models   │  │ Hooks    │  │ (shared)       │ │
│  └──────────┘  └──────────┘  └────────────────┘ │
├──────────────────┬───────────────────────────────┤
│   WEB (React)    │   MOBILE (Expo React Native)  │
│   Tailwind v4    │   NativeWind / RN StyleSheet   │
│   React Router   │   Expo Router (file-based)     │
│   Vite 7+        │   Expo SDK 53+                 │
└──────────────────┴───────────────────────────────┘
```

#### Module Structure (Domain-Driven)
```
modules/
├── athlete/          # VĐV (Vận động viên) management
├── club/             # CLB (Câu lạc bộ) management
├── tournament/       # Giải đấu management
├── federation/       # Liên đoàn management (national + provincial)
├── btc/              # Ban tổ chức (Organizing Committee)
├── parent/           # Phụ huynh (Parent/Guardian)
├── scoring/          # Live scoring & results
├── finance/          # Tài chính (Financial management)
├── report/           # Báo cáo (Reporting & analytics)
└── notification/     # Thông báo (Notifications)
```

### 3. Technical Roadmap Template

```markdown
## Q1 - Foundation
- [ ] Project scaffold & CI/CD
- [ ] Core domain models
- [ ] Supabase Auth + RBAC system
- [ ] Athlete & Club modules
- [ ] Neon database provisioning + branching strategy

## Q2 - Tournament Core
- [ ] Tournament management
- [ ] Registration workflow
- [ ] Live scoring engine
- [ ] Federation module
- [ ] Expo React Native mobile app (MVP)

## Q3 - Advanced Features
- [ ] Real-time notifications (Supabase Realtime)
- [ ] Reporting & analytics
- [ ] Financial management
- [ ] Mobile app full feature parity

## Q4 - Scale & Polish
- [ ] Performance optimization
- [ ] Kubernetes deployment
- [ ] Monitoring & alerting
- [ ] Security audit
- [ ] App Store / Google Play release
```

### 4. Decision Framework

When evaluating technology or architecture decisions, use this checklist:

- [ ] **Performance**: Does it handle 10,000+ concurrent users during tournaments?
- [ ] **Scalability**: Can it scale horizontally for peak events?
- [ ] **Maintainability**: Can a mid-level Go dev understand it in < 1 day?
- [ ] **Security**: Does it meet OWASP Top 10 requirements?
- [ ] **Cost**: Is it cost-effective for a sports federation budget?
- [ ] **Cross-Platform**: Does it work on web AND mobile (Expo)?
- [ ] **Team Fit**: Does the team have expertise or can learn in < 2 weeks?
- [ ] **Vietnam Context**: Does it work with Vietnam's internet infrastructure?

### 5. Code Quality Standards

| Metric | Target | Tool |
|--------|--------|------|
| Test Coverage | ≥ 80% | `go test -cover` |
| Cyclomatic Complexity | ≤ 10 | `gocyclo` |
| Code Duplication | ≤ 3% | `jscpd` |
| Lint Score | 0 errors | `golangci-lint` |
| TypeScript Strict | Enabled | `tsconfig.json` |
| Bundle Size (Web) | < 500KB gzipped | `vite-bundle-visualizer` |
| Mobile App Size | < 50MB | Expo EAS Build |

### 6. Team Structure

```
CTO/Tech Lead
├── System Architect
├── Project Manager
│   └── Business Analyst
├── Backend Team
│   ├── Senior Go Developer
│   └── Go Developer(s)
├── Frontend Team
│   ├── Senior React/Expo Developer
│   └── React/Expo Developer(s)
├── DevOps Engineer
├── DBA
├── QA Engineer
├── Security Engineer
└── UI/UX Designer
```

### 7. Meeting Cadence

| Meeting | Frequency | Duration | Attendees |
|---------|-----------|----------|-----------|
| Daily Standup | Daily | 15 min | All devs |
| Sprint Planning | Bi-weekly | 2 hrs | PM, BA, Leads |
| Architecture Review | Weekly | 1 hr | CTO, Architect, Seniors |
| Code Review Sync | 2x/week | 30 min | All devs |
| Retrospective | Bi-weekly | 1 hr | All team |
| Stakeholder Demo | Bi-weekly | 1 hr | CTO, PM, Stakeholders |

### 8. Risk Assessment Matrix

| Risk | Impact | Probability | Mitigation |
|------|--------|------------|------------|
| DB scaling issues | High | Medium | Neon autoscaling, read replicas, connection pooling |
| Real-time latency | High | Medium | Supabase Realtime + Redis pub/sub |
| Security breach | Critical | Low | Supabase Auth, RLS, regular audits, WAF |
| Team turnover | High | Medium | Documentation, knowledge sharing |
| Vendor lock-in | Medium | Low | Abstract infrastructure layer, standard PostgreSQL |
| Mobile app rejection | Medium | Low | Follow App Store/Play Store guidelines early |
