package domain

import (
	"time"

	"github.com/google/uuid"
)

type Lead struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID            uuid.UUID  `gorm:"type:uuid;index;not null"`
	OwnerID             *uuid.UUID `gorm:"type:uuid;index"`
	Title               string     `gorm:"size:255;not null"`
	Status              string     `gorm:"size:50"`
	Source              string     `gorm:"size:50"`
	Amount              float64    `gorm:"type:decimal(15,2)"`
	ExpectedCloseDate   *time.Time `gorm:"type:date"`
	LifecycleStage      string     `gorm:"size:32"`
	EngagementScore     int16      `gorm:"default:0"`
	LastActivityAt      *time.Time `gorm:"type:timestamptz"`
	Tags                StringArray `gorm:"type:text[]"`
	RelationshipHealth  string     `gorm:"size:16"`
	ConvertedAccountID  *uuid.UUID `gorm:"type:uuid"`
	ConvertedContactID  *uuid.UUID `gorm:"type:uuid"`
	CreatedBy           uuid.UUID  `gorm:"type:uuid"`
	UpdatedBy           uuid.UUID  `gorm:"type:uuid"`
	Timestamps
	SoftDelete
}

func (l *Lead) GetTenantID() uuid.UUID  { return l.TenantID }
func (l *Lead) SetTenantID(id uuid.UUID) { l.TenantID = id }

func (Lead) TableName() string { return "leads" }
