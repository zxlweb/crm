package repository

import (
	"context"

	"crm-backend/internal/pkg/datascope"

	"github.com/google/uuid"
)

func (r *GormAccountRepository) CountScoped(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	var n int64
	err := datascope.OwnerScope(r.base(ctx, tenantID), userID, viewAll).Count(&n).Error
	return n, err
}

func (r *GormAccountRepository) CountLowEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	var n int64
	q := r.base(ctx, tenantID).Where("engagement_score < ?", 40)
	err := datascope.OwnerScope(q, userID, viewAll).Count(&n).Error
	return n, err
}
