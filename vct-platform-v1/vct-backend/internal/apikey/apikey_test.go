package apikey

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreate_And_Validate(t *testing.T) {
	svc := NewService()
	result := svc.Create("Test Key", "user-1", []string{"athletes:read"}, time.Time{})

	if result.RawKey == "" {
		t.Fatal("raw key should be non-empty")
	}
	if result.Key.Prefix != result.RawKey[:8] {
		t.Error("prefix should match first 8 chars")
	}
	if !result.Key.Active {
		t.Error("new key should be active")
	}

	// Validate
	key, err := svc.Validate(result.RawKey)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}
	if key.ID != result.Key.ID {
		t.Error("should return the same key")
	}
}

func TestValidate_InvalidKey(t *testing.T) {
	svc := NewService()
	_, err := svc.Validate("vct_invalid_key_that_does_not_exist")
	if err == nil {
		t.Error("expected error for invalid key")
	}
}

func TestValidate_ExpiredKey(t *testing.T) {
	svc := NewService()
	result := svc.Create("Expired", "user-1", []string{"*"}, time.Now().Add(-time.Hour))

	_, err := svc.Validate(result.RawKey)
	if err == nil {
		t.Error("expected error for expired key")
	}
}

func TestValidate_RevokedKey(t *testing.T) {
	svc := NewService()
	result := svc.Create("Revoked", "user-1", []string{"*"}, time.Time{})

	svc.Revoke(result.Key.ID)

	_, err := svc.Validate(result.RawKey)
	if err == nil {
		t.Error("expected error for revoked key")
	}
}

func TestHasScope(t *testing.T) {
	key := &Key{Scopes: []string{"athletes:read", "matches:*"}}

	if !key.HasScope("athletes:read") {
		t.Error("should have athletes:read")
	}
	if !key.HasScope("matches:write") {
		t.Error("matches:* should match matches:write")
	}
	if key.HasScope("clubs:read") {
		t.Error("should not have clubs:read")
	}

	// Wildcard all
	keyAll := &Key{Scopes: []string{"*"}}
	if !keyAll.HasScope("anything:here") {
		t.Error("* should match everything")
	}
}

func TestList(t *testing.T) {
	svc := NewService()
	svc.Create("Key 1", "user-a", []string{"read"}, time.Time{})
	svc.Create("Key 2", "user-a", []string{"write"}, time.Time{})
	svc.Create("Key 3", "user-b", []string{"*"}, time.Time{})

	keys := svc.List("user-a")
	if len(keys) != 2 {
		t.Errorf("expected 2 keys for user-a, got %d", len(keys))
	}
}

func TestAuthMiddleware_Valid(t *testing.T) {
	svc := NewService()
	result := svc.Create("Test", "user-1", []string{"athletes:read"}, time.Time{})

	handler := AuthMiddleware(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := KeyFromContext(r.Context())
		if key == nil {
			t.Error("key should be in context")
		}
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/api/athletes", nil)
	req.Header.Set("Authorization", "Bearer "+result.RawKey)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAuthMiddleware_NoKey(t *testing.T) {
	svc := NewService()
	handler := AuthMiddleware(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/api/athletes", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestRequireScope_Allowed(t *testing.T) {
	svc := NewService()
	result := svc.Create("Scoped", "user-1", []string{"athletes:read"}, time.Time{})

	handler := RequireScope(svc, "athletes:read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/api/athletes", nil)
	req.Header.Set("Authorization", "ApiKey "+result.RawKey)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRequireScope_Forbidden(t *testing.T) {
	svc := NewService()
	result := svc.Create("Limited", "user-1", []string{"athletes:read"}, time.Time{})

	handler := RequireScope(svc, "system:admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+result.RawKey)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestKeyFromContext_Empty(t *testing.T) {
	if KeyFromContext(context.Background()) != nil {
		t.Error("should return nil from empty context")
	}
}
