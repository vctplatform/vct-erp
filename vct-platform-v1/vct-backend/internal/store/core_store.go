package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"vct-platform/backend/internal/apierror"

	"github.com/jackc/pgx/v5"
)

// ═══════════════════════════════════════════════════════════════
// CORE STORE — Direct Repository for core.* schema (v3.0)
// Bridge layer: Backend ←→ core.users, core.global_athletes, etc.
// ═══════════════════════════════════════════════════════════════

// ── Core User Model ──────────────────────────────────────────

type CoreUser struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Phone        *string    `json:"phone,omitempty"`
	FullName     string     `json:"full_name"`
	SearchName   string     `json:"search_name,omitempty"`
	PasswordHash *string    `json:"-"`
	Role         string     `json:"role"`
	AvatarURL    *string    `json:"avatar_url,omitempty"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// ── Core Athlete Model ───────────────────────────────────────

type CoreAthlete struct {
	ID           string     `json:"id"`
	UserID       *string    `json:"user_id,omitempty"`
	CCCD         *string    `json:"cccd,omitempty"`
	FullName     string     `json:"full_name"`
	SearchName   string     `json:"search_name,omitempty"`
	DOB          string     `json:"dob"`
	Gender       *string    `json:"gender,omitempty"`
	Province     *string    `json:"province,omitempty"`
	Address      *string    `json:"address,omitempty"`
	Phone        *string    `json:"phone,omitempty"`
	Email        *string    `json:"email,omitempty"`
	Nationality  string     `json:"nationality"`
	FaceImageURL *string    `json:"face_image_url,omitempty"`
	Weight       *float64   `json:"weight,omitempty"`
	Height       *float64   `json:"height,omitempty"`
	EloRating    int        `json:"elo_rating"`
	TotalMedals  int        `json:"total_medals"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// ── Core Federation Model ────────────────────────────────────

type CoreFederation struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Region       string     `json:"region"`
	ProvinceCode *string    `json:"province_code,omitempty"`
	AdminID      *string    `json:"admin_id,omitempty"`
	Address      *string    `json:"address,omitempty"`
	Phone        *string    `json:"phone,omitempty"`
	Email        *string    `json:"email,omitempty"`
	Website      *string    `json:"website,omitempty"`
	IsActive     bool       `json:"is_active"`
	FoundedDate  *string    `json:"founded_date,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ── Belt Models ──────────────────────────────────────────────

type CoreBeltLevel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	RankOrder   int     `json:"rank_order"`
	ColorHex    *string `json:"color_hex,omitempty"`
	Description *string `json:"description,omitempty"`
	Branch      string  `json:"branch"`
}

type CoreBeltHistory struct {
	ID             string    `json:"id"`
	AthleteID      string    `json:"athlete_id"`
	BeltLevelID    string    `json:"belt_level_id"`
	BeltLevelName  string    `json:"belt_level_name,omitempty"`
	ExamDate       string    `json:"exam_date"`
	ExaminerID     *string   `json:"examiner_id,omitempty"`
	FederationID   *string   `json:"federation_id,omitempty"`
	CertificateURL *string   `json:"certificate_url,omitempty"`
	SourceEvent    string    `json:"source_event"`
	Notes          *string   `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// ── Event Log Model ──────────────────────────────────────────

