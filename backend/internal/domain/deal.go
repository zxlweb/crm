package domain

import (
	"time"

	"github.com/google/uuid"
)

type Deal struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID          uuid.UUID  `gorm:"type:uuid;index;not null"`
	OwnerID           *uuid.UUID `gorm:"type:uuid;index"`
	Title             string     `gorm:"size:255;not null"`
	Stage             string     `gorm:"size:32;not null;default:qualification"`
	Amount            float64    `gorm:"type:decimal(18,2);default:0"`
	Currency          string     `gorm:"size:8;default:CNY"`
	Probability       int16      `gorm:"default:0"`
	ExpectedCloseDate *time.Time `gorm:"type:date"`
	AccountID         *uuid.UUID `gorm:"type:uuid;index"`
	LeadID            *uuid.UUID `gorm:"type:uuid;index"`
	ContactID         *uuid.UUID `gorm:"type:uuid"`
	Description       string     `gorm:"type:text"`
	LostReason        string     `gorm:"size:500"`
	ClosedAt          *time.Time `gorm:"type:timestamptz"`
	EngagementScore   int16      `gorm:"default:0"`
	LastActivityAt    *time.Time `gorm:"type:timestamptz"`
	Tags              StringArray `gorm:"type:text[]"`
	AuditFields
	Timestamps
	SoftDelete
}

func (d *Deal) GetTenantID() uuid.UUID  { return d.TenantID }
func (d *Deal) SetTenantID(id uuid.UUID) { d.TenantID = id }

func (Deal) TableName() string { return "deals" }
