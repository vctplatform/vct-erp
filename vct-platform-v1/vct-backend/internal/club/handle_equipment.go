package club

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/club"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleClubEquipment(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}

	switch r.Method {
	case http.MethodGet:
		items, err := m.service.ListEquipment(r.Context(), clubID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, items)

	case http.MethodPost:
		var e club.Equipment
		if err := httputil.DecodeJSON(r, &e); err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		if e.ClubID == "" {
			e.ClubID = clubID
		}
		created, err := m.service.CreateEquipment(r.Context(), e)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)
	}
}

func (m *Module) handleClubEquipmentAction(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/club/equipment/")
	switch r.Method {
	case http.MethodGet:
		e, err := m.service.GetEquipment(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "CLUB_404", "Equipment not found")
			return
		}
		httputil.Success(w, http.StatusOK, e)
	case http.MethodPut, http.MethodPatch:
		var patch map[string]any
		if err := httputil.DecodeJSON(r, &patch); err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		if err := m.service.UpdateEquipment(r.Context(), id, patch); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"message": "updated"})
	case http.MethodDelete:
		if err := m.service.DeleteEquipment(r.Context(), id); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusNoContent, nil)
	}
}

func (m *Module) handleClubEquipmentSummary(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}
	summary, err := m.service.GetEquipmentSummary(r.Context(), clubID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, summary)
}

func (m *Module) handleClubEquipmentExport(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}
	csv, err := m.service.ExportEquipmentCSV(r.Context(), clubID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=equipment.csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(csv))
}
