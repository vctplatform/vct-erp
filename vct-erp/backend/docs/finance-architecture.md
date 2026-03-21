# Finance Flow Architecture

```mermaid
flowchart LR
    A[Member Unit\nSaaS / Dojo / Rental] -->|POST /v1/finance/capture\nIdempotency-Key| B[Finance HTTP Adapter]
    B --> C[Capture Use Case]
    C --> D[Idempotency Store]
    C --> E[SaaS Service]
    C --> F[Dojo Service]
    C --> G[Rental Service]

    E --> H[Core Ledger PostEntry]
    F --> H
    G --> H

    E --> I[(saas_contracts + saas_revenue_schedules)]
    F --> J[(dojo_receivables)]
    G --> K[(rental_deposits)]

    H --> L[(journal_entries + journal_items)]
    H --> M[(account_balances)]
    H --> N[(outbox_events)]

    O[Outbox Relay Worker] -->|claim pending / failed| N
    O -->|publish at-least-once| P[Redis Stream / Kafka]

    P --> Q[Dojo Module]
    P --> R[SaaS Module]
    P --> S[Retail / Rental Module]
```

## Notes

- `Capture Use Case` enforces idempotency before dispatching into a business-specific adapter.
- `Core Ledger PostEntry` remains the single posting gate for all business lines.
- `outbox_events` is the boundary between committed accounting data and cross-module propagation.
- The relay worker uses claim-and-retry semantics to keep delivery at-least-once.
