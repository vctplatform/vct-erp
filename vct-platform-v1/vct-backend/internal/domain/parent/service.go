package parent

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// DOMAIN MODELS
// ═══════════════════════════════════════════════════════════════

// LinkStatus represents the state of a parent-athlete link.
type LinkStatus string

const (
	LinkStatusPending  LinkStatus = "pending"
	LinkStatusApproved LinkStatus = "approved"
	LinkStatusRejected LinkStatus = "rejected"
)

// ValidRelations enumerates allowed relation values.
var ValidRelations = map[string]bool{
	"cha": true, "mẹ": true, "người giám hộ": true,
	"ông": true, "bà": true, "anh/chị": true,
}

// ParentLink associates a parent user with an athlete (child).
type ParentLink struct {
	ID          string     `json:"id"`
	ParentID    string     `json:"parent_id"`
	ParentName  string     `json:"parent_name"`
	AthleteID   string     `json:"athlete_id"`
	AthleteName string     `json:"athlete_name"`
	ClubName    string     `json:"club_name"`
	BeltLevel   string     `json:"belt_level"`
	Relation    string     `json:"relation"` // cha/mẹ/người giám hộ
	Status      LinkStatus `json:"status"`
	RequestedAt time.Time  `json:"requested_at"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
}

// ConsentType categorises what the consent covers.
type ConsentType string

const (
	ConsentTournament ConsentType = "tournament"
	ConsentBeltExam   ConsentType = "belt_exam"
	ConsentMedical    ConsentType = "medical"
	ConsentPhotoUsage ConsentType = "photo_usage"
	ConsentTraining   ConsentType = "training"
)

// ValidConsentTypes enumerates allowed consent type values.
var ValidConsentTypes = map[ConsentType]bool{
	ConsentTournament: true, ConsentBeltExam: true, ConsentMedical: true,
	ConsentPhotoUsage: true, ConsentTraining: true,
}

// ConsentStatus tracks whether the consent is active.
type ConsentStatus string

const (
	ConsentActive  ConsentStatus = "active"
	ConsentRevoked ConsentStatus = "revoked"
	ConsentExpired ConsentStatus = "expired"
)

// ConsentRecord stores a parent's signed e-consent.
type ConsentRecord struct {
	ID          string        `json:"id"`
	ParentID    string        `json:"parent_id"`
	AthleteID   string        `json:"athlete_id"`
	AthleteName string        `json:"athlete_name"`
	Type        ConsentType   `json:"type"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      ConsentStatus `json:"status"`
	SignedAt    time.Time     `json:"signed_at"`
	ExpiresAt   *time.Time    `json:"expires_at,omitempty"`
	RevokedAt   *time.Time    `json:"revoked_at,omitempty"`
}

// AttendanceSummary is a lightweight view of attendance for a child.
type AttendanceSummary struct {
	Date    string `json:"date"`
	Session string `json:"session"`
	Status  string `json:"status"` // present/absent/late
	Coach   string `json:"coach"`
}

// ChildResult shows a competition result for the child.
type ChildResult struct {
	Tournament string `json:"tournament"`
	Category   string `json:"category"`
	Result     string `json:"result"` // gold/silver/bronze/eliminated
	Date       string `json:"date"`
}

// Dashboard aggregates the parent overview.
type Dashboard struct {
	ChildrenCount   int           `json:"children_count"`
	PendingConsents int           `json:"pending_consents"`
	ActiveConsents  int           `json:"active_consents"`
	UpcomingEvents  int           `json:"upcoming_events"`
	Children        []ParentLink  `json:"children"`
	RecentResults   []ChildResult `json:"recent_results"`
}

// ═══════════════════════════════════════════════════════════════
// TYPED UPDATE STRUCTS (replaces map[string]interface{})
// ═══════════════════════════════════════════════════════════════

// LinkUpdate carries typed fields for updating a ParentLink.
type LinkUpdate struct {
	Status     *LinkStatus
	ApprovedAt *time.Time
}

// ConsentUpdate carries typed fields for updating a ConsentRecord.
type ConsentUpdate struct {
	Status    *ConsentStatus
	RevokedAt *time.Time
}

