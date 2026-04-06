# Shared Knowledge — VCT ERP

> Thư mục này chứa kiến thức dùng chung cho toàn bộ đội agent trong dự án ERP.

## Cấu trúc

```
shared_knowledge/
├── README.md              ← File này
├── tech_stack.md           ← Technology decisions & standards
├── architecture/           ← Architecture Decision Records (ADRs)
└── domain/                 ← ERP domain knowledge
```

## Quy tắc

1. **Kiến thức chung** — Mọi agent đều đọc được
2. **Chỉ facts, không opinion** — Lưu quyết định đã được approve, không phải đề xuất
3. **Luôn có context** — Mỗi document phải nói rõ WHY, không chỉ WHAT
4. **Cập nhật khi thay đổi** — Khi quyết định thay đổi, update document, không xóa lịch sử

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