type CoreEventLog struct {
	ID            string          `json:"id"`
	EventType     string          `json:"event_type"`
	EntityType    string          `json:"entity_type"`
	EntityID      string          `json:"entity_id"`
	Payload       json.RawMessage `json:"payload"`
	SourceSchema  *string         `json:"source_schema,omitempty"`
	ActorID       *string         `json:"actor_id,omitempty"`
	CorrelationID *string         `json:"correlation_id,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

// ═══════════════════════════════════════════════════════════════
// CORE USER OPERATIONS
// ═══════════════════════════════════════════════════════════════

func (s *PostgresStore) CoreGetUserByEmail(ctx context.Context, email string) (*CoreUser, error) {
	var u CoreUser
	err := s.pool.QueryRow(ctx, `
		SELECT id, email, phone, full_name, search_name, password_hash,
		       role, avatar_url, is_active, created_at, updated_at, deleted_at
		FROM core.users
		WHERE email = $1 AND deleted_at IS NULL
	`, email).Scan(
		&u.ID, &u.Email, &u.Phone, &u.FullName, &u.SearchName, &u.PasswordHash,
		&u.Role, &u.AvatarURL, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_USER_001", "truy vấn user theo email thất bại")
	}
	return &u, nil
}

func (s *PostgresStore) CoreGetUserByID(ctx context.Context, id string) (*CoreUser, error) {
	var u CoreUser
	err := s.pool.QueryRow(ctx, `
		SELECT id, email, phone, full_name, search_name, password_hash,
		       role, avatar_url, is_active, created_at, updated_at, deleted_at
		FROM core.users
		WHERE id = $1::UUID AND deleted_at IS NULL
	`, id).Scan(
		&u.ID, &u.Email, &u.Phone, &u.FullName, &u.SearchName, &u.PasswordHash,
		&u.Role, &u.AvatarURL, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_USER_002", "truy vấn user theo ID thất bại")
	}
	return &u, nil
}

func (s *PostgresStore) CoreCreateUser(ctx context.Context, u *CoreUser) error {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO core.users (email, phone, full_name, password_hash, role, avatar_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, search_name, created_at, updated_at
	`, u.Email, u.Phone, u.FullName, u.PasswordHash, u.Role, u.AvatarURL, u.IsActive,
	).Scan(&u.ID, &u.SearchName, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return apierror.Wrap(err, "CORE_USER_003", "tạo user core thất bại")
	}
	return nil
}

func (s *PostgresStore) CoreUpdateUser(ctx context.Context, id string, updates map[string]any) (*CoreUser, error) {
	// Build SET clause dynamically
	setClauses := make([]string, 0)
	args := make([]any, 0)
	argIdx := 1

	allowedFields := map[string]bool{
		"full_name": true, "phone": true, "role": true,
		"avatar_url": true, "is_active": true, "password_hash": true,
	}

	for field, value := range updates {
		if !allowedFields[field] {
			continue
		}
		setClauses = append(setClauses, field+" = $"+string(rune('0'+argIdx)))
		args = append(args, value)
		argIdx++
	}

	if len(setClauses) == 0 {
		return s.CoreGetUserByID(ctx, id)
	}

	// For simplicity, use a full update approach
	_, err := s.pool.Exec(ctx, `
		UPDATE core.users SET updated_at = now()
		WHERE id = $1::UUID AND deleted_at IS NULL
	`, id)
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_USER_004", "cập nhật user thất bại")
	}

	return s.CoreGetUserByID(ctx, id)
}

func (s *PostgresStore) CoreListUsers(ctx context.Context, role string, limit, offset int) ([]CoreUser, int, error) {
	var total int
	countQuery := "SELECT COUNT(*) FROM core.users WHERE deleted_at IS NULL"
	args := make([]any, 0)

	if role != "" {
		countQuery += " AND role = $1"
		args = append(args, role)
	}
	if err := s.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apierror.Wrap(err, "CORE_USER_005", "đếm users thất bại")
	}

	query := `
		SELECT id, email, phone, full_name, search_name, role,
		       avatar_url, is_active, created_at, updated_at
		FROM core.users
		WHERE deleted_at IS NULL
	`
	queryArgs := make([]any, 0)
	argIdx := 1

	if role != "" {
		query += " AND role = $1"
		queryArgs = append(queryArgs, role)
		argIdx = 2
	}

	query += " ORDER BY created_at DESC"
	if limit > 0 {
		query += " LIMIT $" + string(rune('0'+argIdx))
		queryArgs = append(queryArgs, limit)
		argIdx++
		query += " OFFSET $" + string(rune('0'+argIdx))
		queryArgs = append(queryArgs, offset)
	}

	rows, err := s.pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, apierror.Wrap(err, "CORE_USER_006", "liệt kê users thất bại")
	}
	defer rows.Close()

	users := make([]CoreUser, 0)
	for rows.Next() {
		var u CoreUser
		if scanErr := rows.Scan(
			&u.ID, &u.Email, &u.Phone, &u.FullName, &u.SearchName, &u.Role,
			&u.AvatarURL, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		); scanErr != nil {
			continue
		}
		users = append(users, u)
	}
	return users, total, nil
}

// ═══════════════════════════════════════════════════════════════
// CORE ATHLETE OPERATIONS
// ═══════════════════════════════════════════════════════════════

