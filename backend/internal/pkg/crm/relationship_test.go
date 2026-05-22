package crm

import "testing"

func TestRelationshipHealthFromScore(t *testing.T) {
	cases := []struct {
		score int16
		want  string
	}{
		{80, "high"},
		{70, "high"},
		{69, "medium"},
		{40, "medium"},
		{39, "low"},
		{0, "low"},
	}
	for _, tc := range cases {
		if got := RelationshipHealthFromScore(tc.score); got != tc.want {
			t.Fatalf("score %d: got %q want %q", tc.score, got, tc.want)
		}
	}
}
