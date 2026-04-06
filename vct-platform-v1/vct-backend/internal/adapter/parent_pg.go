package adapter

import (
	"context"
	"database/sql"
	"fmt"

	"vct-platform/backend/internal/domain/parent"
)

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — PARENT MODULE POSTGRESQL ADAPTERS
// Implements ParentLinkStore, ConsentStore, AttendanceStore,
// ResultStore interfaces using PostgreSQL.
// ═══════════════════════════════════════════════════════════════

// ── Parent Link Store ───────────────────────────────────────

type PgParentLinkStore struct {
	db *sql.DB
}

func NewPgParentLinkStore(db *sql.DB) *PgParentLinkStore {
	return &PgParentLinkStore{db: db}
}

func (s *PgParentLinkStore) ListByParent(ctx context.Context, parentID string) ([]parent.ParentLink, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, parent_id, parent_name, athlete_id, athlete_name, club_name, belt_level, relation, status, requested_at, approved_at
		 FROM parent_links WHERE parent_id = $1 ORDER BY requested_at DESC`, parentID,
	)
	if err != nil {
		return nil, fmt.Errorf("parent list links: %w", err)
	}
	defer rows.Close()

	var out []parent.ParentLink
	for rows.Next() {
		var l parent.ParentLink
		var approvedAt sql.NullTime
		if err := rows.Scan(&l.ID, &l.ParentID, &l.ParentName, &l.AthleteID, &l.AthleteName, &l.ClubName, &l.BeltLevel, &l.Relation, &l.Status, &l.RequestedAt, &approvedAt); err != nil {
			return nil, fmt.Errorf("parent scan link: %w", err)
		}
		if approvedAt.Valid {
			l.ApprovedAt = &approvedAt.Time
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (s *PgParentLinkStore) GetByID(ctx context.Context, id string) (*parent.ParentLink, error) {
	var l parent.ParentLink
	var approvedAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, parent_id, parent_name, athlete_id, athlete_name, club_name, belt_level, relation, status, requested_at, approved_at
		 FROM parent_links WHERE id = $1`, id,
	).Scan(&l.ID, &l.ParentID, &l.ParentName, &l.AthleteID, &l.AthleteName, &l.ClubName, &l.BeltLevel, &l.Relation, &l.Status, &l.RequestedAt, &approvedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("parent link %s not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("parent get link: %w", err)
	}
	if approvedAt.Valid {
		l.ApprovedAt = &approvedAt.Time
	}
	return &l, nil
}

func (s *PgParentLinkStore) Create(ctx context.Context, l parent.ParentLink) (*parent.ParentLink, error) {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO parent_links (id, parent_id, parent_name, athlete_id, athlete_name, club_name, belt_level, relation, status, requested_at, approved_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		l.ID, l.ParentID, l.ParentName, l.AthleteID, l.AthleteName, l.ClubName, l.BeltLevel, l.Relation, l.Status, l.RequestedAt, l.ApprovedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("parent create link: %w", err)
	}
	return &l, nil
}

func (s *PgParentLinkStore) Update(ctx context.Context, id string, patch parent.LinkUpdate) error {
	// Build dynamic update
	setClauses := []string{}
	args := []any{}
	argIdx := 1

	if patch.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, string(*patch.Status))
		argIdx++
	}
	if patch.ApprovedAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("approved_at = $%d", argIdx))
		args = append(args, *patch.ApprovedAt)
		argIdx++
	}
	if len(setClauses) == 0 {
		return nil
	}

	query := "UPDATE parent_links SET "
	for i, c := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += c
	}
	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("parent update link: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("parent link %s not found", id)
	}
	return nil
}

func (s *PgParentLinkStore) Delete(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM parent_links WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("parent delete link: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("parent link %s not found", id)
	}
	return nil
}

