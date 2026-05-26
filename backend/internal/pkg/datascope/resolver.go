package datascope

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
)

// RoleNameLister returns role display names for a user in a tenant.
type RoleNameLister interface {
	ListUserRoleNames(ctx context.Context, tenantID, userID uuid.UUID) ([]string, error)
}

// DepartmentReader returns the user's department in a tenant.
type DepartmentReader interface {
	UserDepartment(ctx context.Context, tenantID, userID uuid.UUID) (string, error)
}

// Resolver computes data scope from RBAC + org membership.
type Resolver struct {
	Enforcer *casbin.Enforcer
	Roles    RoleNameLister
	Users    DepartmentReader
}

func (r *Resolver) Params(ctx context.Context, tenantID, userID uuid.UUID) ScopeParams {
	p := ScopeParams{Level: LevelSelf, TenantID: tenantID, UserID: userID}
	if r == nil || r.Enforcer == nil {
		return p
	}
	uid, tid := userID.String(), tenantID.String()
	if CanViewAllTenantData(ctx, r.Enforcer, uid, tid) {
		p.Level = LevelAll
		return p
	}
	// 事业部隔离：非租户管理员且配置了部门 → 仅本部门 owner 数据
	if r.Users != nil {
		if dept, err := r.Users.UserDepartment(ctx, tenantID, userID); err == nil && dept != "" {
			p.Level = LevelDepartment
			p.Department = dept
		}
	}
	return p
}

func (r *Resolver) DataScope(ctx context.Context, tenantID, userID uuid.UUID) string {
	return r.Params(ctx, tenantID, userID).APIScope()
}

// CanViewTeamRanking: 租户管理员看部门排行；销售经理/客户经理看本部门成员排行。
func (r *Resolver) CanViewTeamRanking(ctx context.Context, tenantID, userID uuid.UUID, scope ScopeParams) bool {
	if scope.Level == LevelAll {
		return true
	}
	if scope.Level == LevelDepartment && r != nil && r.Roles != nil {
		return IsManagerRole(ctx, r.Roles, tenantID, userID)
	}
	return false
}

// IsManagerRole identifies roles that may view team ranking within a department.
func IsManagerRole(ctx context.Context, roles RoleNameLister, tenantID, userID uuid.UUID) bool {
	names, err := roles.ListUserRoleNames(ctx, tenantID, userID)
	if err != nil {
		return false
	}
	for _, n := range names {
		switch n {
		case "销售经理", "Sales Manager", "客户经理", "Account Manager":
			return true
		}
	}
	return false
}
