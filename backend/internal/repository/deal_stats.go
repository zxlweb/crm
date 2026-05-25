package repository

import (
	"context"
	"fmt"
	"time"

	"crm-backend/internal/pkg/crm"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DealStatsFilter struct {
	From    *time.Time
	To      *time.Time
	ViewAll bool
	UserID  uuid.UUID
}

type DealStageStat struct {
	Label  string
	Value  int64
	Amount float64
}

type DealWinRatePoint struct {
	Period string
	Won    int64
	Lost   int64
	Rate   float64
}

type DealOwnerMetric struct {
	OwnerID uuid.UUID
	Value   float64
}

func (r *GormDealRepository) scopedStats(ctx context.Context, tenantID uuid.UUID, f DealStatsFilter) *gorm.DB {
	q := r.scoped(ctx, tenantID, f.ViewAll, f.UserID)
	if f.From != nil {
		q = q.Where("created_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("created_at < ?", *f.To)
	}
	return q
}

func (r *GormDealRepository) CountScoped(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	var n int64
	err := r.scoped(ctx, tenantID, viewAll, userID).Count(&n).Error
	return n, err
}

func (r *GormDealRepository) StatsByStage(ctx context.Context, tenantID uuid.UUID, f DealStatsFilter, metric string) ([]DealStageStat, int64, error) {
	type row struct {
		Label  string
		Value  int64
		Amount float64
	}
	q := r.scopedStats(ctx, tenantID, f)
	var rows []row
	selectExpr := "stage AS label, COUNT(*) AS value, COALESCE(SUM(amount), 0) AS amount"
	if metric == "amount" {
		selectExpr = "stage AS label, COALESCE(SUM(amount), 0) AS value, COALESCE(SUM(amount), 0) AS amount"
	}
	err := q.Select(selectExpr).Group("stage").Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	byStage := make(map[string]row, len(rows))
	for _, row := range rows {
		byStage[row.Label] = row
	}
	out := make([]DealStageStat, 0, len(crm.DealPipelineStages))
	var total int64
	for _, stage := range crm.DealPipelineStages {
		row := byStage[stage]
		out = append(out, DealStageStat{Label: stage, Value: row.Value, Amount: row.Amount})
		if metric == "amount" {
			total += int64(row.Amount)
		} else {
			total += row.Value
		}
	}
	return out, total, nil
}

func (r *GormDealRepository) StatsWinRate(ctx context.Context, tenantID uuid.UUID, f DealStatsFilter, granularity string) ([]DealWinRatePoint, error) {
	trunc := "week"
	if granularity == "month" {
		trunc = "month"
	}
	type row struct {
		Bucket time.Time
		Stage  string
		Count  int64
	}
	var rows []row
	err := r.scopedStats(ctx, tenantID, f).
		Where("stage IN ?", []string{crm.DealStageWon, crm.DealStageLost}).
		Where("closed_at IS NOT NULL").
		Select(fmt.Sprintf("date_trunc('%s', closed_at) AS bucket, stage, COUNT(*) AS count", trunc)).
		Group("bucket, stage").
		Order("bucket ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	type agg struct {
		won, lost int64
	}
	byBucket := map[string]*agg{}
	for _, row := range rows {
		key := formatWinRatePeriod(row.Bucket, granularity)
		if byBucket[key] == nil {
			byBucket[key] = &agg{}
		}
		if row.Stage == crm.DealStageWon {
			byBucket[key].won += row.Count
		} else {
			byBucket[key].lost += row.Count
		}
	}
	out := make([]DealWinRatePoint, 0, len(byBucket))
	for period, a := range byBucket {
		rate := 0.0
		if a.won+a.lost > 0 {
			rate = float64(a.won) / float64(a.won+a.lost)
		}
		out = append(out, DealWinRatePoint{Period: period, Won: a.won, Lost: a.lost, Rate: rate})
	}
	// stable order by period string
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Period < out[i].Period {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out, nil
}

func (r *GormDealRepository) DailyCreatedCounts(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID, days int) ([]int64, error) {
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
	err := r.scoped(ctx, tenantID, viewAll, userID).
		Where("created_at >= ?", start).
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

func (r *GormDealRepository) CountByStage(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) ([]LabelCount, error) {
	f := DealStatsFilter{ViewAll: viewAll, UserID: userID}
	rows, _, err := r.StatsByStage(ctx, tenantID, f, "count")
	if err != nil {
		return nil, err
	}
	out := make([]LabelCount, len(rows))
	for i, row := range rows {
		out[i] = LabelCount{Label: row.Label, Count: row.Value}
	}
	return out, nil
}

func (r *GormDealRepository) TeamRanking(ctx context.Context, tenantID uuid.UUID, metric string, limit int) ([]DealOwnerMetric, error) {
	if limit < 1 {
		limit = 10
	}
	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	var rows []DealOwnerMetric
	switch metric {
	case "win_count":
		err := r.base(ctx, tenantID).
			Where("stage = ? AND closed_at >= ?", crm.DealStageWon, monthStart).
			Where("owner_id IS NOT NULL").
			Select("owner_id, COUNT(*)::float8 AS value").
			Group("owner_id").
			Order("value DESC").
			Limit(limit).
			Scan(&rows).Error
		return rows, err
	case "avg_engagement":
		err := r.base(ctx, tenantID).
			Where("owner_id IS NOT NULL").
			Select("owner_id, COALESCE(AVG(engagement_score), 0) AS value").
			Group("owner_id").
			Order("value DESC").
			Limit(limit).
			Scan(&rows).Error
		return rows, err
	default: // won_amount
		err := r.base(ctx, tenantID).
			Where("stage = ? AND closed_at >= ?", crm.DealStageWon, monthStart).
			Where("owner_id IS NOT NULL").
			Select("owner_id, COALESCE(SUM(amount), 0) AS value").
			Group("owner_id").
			Order("value DESC").
			Limit(limit).
			Scan(&rows).Error
		return rows, err
	}
}

func formatWinRatePeriod(t time.Time, granularity string) string {
	if granularity == "month" {
		return t.Format("2006-01")
	}
	y, w := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", y, w)
}