func (s *PgParentLinkStore) IsChildOfParent(ctx context.Context, parentID, athleteID string) bool {
	var count int
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM parent_links WHERE parent_id = $1 AND athlete_id = $2 AND status = 'approved'`,
		parentID, athleteID,
	).Scan(&count)
	return err == nil && count > 0
}

// ── Consent Store ───────────────────────────────────────────

type PgConsentStore struct {
	db *sql.DB
}

func NewPgConsentStore(db *sql.DB) *PgConsentStore {
	return &PgConsentStore{db: db}
}

func (s *PgConsentStore) ListByParent(ctx context.Context, parentID string) ([]parent.ConsentRecord, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, parent_id, athlete_id, athlete_name, type, title, description, status, signed_at, expires_at, revoked_at
		 FROM parent_consents WHERE parent_id = $1 ORDER BY signed_at DESC`, parentID,
	)
	if err != nil {
		return nil, fmt.Errorf("parent list consents: %w", err)
	}
	defer rows.Close()

	var out []parent.ConsentRecord
	for rows.Next() {
		var c parent.ConsentRecord
		var expiresAt, revokedAt sql.NullTime
		if err := rows.Scan(&c.ID, &c.ParentID, &c.AthleteID, &c.AthleteName, &c.Type, &c.Title, &c.Description, &c.Status, &c.SignedAt, &expiresAt, &revokedAt); err != nil {
			return nil, fmt.Errorf("parent scan consent: %w", err)
		}
		if expiresAt.Valid {
			c.ExpiresAt = &expiresAt.Time
		}
		if revokedAt.Valid {
			c.RevokedAt = &revokedAt.Time
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *PgConsentStore) ListByAthlete(ctx context.Context, athleteID string) ([]parent.ConsentRecord, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, parent_id, athlete_id, athlete_name, type, title, description, status, signed_at, expires_at, revoked_at
		 FROM parent_consents WHERE athlete_id = $1 ORDER BY signed_at DESC`, athleteID,
	)
	if err != nil {
		return nil, fmt.Errorf("parent list consents by athlete: %w", err)
	}
	defer rows.Close()

	var out []parent.ConsentRecord
	for rows.Next() {
		var c parent.ConsentRecord
		var expiresAt, revokedAt sql.NullTime
		if err := rows.Scan(&c.ID, &c.ParentID, &c.AthleteID, &c.AthleteName, &c.Type, &c.Title, &c.Description, &c.Status, &c.SignedAt, &expiresAt, &revokedAt); err != nil {
			return nil, fmt.Errorf("parent scan consent: %w", err)
		}
		if expiresAt.Valid {
			c.ExpiresAt = &expiresAt.Time
		}
		if revokedAt.Valid {
			c.RevokedAt = &revokedAt.Time
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *PgConsentStore) GetByID(ctx context.Context, id string) (*parent.ConsentRecord, error) {
	var c parent.ConsentRecord
	var expiresAt, revokedAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, parent_id, athlete_id, athlete_name, type, title, description, status, signed_at, expires_at, revoked_at
		 FROM parent_consents WHERE id = $1`, id,
	).Scan(&c.ID, &c.ParentID, &c.AthleteID, &c.AthleteName, &c.Type, &c.Title, &c.Description, &c.Status, &c.SignedAt, &expiresAt, &revokedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("consent %s not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("parent get consent: %w", err)
	}
	if expiresAt.Valid {
		c.ExpiresAt = &expiresAt.Time
	}
	if revokedAt.Valid {
		c.RevokedAt = &revokedAt.Time
	}
	return &c, nil
}

func (s *PgConsentStore) Create(ctx context.Context, c parent.ConsentRecord) (*parent.ConsentRecord, error) {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO parent_consents (id, parent_id, athlete_id, athlete_name, type, title, description, status, signed_at, expires_at, revoked_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		c.ID, c.ParentID, c.AthleteID, c.AthleteName, c.Type, c.Title, c.Description, c.Status, c.SignedAt, c.ExpiresAt, c.RevokedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("parent create consent: %w", err)
	}
	return &c, nil
}

func (s *PgConsentStore) Update(ctx context.Context, id string, patch parent.ConsentUpdate) error {
	setClauses := []string{}
	args := []any{}
	argIdx := 1

	if patch.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, string(*patch.Status))
		argIdx++
	}
	if patch.RevokedAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("revoked_at = $%d", argIdx))
		args = append(args, *patch.RevokedAt)
		argIdx++
	}
	if len(setClauses) == 0 {
		return nil
	}

	query := "UPDATE parent_consents SET "
	for i, c := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += c
	}
	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("parent update consent: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("consent %s not found", id)
	}
	return nil
}

// ── Attendance Store ────────────────────────────────────────

type PgParentAttendanceStore struct {
	db *sql.DB
}

func NewPgParentAttendanceStore(db *sql.DB) *PgParentAttendanceStore {
	return &PgParentAttendanceStore{db: db}
}

func (s *PgParentAttendanceStore) ListByAthlete(ctx context.Context, athleteID string) ([]parent.AttendanceSummary, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT date, session, status, coach FROM parent_attendance WHERE athlete_id = $1 ORDER BY date DESC`,
		athleteID,
	)
	if err != nil {
		return nil, fmt.Errorf("parent list attendance: %w", err)
	}
	defer rows.Close()

	var out []parent.AttendanceSummary
	for rows.Next() {
		var a parent.AttendanceSummary
		if err := rows.Scan(&a.Date, &a.Session, &a.Status, &a.Coach); err != nil {
			return nil, fmt.Errorf("parent scan attendance: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// ── Results Store ───────────────────────────────────────────

type PgParentResultStore struct {
	db *sql.DB
}

func NewPgParentResultStore(db *sql.DB) *PgParentResultStore {
	return &PgParentResultStore{db: db}
}

func (s *PgParentResultStore) ListByAthlete(ctx context.Context, athleteID string) ([]parent.ChildResult, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT tournament, category, result, date FROM parent_results WHERE athlete_id = $1 ORDER BY date DESC`,
		athleteID,
	)
	if err != nil {
		return nil, fmt.Errorf("parent list results: %w", err)
	}
	defer rows.Close()

	var out []parent.ChildResult
	for rows.Next() {
		var r parent.ChildResult
		if err := rows.Scan(&r.Tournament, &r.Category, &r.Result, &r.Date); err != nil {
			return nil, fmt.Errorf("parent scan result: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
