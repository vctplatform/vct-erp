---
name: system-architect
description: System Architect role - System design, API design, database schema, Clean Architecture patterns, scalability, and integration planning for VCT Platform.
---

# System Architect - VCT Platform

## Role Overview
Designs the overall system architecture ensuring scalability, maintainability, and alignment with Clean Architecture principles. Coordinates cross-platform (Web + Mobile) architecture and cloud-native service integration (Supabase, Neon).

## Core Responsibilities

### 1. System Architecture Overview

```
┌──────────────────────────────────────────────────────────┐
│                    CLIENT LAYER                           │
│  ┌─────────┐  ┌─────────┐  ┌──────────┐  ┌───────────┐ │
│  │ Web App │  │ PWA     │  │ Mobile   │  │ Admin     │ │
│  │(React20)│  │         │  │(Expo RN) │  │ Panel     │ │
│  └────┬────┘  └────┬────┘  └────┬─────┘  └─────┬─────┘ │
└───────┼─────────────┼───────────┼───────────────┼────────┘
        │             │           │               │
┌───────▼─────────────▼───────────▼───────────────▼────────┐
│                SUPABASE LAYER (BaaS)                      │
│  ┌──────────┐  ┌────────────┐  ┌─────────┐  ┌────────┐ │
│  │ Supabase │  │ Supabase   │  │Supabase │  │ Edge   │ │
│  │ Auth     │  │ Realtime   │  │ Storage │  │ Funcs  │ │
│  └──────────┘  └────────────┘  └─────────┘  └────────┘ │
└──────────────────────┬───────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────┐
│                   API GATEWAY (Go 1.26)                   │
│  ┌──────────┐  ┌────────────┐  ┌────────────────────┐   │
│  │ REST API │  │ WebSocket  │  │ Rate Limiter/RBAC  │   │
│  └──────────┘  └────────────┘  └────────────────────┘   │
└──────────────────────┬───────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────┐
│                  SERVICE LAYER                            │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│  │ Athlete  │ │ Club     │ │Tournament│ │Federation│   │
│  │ Service  │ │ Service  │ │ Service  │ │ Service  │   │
│  ├──────────┤ ├──────────┤ ├──────────┤ ├──────────┤   │
│  │ Scoring  │ │ Finance  │ │ Report   │ │ Notif.   │   │
│  │ Service  │ │ Service  │ │ Service  │ │ Service  │   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │
└──────────────────────┬───────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────┐
│                INFRASTRUCTURE LAYER                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐ │
│  │ Neon     │  │  Redis   │  │ Supabase │  │ SMTP    │ │
│  │ PG 18+  │  │  Cache   │  │ Storage  │  │ Email   │ │
│  │(serverl.)│  │          │  │ (S3)     │  │         │ │
│  └──────────┘  └──────────┘  └──────────┘  └─────────┘ │
└──────────────────────────────────────────────────────────┘
```

### 2. Clean Architecture - Module Structure

```
internal/
├── modules/
│   └── {module_name}/
│       ├── domain/
│       │   ├── entity.go          # Domain entities
│       │   ├── value_objects.go   # Value objects
│       │   ├── repository.go     # Repository interface
│       │   ├── service.go        # Domain service interface
│       │   └── errors.go         # Domain-specific errors
│       ├── usecase/
│       │   ├── create.go         # Create use case
│       │   ├── update.go         # Update use case
│       │   ├── delete.go         # Delete use case
│       │   ├── get.go            # Read use cases
│       │   └── list.go           # List with pagination
│       ├── adapter/
│       │   ├── postgres/
│       │   │   ├── repository.go # Neon PostgreSQL implementation
│       │   │   └── queries.go    # SQL queries
│       │   ├── supabase/
│       │   │   ├── auth.go       # Supabase Auth integration
│       │   │   └── storage.go    # Supabase Storage integration
│       │   ├── redis/
│       │   │   └── cache.go      # Redis cache implementation
│       │   └── http/
│       │       ├── handler.go    # HTTP handlers
│       │       ├── request.go    # Request DTOs
│       │       ├── response.go   # Response DTOs
│       │       ├── router.go     # Route definitions
│       │       └── middleware.go # Module-specific middleware
│       └── dto/
│           ├── request.go
│           └── response.go
├── shared/
│   ├── middleware/
│   │   ├── auth.go              # Supabase JWT verification
│   │   ├── cors.go
│   │   ├── logging.go
│   │   └── ratelimit.go
│   ├── errors/
│   │   └── errors.go            # Standardized error responses
│   ├── pagination/
│   │   └── pagination.go
│   └── validator/
│       └── validator.go
└── config/
    ├── config.go
    ├── neon.go                   # Neon connection config
    └── supabase.go               # Supabase config
```

