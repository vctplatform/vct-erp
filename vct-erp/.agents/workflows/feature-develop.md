---
description: Quy trình phát triển tính năng end-to-end — từ PRD đến Deploy
---

# /feature-develop — Feature Development Workflow

## BƯỚC 1: REQUIREMENTS (Skill: erp-command)
// turbo
1. Đọc `.agents/skills/erp-command/SKILL.md`
2. Phân tích yêu cầu → Viết mini PRD:
   - Problem statement
   - User stories + acceptance criteria
   - Scope (in/out)
   - Success metrics

## BƯỚC 2: DESIGN (Skill: erp-command + domain module)
// turbo
1. Đọc SKILL.md của domain module liên quan
2. Database schema design (if needed)
3. API design (endpoints, request/response)
4. UI wireframe considerations

## BƯỚC 3: BACKEND (Skill: backend-engine)
// turbo
1. Đọc `.agents/skills/backend-engine/SKILL.md`
2. Create/update domain entities
3. Define repository interface
4. Implement use case
5. Implement PostgreSQL adapter
6. Create HTTP handler + routes
7. Write unit tests

## BƯỚC 4: FRONTEND (Skill: frontend-craft)
// turbo
1. Đọc `.agents/skills/frontend-craft/SKILL.md`
2. Create page component (App Router)
3. Create form/table components
4. Connect to API
5. Handle loading/error states
6. Test responsive layout

## BƯỚC 5: INTEGRATION TEST
// turbo
1. Test full flow: UI → API → DB → Response
2. Verify edge cases
3. Check error handling

## BƯỚC 6: REVIEW & REPORT
1. Summarize changes
2. List files modified/created
3. Note any follow-up items

// turbo-all

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
