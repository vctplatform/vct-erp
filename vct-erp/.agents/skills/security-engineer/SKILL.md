---
name: security-engineer
description: Security Engineer role - Authentication, authorization, RBAC, OWASP compliance, vulnerability scanning, data protection, and security best practices for VCT Platform.
---

# Security Engineer - VCT Platform

## Role Overview
Ensures the platform is secure from vulnerabilities, implements authentication/authorization via Supabase Auth, manages data protection with Supabase RLS and Neon, and maintains compliance with security standards.

## Core Responsibilities

### 1. Authentication System (Supabase Auth)

#### Supabase Auth Configuration
```
Providers Enabled:
- Email/Password (with email confirmation)
- Google OAuth 2.0
- Facebook OAuth
- Magic Link (passwordless)

Token Lifecycle:
- Access Token (JWT): 1 hour (configurable)
- Refresh Token: 7 days (stored in httpOnly cookie)

MFA: TOTP (Time-based One-Time Password) for admin roles
```

#### Supabase JWT Claims
```go
// Token structure from Supabase Auth
type SupabaseClaims struct {
    Sub      string    `json:"sub"`       // User UUID
    Email    string    `json:"email"`
    Role     string    `json:"role"`      // 'authenticated' or custom
    AppMeta  AppMeta   `json:"app_metadata"`
    UserMeta UserMeta  `json:"user_metadata"`
    jwt.RegisteredClaims
}

type AppMeta struct {
    Provider  string `json:"provider"`
    OrgID     string `json:"org_id"`
    OrgType   string `json:"org_type"`    // "federation", "club", "btc"
    AppRole   string `json:"app_role"`    // "super_admin", "club_admin", etc.
}
```

#### Password Policy (Supabase config)
```
- Minimum 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number
- At least 1 special character
- bcrypt hashing (managed by Supabase)
- Rate limiting on auth endpoints (built-in)
```

### 2. Role-Based Access Control (RBAC)

#### Role Hierarchy
```
super_admin       → Full system access
federation_admin  → National federation management
province_admin    → Provincial federation management
club_admin        → Club management
btc_admin         → Tournament organizing
coach             → Team/athlete management
athlete           → Personal profile & registration
parent            → Child athlete oversight
referee           → Competition scoring
viewer            → Public read-only access
```

#### Custom Claims via Supabase
```sql
-- Set custom claims in Supabase (via database function)
CREATE OR REPLACE FUNCTION set_user_role(user_id UUID, new_role TEXT)
RETURNS VOID AS $$
BEGIN
    UPDATE auth.users
    SET raw_app_meta_data = raw_app_meta_data || jsonb_build_object('app_role', new_role)
    WHERE id = user_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;
```

#### Go Middleware (Supabase JWT)
```go
func RequireRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := GetSupabaseClaimsFromContext(r.Context())
            if claims == nil {
                writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
                return
            }
            
            userRole := claims.AppMeta.AppRole
            allowed := false
            for _, role := range roles {
                if userRole == role {
                    allowed = true
                    break
                }
            }
            
            if !allowed {
                writeError(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

### 3. Row-Level Security (Supabase RLS)

```sql
-- Enable RLS on all data tables
ALTER TABLE athletes ENABLE ROW LEVEL SECURITY;
ALTER TABLE clubs ENABLE ROW LEVEL SECURITY;
ALTER TABLE tournaments ENABLE ROW LEVEL SECURITY;

-- Organization-scoped access
CREATE POLICY "org_members_read" ON athletes
    FOR SELECT
    USING (
        organization_id IN (
            SELECT organization_id FROM user_organizations
            WHERE user_id = auth.uid()
        )
    );

-- Athletes can update own profile
CREATE POLICY "athletes_update_own" ON athletes
    FOR UPDATE 
    USING (user_id = auth.uid())
    WITH CHECK (user_id = auth.uid());

-- Service role bypasses RLS (for backend API)
-- Use SUPABASE_SERVICE_ROLE_KEY only on server-side
```

### 4. OWASP Top 10 Compliance

| # | Vulnerability | Mitigation |
|---|--------------|------------|
| A01 | Broken Access Control | Supabase RLS + RBAC middleware |
| A02 | Cryptographic Failures | AES-256 at rest (Neon/Supabase), TLS 1.3 |
| A03 | Injection | Parameterized queries (pgx), input validation |
| A04 | Insecure Design | Threat modeling, security reviews |
| A05 | Security Misconfiguration | Hardened Docker images, security headers |
| A06 | Vulnerable Components | Dependabot, `go mod tidy`, `npm audit` |
| A07 | Auth Failures | Supabase Auth (MFA, rate limiting, token rotation) |
| A08 | Data Integrity Failures | Signed deployments, integrity checks |
| A09 | Logging Failures | Structured security logging, audit trail |
| A10 | SSRF | URL validation, allowlist external services |

### 5. Security Headers
```go
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
        next.ServeHTTP(w, r)
    })
}
```

### 6. Neon Security

| Feature | Configuration |
|---------|--------------|
| Connection Encryption | SSL/TLS required (enforced by Neon) |
| IP Allowlisting | Restrict DB access to known IPs |
| Branch Isolation | Separate branches for dev/staging/prod |
| Compute Suspension | Auto-suspend idle computes (prevent abuse) |
| Audit Logs | Neon console activity logs |

### 7. Data Protection

#### Sensitive Data Classification
| Level | Data Type | Protection |
|-------|----------|------------|
| **Critical** | Passwords, tokens | Supabase Auth managed, never logged |
| **High** | Email, phone, DOB | Encrypted at rest (Neon/Supabase), masked in logs |
| **Medium** | Name, club membership | Access controlled by RLS + RBAC |
| **Low** | Public results, rankings | Public access OK |

#### Audit Trail
```sql
CREATE TABLE audit_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES auth.users(id),
    action      VARCHAR(50) NOT NULL,  -- CREATE, UPDATE, DELETE, LOGIN, EXPORT
    resource    VARCHAR(100) NOT NULL, -- athletes, tournaments, etc.
    resource_id UUID,
    old_value   JSONB,
    new_value   JSONB,
    ip_address  INET,
    user_agent  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource, resource_id);
```

### 8. Rate Limiting
```go
// API rate limits
var RateLimits = map[string]RateLimit{
    "login":          {Requests: 5,    Window: 15 * time.Minute},
    "register":       {Requests: 3,    Window: 1 * time.Hour},
    "api_general":    {Requests: 100,  Window: 1 * time.Minute},
    "api_export":     {Requests: 10,   Window: 1 * time.Hour},
    "password_reset": {Requests: 3,    Window: 1 * time.Hour},
}
```

### 9. Security Checklist
- [ ] Supabase Auth configured with MFA for admin roles
- [ ] RLS enabled on ALL data tables
- [ ] RBAC enforced at handler level (Go middleware)
- [ ] SQL injection prevention (parameterized queries only)
- [ ] XSS prevention (HTML escaping, CSP headers)
- [ ] CSRF protection (SameSite cookies, CSRF tokens)
- [ ] Rate limiting on auth endpoints (Supabase built-in + custom)
- [ ] Sensitive data encrypted at rest (Neon + Supabase)
- [ ] TLS 1.3 for all traffic (Neon + Supabase enforce SSL)
- [ ] Security headers configured
- [ ] Dependency vulnerabilities scanned (Dependabot)
- [ ] Audit logging for sensitive operations
- [ ] Neon IP allowlisting configured for production
- [ ] Secrets management (no hardcoded secrets)
- [ ] Docker images use non-root user
