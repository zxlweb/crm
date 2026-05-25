package repository

import (
	"context"
	"errors"

	"crm-backend/internal/domain"
	"crm-backend/internal/infrastructure/persistence"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var ErrSegmentNotFound = errors.New("segment not found")

type SegmentRepository interface {
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.SegmentTemplate, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.SegmentTemplate, error)
	UpdateFilter(ctx context.Context, tenantID uuid.UUID, code string, filter datatypes.JSON) error
}

type GormSegmentRepository struct {
	db *gorm.DB
}

func NewSegmentRepository(db *gorm.DB) SegmentRepository {
	return &GormSegmentRepository{db: db}
}

func (r *GormSegmentRepository) base(ctx context.Context, tenantID uuid.UUID) *gorm.DB {
	return persistence.DBFromContext(r.db, ctx).Model(&domain.SegmentTemplate{}).Where("tenant_id = ?", tenantID)
}

func (r *GormSegmentRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.SegmentTemplate, error) {
	var items []domain.SegmentTemplate
	err := r.base(ctx, tenantID).Order("code ASC").Find(&items).Error
	return items, err
}

func (r *GormSegmentRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.SegmentTemplate, error) {
	var s domain.SegmentTemplate
	err := r.base(ctx, tenantID).Where("code = ?", code).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSegmentNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *GormSegmentRepository) UpdateFilter(ctx context.Context, tenantID uuid.UUID, code string, filter datatypes.JSON) error {
	res := r.base(ctx, tenantID).Where("code = ?", code).Update("filter_json", filter)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrSegmentNotFound
	}
	return nil
}
