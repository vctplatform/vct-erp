package auth

import "context"

// UserStore abstracts persistent user storage (PostgreSQL, etc.).
// When nil/unset, the auth Service falls back to in-memory credentials.
type UserStore interface {
	// FindByEmail looks up an active, non-deleted user by email (previously username).
	FindByEmail(ctx context.Context, email string) (*StoredUser, error)
	// Create inserts a new user record.
	Create(ctx context.Context, user *StoredUser) error
	// UpdateLastLogin bumps last_login_at for the given user ID.
	UpdateLastLogin(ctx context.Context, userID string) error
}

// StoredUser is the database representation of a user in core.users.
type StoredUser struct {
	ID           string
	Email        string
	Phone        string
	PasswordHash string
	FullName     string
	Role         UserRole // primary role
	IsActive     bool
}
