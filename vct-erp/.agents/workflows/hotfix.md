---
description: Emergency hotfix workflow for critical production bugs
---

# /hotfix - Emergency Hotfix Workflow

## ⚠️ Only use for CRITICAL production issues that cannot wait for next sprint.

## Criteria for Hotfix
- Production system down or severely degraded
- Security vulnerability discovered
- Data corruption or loss
- Payment/financial processing failure

## Steps

### Step 1: Create hotfix branch from main
```bash
git checkout main
git pull origin main
git checkout -b hotfix/VCT-XXX-description
```

### Step 2: Implement the fix
```bash
# Make minimal, targeted changes
# Focus ONLY on the bug fix, no feature work
```

### Step 3: Test thoroughly
```bash
cd backend && go test -race -count=1 ./...
cd frontend && npm run test -- --run
```

### Step 4: Get expedited review
```bash
# Create PR with [HOTFIX] label
gh pr create \
  --base main \
  --title "[HOTFIX] VCT-XXX: Fix critical issue" \
  --body "## Critical Bug Fix

**Impact**: [Describe user/system impact]
**Root Cause**: [Brief root cause]
**Fix**: [What this PR does]
**Tested**: [How it was verified]

cc: @cto @backend-lead" \
  --label "hotfix,critical" \
  --reviewer cto,backend-lead
```

### Step 5: Merge and deploy
```bash
# After CTO approval
git checkout main
git merge hotfix/VCT-XXX-description --no-ff
git tag -a v$(cat VERSION)-hotfix.1 -m "Hotfix: description"
git push origin main --tags

# Deploy immediately (see /deploy-production)
```

### Step 6: Back-merge to develop
```bash
git checkout develop
git merge main --no-ff
git push origin develop
```

### Step 7: Cleanup
```bash
git branch -d hotfix/VCT-XXX-description
```

### Step 8: Post-mortem
```markdown
## Incident Report - VCT-XXX

**Date**: YYYY-MM-DD HH:MM
**Duration**: X hours
**Impact**: [Users affected, services impacted]
**Root Cause**: [Detailed root cause]
**Fix Applied**: [Description of fix]
**Prevention Measures**:
1. [What we'll do to prevent recurrence]
2. [Additional monitoring/testing]
**Lessons Learned**:
- [Key takeaway]
```
