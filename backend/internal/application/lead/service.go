package lead

import (
	"context"
	"errors"
	"time"

	"crm-backend/internal/application/appscope"
	"crm-backend/internal/application/deal"
	"crm-backend/internal/application/insights"
	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/datascope"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
)

var (
	ErrNotFound                = repository.ErrLeadNotFound
	ErrInvalidLifecycle        = errors.New("invalid lifecycle_stage")
	ErrInvalidStatus           = errors.New("invalid status")
	ErrInvalidStatusTransition = errors.New("invalid_status_transition")
	ErrInvalidSegment          = errors.New("invalid_segment_code")
	ErrConvertNotAllowed       = errors.New("convert_not_allowed")
	ErrConvertMissingAccount   = errors.New("convert_requires_account")
	ErrAlreadyConverted        = errors.New("lead_already_converted")
)

type Service struct {
	repo       repository.LeadRepository
	accounts   repository.AccountRepository
	activities repository.ActivityRepository
	tenants    repository.TenantRepository
	deals      *deal.Service
	enforcer   *casbin.Enforcer
	scope      appscope.Provider
}

func NewService(
	repo repository.LeadRepository,
	accounts repository.AccountRepository,
	activities repository.ActivityRepository,
	tenants repository.TenantRepository,
	enforcer *casbin.Enforcer,
	deals *deal.Service,
	scope appscope.Provider,
) *Service {
	return &Service{repo: repo, accounts: accounts, activities: activities, tenants: tenants, enforcer: enforcer, deals: deals, scope: scope}
}

func (s *Service) dataScope(ctx context.Context, tenantID, userID uuid.UUID) datascope.ScopeParams {
	return s.scope.Params(ctx, tenantID, userID)
}

type LeadDTO struct {
	ID                  uuid.UUID  `json:"id"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	OwnerID             *uuid.UUID `json:"owner_id,omitempty"`
	Title               string     `json:"title"`
	Status              string     `json:"status"`
	Source              string     `json:"source,omitempty"`
	Amount              float64    `json:"amount"`
	ExpectedCloseDate   *string    `json:"expected_close_date,omitempty"`
	LifecycleStage      string     `json:"lifecycle_stage"`
	RelationshipHealth  string     `json:"relationship_health"`
	EngagementScore     int        `json:"engagement_score"`
	LastActivityAt      *time.Time `json:"last_activity_at,omitempty"`
	Tags                []string   `json:"tags"`
	ConvertedAccountID  *uuid.UUID `json:"converted_account_id,omitempty"`
	ConvertedContactID  *uuid.UUID `json:"converted_contact_id,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type ListQuery struct {
	Page               int
	PageSize           int
	Search             string
	Status             string
	Source             string
	LifecycleStage     string
	RelationshipHealth string
	Segment            string
	OwnerID            *uuid.UUID
}

type ListResult struct {
	Items []LeadDTO
	Total int64
	Page  int
	Size  int
}

type CreateInput struct {
	Title             string
	Status            string
	Source            string
	Amount            float64
	ExpectedCloseDate *string
	OwnerID           *uuid.UUID
	LifecycleStage    string
	Tags              []string
}

type UpdateInput struct {
	Title             *string
	Status            *string
	Source            *string
	Amount            *float64
	ExpectedCloseDate *string
	OwnerID           *uuid.UUID
	LifecycleStage    *string
	Tags              []string
	TagsSet           bool
}

type ConvertAccountInput struct {
	Name string
}

type ConvertDealInput struct {
	Title  string
	Amount float64
	Stage  string
}

type ConvertInput struct {
	AccountID     *uuid.UUID
	ContactID     *uuid.UUID
	CreateAccount *ConvertAccountInput
	CreateDeal    *ConvertDealInput
}

type ConvertResult struct {
	LeadDTO
	DealID *uuid.UUID `json:"deal_id,omitempty"`
}

