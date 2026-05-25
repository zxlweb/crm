package insights

import (
	"testing"
	"time"

	"crm-backend/internal/domain"

	"github.com/google/uuid"
)

func TestEvaluateLead_INS001_SilentRisk(t *testing.T) {
	leadID := uuid.New()
	old := time.Now().UTC().Add(-10 * 24 * time.Hour)
	resp := EvaluateLead(LeadEvaluateInput{
		LeadID:         leadID,
		Status:         "new",
		LifecycleStage: "acquire",
		CreatedAt:      time.Now().UTC().Add(-30 * 24 * time.Hour),
		LastActivityAt: &old,
		DaysSilent:     7,
	})
	if len(resp.Items) == 0 || resp.Items[0].RuleID != "INS-001" {
		t.Fatalf("expected INS-001, got %+v", resp.Items)
	}
}

func TestEvaluateLead_INS005_LastTwoNegative(t *testing.T) {
	neg := "negative"
	hes := "hesitant"
	now := time.Now().UTC()
	resp := EvaluateLead(LeadEvaluateInput{
		Status:         "qualified",
		LifecycleStage: "grow",
		CreatedAt:      now.Add(-60 * 24 * time.Hour),
		LastActivityAt: &now,
		Activities: []domain.Activity{
			{OccurredAt: now, Sentiment: &neg},
			{OccurredAt: now.Add(-time.Hour), Sentiment: &hes},
		},
		DaysSilent: 7,
	})
	found := false
	for _, it := range resp.Items {
		if it.RuleID == "INS-005" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected INS-005, got %+v", resp.Items)
	}
}

func TestEvaluateAccount_INS001_NoRecentActivity(t *testing.T) {
	resp := EvaluateAccount(uuid.New(), "acquire", 0, nil)
	if len(resp.Items) != 1 || resp.Items[0].RuleID != "INS-001" {
		t.Fatalf("expected INS-001, got %+v", resp.Items)
	}
}

func TestEvaluateAccount_INS001_NoHitWhenRecentActivity(t *testing.T) {
	last := time.Now().UTC()
	resp := EvaluateAccount(uuid.New(), "grow", 50, &last)
	if len(resp.Items) != 0 {
		t.Fatalf("expected no insights, got %+v", resp.Items)
	}
}

func TestEvaluateContact_MatchesAccountBaseline(t *testing.T) {
	id := uuid.New()
	account := EvaluateAccount(id, "", 0, nil)
	contact := EvaluateContact(id, "", 0, nil)
	if len(account.Items) != len(contact.Items) {
		t.Fatalf("item count mismatch: account=%d contact=%d", len(account.Items), len(contact.Items))
	}
}
