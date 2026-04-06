package club

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/club"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleClubFacilities(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}

	switch r.Method {
	case http.MethodGet:
		items, err := m.service.ListFacilities(r.Context(), clubID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, items)

	case http.MethodPost:
		var f club.Facility
		if err := httputil.DecodeJSON(r, &f); err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		if f.ClubID == "" {
			f.ClubID = clubID
		}
		created, err := m.service.CreateFacility(r.Context(), f)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)
	}
}

func (m *Module) handleClubFacilityAction(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/club/facilities/")
	switch r.Method {
	case http.MethodGet:
		f, err := m.service.GetFacility(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "CLUB_404", "Facility not found")
			return
		}
		httputil.Success(w, http.StatusOK, f)
	case http.MethodPut, http.MethodPatch:
		var patch map[string]any
		if err := httputil.DecodeJSON(r, &patch); err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		if err := m.service.UpdateFacility(r.Context(), id, patch); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"message": "updated"})
	case http.MethodDelete:
		if err := m.service.DeleteFacility(r.Context(), id); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusNoContent, nil)
	}
}

func (m *Module) handleClubFacilitySummary(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}
	summary, err := m.service.GetFacilitySummary(r.Context(), clubID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, summary)
}

func (m *Module) handleClubFacilitiesExport(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}
	csv, err := m.service.ExportFacilitiesCSV(r.Context(), clubID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=facilities.csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(csv))
}
