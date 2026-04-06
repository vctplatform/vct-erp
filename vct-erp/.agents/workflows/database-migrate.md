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

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
