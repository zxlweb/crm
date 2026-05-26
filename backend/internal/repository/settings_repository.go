package repository

import (
	"context"
	"errors"

	"crm-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SettingsRepository interface {
	GetTenant(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error)
	UpdateTenant(ctx context.Context, tenantID uuid.UUID, name *string, config datatypes.JSON) error
}

type GormSettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &GormSettingsRepository{db: db}
}

func (r *GormSettingsRepository) GetTenant(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant
	err := r.db.WithContext(ctx).Where("id = ?", tenantID).First(&tenant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTenantNotFound
	}
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *GormSettingsRepository) UpdateTenant(ctx context.Context, tenantID uuid.UUID, name *string, config datatypes.JSON) error {
	updates := map[string]any{"updated_at": gorm.Expr("NOW()")}
	if name != nil {
		updates["name"] = *name
	}
	if config != nil {
		updates["config"] = config
	}
	res := r.db.WithContext(ctx).Model(&domain.Tenant{}).Where("id = ?", tenantID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrTenantNotFound
	}
	return nil
}
