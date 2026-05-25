package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Activity struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID        uuid.UUID      `gorm:"type:uuid;index;not null"`
	SubjectType     string         `gorm:"size:20;not null"`
	SubjectID       uuid.UUID      `gorm:"type:uuid;not null"`
	EventType       string         `gorm:"size:32;not null"`
	Direction       string         `gorm:"size:16"`
	Body            string         `gorm:"type:text"`
	Metadata        datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
	Sentiment       *string        `gorm:"size:20"`
	SentimentSource *string        `gorm:"size:16"`
	OccurredAt      time.Time      `gorm:"type:timestamptz;not null"`
	CreatedBy       *uuid.UUID     `gorm:"type:uuid"`
	Timestamps
	SoftDelete
}

func (a *Activity) GetTenantID() uuid.UUID  { return a.TenantID }
func (a *Activity) SetTenantID(id uuid.UUID) { a.TenantID = id }

func (Activity) TableName() string { return "activities" }
