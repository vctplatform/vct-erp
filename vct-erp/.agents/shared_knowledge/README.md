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
