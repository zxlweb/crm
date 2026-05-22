package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string         `gorm:"uniqueIndex;not null"`
	PasswordHash string         `gorm:"not null"`
	Name         string         `gorm:"size:100"`
	AvatarURL    string
	IsSuperAdmin bool           `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string { return "users" }

type UserTenant struct {
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (UserTenant) TableName() string { return "user_tenants" }

type Tenant struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string         `gorm:"size:255;not null"`
	Domain    string         `gorm:"size:100;uniqueIndex"`
	LogoURL   string
	Config    datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Tenant) TableName() string { return "tenants" }
