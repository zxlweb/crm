package repository

import (
	"context"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/datascope"

	"github.com/google/uuid"
)

func (r *GormLeadRepository) CountScoped(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	var n int64
	err := datascope.OwnerScope(r.base(ctx, tenantID), userID, viewAll).Count(&n).Error
	return n, err
}

func (r *GormLeadRepository) DailyCreatedCounts(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID, days int) ([]int64, error) {
	if days < 1 {
		days = 7
	}
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -(days - 1))
	type row struct {
		Day   time.Time
		Count int64
	}
	var rows []row
	q := r.base(ctx, tenantID)
	q = datascope.OwnerScope(q, userID, viewAll)
	err := q.Where("created_at >= ?", start).
		Select("date_trunc('day', created_at) AS day, COUNT(*) AS count").
		Group("day").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	byDay := make(map[string]int64, len(rows))
	for _, row := range rows {
		byDay[row.Day.Format("2006-01-02")] = row.Count
	}
	out := make([]int64, days)
	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i).Format("2006-01-02")
		out[i] = byDay[d]
	}
	return out, nil
}

func (r *GormLeadRepository) CountLowEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	var n int64
	q := r.base(ctx, tenantID).Where("engagement_score < ?", 40)
	q = datascope.OwnerScope(q, userID, viewAll)
	err := q.Count(&n).Error
	return n, err
}

func (r *GormLeadRepository) AvgEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (float64, error) {
	var avg *float64
	q := r.base(ctx, tenantID)
	q = datascope.OwnerScope(q, userID, viewAll)
	err := q.Select("COALESCE(AVG(engagement_score), 0)").Scan(&avg).Error
	if err != nil || avg == nil {
		return 0, err
	}
	return *avg, nil
}

func (r *GormLeadRepository) ListPriorityCandidates(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID, limit int) ([]domain.Lead, error) {
	if limit < 1 {
		limit = 20
	}
	q := r.base(ctx, tenantID)
	q = datascope.OwnerScope(q, userID, viewAll)
	var items []domain.Lead
	err := q.Where("status <> ?", "converted").
		Order("engagement_score ASC, updated_at ASC").
		Limit(limit).
		Find(&items).Error
	return items, err
}
