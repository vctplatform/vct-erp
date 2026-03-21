package usecase

import "time"

// AccessMetadata captures who opened a finance dashboard report and under which filters.
type AccessMetadata struct {
	CompanyCode string
	ActorID     string
	ActorRole   string
	IPAddress   string
	UserAgent   string
	Filters     map[string]string
}

// CashRunwayInput controls the executive cash-runway projection request.
type CashRunwayInput struct {
	Access AccessMetadata
	AsOf   time.Time
	Months int
}
