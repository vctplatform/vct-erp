package auditlog

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestLog_And_Query(t *testing.T) {
	logger := New()

	logger.Log(context.Background(), &Entry{
		Action:     ActionCreate,
		Resource:   "athlete",
		ResourceID: "a-1",
		Actor:      Actor{Type: ActorUser, ID: "u-1", Name: "Admin"},
		Success:    true,
	})

	entries := logger.Query(Filter{})
	if len(entries) != 1 {
		t.Fatalf("expected 1, got %d", len(entries))
	}
	if entries[0].Resource != "athlete" {
		t.Error("resource mismatch")
	}
}

func TestAutoID(t *testing.T) {
	logger := New()
	logger.Log(context.Background(), &Entry{Action: ActionCreate})
	logger.Log(context.Background(), &Entry{Action: ActionUpdate})

	entries := logger.Query(Filter{})
	if entries[0].ID == entries[1].ID {
		t.Error("IDs should be unique")
	}
}

func TestQuery_FilterByAction(t *testing.T) {
	logger := New()
	logger.Log(context.Background(), &Entry{Action: ActionCreate, Resource: "athlete"})
	logger.Log(context.Background(), &Entry{Action: ActionUpdate, Resource: "athlete"})
	logger.Log(context.Background(), &Entry{Action: ActionDelete, Resource: "athlete"})

	entries := logger.Query(Filter{Action: ActionUpdate})
	if len(entries) != 1 {
		t.Errorf("expected 1 update, got %d", len(entries))
	}
}

func TestQuery_FilterByResource(t *testing.T) {
	logger := New()
	logger.Log(context.Background(), &Entry{Action: ActionCreate, Resource: "athlete"})
	logger.Log(context.Background(), &Entry{Action: ActionCreate, Resource: "tournament"})

	entries := logger.Query(Filter{Resource: "tournament"})
	if len(entries) != 1 {
		t.Errorf("expected 1 tournament, got %d", len(entries))
	}
}

func TestQuery_FilterByActor(t *testing.T) {
	logger := New()
	logger.Log(context.Background(), &Entry{Actor: Actor{ID: "u-1"}})
	logger.Log(context.Background(), &Entry{Actor: Actor{ID: "u-2"}})

	entries := logger.Query(Filter{ActorID: "u-1"})
	if len(entries) != 1 {
		t.Errorf("expected 1, got %d", len(entries))
	}
}

func TestQuery_DateRange(t *testing.T) {
	logger := New()
	now := time.Now()

	logger.Log(context.Background(), &Entry{
		Timestamp: now.Add(-48 * time.Hour),
		Action:    ActionCreate,
	})
	logger.Log(context.Background(), &Entry{
		Timestamp: now.Add(-1 * time.Hour),
		Action:    ActionUpdate,
	})

	entries := logger.Query(Filter{From: now.Add(-24 * time.Hour)})
	if len(entries) != 1 {
		t.Errorf("expected 1 recent entry, got %d", len(entries))
	}
}

func TestQuery_Limit(t *testing.T) {
	logger := New()
	for i := 0; i < 10; i++ {
		logger.Log(context.Background(), &Entry{Action: ActionRead})
	}

	entries := logger.Query(Filter{Limit: 3})
	if len(entries) != 3 {
		t.Errorf("expected 3, got %d", len(entries))
	}
}

func TestBuilder_Fluent(t *testing.T) {
	logger := New()

	logger.Record(ActionUpdate).
		Resource("athlete", "a-1").
		By(ActorUser, "u-1", "Nguyễn Admin").
		IP("192.168.1.1").
		Tenant("t-1").
		Changed("belt", "Đai Đỏ", "Đai Đen").
		Changed("rank", "3", "2").
		Meta("reason", "Promotion").
		Commit(context.Background())

	entries := logger.Query(Filter{})
	if len(entries) != 1 {
		t.Fatal("expected 1 entry")
	}

	e := entries[0]
	if e.Actor.Name != "Nguyễn Admin" {
		t.Error("actor name mismatch")
	}
	if e.TenantID != "t-1" {
		t.Error("tenant mismatch")
	}
	if len(e.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(e.Changes))
	}
	if e.Changes["belt"].To != "Đai Đen" {
		t.Error("belt change mismatch")
	}
}

func TestBuilder_Failed(t *testing.T) {
	logger := New()

	logger.Record(ActionDelete).
		Resource("tournament", "t-1").
		By(ActorUser, "u-1", "User").
		Failed("permission denied").
		Commit(context.Background())

	entries := logger.Query(Filter{})
	if entries[0].Success {
		t.Error("should be failed")
	}
	if entries[0].Error != "permission denied" {
		t.Error("error message mismatch")
	}
}

func TestResourceHistory(t *testing.T) {
	logger := New()
	logger.Log(context.Background(), &Entry{Action: ActionCreate, Resource: "athlete", ResourceID: "a-1"})
	logger.Log(context.Background(), &Entry{Action: ActionUpdate, Resource: "athlete", ResourceID: "a-1"})
	logger.Log(context.Background(), &Entry{Action: ActionCreate, Resource: "athlete", ResourceID: "a-2"})

	history := logger.ResourceHistory("athlete", "a-1")
	if len(history) != 2 {
		t.Errorf("expected 2 entries for a-1, got %d", len(history))
	}
}

func TestSummarize(t *testing.T) {
	e := &Entry{
		Timestamp:  time.Date(2026, 3, 20, 14, 30, 0, 0, time.UTC),
		Action:     ActionUpdate,
		Resource:   "athlete",
		ResourceID: "a-1",
		Actor:      Actor{Name: "Admin"},
		Changes:    map[string]Change{"belt": {From: "Red", To: "Black"}},
		Success:    true,
	}

	summary := Summarize(e)
	if !strings.Contains(summary, "Admin") || !strings.Contains(summary, "belt") {
		t.Errorf("summary missing info: %s", summary)
	}
}

func TestCount(t *testing.T) {
	logger := New()
	logger.Log(context.Background(), &Entry{})
	logger.Log(context.Background(), &Entry{})

	if logger.Count() != 2 {
		t.Errorf("expected count 2, got %d", logger.Count())
	}
}