func (s *Service) List(ctx context.Context, tenantID, userID uuid.UUID, q ListQuery) (*ListResult, error) {
	scope := s.dataScope(ctx, tenantID, userID)
	page := q.Page
	if page < 1 {
		page = 1
	}
	size := q.PageSize
	if size < 1 {
		size = 20
	}
	if err := validateLeadSegment(q.Segment); err != nil {
		return nil, err
	}
	items, total, err := s.repo.List(ctx, tenantID, repository.LeadListFilter{
		Page:               page,
		PageSize:           size,
		Search:             q.Search,
		Status:             q.Status,
		Source:             q.Source,
		LifecycleStage:     q.LifecycleStage,
		RelationshipHealth: q.RelationshipHealth,
		Segment:            q.Segment,
		SegmentOpts:        s.segmentOpts(ctx, tenantID),
		OwnerID:            q.OwnerID,
		Scope:              scope,
	})
	if err != nil {
		return nil, err
	}
	dtos := make([]LeadDTO, len(items))
	for i := range items {
		dtos[i] = toDTO(&items[i])
	}
	return &ListResult{Items: dtos, Total: total, Page: page, Size: size}, nil
}

func (s *Service) Get(ctx context.Context, tenantID, userID, id uuid.UUID) (*LeadDTO, error) {
	l, err := s.repo.GetByID(ctx, tenantID, id, s.dataScope(ctx, tenantID, userID))
	if err != nil {
		return nil, err
	}
	dto := toDTO(l)
	return &dto, nil
}

func (s *Service) Create(ctx context.Context, tenantID, userID uuid.UUID, in CreateInput) (*LeadDTO, error) {
	stage := in.LifecycleStage
	if stage == "" {
		stage = "acquire"
	}
	if !crm.ValidLifecycleStage(stage) {
		return nil, ErrInvalidLifecycle
	}
	status := in.Status
	if status == "" {
		status = "new"
	}
	if !crm.ValidLeadStatus(status) {
		return nil, ErrInvalidStatus
	}
	owner := in.OwnerID
	if owner == nil {
		owner = &userID
	}
	var expected *time.Time
	if in.ExpectedCloseDate != nil && *in.ExpectedCloseDate != "" {
		t, err := time.Parse("2006-01-02", *in.ExpectedCloseDate)
		if err != nil {
			return nil, errors.New("invalid expected_close_date")
		}
		expected = &t
	}
	tags := domain.StringArray(in.Tags)
	l := &domain.Lead{
		TenantID:           tenantID,
		OwnerID:            owner,
		Title:              in.Title,
		Status:             status,
		Source:             in.Source,
		Amount:             in.Amount,
		ExpectedCloseDate:  expected,
		LifecycleStage:     stage,
		EngagementScore:    0,
		Tags:               tags,
		CreatedBy:          userID,
		UpdatedBy:          userID,
	}
	if err := s.repo.Create(ctx, l); err != nil {
		if err.Error() == "invalid lifecycle_stage" {
			return nil, ErrInvalidLifecycle
		}
		if err.Error() == "invalid status" {
			return nil, ErrInvalidStatus
		}
		return nil, err
	}
	dto := toDTO(l)
	return &dto, nil
}