// ═══════════════════════════════════════════════════════════════
// IN-MEMORY STORES
// ═══════════════════════════════════════════════════════════════

// ── ParentLink Store ─────────────────────────────────────────

type InMemParentLinkStore struct {
	mu    sync.RWMutex
	links map[string]ParentLink
}

func NewInMemParentLinkStore() *InMemParentLinkStore {
	s := &InMemParentLinkStore{links: make(map[string]ParentLink)}
	s.seed()
	return s
}

func (s *InMemParentLinkStore) seed() {
	now := time.Now()
	approved := now.Add(-30 * 24 * time.Hour)
	s.links = map[string]ParentLink{
		"PL-001": {
			ID: "PL-001", ParentID: "PARENT-001", ParentName: "Nguyễn Thị Phụ Huynh",
			AthleteID: "ATH-001", AthleteName: "Nguyễn Văn An", ClubName: "CLB Thanh Long",
			BeltLevel: "Hoàng đai", Relation: "mẹ", Status: LinkStatusApproved,
			RequestedAt: now.Add(-60 * 24 * time.Hour), ApprovedAt: &approved,
		},
		"PL-002": {
			ID: "PL-002", ParentID: "PARENT-001", ParentName: "Nguyễn Thị Phụ Huynh",
			AthleteID: "ATH-002", AthleteName: "Nguyễn Thị Bình", ClubName: "CLB Thanh Long",
			BeltLevel: "Lam đai", Relation: "mẹ", Status: LinkStatusApproved,
			RequestedAt: now.Add(-45 * 24 * time.Hour), ApprovedAt: &approved,
		},
		"PL-003": {
			ID: "PL-003", ParentID: "PARENT-002", ParentName: "Trần Văn Hùng",
			AthleteID: "ATH-003", AthleteName: "Trần Minh Đức", ClubName: "CLB Bạch Hổ",
			BeltLevel: "Vàng đai 1", Relation: "cha", Status: LinkStatusPending,
			RequestedAt: now.Add(-2 * 24 * time.Hour),
		},
	}
}

func (s *InMemParentLinkStore) ListByParent(_ context.Context, parentID string) ([]ParentLink, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []ParentLink
	for _, l := range s.links {
		if l.ParentID == parentID {
			out = append(out, l)
		}
	}
	return out, nil
}

func (s *InMemParentLinkStore) GetByID(_ context.Context, id string) (*ParentLink, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.links[id]
	if !ok {
		return nil, fmt.Errorf("parent link %s not found", id)
	}
	return &l, nil
}

func (s *InMemParentLinkStore) Create(_ context.Context, l ParentLink) (*ParentLink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.links[l.ID] = l
	return &l, nil
}

// Update applies typed patch to a parent link.
func (s *InMemParentLinkStore) Update(_ context.Context, id string, patch LinkUpdate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, ok := s.links[id]
	if !ok {
		return fmt.Errorf("parent link %s not found", id)
	}
	if patch.Status != nil {
		l.Status = *patch.Status
	}
	if patch.ApprovedAt != nil {
		l.ApprovedAt = patch.ApprovedAt
	}
	s.links[id] = l
	return nil
}

func (s *InMemParentLinkStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.links[id]; !ok {
		return fmt.Errorf("parent link %s not found", id)
	}
	delete(s.links, id)
	return nil
}

// IsChildOfParent checks whether a given athlete is linked to the parent.
func (s *InMemParentLinkStore) IsChildOfParent(_ context.Context, parentID, athleteID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, l := range s.links {
		if l.ParentID == parentID && l.AthleteID == athleteID && l.Status == LinkStatusApproved {
			return true
		}
	}
	return false
}

// ── Consent Store ────────────────────────────────────────────

type InMemConsentStore struct {
	mu    sync.RWMutex
	items map[string]ConsentRecord
}

func NewInMemConsentStore() *InMemConsentStore {
	s := &InMemConsentStore{items: make(map[string]ConsentRecord)}
	s.seed()
	return s
}

