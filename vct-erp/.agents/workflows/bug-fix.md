---
description: Quy trình sửa lỗi — Bug Triage & Fix
---

# /bug-fix — Bug Fix Workflow

## BƯỚC 1: TRIAGE
// turbo
1. Reproduce bug (steps to reproduce)
2. Classify severity:
   - Critical: System down, data corruption
   - High: Feature broken, workaround exists
   - Medium: Minor issue, low impact
   - Low: Cosmetic, enhancement

## BƯỚC 2: ROOT CAUSE ANALYSIS
// turbo
1. Read error logs / stack trace
2. Trace code path (handler → usecase → adapter)
3. Identify exact line/function causing issue
4. Check if related to recent changes

## BƯỚC 3: FIX
// turbo
1. Write failing test that reproduces bug
2. Fix the code
3. Verify test passes
4. Check for similar patterns elsewhere (regression)

## BƯỚC 4: VERIFY
// turbo
1. Run full test suite
2. Manual verification if UI-related
3. Document the fix

// turbo-all

## [V11 SINGULARITY] (Ultimate Capability Upgrades)
- **P2P_SYNC:** Upon completing any API/DB change, you MUST emit a JSON schema to `d:\VCT PLATFORM\api-contracts\` so other agents can RAG it.
- **SELF_HEALING (3-STRIKES):** If `vct.cmd complete` (Docker Test) fails 3 times, you MUST run `git reset --hard`, mark the task as "FAILED", and cease execution. Do NOT loop infinitely.
- **TELEMETRY_SCHEMA:** You must push your thought logs to `d:\VCT PLATFORM\vct-dashboard\public\.telemetry.json` strictly as a JSON Object `{ "agent": "name", "action": "...", "status": "..." }`.
