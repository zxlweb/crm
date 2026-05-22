package crm

import "testing"

func TestCanTransitionLeadStatus(t *testing.T) {
	ok := []struct{ from, to string }{
		{"new", "contacted"},
		{"contacted", "qualified"},
		{"contacted", "unqualified"},
		{"qualified", "unqualified"},
		{"new", "new"},
	}
	for _, tc := range ok {
		if !CanTransitionLeadStatus(tc.from, tc.to) {
			t.Fatalf("%s -> %s should be allowed", tc.from, tc.to)
		}
	}
	bad := []struct{ from, to string }{
		{"qualified", "new"},
		{"qualified", "converted"},
		{"new", "qualified"},
		{"converted", "contacted"},
		{"unqualified", "qualified"},
	}
	for _, tc := range bad {
		if CanTransitionLeadStatus(tc.from, tc.to) {
			t.Fatalf("%s -> %s should be denied", tc.from, tc.to)
		}
	}
}

func TestCanConvertLead(t *testing.T) {
	if !CanConvertLead("qualified") {
		t.Fatal("qualified should be convertible")
	}
	for _, s := range []string{"new", "contacted", "converted", "unqualified"} {
		if CanConvertLead(s) {
			t.Fatalf("%s should not be convertible", s)
		}
	}
}