func (s *Service) Update(ctx context.Context, tenantID, userID, id uuid.UUID, in UpdateInput, full bool) (*LeadDTO, error) {
	scope := s.dataScope(ctx, tenantID, userID)
	l, err := s.repo.GetByID(ctx, tenantID, id, scope)
	if err != nil {
		return nil, err
	}

	if full {
		if in.Title == nil || *in.Title == "" {
			return nil, errors.New("title is required")
		}
		l.Title = *in.Title
		if in.Status != nil {
			if err := applyStatusChange(l, *in.Status); err != nil {
				return nil, err
			}
		}
		if in.Source != nil {
			l.Source = *in.Source
		}
		if in.Amount != nil {
			l.Amount = *in.Amount
		}
		if in.ExpectedCloseDate != nil {
			if *in.ExpectedCloseDate == "" {
				l.ExpectedCloseDate = nil
			} else {
				t, err := time.Parse("2006-01-02", *in.ExpectedCloseDate)
				if err != nil {
					return nil, errors.New("invalid expected_close_date")
				}
				l.ExpectedCloseDate = &t
			}
		}
		if in.OwnerID != nil {
			l.OwnerID = in.OwnerID
		}
		if in.LifecycleStage != nil {
			if !crm.ValidLifecycleStage(*in.LifecycleStage) {
				return nil, ErrInvalidLifecycle
			}
			l.LifecycleStage = *in.LifecycleStage
		}
		if in.TagsSet {
			l.Tags = domain.StringArray(in.Tags)
		}
	} else {
		if in.Title != nil {
			l.Title = *in.Title
		}
		if in.Status != nil {
			if err := applyStatusChange(l, *in.Status); err != nil {
				return nil, err
			}
		}
		if in.Source != nil {
			l.Source = *in.Source
		}
		if in.Amount != nil {
			l.Amount = *in.Amount
		}
		if in.ExpectedCloseDate != nil {
			if *in.ExpectedCloseDate == "" {
				l.ExpectedCloseDate = nil
			} else {
				t, err := time.Parse("2006-01-02", *in.ExpectedCloseDate)
				if err != nil {
					return nil, errors.New("invalid expected_close_date")
				}
				l.ExpectedCloseDate = &t
			}
		}
		if in.OwnerID != nil {
			l.OwnerID = in.OwnerID
		}
		if in.LifecycleStage != nil {
			if !crm.ValidLifecycleStage(*in.LifecycleStage) {
				return nil, ErrInvalidLifecycle
			}
			l.LifecycleStage = *in.LifecycleStage
		}
		if in.TagsSet {
			l.Tags = domain.StringArray(in.Tags)
		}
	}
	l.UpdatedBy = userID

	if err := s.repo.Update(ctx, l); err != nil {
		return nil, err
	}
	dto := toDTO(l)
	return &dto, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, userID, id uuid.UUID) error {
	scope := s.dataScope(ctx, tenantID, userID)
	if _, err := s.repo.GetByID(ctx, tenantID, id, scope); err != nil {
		return err
	}
	return s.repo.SoftDelete(ctx, tenantID, id)
}

// Convert marks a qualified lead as converted and links an account (existing or newly created).
func (s *Service) Convert(ctx context.Context, tenantID, userID, id uuid.UUID, in ConvertInput) (*ConvertResult, error) {
	scope := s.dataScope(ctx, tenantID, userID)
	l, err := s.repo.GetByID(ctx, tenantID, id, scope)
	if err != nil {
		return nil, err
	}
	if l.Status == "converted" {
		return nil, ErrAlreadyConverted
	}
	if !crm.CanConvertLead(l.Status) {
		return nil, ErrConvertNotAllowed
	}

	var accountID uuid.UUID
	switch {
	case in.AccountID != nil:
		if s.accounts == nil {
			return nil, errors.New("accounts repository not configured")
		}
		acc, err := s.accounts.GetByID(ctx, tenantID, *in.AccountID, scope)
		if err != nil {
			if errors.Is(err, repository.ErrAccountNotFound) {
				return nil, ErrConvertMissingAccount
			}
			return nil, err
		}
		accountID = acc.ID
	case in.CreateAccount != nil && in.CreateAccount.Name != "":
		if s.accounts == nil {
			return nil, errors.New("accounts repository not configured")
		}
		owner := l.OwnerID
		if owner == nil {
			owner = &userID
		}
		acc := &domain.Account{
			TenantID:       tenantID,
			OwnerID:        owner,
			Name:           in.CreateAccount.Name,
			LifecycleStage: "acquire",
			AuditFields: domain.AuditFields{
				CreatedBy: userID,
				UpdatedBy: userID,
			},
		}
		if err := s.accounts.Create(ctx, acc); err != nil {
			return nil, err
		}
		accountID = acc.ID
	default:
		return nil, ErrConvertMissingAccount
	}

	l.Status = "converted"
	l.ConvertedAccountID = &accountID
	if in.ContactID != nil {
		l.ConvertedContactID = in.ContactID
	}
	l.UpdatedBy = userID

	if err := s.repo.Update(ctx, l); err != nil {
		return nil, err
	}
	dto := toDTO(l)
	result := &ConvertResult{LeadDTO: dto}
	if in.CreateDeal != nil && in.CreateDeal.Title != "" && s.deals != nil {
		created, err := s.deals.CreateFromLead(ctx, tenantID, userID, id, accountID, deal.CreateFromLeadInput{
			Title:  in.CreateDeal.Title,
			Amount: in.CreateDeal.Amount,
			Stage:  in.CreateDeal.Stage,
		})
		if err != nil {
			return nil, err
		}
		result.DealID = &created.ID
	}
	return result, nil
}

