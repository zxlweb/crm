package repository

import (
	"context"
	"errors"

	"crm-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrTenantNotFound = errors.New("tenant not found")

type TenantListFilter struct {
	Page     int
	PageSize int
	Search   string
	IsActive *bool
}

type TenantRepository interface {
	List(ctx context.Context, filter TenantListFilter) ([]domain.Tenant, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
	CountUsers(ctx context.Context, tenantID uuid.UUID) (int64, error)
	CountAll(ctx context.Context) (total int64, active int64, err error)
	CountAllUsers(ctx context.Context) (int64, error)
}

type GormTenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &GormTenantRepository{db: db}
}

func (r *GormTenantRepository) List(ctx context.Context, filter TenantListFilter) ([]domain.Tenant, int64, error) {
	q := r.db.WithContext(ctx).Model(&domain.Tenant{})
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		q = q.Where("name ILIKE ? OR domain ILIKE ?", like, like)
	}
	if filter.IsActive != nil {
		q = q.Where("is_active = ?", *filter.IsActive)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var rows []domain.Tenant
	err := q.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error
	return rows, total, err
}

func (r *GormTenantRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tenant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTenantNotFound
	}
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *GormTenantRepository) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	res := r.db.WithContext(ctx).Model(&domain.Tenant{}).
		Where("id = ?", id).
		Update("is_active", active)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrTenantNotFound
	}
	return nil
}

func (r *GormTenantRepository) CountUsers(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("user_tenants").
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	return count, err
}

func (r *GormTenantRepository) CountAll(ctx context.Context) (int64, int64, error) {
	var total, active int64
	if err := r.db.WithContext(ctx).Model(&domain.Tenant{}).Count(&total).Error; err != nil {
		return 0, 0, err
	}
	if err := r.db.WithContext(ctx).Model(&domain.Tenant{}).Where("is_active = ?", true).Count(&active).Error; err != nil {
		return 0, 0, err
	}
	return total, active, nil
}

func (r *GormTenantRepository) CountAllUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.User{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}
