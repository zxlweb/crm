package domain

import (
	"time"

	"github.com/google/uuid"
)

type Lead struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID           uuid.UUID  `gorm:"type:uuid;index;not null"`
	OwnerID            *uuid.UUID `gorm:"type:uuid;index"`
	Title              string     `gorm:"size:255;not null"`
	Status             string     `gorm:"size:50"`
	Source             string     `gorm:"size:50"`
	Amount             float64    `gorm:"type:decimal(15,2)"`
	ExpectedCloseDate  *time.Time `gorm:"type:date"`
	CreatedBy          uuid.UUID  `gorm:"type:uuid"`
	UpdatedBy          uuid.UUID  `gorm:"type:uuid"`
	Timestamps
	SoftDelete
}

func (l *Lead) GetTenantID() uuid.UUID  { return l.TenantID }
func (l *Lead) SetTenantID(id uuid.UUID) { l.TenantID = id }

func (Lead) TableName() string { return "leads" }
