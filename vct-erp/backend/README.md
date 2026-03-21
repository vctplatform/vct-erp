# VCT Core Ledger Backend

Core ledger scaffold for VCT Group, centered around a double-entry posting flow with:

- PostgreSQL schema for chart of accounts, journal headers/items, real-time balances, and outbox events
- Clean Architecture boundaries under `internal/modules/ledger`
- Finance business adapters for SaaS deferred revenue, dojo receivables, and rental deposits
- Transactional `PostEntry` use case with balance validation and outbox enqueue
- Redis-oriented cache and stream publisher interfaces for chart-of-accounts caching and downstream propagation

## Layout

```text
backend/
├── cmd/api
├── internal/config
├── internal/httpapi
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

- The HTTP bootstrap is intentionally minimal. `GET /healthz` is live and `POST /api/v1/ledger/journal-entries` is wired at the transport level, while infrastructure composition remains open for the real PostgreSQL driver and Redis client chosen by the platform.
- Finance capture routes are prepared at `POST /v1/finance/capture` with idempotency-key support and RBAC hooks for privileged void routes.
- The posting workflow persists the journal entry, journal lines, balance deltas, and outbox row in a single transaction. Redis Stream publication is attempted after commit; if publication fails, the outbox row remains pending for a relay worker.
- The outbox relay worker source lives under `cmd/outbox-relay`; compose it with the concrete Postgres and Redis implementations in the runtime layer.
- A mermaid flow diagram for the cross-business money movement is available in `docs/finance-architecture.md`.

## Runtime Variables

- `APP_ROLE_HEADER` and `IDEMPOTENCY_HEADER` control the inbound headers used by the finance capture API.
- `OUTBOX_RELAY_INTERVAL`, `OUTBOX_RELAY_RETRY_DELAY`, `OUTBOX_RELAY_CLAIM_TTL`, and `OUTBOX_RELAY_BATCH_SIZE` control the relay worker cadence and retry behavior.
