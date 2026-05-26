package lead

import (
	"context"
	"time"

	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

type StatsQuery struct {
	From        *time.Time
	To          *time.Time
	Granularity string
}

type LabelCountDTO struct {
	Label      string  `json:"label"`
	Value      int64   `json:"value"`
	Percentage float64 `json:"percentage,omitempty"`
}

type DistributionDTO struct {
	Items []LabelCountDTO `json:"items"`
	Total int64           `json:"total"`
}

type TrendDTO struct {
	Categories []string         `json:"categories"`
	Series     []TrendSeriesDTO `json:"series"`
}

type TrendSeriesDTO struct {
	Name    string  `json:"name"`
	Data    []int64 `json:"data"`
	Primary bool    `json:"primary,omitempty"`
}

type FunnelStageDTO struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type ConversionRateDTO struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

type FunnelDTO struct {
	Stages          []FunnelStageDTO    `json:"stages"`
	ConversionRates []ConversionRateDTO `json:"conversion_rates"`
}

func (s *Service) statsFilter(ctx context.Context, tenantID, userID uuid.UUID, q StatsQuery) repository.LeadStatsFilter {
	return repository.LeadStatsFilter{
		From:    q.From,
		To:      q.To,
		ViewAll: s.viewAll(ctx, userID.String(), tenantID.String()),
		UserID:  userID,
	}
}

func (s *Service) StatsBySource(ctx context.Context, tenantID, userID uuid.UUID, q StatsQuery) (*DistributionDTO, error) {
	rows, total, err := s.repo.StatsBySource(ctx, tenantID, s.statsFilter(ctx, tenantID, userID, q))
	if err != nil {
		return nil, err
	}
	return toDistribution(rows, total), nil
}

func (s *Service) StatsByStatus(ctx context.Context, tenantID, userID uuid.UUID, q StatsQuery) (*DistributionDTO, error) {
	rows, total, err := s.repo.StatsByStatus(ctx, tenantID, s.statsFilter(ctx, tenantID, userID, q))
	if err != nil {
		return nil, err
	}
	return toDistribution(rows, total), nil
}

func (s *Service) StatsTrend(ctx context.Context, tenantID, userID uuid.UUID, q StatsQuery) (*TrendDTO, error) {
	granularity := q.Granularity
	if granularity == "" {
		granularity = "day"
	}
	points, err := s.repo.StatsTrend(ctx, tenantID, s.statsFilter(ctx, tenantID, userID, q), granularity)
	if err != nil {
		return nil, err
	}
	categories := make([]string, len(points))
	data := make([]int64, len(points))
	for i, p := range points {
		categories[i] = p.Date
		data[i] = p.Count
	}
	return &TrendDTO{
		Categories: categories,
		Series: []TrendSeriesDTO{
			{Name: "new_leads", Data: data, Primary: true},
		},
	}, nil
}

func (s *Service) StatsFunnel(ctx context.Context, tenantID, userID uuid.UUID, q StatsQuery) (*FunnelDTO, error) {
	rows, err := s.repo.StatsFunnel(ctx, tenantID, s.statsFilter(ctx, tenantID, userID, q))
	if err != nil {
		return nil, err
	}
	stages := make([]FunnelStageDTO, len(rows))
	for i, row := range rows {
		stages[i] = FunnelStageDTO{Name: row.Label, Count: row.Count}
	}
	rates := make([]ConversionRateDTO, 0, len(stages)-1)
	for i := 0; i < len(stages)-1; i++ {
		from := stages[i]
		to := stages[i+1]
		rate := 0.0
		if from.Count > 0 {
			rate = float64(to.Count) / float64(from.Count)
		}
		rates = append(rates, ConversionRateDTO{
			From: from.Name,
			To:   to.Name,
			Rate: rate,
		})
	}
	return &FunnelDTO{Stages: stages, ConversionRates: rates}, nil
}

func toDistribution(rows []repository.LabelCount, total int64) *DistributionDTO {
	items := make([]LabelCountDTO, len(rows))
	for i, row := range rows {
		pct := 0.0
		if total > 0 {
			pct = float64(row.Count) / float64(total)
		}
		items[i] = LabelCountDTO{
			Label:      row.Label,
			Value:      row.Count,
			Percentage: pct,
		}
	}
	return &DistributionDTO{Items: items, Total: total}
}
