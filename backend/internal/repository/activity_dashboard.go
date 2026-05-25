package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r *GormActivityRepository) CountSince(ctx context.Context, tenantID uuid.UUID, since time.Time) (int64, error) {
	var n int64
	err := r.base(ctx, tenantID).Where("occurred_at >= ?", since).Count(&n).Error
	return n, err
}

func (r *GormActivityRepository) CountLeadTouchesSince(ctx context.Context, tenantID uuid.UUID, since time.Time, viewAll bool, userID uuid.UUID) (int64, error) {
	var n int64
	q := r.base(ctx, tenantID).Where("occurred_at >= ? AND subject_type = ?", since, "lead")
	if !viewAll {
		q = q.Where(`subject_id IN (
			SELECT id FROM leads WHERE tenant_id = ? AND (owner_id = ? OR owner_id IS NULL) AND deleted_at IS NULL
		)`, tenantID, userID)
	}
	err := q.Count(&n).Error
	return n, err
}

func (r *GormActivityRepository) CountAccountTouchesSince(ctx context.Context, tenantID uuid.UUID, since time.Time, viewAll bool, userID uuid.UUID) (int64, error) {
	var n int64
	q := r.base(ctx, tenantID).Where("occurred_at >= ? AND subject_type = ?", since, "account")
	if !viewAll {
		q = q.Where(`subject_id IN (
			SELECT id FROM accounts WHERE tenant_id = ? AND (owner_id = ? OR owner_id IS NULL) AND deleted_at IS NULL
		)`, tenantID, userID)
	}
	err := q.Count(&n).Error
	return n, err
}
