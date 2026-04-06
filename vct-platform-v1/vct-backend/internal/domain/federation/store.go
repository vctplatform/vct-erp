package federation

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
)

type MasterDataStore interface {
	// Master Belts
	ListMasterBelts(ctx context.Context) ([]MasterBelt, error)
	GetMasterBelt(ctx context.Context, level string) (*MasterBelt, error)
	CreateMasterBelt(ctx context.Context, belt MasterBelt) error
	UpdateMasterBelt(ctx context.Context, belt MasterBelt) error
	DeleteMasterBelt(ctx context.Context, level string) error

	// Master Weight Classes
	ListMasterWeights(ctx context.Context) ([]MasterWeightClass, error)
	GetMasterWeight(ctx context.Context, id string) (*MasterWeightClass, error)
	CreateMasterWeight(ctx context.Context, weight MasterWeightClass) error
	UpdateMasterWeight(ctx context.Context, weight MasterWeightClass) error
	DeleteMasterWeight(ctx context.Context, id string) error

	// Master Age Groups
	ListMasterAges(ctx context.Context) ([]MasterAgeGroup, error)
	GetMasterAge(ctx context.Context, id string) (*MasterAgeGroup, error)
	CreateMasterAge(ctx context.Context, age MasterAgeGroup) error
	UpdateMasterAge(ctx context.Context, age MasterAgeGroup) error
	DeleteMasterAge(ctx context.Context, id string) error

	// Master Competition Contents (Nội dung thi đấu)
	ListMasterContents(ctx context.Context) ([]MasterCompetitionContent, error)
	GetMasterContent(ctx context.Context, id string) (*MasterCompetitionContent, error)
	CreateMasterContent(ctx context.Context, content MasterCompetitionContent) error
	UpdateMasterContent(ctx context.Context, content MasterCompetitionContent) error
	DeleteMasterContent(ctx context.Context, id string) error

	// Approval Workflows
	ListApprovals(ctx context.Context, status string) ([]ApprovalRequest, error)
	GetApproval(ctx context.Context, id string) (ApprovalRequest, error)
	CreateApproval(ctx context.Context, req ApprovalRequest) error
	UpdateApproval(ctx context.Context, req ApprovalRequest) error
}

type MemoryMasterDataStore struct {
	mu        sync.RWMutex
	belts     []MasterBelt
	weights   []MasterWeightClass
	ages      []MasterAgeGroup
	contents  []MasterCompetitionContent
	approvals map[string]ApprovalRequest
}

func NewMemoryMasterDataStore() *MemoryMasterDataStore {
	store := &MemoryMasterDataStore{
		belts:     make([]MasterBelt, 0),
		weights:   make([]MasterWeightClass, 0),
		ages:      make([]MasterAgeGroup, 0),
		contents:  make([]MasterCompetitionContent, 0),
		approvals: make(map[string]ApprovalRequest),
	}
	store.seedData()
	return store
}

// seedData loads effective regulation data (2021 base + 2024 amendment).
// Đai & Nội dung thi đấu: giữ nguyên 2021 (không thay đổi trong Luật 128/2024)
// Hạng cân & Nhóm tuổi: theo Luật 128/2024 (thay thế hoàn toàn 2021)
func (s *MemoryMasterDataStore) seedData() {
	now := time.Now()

	// Đai: giữ nguyên 2021 (Chương II không thay đổi)
	s.belts = NationalBelts()

	// Hạng cân: theo 2024 (Điều 4 sửa đổi)
	s.weights = EffectiveWeightClasses()

	// Nhóm tuổi: theo 2024 (Điều 4 & 25 sửa đổi)
	s.ages = EffectiveAgeGroups()
	s.contents = NationalCompetitionContents()

	// Seed Approvals — Dữ liệu mẫu
	req := ApprovalRequest{
		ID:            "RQ-001",
		WorkflowCode:  "club_registration",
		EntityType:    "club",
		EntityID:      "clb-abc",
		RequesterID:   "usr-1",
		RequesterName: "CLB Rồng Vàng",
		CurrentStep:   1,
		Status:        RequestPending,
		Notes:         "Xin phép thành lập CLB mới tại quận 1",
		SubmittedAt:   now,
		UpdatedAt:     now,
	}
	s.approvals[req.ID] = req
}

// Removed Organization methods

func (s *MemoryMasterDataStore) ListMasterBelts(ctx context.Context) ([]MasterBelt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]MasterBelt{}, s.belts...), nil
}

func (s *MemoryMasterDataStore) GetMasterBelt(ctx context.Context, level string) (*MasterBelt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.belts {
		if fmt.Sprintf("%d", s.belts[i].Level) == level {
			b := s.belts[i]
			return &b, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryMasterDataStore) CreateMasterBelt(ctx context.Context, belt MasterBelt) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	belt.CreatedAt = time.Now()
	s.belts = append(s.belts, belt)
	return nil
}

func (s *MemoryMasterDataStore) UpdateMasterBelt(ctx context.Context, belt MasterBelt) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.belts {
		if s.belts[i].Level == belt.Level {
			belt.CreatedAt = s.belts[i].CreatedAt
			s.belts[i] = belt
			return nil
		}
	}
	return ErrNotFound
}

