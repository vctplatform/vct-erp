package httpapi

import (
	"net/http"
	"testing"
)

type nestedIDResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type directIDResponse struct {
	ID string `json:"id"`
}

func TestClubHandlers_EnforceWriteRolesOnMutations(t *testing.T) {
	server := newTestServer()
	handler := server.Handler()

	clubLeaderToken := loginAccessToken(t, handler, "club-leader", "Club@123", "club_leader")
	coachToken := loginAccessToken(t, handler, "coach", "Coach@123", "coach")

	createMember := requestJSON(t, handler, http.MethodPost, "/api/v1/club/members?club_id=CLB-001", map[string]any{
		"full_name":     "Vo Sinh Test",
		"gender":        "nam",
		"belt_rank":     "lam_dai",
		"member_type":   "student",
		"join_date":     "2026-03-14",
		"date_of_birth": "2012-06-01",
	}, clubLeaderToken)
	if createMember.Code != http.StatusCreated {
		t.Fatalf("expected club leader create member 201, got %d (%s)", createMember.Code, createMember.Body.String())
	}
	member := decodeBody[nestedIDResponse](t, createMember)
	if member.Data.ID == "" {
		t.Fatalf("expected created member id, got %s", createMember.Body.String())
	}

	getMember := requestJSON(t, handler, http.MethodGet, "/api/v1/club/members/"+member.Data.ID, nil, coachToken)
	if getMember.Code != http.StatusOK {
		t.Fatalf("expected coach read member 200, got %d (%s)", getMember.Code, getMember.Body.String())
	}

	updateByCoach := requestJSON(t, handler, http.MethodPut, "/api/v1/club/members/"+member.Data.ID, map[string]any{
		"belt_rank": "hoang_dai",
	}, coachToken)
	if updateByCoach.Code != http.StatusForbidden {
		t.Fatalf("expected coach update member 403, got %d (%s)", updateByCoach.Code, updateByCoach.Body.String())
	}

	approveByCoach := requestJSON(t, handler, http.MethodPost, "/api/v1/club/members/"+member.Data.ID+"/approve", nil, coachToken)
	if approveByCoach.Code != http.StatusForbidden {
		t.Fatalf("expected coach approve member 403, got %d (%s)", approveByCoach.Code, approveByCoach.Body.String())
	}

	approveByLeader := requestJSON(t, handler, http.MethodPost, "/api/v1/club/members/"+member.Data.ID+"/approve", nil, clubLeaderToken)
	if approveByLeader.Code != http.StatusOK {
		t.Fatalf("expected club leader approve member 200, got %d (%s)", approveByLeader.Code, approveByLeader.Body.String())
	}

	createFinanceByCoach := requestJSON(t, handler, http.MethodPost, "/api/v1/club/finance?club_id=CLB-001", map[string]any{
		"type":        "income",
		"category":    "hoi_phi",
		"amount":      500000,
		"description": "Thu hoi phi",
		"date":        "2026-03-14",
		"recorded_by": "coach",
	}, coachToken)
	if createFinanceByCoach.Code != http.StatusForbidden {
		t.Fatalf("expected coach create finance 403, got %d (%s)", createFinanceByCoach.Code, createFinanceByCoach.Body.String())
	}

	createFinanceByLeader := requestJSON(t, handler, http.MethodPost, "/api/v1/club/finance?club_id=CLB-001", map[string]any{
		"type":        "income",
		"category":    "hoi_phi",
		"amount":      500000,
		"description": "Thu hoi phi",
		"date":        "2026-03-14",
		"recorded_by": "club-leader",
	}, clubLeaderToken)
	if createFinanceByLeader.Code != http.StatusCreated {
		t.Fatalf("expected club leader create finance 201, got %d (%s)", createFinanceByLeader.Code, createFinanceByLeader.Body.String())
	}
}

