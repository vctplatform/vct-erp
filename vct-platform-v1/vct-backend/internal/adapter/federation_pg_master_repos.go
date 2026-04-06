package adapter

import (
	"context"
	"fmt"
	"time"

	"vct-platform/backend/internal/domain/federation"
	"vct-platform/backend/internal/store"
)

// ── MasterDataStore (Postgres) ──

type pgMasterDataStore struct {
	beltStore     *StoreAdapter[federation.MasterBelt]
	weightStore   *StoreAdapter[federation.MasterWeightClass]
	ageStore      *StoreAdapter[federation.MasterAgeGroup]
	contentStore  *StoreAdapter[federation.MasterCompetitionContent]
	approvalStore *StoreAdapter[federation.ApprovalRequest]
}

func NewGenericMasterDataStore(ds store.DataStore) federation.MasterDataStore {
	return &pgMasterDataStore{
		beltStore:     NewStoreAdapter[federation.MasterBelt](ds, "federation_master_belts"),
		weightStore:   NewStoreAdapter[federation.MasterWeightClass](ds, "federation_master_weights"),
		ageStore:      NewStoreAdapter[federation.MasterAgeGroup](ds, "federation_master_ages"),
		contentStore:  NewStoreAdapter[federation.MasterCompetitionContent](ds, "federation_master_contents"),
		approvalStore: NewStoreAdapter[federation.ApprovalRequest](ds, "federation_approvals"),
	}
}

// ── Belts ──

func (s *pgMasterDataStore) ListMasterBelts(ctx context.Context) ([]federation.MasterBelt, error) {
	return s.beltStore.List()
}

func (s *pgMasterDataStore) GetMasterBelt(ctx context.Context, level string) (*federation.MasterBelt, error) {
	items, err := s.beltStore.List()
	if err != nil {
		return nil, err
	}
	for i := range items {
		if fmt.Sprintf("%d", items[i].Level) == level {
			return &items[i], nil
		}
	}
	return nil, federation.ErrNotFound
}

func (s *pgMasterDataStore) CreateMasterBelt(ctx context.Context, belt federation.MasterBelt) error {
	belt.CreatedAt = time.Now()
	_, err := s.beltStore.Create(belt)
	return err
}

func (s *pgMasterDataStore) UpdateMasterBelt(ctx context.Context, belt federation.MasterBelt) error {
	_, err := s.beltStore.Update(fmt.Sprintf("%d", belt.Level), map[string]interface{}{
		"name":              belt.Name,
		"color_hex":         belt.ColorHex,
		"required_time_min": belt.RequiredTimeMin,
		"is_dan_level":      belt.IsDanLevel,
		"description":       belt.Description,
		"scope":             belt.Scope,
		"scope_id":          belt.ScopeID,
		"inherits_from":     belt.InheritsFrom,
	})
	return err
}

func (s *pgMasterDataStore) DeleteMasterBelt(ctx context.Context, level string) error {
	return s.beltStore.Delete(level)
}

// ── Weights ──

func (s *pgMasterDataStore) ListMasterWeights(ctx context.Context) ([]federation.MasterWeightClass, error) {
	return s.weightStore.List()
}

func (s *pgMasterDataStore) GetMasterWeight(ctx context.Context, id string) (*federation.MasterWeightClass, error) {
	return s.weightStore.GetByID(id)
}

func (s *pgMasterDataStore) CreateMasterWeight(ctx context.Context, weight federation.MasterWeightClass) error {
	weight.CreatedAt = time.Now()
	_, err := s.weightStore.Create(weight)
	return err
}

func (s *pgMasterDataStore) UpdateMasterWeight(ctx context.Context, weight federation.MasterWeightClass) error {
	_, err := s.weightStore.Update(weight.ID, map[string]interface{}{
		"gender":        weight.Gender,
		"category":      weight.Category,
		"max_weight":    weight.MaxWeight,
		"is_heavy":      weight.IsHeavy,
		"scope":         weight.Scope,
		"scope_id":      weight.ScopeID,
		"inherits_from": weight.InheritsFrom,
	})
	return err
}

