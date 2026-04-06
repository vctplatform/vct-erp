---
name: frontend-craft
description: >-
  Mega-Skill for Next.js frontend development. App Router, TypeScript, shadcn/ui,
  TailwindCSS, data visualization, accessibility, and ERP dashboard patterns.
metadata:
  author: VCT Platform
  version: "1.0.0"
  type: "Mega-Skill"
  locale: vi-VN
---

# FRONTEND-CRAFT — MEGA-SKILL

> Next.js frontend excellence. Beautiful dashboards, data-dense interfaces, premium UX.

---

## 🔹 NĂNG LỰC: FRONTEND-DEVELOPER

> *"The best interface is no interface. The second best is one that feels invisible."*

### Technology Stack
- **Framework**: Next.js 15+ (App Router, Server Components)
- **Language**: TypeScript 5+ (strict mode)
- **UI Components**: shadcn/ui (Radix primitives)
- **Styling**: TailwindCSS 4+
- **State**: React Context / Zustand (minimal global state)
- **Data Fetching**: Server Components + React Suspense
- **Charts**: Recharts / Tremor
- **Tables**: TanStack Table v8
- **Forms**: React Hook Form + Zod

### Project Structure (VCT ERP)
```
frontend/
├── app/
│   ├── (dashboard)/        → Dashboard layout group
│   │   ├── layout.tsx      → Sidebar + Header
│   │   ├── page.tsx        → Home/Overview
│   │   ├── ledger/         → General Ledger
│   │   ├── reports/        → Financial Reports
│   │   ├── reconciliation/ → Bank Reconciliation
│   │   ├── segments/       → Business Segments
│   │   └── control-room/   → System Control
│   ├── api/                → API routes (proxy)
│   ├── layout.tsx          → Root layout
│   └── globals.css         → Global styles
├── components/
│   ├── ui/                 → shadcn/ui primitives
│   ├── layout/             → Sidebar, Header, Footer
│   ├── dashboard/          → Dashboard widgets
│   ├── forms/              → Form components
│   └── data-display/       → Tables, Charts, Cards
├── hooks/                  → Custom hooks
├── lib/                    → Utilities, API client
└── types/                  → TypeScript type definitions
```

### Component Patterns

#### Server Component (Default)
```tsx
// app/(dashboard)/ledger/page.tsx
export default async function LedgerPage() {
  const entries = await getJournalEntries()
  return <JournalTable entries={entries} />
}
```

#### Client Component (When Needed)
```tsx
'use client'
// Only for: interactivity, state, browser APIs, event handlers
export function JournalEntryForm() {
  const form = useForm<JournalEntryInput>({
    resolver: zodResolver(journalEntrySchema),
  })
  // ...
}
```

#### Data Fetching Pattern
```tsx
// Parallel data fetching with Suspense
export default async function DashboardPage() {
  const [balances, recentEntries, kpis] = await Promise.all([
    getAccountBalances(),
    getRecentEntries({ limit: 10 }),
    getDashboardKPIs(),
  ])
  return (
    <div className="grid grid-cols-12 gap-6">
      <KPICards kpis={kpis} />
      <BalanceChart balances={balances} />
      <RecentEntriesTable entries={recentEntries} />
    </div>
  )
}
```

### API Client Pattern
```tsx
const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  })
  if (!res.ok) {
    const error = await res.json()
    throw new APIError(error.error.code, error.error.message)
  }
  return res.json()
}
```

---

## 🔹 NĂNG LỰC: UI-UX-DESIGNER

### Design System (VCT ERP)
```
Tokens:
├── Colors: Primary (brand), Semantic (success/warning/error/info)
├── Typography: Inter (primary), JetBrains Mono (code/numbers)
├── Spacing: 4px grid (4, 8, 12, 16, 24, 32, 48, 64)
├── Shadows: sm, md, lg, xl
├── Borders: radius (4, 8, 12, full), width (1, 2)
└── Motion: ease-out 150ms (micro), 300ms (standard), 500ms (emphasis)

Components:
├── Atoms: Button, Input, Badge, Avatar, Icon
├── Molecules: Search bar, Card, Nav item, Form field, Stat card
├── Organisms: Sidebar, Header, Data table, Modal, Journal entry form
└── Templates: Dashboard layout, Report layout, Settings page
```

### ERP-Specific UI Patterns

#### Large Data Tables
```
├── Server-side pagination (cursor-based)
├── Column sorting, filtering, search
├── Row selection for bulk actions
├── Sticky header and first column
├── Responsive: horizontal scroll on mobile
├── Export to CSV/PDF
└── Print-friendly layout
```

#### Financial Number Formatting
```
├── VND: No decimals, dot thousands separator (1.250.000)
├── Debit: Black text, positive
├── Credit: Red text or parentheses
├── Zero: Dash "—" not "0"
├── Monospace font for alignment (JetBrains Mono)
└── Right-aligned in tables
```

#### Dashboard Widgets
```
├── KPI Card: Value, Trend (↑↓→), Period comparison
├── Line Chart: Revenue/Expense over time
├── Bar Chart: By segment/department
├── Pie/Donut: composition breakdowns
├── Table Widget: Top 10 items with sparklines
└── Calendar Heat Map: Activity density
```

### Accessibility (WCAG 2.1 AA)
| Criterion | Requirement |
|-----------|------------|
| Color Contrast | 4.5:1 text, 3:1 large |
| Keyboard | All interactive reachable via Tab |
| Screen Reader | Semantic HTML + aria labels |
| Focus | Visible :focus-visible ring |
| Touch Targets | Min 44×44px |
| Motion | Respect prefers-reduced-motion |

### Responsive Breakpoints
```
├── Mobile: < 640px (sm)
├── Tablet: 640-1024px (md)
├── Desktop: 1024-1280px (lg)
└── Wide: > 1280px (xl)

ERP Dashboard: Desktop-first (primary users on desktop)
Mobile: Read-only dashboards, approvals
```

---

## 🔹 NĂNG LỰC: DATA-VISUALIZATION

### Chart Selection Guide (ERP)
| Data Type | Chart | Library |
|-----------|-------|---------|
| Time series | Line/Area | Recharts |
| Comparison | Bar (horizontal/vertical) | Recharts |
| Composition | Donut/Stacked bar | Recharts |
| Distribution | Histogram | Recharts |
| KPI | Stat card + sparkline | Custom |
| Table data | Data table | TanStack Table |
| Hierarchy | Treemap | Recharts |

### Report Layout (Print-Ready)
```
VAS Financial Report Format:
├── Header: Company name, report title, period
├── Body: Data table with proper formatting
├── Footer: Signatures (Kế toán trưởng, Giám đốc)
├── Page: A4 portrait/landscape
├── Print CSS: @media print { ... }
└── PDF export: via browser print or jsPDF
```

---

## Development Checklist

- [ ] TypeScript strict mode, no `any`
- [ ] Server Components by default
- [ ] Client Components only when necessary
- [ ] Proper loading.tsx and error.tsx
- [ ] Responsive design tested
- [ ] Keyboard navigation works
- [ ] Color contrast passes AA
- [ ] Financial numbers properly formatted
- [ ] Tables handle empty states gracefully
- [ ] Forms validate before submit

## Trigger Patterns

- "frontend", "UI", "giao diện", "dashboard", "page"
- "component", "layout", "form", "table", "chart"
- "design", "thiết kế", "UX", "responsive"
- "Next.js", "React", "TypeScript", "TailwindCSS"

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
