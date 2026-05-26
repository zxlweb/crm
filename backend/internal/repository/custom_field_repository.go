package repository

import (
	"context"
	"errors"

	"crm-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrCustomFieldNotFound = errors.New("custom field not found")

type CustomFieldListFilter struct {
	EntityType string
	ActiveOnly bool
}

type CustomFieldRepository interface {
	List(ctx context.Context, tenantID uuid.UUID, f CustomFieldListFilter) ([]domain.CustomField, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomField, error)
	Create(ctx context.Context, f *domain.CustomField) error
	Update(ctx context.Context, f *domain.CustomField) error
	Deactivate(ctx context.Context, tenantID, id uuid.UUID) error
}

type GormCustomFieldRepository struct {
	db *gorm.DB
}

func NewCustomFieldRepository(db *gorm.DB) CustomFieldRepository {
	return &GormCustomFieldRepository{db: db}
}

func (r *GormCustomFieldRepository) List(ctx context.Context, tenantID uuid.UUID, f CustomFieldListFilter) ([]domain.CustomField, error) {
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if f.EntityType != "" {
		q = q.Where("entity_type = ?", f.EntityType)
	}
	if f.ActiveOnly {
		q = q.Where("is_active = ?", true)
	}
	var rows []domain.CustomField
	err := q.Order("display_order ASC, created_at ASC").Find(&rows).Error
	return rows, err
}

func (r *GormCustomFieldRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomField, error) {
	var row domain.CustomField
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCustomFieldNotFound
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *GormCustomFieldRepository) Create(ctx context.Context, f *domain.CustomField) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *GormCustomFieldRepository) Update(ctx context.Context, f *domain.CustomField) error {
	res := r.db.WithContext(ctx).Model(f).Where("tenant_id = ? AND id = ?", f.TenantID, f.ID).
		Updates(map[string]any{
			"field_label":    f.FieldLabel,
			"field_type":     f.FieldType,
			"required":       f.Required,
			"options":        f.Options,
			"default_value":  f.DefaultValue,
			"display_order":  f.DisplayOrder,
			"is_active":      f.IsActive,
			"updated_at":     gorm.Expr("NOW()"),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCustomFieldNotFound
	}
	return nil
}

func (r *GormCustomFieldRepository) Deactivate(ctx context.Context, tenantID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Model(&domain.CustomField{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Update("is_active", false)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCustomFieldNotFound
	}
	return nil
}
