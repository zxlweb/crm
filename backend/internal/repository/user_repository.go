package repository

import (
	"context"
	"errors"

	"crm-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type TenantBrief struct {
	ID     uuid.UUID
	Name   string
	Domain string
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	ListActiveTenantsForUser(ctx context.Context, userID uuid.UUID) ([]TenantBrief, error)
	UserBelongsToTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error)
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) ListActiveTenantsForUser(ctx context.Context, userID uuid.UUID) ([]TenantBrief, error) {
	var rows []TenantBrief
	err := r.db.WithContext(ctx).
		Table("tenants t").
		Select("t.id, t.name, t.domain").
		Joins("INNER JOIN user_tenants ut ON ut.tenant_id = t.id").
		Where("ut.user_id = ? AND t.is_active = ?", userID, true).
		Scan(&rows).Error
	return rows, err
}

func (r *GormUserRepository) UserBelongsToTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("user_tenants ut").
		Joins("INNER JOIN tenants t ON t.id = ut.tenant_id").
		Where("ut.user_id = ? AND ut.tenant_id = ? AND t.is_active = ?", userID, tenantID, true).
		Count(&count).Error
	return count > 0, err
}
