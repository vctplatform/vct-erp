package auth

import (
	"vct-platform/backend/internal/auth"
)

// Re-export core authentication types to avoid direct dependency on internal/auth in domain modules.
type Principal = auth.Principal
type AuthUser = auth.AuthUser
type UserRole = auth.UserRole

// Standard Roles (re-exported constants)
const (
	RoleOwner               = auth.RoleOwner
	RoleAdmin               = auth.RoleAdmin
	RoleFederationPresident = auth.RoleFederationPresident
	RoleFederationSecretary = auth.RoleFederationSecretary
	RoleProvincialAdmin     = auth.RoleProvincialAdmin
	RoleTechnicalDirector   = auth.RoleTechnicalDirector
	RoleBTC                 = auth.RoleBTC
	RoleRefereeManager      = auth.RoleRefereeManager
	RoleReferee             = auth.RoleReferee
	RoleCoach               = auth.RoleCoach
	RoleDelegate            = auth.RoleDelegate
	RoleAthlete             = auth.RoleAthlete
	RoleMedicalStaff        = auth.RoleMedicalStaff
)
