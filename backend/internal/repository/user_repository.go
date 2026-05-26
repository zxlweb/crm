package repository

import (
	"context"
	"errors"

	"crm-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailExists   = errors.New("email already exists")
	ErrDomainExists  = errors.New("domain already exists")
)

// RegisterInput 自助注册：创建租户 + 管理员用户 + 默认角色
type RegisterInput struct {
	Email        string
	PasswordHash string
	Name         string
	CompanyName  string
	Domain       string
}

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
	UserDepartment(ctx context.Context, tenantID, userID uuid.UUID) (string, error)
	RegisterWithTenant(ctx context.Context, in RegisterInput) (*domain.User, uuid.UUID, error)
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

func (r *GormUserRepository) UserDepartment(ctx context.Context, tenantID, userID uuid.UUID) (string, error) {
	var row struct {
		Department string
	}
	err := r.db.WithContext(ctx).
		Table("user_tenants").
		Select("department").
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Scan(&row).Error
	return row.Department, err
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

func (r *GormUserRepository) RegisterWithTenant(ctx context.Context, in RegisterInput) (*domain.User, uuid.UUID, error) {
	var user domain.User
	var tenantID uuid.UUID

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var emailCount int64
		if err := tx.Model(&domain.User{}).Where("email = ?", in.Email).Count(&emailCount).Error; err != nil {
			return err
		}
		if emailCount > 0 {
			return ErrEmailExists
		}

		var domainCount int64
		if err := tx.Model(&domain.Tenant{}).Where("domain = ?", in.Domain).Count(&domainCount).Error; err != nil {
			return err
		}
		if domainCount > 0 {
			return ErrDomainExists
		}

		tenant := domain.Tenant{Name: in.CompanyName, Domain: in.Domain, IsActive: true}
		if err := tx.Create(&tenant).Error; err != nil {
			return err
		}
		tenantID = tenant.ID

		user = domain.User{
			Email:        in.Email,
			PasswordHash: in.PasswordHash,
			Name:         in.Name,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Create(&domain.UserTenant{UserID: user.ID, TenantID: tenant.ID}).Error; err != nil {
			return err
		}

		role := domain.Role{
			TenantID:    tenant.ID,
			Name:        "Tenant Admin",
			Description: "租户管理员",
			IsSystem:    true,
		}
		if err := tx.Create(&role).Error; err != nil {
			return err
		}

		var perms []domain.Permission
		if err := tx.Find(&perms).Error; err != nil {
			return err
		}
		for _, p := range perms {
			if err := tx.Create(&domain.RolePermission{RoleID: role.ID, PermissionID: p.ID}).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(&domain.UserRole{UserID: user.ID, RoleID: role.ID, TenantID: tenant.ID}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, uuid.Nil, err
	}
	return &user, tenantID, nil
}
