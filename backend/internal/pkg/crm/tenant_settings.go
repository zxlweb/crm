package crm

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/datatypes"
)

var (
	ErrSettingKeyInvalid      = errors.New("setting_key_invalid")
	ErrCustomFieldKeyConflict = errors.New("custom_field_key_conflict")
	ErrCustomFieldTypeInvalid = errors.New("custom_field_type_invalid")
	fieldKeyPattern           = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
)

const (
	LocaleZhCN = "zh-CN"
	LocaleEnUS = "en-US"
)

// BusinessSwitches tenant feature toggles.
type BusinessSwitches struct {
	AIPreviewEnabled bool   `json:"ai_preview_enabled"`
	LeadImportMode   string `json:"lead_import_mode"`
}

// TenantSettings persisted in tenants.config (Phase 4).
type TenantSettings struct {
	DefaultLocale    string           `json:"default_locale"`
	Timezone         string           `json:"timezone"`
	BusinessSwitches BusinessSwitches `json:"business_switches"`
	SalesQuota       SalesQuota       `json:"sales_quota"`
}

func DefaultTenantSettings() TenantSettings {
	return TenantSettings{
		DefaultLocale: LocaleZhCN,
		Timezone:      "Asia/Shanghai",
		BusinessSwitches: BusinessSwitches{
			AIPreviewEnabled: false,
			LeadImportMode:   "manual_review",
		},
		SalesQuota: SalesQuota{Currency: "CNY", Period: time.Now().Format("2006-01")},
	}
}

func ParseTenantSettings(raw datatypes.JSON) TenantSettings {
	out := DefaultTenantSettings()
	crm := ParseTenantCRMConfig(raw)
	out.SalesQuota = crm.SalesQuota
	if len(raw) == 0 {
		return out
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return out
	}
	if v, ok := m["default_locale"]; ok {
		_ = json.Unmarshal(v, &out.DefaultLocale)
	}
	if v, ok := m["timezone"]; ok {
		_ = json.Unmarshal(v, &out.Timezone)
	}
	if v, ok := m["business_switches"]; ok {
		_ = json.Unmarshal(v, &out.BusinessSwitches)
	}
	out.normalize()
	return out
}

func (s *TenantSettings) normalize() {
	if s.DefaultLocale == "" {
		s.DefaultLocale = LocaleZhCN
	}
	if s.Timezone == "" {
		s.Timezone = "Asia/Shanghai"
	}
	if s.BusinessSwitches.LeadImportMode == "" {
		s.BusinessSwitches.LeadImportMode = "manual_review"
	}
	if s.SalesQuota.Currency == "" {
		s.SalesQuota.Currency = "CNY"
	}
}

func (s TenantSettings) ValidatePatch(name *string, patch TenantSettingsPatch) error {
	if name != nil && strings.TrimSpace(*name) == "" {
		return ErrSettingKeyInvalid
	}
	if patch.DefaultLocale != nil && !ValidLocale(*patch.DefaultLocale) {
		return ErrSettingKeyInvalid
	}
	if patch.Timezone != nil && strings.TrimSpace(*patch.Timezone) == "" {
		return ErrSettingKeyInvalid
	}
	if patch.BusinessSwitches != nil {
		if patch.BusinessSwitches.LeadImportMode != nil &&
			*patch.BusinessSwitches.LeadImportMode != "manual_review" &&
			*patch.BusinessSwitches.LeadImportMode != "auto_merge" {
			return ErrSettingKeyInvalid
		}
	}
	if patch.SalesQuota != nil {
		if patch.SalesQuota.Currency != "" && patch.SalesQuota.Currency != "CNY" && patch.SalesQuota.Currency != "USD" {
			return ErrSettingKeyInvalid
		}
	}
	return nil
}

type BusinessSwitchesPatch struct {
	AIPreviewEnabled *bool
	LeadImportMode   *string
}

