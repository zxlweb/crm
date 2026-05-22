package repository

import (
	"context"
	"errors"

	"crm-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
)

type RoleWithPermissions struct {
	Role          domain.Role
	PermissionIDs []uuid.UUID
	UserCount     int64
}

type RBACRepository interface {
	ListPermissions(ctx context.Context) ([]domain.Permission, error)
	ListRoles(ctx context.Context, tenantID uuid.UUID) ([]RoleWithPermissions, error)
	FindRole(ctx context.Context, tenantID, roleID uuid.UUID) (*RoleWithPermissions, error)
	CreateRole(ctx context.Context, role *domain.Role) error
	UpdateRole(ctx context.Context, tenantID, roleID uuid.UUID, name, description string) error
	SetRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	ListUserRoles(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.Role, error)
	SetUserRoles(ctx context.Context, tenantID, userID uuid.UUID, roleIDs []uuid.UUID) error
	ListUserPermissions(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.Permission, error)
	PermissionsExist(ctx context.Context, ids []uuid.UUID) (bool, error)
	RolesBelongToTenant(ctx context.Context, tenantID uuid.UUID, roleIDs []uuid.UUID) (bool, error)
}

type GormRBACRepository struct {
	db *gorm.DB
}

func NewRBACRepository(db *gorm.DB) RBACRepository {
	return &GormRBACRepository{db: db}
}

func (r *GormRBACRepository) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	var rows []domain.Permission
	err := r.db.WithContext(ctx).Order("resource, action").Find(&rows).Error
	return rows, err
}

func (r *GormRBACRepository) ListRoles(ctx context.Context, tenantID uuid.UUID) ([]RoleWithPermissions, error) {
	var roles []domain.Role
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at").Find(&roles).Error; err != nil {
		return nil, err
	}
	return r.attachRoleMeta(ctx, roles)
}

func (r *GormRBACRepository) FindRole(ctx context.Context, tenantID, roleID uuid.UUID) (*RoleWithPermissions, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, err
	}
	rows, err := r.attachRoleMeta(ctx, []domain.Role{role})
	if err != nil || len(rows) == 0 {
		return nil, ErrRoleNotFound
	}
	return &rows[0], nil
}

func (r *GormRBACRepository) attachRoleMeta(ctx context.Context, roles []domain.Role) ([]RoleWithPermissions, error) {
	out := make([]RoleWithPermissions, 0, len(roles))
	for _, role := range roles {
		var permIDs []uuid.UUID
		if err := r.db.WithContext(ctx).Model(&domain.RolePermission{}).
			Where("role_id = ?", role.ID).
			Pluck("permission_id", &permIDs).Error; err != nil {
			return nil, err
		}
		var userCount int64
		if err := r.db.WithContext(ctx).Model(&domain.UserRole{}).
			Where("role_id = ? AND tenant_id = ?", role.ID, role.TenantID).
			Count(&userCount).Error; err != nil {
			return nil, err
		}
		out = append(out, RoleWithPermissions{
			Role:          role,
			PermissionIDs: permIDs,
			UserCount:     userCount,
		})
	}
	return out, nil
}

func (r *GormRBACRepository) CreateRole(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *GormRBACRepository) UpdateRole(ctx context.Context, tenantID, roleID uuid.UUID, name, description string) error {
	res := r.db.WithContext(ctx).Model(&domain.Role{}).
		Where("id = ? AND tenant_id = ?", roleID, tenantID).
		Updates(map[string]interface{}{"name": name, "description": description})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrRoleNotFound
	}
	return nil
}

func (r *GormRBACRepository) SetRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&domain.RolePermission{}).Error; err != nil {
			return err
		}
		for _, pid := range permissionIDs {
			if err := tx.Create(&domain.RolePermission{RoleID: roleID, PermissionID: pid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormRBACRepository) ListUserRoles(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.WithContext(ctx).
		Table("roles r").
		Select("r.*").
		Joins("INNER JOIN user_roles ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND ur.tenant_id = ?", userID, tenantID).
		Order("r.name").
		Scan(&roles).Error
	return roles, err
}

func (r *GormRBACRepository) SetUserRoles(ctx context.Context, tenantID, userID uuid.UUID, roleIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND tenant_id = ?", userID, tenantID).Delete(&domain.UserRole{}).Error; err != nil {
			return err
		}
		for _, rid := range roleIDs {
			if err := tx.Create(&domain.UserRole{
				UserID: userID, RoleID: rid, TenantID: tenantID,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormRBACRepository) ListUserPermissions(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.Permission, error) {
	var perms []domain.Permission
	err := r.db.WithContext(ctx).
		Table("permissions p").
		Select("DISTINCT p.*").
		Joins("INNER JOIN role_permissions rp ON rp.permission_id = p.id").
		Joins("INNER JOIN user_roles ur ON ur.role_id = rp.role_id").
		Where("ur.user_id = ? AND ur.tenant_id = ?", userID, tenantID).
		Order("p.resource, p.action").
		Scan(&perms).Error
	return perms, err
}

func (r *GormRBACRepository) PermissionsExist(ctx context.Context, ids []uuid.UUID) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Permission{}).Where("id IN ?", ids).Count(&count).Error
	return count == int64(len(ids)), err
}

func (r *GormRBACRepository) RolesBelongToTenant(ctx context.Context, tenantID uuid.UUID, roleIDs []uuid.UUID) (bool, error) {
	if len(roleIDs) == 0 {
		return true, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Role{}).
		Where("tenant_id = ? AND id IN ?", tenantID, roleIDs).
		Count(&count).Error
	return count == int64(len(roleIDs)), err
}