func TestProvincialHandlers_RequireWriteRolesForMutations(t *testing.T) {
	server := newTestServer()
	handler := server.Handler()

	provincialToken := loginAccessToken(t, handler, "provincial", "Prov@123", "provincial_admin")
	athleteToken := loginAccessToken(t, handler, "athlete", "Athlete@123", "athlete")

	listByAthlete := requestJSON(t, handler, http.MethodGet, "/api/v1/provincial/referees", nil, athleteToken)
	if listByAthlete.Code != http.StatusForbidden {
		t.Fatalf("expected athlete list referees 403, got %d (%s)", listByAthlete.Code, listByAthlete.Body.String())
	}

	createReferee := requestJSON(t, handler, http.MethodPost, "/api/v1/provincial/referees", map[string]any{
		"full_name":        "Trong Tai Test",
		"gender":           "nam",
		"referee_rank":     "cap_2",
		"expertise":        "doi_khang",
		"experience_years": 5,
	}, provincialToken)
	if createReferee.Code != http.StatusCreated {
		t.Fatalf("expected provincial admin create referee 201, got %d (%s)", createReferee.Code, createReferee.Body.String())
	}
	referee := decodeBody[directIDResponse](t, createReferee)
	if referee.ID == "" {
		t.Fatalf("expected created referee id, got %s", createReferee.Body.String())
	}

	createCertByAthlete := requestJSON(t, handler, http.MethodPost, "/api/v1/provincial/referees/"+referee.ID+"/certificates", map[string]any{
		"name":       "Chung chi cap II",
		"issuer":     "Lien doan",
		"cert_type":  "referee_license",
		"issue_date": "2026-03-14",
		"status":     "valid",
	}, athleteToken)
	if createCertByAthlete.Code != http.StatusForbidden {
		t.Fatalf("expected athlete create referee certificate 403, got %d (%s)", createCertByAthlete.Code, createCertByAthlete.Body.String())
	}

	createCertByProvincial := requestJSON(t, handler, http.MethodPost, "/api/v1/provincial/referees/"+referee.ID+"/certificates", map[string]any{
		"name":       "Chung chi cap II",
		"issuer":     "Lien doan",
		"cert_type":  "referee_license",
		"issue_date": "2026-03-14",
		"status":     "valid",
	}, provincialToken)
	if createCertByProvincial.Code != http.StatusCreated {
		t.Fatalf("expected provincial admin create referee certificate 201, got %d (%s)", createCertByProvincial.Code, createCertByProvincial.Body.String())
	}

	approveByAthlete := requestJSON(t, handler, http.MethodPost, "/api/v1/provincial/referees/"+referee.ID+"/approve", nil, athleteToken)
	if approveByAthlete.Code != http.StatusForbidden {
		t.Fatalf("expected athlete approve referee 403, got %d (%s)", approveByAthlete.Code, approveByAthlete.Body.String())
	}

	approveByProvincial := requestJSON(t, handler, http.MethodPost, "/api/v1/provincial/referees/"+referee.ID+"/approve", nil, provincialToken)
	if approveByProvincial.Code != http.StatusOK {
		t.Fatalf("expected provincial admin approve referee 200, got %d (%s)", approveByProvincial.Code, approveByProvincial.Body.String())
	}
}

