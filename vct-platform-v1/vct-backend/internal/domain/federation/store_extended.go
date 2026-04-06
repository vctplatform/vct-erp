package federation

import (
	"context"
	"sync"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — IN-MEMORY STORES FOR PR / INTL / WORKFLOW
// Thread-safe with seeded demo data.
// ═══════════════════════════════════════════════════════════════

// ── PR Store ─────────────────────────────────────────────────

type memPRStore struct {
	mu       sync.RWMutex
	articles map[string]NewsArticle
}

func NewMemPRStore() PRStore {
	s := &memPRStore{articles: make(map[string]NewsArticle)}
	now := time.Now()
	seed := []NewsArticle{
		{ID: "art-001", Title: "Giải Vô địch Võ Cổ Truyền Toàn quốc 2024 chính thức khởi tranh", Summary: "Giải đấu quy tụ hơn 500 VĐV từ 42 tỉnh/thành.", Category: "Giải đấu", Author: "Ban TT", Status: ArticleStatusPublished, ViewCount: 3420, CreatedAt: now, UpdatedAt: now},
		{ID: "art-002", Title: "Liên đoàn ký kết hợp tác với Liên đoàn Wushu Trung Quốc", Summary: "Thỏa thuận hợp tác đào tạo HLV và trao đổi VĐV.", Category: "Quốc tế", Author: "Ban ĐN", Status: ArticleStatusPublished, ViewCount: 2180, CreatedAt: now, UpdatedAt: now},
		{ID: "art-003", Title: "Khai mạc lớp tập huấn Trọng tài quốc gia 2024", Summary: "120 trọng tài từ 30 tỉnh/thành tham dự.", Category: "Đào tạo", Author: "Ban TT", Status: ArticleStatusPublished, ViewCount: 1560, CreatedAt: now, UpdatedAt: now},
		{ID: "art-004", Title: "Thông báo sửa đổi Luật thi đấu 128/2024", Summary: "Cập nhật điểm số, quy tắc phạt và hạng cân mới.", Category: "Quy chế", Author: "Ban KHVB", Status: ArticleStatusPublished, ViewCount: 4200, CreatedAt: now, UpdatedAt: now},
		{ID: "art-005", Title: "VĐV Bình Định giành 3 HCV tại SEA Games", Summary: "Đoàn Bình Định xuất sắc thi đấu tại SEA Games.", Category: "Thành tích", Author: "Ban TT", Status: ArticleStatusDraft, ViewCount: 0, CreatedAt: now, UpdatedAt: now},
		{ID: "art-006", Title: "Kế hoạch phát triển phong trào Võ Cổ Truyền 2024-2026", Summary: "Chiến lược 3 năm phát triển VCT toàn quốc.", Category: "Chiến lược", Author: "BCH LĐ", Status: ArticleStatusReview, ViewCount: 0, CreatedAt: now, UpdatedAt: now},
	}
	for _, a := range seed {
		s.articles[a.ID] = a
	}
	return s
}

func (s *memPRStore) ListArticles(_ context.Context) ([]NewsArticle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]NewsArticle, 0, len(s.articles))
	for _, a := range s.articles {
		out = append(out, a)
	}
	return out, nil
}
func (s *memPRStore) GetArticle(_ context.Context, id string) (*NewsArticle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.articles[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &a, nil
}
func (s *memPRStore) CreateArticle(_ context.Context, a NewsArticle) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.articles[a.ID] = a
	return nil
}
func (s *memPRStore) UpdateArticle(_ context.Context, a NewsArticle) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.articles[a.ID]; !ok {
		return ErrNotFound
	}
	s.articles[a.ID] = a
	return nil
}
func (s *memPRStore) DeleteArticle(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.articles, id)
	return nil
}

// ── International Store ──────────────────────────────────────

type memIntlStore struct {
	mu       sync.RWMutex
	partners map[string]InternationalPartner
	events   map[string]InternationalEvent
}