func (s *pgMasterDataStore) DeleteMasterWeight(ctx context.Context, id string) error {
	return s.weightStore.Delete(id)
}

// ── Ages ──

func (s *pgMasterDataStore) ListMasterAges(ctx context.Context) ([]federation.MasterAgeGroup, error) {
	return s.ageStore.List()
}

func (s *pgMasterDataStore) GetMasterAge(ctx context.Context, id string) (*federation.MasterAgeGroup, error) {
	return s.ageStore.GetByID(id)
}

func (s *pgMasterDataStore) CreateMasterAge(ctx context.Context, age federation.MasterAgeGroup) error {
	age.CreatedAt = time.Now()
	_, err := s.ageStore.Create(age)
	return err
}

func (s *pgMasterDataStore) UpdateMasterAge(ctx context.Context, age federation.MasterAgeGroup) error {
	_, err := s.ageStore.Update(age.ID, map[string]interface{}{
		"name":          age.Name,
		"min_age":       age.MinAge,
		"max_age":       age.MaxAge,
		"scope":         age.Scope,
		"scope_id":      age.ScopeID,
		"inherits_from": age.InheritsFrom,
	})
	return err
}

func (s *pgMasterDataStore) DeleteMasterAge(ctx context.Context, id string) error {
	return s.ageStore.Delete(id)
}

// ── Competition Contents ──

func (s *pgMasterDataStore) ListMasterContents(ctx context.Context) ([]federation.MasterCompetitionContent, error) {
	return s.contentStore.List()
}

func (s *pgMasterDataStore) GetMasterContent(ctx context.Context, id string) (*federation.MasterCompetitionContent, error) {
	return s.contentStore.GetByID(id)
}

func (s *pgMasterDataStore) CreateMasterContent(ctx context.Context, content federation.MasterCompetitionContent) error {
	content.CreatedAt = time.Now()
	_, err := s.contentStore.Create(content)
	return err
}

func (s *pgMasterDataStore) UpdateMasterContent(ctx context.Context, content federation.MasterCompetitionContent) error {
	_, err := s.contentStore.Update(content.ID, map[string]interface{}{
		"code":            content.Code,
		"name":            content.Name,
		"description":     content.Description,
		"requires_weight": content.RequiresWeight,
		"is_team_event":   content.IsTeamEvent,
		"min_athletes":    content.MinAthletes,
		"max_athletes":    content.MaxAthletes,
		"has_weapon":      content.HasWeapon,
		"scope":           content.Scope,
		"scope_id":        content.ScopeID,
	})
	return err
}

func (s *pgMasterDataStore) DeleteMasterContent(ctx context.Context, id string) error {
	return s.contentStore.Delete(id)
}

// ── Approvals ──

func (s *pgMasterDataStore) ListApprovals(ctx context.Context, status string) ([]federation.ApprovalRequest, error) {
	items, err := s.approvalStore.List()
	if err != nil {
		return nil, err
	}
	var res []federation.ApprovalRequest
	for _, a := range items {
		if status == "" || string(a.Status) == status {
			res = append(res, a)
		}
	}
	return res, nil
}

func (s *pgMasterDataStore) GetApproval(ctx context.Context, id string) (federation.ApprovalRequest, error) {
	item, err := s.approvalStore.GetByID(id)
	if err != nil {
		return federation.ApprovalRequest{}, err
	}
	if item == nil {
		return federation.ApprovalRequest{}, federation.ErrNotFound
	}
	return *item, nil
}

func (s *pgMasterDataStore) CreateApproval(ctx context.Context, req federation.ApprovalRequest) error {
	req.SubmittedAt = time.Now()
	req.UpdatedAt = time.Now()
	_, err := s.approvalStore.Create(req)
	return err
}

func (s *pgMasterDataStore) UpdateApproval(ctx context.Context, req federation.ApprovalRequest) error {
	_, err := s.approvalStore.Update(req.ID, map[string]interface{}{
		"status":     req.Status,
		"notes":      req.Notes,
		"updated_at": time.Now(),
	})
	return err
}
