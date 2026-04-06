---
name: module-commercial
description: >-
  Mega-Skill for Sales, Marketing & CRM ERP module. Sales pipeline, quotation,
  invoicing, campaign management, and customer service for VCT business.
metadata:
  author: VCT Platform
  version: "1.0.0"
  type: "Mega-Skill"
  locale: vi-VN
---

# MODULE-COMMERCIAL — MEGA-SKILL

> Domain expertise cho module Kinh doanh, Marketing & CRM.

---

## 🔹 NĂNG LỰC: SALES MANAGEMENT

### Sales Pipeline (VCT Business)
```
Stages:
├── 1. Lead (Liên hệ ban đầu)
├── 2. Qualified (Xác nhận nhu cầu)
├── 3. Proposal (Gửi báo giá / demo)
├── 4. Negotiation (Đàm phán điều khoản)
├── 5. Closed Won (Ký hợp đồng)
└── 6. Closed Lost (Mất deal → phân tích lý do)

VCT Revenue Streams:
├── SaaS Subscriptions: Platform access for Federations/Dojos
├── Dojo Management Fees: Monthly per-student fees
├── Equipment Rental: Scoring devices, mats
├── Training Events: Tournament management
└── Consulting: Digital transformation advisory
```

### Data Model
```
├── leads: source, contact_info, status, score
├── opportunities: lead_id, value, stage, probability, close_date
├── quotations: opportunity_id, items[], total, valid_until
├── sales_orders: quotation_id, payment_terms, delivery
├── invoices: order_id, items[], tax, total, due_date, status
└── payments: invoice_id, amount, method, date
```

### Quotation → Invoice → Journal Entry Flow
```
1. Sales creates Quotation
2. Customer approves → Sales Order
3. Service delivered → Invoice generated
4. Invoice posted → Journal Entry auto-created:
   Nợ 131 (Phải thu KH)         [Total incl. VAT]
       Có 5113 (DT cung cấp DV)          [Amount excl. VAT]
       Có 33311 (Thuế GTGT đầu ra)       [VAT amount]
5. Payment received → Settlement entry:
   Nợ 1121 (TGNH)               [Amount received]
       Có 131 (Phải thu KH)               [Amount received]
```

### Invoice (Hóa đơn GTGT)
```
Required fields per Vietnam e-invoice law:
├── Tên, địa chỉ, MST người bán
├── Tên, địa chỉ, MST người mua (nếu > 200K)
├── Tên hàng hóa, dịch vụ
├── Đơn vị tính, số lượng, đơn giá
├── Thuế suất GTGT (0%, 5%, 8%, 10%)
├── Tổng tiền chưa thuế, tiền thuế, tổng thanh toán
├── Chữ ký số (e-invoice)
└── Mã CQT (authority code)
```

---

## 🔹 NĂNG LỰC: MARKETING & CAMPAIGNS

### Campaign Management
```
Data Model:
├── campaigns: name, type, budget, start/end, goals
├── campaign_channels: campaign_id, channel, spend, metrics
├── campaign_leads: campaign_id, lead_id (attribution)
└── campaign_content: assets, copy, creatives

Types:
├── Digital: Social media, SEO, SEM, Email
├── Event: Tournament sponsorship, Demo day
├── Content: Blog, Video, Webinar
└── Partnership: Federation co-marketing

Metrics per campaign:
├── Reach / Impressions
├── Engagement (likes, shares, comments)
├── Leads generated
├── Cost per Lead (CPL)
├── Conversion rate (Lead → Customer)
└── ROI = (Revenue - Cost) / Cost
```

### Customer Segmentation (VCT)
```
├── Federations (Liên đoàn): Enterprise, multi-year contracts
├── Dojos (Võ đường): SMB, monthly subscriptions
├── Athletes (VĐV): Freemium → Premium individual
├── Event Organizers: Pay-per-event
└── Equipment Buyers: One-time + recurring
```

---

## 🔹 NĂNG LỰC: CUSTOMER SERVICE

### Ticket System
```
Data Model:
├── tickets: customer_id, subject, description, priority, status
├── ticket_messages: ticket_id, sender, content, timestamp
├── ticket_assignments: ticket_id, agent_id, assigned_at
└── sla_configs: priority → response_time, resolution_time

Priority & SLA:
├── Critical: 1h response, 4h resolution (system down)
├── High: 4h response, 24h resolution (feature broken)
├── Medium: 8h response, 48h resolution (question/request)
└── Low: 24h response, 1 week resolution (enhancement)

Status: new → assigned → in_progress → resolved → closed
```

### Customer Health Score
```
Score (0-100) based on:
├── Usage frequency (30%)
├── Feature adoption (20%)
├── Support tickets (15%, inverse)
├── Payment history (15%)
├── Engagement (10%)
└── NPS/Satisfaction (10%)

Actions:
├── 80-100: Expand/upsell candidates
├── 60-79: Healthy, maintain
├── 40-59: At risk, proactive outreach
└── 0-39: Critical, escalate to leadership
```

---

## Trigger Patterns

- "bán hàng", "sales", "kinh doanh", "doanh thu"
- "khách hàng", "customer", "CRM", "lead"
- "hóa đơn", "invoice", "báo giá", "quotation"
- "marketing", "chiến dịch", "campaign"
- "hỗ trợ", "support", "ticket"

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