func NewMemIntlStore() InternationalStore {
	s := &memIntlStore{
		partners: make(map[string]InternationalPartner),
		events:   make(map[string]InternationalEvent),
	}
	now := time.Now()
	partners := []InternationalPartner{
		{ID: "ip-001", Name: "World Martial Arts Union", Abbreviation: "WoMAU", Country: "Hàn Quốc", CountryCode: "KR", Type: "Liên đoàn Quốc tế", Status: PartnerStatusActive, PartnerSince: "2018", CreatedAt: now, UpdatedAt: now},
		{ID: "ip-002", Name: "Asian Martial Arts Federation", Country: "Nhật Bản", CountryCode: "JP", Type: "Liên đoàn Châu Á", Status: PartnerStatusActive, PartnerSince: "2019", CreatedAt: now, UpdatedAt: now},
		{ID: "ip-003", Name: "Chinese Wushu Association", Country: "Trung Quốc", CountryCode: "CN", Type: "Lưỡng phương", Status: PartnerStatusActive, PartnerSince: "2023", CreatedAt: now, UpdatedAt: now},
		{ID: "ip-004", Name: "SEA Games Federation", Country: "Đông Nam Á", CountryCode: "ASEAN", Type: "Đa phương", Status: PartnerStatusActive, PartnerSince: "2015", CreatedAt: now, UpdatedAt: now},
		{ID: "ip-005", Name: "French Martial Arts Federation", Country: "Pháp", CountryCode: "FR", Type: "Lưỡng phương", Status: PartnerStatusPending, PartnerSince: "2024", CreatedAt: now, UpdatedAt: now},
	}
	events := []InternationalEvent{
		{ID: "ie-001", Name: "SEA Games 2025 — Võ Cổ Truyền", Location: "Bangkok", Country: "Thái Lan", StartDate: "2025-12-01", EndDate: "2025-12-07", AthleteCount: 12, Status: IntlEventPlanning, CreatedAt: now, UpdatedAt: now},
		{ID: "ie-002", Name: "Asian Martial Arts Championship", Location: "Seoul", Country: "Hàn Quốc", StartDate: "2024-08-15", EndDate: "2024-08-20", AthleteCount: 8, MedalGold: 2, MedalSilver: 3, MedalBronze: 1, Status: IntlEventCompleted, CreatedAt: now, UpdatedAt: now},
		{ID: "ie-003", Name: "World Martial Arts Festival", Location: "Chungju", Country: "Hàn Quốc", StartDate: "2024-10-10", EndDate: "2024-10-15", AthleteCount: 15, Status: IntlEventConfirmed, CreatedAt: now, UpdatedAt: now},
	}
	for _, p := range partners {
		s.partners[p.ID] = p
	}
	for _, e := range events {
		s.events[e.ID] = e
	}
	return s
}

func (s *memIntlStore) ListPartners(_ context.Context) ([]InternationalPartner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]InternationalPartner, 0, len(s.partners))
	for _, p := range s.partners {
		out = append(out, p)
	}
	return out, nil
}
func (s *memIntlStore) GetPartner(_ context.Context, id string) (*InternationalPartner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.partners[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &p, nil
}
func (s *memIntlStore) CreatePartner(_ context.Context, p InternationalPartner) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.partners[p.ID] = p
	return nil
}
func (s *memIntlStore) UpdatePartner(_ context.Context, p InternationalPartner) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.partners[p.ID]; !ok {
		return ErrNotFound
	}
	s.partners[p.ID] = p
	return nil
}
func (s *memIntlStore) DeletePartner(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.partners, id)
	return nil
}
func (s *memIntlStore) ListEvents(_ context.Context) ([]InternationalEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]InternationalEvent, 0, len(s.events))
	for _, e := range s.events {
		out = append(out, e)
	}
	return out, nil
}
func (s *memIntlStore) GetEvent(_ context.Context, id string) (*InternationalEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.events[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &e, nil
}
func (s *memIntlStore) CreateEvent(_ context.Context, e InternationalEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.ID] = e
	return nil
}
func (s *memIntlStore) UpdateEvent(_ context.Context, e InternationalEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[e.ID]; !ok {
		return ErrNotFound
	}
	s.events[e.ID] = e
	return nil
}
func (s *memIntlStore) DeleteEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, id)
	return nil
}

// ── Workflow Store ───────────────────────────────────────────

type memWorkflowStore struct {
	mu        sync.RWMutex
	workflows map[string]WorkflowDefinition
}