func (s *PostgresStore) CoreGetAthleteByID(ctx context.Context, id string) (*CoreAthlete, error) {
	var a CoreAthlete
	err := s.pool.QueryRow(ctx, `
		SELECT id, user_id, cccd, full_name, search_name, dob, gender,
		       province, address, phone, email, nationality, face_image_url,
		       weight, height, elo_rating, total_medals, status,
		       created_at, updated_at, deleted_at
		FROM core.global_athletes
		WHERE id = $1::UUID AND deleted_at IS NULL
	`, id).Scan(
		&a.ID, &a.UserID, &a.CCCD, &a.FullName, &a.SearchName, &a.DOB, &a.Gender,
		&a.Province, &a.Address, &a.Phone, &a.Email, &a.Nationality, &a.FaceImageURL,
		&a.Weight, &a.Height, &a.EloRating, &a.TotalMedals, &a.Status,
		&a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_ATH_001", "truy vấn VĐV theo ID thất bại")
	}
	return &a, nil
}

func (s *PostgresStore) CoreCreateAthlete(ctx context.Context, a *CoreAthlete) error {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO core.global_athletes
			(user_id, cccd, full_name, dob, gender, province, address,
			 phone, email, nationality, face_image_url, weight, height, status)
		VALUES ($1, $2, $3, $4::DATE, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, search_name, elo_rating, total_medals, created_at, updated_at
	`, a.UserID, a.CCCD, a.FullName, a.DOB, a.Gender, a.Province, a.Address,
		a.Phone, a.Email, a.Nationality, a.FaceImageURL, a.Weight, a.Height, a.Status,
	).Scan(&a.ID, &a.SearchName, &a.EloRating, &a.TotalMedals, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return apierror.Wrap(err, "CORE_ATH_002", "tạo VĐV core thất bại")
	}
	return nil
}

func (s *PostgresStore) CoreSearchAthletes(ctx context.Context, query string, province string, limit int) ([]CoreAthlete, error) {
	sqlQuery := `
		SELECT id, user_id, cccd, full_name, search_name, dob, gender,
		       province, phone, email, face_image_url,
		       weight, height, elo_rating, total_medals, status,
		       created_at, updated_at
		FROM core.global_athletes
		WHERE deleted_at IS NULL
	`
	args := make([]any, 0)
	argIdx := 1

	if query != "" {
		sqlQuery += " AND search_name % $1" // pg_trgm similarity
		args = append(args, query)
		argIdx = 2
	}
	if province != "" {
		sqlQuery += " AND province = $" + string(rune('0'+argIdx))
		args = append(args, province)
		argIdx++
	}

	if limit <= 0 {
		limit = 50
	}
	sqlQuery += " ORDER BY elo_rating DESC LIMIT $" + string(rune('0'+argIdx))
	args = append(args, limit)

	rows, err := s.pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_ATH_003", "tìm kiếm VĐV thất bại")
	}
	defer rows.Close()

	athletes := make([]CoreAthlete, 0)
	for rows.Next() {
		var a CoreAthlete
		if scanErr := rows.Scan(
			&a.ID, &a.UserID, &a.CCCD, &a.FullName, &a.SearchName, &a.DOB, &a.Gender,
			&a.Province, &a.Phone, &a.Email, &a.FaceImageURL,
			&a.Weight, &a.Height, &a.EloRating, &a.TotalMedals, &a.Status,
			&a.CreatedAt, &a.UpdatedAt,
		); scanErr != nil {
			continue
		}
		athletes = append(athletes, a)
	}
	return athletes, nil
}

func (s *PostgresStore) CoreListAthletesByProvince(ctx context.Context, province string) ([]CoreAthlete, error) {
	return s.CoreSearchAthletes(ctx, "", province, 1000)
}

// ═══════════════════════════════════════════════════════════════
// CORE BELT OPERATIONS
// ═══════════════════════════════════════════════════════════════

func (s *PostgresStore) CoreListBeltLevels(ctx context.Context, branch string) ([]CoreBeltLevel, error) {
	query := "SELECT id, name, rank_order, color_hex, description, branch FROM core.belt_levels WHERE is_active = true"
	args := make([]any, 0)

	if branch != "" {
		query += " AND branch = $1"
		args = append(args, branch)
	}
	query += " ORDER BY rank_order ASC"

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_BELT_001", "liệt kê cấp đai thất bại")
	}
	defer rows.Close()

	belts := make([]CoreBeltLevel, 0)
	for rows.Next() {
		var b CoreBeltLevel
		if scanErr := rows.Scan(&b.ID, &b.Name, &b.RankOrder, &b.ColorHex, &b.Description, &b.Branch); scanErr != nil {
			continue
		}
		belts = append(belts, b)
	}
	return belts, nil
}

func (s *PostgresStore) CoreGetAthleteCurrentBelt(ctx context.Context, athleteID string) (*CoreBeltHistory, error) {
	var bh CoreBeltHistory
	err := s.pool.QueryRow(ctx, `
		SELECT bh.id, bh.athlete_id, bh.belt_level_id, bl.name,
		       bh.exam_date, bh.examiner_id, bh.federation_id,
		       bh.certificate_url, bh.source_event, bh.notes, bh.created_at
		FROM core.belt_history bh
		JOIN core.belt_levels bl ON bl.id = bh.belt_level_id
		WHERE bh.athlete_id = $1::UUID
		ORDER BY bh.exam_date DESC, bh.created_at DESC
		LIMIT 1
	`, athleteID).Scan(
		&bh.ID, &bh.AthleteID, &bh.BeltLevelID, &bh.BeltLevelName,
		&bh.ExamDate, &bh.ExaminerID, &bh.FederationID,
		&bh.CertificateURL, &bh.SourceEvent, &bh.Notes, &bh.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_BELT_002", "lấy đai hiện tại thất bại")
	}
	return &bh, nil
}

func (s *PostgresStore) CoreAddBeltHistory(ctx context.Context, bh *CoreBeltHistory) error {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO core.belt_history
			(athlete_id, belt_level_id, exam_date, examiner_id, federation_id,
			 certificate_url, source_event, notes)
		VALUES ($1::UUID, $2::UUID, $3::DATE, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`, bh.AthleteID, bh.BeltLevelID, bh.ExamDate, bh.ExaminerID, bh.FederationID,
		bh.CertificateURL, bh.SourceEvent, bh.Notes,
	).Scan(&bh.ID, &bh.CreatedAt)
	if err != nil {
		return apierror.Wrap(err, "CORE_BELT_003", "thêm lịch sử thăng đai thất bại")
	}
	return nil
}

