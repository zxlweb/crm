package auditstats

import (
	"context"
	"errors"
	"sync"
	"time"

	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

var ErrExportRateLimited = errors.New("audit_export_rate_limited")

type Service struct {
	repo      repository.AuditStatsRepository
	exportMu  sync.Mutex
	exportHit map[string]time.Time
}

func NewService(repo repository.AuditStatsRepository) *Service {
	return &Service{
		repo:      repo,
		exportHit: map[string]time.Time{},
	}
}

type ByActionDTO struct {
	Items []repository.ActionCountRow `json:"items"`
	Total int64                       `json:"total"`
}

type TrendDTO struct {
	Items []repository.TrendRow `json:"items"`
}

type TopActorsDTO struct {
	Items []repository.ActorCountRow `json:"items"`
}

func (s *Service) ByAction(ctx context.Context, tenantID uuid.UUID, f repository.AuditStatsFilter) (*ByActionDTO, error) {
	rows, total, err := s.repo.CountByAction(ctx, tenantID, f)
	if err != nil {
		return nil, err
	}
	return &ByActionDTO{Items: rows, Total: total}, nil
}

func (s *Service) Trend(ctx context.Context, tenantID uuid.UUID, f repository.AuditStatsFilter, granularity string) (*TrendDTO, error) {
	if granularity == "" {
		granularity = "day"
	}
	rows, err := s.repo.Trend(ctx, tenantID, f, granularity)
	if err != nil {
		return nil, err
	}
	return &TrendDTO{Items: rows}, nil
}

func (s *Service) TopActors(ctx context.Context, tenantID uuid.UUID, f repository.AuditStatsFilter) (*TopActorsDTO, error) {
	rows, err := s.repo.TopActors(ctx, tenantID, f)
	if err != nil {
		return nil, err
	}
	return &TopActorsDTO{Items: rows}, nil
}

func (s *Service) ExportCSV(ctx context.Context, tenantID, userID uuid.UUID, f repository.AuditStatsFilter) (string, error) {
	key := tenantID.String() + ":" + userID.String()
	now := time.Now()
	s.exportMu.Lock()
	if last, ok := s.exportHit[key]; ok && now.Sub(last) < time.Minute {
		s.exportMu.Unlock()
		return "", ErrExportRateLimited
	}
	s.exportHit[key] = now
	s.exportMu.Unlock()

	rows, err := s.repo.ExportRows(ctx, tenantID, f)
	if err != nil {
		return "", err
	}
	return repository.FormatAuditCSV(rows), nil
}
