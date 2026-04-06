package tournament

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — In-Memory MgmtRepository (Test Only)
// Lightweight implementation for domain-level tests.
// ═══════════════════════════════════════════════════════════════

import (
	"context"
	"fmt"
	"sync"
)

// inMemMgmtRepo is a minimal in-memory MgmtRepository for testing.
type inMemMgmtRepo struct {
	mu            sync.RWMutex
	categories    map[string]*Category
	registrations map[string]*Registration
	regAthletes   map[string]*RegistrationAthlete
	scheduleSlots map[string]*ScheduleSlot
	arenaAssigns  map[string]*ArenaAssignment
	results       map[string]*TournamentResult
	standings     map[string]*TeamStanding
}

// NewInMemMgmtRepo creates an in-memory repository for tests.
func NewInMemMgmtRepo() MgmtRepository {
	return &inMemMgmtRepo{
		categories:    make(map[string]*Category),
		registrations: make(map[string]*Registration),
		regAthletes:   make(map[string]*RegistrationAthlete),
		scheduleSlots: make(map[string]*ScheduleSlot),
		arenaAssigns:  make(map[string]*ArenaAssignment),
		results:       make(map[string]*TournamentResult),
		standings:     make(map[string]*TeamStanding),
	}
}

// ── Categories ──
func (r *inMemMgmtRepo) CreateCategory(_ context.Context, c *Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.categories[c.ID] = c
	return nil
}
func (r *inMemMgmtRepo) GetCategory(_ context.Context, id string) (*Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.categories[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return c, nil
}
func (r *inMemMgmtRepo) UpdateCategory(_ context.Context, c *Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.categories[c.ID] = c
	return nil
}
func (r *inMemMgmtRepo) DeleteCategory(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.categories, id)
	return nil
}
func (r *inMemMgmtRepo) ListCategories(_ context.Context, tid string) ([]*Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*Category
	for _, c := range r.categories {
		if c.TournamentID == tid {
			out = append(out, c)
		}
	}
	return out, nil
}

// ── Registrations ──
func (r *inMemMgmtRepo) CreateRegistration(_ context.Context, reg *Registration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registrations[reg.ID] = reg
	return nil
}
func (r *inMemMgmtRepo) GetRegistration(_ context.Context, id string) (*Registration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	reg, ok := r.registrations[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return reg, nil
}
func (r *inMemMgmtRepo) UpdateRegistration(_ context.Context, reg *Registration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registrations[reg.ID] = reg
	return nil
}
func (r *inMemMgmtRepo) ListRegistrations(_ context.Context, tid string) ([]*Registration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*Registration
	for _, reg := range r.registrations {
		if reg.TournamentID == tid {
			out = append(out, reg)
		}
	}
	return out, nil
}

// ── Registration Athletes ──
func (r *inMemMgmtRepo) AddRegistrationAthlete(_ context.Context, a *RegistrationAthlete) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.regAthletes[a.ID] = a
	return nil
}
func (r *inMemMgmtRepo) ListRegistrationAthletes(_ context.Context, regID string) ([]*RegistrationAthlete, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*RegistrationAthlete
	for _, a := range r.regAthletes {
		if a.RegistrationID == regID {
			out = append(out, a)
		}
	}
	return out, nil
}

// ── Schedule ──
func (r *inMemMgmtRepo) CreateScheduleSlot(_ context.Context, s *ScheduleSlot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.scheduleSlots[s.ID] = s
	return nil
}
func (r *inMemMgmtRepo) GetScheduleSlot(_ context.Context, id string) (*ScheduleSlot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.scheduleSlots[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return s, nil
}
func (r *inMemMgmtRepo) UpdateScheduleSlot(_ context.Context, s *ScheduleSlot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.scheduleSlots[s.ID] = s
	return nil
}
func (r *inMemMgmtRepo) DeleteScheduleSlot(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.scheduleSlots, id)
	return nil
}
func (r *inMemMgmtRepo) ListScheduleSlots(_ context.Context, tid string) ([]*ScheduleSlot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*ScheduleSlot
	for _, s := range r.scheduleSlots {
		if s.TournamentID == tid {
			out = append(out, s)
		}
	}
	return out, nil
}

// ── Arena Assignments ──
func (r *inMemMgmtRepo) CreateArenaAssignment(_ context.Context, a *ArenaAssignment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.arenaAssigns[a.ID] = a
	return nil
}
func (r *inMemMgmtRepo) ListArenaAssignments(_ context.Context, tid string) ([]*ArenaAssignment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*ArenaAssignment
	for _, a := range r.arenaAssigns {
		if a.TournamentID == tid {
			out = append(out, a)
		}
	}
	return out, nil
}
func (r *inMemMgmtRepo) DeleteArenaAssignment(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.arenaAssigns, id)
	return nil
}

// ── Results ──
func (r *inMemMgmtRepo) RecordResult(_ context.Context, res *TournamentResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.results[res.ID] = res
	return nil
}
func (r *inMemMgmtRepo) GetResult(_ context.Context, id string) (*TournamentResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res, ok := r.results[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return res, nil
}
func (r *inMemMgmtRepo) UpdateResult(_ context.Context, res *TournamentResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.results[res.ID] = res
	return nil
}
func (r *inMemMgmtRepo) ListResults(_ context.Context, tid string) ([]*TournamentResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*TournamentResult
	for _, res := range r.results {
		if res.TournamentID == tid {
			out = append(out, res)
		}
	}
	return out, nil
}

// ── Team Standings ──
func (r *inMemMgmtRepo) UpsertTeamStanding(_ context.Context, ts *TeamStanding) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.standings[ts.ID] = ts
	return nil
}
func (r *inMemMgmtRepo) ListTeamStandings(_ context.Context, tid string) ([]*TeamStanding, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*TeamStanding
	for _, ts := range r.standings {
		if ts.TournamentID == tid {
			out = append(out, ts)
		}
	}
	return out, nil
}
