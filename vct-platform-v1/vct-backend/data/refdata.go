package refdata

import (
	"embed"
	"encoding/json"
	"fmt"
)

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — Reference Data Loader
// Loads default/reference JSON data from backend/data/ directory.
// All JSON files are editable without recompiling.
// ═══════════════════════════════════════════════════════════════

//go:embed *.json
var dataFS embed.FS

// ── Generic JSON loader ─────────────────────────────────────

func loadJSON[T any](filename string) (T, error) {
	var result T
	data, err := dataFS.ReadFile(filename)
	if err != nil {
		return result, fmt.Errorf("refdata: read %s: %w", filename, err)
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("refdata: parse %s: %w", filename, err)
	}
	return result, nil
}

// ── Belt Ranks ──────────────────────────────────────────────

type BeltRank struct {
	Code              string `json:"code"`
	NameVI            string `json:"name_vi"`
	NameEN            string `json:"name_en"`
	ColorHex          string `json:"color_hex"`
	Order             int    `json:"order"`
	MinTrainingMonths int    `json:"min_training_months"`
	Description       string `json:"description"`
}

type beltRanksFile struct {
	BeltRanks []BeltRank `json:"belt_ranks"`
}

func LoadBeltRanks() ([]BeltRank, error) {
	f, err := loadJSON[beltRanksFile]("belt_ranks.json")
	if err != nil {
		return nil, err
	}
	return f.BeltRanks, nil
}

// ── Weight Classes ──────────────────────────────────────────

type WeightClass struct {
	Code  string   `json:"code"`
	Name  string   `json:"name"`
	MinKg float64  `json:"min_kg"`
	MaxKg *float64 `json:"max_kg"`
	Order int      `json:"order"`
}

type weightClassesFile struct {
	WeightClasses map[string][]WeightClass `json:"weight_classes"`
}

func LoadWeightClasses() (map[string][]WeightClass, error) {
	f, err := loadJSON[weightClassesFile]("standard_weight_classes.json")
	if err != nil {
		return nil, err
	}
	return f.WeightClasses, nil
}

// ── Age Groups ──────────────────────────────────────────────

type AgeGroup struct {
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	MinAge        int      `json:"min_age"`
	MaxAge        *int     `json:"max_age"`
	Description   string   `json:"description"`
	AllowedEvents []string `json:"allowed_events"`
}

type ageGroupsFile struct {
	AgeGroups []AgeGroup `json:"age_groups"`
}

func LoadAgeGroups() ([]AgeGroup, error) {
	f, err := loadJSON[ageGroupsFile]("standard_age_groups.json")
	if err != nil {
		return nil, err
	}
	return f.AgeGroups, nil
}

// ── Standard Forms ──────────────────────────────────────────

type StandardForm struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Origin     string `json:"origin,omitempty"`
	Weapon     string `json:"weapon,omitempty"`
	Difficulty int    `json:"difficulty"`
	Gender     string `json:"gender"`
	Type       string `json:"type"`
	TeamSize   int    `json:"team_size,omitempty"`
}

type formCategory struct {
	Name  string         `json:"name"`
	Forms []StandardForm `json:"forms"`
}

type standardFormsFile struct {
	Categories map[string]formCategory `json:"categories"`
}

func LoadStandardForms() (map[string][]StandardForm, error) {
	f, err := loadJSON[standardFormsFile]("standard_forms.json")
	if err != nil {
		return nil, err
	}
	result := make(map[string][]StandardForm, len(f.Categories))
	for key, cat := range f.Categories {
		result[key] = cat.Forms
	}
	return result, nil
}

// ── System Config ───────────────────────────────────────────

type TournamentDefaults struct {
	MaxAthletesPerEvent int     `json:"max_athletes_per_team_per_event"`
	MaxEventsPerAthlete int     `json:"max_events_per_athlete"`
	WeighInToleranceKg  float64 `json:"weigh_in_tolerance_kg"`
	MinAgeCompetition   int     `json:"min_age_competition"`
	MaxRoundsCombat     int     `json:"max_rounds_combat"`
	RoundDuration       int     `json:"round_duration_seconds"`
	BreakBetweenRounds  int     `json:"break_between_rounds_seconds"`
	MinJudgesForm       int     `json:"min_judges_form"`
	MaxJudgesForm       int     `json:"max_judges_form"`
	ScoreDecimalPlaces  int     `json:"score_decimal_places"`
	DropHighestScore    bool    `json:"drop_highest_score"`
	DropLowestScore     bool    `json:"drop_lowest_score"`
}

type SessionConfig struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type PlatformConfig struct {
	Name             string `json:"name"`
	FullName         string `json:"full_name"`
	DefaultLanguage  string `json:"default_language"`
	Timezone         string `json:"timezone"`
	Currency         string `json:"currency"`
	CurrencySymbol   string `json:"currency_symbol"`
	DateFormat       string `json:"date_format"`
	PhoneCountryCode string `json:"phone_country_code"`
}

type SystemConfig struct {
	Platform   PlatformConfig `json:"platform"`
	Tournament struct {
		TournamentDefaults
		Sessions []SessionConfig `json:"sessions"`
	} `json:"tournament"`
}

func LoadSystemConfig() (*SystemConfig, error) {
	return loadJSON[*SystemConfig]("system_configs.json")
}

// ── Roles & Permissions ─────────────────────────────────────

