package domain

import (
	"time"

	"github.com/google/uuid"
)

// Account 公司（客户）
type Account struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID         uuid.UUID  `gorm:"type:uuid;index;not null"`
	OwnerID          *uuid.UUID `gorm:"type:uuid;index"`
	Name             string     `gorm:"size:255;not null"`
	Industry         string     `gorm:"size:100"`
	Website          string     `gorm:"type:text"`
	LifecycleStage   string     `gorm:"size:32;not null;default:acquire"`
	EngagementScore  int16      `gorm:"not null;default:0"`
	LastActivityAt   *time.Time `gorm:"type:timestamptz"`
	Tags             StringArray `gorm:"type:text[]"`
	AuditFields
	Timestamps
	SoftDelete
}

func (a *Account) GetTenantID() uuid.UUID  { return a.TenantID }
func (a *Account) SetTenantID(id uuid.UUID) { a.TenantID = id }

func (Account) TableName() string { return "accounts" }
