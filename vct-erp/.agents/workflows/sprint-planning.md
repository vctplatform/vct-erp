---
description: Sprint planning, task breakdown, estimation, and assignment for VCT Platform
---

# /sprint-planning - Sprint Planning Workflow

## When to Trigger
- Start of every 2-week sprint
- After stakeholder meetings with new requirements

## Step 1: Review Previous Sprint
- [ ] Review velocity (planned vs completed story points)
- [ ] Review carried-over tasks
- [ ] Document lessons learned

## Step 2: Prioritize Backlog

### Priority Levels
| Priority | Label | Description | SLA |
|----------|-------|-------------|-----|
| P0 | 🔴 Critical | Blocks production/users | Must fix this sprint |
| P1 | 🟠 High | Core feature/major bug | Should fix this sprint |
| P2 | 🟡 Medium | Enhancement/minor bug | Plan for this sprint |
| P3 | 🟢 Low | Nice to have | Backlog |

## Step 3: Break Down Epics into Tasks

### Task Template
```markdown
### [VCT-XXX] [Task Title]
- **Epic**: [Parent epic]
- **Type**: Feature / Bug / Chore / Improvement
- **Priority**: P0 / P1 / P2 / P3
- **Estimate**: [story points]
- **Assignee**: @developer
- **Acceptance Criteria**:
  - [ ] ...
- **Dependencies**: [VCT-YYY if any]
```

### Estimation Guide
| Points | Time | Example |
|--------|------|---------|
| 1 | < 2 hours | Config change, text update |
| 2 | 2-4 hours | Simple CRUD endpoint |
| 3 | 4-8 hours | New component with API hook |
| 5 | 1-2 days | Full feature page |
| 8 | 2-3 days | Complex module |
| 13 | 3-5 days | Cross-module integration |

## Step 4: Assign Tasks

### Team Capacity
```
Sprint Capacity = (Team Size × Working Days × 6 hours) / Average Velocity Factor

Example:
- 4 developers × 10 days × 6 productive hours = 240 hours
- Velocity factor: 0.7 (meetings, reviews, etc.)
- Effective capacity: 168 hours ≈ 50-60 story points
```

### Assignment Rules
- No developer exceeds 80% capacity
- Each task has ONE primary owner
- Junior devs pair with seniors on complex tasks
- Balance frontend/backend work per developer

## Step 5: Sprint Board Setup

```
┌──────────┬──────────┬──────────┬──────────┬──────────┐
│ BACKLOG  │  TODO    │ IN PROG  │ REVIEW   │  DONE    │
├──────────┼──────────┼──────────┼──────────┼──────────┤
│          │ VCT-101  │ VCT-098  │ VCT-095  │ VCT-090  │
│          │ VCT-102  │ VCT-099  │          │ VCT-091  │
│          │ VCT-103  │          │          │ VCT-092  │
└──────────┴──────────┴──────────┴──────────┴──────────┘
```

## Step 6: Document Sprint Goals
```markdown
## Sprint [N] Goals
1. [Primary goal - must achieve]
2. [Secondary goal - should achieve]
3. [Stretch goal - nice to have]

### Risks
- [Risk 1]: [Mitigation]
- [Risk 2]: [Mitigation]
```
