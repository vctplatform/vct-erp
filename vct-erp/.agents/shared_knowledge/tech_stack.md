# Technology Stack — VCT ERP

> Tài liệu này ghi nhận các quyết định công nghệ đã được chấp thuận cho dự án ERP.
> Cập nhật lần cuối: 2026-03-29

---

## Core Stack

### Backend
| Component | Choice | Version | Rationale |
|-----------|--------|---------|-----------|
| Language | **Go** | 1.26 | Performance, simplicity, strong stdlib, excellent concurrency |
| HTTP Router | **Chi v5** / stdlib | Latest | Lightweight, idiomatic, middleware chain support |
| Database | **PostgreSQL** | 18+ | ACID compliance, VAS accounting needs, JSON support, partitioning |
| DB Driver | **pgx v5** | Latest | Best Go PostgreSQL driver, native types, pooling |
| Cache | **Redis** | 7+ | Session, rate limiting, pub/sub, streams for outbox relay |
| Audit DB | **MongoDB** | 7+ | Flexible schema for inspection-grade audit logs |
| Logger | **zerolog** | Latest | Zero-allocation JSON logger, best performance |
| Validation | **go-playground/validator** | v10 | Struct tag validation, custom validators |

### Frontend
| Component | Choice | Version | Rationale |
|-----------|--------|---------|-----------|
| Framework | **Next.js** | 15+ | Server Components, App Router, SSR/SSG, React ecosystem |
| Language | **TypeScript** | 5+ | Type safety, better DX, catch bugs at compile time |
| UI Library | **shadcn/ui** | Latest | Radix primitives, accessible, customizable, no vendor lock |
| Styling | **TailwindCSS** | 4+ | Utility-first, consistent design, great DX |
| Tables | **TanStack Table** | v8 | Headless, powerful sorting/filtering/pagination |
| Charts | **Recharts** | Latest | React-native charts, good for financial data |
| Forms | **React Hook Form + Zod** | Latest | Performance, validation, type inference |

### Infrastructure
| Component | Choice | Rationale |
|-----------|--------|-----------|
| Container | **Docker** | Industry standard, reproducible builds |
| CI/CD | **GitHub Actions** | Already using GitHub, good marketplace |
| Monitoring | **Prometheus + Grafana** | Open source, powerful alerting |
| Logging | **Loki** | Log aggregation, integrates with Grafana |

---

## Architecture Decisions

### ADR-001: Clean Architecture (Accepted)
- **Decision**: Follow Clean Architecture with domain/usecase/adapter layers
- **Rationale**: Clear boundaries, testable, swappable implementations
- **Trade-off**: More files/boilerplate, but better maintainability

### ADR-002: Double-Entry Accounting (Accepted)
- **Decision**: Implement proper double-entry bookkeeping with debit/credit balance checking
- **Rationale**: VAS TT200 compliance, financial accuracy, audit requirements
- **Trade-off**: More complex than simple transaction recording

### ADR-003: Outbox Pattern (Accepted)
- **Decision**: Use transactional outbox for reliable event publishing
- **Rationale**: Guarantee consistency between DB writes and event publishing
- **Trade-off**: Additional complexity, relay worker needed

### ADR-004: Table Partitioning (Accepted)
- **Decision**: Partition journal_items by quarter using RANGE on created_at
- **Rationale**: Performance for time-range queries, manageable partition sizes
- **Trade-off**: More complex queries, partition maintenance needed

### ADR-005: UUID Primary Keys (Accepted)
- **Decision**: Use UUID v4 for all table primary keys
- **Rationale**: No conflicts in distributed systems, no information leakage
- **Trade-off**: Larger storage, cannot sort by ID for time ordering

---

## Coding Conventions

### Go
```
├── Package naming: lowercase, no underscores (domain, usecase, postgres)
├── Exported functions: PascalCase with doc comments
├── Internal functions: camelCase
├── Error handling: Always wrap with context fmt.Errorf("verb noun: %w", err)
├── Context: First parameter for all I/O functions
├── Tests: filename_test.go, table-driven tests
├── Linting: golangci-lint with default + custom rules
└── Formatting: gofmt (non-negotiable)
```

### TypeScript/React
```
├── Components: PascalCase (JournalTable.tsx)
├── Hooks: camelCase with "use" prefix (useJournalEntries.ts)
├── Utils: camelCase (formatCurrency.ts)
├── Types: PascalCase (JournalEntry, AccountBalance)
├── API responses: snake_case (matching Go JSON tags)
├── CSS: TailwindCSS utility classes
├── Server Components: default (no 'use client')
└── Client Components: only when needed (add 'use client')
```

### SQL
```
├── Table names: snake_case, plural (journal_entries)
├── Column names: snake_case (posting_date)
├── Constraints: descriptive prefix (uq_, ck_, fk_)
├── Indexes: idx_{table}_{columns}
├── Enums: snake_case values ('posted', 'reversed')
├── Timestamps: TIMESTAMPTZ (always timezone-aware)
└── Money: NUMERIC(20, 4) (20 digits, 4 decimals)
```
