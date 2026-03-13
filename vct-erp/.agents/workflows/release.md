---
description: Version release process with semantic versioning, changelog, and tagging
---

# /release - Version Release Process

## Version Strategy: Semantic Versioning

```
v{MAJOR}.{MINOR}.{PATCH}
  │        │        └── Bug fixes, patches (backward compatible)
  │        └── New features (backward compatible)
  └── Breaking changes
```

## Steps

### Step 1: Determine version bump
```bash
# Check commits since last tag
git log $(git describe --tags --abbrev=0)..HEAD --oneline

# Determine bump type based on commits:
# feat: → MINOR
# fix:  → PATCH
# feat!: or BREAKING CHANGE: → MAJOR
```

### Step 2: Create release branch
```bash
NEW_VERSION="1.2.0"  # Set your version
git checkout develop
git pull origin develop
git checkout -b release/v${NEW_VERSION}
```

### Step 3: Update version files
```bash
echo "${NEW_VERSION}" > VERSION

# Update backend version
sed -i "s/Version = \".*\"/Version = \"${NEW_VERSION}\"/" backend/internal/shared/config/version.go

# Update frontend version
cd frontend && npm version ${NEW_VERSION} --no-git-tag-version
```

### Step 4: Update CHANGELOG.md
```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- feat(module): description of new features

### Changed
- refactor(module): description of changes

### Fixed
- fix(module): description of bug fixes

### Security
- security(module): description of security fixes
```

### Step 5: Final tests
```bash
cd backend && go test -race ./...
cd frontend && npm run test -- --run && npm run build
```

### Step 6: Commit release changes
```bash
git add .
git commit -m "chore(release): v${NEW_VERSION}"
```

### Step 7: Merge to main
```bash
git checkout main
git merge release/v${NEW_VERSION} --no-ff -m "Release v${NEW_VERSION}"
git tag -a v${NEW_VERSION} -m "Release v${NEW_VERSION}"
git push origin main --tags
```

### Step 8: Merge back to develop
```bash
git checkout develop
git merge main --no-ff
git push origin develop
```

### Step 9: Create GitHub release
```bash
gh release create v${NEW_VERSION} \
  --title "v${NEW_VERSION}" \
  --notes-file CHANGELOG.md \
  --latest
```

### Step 10: Cleanup
```bash
git branch -d release/v${NEW_VERSION}
git push origin --delete release/v${NEW_VERSION}
```

### Step 11: Deploy (see /deploy-production)
