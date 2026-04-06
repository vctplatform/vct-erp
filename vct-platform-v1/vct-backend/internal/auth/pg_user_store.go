package auth

import (
	"context"
	"strings"

	"vct-platform/backend/internal/apierror"
	"vct-platform/backend/internal/store"
)

// PgUserStore implements UserStore using the v3.0 core_store.go Data Access Layer.
type PgUserStore struct {
	dbStore *store.PostgresStore
}

// NewPgUserStore creates a PostgreSQL-backed user store using the new bridge layer.
func NewPgUserStore(dbStore *store.PostgresStore) *PgUserStore {
	return &PgUserStore{dbStore: dbStore}
}

// FindByEmail looks up a user by email in core.users.
func (s *PgUserStore) FindByEmail(ctx context.Context, email string) (*StoredUser, error) {
	emailStr := strings.ToLower(strings.TrimSpace(email))

	u, err := s.dbStore.CoreGetUserByEmail(ctx, emailStr)
	if err != nil {
		return nil, apierror.Wrap(err, "AUTH_500_DB", "lỗi truy vấn người dùng")
	}
	if u == nil {
		return nil, nil // not found
	}

	phone := ""
	if u.Phone != nil {
		phone = *u.Phone
	}

	return &StoredUser{
		ID:           u.ID,
		Email:        u.Email,
		Phone:        phone,
		PasswordHash: *u.PasswordHash,
		FullName:     u.FullName,
		Role:         UserRole(u.Role),
		IsActive:     u.IsActive,
	}, nil
}

// Create inserts a new user into core.users via bridge.
func (s *PgUserStore) Create(ctx context.Context, user *StoredUser) error {
	email := strings.ToLower(strings.TrimSpace(user.Email))

	// Ensure email uniqueness
	existing, err := s.dbStore.CoreGetUserByEmail(ctx, email)
	if err != nil {
		return apierror.Wrap(err, "AUTH_500_DB", "lỗi kiểm tra người dùng tồn tại")
	}
	if existing != nil {
		return apierror.New("AUTH_409_EXISTS", "email đăng nhập đã tồn tại")
	}

	phonePtr := &user.Phone
	if user.Phone == "" {
		phonePtr = nil
	}
	
	hashPtr := &user.PasswordHash

	coreUser := &store.CoreUser{
		Email:        email,
		Phone:        phonePtr,
		FullName:     user.FullName,
		PasswordHash: hashPtr,
		Role:         string(user.Role),
		IsActive:     user.IsActive,
	}

	if err := s.dbStore.CoreCreateUser(ctx, coreUser); err != nil {
		return apierror.Wrap(err, "AUTH_500_DB", "lỗi tạo người dùng")
	}

	user.ID = coreUser.ID
	return nil
}

// UpdateLastLogin is a no-op right now for v3.0 since we moved tracking to sessions/event_logs.
// If needed, we can add last_login_at back to core.users or rely on core.sessions.
func (s *PgUserStore) UpdateLastLogin(ctx context.Context, userID string) error {
	// Not needed in v3.0 core.users anymore, we have core.sessions.
	return nil
}

