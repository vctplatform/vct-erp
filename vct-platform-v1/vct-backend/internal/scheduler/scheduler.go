// Package scheduler provides a lightweight cron scheduler with expression
// parsing, named job registry, concurrent execution, and job statistics.
package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// Cron Expression Parser
// ═══════════════════════════════════════════════════════════════

// Schedule represents a parsed cron expression.
// Format: minute hour day-of-month month day-of-week
type Schedule struct {
	Minutes  []int // 0-59
	Hours    []int // 0-23
	Days     []int // 1-31
	Months   []int // 1-12
	Weekdays []int // 0-6 (Sunday=0)
}

// ParseCron parses a cron expression (5 fields).
// Supports: * (all), */n (step), n-m (range), n,m (list), n (exact).
func ParseCron(expr string) (*Schedule, error) {
	// Shortcuts
	switch expr {
	case "@yearly", "@annually":
		expr = "0 0 1 1 *"
	case "@monthly":
		expr = "0 0 1 * *"
	case "@weekly":
		expr = "0 0 * * 0"
	case "@daily", "@midnight":
		expr = "0 0 * * *"
	case "@hourly":
		expr = "0 * * * *"
	}

	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("cron: expected 5 fields, got %d", len(parts))
	}

	s := &Schedule{}
	var err error

	if s.Minutes, err = parseField(parts[0], 0, 59); err != nil {
		return nil, fmt.Errorf("cron minute: %w", err)
	}
	if s.Hours, err = parseField(parts[1], 0, 23); err != nil {
		return nil, fmt.Errorf("cron hour: %w", err)
	}
	if s.Days, err = parseField(parts[2], 1, 31); err != nil {
		return nil, fmt.Errorf("cron day: %w", err)
	}
	if s.Months, err = parseField(parts[3], 1, 12); err != nil {
		return nil, fmt.Errorf("cron month: %w", err)
	}
	if s.Weekdays, err = parseField(parts[4], 0, 6); err != nil {
		return nil, fmt.Errorf("cron weekday: %w", err)
	}

	return s, nil
}

// Matches checks if a time matches this schedule.
func (s *Schedule) Matches(t time.Time) bool {
	return contains(s.Minutes, t.Minute()) &&
		contains(s.Hours, t.Hour()) &&
		contains(s.Days, t.Day()) &&
		contains(s.Months, int(t.Month())) &&
		contains(s.Weekdays, int(t.Weekday()))
}

func contains(vals []int, v int) bool {
	if vals == nil {
		return true // wildcard
	}
	for _, x := range vals {
		if x == v {
			return true
		}
	}
	return false
}

func parseField(field string, min, max int) ([]int, error) {
	if field == "*" {
		return nil, nil // wildcard
	}

	var result []int

	for _, part := range strings.Split(field, ",") {
		// Step: */n or n-m/s
		if strings.Contains(part, "/") {
			sub := strings.SplitN(part, "/", 2)
			step, err := strconv.Atoi(sub[1])
			if err != nil || step <= 0 {
				return nil, fmt.Errorf("invalid step: %s", part)
			}
			start, end := min, max
			if sub[0] != "*" {
				rangeParts := strings.SplitN(sub[0], "-", 2)
				start, err = strconv.Atoi(rangeParts[0])
				if err != nil {
					return nil, fmt.Errorf("invalid range start: %s", part)
				}
				if len(rangeParts) == 2 {
					end, err = strconv.Atoi(rangeParts[1])
					if err != nil {
						return nil, fmt.Errorf("invalid range end: %s", part)
					}
				}
			}
			for i := start; i <= end; i += step {
				result = append(result, i)
			}
			continue
		}

		// Range: n-m
		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			continue
		}

		// Exact value
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid value: %s", part)
		}
		if val < min || val > max {
			return nil, fmt.Errorf("value %d out of range [%d, %d]", val, min, max)
		}
		result = append(result, val)
	}

	return result, nil
}

// ═══════════════════════════════════════════════════════════════
// Job
// ═══════════════════════════════════════════════════════════════

