package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AuditLog 审计日志
type AuditLog struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID     uuid.UUID      `gorm:"type:uuid;index;not null"`
	UserID       *uuid.UUID     `gorm:"type:uuid"`
	Action       string         `gorm:"size:100"`
	ResourceType string         `gorm:"size:100"`
	ResourceID   *uuid.UUID     `gorm:"type:uuid"`
	OldValue     datatypes.JSON `gorm:"type:jsonb"`
	NewValue     datatypes.JSON `gorm:"type:jsonb"`
	IPAddress    string         `gorm:"size:45"`
	CreatedAt    time.Time
}

func (AuditLog) TableName() string { return "audit_logs" }
