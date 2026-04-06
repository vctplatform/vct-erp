package divisions

import (
	"embed"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

//go:embed data/vietnam_divisions.json
var dataFS embed.FS

// Province represents a Vietnamese province or municipality.
type Province struct {
	Name         string `json:"name"`
	Code         int    `json:"code"`
	DivisionType string `json:"division_type"`
	Codename     string `json:"codename"`
	PhoneCode    int    `json:"phone_code"`
	Wards        []Ward `json:"wards"`
}

// ProvinceInfo is a lightweight version without wards for list endpoints.
type ProvinceInfo struct {
	Name         string `json:"name"`
	Code         int    `json:"code"`
	DivisionType string `json:"division_type"`
	Codename     string `json:"codename"`
	PhoneCode    int    `json:"phone_code"`
	WardCount    int    `json:"ward_count"`
}

// Ward represents a ward/commune/township.
type Ward struct {
	Name         string `json:"name"`
	Code         int    `json:"code"`
	DivisionType string `json:"division_type"`
	Codename     string `json:"codename"`
	ProvinceCode int    `json:"province_code"`
}

// Store holds the loaded administrative division data.
type Store struct {
	provinces   []Province
	provinceMap map[int]*Province // code → Province
	once        sync.Once
}

var defaultStore = &Store{}

// Default returns the singleton Store instance with data loaded.
func Default() *Store {
	defaultStore.once.Do(func() {
		defaultStore.load()
	})
	return defaultStore
}

func (s *Store) load() {
	data, err := dataFS.ReadFile("data/vietnam_divisions.json")
	if err != nil {
		slog.Error("divisions: failed to read embedded data", slog.String("error", err.Error()))
		return
	}
	if err := json.Unmarshal(data, &s.provinces); err != nil {
		slog.Error("divisions: failed to parse data", slog.String("error", err.Error()))
		return
	}
	s.provinceMap = make(map[int]*Province, len(s.provinces))
	for i := range s.provinces {
		s.provinceMap[s.provinces[i].Code] = &s.provinces[i]
	}
	slog.Info("divisions loaded", slog.Int("provinces", len(s.provinces)), slog.Int("wards", s.TotalWards()))
}

// Provinces returns all provinces as lightweight ProvinceInfo (without wards).
func (s *Store) Provinces() []ProvinceInfo {
	out := make([]ProvinceInfo, len(s.provinces))
	for i, p := range s.provinces {
		out[i] = ProvinceInfo{
			Name:         p.Name,
			Code:         p.Code,
			DivisionType: p.DivisionType,
			Codename:     p.Codename,
			PhoneCode:    p.PhoneCode,
			WardCount:    len(p.Wards),
		}
	}
	return out
}

// SearchProvinces returns provinces matching the query string.
func (s *Store) SearchProvinces(q string) []ProvinceInfo {
	q = normalizeQuery(q)
	if q == "" {
		return s.Provinces()
	}
	var out []ProvinceInfo
	for _, p := range s.provinces {
		if containsNormalized(p.Name, q) || containsNormalized(p.Codename, q) {
			out = append(out, ProvinceInfo{
				Name:         p.Name,
				Code:         p.Code,
				DivisionType: p.DivisionType,
				Codename:     p.Codename,
				PhoneCode:    p.PhoneCode,
				WardCount:    len(p.Wards),
			})
		}
	}
	return out
}

// Province returns a specific province by code, or nil if not found.
func (s *Store) Province(code int) *Province {
	return s.provinceMap[code]
}

// Wards returns all wards of a given province.
func (s *Store) Wards(provinceCode int) []Ward {
	p := s.provinceMap[provinceCode]
	if p == nil {
		return nil
	}
	return p.Wards
}

// SearchWards returns wards of a province matching the query string.
func (s *Store) SearchWards(provinceCode int, q string) []Ward {
	p := s.provinceMap[provinceCode]
	if p == nil {
		return nil
	}
	q = normalizeQuery(q)
	if q == "" {
		return p.Wards
	}
	var out []Ward
	for _, w := range p.Wards {
		if containsNormalized(w.Name, q) || containsNormalized(w.Codename, q) {
			out = append(out, w)
		}
	}
	return out
}

// TotalWards returns the total number of wards across all provinces.
func (s *Store) TotalWards() int {
	total := 0
	for _, p := range s.provinces {
		total += len(p.Wards)
	}
	return total
}

// ── Vietnamese Text Normalization ────────────────────────────

// normalizeQuery lowercases and strips diacritics for fuzzy matching.
func normalizeQuery(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// containsNormalized checks if target contains query (case-insensitive).
func containsNormalized(target, query string) bool {
	target = strings.ToLower(target)
	if strings.Contains(target, query) {
		return true
	}
	// Also try without diacritics
	return strings.Contains(removeDiacritics(target), removeDiacritics(query))
}

// removeDiacritics strips Vietnamese diacritical marks for search.
func removeDiacritics(s string) string {
	t := norm.NFD.String(s)
	var b strings.Builder
	b.Grow(len(t))
	for _, r := range t {
		if !unicode.Is(unicode.Mn, r) { // Mn = Mark, Nonspacing
			b.WriteRune(r)
		}
	}
	// Handle special Vietnamese chars
	result := b.String()
	result = strings.ReplaceAll(result, "đ", "d")
	result = strings.ReplaceAll(result, "Đ", "D")
	return result
}
