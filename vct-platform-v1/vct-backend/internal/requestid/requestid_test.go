package requestid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewID(t *testing.T) {
	id := NewID()
	if len(id) != 16 {
		t.Errorf("expected 16 chars, got %d: %s", len(id), id)
	}

	// Should be unique
	id2 := NewID()
	if id == id2 {
		t.Error("IDs should be unique")
	}
}

func TestNewLongID(t *testing.T) {
	id := NewLongID()
	if len(id) != 32 {
		t.Errorf("expected 32 chars, got %d", len(id))
	}
}

func TestContextHelpers(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-123")
	ctx = WithCorrelationID(ctx, "corr-456")

	if GetRequestID(ctx) != "req-123" {
		t.Error("request ID mismatch")
	}
	if GetCorrelationID(ctx) != "corr-456" {
		t.Error("correlation ID mismatch")
	}

	// Empty context
	if GetRequestID(context.Background()) != "" {
		t.Error("expected empty")
	}
}

func TestMiddleware_GeneratesIDs(t *testing.T) {
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GetRequestID(r.Context())
		corrID := GetCorrelationID(r.Context())
		if reqID == "" {
			t.Error("request ID should be generated")
		}
		if corrID == "" {
			t.Error("correlation ID should be generated")
		}
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get(HeaderRequestID) == "" {
		t.Error("response should have X-Request-ID")
	}
	if rec.Header().Get(HeaderCorrelationID) == "" {
		t.Error("response should have X-Correlation-ID")
	}
}

func TestMiddleware_ReusesClientIDs(t *testing.T) {
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetRequestID(r.Context()) != "client-req-1" {
			t.Error("should reuse client request ID")
		}
		if GetCorrelationID(r.Context()) != "client-corr-1" {
			t.Error("should reuse client correlation ID")
		}
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(HeaderRequestID, "client-req-1")
	req.Header.Set(HeaderCorrelationID, "client-corr-1")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get(HeaderRequestID) != "client-req-1" {
		t.Error("response should echo client request ID")
	}
}

func TestInjectOutgoing(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-abc")
	ctx = WithCorrelationID(ctx, "corr-xyz")

	outReq, _ := http.NewRequest("GET", "http://downstream/api", nil)
	InjectOutgoing(ctx, outReq)

	if outReq.Header.Get(HeaderRequestID) != "req-abc" {
		t.Error("should inject request ID")
	}
	if outReq.Header.Get(HeaderCorrelationID) != "corr-xyz" {
		t.Error("should inject correlation ID")
	}
}

func TestSanitize(t *testing.T) {
	if sanitize("  abc  ") != "abc" {
		t.Error("should trim spaces")
	}

	long := ""
	for i := 0; i < 100; i++ {
		long += "x"
	}
	if len(sanitize(long)) != 64 {
		t.Error("should truncate to 64 chars")
	}
}
