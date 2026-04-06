package athlete

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vct-platform/backend/internal/domain"
	"vct-platform/backend/internal/domain/athlete"
)

type repoMem struct {
	mu       sync.RWMutex
	athletes map[string]domain.Athlete
}

func newMemRepository() athlete.Repository {
	return &repoMem{
		athletes: make(map[string]domain.Athlete),
	}
}

func (r *repoMem) List(ctx context.Context) ([]domain.Athlete, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]domain.Athlete, 0, len(r.athletes))
	for _, a := range r.athletes {
		res = append(res, a)
	}
	return res, nil
}

func (r *repoMem) GetByID(ctx context.Context, id string) (*domain.Athlete, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.athletes[id]
	if !ok {
		return nil, fmt.Errorf("athlete not found")
	}
	return &a, nil
}

func (r *repoMem) Create(ctx context.Context, a domain.Athlete) (*domain.Athlete, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if a.ID == "" {
		a.ID = fmt.Sprintf("vdv_%d", time.Now().UnixNano())
	}
	r.athletes[a.ID] = a
	return &a, nil
}

func (r *repoMem) Update(ctx context.Context, id string, patch map[string]interface{}) (*domain.Athlete, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.athletes[id]
	if !ok {
		return nil, fmt.Errorf("athlete not found")
	}

	// Simple patch implementation for memory repo
	if ht, ok := patch["ho_ten"].(string); ok {
		a.HoTen = ht
	}
	if st, ok := patch["trang_thai"].(string); ok {
		a.TrangThai = domain.TrangThaiVDV(st)
	}
	// ... add more as needed for testing

	r.athletes[id] = a
	return &a, nil
}

func (r *repoMem) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.athletes, id)
	return nil
}

func (r *repoMem) ListByTeam(ctx context.Context, teamID string) ([]domain.Athlete, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := []domain.Athlete{}
	for _, a := range r.athletes {
		if a.DoanID == teamID {
			res = append(res, a)
		}
	}
	return res, nil
}

func (r *repoMem) ListByTournament(ctx context.Context, tournamentID string) ([]domain.Athlete, error) {
	// In-memory doesn't track tournament relationship currently, return all for now or empty
	return []domain.Athlete{}, nil
}
