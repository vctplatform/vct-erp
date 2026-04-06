package tenant

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func seedResolver() *MemoryResolver {
	r := NewMemoryResolver()
	r.Add(&Tenant{ID: "t1", Name: "VN Martial Arts Federation", Slug: "vnmaf", Plan: "pro", Active: true, CreatedAt: time.Now()})
	r.Add(&Tenant{ID: "t2", Name: "HCM City Club", Slug: "hcm-club", Plan: "free", Active: true, CreatedAt: time.Now()})
	r.Add(&Tenant{ID: "t3", Name: "Suspended Org", Slug: "suspended", Plan: "free", Active: false, CreatedAt: time.Now()})
	return r
}

func TestContext(t *testing.T) {
	ten := &Tenant{ID: "t1", Name: "Test"}
	ctx := WithTenant(context.Background(), ten)

	got := FromContext(ctx)
	if got == nil || got.ID != "t1" {
		t.Error("expected tenant in context")
	}

	if IDFromContext(ctx) != "t1" {
		t.Error("IDFromContext mismatch")
	}

	if FromContext(context.Background()) != nil {
		t.Error("expected nil from empty context")
	}
}

func TestMustFromContext_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	MustFromContext(context.Background())
}

func TestResolverMiddleware_Header(t *testing.T) {
	resolver := seedResolver()
	handler := ResolverMiddleware(resolver, StrategyHeader)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ten := FromContext(r.Context())
		if ten == nil || ten.ID != "t1" {
			t.Error("expected tenant t1")
		}
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/api/athletes", nil)
	req.Header.Set(HeaderTenantID, "t1")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestResolverMiddleware_MissingHeader(t *testing.T) {
	resolver := seedResolver()
	handler := ResolverMiddleware(resolver, StrategyHeader)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/api/athletes", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestResolverMiddleware_NotFound(t *testing.T) {
	resolver := seedResolver()
	handler := ResolverMiddleware(resolver, StrategyHeader)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/api/athletes", nil)
	req.Header.Set(HeaderTenantID, "nonexistent")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestResolverMiddleware_Inactive(t *testing.T) {
	resolver := seedResolver()
	handler := ResolverMiddleware(resolver, StrategyHeader)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/api/athletes", nil)
	req.Header.Set(HeaderTenantID, "t3")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 for inactive tenant, got %d", rec.Code)
	}
}

func TestResolverMiddleware_Subdomain(t *testing.T) {
	resolver := seedResolver()
	handler := ResolverMiddleware(resolver, StrategySubdomain)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ten := FromContext(r.Context())
		if ten == nil || ten.Slug != "vnmaf" {
			t.Error("expected vnmaf tenant")
		}
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/api/athletes", nil)
	req.Host = "vnmaf.vctplatform.com"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestResolverMiddleware_PathPrefix(t *testing.T) {
	resolver := seedResolver()
	handler := ResolverMiddleware(resolver, StrategyPathPrefix)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ten := FromContext(r.Context())
		if ten == nil || ten.Slug != "hcm-club" {
			t.Error("expected hcm-club tenant")
		}
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/t/hcm-club/athletes", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestScopeQuery_WithTenant(t *testing.T) {
	ctx := WithTenant(context.Background(), &Tenant{ID: "t1"})

	q, args := ScopeQuery(ctx, "SELECT * FROM athletes WHERE age > $1", []any{18})
	expected := "SELECT * FROM athletes WHERE age > $1 AND tenant_id = $2"
	if q != expected {
		t.Errorf("expected %q, got %q", expected, q)
	}
	if len(args) != 2 || args[1] != "t1" {
		t.Errorf("args: expected [18, t1], got %v", args)
	}
}

func TestScopeQuery_NoTenant(t *testing.T) {
	q, args := ScopeQuery(context.Background(), "SELECT * FROM athletes", nil)
	if q != "SELECT * FROM athletes" {
		t.Error("should not modify query without tenant")
	}
	if len(args) != 0 {
		t.Error("should not add args")
	}
}

func TestScopeQuery_NoWhereClause(t *testing.T) {
	ctx := WithTenant(context.Background(), &Tenant{ID: "t2"})
	q, args := ScopeQuery(ctx, "SELECT * FROM athletes", nil)
	expected := "SELECT * FROM athletes WHERE tenant_id = $1"
	if q != expected {
		t.Errorf("expected %q, got %q", expected, q)
	}
	if len(args) != 1 {
		t.Error("should add tenant_id arg")
	}
}