func applyStatusChange(l *domain.Lead, to string) error {
	if !crm.ValidLeadStatus(to) {
		return ErrInvalidStatus
	}
	if l.Status == "converted" && to != l.Status {
		return ErrInvalidStatusTransition
	}
	if to != l.Status && !crm.CanTransitionLeadStatus(l.Status, to) {
		return ErrInvalidStatusTransition
	}
	l.Status = to
	return nil
}

func validateLeadSegment(code string) error {
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

func (s *Service) EvaluateInsights(ctx context.Context, tenantID, userID, leadID uuid.UUID) (*insights.EvaluateResponse, error) {
	lead, err := s.repo.GetByID(ctx, tenantID, leadID, s.dataScope(ctx, tenantID, userID))
	if err != nil {
		return nil, err
	}
	acts, _, err := s.activities.List(ctx, tenantID, repository.ActivityListFilter{
		SubjectType: "lead",
		SubjectID:   leadID,
		Page:        1,
		PageSize:    50,
	})
	if err != nil {
		return nil, err
	}
	daysSilent := 7
	if s.tenants != nil {
		if t, err := s.tenants.FindByID(ctx, tenantID); err == nil && t != nil {
			cfg := crm.ParseTenantCRMConfig(t.Config)
			daysSilent = cfg.InsightThresholds.DaysSilent
		}
	}
	resp := insights.EvaluateLead(insights.LeadEvaluateInput{
		LeadID:         leadID,
		Status:         lead.Status,
		LifecycleStage: lead.LifecycleStage,
		CreatedAt:      lead.CreatedAt,
		LastActivityAt: lead.LastActivityAt,
		Activities:     acts,
		DaysSilent:     daysSilent,
	})
	lastAt := lead.LastActivityAt
	if err := s.repo.UpdateEngagementFromActivity(ctx, tenantID, leadID, userID, lastAt, int16(resp.EngagementScore)); err != nil {
		return nil, err
	}
	return &resp, nil
}

func toDTO(l *domain.Lead) LeadDTO {
	tags := []string(l.Tags)
	if tags == nil {
		tags = []string{}
	}
	var expected *string
	if l.ExpectedCloseDate != nil {
		s := l.ExpectedCloseDate.Format("2006-01-02")
		expected = &s
	}
	return LeadDTO{
		ID:                 l.ID,
		TenantID:           l.TenantID,
		OwnerID:            l.OwnerID,
		Title:              l.Title,
		Status:             l.Status,
		Source:             l.Source,
		Amount:             l.Amount,
		ExpectedCloseDate: expected,
		LifecycleStage:     l.LifecycleStage,
		RelationshipHealth: crm.RelationshipHealthFromScore(l.EngagementScore),
		EngagementScore:    int(l.EngagementScore),
		LastActivityAt:     l.LastActivityAt,
		Tags:               tags,
		ConvertedAccountID: l.ConvertedAccountID,
		ConvertedContactID: l.ConvertedContactID,
		CreatedAt:          l.CreatedAt,
		UpdatedAt:          l.UpdatedAt,
	}
}