func TestFederationHandlers_ProtectDocumentAndDisciplineMutations(t *testing.T) {
	server := newTestServer()
	handler := server.Handler()

	adminToken := loginAccessToken(t, handler, "admin", "Admin@123", "admin")
	athleteToken := loginAccessToken(t, handler, "athlete", "Athlete@123", "athlete")

	createDocument := requestJSON(t, handler, http.MethodPost, "/api/v1/documents", map[string]any{
		"number": "QD-2026-001",
		"title":  "Quyet dinh test",
		"type":   "decision",
	}, adminToken)
	if createDocument.Code != http.StatusCreated {
		t.Fatalf("expected admin create document 201, got %d (%s)", createDocument.Code, createDocument.Body.String())
	}
	document := decodeBody[directIDResponse](t, createDocument)
	if document.ID == "" {
		t.Fatalf("expected created document id, got %s", createDocument.Body.String())
	}

	readDocumentByAthlete := requestJSON(t, handler, http.MethodGet, "/api/v1/documents/"+document.ID, nil, athleteToken)
	if readDocumentByAthlete.Code != http.StatusForbidden {
		t.Fatalf("expected athlete read document 403, got %d (%s)", readDocumentByAthlete.Code, readDocumentByAthlete.Body.String())
	}

	submitByAthlete := requestJSON(t, handler, http.MethodPost, "/api/v1/documents/"+document.ID+"/submit", nil, athleteToken)
	if submitByAthlete.Code != http.StatusForbidden {
		t.Fatalf("expected athlete submit document 403, got %d (%s)", submitByAthlete.Code, submitByAthlete.Body.String())
	}

	submitByAdmin := requestJSON(t, handler, http.MethodPost, "/api/v1/documents/"+document.ID+"/submit", nil, adminToken)
	if submitByAdmin.Code != http.StatusOK {
		t.Fatalf("expected admin submit document 200, got %d (%s)", submitByAdmin.Code, submitByAdmin.Body.String())
	}

	publishByAthlete := requestJSON(t, handler, http.MethodPost, "/api/v1/documents/"+document.ID+"/publish", nil, athleteToken)
	if publishByAthlete.Code != http.StatusForbidden {
		t.Fatalf("expected athlete publish document 403, got %d (%s)", publishByAthlete.Code, publishByAthlete.Body.String())
	}

	approveByAdmin := requestJSON(t, handler, http.MethodPost, "/api/v1/documents/"+document.ID+"/approve", nil, adminToken)
	if approveByAdmin.Code != http.StatusOK {
		t.Fatalf("expected admin approve document 200, got %d (%s)", approveByAdmin.Code, approveByAdmin.Body.String())
	}

	publishByAdmin := requestJSON(t, handler, http.MethodPost, "/api/v1/documents/"+document.ID+"/publish", nil, adminToken)
	if publishByAdmin.Code != http.StatusOK {
		t.Fatalf("expected admin publish document 200, got %d (%s)", publishByAdmin.Code, publishByAdmin.Body.String())
	}

	revokeByAthlete := requestJSON(t, handler, http.MethodPost, "/api/v1/documents/"+document.ID+"/revoke", map[string]any{
		"reason": "khong hop le",
	}, athleteToken)
	if revokeByAthlete.Code != http.StatusForbidden {
		t.Fatalf("expected athlete revoke document 403, got %d (%s)", revokeByAthlete.Code, revokeByAthlete.Body.String())
	}

	revokeByAdmin := requestJSON(t, handler, http.MethodPost, "/api/v1/documents/"+document.ID+"/revoke", map[string]any{
		"reason": "het hieu luc",
	}, adminToken)
	if revokeByAdmin.Code != http.StatusOK {
		t.Fatalf("expected admin revoke document 200, got %d (%s)", revokeByAdmin.Code, revokeByAdmin.Body.String())
	}

	createCase := requestJSON(t, handler, http.MethodPost, "/api/v1/discipline/cases", map[string]any{
		"title":          "Vu viec test",
		"violation_type": "other",
		"subject_type":   "athlete",
		"subject_id":     "athlete-001",
		"subject_name":   "Vo Sinh Test",
	}, adminToken)
	if createCase.Code != http.StatusCreated {
		t.Fatalf("expected admin create discipline case 201, got %d (%s)", createCase.Code, createCase.Body.String())
	}
	disciplineCase := decodeBody[directIDResponse](t, createCase)
	if disciplineCase.ID == "" {
		t.Fatalf("expected created discipline case id, got %s", createCase.Body.String())
	}

	dismissByAthlete := requestJSON(t, handler, http.MethodPost, "/api/v1/discipline/cases/"+disciplineCase.ID+"/dismiss", nil, athleteToken)
	if dismissByAthlete.Code != http.StatusForbidden {
		t.Fatalf("expected athlete dismiss discipline case 403, got %d (%s)", dismissByAthlete.Code, dismissByAthlete.Body.String())
	}

	dismissByAdmin := requestJSON(t, handler, http.MethodPost, "/api/v1/discipline/cases/"+disciplineCase.ID+"/dismiss", nil, adminToken)
	if dismissByAdmin.Code != http.StatusOK {
		t.Fatalf("expected admin dismiss discipline case 200, got %d (%s)", dismissByAdmin.Code, dismissByAdmin.Body.String())
	}
}

func TestFederationHandlers_ProtectCertificationMutations(t *testing.T) {
	server := newTestServer()
	handler := server.Handler()

	adminToken := loginAccessToken(t, handler, "admin", "Admin@123", "admin")
	athleteToken := loginAccessToken(t, handler, "athlete", "Athlete@123", "athlete")

	createCert := requestJSON(t, handler, http.MethodPost, "/api/v1/certifications", map[string]any{
		"type":        "coach_license",
		"holder_type": "person",
		"holder_id":   "holder-001",
		"holder_name": "Huynh Luyen Vien Test",
	}, adminToken)
	if createCert.Code != http.StatusCreated {
		t.Fatalf("expected admin create certification 201, got %d (%s)", createCert.Code, createCert.Body.String())
	}
	cert := decodeBody[directIDResponse](t, createCert)
	if cert.ID == "" {
		t.Fatalf("expected created certification id, got %s", createCert.Body.String())
	}

	readCertByAthlete := requestJSON(t, handler, http.MethodGet, "/api/v1/certifications/"+cert.ID, nil, athleteToken)
	if readCertByAthlete.Code != http.StatusForbidden {
		t.Fatalf("expected athlete read certification 403, got %d (%s)", readCertByAthlete.Code, readCertByAthlete.Body.String())
	}

	renewByAthlete := requestJSON(t, handler, http.MethodPost, "/api/v1/certifications/"+cert.ID+"/renew", map[string]any{
		"valid_until": "2027-12-31",
	}, athleteToken)
	if renewByAthlete.Code != http.StatusForbidden {
		t.Fatalf("expected athlete renew certification 403, got %d (%s)", renewByAthlete.Code, renewByAthlete.Body.String())
	}

	renewByAdmin := requestJSON(t, handler, http.MethodPost, "/api/v1/certifications/"+cert.ID+"/renew", map[string]any{
		"valid_until": "2027-12-31",
	}, adminToken)
	if renewByAdmin.Code != http.StatusOK {
		t.Fatalf("expected admin renew certification 200, got %d (%s)", renewByAdmin.Code, renewByAdmin.Body.String())
	}
}
