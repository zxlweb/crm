package insights

import (
	"sort"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/crm"

	"github.com/google/uuid"
)

// EvaluateResponse matches phase-2-crm-ai.md §6.3.
type EvaluateResponse struct {
	Items                    []InsightItem `json:"items"`
	EngagementScore          int           `json:"engagement_score"`
	EngagementExplanationKey string        `json:"engagement_explanation_key"`
	LifecycleStage           string        `json:"lifecycle_stage"`
	RelationshipHealth       string        `json:"relationship_health"`
}

type InsightItem struct {
	ID              string           `json:"id"`
	Priority        int              `json:"priority"`
	TitleKey        string           `json:"title_key"`
	BodyKey         string           `json:"body_key"`
	SuggestedAction *SuggestedAction `json:"suggested_action,omitempty"`
	RuleID          string           `json:"rule_id"`
	ExplanationKey  string           `json:"explanation_key"`
}

type SuggestedAction struct {
	ActivityEventType string `json:"activity_event_type"`
	ActivityDirection string `json:"activity_direction"`
	Title             string `json:"title"`
}

// LeadEvaluateInput bundles lead + activities for INS-001～006.
type LeadEvaluateInput struct {
	LeadID         uuid.UUID
	Status         string
	LifecycleStage string
	CreatedAt      time.Time
	LastActivityAt *time.Time
	Activities     []domain.Activity
	DaysSilent     int
}

// EvaluateLead runs INS-001～006 and computes engagement score.
func EvaluateLead(in LeadEvaluateInput) EvaluateResponse {
	lifecycle := in.LifecycleStage
	if lifecycle == "" {
		lifecycle = "acquire"
	}
	daysSilent := in.DaysSilent
	if daysSilent < 1 {
		daysSilent = 7
	}
	acts := sortedActivities(in.Activities)
	count7d := countActivitiesSince(acts, time.Now().UTC().AddDate(0, 0, -7))
	score := crm.ComputeEngagementScore(in.Status, lifecycle, count7d)
	health := crm.RelationshipHealthFromScore(score)

	items := evaluateLeadRules(in, acts, daysSilent)
	if len(items) > 3 {
		items = items[:3]
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Priority < items[j].Priority })

	explanation := "engagement.last_7_days"
	if count7d == 0 {
		explanation = "engagement.baseline"
	}
	_ = in.LeadID
	return EvaluateResponse{
		Items:                    items,
		EngagementScore:          int(score),
		EngagementExplanationKey: explanation,
		LifecycleStage:           lifecycle,
		RelationshipHealth:       health,
	}
}

