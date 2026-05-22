package domain

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// StringArray maps PostgreSQL text[].
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil || len(a) == 0 {
		return "{}", nil
	}
	return "{" + strings.Join(quoteArrayElems([]string(a)), ",") + "}", nil
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}
	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}
	*a = parseTextArray(s)
	return nil
}

func parseTextArray(s string) StringArray {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "{}")
	if s == "" {
		return StringArray{}
	}
	parts := strings.Split(s, ",")
	out := make(StringArray, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, `"`)
		out = append(out, p)
	}
	return out
}

func quoteArrayElems(elems []string) []string {
	out := make([]string, len(elems))
	for i, e := range elems {
		escaped := strings.ReplaceAll(e, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		out[i] = `"` + escaped + `"`
	}
	return out
}
