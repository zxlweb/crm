package domain

import (
	"time"

	"github.com/google/uuid"
)

// Contact 联系人
type Contact struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID        uuid.UUID  `gorm:"type:uuid;index;not null"`
	AccountID       *uuid.UUID `gorm:"type:uuid;index"`
	OwnerID         *uuid.UUID `gorm:"type:uuid;index"`
	FirstName       string     `gorm:"size:100"`
	LastName        string     `gorm:"size:100"`
	Email           string     `gorm:"size:255"`
	Phone           string     `gorm:"size:50"`
	IsPrimary       bool       `gorm:"not null;default:false"`
	LifecycleStage  string     `gorm:"size:32;not null;default:acquire"`
	EngagementScore int16      `gorm:"not null;default:0"`
	LastActivityAt  *time.Time `gorm:"type:timestamptz"`
	Tags            StringArray `gorm:"type:text[]"`
	AuditFields
	Timestamps
	SoftDelete
}

func (c *Contact) GetTenantID() uuid.UUID  { return c.TenantID }
func (c *Contact) SetTenantID(id uuid.UUID) { c.TenantID = id }

func (Contact) TableName() string { return "contacts" }
