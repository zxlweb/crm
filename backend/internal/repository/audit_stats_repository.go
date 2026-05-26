package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditStatsFilter struct {
	From      time.Time
	To        time.Time
	Module    string
	ActorRole string
	Action    string
	ActorID   *uuid.UUID
	Limit     int
}

type ActionCountRow struct {
	Action string
	Count  int64
}

type TrendRow struct {
	Bucket string
	Count  int64
}

type ActorCountRow struct {
	ActorID   *uuid.UUID
	ActorName string
	Count     int64
}

type AuditLogExportRow struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	Action       string
	ResourceType string
	ResourceID   *uuid.UUID
	UserID       *uuid.UUID
	UserEmail    string
	IPAddress    string
}

type AuditStatsRepository interface {
	CountByAction(ctx context.Context, tenantID uuid.UUID, f AuditStatsFilter) ([]ActionCountRow, int64, error)
	Trend(ctx context.Context, tenantID uuid.UUID, f AuditStatsFilter, granularity string) ([]TrendRow, error)
	TopActors(ctx context.Context, tenantID uuid.UUID, f AuditStatsFilter) ([]ActorCountRow, error)
	ExportRows(ctx context.Context, tenantID uuid.UUID, f AuditStatsFilter) ([]AuditLogExportRow, error)
}

type GormAuditStatsRepository struct {
	db *gorm.DB
}

func NewAuditStatsRepository(db *gorm.DB) AuditStatsRepository {
	return &GormAuditStatsRepository{db: db}
}

func (r *GormAuditStatsRepository) baseQuery(ctx context.Context, tenantID uuid.UUID, f AuditStatsFilter) *gorm.DB {
	q := r.db.WithContext(ctx).Table("audit_logs").Where("tenant_id = ?", tenantID)
	if !f.From.IsZero() {
		q = q.Where("created_at >= ?", f.From)
	}
	if !f.To.IsZero() {
		q = q.Where("created_at <= ?", f.To)
	}
	if f.Module != "" {
		q = q.Where("resource_type = ?", f.Module)
	}
	if f.Action != "" {
		q = q.Where("action = ?", f.Action)
	}
	if f.ActorID != nil {
		q = q.Where("user_id = ?", *f.ActorID)
	}
	if f.ActorRole != "" {
		q = q.Where(`EXISTS (
			SELECT 1 FROM user_roles ur
			JOIN roles ro ON ro.id = ur.role_id
			WHERE ur.user_id = audit_logs.user_id
			  AND ur.tenant_id = audit_logs.tenant_id
			  AND ro.name ILIKE ?
		)`, "%"+f.ActorRole+"%")
	}
	return q
}

func (r *GormAuditStatsRepository) CountByAction(ctx context.Context, tenantID uuid.UUID, f AuditStatsFilter) ([]ActionCountRow, int64, error) {
	var rows []ActionCountRow
	err := r.baseQuery(ctx, tenantID, f).
		Select("action, COUNT(*)::bigint AS count").
		Group("action").
		Order("count DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	for _, row := range rows {
		total += row.Count
	}
	return rows, total, nil
}

func (r *GormAuditStatsRepository) Trend(ctx context.Context, tenantID uuid.UUID, f AuditStatsFilter, granularity string) ([]TrendRow, error) {
	trunc := "day"
	if granularity == "week" {
		trunc = "week"
	}
	var rows []TrendRow
	err := r.baseQuery(ctx, tenantID, f).
		Select(fmt.Sprintf("to_char(date_trunc('%s', created_at), 'YYYY-MM-DD') AS bucket, COUNT(*)::bigint AS count", trunc)).
		Group("bucket").
		Order("bucket ASC").
		Scan(&rows).Error
	return rows, err
}

func (r *GormAuditStatsRepository) TopActors(ctx context.Context, tenantID uuid.UUID, f AuditStatsFilter) ([]ActorCountRow, error) {
	limit := f.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	var rows []ActorCountRow
	err := r.baseQuery(ctx, tenantID, f).
		Select("audit_logs.user_id AS actor_id, COALESCE(users.name, users.email, '') AS actor_name, COUNT(*)::bigint AS count").
		Joins("LEFT JOIN users ON users.id = audit_logs.user_id").
		Group("audit_logs.user_id, users.name, users.email").
		Order("count DESC").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}

func (r *GormAuditStatsRepository) ExportRows(ctx context.Context, tenantID uuid.UUID, f AuditStatsFilter) ([]AuditLogExportRow, error) {
	var rows []AuditLogExportRow
	err := r.baseQuery(ctx, tenantID, f).
		Select(`audit_logs.id, audit_logs.created_at, audit_logs.action, audit_logs.resource_type,
			audit_logs.resource_id, audit_logs.user_id, COALESCE(users.email, '') AS user_email, audit_logs.ip_address`).
		Joins("LEFT JOIN users ON users.id = audit_logs.user_id").
		Order("audit_logs.created_at DESC").
		Limit(10000).
		Scan(&rows).Error
	return rows, err
}

// FormatAuditCSV builds a minimal CSV export.
func FormatAuditCSV(rows []AuditLogExportRow) string {
	var b strings.Builder
	b.WriteString("id,created_at,action,resource_type,resource_id,user_email,ip_address\n")
	for _, row := range rows {
		rid := ""
		if row.ResourceID != nil {
			rid = row.ResourceID.String()
		}
		b.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n",
			row.ID, row.CreatedAt.UTC().Format(time.RFC3339), row.Action, row.ResourceType, rid, row.UserEmail, row.IPAddress))
	}
	return b.String()
}
