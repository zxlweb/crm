package appscope

import (
	"context"

	"crm-backend/internal/pkg/datascope"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
)

// Provider resolves CRM data scope (all / department / self) for application services.
type Provider struct {
	Resolver *datascope.Resolver
	Enforcer *casbin.Enforcer
}

func (p *Provider) Params(ctx context.Context, tenantID, userID uuid.UUID) datascope.ScopeParams {
	if p != nil && p.Resolver != nil {
		return p.Resolver.Params(ctx, tenantID, userID)
	}
	sp := datascope.ScopeParams{Level: datascope.LevelSelf, TenantID: tenantID, UserID: userID}
	if p != nil && p.Enforcer != nil &&
		datascope.CanViewAllTenantData(ctx, p.Enforcer, userID.String(), tenantID.String()) {
		sp.Level = datascope.LevelAll
	}
	return sp
}

func (p *Provider) CanViewTeamRanking(ctx context.Context, tenantID, userID uuid.UUID, scope datascope.ScopeParams) bool {
	if p != nil && p.Resolver != nil {
		return p.Resolver.CanViewTeamRanking(ctx, tenantID, userID, scope)
	}
	return scope.Level == datascope.LevelAll
}
