package crm

import (
	"encoding/json"
	"regexp"
	"time"
)

var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

type SelectOption struct {
	Value string            `json:"value"`
	Label map[string]string `json:"label"`
}

func ValidateCustomFieldInput(fieldType string, required bool, optionsJSON []byte, defaultVal json.RawMessage) error {
	if !ValidFieldType(fieldType) {
		return ErrCustomFieldTypeInvalid
	}
	switch fieldType {
	case "select":
		if len(optionsJSON) == 0 {
			return ErrCustomFieldTypeInvalid
		}
		var opts []SelectOption
		if err := json.Unmarshal(optionsJSON, &opts); err != nil || len(opts) == 0 {
			return ErrCustomFieldTypeInvalid
		}
		seen := map[string]struct{}{}
		for _, o := range opts {
			if o.Value == "" {
				return ErrCustomFieldTypeInvalid
			}
			if _, dup := seen[o.Value]; dup {
				return ErrCustomFieldTypeInvalid
			}
			seen[o.Value] = struct{}{}
		}
		if len(defaultVal) > 0 {
			var dv string
			if err := json.Unmarshal(defaultVal, &dv); err != nil || !optionContains(opts, dv) {
				return ErrCustomFieldTypeInvalid
			}
		}
	case "date":
		if len(defaultVal) > 0 {
			var dv string
			if err := json.Unmarshal(defaultVal, &dv); err != nil || !datePattern.MatchString(dv) {
				return ErrCustomFieldTypeInvalid
			}
			if _, err := time.Parse("2006-01-02", dv); err != nil {
				return ErrCustomFieldTypeInvalid
			}
		}
	case "text":
		if len(defaultVal) > 0 {
			var dv string
			if err := json.Unmarshal(defaultVal, &dv); err != nil {
				return ErrCustomFieldTypeInvalid
			}
			if len(dv) > 255 {
				return ErrCustomFieldTypeInvalid
			}
		}
	}
	_ = required
	return nil
}

func optionContains(opts []SelectOption, value string) bool {
	for _, o := range opts {
		if o.Value == value {
			return true
		}
	}
	return false
}
