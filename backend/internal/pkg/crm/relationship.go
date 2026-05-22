package crm

// RelationshipHealthFromScore maps engagement_score (0–100) to relationship_health.
func RelationshipHealthFromScore(score int16) string {
	switch {
	case score >= 70:
		return "high"
	case score >= 40:
		return "medium"
	default:
		return "low"
	}
}

// ValidLifecycleStages for accounts / contacts / leads.
var ValidLifecycleStages = map[string]bool{
	"acquire": true, "activate": true, "grow": true, "retain": true, "revive": true,
}

func ValidLifecycleStage(stage string) bool {
	return ValidLifecycleStages[stage]
}