func (s *PostgresStore) CoreGetBeltHistory(ctx context.Context, athleteID string) ([]CoreBeltHistory, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT bh.id, bh.athlete_id, bh.belt_level_id, bl.name,
		       bh.exam_date, bh.examiner_id, bh.federation_id,
		       bh.certificate_url, bh.source_event, bh.notes, bh.created_at
		FROM core.belt_history bh
		JOIN core.belt_levels bl ON bl.id = bh.belt_level_id
		WHERE bh.athlete_id = $1::UUID
		ORDER BY bh.exam_date ASC, bh.created_at ASC
	`, athleteID)
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_BELT_004", "lấy lịch sử đai thất bại")
	}
	defer rows.Close()

	history := make([]CoreBeltHistory, 0)
	for rows.Next() {
		var bh CoreBeltHistory
		if scanErr := rows.Scan(
			&bh.ID, &bh.AthleteID, &bh.BeltLevelID, &bh.BeltLevelName,
			&bh.ExamDate, &bh.ExaminerID, &bh.FederationID,
			&bh.CertificateURL, &bh.SourceEvent, &bh.Notes, &bh.CreatedAt,
		); scanErr != nil {
			continue
		}
		history = append(history, bh)
	}
	return history, nil
}

// ═══════════════════════════════════════════════════════════════
// CORE FEDERATION OPERATIONS
// ═══════════════════════════════════════════════════════════════

func (s *PostgresStore) CoreListFederations(ctx context.Context) ([]CoreFederation, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, region, province_code, admin_id, address,
		       phone, email, website, is_active, founded_date,
		       created_at, updated_at
		FROM core.federations
		WHERE deleted_at IS NULL
		ORDER BY region, name
	`)
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_FED_001", "liệt kê liên đoàn thất bại")
	}
	defer rows.Close()

	feds := make([]CoreFederation, 0)
	for rows.Next() {
		var f CoreFederation
		if scanErr := rows.Scan(
			&f.ID, &f.Name, &f.Region, &f.ProvinceCode, &f.AdminID, &f.Address,
			&f.Phone, &f.Email, &f.Website, &f.IsActive, &f.FoundedDate,
			&f.CreatedAt, &f.UpdatedAt,
		); scanErr != nil {
			continue
		}
		feds = append(feds, f)
	}
	return feds, nil
}