func NewMemWorkflowStore() WorkflowStore {
	s := &memWorkflowStore{workflows: make(map[string]WorkflowDefinition)}
	now := time.Now()
	seed := []WorkflowDefinition{
		{ID: "wf-001", Code: "club_registration", Name: "Đăng ký CLB mới", Description: "Quy trình phê duyệt thành lập CLB Võ Cổ Truyền", Category: "CLB", IsActive: true, Steps: []WorkflowStep{
			{Order: 1, Name: "Nộp hồ sơ", RoleCode: "club_admin"},
			{Order: 2, Name: "Xét duyệt cấp tỉnh", RoleCode: "provincial_admin"},
			{Order: 3, Name: "Phê duyệt liên đoàn", RoleCode: "federation_secretary"},
		}, CreatedAt: now, UpdatedAt: now},
		{ID: "wf-002", Code: "belt_promotion", Name: "Thi thăng đai", Description: "Quy trình xét duyệt kết quả thi đai từ CLB → Tỉnh → LĐ", Category: "Đai", IsActive: true, Steps: []WorkflowStep{
			{Order: 1, Name: "CLB đề nghị", RoleCode: "club_admin"},
			{Order: 2, Name: "Hội đồng thi", RoleCode: "national_referee"},
			{Order: 3, Name: "Xác nhận tỉnh", RoleCode: "provincial_admin"},
			{Order: 4, Name: "Phê duyệt LĐ", RoleCode: "federation_president"},
		}, CreatedAt: now, UpdatedAt: now},
		{ID: "wf-003", Code: "coach_cert", Name: "Cấp chứng chỉ HLV", Description: "Quy trình xét duyệt và cấp chứng chỉ huấn luyện viên", Category: "HLV", IsActive: true, Steps: []WorkflowStep{
			{Order: 1, Name: "Nộp hồ sơ", RoleCode: "club_admin"},
			{Order: 2, Name: "Kiểm tra năng lực", RoleCode: "national_coach"},
			{Order: 3, Name: "Phê duyệt", RoleCode: "federation_secretary"},
		}, CreatedAt: now, UpdatedAt: now},
		{ID: "wf-004", Code: "referee_cert", Name: "Cấp thẻ Trọng tài", Description: "Quy trình đào tạo và cấp thẻ trọng tài quốc gia", Category: "Trọng tài", IsActive: true, Steps: []WorkflowStep{
			{Order: 1, Name: "Đăng ký", RoleCode: "club_admin"},
			{Order: 2, Name: "Tập huấn", RoleCode: "national_referee"},
			{Order: 3, Name: "Thi sát hạch", RoleCode: "national_referee"},
			{Order: 4, Name: "Xét duyệt", RoleCode: "federation_secretary"},
			{Order: 5, Name: "Cấp thẻ", RoleCode: "federation_president"},
		}, CreatedAt: now, UpdatedAt: now},
		{ID: "wf-005", Code: "tournament_approval", Name: "Phê duyệt Giải đấu", Description: "Quy trình phê duyệt tổ chức giải đấu cấp tỉnh/quốc gia", Category: "Giải đấu", IsActive: true, Steps: []WorkflowStep{
			{Order: 1, Name: "Nộp kế hoạch", RoleCode: "provincial_admin"},
			{Order: 2, Name: "Rà soát kỹ thuật", RoleCode: "national_referee"},
			{Order: 3, Name: "Xét duyệt", RoleCode: "federation_secretary"},
			{Order: 4, Name: "Phê duyệt", RoleCode: "federation_president"},
		}, CreatedAt: now, UpdatedAt: now},
		{ID: "wf-006", Code: "discipline_case", Name: "Xử lý Kỷ luật", Description: "Quy trình điều tra, xét xử và ra quyết định kỷ luật", Category: "Kỷ luật", IsActive: false, Steps: []WorkflowStep{
			{Order: 1, Name: "Tiếp nhận tố cáo", RoleCode: "federation_secretary"},
			{Order: 2, Name: "Điều tra sơ bộ", RoleCode: "federation_secretary"},
			{Order: 3, Name: "Lập hội đồng", RoleCode: "federation_president"},
			{Order: 4, Name: "Xét xử", RoleCode: "federation_president"},
			{Order: 5, Name: "Ra quyết định", RoleCode: "federation_president"},
			{Order: 6, Name: "Lưu trữ", RoleCode: "federation_secretary"},
		}, CreatedAt: now, UpdatedAt: now},
		{ID: "wf-007", Code: "document_publish", Name: "Ban hành Văn bản", Description: "Quy trình soạn thảo, duyệt và ban hành công văn", Category: "Văn bản", IsActive: true, Steps: []WorkflowStep{
			{Order: 1, Name: "Soạn thảo", RoleCode: "federation_secretary"},
			{Order: 2, Name: "Rà soát", RoleCode: "federation_secretary"},
			{Order: 3, Name: "Ký ban hành", RoleCode: "federation_president"},
		}, CreatedAt: now, UpdatedAt: now},
	}
	for _, w := range seed {
		s.workflows[w.ID] = w
	}
	return s
}

func (s *memWorkflowStore) ListWorkflows(_ context.Context) ([]WorkflowDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]WorkflowDefinition, 0, len(s.workflows))
	for _, w := range s.workflows {
		out = append(out, w)
	}
	return out, nil
}
func (s *memWorkflowStore) GetWorkflow(_ context.Context, id string) (*WorkflowDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.workflows[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &w, nil
}
func (s *memWorkflowStore) CreateWorkflow(_ context.Context, w WorkflowDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workflows[w.ID] = w
	return nil
}
func (s *memWorkflowStore) UpdateWorkflow(_ context.Context, w WorkflowDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.workflows[w.ID]; !ok {
		return ErrNotFound
	}
	s.workflows[w.ID] = w
	return nil
}
func (s *memWorkflowStore) DeleteWorkflow(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.workflows, id)
	return nil
}
