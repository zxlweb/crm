package crm

import (
	"encoding/json"
	"testing"

	"gorm.io/datatypes"
)

func TestMergeTenantConfig_BusinessSwitchesPartial(t *testing.T) {
	raw, _ := json.Marshal(map[string]any{
		"business_switches": map[string]any{
			"ai_preview_enabled": true,
			"lead_import_mode":   "manual_review",
		},
	})
	mode := "auto_merge"
	cfg, err := MergeTenantConfig(datatypes.JSON(raw), nil, TenantSettingsPatch{
		BusinessSwitches: &BusinessSwitchesPatch{LeadImportMode: &mode},
	})
	if err != nil {
		t.Fatal(err)
	}
	settings := ParseTenantSettings(cfg)
	if !settings.BusinessSwitches.AIPreviewEnabled {
		t.Fatalf("ai preview should remain true")
	}
	if settings.BusinessSwitches.LeadImportMode != "auto_merge" {
		t.Fatalf("import mode: %s", settings.BusinessSwitches.LeadImportMode)
	}
}

func TestValidateCustomFieldSelect(t *testing.T) {
	opts, _ := json.Marshal([]SelectOption{{Value: "a", Label: map[string]string{"zh-CN": "A"}}})
	if err := ValidateCustomFieldInput("select", false, opts, nil); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCustomFieldInput("select", false, nil, nil); err == nil {
		t.Fatal("expected error for missing options")
	}
}