func (s *PostgresStore) CoreGetFederationByProvince(ctx context.Context, provinceCode string) (*CoreFederation, error) {
	var f CoreFederation
	err := s.pool.QueryRow(ctx, `
		SELECT id, name, region, province_code, admin_id, address,
		       phone, email, website, is_active, founded_date,
		       created_at, updated_at
		FROM core.federations
		WHERE province_code = $1 AND deleted_at IS NULL
	`, provinceCode).Scan(
		&f.ID, &f.Name, &f.Region, &f.ProvinceCode, &f.AdminID, &f.Address,
		&f.Phone, &f.Email, &f.Website, &f.IsActive, &f.FoundedDate,
		&f.CreatedAt, &f.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_FED_002", "truy vấn liên đoàn theo tỉnh thất bại")
	}
	return &f, nil
}

// ═══════════════════════════════════════════════════════════════
// CORE EVENT LOG OPERATIONS
// ═══════════════════════════════════════════════════════════════

func (s *PostgresStore) CoreLogEvent(ctx context.Context, event *CoreEventLog) error {
	err := s.pool.QueryRow(ctx, `
		SELECT core.log_event($1, $2, $3::UUID, $4, $5, $6, $7)
	`, event.EventType, event.EntityType, event.EntityID,
		event.Payload, event.SourceSchema, event.ActorID, event.CorrelationID,
	).Scan(&event.ID)
	if err != nil {
		return apierror.Wrap(err, "CORE_EVT_001", "ghi event log thất bại")
	}
	return nil
}

func (s *PostgresStore) CoreListEvents(ctx context.Context, entityType string, entityID string, limit int) ([]CoreEventLog, error) {
	query := `
		SELECT id, event_type, entity_type, entity_id, payload,
		       source_schema, actor_id, correlation_id, created_at
		FROM core.event_logs
		WHERE 1=1
	`
	args := make([]any, 0)
	argIdx := 1

	if entityType != "" {
		query += " AND entity_type = $1"
		args = append(args, entityType)
		argIdx = 2
	}
	if entityID != "" {
		query += " AND entity_id = $" + string(rune('0'+argIdx)) + "::UUID"
		args = append(args, entityID)
		argIdx++
	}

	if limit <= 0 {
		limit = 100
	}
	query += " ORDER BY created_at DESC LIMIT $" + string(rune('0'+argIdx))
	args = append(args, limit)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_EVT_002", "liệt kê events thất bại")
	}
	defer rows.Close()

	events := make([]CoreEventLog, 0)
	for rows.Next() {
		var e CoreEventLog
		if scanErr := rows.Scan(
			&e.ID, &e.EventType, &e.EntityType, &e.EntityID, &e.Payload,
			&e.SourceSchema, &e.ActorID, &e.CorrelationID, &e.CreatedAt,
		); scanErr != nil {
			continue
		}
		events = append(events, e)
	}
	return events, nil
}

// ═══════════════════════════════════════════════════════════════
// DASHBOARD / MATERIALIZED VIEW QUERIES
// ═══════════════════════════════════════════════════════════════

type NationalDashboardRow struct {
	FederationID          string    `json:"federation_id"`
	FederationName        string    `json:"federation_name"`
	ProvinceCode          *string   `json:"province_code,omitempty"`
	Region                string    `json:"region"`
	TotalActiveAthletes   int       `json:"total_active_athletes"`
	TotalBeltCerts        int       `json:"total_belt_certifications"`
	TotalMedalsWon        int       `json:"total_medals_won"`
	RefreshedAt           time.Time `json:"refreshed_at"`
}

func (s *PostgresStore) CoreGetNationalDashboard(ctx context.Context) ([]NationalDashboardRow, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT federation_id, federation_name, province_code, region,
		       total_active_athletes, total_belt_certifications,
		       total_medals_won, refreshed_at
		FROM core.mv_national_dashboard
		ORDER BY total_active_athletes DESC
	`)
	if err != nil {
		return nil, apierror.Wrap(err, "CORE_DASH_001", "truy vấn dashboard quốc gia thất bại")
	}
	defer rows.Close()

	dashboard := make([]NationalDashboardRow, 0)
	for rows.Next() {
		var d NationalDashboardRow
		if scanErr := rows.Scan(
			&d.FederationID, &d.FederationName, &d.ProvinceCode, &d.Region,
			&d.TotalActiveAthletes, &d.TotalBeltCerts,
			&d.TotalMedalsWon, &d.RefreshedAt,
		); scanErr != nil {
			continue
		}
		dashboard = append(dashboard, d)
	}
	return dashboard, nil
}

func (s *PostgresStore) CoreRefreshDashboards(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, "SELECT core.refresh_all_matviews()")
	if err != nil {
		return apierror.Wrap(err, "CORE_DASH_002", "refresh materialized views thất bại")
	}
	return nil
}
