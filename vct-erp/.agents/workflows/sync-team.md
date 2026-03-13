---
description: Sync codebase across team members - pull, rebase, resolve conflicts
---

# /sync-team - Team Code Synchronization

// turbo-all

## Daily Sync Routine (Start of Day)

### Step 1: Stash any work in progress
```bash
git stash save "WIP: $(date +%Y%m%d)"
```

### Step 2: Fetch latest from remote
```bash
git fetch origin --prune
```

### Step 3: Update develop branch
```bash
git checkout develop
git pull origin develop --rebase
```

### Step 4: Rebase current feature branch
```bash
git checkout feature/VCT-xxx-your-branch
git rebase origin/develop
```

### Step 5: Restore stashed work
```bash
git stash pop
```

### Step 6: Resolve conflicts (if any)
```bash
# If rebase conflicts occur:
# 1. Check conflicting files
git status

# 2. Open files with conflicts and resolve
# Look for <<<<<<< HEAD markers

# 3. After resolving:
git add .
git rebase --continue

# 4. If you need to abort:
git rebase --abort
```

### Step 7: Force push rebased branch (if needed)
```bash
git push origin HEAD --force-with-lease
```

## Sync Backend Dependencies
```bash
cd backend
go mod tidy
go mod download
```

## Sync Frontend Dependencies
```bash
cd frontend
npm ci  # Clean install from lock file
```

## Sync Database
```bash
# Apply any new migrations
cd backend
migrate -path migrations -database "${DATABASE_URL}" up
```

## Full Sync Checklist
- [ ] Git fetch and pull latest
- [ ] Rebase feature branch on develop
- [ ] Resolve any conflicts
- [ ] Update Go dependencies (`go mod tidy`)
- [ ] Update Node dependencies (`npm ci`)
- [ ] Apply new database migrations
- [ ] Rebuild Docker images if Dockerfile changed
- [ ] Run quick tests to verify
