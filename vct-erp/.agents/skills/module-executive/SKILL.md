---
name: module-executive
description: >-
  Mega-Skill for CEO Dashboard, Analytics & Reporting ERP module. Real-time KPIs,
  financial analysis, cross-module reporting, OKR tracking, and business intelligence.
metadata:
  author: VCT Platform
  version: "1.0.0"
  type: "Mega-Skill"
  locale: vi-VN
---

# MODULE-EXECUTIVE — MEGA-SKILL

> Domain expertise cho Dashboard Điều hành, Analytics & Báo cáo BI.

---

## 🔹 NĂNG LỰC: CEO DASHBOARD

### KPI Framework (North Star + Input Metrics)
```
North Star: Monthly Recurring Revenue (MRR)

Input Metrics:
├── Breadth: # Active dojos / federations
├── Depth: Average revenue per customer (ARPC)
├── Frequency: Monthly engagement rate
├── Efficiency: Customer Acquisition Cost (CAC)
└── Retention: Monthly churn rate

Financial KPIs:
├── Revenue (Doanh thu)
├── Gross Profit (Lợi nhuận gộp)
├── Net Profit (LNST)
├── Cash Balance (Tiền mặt + TGNH)
├── Accounts Receivable (Phải thu)
├── Runway (months of cash remaining)
└── Burn Rate (monthly net cash outflow)

Operational KPIs:
├── # Active Customers
├── # New Customers (this month)
├── Customer Churn Rate
├── Support Ticket Resolution Time
├── Employee Count / Turnover
└── Sprint Velocity (engineering)
```

### Dashboard Layout
```
┌─────────────────────────────────────────────────────┐
│  KPI Cards (Revenue, Profit, Cash, Customers)       │
├─────────────────────────────┬───────────────────────┤
│  Revenue Trend (Line Chart) │  P&L Summary          │
│  12-month rolling           │  This month vs last    │
├─────────────────────────────┼───────────────────────┤
│  Revenue by Segment         │  Top 10 Customers     │
│  (Stacked Bar)              │  by revenue            │
├─────────────────────────────┼───────────────────────┤
│  Cash Flow Forecast         │  Recent Journal       │
│  (Area Chart)               │  Entries              │
├─────────────────────────────┴───────────────────────┤
│  Alerts & Action Items                               │
└─────────────────────────────────────────────────────┘
```

### Real-time Updates
```
├── Balance changes: After journal entry posting
├── KPI refresh: Every 5 minutes (SSE)
├── Alerts: Immediate (WebSocket)
└── Reports: On-demand with cache (5min TTL)
```

---

## 🔹 NĂNG LỰC: BUSINESS ANALYTICS

### Financial Analysis
```
Profitability Analysis:
├── Gross margin by segment (v_gross_profit_by_segment view)
├── Operating margin trend
├── Customer lifetime value (LTV)
├── Unit economics: LTV/CAC ratio (target > 3x)
└── Cohort analysis: Revenue retention

Cash Flow Analysis:
├── Operating cash flow
├── Investing cash flow
├── Financing cash flow
├── Free cash flow = Operating - CapEx
└── Cash conversion cycle

Variance Analysis:
├── Budget vs Actual (monthly)
├── Revenue forecast vs actual
├── Expense trend analysis
└── Department-level P&L
```

### Cross-Module Reporting
```
Finance + HR:
├── Cost per employee
├── Revenue per employee
├── Department P&L
└── Payroll as % of revenue

Finance + Sales:
├── Revenue by customer segment
├── Sales pipeline value vs target
├── Invoice aging report
├── Days Sales Outstanding (DSO)

Finance + Operations:
├── Gross margin by product/service
├── Break-even analysis
└── Working capital ratio
```

### Report Builder
```
Standard Reports:
├── Bảng Cân đối Kế toán (Balance Sheet) — B01-DN
├── Báo cáo KQKD (Income Statement) — B02-DN
├── Báo cáo Lưu chuyển Tiền tệ — B03-DN
├── Bảng Cân đối Phát sinh (Trial Balance)
├── Sổ Cái (General Ledger)
├── Sổ Nhật ký (General Journal)
├── Báo cáo Công nợ Phải thu
├── Báo cáo Công nợ Phải trả
└── Báo cáo Thuế GTGT

Custom Reports:
├── Drag-and-drop report builder
├── Filter by: date range, company, department, segment
├── Group by: account, month, quarter, year
├── Export: PDF, Excel, CSV
└── Schedule: Daily, weekly, monthly auto-send
```

---

## 🔹 NĂNG LỰC: OKR & STRATEGIC PLANNING

### OKR Tracking
```
Data Model:
├── objectives: title, owner, period, status
├── key_results: objective_id, metric, target, current, score
├── initiatives: key_result_id, title, owner, status, deadline
└── check_ins: key_result_id, date, value, notes

Score: 0.0 - 1.0
├── 0.7 - 1.0: Achieved (set harder next time)
├── 0.4 - 0.6: Good progress (stretch goals)
└── 0.0 - 0.3: Failed (root cause analysis)

Cadence:
├── Quarterly: Set OKRs
├── Weekly: Check-in scores
├── Monthly: Progress review
└── Quarterly: Score + retrospective
```

### Strategic Dashboard
```
├── OKR Progress (company-level)
├── Market position metrics
├── Competitive landscape changes
├── Risk register (top 5 risks)
└── Strategic initiative tracker
```

---

## Trigger Patterns

- "dashboard", "KPI", "báo cáo", "report"
- "analytics", "phân tích", "thống kê"
- "OKR", "mục tiêu", "chiến lược"
- "doanh số", "doanh thu", "lợi nhuận"
- "CEO", "điều hành", "tổng quan"

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
