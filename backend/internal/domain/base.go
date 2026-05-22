package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TenantScoped 标记需要租户隔离的业务实体
type TenantScoped interface {
	GetTenantID() uuid.UUID
	SetTenantID(uuid.UUID)
}

// Timestamps 通用时间戳字段
type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SoftDelete 软删除
type SoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// AuditFields 审计字段（创建人/更新人）
type AuditFields struct {
	CreatedBy uuid.UUID `gorm:"type:uuid"`
	UpdatedBy uuid.UUID `gorm:"type:uuid"`
}
