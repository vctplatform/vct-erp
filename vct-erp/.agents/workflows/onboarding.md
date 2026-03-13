---
description: New team member onboarding workflow for VCT Platform
---

# /onboarding - New Team Member Onboarding

## Day 1: Environment Setup

### Step 1: Access & Accounts
- [ ] GitHub organization invite accepted
- [ ] Slack workspace joined (#vct-dev, #vct-alerts)
- [ ] Google Meet access verified
- [ ] Task tracking tool access (GitHub Issues/Projects)

### Step 2: Clone repository
```bash
git clone https://github.com/{org}/vct-platform.git
cd vct-platform
```

### Step 3: Run full environment setup
Follow `/setup-environment` workflow to install all tools.

### Step 4: Configure Git identity
```bash
git config user.name "Full Name"
git config user.email "email@organization.com"
```

### Step 5: Verify development environment
```bash
# Start services
docker compose up -d

# Backend
cd backend && go test ./...

# Frontend
cd frontend && npm run dev
```

## Day 1-2: Codebase Orientation

### Step 6: Read key documentation
- [ ] `README.md` - Project overview
- [ ] `docs/architecture/` - System architecture
- [ ] `docs/api/` - API documentation
- [ ] `.agents/skills/system-architect/SKILL.md` - Architecture patterns
- [ ] `.agents/skills/backend-developer/SKILL.md` - Go patterns (if backend)
- [ ] `.agents/skills/frontend-developer/SKILL.md` - React patterns (if frontend)

### Step 7: Understand project structure
```bash
# Explore backend structure
find backend/internal -type f -name "*.go" | head -30

# Explore frontend structure
find frontend/src -type f -name "*.tsx" | head -30
```

### Step 8: Read existing code
Start with a simple module (e.g., `athlete`) and trace the flow:
1. `adapter/http/router.go` → Route definitions
2. `adapter/http/handler.go` → HTTP handlers
3. `usecase/create.go` → Business logic
4. `domain/entity.go` → Domain model
5. `adapter/postgres/repository.go` → Database layer

## Day 2-3: First Contribution

### Step 9: Pick a starter task
- Look for issues labeled `good-first-issue`
- Or fix a simple bug / add a test

### Step 10: Follow the git workflow
Follow `/git-workflow` for branching and commit conventions.

### Step 11: Submit first PR
Follow `/code-review` checklist before requesting review.

## Week 1: Integration

### Step 12: Team introductions
- [ ] Meet team leads (backend, frontend, DevOps)
- [ ] Attend daily standup
- [ ] Shadow PR reviews

### Step 13: Domain knowledge
- [ ] Read Business Analyst skill for VCT domain understanding
- [ ] Understand Vietnamese cycling/triathlon federation structure
- [ ] Review key Vietnamese terms used in the system

## Onboarding Checklist Summary
- [ ] Development environment running
- [ ] Can build and test locally
- [ ] Understands Clean Architecture
- [ ] Read relevant skill documents
- [ ] First PR submitted and merged
- [ ] Attended sprint ceremonies
- [ ] Understands git workflow and commit conventions
