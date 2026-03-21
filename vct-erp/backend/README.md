# VCT Core Ledger Backend

Core ledger scaffold for VCT Group, centered around a double-entry posting flow with:

- PostgreSQL schema for chart of accounts, journal headers/items, real-time balances, and outbox events
- Clean Architecture boundaries under `internal/modules/ledger`
- Finance business adapters for SaaS deferred revenue, dojo receivables, and rental deposits
- Transactional `PostEntry` use case with balance validation and outbox enqueue
- Redis-oriented cache and stream publisher interfaces for chart-of-accounts caching and downstream propagation
- VAS-oriented voucher numbering, reverse-entry flow, MongoDB audit trail, and bank reconciliation primitives

## Layout

```text
backend/
├── cmd/api
├── cmd/outbox-relay
├── migrations
├── sql/reports
├── internal/config
├── internal/infra
├── internal/httpapi
├── internal/modules/finance
├── internal/modules/ledger
│   ├── adapter
│   │   ├── cache
│   │   ├── outbox
│   │   └── postgres
│   ├── domain
│   └── usecase
├── internal/shared
│   ├── id
│   └── repository
└── schema.sql
```

## Notes

- `cmd/api` now composes the pgx-backed PostgreSQL connection, go-redis cache/stream publisher, and MongoDB audit collection.
- Finance capture routes are wired at `POST /v1/finance/capture`, while privileged reversal is exposed under `POST /v1/finance/journal-entries/{id}/void`.
- The posting workflow persists the journal entry, journal lines, balance deltas, and outbox row in a single transaction. Redis Stream publication is attempted after commit; if publication fails, the outbox row remains pending for a relay worker.
- Voucher numbering follows the VAS-style formats `PT-0001/03-26`, `PC-0001/03-26`, and `PK-0001/03-26`.
- The outbox relay worker source lives under `cmd/outbox-relay` and drains Redis Stream events with an at-least-once guarantee plus graceful batch completion on shutdown.
- A mermaid flow diagram for the cross-business money movement is available in `docs/finance-architecture.md`.
- SQL report templates live under `sql/reports` and use Postgres window functions for trial balance, B02-DN P&L, and general journal outputs.

## Runtime Variables

- `APP_ROLE_HEADER`, `APP_ACTOR_HEADER`, and `IDEMPOTENCY_HEADER` control the inbound headers used by the finance capture and privileged void APIs.
- `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, and `DB_CONN_MAX_LIFETIME` tune the pgx-backed connection pool.
- `MONGO_URI`, `MONGO_DATABASE`, and `MONGO_AUDIT_COLLECTION` configure the inspection-grade audit log store.
- `OUTBOX_RELAY_INTERVAL`, `OUTBOX_RELAY_RETRY_DELAY`, `OUTBOX_RELAY_CLAIM_TTL`, and `OUTBOX_RELAY_BATCH_SIZE` control the relay worker cadence and retry behavior.
- `OUTBOX_RELAY_RUN_TIMEOUT` bounds one relay batch so pod shutdowns can drain cleanly on Kubernetes.
