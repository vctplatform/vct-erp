---
description: Quy trình xử lý sự cố — Incident Response & Post-mortem
---

# /incident-response — Incident Handling Workflow

## BƯỚC 1: DETECT & CLASSIFY
1. Classify severity:
   - SEV1: System down → ALL HANDS, < 15min
   - SEV2: Major feature broken → < 30min
   - SEV3: Minor issue → < 2h
   - SEV4: Cosmetic → Next day

## BƯỚC 2: CONTAIN
// turbo
1. Identify affected services/users
2. Rollback if recent deployment caused it
3. Apply temporary workaround if possible
4. Communicate status to Chairman

## BƯỚC 3: FIX
// turbo
1. Root cause analysis (5 Whys)
2. Implement permanent fix
3. Test fix thoroughly
4. Deploy fix (with approval if needed)

## BƯỚC 4: POST-MORTEM
// turbo
1. Write post-mortem:
   - Timeline (minute by minute)
   - Impact (users, duration, data)
   - Root cause (5 Whys)
   - What went well
   - What went poorly
   - Action items (owners + deadlines)
2. Share with Chairman
3. Update runbook if needed

// turbo-all
