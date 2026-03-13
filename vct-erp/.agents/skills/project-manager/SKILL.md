---
name: project-manager
description: Project Manager role - Sprint planning, task tracking, release management, risk assessment, stakeholder communication for VCT Platform.
---

# Project Manager - VCT Platform

## Role Overview
Manages project timeline, team coordination, sprint cycles, and ensures delivery milestones are met for the VCT Platform. Coordinates cross-platform delivery (Web + Mobile).

## Core Responsibilities

### 1. Sprint Management

#### Sprint Cycle (2-week sprints)
```
Day 1:  Sprint Planning → Break down epics → Estimate story points
Day 2-9:  Development → Daily standups → Track progress
Day 10: Code freeze → QA testing → Bug fixes
Day 11: Sprint Review → Demo to stakeholders
Day 12: Retrospective → Process improvements
```

#### Story Point Estimation Guide
| Points | Complexity | Example |
|--------|-----------|---------|
| 1 | Trivial | Fix typo, update text |
| 2 | Simple | Add field to form, simple validation |
| 3 | Small | CRUD endpoint, basic component |
| 5 | Medium | New page with Supabase API integration |
| 8 | Large | New module with business logic + RLS policies |
| 13 | Very Large | Complex integration (scoring engine + Supabase Realtime) |
| 21 | Epic | Architecture migration, Expo mobile app MVP |

### 2. Task Tracking Template

```markdown
## Sprint [N] - [Start Date] to [End Date]

### Goals
1. [Primary goal]
2. [Secondary goal]

### Backlog
| ID | Task | Owner | Points | Status | Platform | Notes |
|----|------|-------|--------|--------|----------|-------|
| VCT-001 | ... | @dev | 5 | TODO | Web+API | ... |
| VCT-002 | ... | @dev | 3 | TODO | Mobile | ... |

### Velocity
- Planned: [X] points
- Completed: [Y] points
- Carried over: [Z] points
```

### 3. Release Management

#### Version Naming Convention
```
v{MAJOR}.{MINOR}.{PATCH}[-{pre-release}]
Examples: v1.0.0, v1.2.0-beta.1, v2.0.0-rc.1
Mobile: v1.0.0 (build 42) → App Store / Play Store
```

#### Release Checklist
- [ ] All sprint tasks completed or deferred
- [ ] All tests passing (unit + integration + E2E)
- [ ] Code review approved for all PRs
- [ ] CHANGELOG.md updated
- [ ] Version bumped in `go.mod`, `package.json`, `app.json` (Expo)
- [ ] Database migrations reviewed (Neon branch tested)
- [ ] Supabase migrations synced (`supabase db push`)
- [ ] Staging deployment verified
- [ ] Mobile build tested (Expo EAS Build)
- [ ] Performance benchmarks acceptable
- [ ] Security scan passed
- [ ] Stakeholder sign-off received
- [ ] Rollback plan documented (Neon PITR)
- [ ] Production deployment executed
- [ ] Post-deployment health check passed
- [ ] Release notes published
- [ ] Mobile app submitted to App Store / Play Store (if applicable)

### 4. Risk Register

| ID | Risk | Category | Impact | Likelihood | Mitigation | Owner | Status |
|----|------|----------|--------|-----------|------------|-------|--------|
| R-001 | Tournament peak traffic | Technical | High | High | Neon autoscaling, load testing | DevOps | Active |
| R-002 | Data loss | Technical | Critical | Low | Neon PITR, Supabase backups | DBA | Active |
| R-003 | Scope creep | Process | Medium | High | Change request process | PM | Active |
| R-004 | Key person dependency | People | High | Medium | Cross-training, docs | PM | Active |
| R-005 | App Store rejection | Technical | Medium | Low | Follow guidelines early | Frontend | Active |
| R-006 | Vendor lock-in (Supabase/Neon) | Technical | Medium | Low | Abstract infra layer | Architect | Active |

### 5. Communication Plan

| Channel | Purpose | Frequency | Tool |
|---------|---------|-----------|------|
| Slack #vct-dev | Daily dev discussion | Real-time | Slack |
| Slack #vct-alerts | Production alerts | Real-time | Slack/PagerDuty |
| GitHub Issues | Task tracking | Ongoing | GitHub |
| GitHub PRs | Code review | Per feature | GitHub |
| Email | Stakeholder updates | Weekly | Email |
| Google Meet | Sprint ceremonies | Bi-weekly | Google Meet |

### 6. Definition of Done (DoD)

A task is "Done" when:
- [ ] Code implemented and self-reviewed
- [ ] Unit tests written (≥ 80% coverage for new code)
- [ ] Integration tests passing (Neon branch)
- [ ] Code review approved by ≥ 1 reviewer
- [ ] No lint errors (`golangci-lint`, `eslint`)
- [ ] TypeScript compiles without errors
- [ ] API documentation updated (if API changed)
- [ ] i18n translations added (Vietnamese + English)
- [ ] Responsive design verified (mobile + desktop)
- [ ] Expo mobile tested (if cross-platform feature)
- [ ] Supabase RLS policies verified
- [ ] Accessibility checked (WCAG 2.1 AA)
- [ ] Deployed to staging and manually verified

### 7. Escalation Path

```
Level 1: Developer → Team Lead (< 4 hours)
Level 2: Team Lead → PM (< 8 hours)
Level 3: PM → CTO (< 24 hours)
Level 4: CTO → Stakeholders (critical issues)
```
