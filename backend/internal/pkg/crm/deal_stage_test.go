package crm

import "testing"

func TestCanTransitionDealStage(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want bool
	}{
		{DealStageQualification, DealStageProposal, true},
		{DealStageProposal, DealStageQualification, true},
		{DealStageProposal, DealStageNegotiation, true},
		{DealStageNegotiation, DealStageProposal, true},
		{DealStageNegotiation, DealStageWon, true},
		{DealStageNegotiation, DealStageLost, true},
		{DealStageQualification, DealStageNegotiation, false},
		{DealStageQualification, DealStageWon, false},
		{DealStageProposal, DealStageWon, false},
		{DealStageWon, DealStageProposal, false},
		{DealStageLost, DealStageQualification, false},
		{DealStageWon, DealStageWon, true},
	}
	for _, tt := range tests {
		got := CanTransitionDealStage(tt.from, tt.to)
		if got != tt.want {
			t.Errorf("CanTransitionDealStage(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
		}
	}
}
