---
description: Quy trình thay đổi database schema — Migration Workflow
---

# /database-migrate — Database Migration Workflow

## BƯỚC 1: PLAN
// turbo
1. Đọc `.agents/skills/backend-engine/SKILL.md` phần Database
2. Document current schema state
3. Design new schema changes
4. Write migration SQL (UP + DOWN)

## BƯỚC 2: CREATE MIGRATION
// turbo
1. Create migration file: `migrations/XXXXXX_description.sql`
2. Include both UP and DOWN statements
3. Use transactions for DDL when possible
4. Add NOT NULL with DEFAULT to avoid table locks
5. Create indexes CONCURRENTLY for large tables

## BƯỚC 3: TEST
// turbo
1. Apply migration on dev database
2. Verify schema matches expectations
3. Test rollback (DOWN migration)
4. Test with seed data

## BƯỚC 4: UPDATE CODE
// turbo
1. Update Go entities/domain if needed
2. Update repository queries
3. Update API if schema affects response
4. Run full test suite

## CHƯA LÀM: Deploy migration to staging/production (Cấp 3 — cần Chairman approve)

// turbo-all