func (s *InMemConsentStore) seed() {
	now := time.Now()
	expires := now.Add(365 * 24 * time.Hour)
	s.items = map[string]ConsentRecord{
		"CS-001": {
			ID: "CS-001", ParentID: "PARENT-001", AthleteID: "ATH-001", AthleteName: "Nguyễn Văn An",
			Type: ConsentTournament, Title: "Đồng ý tham gia Giải Vovinam Toàn Quốc 2026",
			Description: "Cho phép con em tham gia giải đấu Vovinam Toàn Quốc 2026 tại TP.HCM",
			Status:      ConsentActive, SignedAt: now.Add(-10 * 24 * time.Hour), ExpiresAt: &expires,
		},
		"CS-002": {
			ID: "CS-002", ParentID: "PARENT-001", AthleteID: "ATH-001", AthleteName: "Nguyễn Văn An",
			Type: ConsentMedical, Title: "Đồng ý khám sức khỏe & sơ cứu y tế",
			Description: "Cho phép nhân viên y tế giải thực hiện sơ cứu và khám sức khỏe cho con em",
			Status:      ConsentActive, SignedAt: now.Add(-10 * 24 * time.Hour), ExpiresAt: &expires,
		},
		"CS-003": {
			ID: "CS-003", ParentID: "PARENT-001", AthleteID: "ATH-002", AthleteName: "Nguyễn Thị Bình",
			Type: ConsentBeltExam, Title: "Đồng ý thi lên đai Lam đai 2",
			Description: "Cho phép con em tham gia kỳ thi thăng đai Lam đai 2 tại CLB Thanh Long",
			Status:      ConsentActive, SignedAt: now.Add(-5 * 24 * time.Hour), ExpiresAt: &expires,
		},
		"CS-004": {
			ID: "CS-004", ParentID: "PARENT-001", AthleteID: "ATH-001", AthleteName: "Nguyễn Văn An",
			Type: ConsentPhotoUsage, Title: "Sử dụng hình ảnh thi đấu",
			Description: "Cho phép BTC sử dụng hình ảnh/video thi đấu của con em cho mục đích truyền thông",
			Status:      ConsentRevoked, SignedAt: now.Add(-30 * 24 * time.Hour),
		},
	}
}

func (s *InMemConsentStore) ListByParent(_ context.Context, parentID string) ([]ConsentRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []ConsentRecord
	for _, c := range s.items {
		if c.ParentID == parentID {
			out = append(out, c)
		}
	}
	return out, nil
}

func (s *InMemConsentStore) ListByAthlete(_ context.Context, athleteID string) ([]ConsentRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []ConsentRecord
	for _, c := range s.items {
		if c.AthleteID == athleteID {
			out = append(out, c)
		}
	}
	return out, nil
}

func (s *InMemConsentStore) GetByID(_ context.Context, id string) (*ConsentRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.items[id]
	if !ok {
		return nil, fmt.Errorf("consent %s not found", id)
	}
	return &c, nil
}

func (s *InMemConsentStore) Create(_ context.Context, c ConsentRecord) (*ConsentRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[c.ID] = c
	return &c, nil
}

// Update applies a typed patch to a consent record.
func (s *InMemConsentStore) Update(_ context.Context, id string, patch ConsentUpdate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.items[id]
	if !ok {
		return fmt.Errorf("consent %s not found", id)
	}
	if patch.Status != nil {
		c.Status = *patch.Status
	}
	if patch.RevokedAt != nil {
		c.RevokedAt = patch.RevokedAt
	}
	s.items[id] = c
	return nil
}

// ── Attendance Store (seeded) ────────────────────────────────

type InMemAttendanceStore struct {
	mu      sync.RWMutex
	records map[string][]AttendanceSummary // key: athleteID
}

func NewInMemAttendanceStore() *InMemAttendanceStore {
	s := &InMemAttendanceStore{records: make(map[string][]AttendanceSummary)}
	s.seed()
	return s
}