### 3. API Design Standards

#### RESTful Conventions
```
GET    /api/v1/{resource}          # List (with pagination)
GET    /api/v1/{resource}/{id}     # Get by ID
POST   /api/v1/{resource}          # Create
PUT    /api/v1/{resource}/{id}     # Full update
PATCH  /api/v1/{resource}/{id}     # Partial update
DELETE /api/v1/{resource}/{id}     # Soft delete

# Nested resources
GET    /api/v1/clubs/{id}/athletes
POST   /api/v1/tournaments/{id}/registrations

# Actions
POST   /api/v1/tournaments/{id}/publish
POST   /api/v1/scoring/{race_id}/start
```

#### Standard Response Format
```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

#### Standard Error Format
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      { "field": "email", "message": "Invalid email format" }
    ]
  }
}
```

### 4. Database Schema Patterns

#### Naming Conventions
- Tables: `snake_case`, plural (e.g., `athletes`, `club_members`)
- Columns: `snake_case` (e.g., `first_name`, `created_at`)
- Primary keys: `id` (UUID)
- Foreign keys: `{referenced_table_singular}_id`
- Indexes: `idx_{table}_{columns}`
- Constraints: `{type}_{table}_{columns}` (e.g., `uq_athletes_email`)

#### Standard Columns (every table)
```sql
id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
deleted_at  TIMESTAMPTZ,         -- Soft delete
created_by  UUID REFERENCES auth.users(id),  -- Supabase auth
updated_by  UUID REFERENCES auth.users(id)
```

#### Multi-tenancy Pattern (Supabase RLS)
```sql
-- Provincial federation data isolation via Supabase RLS
ALTER TABLE athletes ENABLE ROW LEVEL SECURITY;

CREATE POLICY org_isolation ON athletes
    FOR ALL
    USING (
        organization_id IN (
            SELECT organization_id FROM user_organizations
            WHERE user_id = auth.uid()
        )
    );
```

### 5. Integration Patterns

| Pattern | Use Case | Implementation |
|---------|----------|---------------|
| **Supabase Realtime** | Live scoring updates | `postgres_changes` channel |
| **CQRS** | Read-heavy dashboards | Separate read/write models |
| **Pub/Sub** | Real-time notifications | Supabase Realtime + Redis |
| **Circuit Breaker** | External API calls | Retry with exponential backoff |
| **Saga** | Tournament registration | Distributed transaction coordination |
| **Neon Branching** | Preview environments | Per-PR database branches |

### 6. Scalability Checklist

- [ ] Neon autoscaling configured (auto-scale compute 0.25→8 CU)
- [ ] Redis caching for hot data (athlete profiles, rankings)
- [ ] Supabase Storage for static assets (images, documents)
- [ ] CDN for frontend assets (Cloudflare)
- [ ] Horizontal scaling via stateless Go services
- [ ] Neon read replicas for report queries
- [ ] WebSocket connection limits and heartbeat (Supabase Realtime)
- [ ] File uploads via Supabase Storage (pre-signed URLs)
- [ ] Background job queue for heavy processing (email, reports)
- [ ] Expo OTA updates for mobile app patches
