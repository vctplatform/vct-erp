package federation

import (
	"context"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — PR / INTERNATIONAL / WORKFLOW MODELS + STORES
// ═══════════════════════════════════════════════════════════════

// ── PR: News Articles ────────────────────────────────────────

type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"
	ArticleStatusReview    ArticleStatus = "review"
	ArticleStatusPublished ArticleStatus = "published"
)

type NewsArticle struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Summary     string        `json:"summary"`
	Content     string        `json:"content"`
	Category    string        `json:"category"` // Giải đấu, Đào tạo, Quy chế, Thành tích, Quốc tế, Chiến lược
	ImageURL    string        `json:"image_url"`
	Author      string        `json:"author"`
	AuthorID    string        `json:"author_id"`
	Status      ArticleStatus `json:"status"`
	PublishedAt *time.Time    `json:"published_at,omitempty"`
	ViewCount   int           `json:"view_count"`
	Tags        []string      `json:"tags,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// ── International: Partners ──────────────────────────────────

type PartnerStatus string

const (
	PartnerStatusActive  PartnerStatus = "active"
	PartnerStatusPending PartnerStatus = "pending"
	PartnerStatusExpired PartnerStatus = "expired"
)

type InternationalPartner struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`         // e.g. "World Martial Arts Union"
	Abbreviation string        `json:"abbreviation"` // e.g. "WoMAU"
	Country      string        `json:"country"`      // e.g. "Hàn Quốc"
	CountryCode  string        `json:"country_code"` // e.g. "KR"
	Type         string        `json:"type"`         // Liên đoàn Quốc tế, Lưỡng phương, Đa phương
	ContactName  string        `json:"contact_name"`
	ContactEmail string        `json:"contact_email"`
	Website      string        `json:"website"`
	Status       PartnerStatus `json:"status"`
	PartnerSince string        `json:"partner_since"` // YYYY
	Notes        string        `json:"notes,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// ── International: Events ────────────────────────────────────

type IntlEventStatus string

const (
	IntlEventPlanning  IntlEventStatus = "planning"
	IntlEventConfirmed IntlEventStatus = "confirmed"
	IntlEventOngoing   IntlEventStatus = "ongoing"
	IntlEventCompleted IntlEventStatus = "completed"
)

type InternationalEvent struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Location     string          `json:"location"` // City, Country
	Country      string          `json:"country"`
	StartDate    string          `json:"start_date"` // YYYY-MM-DD
	EndDate      string          `json:"end_date"`
	AthleteCount int             `json:"athlete_count"` // VĐV tham gia
	CoachCount   int             `json:"coach_count"`
	MedalGold    int             `json:"medal_gold"`
	MedalSilver  int             `json:"medal_silver"`
	MedalBronze  int             `json:"medal_bronze"`
	Status       IntlEventStatus `json:"status"`
	Description  string          `json:"description,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ── Workflow Definitions ─────────────────────────────────────

type WorkflowDefinition struct {
	ID          string         `json:"id"`
	Code        string         `json:"code"` // e.g. "club_registration"
	Name        string         `json:"name"` // e.g. "Đăng ký CLB mới"
	Description string         `json:"description"`
	Category    string         `json:"category"` // CLB, Đai, HLV, Trọng tài, Giải đấu, Kỷ luật, Văn bản
	Steps       []WorkflowStep `json:"steps"`
	IsActive    bool           `json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type WorkflowStep struct {
	Order       int    `json:"order"` // 1, 2, 3...
	Name        string `json:"name"`  // e.g. "Nộp hồ sơ"
	Description string `json:"description"`
	RoleCode    string `json:"role_code"` // required role to execute
	AutoApprove bool   `json:"auto_approve"`
}

// ── Store Interfaces ─────────────────────────────────────────

type PRStore interface {
	ListArticles(ctx context.Context) ([]NewsArticle, error)
	GetArticle(ctx context.Context, id string) (*NewsArticle, error)
	CreateArticle(ctx context.Context, a NewsArticle) error
	UpdateArticle(ctx context.Context, a NewsArticle) error
	DeleteArticle(ctx context.Context, id string) error
}

type InternationalStore interface {
	ListPartners(ctx context.Context) ([]InternationalPartner, error)
	GetPartner(ctx context.Context, id string) (*InternationalPartner, error)
	CreatePartner(ctx context.Context, p InternationalPartner) error
	UpdatePartner(ctx context.Context, p InternationalPartner) error
	DeletePartner(ctx context.Context, id string) error
	ListEvents(ctx context.Context) ([]InternationalEvent, error)
	GetEvent(ctx context.Context, id string) (*InternationalEvent, error)
	CreateEvent(ctx context.Context, e InternationalEvent) error
	UpdateEvent(ctx context.Context, e InternationalEvent) error
	DeleteEvent(ctx context.Context, id string) error
}

type WorkflowStore interface {
	ListWorkflows(ctx context.Context) ([]WorkflowDefinition, error)
	GetWorkflow(ctx context.Context, id string) (*WorkflowDefinition, error)
	CreateWorkflow(ctx context.Context, w WorkflowDefinition) error
	UpdateWorkflow(ctx context.Context, w WorkflowDefinition) error
	DeleteWorkflow(ctx context.Context, id string) error
}
