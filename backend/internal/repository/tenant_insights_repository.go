package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantHealthRow struct {
	TenantID   uuid.UUID
	TenantName string
	Scores     map[string]int
	Overall    int
}

type PlanCountRow struct {
	Plan  string
	Count int64
}

type TopTenantRow struct {
	TenantID   uuid.UUID
	TenantName string
	Value      float64
}

type TenantInsightsRepository interface {
	TenantHealth(ctx context.Context) ([]TenantHealthRow, error)
	PlanDistribution(ctx context.Context, from, to time.Time) ([]PlanCountRow, error)
	TopTenants(ctx context.Context, metric string, limit int) ([]TopTenantRow, error)
}

type GormTenantInsightsRepository struct {
	db *gorm.DB
}

func NewTenantInsightsRepository(db *gorm.DB) TenantInsightsRepository {
	return &GormTenantInsightsRepository{db: db}
}

func (r *GormTenantInsightsRepository) TenantHealth(ctx context.Context) ([]TenantHealthRow, error) {
	type raw struct {
		TenantID          uuid.UUID `gorm:"column:tenant_id"`
		TenantName        string    `gorm:"column:tenant_name"`
		ActivityCount     int64     `gorm:"column:activity_count"`
		HasLocale         bool      `gorm:"column:has_locale"`
		HasTimezone       bool      `gorm:"column:has_timezone"`
		HasQuota          bool      `gorm:"column:has_quota"`
		CustomFieldCount  int64     `gorm:"column:custom_field_count"`
		RiskActions       int64     `gorm:"column:risk_actions"`
		LastEntityTouch   *time.Time `gorm:"column:last_entity_touch"`
		DistinctModules   int64     `gorm:"column:distinct_modules"`
	}
	var rows []raw
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			t.id AS tenant_id,
			t.name AS tenant_name,
			COALESCE(al.c, 0) AS activity_count,
			(t.config ? 'default_locale') AS has_locale,
			(t.config ? 'timezone') AS has_timezone,
			(t.config ? 'sales_quota') AS has_quota,
			COALESCE(cf.c, 0) AS custom_field_count,
			COALESCE(risk.c, 0) AS risk_actions,
			GREATEST(lu.last_u, au.last_u, du.last_u) AS last_entity_touch,
			COALESCE(mod.c, 0) AS distinct_modules
		FROM tenants t
		LEFT JOIN (
			SELECT tenant_id, COUNT(*)::bigint AS c
			FROM audit_logs
			WHERE created_at >= NOW() - INTERVAL '30 days'
			GROUP BY tenant_id
		) al ON al.tenant_id = t.id
		LEFT JOIN (
			SELECT tenant_id, COUNT(*)::bigint AS c
			FROM custom_fields
			WHERE is_active = true
			GROUP BY tenant_id
		) cf ON cf.tenant_id = t.id
		LEFT JOIN (
			SELECT tenant_id, COUNT(*)::bigint AS c
			FROM audit_logs
			WHERE created_at >= NOW() - INTERVAL '30 days'
			  AND action IN ('lead.delete', 'deal.delete', 'account.delete', 'rbac.role.assign', 'settings.update')
			GROUP BY tenant_id
		) risk ON risk.tenant_id = t.id
		LEFT JOIN (
			SELECT tenant_id, MAX(updated_at) AS last_u FROM leads WHERE deleted_at IS NULL GROUP BY tenant_id
		) lu ON lu.tenant_id = t.id
		LEFT JOIN (
			SELECT tenant_id, MAX(updated_at) AS last_u FROM accounts WHERE deleted_at IS NULL GROUP BY tenant_id
		) au ON au.tenant_id = t.id
		LEFT JOIN (
			SELECT tenant_id, MAX(updated_at) AS last_u FROM deals WHERE deleted_at IS NULL GROUP BY tenant_id
		) du ON du.tenant_id = t.id
		LEFT JOIN (
			SELECT tenant_id, COUNT(DISTINCT resource_type)::bigint AS c
			FROM audit_logs
			WHERE created_at >= NOW() - INTERVAL '30 days'
			GROUP BY tenant_id
		) mod ON mod.tenant_id = t.id
		WHERE t.is_active = true
		ORDER BY t.name
	`).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]TenantHealthRow, 0, len(rows))
	for _, row := range rows {
		cfgScore := 0
		if row.HasLocale {
			cfgScore += 25
		}
		if row.HasTimezone {
			cfgScore += 25
		}
		if row.HasQuota {
			cfgScore += 25
		}
		if row.CustomFieldCount > 0 {
			cfgScore += 25
		}
		activity := int(clampScore(row.ActivityCount * 5))
		auditRisk := int(clampScore(100 - row.RiskActions*15))
		freshness := 40
		if row.LastEntityTouch != nil {
			days := time.Since(*row.LastEntityTouch).Hours() / 24
			switch {
			case days <= 3:
				freshness = 95
			case days <= 7:
				freshness = 80
			case days <= 14:
				freshness = 60
			case days <= 30:
				freshness = 40
			default:
				freshness = 20
			}
		}
		adoption := int(clampScore(row.DistinctModules * 20))
		scores := map[string]int{
			"activity":             activity,
			"config_completeness":  cfgScore,
			"audit_risk":           auditRisk,
			"data_freshness":       freshness,
			"feature_adoption":     adoption,
		}
		overall := (activity + cfgScore + auditRisk + freshness + adoption) / 5
		out = append(out, TenantHealthRow{
			TenantID: row.TenantID, TenantName: row.TenantName,
			Scores: scores, Overall: overall,
		})
	}
	return out, nil
}

func clampScore(v int64) int64 {
	if v > 100 {
		return 100
	}
	if v < 0 {
		return 0
	}
	return v
}

func (r *GormTenantInsightsRepository) PlanDistribution(ctx context.Context, from, to time.Time) ([]PlanCountRow, error) {
	q := r.db.WithContext(ctx).Model(&struct{}{}).Table("tenants").Where("is_active = ?", true)
	if !from.IsZero() {
		q = q.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		q = q.Where("created_at <= ?", to)
	}
	var rows []PlanCountRow
	err := q.Select("COALESCE(NULLIF(plan, ''), 'standard') AS plan, COUNT(*)::bigint AS count").
		Group("plan").
		Order("count DESC").
		Scan(&rows).Error
	return rows, err
}

func (r *GormTenantInsightsRepository) TopTenants(ctx context.Context, metric string, limit int) ([]TopTenantRow, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	var rows []TopTenantRow
	switch metric {
	case "revenue":
		err := r.db.WithContext(ctx).Raw(`
			SELECT t.id AS tenant_id, t.name AS tenant_name, COALESCE(SUM(d.amount), 0)::float8 AS value
			FROM tenants t
			LEFT JOIN deals d ON d.tenant_id = t.id AND d.deleted_at IS NULL AND d.stage = 'closed_won'
			WHERE t.is_active = true
			GROUP BY t.id, t.name
			ORDER BY value DESC
			LIMIT ?
		`, limit).Scan(&rows).Error
		return rows, err
	case "risk":
		err := r.db.WithContext(ctx).Raw(`
			SELECT t.id AS tenant_id, t.name AS tenant_name, COALESCE(risk.c, 0)::float8 AS value
			FROM tenants t
			LEFT JOIN (
				SELECT tenant_id, COUNT(*)::bigint AS c
				FROM audit_logs
				WHERE created_at >= NOW() - INTERVAL '30 days'
				  AND action IN ('lead.delete', 'deal.delete', 'account.delete', 'rbac.role.assign')
				GROUP BY tenant_id
			) risk ON risk.tenant_id = t.id
			WHERE t.is_active = true
			ORDER BY value DESC
			LIMIT ?
		`, limit).Scan(&rows).Error
		return rows, err
	default:
		err := r.db.WithContext(ctx).Raw(`
			SELECT t.id AS tenant_id, t.name AS tenant_name, COALESCE(al.c, 0)::float8 AS value
			FROM tenants t
			LEFT JOIN (
				SELECT tenant_id, COUNT(*)::bigint AS c
				FROM audit_logs
				WHERE created_at >= NOW() - INTERVAL '30 days'
				GROUP BY tenant_id
			) al ON al.tenant_id = t.id
			WHERE t.is_active = true
			ORDER BY value DESC
			LIMIT ?
		`, limit).Scan(&rows).Error
		return rows, err
	}
}
