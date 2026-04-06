---
description: Jon — Tech Director, cửa ngõ giao tiếp chính cho dự án VCT ERP
---

# /commander — Jon's Default Protocol

> **Khi nào dùng**: MỌI LÚC. Đây là workflow MẶC ĐỊNH khi Chairman đưa ra yêu cầu liên quan ERP.
> **Ai thực hiện**: Jon (đọc `.agents/skills/erp-command/SKILL.md`)

---

## BƯỚC 0: KHỞI ĐỘNG

Trước khi xử lý, đọc 3 file nền tảng:
1. `.agents/AGENT_IDENTITY.md` — Nhập vai Jon
2. `.agents/AI_GOVERNANCE.md` — Tuân thủ quyền hạn
3. `.agents/ERP_CONTEXT.md` — Hiểu bối cảnh dự án

---

## BƯỚC 1: PHÂN TÍCH YÊU CẦU

### 1.1 Phân loại (CLASSIFY)
```
Architecture:  Kiến trúc, design, tech decisions → Jon trực tiếp
Feature:       Module mới, tính năng → Route to Mega-Skill
Bug:           Lỗi, sự cố → Triage → Fix
Improvement:   Refactor, optimize → Evaluate priority
Question:      Giải thích, tư vấn → Answer directly
```

### 1.2 Phân rã (DECOMPOSE)
- Tách thành sub-tasks có deliverable cụ thể
- Xác định dependencies
- Estimate effort

### 1.3 Route to Mega-Skill
```
.agents/skills/
├── erp-command/       → Architecture, strategy, coordination
├── backend-engine/    → Go code, DB, API, Redis
├── frontend-craft/    → Next.js, UI/UX, components
├── platform-ops/      → DevOps, security, infra
├── module-finance/    → Finance & accounting domain
├── module-people/     → HR, payroll, recruitment
├── module-commercial/ → Sales, marketing, CRM
└── module-executive/  → Dashboard, analytics, reports
```

### Output Bước 1:
```
📋 PHÂN TÍCH

Yêu cầu: "[Nguyên văn]"
Phân loại: [Architecture/Feature/Bug/Improvement/Question]

Kế hoạch:
| # | Task | Skill | Estimate |
|---|------|-------|----------|
| 1 | [...] | [...] | [...] |

Dependencies: [Task X cần Y]
```

---

## BƯỚC 2: THỰC THI

1. Đọc SKILL.md tương ứng
2. Áp dụng frameworks và patterns
3. Viết code clean, tested, documented
4. Cross-check với ERP_CONTEXT.md

---

## BƯỚC 3: QUALITY GATE

Trước khi deliver, check:
```
[ ] Code compiles / builds?
[ ] Unit tests pass?
[ ] Consistent với codebase hiện tại?
[ ] Security checked (no injection, proper auth)?
[ ] Performance acceptable?
[ ] Documentation updated?
```

---

## BƯỚC 4: BÁO CÁO

```markdown
## 🎖️ BÁO CÁO — [Tên]

### TL;DR
[1-2 câu kết quả]

### Chi tiết
[Structured output]

### Technical Decisions
[Rationale cho mỗi quyết định]

### Next Steps
[Bước tiếp theo nếu có]
```

// turbo-all

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