func (s *MemoryMasterDataStore) DeleteMasterBelt(ctx context.Context, level string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.belts {
		if fmt.Sprintf("%d", s.belts[i].Level) == level {
			s.belts = append(s.belts[:i], s.belts[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (s *MemoryMasterDataStore) ListMasterWeights(ctx context.Context) ([]MasterWeightClass, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]MasterWeightClass{}, s.weights...), nil
}

func (s *MemoryMasterDataStore) GetMasterWeight(ctx context.Context, id string) (*MasterWeightClass, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.weights {
		if s.weights[i].ID == id {
			w := s.weights[i]
			return &w, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryMasterDataStore) CreateMasterWeight(ctx context.Context, weight MasterWeightClass) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if weight.ID == "" {
		weight.ID = fmt.Sprintf("wt-%d", time.Now().UnixNano())
	}
	weight.CreatedAt = time.Now()
	s.weights = append(s.weights, weight)
	return nil
}

func (s *MemoryMasterDataStore) UpdateMasterWeight(ctx context.Context, weight MasterWeightClass) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.weights {
		if s.weights[i].ID == weight.ID {
			weight.CreatedAt = s.weights[i].CreatedAt
			s.weights[i] = weight
			return nil
		}
	}
	return ErrNotFound
}

func (s *MemoryMasterDataStore) DeleteMasterWeight(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.weights {
		if s.weights[i].ID == id {
			s.weights = append(s.weights[:i], s.weights[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (s *MemoryMasterDataStore) ListMasterAges(ctx context.Context) ([]MasterAgeGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]MasterAgeGroup{}, s.ages...), nil
}

func (s *MemoryMasterDataStore) GetMasterAge(ctx context.Context, id string) (*MasterAgeGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.ages {
		if s.ages[i].ID == id {
			a := s.ages[i]
			return &a, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryMasterDataStore) CreateMasterAge(ctx context.Context, age MasterAgeGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if age.ID == "" {
		age.ID = fmt.Sprintf("age-%d", time.Now().UnixNano())
	}
	age.CreatedAt = time.Now()
	s.ages = append(s.ages, age)
	return nil
}

func (s *MemoryMasterDataStore) UpdateMasterAge(ctx context.Context, age MasterAgeGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.ages {
		if s.ages[i].ID == age.ID {
			age.CreatedAt = s.ages[i].CreatedAt
			s.ages[i] = age
			return nil
		}
	}
	return ErrNotFound
}

func (s *MemoryMasterDataStore) DeleteMasterAge(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.ages {
		if s.ages[i].ID == id {
			s.ages = append(s.ages[:i], s.ages[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

// ── Competition Contents ─────────────────────────────────────

func (s *MemoryMasterDataStore) ListMasterContents(ctx context.Context) ([]MasterCompetitionContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]MasterCompetitionContent{}, s.contents...), nil
}

func (s *MemoryMasterDataStore) GetMasterContent(ctx context.Context, id string) (*MasterCompetitionContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.contents {
		if s.contents[i].ID == id {
			c := s.contents[i]
			return &c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryMasterDataStore) CreateMasterContent(ctx context.Context, content MasterCompetitionContent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if content.ID == "" {
		content.ID = fmt.Sprintf("nd-%d", time.Now().UnixNano())
	}
	content.CreatedAt = time.Now()
	s.contents = append(s.contents, content)
	return nil
}

func (s *MemoryMasterDataStore) UpdateMasterContent(ctx context.Context, content MasterCompetitionContent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.contents {
		if s.contents[i].ID == content.ID {
			content.CreatedAt = s.contents[i].CreatedAt
			s.contents[i] = content
			return nil
		}
	}
	return ErrNotFound
}

func (s *MemoryMasterDataStore) DeleteMasterContent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.contents {
		if s.contents[i].ID == id {
			s.contents = append(s.contents[:i], s.contents[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (s *MemoryMasterDataStore) ListApprovals(ctx context.Context, status string) ([]ApprovalRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var res []ApprovalRequest
	for _, req := range s.approvals {
		if status == "" || string(req.Status) == status {
			res = append(res, req)
		}
	}
	return res, nil
}

func (s *MemoryMasterDataStore) GetApproval(ctx context.Context, id string) (ApprovalRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	req, ok := s.approvals[id]
	if !ok {
		return ApprovalRequest{}, ErrNotFound
	}
	return req, nil
}

func (s *MemoryMasterDataStore) CreateApproval(ctx context.Context, req ApprovalRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req.ID == "" {
		req.ID = fmt.Sprintf("rq-%d", time.Now().UnixNano())
	}
	req.SubmittedAt = time.Now()
	req.UpdatedAt = time.Now()
	s.approvals[req.ID] = req
	return nil
}

func (s *MemoryMasterDataStore) UpdateApproval(ctx context.Context, req ApprovalRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.approvals[req.ID]
	if !ok {
		return ErrNotFound
	}
	req.SubmittedAt = existing.SubmittedAt
	req.UpdatedAt = time.Now()
	s.approvals[req.ID] = req
	return nil
}
