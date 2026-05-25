package crm

const (
	DealStageQualification = "qualification"
	DealStageProposal      = "proposal"
	DealStageNegotiation   = "negotiation"
	DealStageWon           = "won"
	DealStageLost          = "lost"
)

var dealStageOrder = map[string]int{
	DealStageQualification: 0,
	DealStageProposal:      1,
	DealStageNegotiation:   2,
	DealStageWon:           3,
	DealStageLost:          3,
}

// DealPipelineStages is the canonical order for pipeline / funnel responses.
var DealPipelineStages = []string{
	DealStageQualification,
	DealStageProposal,
	DealStageNegotiation,
	DealStageWon,
	DealStageLost,
}

// ValidDealStage reports whether stage is a known deal stage code.
func ValidDealStage(stage string) bool {
	_, ok := dealStageOrder[stage]
	return ok
}

// ValidDealCurrency reports supported currency codes (Phase 3).
func ValidDealCurrency(currency string) bool {
	return currency == "" || currency == "CNY" || currency == "USD"
}

// IsDealTerminal reports won/lost stages that cannot change again.
func IsDealTerminal(stage string) bool {
	return stage == DealStageWon || stage == DealStageLost
}

// IsDealOpen reports non-terminal pipeline stages.
func IsDealOpen(stage string) bool {
	return stage == DealStageQualification || stage == DealStageProposal || stage == DealStageNegotiation
}

// CanTransitionDealStage validates pipeline moves per phase-3-deals-dashboard-api.md §2.1.
func CanTransitionDealStage(from, to string) bool {
	if from == to {
		return true
	}
	if !ValidDealStage(from) || !ValidDealStage(to) {
		return false
	}
	if IsDealTerminal(from) {
		return false
	}
	if to == DealStageWon || to == DealStageLost {
		return from == DealStageNegotiation
	}
	fromIdx, okFrom := dealStageOrder[from]
	toIdx, okTo := dealStageOrder[to]
	if !okFrom || !okTo {
		return false
	}
	if toIdx < fromIdx {
		return true
	}
	return toIdx == fromIdx+1
}
