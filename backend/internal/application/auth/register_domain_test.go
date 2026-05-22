package auth

import "testing"

func TestNormalizeDomain(t *testing.T) {
	if got := normalizeDomain("acme-corp", ""); got != "acme-corp" {
		t.Fatalf("explicit: %s", got)
	}
	if got := normalizeDomain("", "Acme Corp"); got != "acme-corp" {
		t.Fatalf("from company: %s", got)
	}
}
