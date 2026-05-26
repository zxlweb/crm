package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// CustomField tenant-scoped custom field metadata (Phase 4).
type CustomField struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID     uuid.UUID      `gorm:"type:uuid;index;not null"`
	EntityType   string         `gorm:"size:20;not null"`
	FieldKey     string         `gorm:"size:100;not null"`
	FieldLabel   datatypes.JSON `gorm:"type:jsonb;not null"`
	FieldType    string         `gorm:"size:20;not null"`
	Required     bool           `gorm:"default:false"`
	Options      datatypes.JSON `gorm:"type:jsonb"`
	DefaultValue datatypes.JSON `gorm:"type:jsonb"`
	DisplayOrder int            `gorm:"default:100"`
	IsActive     bool           `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (CustomField) TableName() string { return "custom_fields" }
