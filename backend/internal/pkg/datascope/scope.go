package datascope

import (
	"context"

	"crm-backend/internal/pkg/rbacutil"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CanViewAllTenantData managers / tenant admins see all records in tenant.
func CanViewAllTenantData(ctx context.Context, enforcer *casbin.Enforcer, userID, tenantID string) bool {
	ok, err := rbacutil.Enforce(ctx, enforcer, userID, tenantID, "rbac", "manage")
	return err == nil && ok
}

// OwnerScope limits queries to owner_id = user when view-all is false.
func OwnerScope(db *gorm.DB, userID uuid.UUID, viewAll bool) *gorm.DB {
	if viewAll {
		return db
	}
	return db.Where("(owner_id = ? OR owner_id IS NULL)", userID)
}

// ResolveDataScope maps Casbin policies to API data_scope (legacy: prefer Resolver.DataScope).
func ResolveDataScope(ctx context.Context, enforcer *casbin.Enforcer, userID, tenantID string) string {
	if CanViewAllTenantData(ctx, enforcer, userID, tenantID) {
		return LevelAll
	}
	return LevelSelf
}

// CanViewTeamData is deprecated; use Resolver.CanViewTeamRanking for role-aware checks.
func CanViewTeamData(p ScopeParams) bool {
	return p.Level == LevelAll || p.Level == LevelDepartment
}

// CanAccessDashboard reports dashboard:view or module view for summary.
func CanAccessDashboard(ctx context.Context, enforcer *casbin.Enforcer, userID, tenantID string) bool {
	if enforcer == nil {
		return false
	}
	for _, pair := range [][2]string{
		{"dashboard", "view"},
		{"leads", "view"},
		{"deals", "view"},
	} {
		ok, err := rbacutil.Enforce(ctx, enforcer, userID, tenantID, pair[0], pair[1])
		if err == nil && ok {
			return true
		}
	}
	return false
}
