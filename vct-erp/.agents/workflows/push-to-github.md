---
description: Push code to GitHub with proper repository setup, branch protection, and CI/CD
---

# /push-to-github - Push Code to GitHub

// turbo-all

## Step 1: Create GitHub repository (first time only)
```bash
# Via GitHub CLI
gh repo create vct-platform --private --description "VCT Platform - Vietnam Cycling & Triathlon Management System"

# Or manually:
# 1. Go to https://github.com/new
# 2. Name: vct-platform
# 3. Visibility: Private
# 4. Do NOT initialize with README (we already have code)
```

## Step 2: Add remote origin (first time only)
```bash
git remote add origin https://github.com/{org}/vct-platform.git
# Verify:
git remote -v
```

## Step 3: Configure branch protection (first time only)
```bash
# Via GitHub CLI - protect main branch
gh api repos/{owner}/{repo}/branches/main/protection -X PUT -f '{
  "required_status_checks": {
    "strict": true,
    "contexts": ["backend-test", "frontend-test", "docker-build"]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "required_approving_review_count": 2
  },
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false
}'
```

## Step 4: Ensure code is clean before pushing
```bash
# Backend checks
cd backend
go vet ./...
golangci-lint run ./...
go test ./... -count=1

# Frontend checks
cd frontend
npm run lint
npm run type-check
npm run build
```

## Step 5: Stage and commit
```bash
git add .
git status  # Review what's being committed
git commit -m "feat(scope): descriptive commit message"
```

## Step 6: Push to remote
```bash
# Push current branch
git push origin HEAD

# Push and set upstream (first push)
git push -u origin HEAD
```

## Step 7: Create Pull Request (if pushing feature branch)
```bash
gh pr create \
  --base develop \
  --title "feat(scope): description" \
  --body "## Changes\n- List of changes\n\n## Testing\n- How it was tested" \
  --assignee @me \
  --reviewer teammate1,teammate2
```

## Step 8: Verify CI pipeline
```bash
# Check CI status
gh pr checks

# View workflow runs
gh run list --limit 5

# Watch a specific run
gh run watch
```

## Step 9: Tag releases (for production pushes)
```bash
# Create annotated tag
git tag -a v1.0.0 -m "Release v1.0.0 - Initial release with athlete and club modules"

# Push tag
git push origin v1.0.0

# Create GitHub release
gh release create v1.0.0 --title "v1.0.0" --notes "## What's New\n- Athlete management\n- Club management\n- Authentication system"
```

## Repository Structure

```
.github/
├── workflows/
│   ├── ci.yml          # CI: lint, test, build
│   ├── staging.yml     # Deploy to staging
│   ├── production.yml  # Deploy to production
│   └── release.yml     # Release automation
├── PULL_REQUEST_TEMPLATE.md
├── ISSUE_TEMPLATE/
│   ├── bug_report.md
│   └── feature_request.md
└── CODEOWNERS
```

## CODEOWNERS file
```
# Backend
/backend/ @backend-lead @cto

# Frontend
/frontend/ @frontend-lead @cto

# Infrastructure
/deploy/ @devops-lead
/.github/workflows/ @devops-lead

# Database
/backend/migrations/ @dba @backend-lead
```
