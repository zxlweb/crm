package deal

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
	Metric      string
}

type StageStatDTO struct {
	Label  string  `json:"label"`
	Value  int64   `json:"value"`
	Amount float64 `json:"amount,omitempty"`
}

type ByStageDTO struct {
	Items []StageStatDTO `json:"items"`
	Total int64          `json:"total"`
}

type WinRateItemDTO struct {
	Period string  `json:"period"`
	Won    int64   `json:"won"`
	Lost   int64   `json:"lost"`
	Rate   float64 `json:"rate"`
}

type WinRateDTO struct {
	Items []WinRateItemDTO `json:"items"`
}

func (s *Service) statsFilter(ctx context.Context, tenantID, userID uuid.UUID, q StatsQuery) repository.DealStatsFilter {
	return repository.DealStatsFilter{
		From:    q.From,
		To:      q.To,
		ViewAll: s.viewAll(userID.String(), tenantID.String()),
		UserID:  userID,
	}
}

func (s *Service) StatsByStage(ctx context.Context, tenantID, userID uuid.UUID, q StatsQuery) (*ByStageDTO, error) {
	metric := q.Metric
	if metric == "" {
		metric = "count"
	}
	rows, total, err := s.repo.StatsByStage(ctx, tenantID, s.statsFilter(ctx, tenantID, userID, q), metric)
	if err != nil {
		return nil, err
	}
	items := make([]StageStatDTO, len(rows))
	for i, row := range rows {
		items[i] = StageStatDTO{Label: row.Label, Value: row.Value, Amount: row.Amount}
	}
	return &ByStageDTO{Items: items, Total: total}, nil
}

func (s *Service) StatsWinRate(ctx context.Context, tenantID, userID uuid.UUID, q StatsQuery) (*WinRateDTO, error) {
	granularity := q.Granularity
	if granularity == "" {
		granularity = "week"
	}
	rows, err := s.repo.StatsWinRate(ctx, tenantID, s.statsFilter(ctx, tenantID, userID, q), granularity)
	if err != nil {
		return nil, err
	}
	items := make([]WinRateItemDTO, len(rows))
	for i, row := range rows {
		items[i] = WinRateItemDTO{Period: row.Period, Won: row.Won, Lost: row.Lost, Rate: row.Rate}
	}
	return &WinRateDTO{Items: items}, nil
}
