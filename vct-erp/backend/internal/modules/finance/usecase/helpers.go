package usecase

import (
	"strings"
	"time"
)

func normalizeCurrency(currency string) string {
	if trimmed := strings.ToUpper(strings.TrimSpace(currency)); trimmed != "" {
		return trimmed
	}
	return "VND"
}

func monthStart(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