func (s *InMemAttendanceStore) seed() {
	s.records["ATH-001"] = []AttendanceSummary{
		{Date: "2026-03-10", Session: "Sáng 07:00–09:00", Status: "present", Coach: "HLV Trần Văn Minh"},
		{Date: "2026-03-08", Session: "Chiều 16:00–18:00", Status: "present", Coach: "HLV Trần Văn Minh"},
		{Date: "2026-03-06", Session: "Sáng 07:00–09:00", Status: "late", Coach: "HLV Lê Thị Hoa"},
		{Date: "2026-03-04", Session: "Chiều 16:00–18:00", Status: "absent", Coach: "HLV Trần Văn Minh"},
		{Date: "2026-03-02", Session: "Sáng 07:00–09:00", Status: "present", Coach: "HLV Lê Thị Hoa"},
		{Date: "2026-02-28", Session: "Chiều 16:00–18:00", Status: "present", Coach: "HLV Trần Văn Minh"},
		{Date: "2026-02-26", Session: "Sáng 07:00–09:00", Status: "present", Coach: "HLV Lê Thị Hoa"},
		{Date: "2026-02-24", Session: "Chiều 16:00–18:00", Status: "present", Coach: "HLV Trần Văn Minh"},
	}
	s.records["ATH-002"] = []AttendanceSummary{
		{Date: "2026-03-10", Session: "Sáng 07:00–09:00", Status: "present", Coach: "HLV Lê Thị Hoa"},
		{Date: "2026-03-08", Session: "Chiều 16:00–18:00", Status: "present", Coach: "HLV Lê Thị Hoa"},
		{Date: "2026-03-06", Session: "Sáng 07:00–09:00", Status: "present", Coach: "HLV Lê Thị Hoa"},
		{Date: "2026-03-04", Session: "Chiều 16:00–18:00", Status: "late", Coach: "HLV Lê Thị Hoa"},
		{Date: "2026-03-02", Session: "Sáng 07:00–09:00", Status: "present", Coach: "HLV Lê Thị Hoa"},
	}
}

func (s *InMemAttendanceStore) ListByAthlete(_ context.Context, athleteID string) ([]AttendanceSummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	records := s.records[athleteID]
	if records == nil {
		return []AttendanceSummary{}, nil
	}
	return records, nil
}

// ── Results Store (seeded) ───────────────────────────────────

type InMemResultStore struct {
	mu      sync.RWMutex
	records map[string][]ChildResult // key: athleteID
}

func NewInMemResultStore() *InMemResultStore {
	s := &InMemResultStore{records: make(map[string][]ChildResult)}
	s.seed()
	return s
}

func (s *InMemResultStore) seed() {
	s.records["ATH-001"] = []ChildResult{
		{Tournament: "Giải Vovinam TP.HCM Mở rộng 2026", Category: "Đối kháng Nam 52kg", Result: "🥇 Huy chương vàng", Date: "2026-02-15"},
		{Tournament: "Giải Vovinam Học sinh 2025", Category: "Quyền Nam Thiếu niên", Result: "🥈 Huy chương bạc", Date: "2025-11-20"},
		{Tournament: "Giải CLB Thanh Long 2025", Category: "Đối kháng Nam 48kg", Result: "🥉 Huy chương đồng", Date: "2025-09-10"},
	}
	s.records["ATH-002"] = []ChildResult{
		{Tournament: "Giải Vovinam Học sinh 2025", Category: "Quyền Nữ Thiếu niên", Result: "🥈 Huy chương bạc", Date: "2025-11-20"},
	}
}

func (s *InMemResultStore) ListByAthlete(_ context.Context, athleteID string) ([]ChildResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	records := s.records[athleteID]
	if records == nil {
		return []ChildResult{}, nil
	}
	return records, nil
}

// ═══════════════════════════════════════════════════════════════
// STORE INTERFACES (replaceable with PostgreSQL adapters)
// ═══════════════════════════════════════════════════════════════

// ParentLinkStore persists parent-athlete links.
type ParentLinkStore interface {
	ListByParent(ctx context.Context, parentID string) ([]ParentLink, error)
	GetByID(ctx context.Context, id string) (*ParentLink, error)
	Create(ctx context.Context, l ParentLink) (*ParentLink, error)
	Update(ctx context.Context, id string, patch LinkUpdate) error
	Delete(ctx context.Context, id string) error
	IsChildOfParent(ctx context.Context, parentID, athleteID string) bool
}

