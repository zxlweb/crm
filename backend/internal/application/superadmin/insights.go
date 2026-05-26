package superadmin

import (
	"context"
	"time"
)

type TenantHealthDTO struct {
	Dimensions []string              `json:"dimensions"`
	Items      []TenantHealthItemDTO `json:"items"`
}

type TenantHealthItemDTO struct {
	TenantID     string         `json:"tenant_id"`
	TenantName   string         `json:"tenant_name"`
	Scores       map[string]int `json:"scores"`
	OverallScore int            `json:"overall_score"`
}

type PlanDistributionDTO struct {
	Items []PlanSliceDTO `json:"items"`
	Total int64          `json:"total"`
}

type PlanSliceDTO struct {
	Plan  string `json:"plan"`
	Count int64  `json:"count"`
}

type TopTenantsDTO struct {
	Metric string           `json:"metric"`
	Items  []TopTenantItem  `json:"items"`
}

type TopTenantItem struct {
	TenantID   string  `json:"tenant_id"`
	TenantName string  `json:"tenant_name"`
	Value      float64 `json:"value"`
}

func (s *Service) TenantHealth(ctx context.Context) (*TenantHealthDTO, error) {
	if s.insights == nil {
		return &TenantHealthDTO{
			Dimensions: []string{"activity", "config_completeness", "audit_risk", "data_freshness", "feature_adoption"},
			Items:      []TenantHealthItemDTO{},
		}, nil
	}
	rows, err := s.insights.TenantHealth(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]TenantHealthItemDTO, 0, len(rows))
	for _, row := range rows {
		items = append(items, TenantHealthItemDTO{
			TenantID:     row.TenantID.String(),
			TenantName:   row.TenantName,
			Scores:       row.Scores,
			OverallScore: row.Overall,
		})
	}
	return &TenantHealthDTO{
		Dimensions: []string{"activity", "config_completeness", "audit_risk", "data_freshness", "feature_adoption"},
		Items:      items,
	}, nil
}

func (s *Service) PlanDistribution(ctx context.Context, from, to time.Time) (*PlanDistributionDTO, error) {
	if s.insights == nil {
		return &PlanDistributionDTO{Items: []PlanSliceDTO{}}, nil
	}
	rows, err := s.insights.PlanDistribution(ctx, from, to)
	if err != nil {
		return nil, err
	}
	var total int64
	items := make([]PlanSliceDTO, 0, len(rows))
	for _, row := range rows {
		total += row.Count
		items = append(items, PlanSliceDTO{Plan: row.Plan, Count: row.Count})
	}
	return &PlanDistributionDTO{Items: items, Total: total}, nil
}

func (s *Service) TopTenants(ctx context.Context, metric string, limit int) (*TopTenantsDTO, error) {
	if metric == "" {
		metric = "activity"
	}
	if s.insights == nil {
		return &TopTenantsDTO{Metric: metric, Items: []TopTenantItem{}}, nil
	}
	rows, err := s.insights.TopTenants(ctx, metric, limit)
	if err != nil {
		return nil, err
	}
	items := make([]TopTenantItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, TopTenantItem{
			TenantID:   row.TenantID.String(),
			TenantName: row.TenantName,
			Value:      row.Value,
		})
	}
	return &TopTenantsDTO{Metric: metric, Items: items}, nil
}
