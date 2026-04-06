package federation

import (
	"context"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — SERVICE METHODS FOR PR / INTL / WORKFLOW
// Extends the main Service with optional stores.
// ═══════════════════════════════════════════════════════════════

// SetExtendedStores injects PR, International, and Workflow stores after construction.
func (s *Service) SetExtendedStores(pr PRStore, intl InternationalStore, wf WorkflowStore) {
	s.prStore = pr
	s.intlStore = intl
	s.workflowStore = wf
}

// ── PR Service ───────────────────────────────────────────────

func (s *Service) ListArticles(ctx context.Context) ([]NewsArticle, error) {
	if s.prStore == nil {
		return nil, nil
	}
	return s.prStore.ListArticles(ctx)
}

func (s *Service) GetArticle(ctx context.Context, id string) (*NewsArticle, error) {
	return s.prStore.GetArticle(ctx, id)
}

func (s *Service) CreateArticle(ctx context.Context, a NewsArticle) error {
	a.ID = s.idGen()
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now
	if a.Status == "" {
		a.Status = ArticleStatusDraft
	}
	return s.prStore.CreateArticle(ctx, a)
}

func (s *Service) UpdateArticle(ctx context.Context, a NewsArticle) error {
	a.UpdatedAt = time.Now().UTC()
	return s.prStore.UpdateArticle(ctx, a)
}

func (s *Service) DeleteArticle(ctx context.Context, id string) error {
	return s.prStore.DeleteArticle(ctx, id)
}

func (s *Service) PublishArticle(ctx context.Context, id string) error {
	a, err := s.prStore.GetArticle(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	a.Status = ArticleStatusPublished
	a.PublishedAt = &now
	a.UpdatedAt = now
	return s.prStore.UpdateArticle(ctx, *a)
}

// ── International Service ────────────────────────────────────

func (s *Service) ListPartners(ctx context.Context) ([]InternationalPartner, error) {
	if s.intlStore == nil {
		return nil, nil
	}
	return s.intlStore.ListPartners(ctx)
}

func (s *Service) GetPartner(ctx context.Context, id string) (*InternationalPartner, error) {
	return s.intlStore.GetPartner(ctx, id)
}

func (s *Service) CreatePartner(ctx context.Context, p InternationalPartner) error {
	p.ID = s.idGen()
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	return s.intlStore.CreatePartner(ctx, p)
}

func (s *Service) UpdatePartner(ctx context.Context, p InternationalPartner) error {
	p.UpdatedAt = time.Now().UTC()
	return s.intlStore.UpdatePartner(ctx, p)
}

func (s *Service) DeletePartner(ctx context.Context, id string) error {
	return s.intlStore.DeletePartner(ctx, id)
}

func (s *Service) ListIntlEvents(ctx context.Context) ([]InternationalEvent, error) {
	if s.intlStore == nil {
		return nil, nil
	}
	return s.intlStore.ListEvents(ctx)
}

func (s *Service) GetIntlEvent(ctx context.Context, id string) (*InternationalEvent, error) {
	return s.intlStore.GetEvent(ctx, id)
}

func (s *Service) CreateIntlEvent(ctx context.Context, e InternationalEvent) error {
	e.ID = s.idGen()
	now := time.Now().UTC()
	e.CreatedAt = now
	e.UpdatedAt = now
	return s.intlStore.CreateEvent(ctx, e)
}

func (s *Service) UpdateIntlEvent(ctx context.Context, e InternationalEvent) error {
	e.UpdatedAt = time.Now().UTC()
	return s.intlStore.UpdateEvent(ctx, e)
}

func (s *Service) DeleteIntlEvent(ctx context.Context, id string) error {
	return s.intlStore.DeleteEvent(ctx, id)
}

// ── Workflow Service ─────────────────────────────────────────

func (s *Service) ListWorkflows(ctx context.Context) ([]WorkflowDefinition, error) {
	if s.workflowStore == nil {
		return nil, nil
	}
	return s.workflowStore.ListWorkflows(ctx)
}

func (s *Service) GetWorkflow(ctx context.Context, id string) (*WorkflowDefinition, error) {
	return s.workflowStore.GetWorkflow(ctx, id)
}

func (s *Service) CreateWorkflow(ctx context.Context, w WorkflowDefinition) error {
	w.ID = s.idGen()
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	return s.workflowStore.CreateWorkflow(ctx, w)
}

func (s *Service) UpdateWorkflow(ctx context.Context, w WorkflowDefinition) error {
	w.UpdatedAt = time.Now().UTC()
	return s.workflowStore.UpdateWorkflow(ctx, w)
}

func (s *Service) DeleteWorkflow(ctx context.Context, id string) error {
	return s.workflowStore.DeleteWorkflow(ctx, id)
}

func (s *Service) ToggleWorkflow(ctx context.Context, id string, active bool) error {
	w, err := s.workflowStore.GetWorkflow(ctx, id)
	if err != nil {
		return err
	}
	w.IsActive = active
	w.UpdatedAt = time.Now().UTC()
	return s.workflowStore.UpdateWorkflow(ctx, *w)
}
