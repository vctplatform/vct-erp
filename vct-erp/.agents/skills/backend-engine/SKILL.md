---
name: backend-engine
description: >-
  Mega-Skill for Go backend development. Clean Architecture, REST API, PostgreSQL,
  Redis, MongoDB, testing, and database patterns for VCT ERP.
metadata:
  author: VCT Platform
  version: "1.0.0"
  type: "Mega-Skill"
  locale: vi-VN
---

# BACKEND-ENGINE — MEGA-SKILL

> Go backend powerhouse. Clean Architecture, database mastery, API excellence.

---

## 🔹 NĂNG LỰC: BACKEND-DEVELOPER

> *"Write code for humans first, computers second."*

### Technology Stack
- **Language**: Go 1.26 (range-over-func, enhanced generics, iterator patterns)
- **Router**: Chi v5 / stdlib enhanced ServeMux
- **DB Driver**: pgx v5 (PostgreSQL 18+)
- **Cache**: go-redis v9
- **Logger**: zerolog / slog
- **Validator**: go-playground/validator v10
- **Testing**: stdlib testing + testify

### Project Layout (VCT ERP)
```
cmd/
├── api/main.go              → HTTP server entry
└── outbox-relay/main.go     → Background worker
internal/
├── config/                  → Configuration
├── infra/                   → Infrastructure adapters
├── httpapi/                 → HTTP routing & middleware
├── modules/
│   └── [module]/
│       ├── domain/          → Entities, interfaces, errors
│       ├── usecase/         → Business logic
│       └── adapter/
│           ├── postgres/    → DB implementation
│           ├── cache/       → Redis cache
│           └── outbox/      → Event outbox
└── shared/                  → Shared utilities
    ├── id/                  → UUID helpers
    └── repository/          → Base repository
```

### Entity Pattern
```go
type Entity struct {
    ID        uuid.UUID `json:"id"`
    // fields with validation tags
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Repository Interface Pattern
```go
type Repository interface {
    Create(ctx context.Context, entity *Entity) error
    GetByID(ctx context.Context, id uuid.UUID) (*Entity, error)
    Update(ctx context.Context, entity *Entity) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, filter Filter) ([]Entity, int64, error)
    // Go 1.26 iterator for streaming
    Stream(ctx context.Context, filter Filter) iter.Seq2[Entity, error]
}
```

### Use Case Pattern
```go
type CreateUseCase struct {
    repo domain.Repository
}

func (uc *CreateUseCase) Execute(ctx context.Context, input Input) (*Entity, error) {
    entity := input.ToEntity()
    if err := uc.repo.Create(ctx, entity); err != nil {
        return nil, fmt.Errorf("create: %w", err)
    }
    return entity, nil
}
```

### HTTP Handler Pattern
```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    }
    if err := req.Validate(); err != nil {
        writeValidationError(w, err)
        return
    }
    result, err := h.useCase.Execute(r.Context(), req.ToInput())
    if err != nil {
        writeError(w, mapDomainError(err))
        return
    }
    writeJSON(w, http.StatusCreated, result)
}
```

### Error Handling Standards
```go
// Domain errors
var (
    ErrNotFound      = errors.New("not found")
    ErrDuplicate     = errors.New("duplicate")
    ErrUnbalanced    = errors.New("journal entry unbalanced")
)

// Wrap with context
fmt.Errorf("post journal entry %s: %w", entryID, err)
```

### Testing Standards
```go
func TestUseCase(t *testing.T) {
    tests := []struct {
        name    string
        input   Input
        mock    func(*MockRepo)
        wantErr bool
    }{
        // table-driven tests
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // arrange → act → assert
        })
    }
}
```

---

## 🔹 NĂNG LỰC: DATABASE-ARCHITECT

### PostgreSQL Mastery (VCT ERP Patterns)

#### Schema Design Rules
```
├── UUIDs as primary keys (gen_random_uuid())
├── created_at/updated_at on EVERY table
├── Trigger for auto-updating updated_at
├── ENUM types for status fields
├── CHECK constraints for data integrity
├── UNIQUE constraints for business rules
├── Foreign keys with appropriate ON DELETE
├── JSONB for flexible metadata
└── Partitioning for high-volume tables
```

#### VCT ERP Specific Patterns
```
Double-Entry Accounting:
├── journal_entries: Header (entry_no, status, posting_date)
├── journal_items: Lines (account_id, debit/credit, amount)
├── account_balances: Running balances (optimistic)
├── Constraint: SUM(debit) = SUM(credit) per entry (DEFERRED)
├── Partition: journal_items BY RANGE (created_at) quarterly
└── Voucher: VAS-style PT-0001/03-26, PC-0001/03-26, PK-0001/03-26

Outbox Pattern:
├── outbox_events: Reliable event publishing
├── Status: pending → processing → published | failed
├── Relay worker drains via Redis Stream
└── At-least-once guarantee

Idempotency:
├── idempotency_keys: Prevent duplicate processing
├── Scope + key unique constraint
└── Request hash for drift detection
```

#### Query Optimization
```
├── Use EXPLAIN ANALYZE for every critical query
├── Create indexes for WHERE, JOIN, ORDER BY columns
├── Use BRIN indexes for time-series data
├── Use partial indexes (WHERE clause) for status filters
├── Use cursor-based pagination for large datasets
├── Avoid SELECT * — specify columns
├── Use CTEs for readability, subqueries for performance
└── Connection pooling: pgxpool with Neon-optimized settings
```

#### Migration Rules
```
├── One migration per schema change
├── Always include UP and DOWN
├── Never modify existing migrations
├── Test migration on staging before production
├── Use transactions for DDL when possible
├── Add NOT NULL with DEFAULT to avoid table locks
└── Create indexes CONCURRENTLY on large tables
```

---

## 🔹 NĂNG LỰC: DATA-ENGINEER

### Redis Patterns (VCT ERP)
```
Cache:
├── Chart of accounts cache (read-heavy)
├── Session storage
├── Rate limiting counters
└── TTL-based expiration

Streams:
├── Journal entry events (outbox relay)
├── Balance update notifications
├── Cross-module event propagation
└── Consumer groups for reliable delivery
```

### Outbox Relay Worker
```
Cycle:
├── Poll outbox_events WHERE status = 'pending'
├── Claim batch (SET status = 'processing', locked_by)
├── Publish to Redis Stream
├── Mark published (SET status = 'published')
├── On failure: increment attempt_count, reset status
└── Graceful shutdown: drain current batch
```

---

## Development Checklist

- [ ] Clean Architecture boundaries enforced
- [ ] All exported functions have doc comments
- [ ] Unit tests for use cases (mock repository)
- [ ] Integration tests for repository (test DB)
- [ ] `golangci-lint` zero errors
- [ ] SQL queries parameterized (no injection)
- [ ] Context propagation in all DB/Redis calls
- [ ] Error wrapping with `fmt.Errorf("ctx: %w", err)`
- [ ] Proper HTTP status codes
- [ ] Request validation before business logic

## Trigger Patterns

- "backend", "API", "endpoint", "Go", "golang"
- "database", "schema", "migration", "query", "SQL"
- "Redis", "cache", "stream", "outbox"
- "repository", "use case", "handler", "middleware"

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
