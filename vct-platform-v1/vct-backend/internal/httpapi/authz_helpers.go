package httpapi

import (
	"net/http"

	"vct-platform/backend/internal/auth"
)

func requireFederationRead(w http.ResponseWriter, p auth.Principal) bool {
	return requireRole(w, p, federationReadRoles...)
}

func requireFederationWrite(w http.ResponseWriter, p auth.Principal) bool {
	return requireRole(w, p, federationWriteRoles...)
}

func requireClubRead(w http.ResponseWriter, p auth.Principal) bool {
	return requireRole(w, p, clubReadRoles...)
}

func requireClubWrite(w http.ResponseWriter, p auth.Principal) bool {
	return requireRole(w, p, clubWriteRoles...)
}

func requireProvincialRead(w http.ResponseWriter, p auth.Principal) bool {
	return requireRole(w, p, provincialReadRoles...)
}

func requireProvincialWrite(w http.ResponseWriter, p auth.Principal) bool {
	return requireRole(w, p, provincialWriteRoles...)
}
