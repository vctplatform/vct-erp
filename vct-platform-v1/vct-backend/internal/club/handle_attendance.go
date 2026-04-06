package club

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/club"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleClubAttendance(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}

	switch r.Method {
	case http.MethodGet:
		classID := r.URL.Query().Get("class_id")
		date := r.URL.Query().Get("date")

		var records []club.Attendance
		var err error
		if date != "" {
			records, err = m.service.ListAttendanceByDate(r.Context(), clubID, date)
		} else if classID != "" {
			records, err = m.service.ListAttendanceByClass(r.Context(), clubID, classID)
		} else {
			records, err = m.service.ListAttendance(r.Context(), clubID)
		}
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, records)

	case http.MethodPost:
		var a club.Attendance
		if err := httputil.DecodeJSON(r, &a); err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		if a.ClubID == "" {
			a.ClubID = clubID
		}
		created, err := m.service.RecordAttendance(r.Context(), a)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)
	}
}

func (m *Module) handleClubAttendanceAction(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/club/attendance/")
	if r.Method == http.MethodDelete {
		if err := m.service.DeleteAttendance(r.Context(), id); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"message": "deleted"})
		return
	}
	httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
}

func (m *Module) handleClubAttendanceSummary(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}
	summary, err := m.service.GetAttendanceSummary(r.Context(), clubID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, summary)
}

func (m *Module) handleClubAttendanceBulk(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Records []club.Attendance `json:"records"`
	}
	if err := httputil.DecodeJSON(r, &body); err != nil {
		httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
		return
	}
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}
	for i := range body.Records {
		if body.Records[i].ClubID == "" {
			body.Records[i].ClubID = clubID
		}
	}
	created, err := m.service.BulkRecordAttendance(r.Context(), body.Records)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleClubAttendanceExport(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}
	csv, err := m.service.ExportAttendanceCSV(r.Context(), clubID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=attendance.csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(csv))
}