type Role struct {
	Code        string   `json:"code"`
	NameVI      string   `json:"name_vi"`
	NameEN      string   `json:"name_en"`
	Level       int      `json:"level"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type rolesFile struct {
	Roles []Role `json:"roles"`
}

func LoadRoles() ([]Role, error) {
	f, err := loadJSON[rolesFile]("roles_permissions.json")
	if err != nil {
		return nil, err
	}
	return f.Roles, nil
}

// ── Violation Types ─────────────────────────────────────────

type ViolationType struct {
	Code             string   `json:"code"`
	NameVI           string   `json:"name_vi"`
	NameEN           string   `json:"name_en"`
	Severity         string   `json:"severity"`
	ApplicableTo     []string `json:"applicable_to"`
	DefaultSanctions []string `json:"default_sanctions"`
	Description      string   `json:"description"`
}

type SanctionType struct {
	Code        string `json:"code"`
	NameVI      string `json:"name_vi"`
	NameEN      string `json:"name_en"`
	Description string `json:"description"`
}

type violationFile struct {
	ViolationTypes []ViolationType `json:"violation_types"`
	SanctionTypes  []SanctionType  `json:"sanction_types"`
}

func LoadViolationTypes() ([]ViolationType, error) {
	f, err := loadJSON[violationFile]("violation_types.json")
	if err != nil {
		return nil, err
	}
	return f.ViolationTypes, nil
}

func LoadSanctionTypes() ([]SanctionType, error) {
	f, err := loadJSON[violationFile]("violation_types.json")
	if err != nil {
		return nil, err
	}
	return f.SanctionTypes, nil
}

// ── Notification Templates ──────────────────────────────────

type NotificationTemplate struct {
	Code          string   `json:"code"`
	Type          string   `json:"type"`
	Category      string   `json:"category"`
	TitleTemplate string   `json:"title_template"`
	BodyTemplate  string   `json:"body_template"`
	Recipients    []string `json:"recipients"`
	Channels      []string `json:"channels"`
}

type notifFile struct {
	Templates []NotificationTemplate `json:"notification_templates"`
}

func LoadNotificationTemplates() ([]NotificationTemplate, error) {
	f, err := loadJSON[notifFile]("notification_templates.json")
	if err != nil {
		return nil, err
	}
	return f.Templates, nil
}

// ── Email Templates ─────────────────────────────────────────

type EmailTemplate struct {
	Code     string `json:"code"`
	Subject  string `json:"subject"`
	Category string `json:"category"`
	BodyHTML string `json:"body_html"`
}

type emailFile struct {
	Templates []EmailTemplate `json:"email_templates"`
}

func LoadEmailTemplates() ([]EmailTemplate, error) {
	f, err := loadJSON[emailFile]("email_templates.json")
	if err != nil {
		return nil, err
	}
	return f.Templates, nil
}

// ── Sample Clubs ────────────────────────────────────────────

type SampleClub struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	ProvinceCode string `json:"province_code"`
	ProvinceName string `json:"province_name"`
	SchoolStyle  string `json:"school_style"`
	FoundedYear  int    `json:"founded_year"`
	MemberCount  int    `json:"member_count"`
	CoachName    string `json:"coach_name"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Status       string `json:"status"`
	Grade        string `json:"grade"`
}

type clubsFile struct {
	Clubs []SampleClub `json:"sample_clubs"`
}

func LoadSampleClubs() ([]SampleClub, error) {
	f, err := loadJSON[clubsFile]("sample_clubs.json")
	if err != nil {
		return nil, err
	}
	return f.Clubs, nil
}

// ── Scoring Criteria ────────────────────────────────────────

type ScoringCriterion struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	WeightPct   int    `json:"weight_pct"`
	Description string `json:"description"`
}

type ScoringDeduction struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Deduction float64 `json:"deduction"`
}

func LoadScoringCriteria() (map[string]any, error) {
	return loadJSON[map[string]any]("scoring_criteria.json")
}

// ── Fee Schedule ────────────────────────────────────────────

type TournamentFeeLevel struct {
	Level         string `json:"level"`
	Name          string `json:"name"`
	FeePerAthlete int    `json:"fee_per_athlete"`
	FeePerTeam    int    `json:"fee_per_team"`
	Deposit       int    `json:"deposit"`
}

func LoadFeeSchedule() (map[string]any, error) {
	return loadJSON[map[string]any]("fee_schedule.json")
}

// ── Countries ───────────────────────────────────────────────

type Country struct {
	Code          string `json:"code"`
	NameVI        string `json:"name_vi"`
	NameEN        string `json:"name_en"`
	Region        string `json:"region"`
	HasFederation bool   `json:"has_federation"`
	Vovinam       bool   `json:"vovinam"`
	VCT           bool   `json:"vct"`
	PhoneCode     string `json:"phone_code"`
}

type countriesFile struct {
	Countries []Country `json:"countries"`
}

func LoadCountries() ([]Country, error) {
	f, err := loadJSON[countriesFile]("countries.json")
	if err != nil {
		return nil, err
	}
	return f.Countries, nil
}

// ── Equipment Standards ─────────────────────────────────────

func LoadEquipmentStandards() (map[string]any, error) {
	return loadJSON[map[string]any]("equipment_standards.json")
}

// ── Training Syllabus ───────────────────────────────────────

type BeltRequirement struct {
	BeltCode          string   `json:"belt_code"`
	BeltName          string   `json:"belt_name"`
	MinTrainingMonths int      `json:"min_training_months"`
	Theory            []string `json:"theory"`
}

type trainingSyllabusFile struct {
	BeltRequirements []BeltRequirement `json:"belt_requirements"`
}

func LoadTrainingSyllabus() ([]BeltRequirement, error) {
	f, err := loadJSON[trainingSyllabusFile]("training_syllabus.json")
	if err != nil {
		return nil, err
	}
	return f.BeltRequirements, nil
}
