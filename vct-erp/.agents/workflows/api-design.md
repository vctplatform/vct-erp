---
description: Quy trình thiết kế API — Spec, Validate, Document
---

# /api-design — API Design Workflow

## BƯỚC 1: REQUIREMENTS
// turbo
1. Identify resource (noun): /journal-entries, /employees, /invoices
2. Define operations: CRUD + custom actions
3. List request/response fields
4. Define error codes

## BƯỚC 2: DESIGN
// turbo
1. Đọc `.agents/skills/backend-engine/SKILL.md` phần API Design
2. Follow REST conventions:
   - GET /v1/{resources} — List (with pagination)
   - GET /v1/{resources}/{id} — Get one
   - POST /v1/{resources} — Create
   - PUT /v1/{resources}/{id} — Full update
   - PATCH /v1/{resources}/{id} — Partial update
   - DELETE /v1/{resources}/{id} — Delete
3. Define query params (filter, sort, pagination)
4. Define headers (Authorization, Idempotency-Key)

## BƯỚC 3: SPEC
// turbo
1. Write OpenAPI 3.1 spec or inline doc
2. Define JSON schema for request/response
3. Define error response format
4. Add examples

## BƯỚC 4: IMPLEMENT
// turbo
1. Create handler, request, response files
2. Implement with proper validation
3. Write tests
4. Update API documentation

// turbo-all
