---
description: Security audit workflow for VCT Platform - OWASP, dependency vulnerabilities, Supabase RLS, and Neon security
---

# /security-audit - Security Audit Workflow

## When to Use
- Before each release
- Monthly scheduled audit
- After adding new tables/RLS policies
- After dependency updates
- After security incident

## Steps

### Step 1: Go vulnerability scan
```bash
cd backend
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### Step 2: npm dependency audit
```bash
# Web
cd apps/web && npm audit --audit-level=high

# Mobile
cd apps/mobile && npm audit --audit-level=high
```

### Step 3: Docker image scan
```bash
# Scan backend image
docker scout cves ghcr.io/vct-platform/api:latest

# Scan frontend image
docker scout cves ghcr.io/vct-platform/web:latest
```

### Step 4: Supabase RLS audit
```sql
-- Check tables WITHOUT RLS enabled
SELECT schemaname, tablename
FROM pg_tables
WHERE schemaname = 'public'
AND tablename NOT IN (
    SELECT tablename FROM pg_tables t
    JOIN pg_class c ON c.relname = t.tablename
    WHERE c.relrowsecurity = true
);

-- Check tables with RLS enabled but NO policies
SELECT c.relname AS table_name
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = 'public'
AND c.relkind = 'r'
AND c.relrowsecurity = true
AND NOT EXISTS (
    SELECT 1 FROM pg_policy p WHERE p.polrelid = c.oid
);

-- List all RLS policies
SELECT schemaname, tablename, policyname, permissive, roles, cmd, qual
FROM pg_policies
WHERE schemaname = 'public'
ORDER BY tablename, policyname;
```

### Step 5: Supabase Auth configuration audit
```bash
# Check auth providers
supabase functions list

# Verify JWT settings
echo "Check in Supabase Dashboard → Authentication → Settings:"
echo "  - JWT expiry ≤ 3600s (1 hour)"
echo "  - Refresh token rotation enabled"
echo "  - MFA available for admin roles"
echo "  - Email confirmations enabled"
echo "  - Rate limiting configured"
```

### Step 6: Neon security audit
```bash
# Check IP allowlist
echo "Verify in Neon Dashboard → Project Settings → IP Allow:"
echo "  - Only known IPs/CIDRs allowed"
echo "  - No 0.0.0.0/0 in production"

# Check compute configuration
echo "Verify in Neon Dashboard → Branches → main:"
echo "  - Suspend after idle: configured"
echo "  - SSL mode: require"
```

### Step 7: OWASP Top 10 checklist
```markdown
- [ ] **A01: Broken Access Control** → Supabase RLS + RBAC middleware
- [ ] **A02: Cryptographic Failures** → TLS 1.3, AES-256 at rest (Neon/Supabase)
- [ ] **A03: Injection** → pgx parameterized queries, input validation
- [ ] **A04: Insecure Design** → Clean Architecture, threat modeling
- [ ] **A05: Security Misconfiguration** → No default passwords, CORS restricted
- [ ] **A06: Vulnerable Components** → govulncheck, npm audit
- [ ] **A07: Auth Failures** → Supabase Auth, MFA, rate limiting
- [ ] **A08: Data Integrity** → CSRF tokens, signed uploads (Supabase Storage)
- [ ] **A09: Logging Failures** → zerolog, audit trail, Neon audit logs
- [ ] **A10: SSRF** → URL validation, no user-controlled URLs in server requests
```

### Step 8: API security checks
```bash
# Check CORS configuration
curl -v -X OPTIONS https://api.vct-platform.com/api/v1/athletes \
  -H "Origin: https://evil.com" \
  -H "Access-Control-Request-Method: GET" 2>&1 | grep "Access-Control"

# Check rate limiting
for i in $(seq 1 20); do
  curl -s -o /dev/null -w "%{http_code}\n" https://api.vct-platform.com/api/v1/athletes
done

# Check security headers
curl -sI https://api.vct-platform.com | grep -iE "(strict-transport|x-content-type|x-frame|content-security)"
```

### Step 9: Secrets audit
```bash
# Check for leaked secrets in codebase
grep -rn "eyJ" . --include="*.go" --include="*.ts" --include="*.tsx" --include="*.env" | grep -v node_modules | grep -v ".git"
grep -rn "password\s*=" . --include="*.go" --include="*.ts" | grep -v node_modules | grep -v "_test" | grep -v ".example"

# Check .env files not in .gitignore
git ls-files | grep -E "\.env$|\.env\.local$"
```

### Step 10: Generate security report
```bash
echo "=== Security Audit Report ==="
echo "Date: $(date)"
echo "Auditor: $(git config user.name)"
echo ""
echo "Go Vulnerabilities: $(govulncheck ./... 2>&1 | grep -c 'Vulnerability')"
echo "npm High/Critical: $(npm audit 2>/dev/null | grep -c 'high\|critical')"
echo ""
echo "Submit report to: security@vct-platform.com"
```

## Severity Matrix

| Finding | Severity | SLA |
|---------|----------|-----|
| RLS missing on public table | 🔴 Critical | Fix immediately |
| Known CVE in dependency | 🔴 Critical | 24 hours |
| Exposed secret | 🔴 Critical | Rotate immediately |
| Missing security header | 🟡 Medium | Next sprint |
| Outdated dependency (no CVE) | 🟢 Low | Next dependency update cycle |
