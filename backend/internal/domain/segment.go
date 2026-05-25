package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SegmentTemplate struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID    uuid.UUID      `gorm:"type:uuid;index;not null"`
	Code        string         `gorm:"size:50;not null"`
	NameI18nKey string        `gorm:"size:100;not null"`
	FilterJSON  datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
	IsSystem    bool           `gorm:"not null;default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SegmentTemplate) TableName() string { return "segment_templates" }
