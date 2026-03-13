---
description: Code review checklist and process for VCT Platform pull requests
---

# /code-review - Code Review Process

## When to Trigger
- When a Pull Request is created or updated
- Before merging any branch into `develop` or `main`

## Review Checklist

### 1. Architecture & Design
- [ ] Follows Clean Architecture boundaries (dependencies point inward)
- [ ] Domain layer has no external dependencies
- [ ] Use cases orchestrate business logic, not handlers
- [ ] No business logic in HTTP handlers
- [ ] New code fits existing module structure

### 2. Code Quality
- [ ] `golangci-lint` passes with zero errors
- [ ] `eslint` and TypeScript compile with zero errors
- [ ] No `TODO` or `FIXME` without linked issue
- [ ] No commented-out code
- [ ] No hardcoded values (use constants/config)
- [ ] Functions are concise (< 50 lines)
- [ ] Cyclomatic complexity ≤ 10

### 3. Naming
- [ ] Go: follows Go naming conventions (PascalCase exports, camelCase private)
- [ ] TypeScript: follows project conventions (PascalCase components, camelCase functions)
- [ ] SQL: follows naming conventions (snake_case tables/columns)
- [ ] Variables/functions named clearly (no single-letter except loops)

### 4. Testing
- [ ] New code has unit tests
- [ ] Test coverage ≥ 80% for changed code
- [ ] Edge cases covered (nil, empty, boundary values)
- [ ] Mocks used appropriately (not over-mocking)
- [ ] Table-driven tests for Go

### 5. Security
- [ ] No SQL injection vulnerabilities (parameterized queries)
- [ ] No hardcoded secrets or credentials
- [ ] Auth/permission checks on all protected endpoints
- [ ] Input validation on all user inputs
- [ ] No sensitive data in logs

### 6. API Changes
- [ ] API documentation updated (Swagger/OpenAPI)
- [ ] Backward compatible (or breaking change documented)
- [ ] Standard response format used
- [ ] Proper HTTP status codes
- [ ] Pagination for list endpoints

### 7. Database Changes
- [ ] Migration files have UP and DOWN
- [ ] Indexes added for query patterns
- [ ] No N+1 query issues
- [ ] Soft delete pattern used

### 8. Frontend
- [ ] All text uses i18n (no hardcoded Vietnamese/English)
- [ ] Responsive design works (320px to 1440px)
- [ ] Loading states for async operations
- [ ] Error handling with user-friendly messages
- [ ] Accessibility: proper ARIA, keyboard nav

### 9. Performance
- [ ] No unnecessary database queries in loops
- [ ] Pagination used for all list queries
- [ ] Heavy operations are async/background
- [ ] No memory leaks (goroutine leaks, unclosed resources)

## Review Process

1. **Self-review first**: Author checks all items above
2. **Request review**: Assign ≥ 1 reviewer
3. **Respond to feedback**: Address all comments
4. **Get approval**: At least 1 approval required
5. **Merge**: Squash merge into target branch
