---
description: Quy trình review code — PR Review Process
---

# /code-review — Code Review Workflow

## CHECKLIST
// turbo
Đọc `.agents/skills/erp-command/SKILL.md` phần Code Review Standards.

Per file/change, check:
1. **Correctness**: Logic đúng? Edge cases handled?
2. **Readability**: Tên biến/function rõ nghĩa?
3. **Simplicity**: Có cách đơn giản hơn?
4. **Testing**: Unit test đi kèm?
5. **Security**: SQL injection? Auth check?
6. **Performance**: N+1? Unnecessary computation?
7. **Consistency**: Follow project patterns?
8. **ERP-specific**: Double-entry cân bằng? VAS compliance?

## OUTPUT FORMAT
```markdown
## 📝 CODE REVIEW

### Summary
[Tổng quan thay đổi]

### ✅ Approved / ❌ Changes Requested / 💬 Comments

### Findings
| # | File | Line | Type | Comment |
|---|------|------|------|---------|

### Suggestions
[Improvement opportunities]
```

// turbo-all

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
