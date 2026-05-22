package datascope

import (
	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CanViewAllTenantData managers / tenant admins see all records in tenant.
func CanViewAllTenantData(enforcer *casbin.Enforcer, userID, tenantID string) bool {
	if enforcer == nil {
		return false
	}
	ok, err := enforcer.Enforce(userID, tenantID, "rbac", "manage")
	return err == nil && ok
}

// OwnerScope limits queries to owner_id = user when view-all is false.
func OwnerScope(db *gorm.DB, userID uuid.UUID, viewAll bool) *gorm.DB {
	if viewAll {
		return db
	}
	return db.Where("(owner_id = ? OR owner_id IS NULL)", userID)
}