func evaluateLeadRules(in LeadEvaluateInput, acts []domain.Activity, daysSilent int) []InsightItem {
	var items []InsightItem
	now := time.Now().UTC()
	daysSinceActivity := crm.DaysSince(in.LastActivityAt)

	// INS-001: days_since_last_activity > N
	if daysSinceActivity > daysSilent {
		items = append(items, insightSilentRisk())
	}

	// INS-002: qualified + 7 days no follow-up
	if in.Status == "qualified" && daysSinceActivity > daysSilent {
		items = append(items, InsightItem{
			ID:             "INS-002",
			Priority:       2,
			TitleKey:       "insight.qualified_cooling.title",
			BodyKey:        "insight.qualified_cooling.body",
			RuleID:         "INS-002",
			ExplanationKey: "insight.qualified_cooling.explanation",
			SuggestedAction: &SuggestedAction{
				ActivityEventType: "call",
				ActivityDirection: "outbound",
				Title:             "经理介入回访",
			},
		})
	}

	// INS-003: lifecycle = revive
	if in.LifecycleStage == "revive" {
		items = append(items, InsightItem{
			ID:             "INS-003",
			Priority:       3,
			TitleKey:       "insight.revive_window.title",
			BodyKey:        "insight.revive_window.body",
			RuleID:         "INS-003",
			ExplanationKey: "insight.revive_window.explanation",
			SuggestedAction: &SuggestedAction{
				ActivityEventType: "email",
				ActivityDirection: "outbound",
				Title:             "发送唤醒权益",
			},
		})
	}

	// INS-004: created > 3 days ago, no first touch within 3 days
	activationDeadline := in.CreatedAt.UTC().Add(3 * 24 * time.Hour)
	if now.After(activationDeadline) && !hasActivityBefore(acts, activationDeadline) {
		items = append(items, InsightItem{
			ID:             "INS-004",
			Priority:       4,
			TitleKey:       "insight.activation_failed.title",
			BodyKey:        "insight.activation_failed.body",
			RuleID:         "INS-004",
			ExplanationKey: "insight.activation_failed.explanation",
			SuggestedAction: &SuggestedAction{
				ActivityEventType: "call",
				ActivityDirection: "outbound",
				Title:             "分配复核",
			},
		})
	}

	// INS-005: last 2 activities sentiment in {negative, hesitant}
	if lastNSentimentsMatch(acts, 2, "negative", "hesitant") {
		items = append(items, InsightItem{
			ID:             "INS-005",
			Priority:       5,
			TitleKey:       "insight.sentiment_cooling.title",
			BodyKey:        "insight.sentiment_cooling.body",
			RuleID:         "INS-005",
			ExplanationKey: "insight.sentiment_cooling.explanation",
			SuggestedAction: &SuggestedAction{
				ActivityEventType: "call",
				ActivityDirection: "outbound",
				Title:             "共情式回访",
			},
		})
	}

	// INS-006: days_since_positive > 30 && lifecycle = grow
	if in.LifecycleStage == "grow" && daysSinceSentiment(acts, "positive") > 30 {
		items = append(items, InsightItem{
			ID:             "INS-006",
			Priority:       6,
			TitleKey:       "insight.grow_enthusiasm_fade.title",
			BodyKey:        "insight.grow_enthusiasm_fade.body",
			RuleID:         "INS-006",
			ExplanationKey: "insight.grow_enthusiasm_fade.explanation",
			SuggestedAction: &SuggestedAction{
				ActivityEventType: "meeting",
				ActivityDirection: "outbound",
				Title:             "推进下一阶段",
			},
		})
	}

	return items
}

func insightSilentRisk() InsightItem {
	return InsightItem{
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
	}
}

// EvaluateAccount runs Phase 2 rule baseline for accounts.
func EvaluateAccount(subjectID uuid.UUID, lifecycle string, engagement int16, lastActivityAt *time.Time) EvaluateResponse {
	health := crm.RelationshipHealthFromScore(engagement)
	items := []InsightItem{}
	if crm.DaysSince(lastActivityAt) > 7 {
		items = append(items, insightSilentRisk())
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

// EvaluateContact runs the same Phase 2 rule baseline for contacts.
func EvaluateContact(subjectID uuid.UUID, lifecycle string, engagement int16, lastActivityAt *time.Time) EvaluateResponse {
	return EvaluateAccount(subjectID, lifecycle, engagement, lastActivityAt)
}

func sortedActivities(acts []domain.Activity) []domain.Activity {
	out := make([]domain.Activity, len(acts))
	copy(out, acts)
	sort.Slice(out, func(i, j int) bool {
		return out[i].OccurredAt.After(out[j].OccurredAt)
	})
	return out
}

func countActivitiesSince(acts []domain.Activity, since time.Time) int {
	n := 0
	for _, a := range acts {
		if !a.OccurredAt.Before(since) {
			n++
		}
	}
	return n
}

func hasActivityBefore(acts []domain.Activity, deadline time.Time) bool {
	for _, a := range acts {
		if a.OccurredAt.Before(deadline) || a.OccurredAt.Equal(deadline) {
			return true
		}
	}
	return false
}

func lastNSentimentsMatch(acts []domain.Activity, n int, sentiments ...string) bool {
	if len(acts) < n {
		return false
	}
	allowed := map[string]bool{}
	for _, s := range sentiments {
		allowed[s] = true
	}
	sorted := sortedActivities(acts)
	for i := 0; i < n; i++ {
		if sorted[i].Sentiment == nil || !allowed[*sorted[i].Sentiment] {
			return false
		}
	}
	return true
}

func daysSinceSentiment(acts []domain.Activity, target string) int {
	sorted := sortedActivities(acts)
	for _, a := range sorted {
		if a.Sentiment != nil && *a.Sentiment == target {
			return crm.DaysSince(&a.OccurredAt)
		}
	}
	return 9999
}
