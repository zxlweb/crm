package http_test

import (
	"context"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type memSegmentRepo struct {
	byTenant map[uuid.UUID][]domain.SegmentTemplate
}

func presetSegments(tenantID uuid.UUID) []domain.SegmentTemplate {
	codes := []string{"high_value", "churn_risk", "new_potential", "needs_activation", "revive_pool"}
	out := make([]domain.SegmentTemplate, len(codes))
	for i, code := range codes {
		out[i] = domain.SegmentTemplate{
			ID:          uuid.New(),
			TenantID:    tenantID,
			Code:        code,
			NameI18nKey: "segments." + code + ".name",
			FilterJSON:  datatypes.JSON([]byte(`{}`)),
			IsSystem:    true,
		}
	}
	return out
}

func (m *memSegmentRepo) ListByTenant(_ context.Context, tenantID uuid.UUID) ([]domain.SegmentTemplate, error) {
	if m.byTenant == nil {
		m.byTenant = map[uuid.UUID][]domain.SegmentTemplate{}
	}
	if _, ok := m.byTenant[tenantID]; !ok {
		m.byTenant[tenantID] = presetSegments(tenantID)
	}
	return m.byTenant[tenantID], nil
}

func (m *memSegmentRepo) GetByCode(_ context.Context, tenantID uuid.UUID, code string) (*domain.SegmentTemplate, error) {
	items, err := m.ListByTenant(context.Background(), tenantID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].Code == code {
			cp := items[i]
			return &cp, nil
		}
	}
	return nil, repository.ErrSegmentNotFound
}

func (m *memSegmentRepo) UpdateFilter(context.Context, uuid.UUID, string, datatypes.JSON) error {
	return nil
}

func leadMatchesSegment(l *domain.Lead, code string, opts crm.SegmentApplyOpts) bool {
	if code == "" {
		return true
	}
	now := time.Now().UTC()
	switch code {
	case "high_value":
		threshold := opts.HighValueAmount
		if threshold <= 0 {
			threshold = 100000
		}
		return l.Amount > threshold
	case "churn_risk":
		days := opts.DaysSilent
		if days < 1 {
			days = 7
		}
		cutoff := now.AddDate(0, 0, -days)
		return l.LastActivityAt == nil || l.LastActivityAt.Before(cutoff)
	case "new_potential":
		since := now.AddDate(0, 0, -7)
		return !l.CreatedAt.Before(since) && l.Status != "qualified"
	case "needs_activation":
		return l.LifecycleStage == "acquire"
	case "revive_pool":
		return l.LifecycleStage == "revive"
	default:
		return false
	}
}

func accountMatchesSegment(a *domain.Account, code string, opts crm.SegmentApplyOpts) bool {
	if code == "" {
		return true
	}
	now := time.Now().UTC()
	switch code {
	case "high_value":
		return a.EngagementScore >= 70
	case "churn_risk":
		days := opts.DaysSilent
		if days < 1 {
			days = 7
		}
		cutoff := now.AddDate(0, 0, -days)
		return a.LastActivityAt == nil || a.LastActivityAt.Before(cutoff)
	case "new_potential":
		since := now.AddDate(0, 0, -7)
		return !a.CreatedAt.Before(since) && a.LifecycleStage == "acquire"
	case "needs_activation":
		return a.LifecycleStage == "acquire"
	case "revive_pool":
		return a.LifecycleStage == "revive"
	default:
		return false
	}
}