// Job is a scheduled task.
type Job struct {
	Name     string
	Schedule *Schedule
	Fn       func(ctx context.Context) error

	// Stats
	runs     atomic.Int64
	failures atomic.Int64
	lastRun  time.Time
	lastErr  error
	mu       sync.Mutex
}

// JobStats holds job execution statistics.
type JobStats struct {
	Name     string    `json:"name"`
	Runs     int64     `json:"runs"`
	Failures int64     `json:"failures"`
	LastRun  time.Time `json:"last_run,omitempty"`
	LastErr  string    `json:"last_error,omitempty"`
}

func (j *Job) Stats() JobStats {
	j.mu.Lock()
	defer j.mu.Unlock()
	s := JobStats{
		Name:     j.Name,
		Runs:     j.runs.Load(),
		Failures: j.failures.Load(),
		LastRun:  j.lastRun,
	}
	if j.lastErr != nil {
		s.LastErr = j.lastErr.Error()
	}
	return s
}

// ═══════════════════════════════════════════════════════════════
// Scheduler
// ═══════════════════════════════════════════════════════════════

// Scheduler runs jobs on their cron schedules.
type Scheduler struct {
	jobs   []*Job
	logger *slog.Logger
	mu     sync.Mutex
	cancel context.CancelFunc
}

// New creates a scheduler.
func New(logger *slog.Logger) *Scheduler {
	return &Scheduler{
		logger: logger.With(slog.String("component", "scheduler")),
	}
}

// Add registers a job with a cron expression.
func (s *Scheduler) Add(name, cronExpr string, fn func(ctx context.Context) error) error {
	sched, err := ParseCron(cronExpr)
	if err != nil {
		return fmt.Errorf("job %q: %w", name, err)
	}

	s.mu.Lock()
	s.jobs = append(s.jobs, &Job{
		Name:     name,
		Schedule: sched,
		Fn:       fn,
	})
	s.mu.Unlock()

	s.logger.Info("job registered", "name", name, "cron", cronExpr)
	return nil
}

// Start begins the scheduler loop. Ticks every minute.
func (s *Scheduler) Start(ctx context.Context) {
	ctx, s.cancel = context.WithCancel(ctx)

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				s.logger.Info("scheduler stopped")
				return
			case t := <-ticker.C:
				s.tick(ctx, t)
			}
		}
	}()

	s.logger.Info("scheduler started", "jobs", len(s.jobs))
}

// Stop halts the scheduler.
func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// RunNow executes a job immediately by name (for testing/manual trigger).
func (s *Scheduler) RunNow(ctx context.Context, name string) error {
	s.mu.Lock()
	var job *Job
	for _, j := range s.jobs {
		if j.Name == name {
			job = j
			break
		}
	}
	s.mu.Unlock()

	if job == nil {
		return fmt.Errorf("job %q not found", name)
	}

	s.executeJob(ctx, job)
	return nil
}

// AllStats returns stats for all jobs.
func (s *Scheduler) AllStats() []JobStats {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats := make([]JobStats, len(s.jobs))
	for i, j := range s.jobs {
		stats[i] = j.Stats()
	}
	return stats
}

func (s *Scheduler) tick(ctx context.Context, t time.Time) {
	s.mu.Lock()
	jobs := make([]*Job, len(s.jobs))
	copy(jobs, s.jobs)
	s.mu.Unlock()

	for _, job := range jobs {
		if job.Schedule.Matches(t) {
			go s.executeJob(ctx, job)
		}
	}
}

func (s *Scheduler) executeJob(ctx context.Context, job *Job) {
	job.runs.Add(1)
	job.mu.Lock()
	job.lastRun = time.Now()
	job.mu.Unlock()

	if err := job.Fn(ctx); err != nil {
		job.failures.Add(1)
		job.mu.Lock()
		job.lastErr = err
		job.mu.Unlock()
		s.logger.Error("job failed", "name", job.Name, "error", err)
	}
}
