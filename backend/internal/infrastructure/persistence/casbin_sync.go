package persistence

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type policyRow struct {
	UserID   string
	TenantID string
	RoleID   string
	Resource string
	Action   string
}

// SyncCasbinPolicies 从数据库加载 RBAC 策略到 Casbin
func SyncCasbinPolicies(db *gorm.DB, enforcer *casbin.Enforcer) error {
	enforcer.ClearPolicy()

	var rows []policyRow
	err := db.Raw(`
		SELECT
			ur.user_id::text AS user_id,
			ur.tenant_id::text AS tenant_id,
			ur.role_id::text AS role_id,
			p.resource,
			p.action
		FROM user_roles ur
		JOIN role_permissions rp ON rp.role_id = ur.role_id
		JOIN permissions p ON p.id = rp.permission_id
	`).Scan(&rows).Error
	if err != nil {
		return fmt.Errorf("load rbac policies: %w", err)
	}

	for _, row := range rows {
		_, _ = enforcer.AddGroupingPolicy(row.UserID, row.RoleID, row.TenantID)
		_, _ = enforcer.AddPolicy(row.RoleID, row.TenantID, row.Resource, row.Action)
	}

	log.Printf("✅ Casbin 已同步 %d 条策略", len(rows))
	return nil
}