// ConsentStore persists parent consent records.
type ConsentStore interface {
	ListByParent(ctx context.Context, parentID string) ([]ConsentRecord, error)
	ListByAthlete(ctx context.Context, athleteID string) ([]ConsentRecord, error)
	GetByID(ctx context.Context, id string) (*ConsentRecord, error)
	Create(ctx context.Context, c ConsentRecord) (*ConsentRecord, error)
	Update(ctx context.Context, id string, patch ConsentUpdate) error
}

// AttendanceStore reads child attendance records.
type AttendanceStore interface {
	ListByAthlete(ctx context.Context, athleteID string) ([]AttendanceSummary, error)
}

// ResultStore reads child competition results.
type ResultStore interface {
	ListByAthlete(ctx context.Context, athleteID string) ([]ChildResult, error)
}

// ═══════════════════════════════════════════════════════════════
// SERVICE
// ═══════════════════════════════════════════════════════════════

type Service struct {
	linkStore       ParentLinkStore
	consentStore    ConsentStore
	attendanceStore AttendanceStore
	resultStore     ResultStore
	genID           func() string
}

func NewService(
	linkStore ParentLinkStore,
	consentStore ConsentStore,
	attendanceStore AttendanceStore,
	resultStore ResultStore,
	genID func() string,
) *Service {
	return &Service{
		linkStore:       linkStore,
		consentStore:    consentStore,
		attendanceStore: attendanceStore,
		resultStore:     resultStore,
		genID:           genID,
	}
}

// ── Link Management ──────────────────────────────────────────

// ListMyChildren returns all approved athletes linked to a parent.
func (svc *Service) ListMyChildren(ctx context.Context, parentID string) ([]ParentLink, error) {
	links, err := svc.linkStore.ListByParent(ctx, parentID)
	if err != nil {
		return nil, err
	}
	var approved []ParentLink
	for _, l := range links {
		if l.Status == LinkStatusApproved {
			approved = append(approved, l)
		}
	}
	return approved, nil
}

// ListAllLinks returns all links (including pending) for a parent.
func (svc *Service) ListAllLinks(ctx context.Context, parentID string) ([]ParentLink, error) {
	return svc.linkStore.ListByParent(ctx, parentID)
}

// IsChildOfParent checks if an athlete belongs to a parent's approved links.
func (svc *Service) IsChildOfParent(ctx context.Context, parentID, athleteID string) bool {
	return svc.linkStore.IsChildOfParent(ctx, parentID, athleteID)
}

// RequestLink creates a new pending link between parent & athlete with validation.
func (svc *Service) RequestLink(ctx context.Context, link ParentLink) (*ParentLink, error) {
	// Validate required fields
	if strings.TrimSpace(link.AthleteID) == "" {
		return nil, fmt.Errorf("athlete_id is required")
	}
	if strings.TrimSpace(link.AthleteName) == "" {
		return nil, fmt.Errorf("athlete_name is required")
	}
	if !ValidRelations[link.Relation] {
		return nil, fmt.Errorf("invalid relation %q; allowed: cha, mẹ, người giám hộ, ông, bà, anh/chị", link.Relation)
	}

	link.ID = svc.genID()
	link.Status = LinkStatusPending
	link.RequestedAt = time.Now()
	return svc.linkStore.Create(ctx, link)
}

// ApproveLink sets a pending link to approved.
func (svc *Service) ApproveLink(ctx context.Context, linkID string) error {
	now := time.Now()
	status := LinkStatusApproved
	return svc.linkStore.Update(ctx, linkID, LinkUpdate{
		Status:     &status,
		ApprovedAt: &now,
	})
}

// GetLinkByID returns a single parent-child link by ID.
func (svc *Service) GetLinkByID(ctx context.Context, linkID string) (*ParentLink, error) {
	return svc.linkStore.GetByID(ctx, linkID)
}

