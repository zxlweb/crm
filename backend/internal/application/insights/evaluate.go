package insights

import (
	"time"

	"crm-backend/internal/pkg/crm"

	"github.com/google/uuid"
)

// EvaluateResponse matches phase-2-crm-ai.md §6.3 (rule engine baseline).
type EvaluateResponse struct {
	Items                    []InsightItem `json:"items"`
	EngagementScore          int           `json:"engagement_score"`
	EngagementExplanationKey string        `json:"engagement_explanation_key"`
	LifecycleStage           string        `json:"lifecycle_stage"`
	RelationshipHealth       string        `json:"relationship_health"`
}

type InsightItem struct {
	ID               string           `json:"id"`
	Priority         int              `json:"priority"`
	TitleKey         string           `json:"title_key"`
	BodyKey          string           `json:"body_key"`
	SuggestedAction  *SuggestedAction `json:"suggested_action,omitempty"`
	RuleID           string           `json:"rule_id"`
	ExplanationKey   string           `json:"explanation_key"`
}

type SuggestedAction struct {
	ActivityEventType    string `json:"activity_event_type"`
	ActivityDirection    string `json:"activity_direction"`
	Title                string `json:"title"`
}

// EvaluateAccount runs Phase 2 rule baseline (INS-001 silent risk when no recent activity).
func EvaluateAccount(subjectID uuid.UUID, lifecycle string, engagement int16, lastActivityAt *time.Time) EvaluateResponse {
	health := crm.RelationshipHealthFromScore(engagement)
	items := []InsightItem{}

	if lastActivityAt == nil && engagement < 40 {
		items = append(items, InsightItem{
			ID:             "INS-001",
			Priority:       1,
			TitleKey:       "insight.silent_risk.title",
			BodyKey:        "insight.silent_risk.body",
			RuleID:         "INS-001",
			ExplanationKey: "insight.silent_risk.explanation",
			SuggestedAction: &SuggestedAction{
				ActivityEventType: "call",
				ActivityDirection: "outbound",
				Title:             "今日回访",
			},
		})
	}
	if len(items) > 3 {
		items = items[:3]
	}
	if lifecycle == "" {
		lifecycle = "acquire"
	}
	_ = subjectID
	return EvaluateResponse{
		Items:                    items,
		EngagementScore:          int(engagement),
		EngagementExplanationKey: "engagement.baseline",
		LifecycleStage:           lifecycle,
		RelationshipHealth:       health,
	}
}
