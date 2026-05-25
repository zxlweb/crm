package repository

import (
	"context"
	"fmt"
	"time"

	"crm-backend/internal/pkg/crm"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeadStatsFilter struct {
	From    *time.Time
	To      *time.Time
	ViewAll bool
	UserID  uuid.UUID
}

type LabelCount struct {
	Label string
	Count int64
}

type TrendPoint struct {
	Date  string
	Count int64
}

func (r *GormLeadRepository) scopedStats(ctx context.Context, tenantID uuid.UUID, f LeadStatsFilter) *gorm.DB {
	q := r.base(ctx, tenantID)
	if !f.ViewAll {
		q = q.Where("(owner_id = ? OR owner_id IS NULL)", f.UserID)
	}
	if f.From != nil {
		q = q.Where("created_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("created_at < ?", *f.To)
	}
	return q
}

func (r *GormLeadRepository) StatsBySource(ctx context.Context, tenantID uuid.UUID, f LeadStatsFilter) ([]LabelCount, int64, error) {
	return r.statsByColumn(ctx, tenantID, f, "COALESCE(NULLIF(source, ''), 'unknown')")
}

func (r *GormLeadRepository) StatsByStatus(ctx context.Context, tenantID uuid.UUID, f LeadStatsFilter) ([]LabelCount, int64, error) {
	return r.statsByColumn(ctx, tenantID, f, "status")
}

func (r *GormLeadRepository) statsByColumn(ctx context.Context, tenantID uuid.UUID, f LeadStatsFilter, expr string) ([]LabelCount, int64, error) {
	type row struct {
		Label string
		Count int64
	}
	var rows []row
	err := r.scopedStats(ctx, tenantID, f).
		Select(fmt.Sprintf("%s AS label, COUNT(*) AS count", expr)).
		Group("label").
		Order("count DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	out := make([]LabelCount, len(rows))
	for i, row := range rows {
		out[i] = LabelCount{Label: row.Label, Count: row.Count}
		total += row.Count
	}
	return out, total, nil
}

func (r *GormLeadRepository) StatsTrend(ctx context.Context, tenantID uuid.UUID, f LeadStatsFilter, granularity string) ([]TrendPoint, error) {
	trunc := "day"
	if granularity == "week" {
		trunc = "week"
	}
	type row struct {
		Bucket time.Time
		Count  int64
	}
	var rows []row
	err := r.scopedStats(ctx, tenantID, f).
		Select(fmt.Sprintf("date_trunc('%s', created_at) AS bucket, COUNT(*) AS count", trunc)).
		Group("bucket").
		Order("bucket ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]TrendPoint, len(rows))
	for i, row := range rows {
		out[i] = TrendPoint{
			Date:  row.Bucket.Format("2006-01-02"),
			Count: row.Count,
		}
	}
	return out, nil
}

func (r *GormLeadRepository) StatsFunnel(ctx context.Context, tenantID uuid.UUID, f LeadStatsFilter) ([]LabelCount, error) {
	counts, _, err := r.statsByColumn(ctx, tenantID, f, "status")
	if err != nil {
		return nil, err
	}
	byStatus := make(map[string]int64, len(counts))
	for _, c := range counts {
		byStatus[c.Label] = c.Count
	}
	order := []string{"new", "contacted", "qualified", "unqualified", "converted"}
	out := make([]LabelCount, 0, len(order))
	for _, status := range order {
		if !crm.ValidLeadStatus(status) {
			continue
		}
		out = append(out, LabelCount{Label: status, Count: byStatus[status]})
	}
	return out, nil
}
