package segment

import (
	"context"
	"errors"

	"crm-backend/internal/application/appscope"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/datascope"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

var (
	ErrNotFound        = repository.ErrSegmentNotFound
	ErrInvalidSegment  = errors.New("invalid_segment_code")
	ErrInvalidFilter   = errors.New("invalid_filter_json")
)

type Service struct {
	segments repository.SegmentRepository
	leads    repository.LeadRepository
	accounts repository.AccountRepository
	tenants  repository.TenantRepository
	enforcer *casbin.Enforcer
	scope    appscope.Provider
}

func NewService(
	segments repository.SegmentRepository,
	leads repository.LeadRepository,
	accounts repository.AccountRepository,
	tenants repository.TenantRepository,
	enforcer *casbin.Enforcer,
	scope appscope.Provider,
) *Service {
	return &Service{segments: segments, leads: leads, accounts: accounts, tenants: tenants, enforcer: enforcer, scope: scope}
}

func (s *Service) dataScope(ctx context.Context, tenantID, userID uuid.UUID) datascope.ScopeParams {
	return s.scope.Params(ctx, tenantID, userID)
}

type SegmentDTO struct {
	Code           string         `json:"code"`
	NameKey        string         `json:"name_key"`
	DescriptionKey string         `json:"description_key"`
	Filter         datatypes.JSON `json:"filter"`
	Count          *int64         `json:"count,omitempty"`
}

type CountDTO struct {
	Code  string `json:"code"`
	Count int64  `json:"count"`
}

func (s *Service) List(ctx context.Context, tenantID, userID uuid.UUID, withCount bool) ([]SegmentDTO, error) {
	templates, err := s.segments.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	opts := s.segmentOpts(ctx, tenantID)
	scope := s.dataScope(ctx, tenantID, userID)
	out := make([]SegmentDTO, len(templates))
	for i, t := range templates {
		dto := SegmentDTO{
			Code:           t.Code,
			NameKey:        t.NameI18nKey,
			DescriptionKey: "segments." + t.Code + ".description",
			Filter:         t.FilterJSON,
		}
		if withCount {
			n, err := s.countForEntity(ctx, tenantID, userID, t.Code, "lead", opts, scope)
			if err != nil {
				return nil, err
			}
			dto.Count = &n
		}
		out[i] = dto
	}
	return out, nil
}

func (s *Service) Count(ctx context.Context, tenantID, userID uuid.UUID, code, entity string) (*CountDTO, error) {
	if !crm.ValidSegmentCode(code) {
		return nil, ErrInvalidSegment
	}
	if _, err := s.segments.GetByCode(ctx, tenantID, code); err != nil {
		return nil, err
	}
	opts := s.segmentOpts(ctx, tenantID)
	scope := s.dataScope(ctx, tenantID, userID)
	n, err := s.countForEntity(ctx, tenantID, userID, code, entity, opts, scope)
	if err != nil {
		return nil, err
	}
	return &CountDTO{Code: code, Count: n}, nil
}

func (s *Service) ResolveLeadSegment(code string) error {
	return validateForLeads(code)
}

func (s *Service) ResolveAccountSegment(code string) error {
	return validateForAccounts(code)
}

func (s *Service) LeadSegmentOpts(ctx context.Context, tenantID uuid.UUID) crm.SegmentApplyOpts {
	return s.segmentOpts(ctx, tenantID)
}

func (s *Service) countForEntity(ctx context.Context, tenantID, userID uuid.UUID, code, entity string, opts crm.SegmentApplyOpts, scope datascope.ScopeParams) (int64, error) {
	switch entity {
	case "", "lead", "leads":
		_, n, err := s.leads.List(ctx, tenantID, repository.LeadListFilter{
			Page: 1, PageSize: 1, Segment: code, SegmentOpts: opts, Scope: scope,
		})
		return n, err
	case "account", "accounts":
		_, n, err := s.accounts.List(ctx, tenantID, repository.AccountListFilter{
			Page: 1, PageSize: 1, Segment: code, SegmentOpts: opts, Scope: scope,
		})
		return n, err
	default:
		return 0, ErrInvalidSegment
	}
}

func (s *Service) segmentOpts(ctx context.Context, tenantID uuid.UUID) crm.SegmentApplyOpts {
	opts := crm.SegmentApplyOpts{DaysSilent: 7, HighValueAmount: 100000}
	if s.tenants == nil {
		return opts
	}
	t, err := s.tenants.FindByID(ctx, tenantID)
	if err != nil || t == nil {
		return opts
	}
	cfg := crm.ParseTenantCRMConfig(t.Config)
	opts.DaysSilent = cfg.InsightThresholds.DaysSilent
	opts.HighValueAmount = cfg.InsightThresholds.HighValueAmount
	return opts
}

func validateForLeads(code string) error {
	if code == "" {
		return nil
	}
	if !crm.ValidSegmentCode(code) {
		return ErrInvalidSegment
	}
	ok, _, err := crm.SegmentEntityForCode(code)
	if err != nil || !ok {
		return ErrInvalidSegment
	}
	return nil
}

func validateForAccounts(code string) error {
	if code == "" {
		return nil
	}
	if !crm.ValidSegmentCode(code) {
		return ErrInvalidSegment
	}
	_, ok, err := crm.SegmentEntityForCode(code)
	if err != nil || !ok {
		return ErrInvalidSegment
	}
	return nil
}
