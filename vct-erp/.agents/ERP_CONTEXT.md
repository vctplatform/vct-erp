# VCT ERP — Project Context

> Tài liệu này mô tả bối cảnh dự án ERP để mọi agent hiểu đúng scope và constraints.

---

## 1. MỤC ĐÍCH DỰ ÁN

Xây dựng hệ thống **ERP toàn diện** cho VCT Group — nền tảng quản trị doanh nghiệp tích hợp:
- Tài chính & Kế toán (VAS-compliant)
- Nhân sự & Chấm công
- Kinh doanh & CRM
- Kho vận & Mua sắm
- Dashboard & Báo cáo điều hành

## 2. TECH STACK

### Backend
| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.26 |
| HTTP Router | Chi v5 / stdlib | Latest |
| Database | PostgreSQL | 18+ |
| DB Driver | pgx v5 | Latest |
| Cache | Redis | 7+ |
| Audit | MongoDB | 7+ |
| Logger | zerolog / slog | Latest |
| Config | envconfig | Latest |

### Frontend
| Component | Technology | Version |
|-----------|-----------|---------|
| Framework | Next.js | 15+ |
| Language | TypeScript | 5+ |
| UI Library | shadcn/ui | Latest |
| Styling | TailwindCSS | 4+ |
| State | React Context / Zustand | Latest |
| Charts | Recharts / Tremor | Latest |

### Infrastructure
| Component | Technology |
|-----------|-----------|
| Container | Docker |
| Orchestration | Kubernetes |
| CI/CD | GitHub Actions |
| Monitoring | Prometheus + Grafana |
| Logging | Loki |

## 3. ARCHITECTURE

### Clean Architecture Pattern
```
cmd/                    ← Entry points
├── api/               ← HTTP server
└── outbox-relay/      ← Background worker

internal/
├── config/            ← Configuration
├── infra/             ← Infrastructure adapters
├── httpapi/           ← HTTP routing
├── modules/           ← Business modules
│   ├── finance/       ← Finance capture & void
│   ├── ledger/        ← Core double-entry engine
│   │   ├── domain/    ← Entities, interfaces
│   │   ├── usecase/   ← Business logic
│   │   └── adapter/   ← Implementations
│   │       ├── postgres/
│   │       ├── cache/
│   │       └── outbox/
│   └── analytics/     ← Reporting & BI
└── shared/            ← Shared utilities
    ├── id/
    └── repository/
```

### Database Architecture
- **Double-entry accounting** (debit = credit always)
- **VAS TT200 compliant** chart of accounts
- **Partitioned journal_items** by quarter
- **Outbox pattern** for reliable event publishing
- **Idempotency keys** for safe retries
- **Multi-company** via `company_code` (default: VCT_GROUP)

### Key APIs
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/finance/capture` | POST | Create journal entry |
| `/v1/finance/journal-entries/{id}/void` | POST | Void/reverse entry |

### Voucher Types (VAS Style)
- `PT` — Phiếu Thu (Receipt)
- `PC` — Phiếu Chi (Payment)
- `PK` — Phiếu Kế toán (General Journal)

## 4. BUSINESS DOMAIN

### VCT Group — Công ty Cổ phần VCT
- **Ngành**: Công nghệ thể thao — Chuyển đổi số Võ Cổ Truyền Việt Nam
- **Sản phẩm chính**: VCT Platform (SaaS cho Liên đoàn, Võ đường, VĐV)
- **Doanh thu**: SaaS subscriptions + Dojo management fees + Equipment rental
- **Currency**: VND (primary)

### ERP Modules Roadmap
| # | Module | Priority | Status |
|---|--------|----------|--------|
| 1 | Finance & Accounting | 🔴 Critical | ✅ Core built |
| 2 | General Ledger | 🔴 Critical | ✅ Core built |
| 3 | Analytics & Reports | 🟡 High | 🔨 In progress |
| 4 | HR & Payroll | 🟡 High | 📋 Planned |
| 5 | Sales & CRM | 🟡 High | 📋 Planned |
| 6 | Procurement | 🟢 Medium | 📋 Planned |
| 7 | CEO Dashboard | 🟢 Medium | 📋 Planned |

## 5. CODING STANDARDS

### Go Backend
- Clean Architecture boundaries strictly enforced
- All exported functions have doc comments
- Unit tests > 80% coverage
- `golangci-lint` passes with zero errors
- SQL queries use parameterized inputs (no SQL injection)
- Context propagation in all DB calls
- Proper error wrapping: `fmt.Errorf("context: %w", err)`

### Next.js Frontend
- App Router with server components by default
- Client components only when needed (interactivity, state)
- TypeScript strict mode
- Responsive design (mobile-first)
- Accessibility (WCAG 2.1 AA)

### Database
- All tables have `created_at`, `updated_at` timestamps
- UUIDs as primary keys
- Constraints for data integrity
- Indexes for query performance
- Partitioning for large tables

## 6. LIÊN KẾT VỚI DỰ ÁN KHÁC

```
VCT PLATFORM/
├── vct-agent-business/  ← Jen quản lý (business operations)
├── vct-erp/             ← Jon quản lý (THIS PROJECT)
├── vct-platform/        ← Javis quản lý (main platform)
└── vct-cms-website/     ← Jack quản lý (marketing site)
```

Jon phối hợp với Jen qua Chairman khi cần cross-project coordination.
