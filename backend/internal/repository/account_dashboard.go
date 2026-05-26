package repository

import (
	"context"

	"crm-backend/internal/pkg/datascope"

	"github.com/google/uuid"
)

func (r *GormAccountRepository) CountScoped(ctx context.Context, tenantID uuid.UUID, scope datascope.ScopeParams) (int64, error) {
	var n int64
	err := datascope.ApplyOwnerScope(r.base(ctx, tenantID), scope).Count(&n).Error
	return n, err
}

func (r *GormAccountRepository) CountLowEngagement(ctx context.Context, tenantID uuid.UUID, scope datascope.ScopeParams) (int64, error) {
	var n int64
	q := r.base(ctx, tenantID).Where("engagement_score < ?", 40)
	err := datascope.ApplyOwnerScope(q, scope).Count(&n).Error
	return n, err
}
