package persistence

import (
	"context"

	"crm-backend/internal/pkg/tenant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TenantScope 为查询自动附加 tenant_id 过滤
func TenantScope(tenantID uuid.UUID) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if tenantID == uuid.Nil {
			return db
		}
		return db.Where("tenant_id = ?", tenantID)
	}
}

// DBFromContext 返回附带租户 Scope 的 GORM 会话（从 context 读取 tenant_id）
func DBFromContext(db *gorm.DB, ctx context.Context) *gorm.DB {
	session := db.WithContext(ctx)
	if tid, ok := tenant.IDFromContext(ctx); ok {
		return session.Scopes(TenantScope(tid))
	}
	return session
}
