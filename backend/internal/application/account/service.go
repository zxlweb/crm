package account

import (
	"context"
	"errors"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/datascope"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
)

var (
	ErrNotFound         = repository.ErrAccountNotFound
	ErrInvalidLifecycle = errors.New("invalid lifecycle_stage")
	ErrInvalidSegment   = errors.New("invalid_segment_code")
)

type Service struct {
	repo     repository.AccountRepository
	tenants  repository.TenantRepository
	enforcer *casbin.Enforcer
}

func NewService(repo repository.AccountRepository, tenants repository.TenantRepository, enforcer *casbin.Enforcer) *Service {
	return &Service{repo: repo, tenants: tenants, enforcer: enforcer}
}

type AccountDTO struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	Industry           string     `json:"industry,omitempty"`
	Website            string     `json:"website,omitempty"`
	OwnerID            *uuid.UUID `json:"owner_id,omitempty"`
	LifecycleStage     string     `json:"lifecycle_stage"`
	RelationshipHealth string     `json:"relationship_health"`
	EngagementScore    int        `json:"engagement_score"`
	LastActivityAt     *time.Time `json:"last_activity_at,omitempty"`
	Tags               []string   `json:"tags"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type ListQuery struct {
	Page               int
	PageSize           int
	Search             string
	LifecycleStage     string
	RelationshipHealth string
	Segment            string
	OwnerID            *uuid.UUID
}

type ListResult struct {
	Items []AccountDTO
	Total int64
	Page  int
	Size  int
}

type CreateInput struct {
	Name           string
	Industry       string
	Website        string
	OwnerID        *uuid.UUID
	LifecycleStage string
	Tags           []string
}

type UpdateInput struct {
	Name           *string
	Industry       *string
	Website        *string
	OwnerID        *uuid.UUID
	LifecycleStage *string
	Tags           []string
	TagsSet        bool
}

func (s *Service) List(ctx context.Context, tenantID, userID uuid.UUID, q ListQuery) (*ListResult, error) {
	viewAll := s.viewAll(ctx, userID.String(), tenantID.String())
	page := q.Page
	if page < 1 {
		page = 1
	}
	size := q.PageSize
	if size < 1 {
		size = 20
	}
	if err := validateAccountSegment(q.Segment); err != nil {
		return nil, err
	}
	items, total, err := s.repo.List(ctx, tenantID, repository.AccountListFilter{
		Page:               page,
		PageSize:           size,
		Search:             q.Search,
		LifecycleStage:     q.LifecycleStage,
		RelationshipHealth: q.RelationshipHealth,
		Segment:            q.Segment,
		SegmentOpts:        s.segmentOpts(ctx, tenantID),
		OwnerID:            q.OwnerID,
		ViewAll:            viewAll,
		UserID:             userID,
	})
	if err != nil {
		return nil, err
	}
	dtos := make([]AccountDTO, len(items))
	for i := range items {
		dtos[i] = toDTO(&items[i])
	}
	return &ListResult{Items: dtos, Total: total, Page: page, Size: size}, nil
}

func (s *Service) Get(ctx context.Context, tenantID, userID, id uuid.UUID) (*AccountDTO, error) {
	a, err := s.repo.GetByID(ctx, tenantID, id, s.viewAll(ctx, userID.String(), tenantID.String()), userID)
	if err != nil {
		return nil, err
	}
	dto := toDTO(a)
	return &dto, nil
}

func (s *Service) Create(ctx context.Context, tenantID, userID uuid.UUID, in CreateInput) (*AccountDTO, error) {
	stage := in.LifecycleStage
	if stage == "" {
		stage = "acquire"
	}
	if !crm.ValidLifecycleStage(stage) {
		return nil, ErrInvalidLifecycle
	}
	owner := in.OwnerID
	if owner == nil {
		owner = &userID
	}
	tags := domain.StringArray(in.Tags)
	a := &domain.Account{
		TenantID:        tenantID,
		OwnerID:         owner,
		Name:            in.Name,
		Industry:        in.Industry,
		Website:         in.Website,
		LifecycleStage:  stage,
		EngagementScore: 0,
		Tags:            tags,
		AuditFields: domain.AuditFields{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
	}
	if err := s.repo.Create(ctx, a); err != nil {
		if err.Error() == "invalid lifecycle_stage" {
			return nil, ErrInvalidLifecycle
		}
		return nil, err
	}
	dto := toDTO(a)
	return &dto, nil
}

func (s *Service) Update(ctx context.Context, tenantID, userID, id uuid.UUID, in UpdateInput, full bool) (*AccountDTO, error) {
	viewAll := s.viewAll(ctx, userID.String(), tenantID.String())
	a, err := s.repo.GetByID(ctx, tenantID, id, viewAll, userID)
	if err != nil {
		return nil, err
	}
	oldStage := a.LifecycleStage

	if full {
		if in.Name == nil || *in.Name == "" {
			return nil, errors.New("name is required")
		}
		a.Name = *in.Name
		if in.Industry != nil {
			a.Industry = *in.Industry
		} else {
			a.Industry = ""
		}
		if in.Website != nil {
			a.Website = *in.Website
		} else {
			a.Website = ""
		}
		if in.OwnerID != nil {
			a.OwnerID = in.OwnerID
		}
		if in.LifecycleStage != nil {
			if !crm.ValidLifecycleStage(*in.LifecycleStage) {
				return nil, ErrInvalidLifecycle
			}
			a.LifecycleStage = *in.LifecycleStage
		}
		if in.TagsSet {
			a.Tags = domain.StringArray(in.Tags)
		}
	} else {
		if in.Name != nil {
			a.Name = *in.Name
		}
		if in.Industry != nil {
			a.Industry = *in.Industry
		}
		if in.Website != nil {
			a.Website = *in.Website
		}
		if in.OwnerID != nil {
			a.OwnerID = in.OwnerID
		}
		if in.LifecycleStage != nil {
			if !crm.ValidLifecycleStage(*in.LifecycleStage) {
				return nil, ErrInvalidLifecycle
			}
			a.LifecycleStage = *in.LifecycleStage
		}
		if in.TagsSet {
			a.Tags = domain.StringArray(in.Tags)
		}
	}
	a.UpdatedBy = userID

	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	_ = oldStage // lifecycle history hook in Phase 2.12
	dto := toDTO(a)
	return &dto, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, userID, id uuid.UUID) error {
	viewAll := s.viewAll(ctx, userID.String(), tenantID.String())
	if _, err := s.repo.GetByID(ctx, tenantID, id, viewAll, userID); err != nil {
		return err
	}
	return s.repo.SoftDelete(ctx, tenantID, id)
}

func validateAccountSegment(code string) error {
	if code == "" {
		return nil
	}
	if !crm.ValidSegmentCode(code) {
		return ErrInvalidSegment
	}
	return nil
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

func (s *Service) viewAll(ctx context.Context, userID, tenantID string) bool {
	return datascope.CanViewAllTenantData(ctx, s.enforcer, userID, tenantID)
}

func toDTO(a *domain.Account) AccountDTO {
	tags := []string(a.Tags)
	if tags == nil {
		tags = []string{}
	}
	return AccountDTO{
		ID:                 a.ID,
		Name:               a.Name,
		Industry:           a.Industry,
		Website:            a.Website,
		OwnerID:            a.OwnerID,
		LifecycleStage:     a.LifecycleStage,
		RelationshipHealth: crm.RelationshipHealthFromScore(a.EngagementScore),
		EngagementScore:    int(a.EngagementScore),
		LastActivityAt:     a.LastActivityAt,
		Tags:               tags,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}
