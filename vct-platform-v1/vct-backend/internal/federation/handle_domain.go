package federation

import (
	"encoding/json"
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/certification"
	"vct-platform/backend/internal/domain/discipline"
	"vct-platform/backend/internal/domain/document"
	"vct-platform/backend/internal/shared/auth"
	"vct-platform/backend/internal/shared/httputil"
)

// ── Document Handlers ────────────────────────────────────────

func (m *Module) handleDocumentRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireFederationRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý văn bản")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/documents")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	id := parts[0]
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}

	switch {
	case id == "":
		if r.Method == http.MethodGet {
			m.handleListDocuments(w, r)
		} else if r.Method == http.MethodPost {
			m.handleCreateDocument(w, r, p)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
	case action == "":
		if r.Method == http.MethodGet {
			m.handleGetDocument(w, r, id)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
	case action == "submit":
		m.handleDocumentAction(w, r, id, "submit", p)
	case action == "approve":
		m.handleDocumentAction(w, r, id, "approve", p)
	case action == "publish":
		m.handleDocumentAction(w, r, id, "publish", p)
	case action == "revoke":
		m.handleDocumentAction(w, r, id, "revoke", p)
	default:
		httputil.Error(w, http.StatusNotFound, "DOCUMENT_404", "Không tìm thấy tài nguyên")
	}
}

func (m *Module) handleListDocuments(w http.ResponseWriter, r *http.Request) {
	docs, err := m.document.ListDocuments(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"documents": docs, "total": len(docs)})
}

func (m *Module) handleGetDocument(w http.ResponseWriter, r *http.Request, id string) {
	doc, err := m.document.GetDocument(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "DOCUMENT_404", "Không tìm thấy văn bản")
		return
	}
	httputil.Success(w, http.StatusOK, doc)
}

func (m *Module) handleCreateDocument(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !requireFederationWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo văn bản")
		return
	}
	var doc document.OfficialDocument
	if err := httputil.DecodeJSON(r, &doc); err != nil {
		httputil.Error(w, http.StatusBadRequest, "DOCUMENT_400", err.Error())
		return
	}
	doc.IssuedBy = p.User.ID
	created, err := m.document.CreateDraft(r.Context(), doc)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "DOCUMENT_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleDocumentAction(w http.ResponseWriter, r *http.Request, id, action string, p auth.Principal) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	var err error
	switch action {
	case "submit":
		err = m.document.SubmitForApproval(r.Context(), id)
	case "approve":
		err = m.document.Approve(r.Context(), id, p.User.ID, p.User.DisplayName)
	case "publish":
		err = m.document.Publish(r.Context(), id)
	case "revoke":
		var body struct {
			Reason string `json:"reason"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		err = m.document.Revoke(r.Context(), id, body.Reason)
	}

	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "DOCUMENT_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"status": action + "_done"})
}

// ── Discipline Handlers ──────────────────────────────────────

func (m *Module) handleDisciplineRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireFederationRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý kỷ luật")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/discipline/cases")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	id := parts[0]
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}

	switch {
	case id == "":
		if r.Method == http.MethodGet {
			m.handleListCases(w, r)
		} else if r.Method == http.MethodPost {
			m.handleCreateCase(w, r, p)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
	case action == "":
		if r.Method == http.MethodGet {
			m.handleGetCase(w, r, id)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
	case action == "dismiss":
		m.handleCaseDismiss(w, r, id, p)
	default:
		httputil.Error(w, http.StatusNotFound, "DISCIPLINE_404", "Không tìm thấy hồ sơ kỷ luật")
	}
}

func (m *Module) handleListCases(w http.ResponseWriter, r *http.Request) {
	list, err := m.discipline.ListCases(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"cases": list, "total": len(list)})
}

func (m *Module) handleGetCase(w http.ResponseWriter, r *http.Request, id string) {
	dc, err := m.discipline.GetCase(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "DISCIPLINE_404", "Không tìm thấy hồ sơ kỷ luật")
		return
	}
	httputil.Success(w, http.StatusOK, dc)
}

func (m *Module) handleCreateCase(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !requireFederationWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền báo cáo vi phạm")
		return
	}
	var dc discipline.DisciplineCase
	if err := httputil.DecodeJSON(r, &dc); err != nil {
		httputil.Error(w, http.StatusBadRequest, "DISCIPLINE_400", err.Error())
		return
	}
	dc.ReportedBy = p.User.ID
	created, err := m.discipline.ReportViolation(r.Context(), dc)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "DISCIPLINE_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleCaseDismiss(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	if err := m.discipline.DismissCase(r.Context(), id, p.User.ID); err != nil {
		httputil.Error(w, http.StatusBadRequest, "DISCIPLINE_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"status": "dismissed"})
}

// ── Certification Handlers ───────────────────────────────────

func (m *Module) handleCertificationRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireFederationRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý chứng chỉ")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/certifications")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	id := parts[0]
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}

	switch {
	case id == "":
		if r.Method == http.MethodGet {
			m.handleListCertificates(w, r)
		} else if r.Method == http.MethodPost {
			m.handleIssueCertificate(w, r, p)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
	case action == "":
		if r.Method == http.MethodGet {
			m.handleGetCertificate(w, r, id)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
	default:
		httputil.Error(w, http.StatusNotFound, "CERTIFICATION_404", "Không tìm thấy chứng chỉ")
	}
}

func (m *Module) handleListCertificates(w http.ResponseWriter, r *http.Request) {
	list, err := m.certification.ListByHolder(r.Context(), "", "")
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"certifications": list, "total": len(list)})
}

func (m *Module) handleGetCertificate(w http.ResponseWriter, r *http.Request, id string) {
	cert, err := m.certification.GetCertificate(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "CERTIFICATION_404", "Không tìm thấy chứng chỉ")
		return
	}
	httputil.Success(w, http.StatusOK, cert)
}

func (m *Module) handleIssueCertificate(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !requireFederationWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền cấp chứng chỉ")
		return
	}
	var cert certification.Certificate
	if err := httputil.DecodeJSON(r, &cert); err != nil {
		httputil.Error(w, http.StatusBadRequest, "CERTIFICATION_400", err.Error())
		return
	}
	cert.IssuedBy = p.User.ID
	created, err := m.certification.Issue(r.Context(), cert)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "CERTIFICATION_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleCertVerifyPublic(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/api/v1/certifications/verify/")
	if code == "" {
		httputil.Error(w, http.StatusBadRequest, "CERT_400", "Mã xác thực là bắt buộc")
		return
	}
	cert, err := m.certification.Verify(r.Context(), code)
	if err != nil {
		httputil.Success(w, http.StatusOK, map[string]any{"found": false, "code": code})
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"found": true, "cert": cert})
}
