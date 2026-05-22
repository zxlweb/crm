package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role 租户内角色
type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Name        string    `gorm:"size:100;not null"`
	Description string
	IsSystem    bool `gorm:"default:false"`
	CreatedAt   time.Time
}

func (Role) TableName() string { return "roles" }

// Permission 全局权限资源（resource:action）
type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Resource    string    `gorm:"size:100;not null;uniqueIndex:idx_perm_resource_action"`
	Action      string    `gorm:"size:50;not null;uniqueIndex:idx_perm_resource_action"`
	Description string
}

func (Permission) TableName() string { return "permissions" }

// RolePermission 角色-权限关联
type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (RolePermission) TableName() string { return "role_permissions" }

// UserRole 用户-角色关联（租户内）
type UserRole struct {
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoleID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (UserRole) TableName() string { return "user_roles" }