// DeleteLink removes a parent-child link.
func (svc *Service) DeleteLink(ctx context.Context, linkID string) error {
	return svc.linkStore.Delete(ctx, linkID)
}

// ── Consent Management ───────────────────────────────────────

// ListConsents returns all consent records for a parent.
func (svc *Service) ListConsents(ctx context.Context, parentID string) ([]ConsentRecord, error) {
	return svc.consentStore.ListByParent(ctx, parentID)
}

// CreateConsent signs a new e-consent record with validation.
func (svc *Service) CreateConsent(ctx context.Context, c ConsentRecord) (*ConsentRecord, error) {
	if strings.TrimSpace(c.AthleteID) == "" {
		return nil, fmt.Errorf("athlete_id is required")
	}
	if strings.TrimSpace(c.Title) == "" {
		return nil, fmt.Errorf("title is required")
	}
	if !ValidConsentTypes[c.Type] {
		return nil, fmt.Errorf("invalid consent type %q", c.Type)
	}

	c.ID = svc.genID()
	c.Status = ConsentActive
	c.SignedAt = time.Now()
	return svc.consentStore.Create(ctx, c)
}

// RevokeConsent marks a consent as revoked.
func (svc *Service) RevokeConsent(ctx context.Context, consentID, parentID string) error {
	// Verify ownership
	consent, err := svc.consentStore.GetByID(ctx, consentID)
	if err != nil {
		return err
	}
	if consent.ParentID != parentID {
		return fmt.Errorf("consent %s does not belong to parent %s", consentID, parentID)
	}
	if consent.Status != ConsentActive {
		return fmt.Errorf("consent %s is not active (current: %s)", consentID, consent.Status)
	}

	now := time.Now()
	status := ConsentRevoked
	return svc.consentStore.Update(ctx, consentID, ConsentUpdate{
		Status:    &status,
		RevokedAt: &now,
	})
}

// ── Child Data Access ────────────────────────────────────────

// GetChildAttendance returns attendance records from store for a given athlete.
func (svc *Service) GetChildAttendance(ctx context.Context, athleteID string) ([]AttendanceSummary, error) {
	return svc.attendanceStore.ListByAthlete(ctx, athleteID)
}

// GetChildResults returns competition results from store for a given athlete.
func (svc *Service) GetChildResults(ctx context.Context, athleteID string) ([]ChildResult, error) {
	return svc.resultStore.ListByAthlete(ctx, athleteID)
}

// ── Dashboard ────────────────────────────────────────────────

// GetDashboard returns aggregated data for the parent with proper error handling.
func (svc *Service) GetDashboard(ctx context.Context, parentID string) (*Dashboard, error) {
	children, err := svc.ListMyChildren(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("list children: %w", err)
	}

	consents, err := svc.ListConsents(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("list consents: %w", err)
	}

	// Aggregate results from all children
	var allResults []ChildResult
	for _, child := range children {
		results, err := svc.GetChildResults(ctx, child.AthleteID)
		if err != nil {
			return nil, fmt.Errorf("get results for %s: %w", child.AthleteID, err)
		}
		allResults = append(allResults, results...)
	}

	// Count consent statuses
	activeConsents := 0
	pendingConsents := 0
	for _, c := range consents {
		switch c.Status {
		case ConsentActive:
			activeConsents++
		}
	}

	// Count pending links
	allLinks, err := svc.ListAllLinks(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("list all links: %w", err)
	}
	for _, l := range allLinks {
		if l.Status == LinkStatusPending {
			pendingConsents++
		}
	}

	// Count upcoming events (from active tournament consents)
	upcomingEvents := 0
	for _, c := range consents {
		if c.Status == ConsentActive && (c.Type == ConsentTournament || c.Type == ConsentBeltExam) {
			upcomingEvents++
		}
	}

	return &Dashboard{
		ChildrenCount:   len(children),
		PendingConsents: pendingConsents,
		ActiveConsents:  activeConsents,
		UpcomingEvents:  upcomingEvents,
		Children:        children,
		RecentResults:   allResults,
	}, nil
}
