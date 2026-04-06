# VCT ERP Domain Knowledge

> Kiến thức domain cốt lõi về hệ thống ERP cho VCT Group.

---

## VCT Group — Tổng quan

### Mô hình Kinh doanh
```
VCT Group (Công ty CP VCT)
├── VCT Platform (SaaS)
│   ├── Federation Management (Quản lý Liên đoàn)
│   ├── Dojo Management (Quản lý Võ đường)
│   ├── Athlete Management (Quản lý VĐV)
│   ├── Tournament Scoring (Chấm điểm thi đấu)
│   └── Digital Certificates (Chứng chỉ số)
├── Equipment Rental (Cho thuê thiết bị)
│   ├── Scoring devices
│   └── Competition mats
└── Consulting (Tư vấn chuyển đổi số)
```

### Dòng Doanh thu
| Stream | Type | Frequency | Account |
|--------|------|-----------|---------|
| SaaS Subscription | Recurring | Monthly/Annual | 5113 (DT CCDV) |
| Dojo Fees | Recurring | Monthly | 5113 (DT CCDV) |
| Equipment Rental | Recurring | Per-event | 5113 (DT CCDV) |
| Tournament Revenue | One-time | Per-event | 5111 (DT bán HH) |
| Consulting | Project | Per-project | 5113 (DT CCDV) |

---

## ERP Module Map

### Hiện trạng (As of 2026-03)

```
┌──────────────────────────────────────────────────┐
│                    CEO DASHBOARD                  │
│          (KPIs, Analytics, OKR Tracking)          │
├───────────────┬────────────────┬──────────────────┤
│   FINANCE     │    COMMERCIAL  │     PEOPLE       │
│   ┌────────┐  │  ┌──────────┐  │  ┌────────────┐  │
│   │ Ledger │  │  │  Sales   │  │  │     HR     │  │
│   │ ✅ Built│  │  │ 📋 Plan  │  │  │  📋 Plan   │  │
│   ├────────┤  │  ├──────────┤  │  ├────────────┤  │
│   │ Finance│  │  │ Marketing│  │  │  Payroll   │  │
│   │ ✅ Built│  │  │ 📋 Plan  │  │  │  📋 Plan   │  │
│   ├────────┤  │  ├──────────┤  │  ├────────────┤  │
│   │ Reports│  │  │   CRM    │  │  │ Recruitment│  │
│   │ 🔨 WIP │  │  │ 📋 Plan  │  │  │  📋 Plan   │  │
│   └────────┘  │  └──────────┘  │  └────────────┘  │
├───────────────┴────────────────┴──────────────────┤
│              PLATFORM (DevOps, Security)           │
│              BACKEND ENGINE (Go, PostgreSQL)        │
│              FRONTEND CRAFT (Next.js, TypeScript)   │
└──────────────────────────────────────────────────┘
```

### Finance Module (✅ Built)
```
Core Tables:
├── accounts: Chart of Accounts (VAS TT200)
├── journal_entries: Transaction headers
├── journal_items: Transaction lines (partitioned quarterly)
├── account_balances: Running balance cache
├── voucher_sequences: VAS voucher numbering
├── outbox_events: Reliable event publishing
├── idempotency_keys: Safe retry mechanism

Business Adapters:
├── saas_contracts: SaaS subscription tracking
├── saas_revenue_schedules: Deferred revenue recognition
├── dojo_receivables: Dojo fee billing/collection
├── rental_deposits: Equipment rental deposits
└── bank_statement_lines: Bank reconciliation

Views:
├── v_finance_reconciliation: Account balance audit
└── v_gross_profit_by_segment: Segment profitability
```

---

## Vietnamese Accounting Standards (VAS)

### Key Regulations
- **TT200/2014/TT-BTC**: Chế độ kế toán doanh nghiệp (primary)
- **NĐ 174/2016/NĐ-CP**: Hướng dẫn luật kế toán
- **Luật Kế toán 88/2015**: Luật kế toán hiện hành

### Fiscal Year
- Vietnam: Jan 1 — Dec 31 (default)
- Tax filing: Quarterly CIT advance, Annual CIT final

### E-Invoice Requirements (NĐ 123/2020)
- Mandatory electronic invoices from 2022
- Must have CQT (Tax Authority) code
- XML format with digital signature
- Real-time verification via tax portal

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
