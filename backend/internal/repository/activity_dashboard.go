package repository

import (
	"context"
	"time"

	"crm-backend/internal/infrastructure/persistence"
	"crm-backend/internal/pkg/datascope"

	"github.com/google/uuid"
)

func (r *GormActivityRepository) CountSince(ctx context.Context, tenantID uuid.UUID, since time.Time) (int64, error) {
	var n int64
	err := r.base(ctx, tenantID).Where("occurred_at >= ?", since).Count(&n).Error
	return n, err
}

func (r *GormActivityRepository) CountLeadTouchesSince(ctx context.Context, tenantID uuid.UUID, since time.Time, scope datascope.ScopeParams) (int64, error) {
	var n int64
	q := r.base(ctx, tenantID).Where("occurred_at >= ? AND subject_type = ?", since, "lead")
	if scope.Level != datascope.LevelAll {
		sub := persistence.DBFromContext(r.db, ctx).
			Table("leads").
			Select("id").
			Where("tenant_id = ? AND deleted_at IS NULL", tenantID)
		sub = datascope.ApplyOwnerScope(sub, scope)
		q = q.Where("subject_id IN (?)", sub)
	}
	err := q.Count(&n).Error
	return n, err
}

func (r *GormActivityRepository) CountAccountTouchesSince(ctx context.Context, tenantID uuid.UUID, since time.Time, scope datascope.ScopeParams) (int64, error) {
	var n int64
	q := r.base(ctx, tenantID).Where("occurred_at >= ? AND subject_type = ?", since, "account")
	if scope.Level != datascope.LevelAll {
		sub := persistence.DBFromContext(r.db, ctx).
			Table("accounts").
			Select("id").
			Where("tenant_id = ? AND deleted_at IS NULL", tenantID)
		sub = datascope.ApplyOwnerScope(sub, scope)
		q = q.Where("subject_id IN (?)", sub)
	}
	err := q.Count(&n).Error
	return n, err
}
