---
description: Git branching strategy, commit conventions, and PR workflow for VCT Platform
---

# /git-workflow - Git Branching & Commit Strategy

## Branch Strategy (Git Flow Simplified)

```
main ─────────────────────────────────────────────── (production)
  │
  └── develop ─────────────────────────────────────── (integration)
        │          │           │
        └── feature/VCT-123-athlete-crud             (feature)
        └── feature/VCT-456-tournament-scoring       (feature)
        └── bugfix/VCT-789-login-fix                 (bugfix)
```

### Branch Naming Convention
```
feature/VCT-{ticket}-{short-description}
bugfix/VCT-{ticket}-{short-description}
hotfix/VCT-{ticket}-{short-description}
release/v{major}.{minor}.{patch}
```

**Examples:**
```
feature/VCT-001-athlete-registration
feature/VCT-045-live-scoring-websocket
bugfix/VCT-123-pagination-offset
hotfix/VCT-200-auth-token-expiry
release/v1.2.0
```

## Commit Convention (Conventional Commits)

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types
| Type | Description | Example |
|------|------------|---------|
| `feat` | New feature | `feat(athlete): add registration API` |
| `fix` | Bug fix | `fix(auth): correct token expiry calculation` |
| `docs` | Documentation | `docs(api): update swagger annotations` |
| `style` | Formatting only | `style(frontend): fix indentation` |
| `refactor` | Code cleanup | `refactor(club): extract validation logic` |
| `test` | Add/modify tests | `test(athlete): add unit tests for create` |
| `chore` | Build/tooling | `chore(ci): update Go version in workflow` |
| `perf` | Performance | `perf(query): add index for athlete search` |
| `build` | Build system | `build(docker): optimize multi-stage build` |
| `ci` | CI/CD changes | `ci(github): add deployment workflow` |

### Scopes
```
athlete, club, tournament, federation, btc, parent,
scoring, finance, report, notification,
auth, api, db, frontend, backend, ci, docker, deps
```

## Daily Workflow Steps

// turbo-all

### Step 1: Start your day - sync with develop
```bash
git checkout develop
git pull origin develop
```

### Step 2: Create feature branch
```bash
git checkout -b feature/VCT-xxx-description
```

### Step 3: Make changes and commit frequently
```bash
git add .
git commit -m "feat(athlete): implement create athlete endpoint"
```

### Step 4: Keep branch up to date with develop
```bash
git fetch origin develop
git rebase origin/develop
# Resolve any conflicts if needed
```

### Step 5: Push to remote
```bash
git push origin feature/VCT-xxx-description
```

### Step 6: Create Pull Request
```bash
# Via GitHub CLI
gh pr create --base develop --title "feat(athlete): implement athlete CRUD API" --body "
## Description
Implements the athlete CRUD endpoints following Clean Architecture.

## Changes
- Added athlete domain entity and repository interface
- Implemented PostgreSQL repository
- Created HTTP handlers with validation
- Added unit and integration tests

## Checklist
- [x] Tests pass
- [x] Lint clean
- [x] API docs updated
- [x] i18n translations added
"
```

### Step 7: After PR is merged, cleanup
```bash
git checkout develop
git pull origin develop
git branch -d feature/VCT-xxx-description
```

## Merge Rules

| Target Branch | Source Branch | Strategy | Approval |
|--------------|-------------|----------|----------|
| `develop` | `feature/*` | Squash merge | 1 reviewer |
| `develop` | `bugfix/*` | Squash merge | 1 reviewer |
| `main` | `develop` | Merge commit | 2 reviewers + PM |
| `main` | `hotfix/*` | Merge commit | CTO approval |

## Protected Branch Rules (GitHub)
- `main`: Require PR, 2 approvals, CI passing, no force push
- `develop`: Require PR, 1 approval, CI passing
