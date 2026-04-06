package bulkhead

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestExecute_Success(t *testing.T) {
	b := New("test", 5, time.Second)
	var ran atomic.Int32

	err := b.Execute(context.Background(), func() error {
		ran.Add(1)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
	if ran.Load() != 1 {
		t.Error("function should have run")
	}
	if b.Active() != 0 {
		t.Error("no active after completion")
	}
}

func TestExecute_ConcurrencyLimit(t *testing.T) {
	b := New("test", 2, 50*time.Millisecond)
	var started sync.WaitGroup
	var blocked atomic.Int32

	// Fill the bulkhead
	started.Add(2)
	for i := 0; i < 2; i++ {
		go b.Execute(context.Background(), func() error {
			started.Done()
			time.Sleep(200 * time.Millisecond)
			return nil
		})
	}
	started.Wait()

	// This one should be rejected (timeout)
	err := b.Execute(context.Background(), func() error {
		blocked.Add(1)
		return nil
	})

	if !errors.Is(err, ErrBulkheadFull) {
		t.Errorf("expected ErrBulkheadFull, got %v", err)
	}
	if blocked.Load() != 0 {
		t.Error("blocked function should not have run")
	}
}

func TestExecute_ContextCancellation(t *testing.T) {
	b := New("test", 1, 5*time.Second)

	// Fill the bulkhead
	done := make(chan struct{})
	go b.Execute(context.Background(), func() error {
		<-done
		return nil
	})
	time.Sleep(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := b.Execute(ctx, func() error { return nil })
	close(done)

	if err == nil {
		t.Error("expected cancellation error")
	}
}

func TestStats(t *testing.T) {
	b := New("api", 10, time.Second)

	b.Execute(context.Background(), func() error { return nil })
	b.Execute(context.Background(), func() error { return nil })

	stats := b.Stats()
	if stats.Name != "api" {
		t.Error("name mismatch")
	}
	if stats.Total != 2 {
		t.Errorf("expected 2 total, got %d", stats.Total)
	}
	if stats.Completed != 2 {
		t.Errorf("expected 2 completed, got %d", stats.Completed)
	}
	if stats.Active != 0 {
		t.Error("should have 0 active")
	}
}

func TestGroup(t *testing.T) {
	g := NewGroup(5, time.Second)

	b1 := g.Get("api")
	b2 := g.Get("api")
	if b1 != b2 {
		t.Error("same name should return same bulkhead")
	}

	b3 := g.Get("db")
	if b1 == b3 {
		t.Error("different name should return different bulkhead")
	}

	stats := g.AllStats()
	if len(stats) != 2 {
		t.Errorf("expected 2 bulkheads, got %d", len(stats))
	}
}

func TestHTTPMiddleware_Allowed(t *testing.T) {
	b := New("http", 10, time.Second)

	handler := HTTPMiddleware(b)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHTTPMiddleware_Full(t *testing.T) {
	b := New("http", 1, 50*time.Millisecond)

	handler := HTTPMiddleware(b)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(200)
	}))

	// Fill bulkhead
	go func() {
		req := httptest.NewRequest("GET", "/slow", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}()
	time.Sleep(10 * time.Millisecond)

	// This should get 429
	req := httptest.NewRequest("GET", "/blocked", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", rec.Code)
	}
}
