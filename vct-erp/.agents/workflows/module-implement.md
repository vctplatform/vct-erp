---
description: Quy trình xây dựng module ERP mới — Blueprint
---

# /module-implement — New ERP Module Implementation

## BƯỚC 1: DOMAIN ANALYSIS
// turbo
1. Đọc domain SKILL.md tương ứng (module-finance, module-people, etc.)
2. Identify core entities and relationships
3. Define business rules and constraints
4. Map integration points with other modules

## BƯỚC 2: DATABASE SCHEMA
// turbo
1. Đọc `.agents/skills/backend-engine/SKILL.md` phần Database
2. Design tables with proper types, constraints, indexes
3. Create migration file
4. Add seed data if needed

## BƯỚC 3: BACKEND DOMAIN LAYER
// turbo
1. Create `internal/modules/{module}/domain/`
   - entity.go (structs)
   - repository.go (interfaces)
   - errors.go (domain errors)
2. Create `internal/modules/{module}/usecase/`
   - CRUD use cases
   - Business logic use cases
3. Create `internal/modules/{module}/adapter/postgres/`
   - Repository implementation

## BƯỚC 4: BACKEND API LAYER
// turbo
1. Create HTTP handlers
2. Define request/response DTOs
3. Register routes in httpapi
4. Add middleware (auth, validation)

## BƯỚC 5: FRONTEND PAGES
// turbo
1. Create `app/(dashboard)/{module}/` directory
2. Create page.tsx (list view)
3. Create [id]/page.tsx (detail view)
4. Create components (forms, tables)
5. Connect to API

## BƯỚC 6: TESTING
// turbo
1. Unit tests for use cases
2. Integration tests for repository
3. API tests for handlers
4. Frontend component tests

## BƯỚC 7: DOCUMENTATION
// turbo
1. Update ERP_CONTEXT.md module status
2. Create/update API docs
3. Add ADR if architectural decision made

// turbo-all
