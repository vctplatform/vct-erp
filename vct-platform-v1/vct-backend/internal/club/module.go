package club

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/club"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Club module.
type Module struct {
	service     *club.Service
	logger      *slog.Logger
	broadcaster httputil.EventBroadcaster
}

// Deps holds the dependencies for the Club module.
type Deps struct {
	Service     *club.Service
	Logger      *slog.Logger
	Broadcaster httputil.EventBroadcaster
}

// New creates a new Club module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service:     deps.Service,
		logger:      deps.Logger.With(slog.String("module", "club")),
		broadcaster: deps.Broadcaster,
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers all club-v2 routes.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/club/dashboard", m.handleClubDashboard)
	mux.HandleFunc("/api/v1/club/attendance", m.handleClubAttendance)
	mux.HandleFunc("/api/v1/club/attendance/", m.handleClubAttendanceAction)
	mux.HandleFunc("/api/v1/club/attendance/bulk", m.handleClubAttendanceBulk)
	mux.HandleFunc("/api/v1/club/attendance/summary", m.handleClubAttendanceSummary)
	mux.HandleFunc("/api/v1/club/attendance/export", m.handleClubAttendanceExport)

	mux.HandleFunc("/api/v1/club/equipment", m.handleClubEquipment)
	mux.HandleFunc("/api/v1/club/equipment/", m.handleClubEquipmentAction)
	mux.HandleFunc("/api/v1/club/equipment/summary", m.handleClubEquipmentSummary)
	mux.HandleFunc("/api/v1/club/equipment/export", m.handleClubEquipmentExport)

	mux.HandleFunc("/api/v1/club/facilities", m.handleClubFacilities)
	mux.HandleFunc("/api/v1/club/facilities/", m.handleClubFacilityAction)
	mux.HandleFunc("/api/v1/club/facilities/summary", m.handleClubFacilitySummary)
	mux.HandleFunc("/api/v1/club/facilities/export", m.handleClubFacilitiesExport)

	m.logger.Info("club module routes registered")
}
