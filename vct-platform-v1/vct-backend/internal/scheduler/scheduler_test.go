package scheduler

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestParseCron_Wildcard(t *testing.T) {
	s, err := ParseCron("* * * * *")
	if err != nil {
		t.Fatal(err)
	}
	// All wildcards = matches any time
	if !s.Matches(time.Now()) {
		t.Error("* * * * * should match any time")
	}
}

func TestParseCron_Exact(t *testing.T) {
	s, err := ParseCron("30 14 * * *")
	if err != nil {
		t.Fatal(err)
	}

	match := time.Date(2026, 3, 20, 14, 30, 0, 0, time.UTC)
	if !s.Matches(match) {
		t.Error("should match 14:30")
	}

	noMatch := time.Date(2026, 3, 20, 14, 31, 0, 0, time.UTC)
	if s.Matches(noMatch) {
		t.Error("should not match 14:31")
	}
}

func TestParseCron_Step(t *testing.T) {
	s, err := ParseCron("*/15 * * * *")
	if err != nil {
		t.Fatal(err)
	}

	// Minutes: 0, 15, 30, 45
	for _, m := range []int{0, 15, 30, 45} {
		tt := time.Date(2026, 1, 1, 0, m, 0, 0, time.UTC)
		if !s.Matches(tt) {
			t.Errorf("should match minute %d", m)
		}
	}

	tt := time.Date(2026, 1, 1, 0, 10, 0, 0, time.UTC)
	if s.Matches(tt) {
		t.Error("should not match minute 10")
	}
}

func TestParseCron_Range(t *testing.T) {
	s, err := ParseCron("0 9-17 * * *")
	if err != nil {
		t.Fatal(err)
	}

	for h := 9; h <= 17; h++ {
		tt := time.Date(2026, 1, 1, h, 0, 0, 0, time.UTC)
		if !s.Matches(tt) {
			t.Errorf("should match hour %d", h)
		}
	}

	tt := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)
	if s.Matches(tt) {
		t.Error("should not match hour 8")
	}
}

func TestParseCron_List(t *testing.T) {
	s, err := ParseCron("0 0 1,15 * *")
	if err != nil {
		t.Fatal(err)
	}

	d1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	d15 := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	d10 := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)

	if !s.Matches(d1) || !s.Matches(d15) {
		t.Error("should match day 1 and 15")
	}
	if s.Matches(d10) {
		t.Error("should not match day 10")
	}
}

func TestParseCron_Shortcuts(t *testing.T) {
	tests := []struct {
		expr string
	}{
		{"@yearly"},
		{"@monthly"},
		{"@weekly"},
		{"@daily"},
		{"@hourly"},
	}

	for _, tt := range tests {
		_, err := ParseCron(tt.expr)
		if err != nil {
			t.Errorf("%s: %v", tt.expr, err)
		}
	}
}

func TestParseCron_Invalid(t *testing.T) {
	_, err := ParseCron("bad expression")
	if err == nil {
		t.Error("expected error")
	}
}

func TestScheduler_AddAndRunNow(t *testing.T) {
	sched := New(testLogger())
	var ran atomic.Int32

	sched.Add("test-job", "* * * * *", func(ctx context.Context) error {
		ran.Add(1)
		return nil
	})

	sched.RunNow(context.Background(), "test-job")

	if ran.Load() != 1 {
		t.Errorf("expected 1 run, got %d", ran.Load())
	}
}

func TestScheduler_RunNow_NotFound(t *testing.T) {
	sched := New(testLogger())
	err := sched.RunNow(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for missing job")
	}
}

func TestScheduler_JobStats(t *testing.T) {
	sched := New(testLogger())
	sched.Add("success-job", "@daily", func(ctx context.Context) error {
		return nil
	})
	sched.Add("fail-job", "@daily", func(ctx context.Context) error {
		return errors.New("database down")
	})

	sched.RunNow(context.Background(), "success-job")
	sched.RunNow(context.Background(), "fail-job")

	stats := sched.AllStats()
	if len(stats) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(stats))
	}

	for _, s := range stats {
		if s.Name == "success-job" && s.Failures != 0 {
			t.Error("success-job should have 0 failures")
		}
		if s.Name == "fail-job" && s.Failures != 1 {
			t.Error("fail-job should have 1 failure")
		}
	}
}

func TestScheduler_WeekdayMatch(t *testing.T) {
	// "0 9 * * 1" = Monday at 9:00
	s, err := ParseCron("0 9 * * 1")
	if err != nil {
		t.Fatal(err)
	}

	// 2026-03-23 is a Monday
	monday := time.Date(2026, 3, 23, 9, 0, 0, 0, time.UTC)
	if !s.Matches(monday) {
		t.Error("should match Monday 9:00")
	}

	// 2026-03-24 is a Tuesday
	tuesday := time.Date(2026, 3, 24, 9, 0, 0, 0, time.UTC)
	if s.Matches(tuesday) {
		t.Error("should not match Tuesday")
	}
}