type TenantSettingsPatch struct {
	DefaultLocale    *string
	Timezone         *string
	BusinessSwitches *BusinessSwitchesPatch
	SalesQuota       *SalesQuota
}

func ValidLocale(locale string) bool {
	return locale == LocaleZhCN || locale == LocaleEnUS
}

func ValidEntityType(t string) bool {
	switch t {
	case "account", "contact", "lead", "deal":
		return true
	default:
		return false
	}
}

func ValidFieldType(t string) bool {
	switch t {
	case "text", "select", "date":
		return true
	default:
		return false
	}
}

func ValidFieldKey(key string) bool {
	return fieldKeyPattern.MatchString(key)
}

func MergeTenantConfig(raw datatypes.JSON, name *string, patch TenantSettingsPatch) (datatypes.JSON, error) {
	if err := DefaultTenantSettings().ValidatePatch(name, patch); err != nil {
		return nil, err
	}
	var m map[string]any
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &m); err != nil {
			m = map[string]any{}
		}
	}
	if m == nil {
		m = map[string]any{}
	}
	if patch.DefaultLocale != nil {
		m["default_locale"] = *patch.DefaultLocale
	}
	if patch.Timezone != nil {
		m["timezone"] = *patch.Timezone
	}
	if patch.BusinessSwitches != nil {
		bs, _ := m["business_switches"].(map[string]any)
		if bs == nil {
			bs = map[string]any{}
		}
		if patch.BusinessSwitches.AIPreviewEnabled != nil {
			bs["ai_preview_enabled"] = *patch.BusinessSwitches.AIPreviewEnabled
		}
		if patch.BusinessSwitches.LeadImportMode != nil {
			bs["lead_import_mode"] = *patch.BusinessSwitches.LeadImportMode
		}
		m["business_switches"] = bs
	}
	if patch.SalesQuota != nil {
		sq, _ := m["sales_quota"].(map[string]any)
		if sq == nil {
			sq = map[string]any{}
		}
		if patch.SalesQuota.Amount > 0 {
			sq["amount"] = patch.SalesQuota.Amount
		}
		if patch.SalesQuota.Currency != "" {
			sq["currency"] = patch.SalesQuota.Currency
		}
		if patch.SalesQuota.Period != "" {
			sq["period"] = patch.SalesQuota.Period
		}
		m["sales_quota"] = sq
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(b), nil
}

// FeatureDefinition describes a tenant setting switch for GET /settings/features.
type FeatureDefinition struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Default     any    `json:"default"`
	Description string `json:"description"`
}

func SettingsFeatureCatalog(current TenantSettings) []map[string]any {
	defs := []FeatureDefinition{
		{Key: "default_locale", Type: "string", Default: LocaleZhCN, Description: "Default locale"},
		{Key: "timezone", Type: "string", Default: "Asia/Shanghai", Description: "Tenant timezone"},
		{Key: "ai_preview_enabled", Type: "boolean", Default: false, Description: "AI Preview master switch"},
		{Key: "lead_import_mode", Type: "string", Default: "manual_review", Description: "Lead import policy"},
		{Key: "sales_quota.amount", Type: "number", Default: 0, Description: "Sales quota amount"},
		{Key: "sales_quota.currency", Type: "string", Default: "CNY", Description: "Quota currency"},
	}
	values := map[string]any{
		"default_locale":       current.DefaultLocale,
		"timezone":             current.Timezone,
		"ai_preview_enabled":   current.BusinessSwitches.AIPreviewEnabled,
		"lead_import_mode":     current.BusinessSwitches.LeadImportMode,
		"sales_quota.amount":   current.SalesQuota.Amount,
		"sales_quota.currency": current.SalesQuota.Currency,
	}
	out := make([]map[string]any, 0, len(defs))
	for _, d := range defs {
		out = append(out, map[string]any{
			"key":         d.Key,
			"type":        d.Type,
			"default":     d.Default,
			"description": d.Description,
			"value":       values[d.Key],
		})
	}
	return out
}
