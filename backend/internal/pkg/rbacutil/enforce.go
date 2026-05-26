package rbacutil

import (
	"context"

	"crm-backend/internal/pkg/activerole"

	"github.com/casbin/casbin/v2"
)

// SubjectForEnforce uses active role ID when present, otherwise user ID (Casbin g policies).
func SubjectForEnforce(ctx context.Context, userID string) string {
	if rid, ok := activerole.FromContext(ctx); ok {
		return rid
	}
	return userID
}

// Enforce checks permission for the request subject (active role or user union via g).
func Enforce(ctx context.Context, enforcer *casbin.Enforcer, userID, tenantID, resource, action string) (bool, error) {
	if enforcer == nil {
		return false, nil
	}
	sub := SubjectForEnforce(ctx, userID)
	return enforcer.Enforce(sub, tenantID, resource, action)
}
