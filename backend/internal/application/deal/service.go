package deal

import (
	"context"
	"errors"
	"time"

	"crm-backend/internal/application/appscope"
	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/datascope"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
)

var (
	ErrNotFound                = repository.ErrDealNotFound
	ErrInvalidStage            = errors.New("invalid stage")
	ErrInvalidStageTransition  = errors.New("invalid_stage_transition")
	ErrDealClosedReadonly      = errors.New("deal_closed_readonly")
	ErrInvalidCurrency         = errors.New("invalid currency")
	ErrInvalidAmount           = errors.New("invalid amount")
	ErrInvalidProbability      = errors.New("invalid probability")
)

type Service struct {
	repo     repository.DealRepository
	accounts repository.AccountRepository
	enforcer *casbin.Enforcer
	scope    appscope.Provider
}

func NewService(repo repository.DealRepository, accounts repository.AccountRepository, enforcer *casbin.Enforcer, scope appscope.Provider) *Service {
	return &Service{repo: repo, accounts: accounts, enforcer: enforcer, scope: scope}
}

func (s *Service) dataScope(ctx context.Context, tenantID, userID uuid.UUID) datascope.ScopeParams {
	return s.scope.Params(ctx, tenantID, userID)
}

type DealDTO struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	OwnerID           *uuid.UUID `json:"owner_id,omitempty"`
	Title             string     `json:"title"`
	Stage             string     `json:"stage"`
	Amount            float64    `json:"amount"`
	Currency          string     `json:"currency"`
	Probability       int        `json:"probability"`
	ExpectedCloseDate *string    `json:"expected_close_date,omitempty"`
	AccountID         *uuid.UUID `json:"account_id,omitempty"`
	LeadID            *uuid.UUID `json:"lead_id,omitempty"`
	ContactID         *uuid.UUID `json:"contact_id,omitempty"`
	Description       string     `json:"description,omitempty"`
	Tags              []string   `json:"tags"`
	LostReason        string     `json:"lost_reason,omitempty"`
	ClosedAt          *time.Time `json:"closed_at,omitempty"`
	EngagementScore   int        `json:"engagement_score"`
	LastActivityAt    *time.Time `json:"last_activity_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type ListQuery struct {
	Page              int
	PageSize          int
	Search            string
	Stage             string
	Stages            []string
	OwnerID           *uuid.UUID
	AccountID         *uuid.UUID
	LeadID            *uuid.UUID
	ExpectedCloseFrom *string
	ExpectedCloseTo   *string
	MinAmount         *float64
	MaxAmount         *float64
}

type ListResult struct {
	Items []DealDTO
	Total int64
	Page  int
	Size  int
}

type CreateInput struct {
	Title             string
	Stage             string
	Amount            float64
	Currency          string
	Probability       *int
	ExpectedCloseDate *string
	AccountID         *uuid.UUID
	LeadID            *uuid.UUID
	ContactID         *uuid.UUID
	OwnerID           *uuid.UUID
	Description       string
	Tags              []string
}

type UpdateInput struct {
	Title             *string
	Stage             *string
	Amount            *float64
	Currency          *string
	Probability       *int
	ExpectedCloseDate *string
	AccountID         *uuid.UUID
	LeadID            *uuid.UUID
	ContactID         *uuid.UUID
	OwnerID           *uuid.UUID
	Description       *string
	LostReason        *string
	Tags              []string
	TagsSet           bool
}

type StageInput struct {
	Stage      string
	LostReason *string
	Note       string
}

type PipelineQuery struct {
	OwnerID   *uuid.UUID
	AccountID *uuid.UUID
}

type PipelineItemDTO struct {
	ID                uuid.UUID  `json:"id"`
	Title             string     `json:"title"`
	Amount            float64    `json:"amount"`
	Currency          string     `json:"currency"`
	Probability       int        `json:"probability"`
	ExpectedCloseDate *string    `json:"expected_close_date,omitempty"`
	AccountID         *uuid.UUID `json:"account_id,omitempty"`
	OwnerID           *uuid.UUID `json:"owner_id,omitempty"`
}

type PipelineStageDTO struct {
	Stage       string            `json:"stage"`
	Count       int64             `json:"count"`
	AmountTotal float64           `json:"amount_total"`
	Items       []PipelineItemDTO `json:"items"`
}

type PipelineResult struct {
	Stages  []PipelineStageDTO `json:"stages"`
	Summary PipelineSummaryDTO `json:"summary"`
}

type PipelineSummaryDTO struct {
	OpenCount    int64   `json:"open_count"`
	OpenAmount   float64 `json:"open_amount"`
	WonCountMTD  int64   `json:"won_count_mtd"`
	WonAmountMTD float64 `json:"won_amount_mtd"`
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
	closeFrom, closeTo, err := parseDateRange(q.ExpectedCloseFrom, q.ExpectedCloseTo)
	if err != nil {
		return nil, err
	}
	items, total, err := s.repo.List(ctx, tenantID, repository.DealListFilter{
		Page:              page,
		PageSize:          size,
		Search:            q.Search,
		Stage:             q.Stage,
		Stages:            q.Stages,
		OwnerID:           q.OwnerID,
		AccountID:         q.AccountID,
		LeadID:            q.LeadID,
		ExpectedCloseFrom: closeFrom,
		ExpectedCloseTo:   closeTo,
		MinAmount:         q.MinAmount,
		MaxAmount:         q.MaxAmount,
		Scope:             scope,
	})
	if err != nil {
		return nil, err
	}
	dtos := make([]DealDTO, len(items))
	for i := range items {
		dtos[i] = toDTO(&items[i])
	}
	return &ListResult{Items: dtos, Total: total, Page: page, Size: size}, nil
}

func (s *Service) Get(ctx context.Context, tenantID, userID, id uuid.UUID) (*DealDTO, error) {
	d, err := s.repo.GetByID(ctx, tenantID, id, s.dataScope(ctx, tenantID, userID))
	if err != nil {
		return nil, err
	}
	dto := toDTO(d)
	return &dto, nil
}

func (s *Service) Create(ctx context.Context, tenantID, userID uuid.UUID, in CreateInput) (*DealDTO, error) {
	stage := in.Stage
	if stage == "" {
		stage = crm.DealStageQualification
	}
	if !crm.ValidDealStage(stage) {
		return nil, ErrInvalidStage
	}
	currency := in.Currency
	if currency == "" {
		currency = "CNY"
	}
	if !crm.ValidDealCurrency(currency) {
		return nil, ErrInvalidCurrency
	}
	if in.Amount < 0 {
		return nil, ErrInvalidAmount
	}
	prob := int16(0)
	if in.Probability != nil {
		if *in.Probability < 0 || *in.Probability > 100 {
			return nil, ErrInvalidProbability
		}
		prob = int16(*in.Probability)
	}
	if err := s.validateAccount(ctx, tenantID, userID, in.AccountID); err != nil {
		return nil, err
	}
	expected, err := parseOptionalDate(in.ExpectedCloseDate)
	if err != nil {
		return nil, err
	}
	owner := in.OwnerID
	if owner == nil {
		owner = &userID
	}
	tags := domain.StringArray(in.Tags)
	d := &domain.Deal{
		TenantID:          tenantID,
		OwnerID:           owner,
		Title:             in.Title,
		Stage:             stage,
		Amount:            in.Amount,
		Currency:          currency,
		Probability:       prob,
		ExpectedCloseDate: expected,
		AccountID:         in.AccountID,
		LeadID:            in.LeadID,
		ContactID:         in.ContactID,
		Description:       in.Description,
		Tags:              tags,
		AuditFields: domain.AuditFields{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
	}
	if err := s.repo.Create(ctx, d); err != nil {
		if err.Error() == "invalid stage" {
			return nil, ErrInvalidStage
		}
		if err.Error() == "invalid currency" {
			return nil, ErrInvalidCurrency
		}
		return nil, err
	}
	dto := toDTO(d)
	return &dto, nil
}

func (s *Service) Update(ctx context.Context, tenantID, userID, id uuid.UUID, in UpdateInput, full bool) (*DealDTO, error) {
	scope := s.dataScope(ctx, tenantID, userID)
	d, err := s.repo.GetByID(ctx, tenantID, id, scope)
	if err != nil {
		return nil, err
	}

	if full {
		if in.Title == nil || *in.Title == "" {
			return nil, errors.New("title is required")
		}
		d.Title = *in.Title
		if in.Stage != nil {
			if err := applyStageChange(d, *in.Stage, in.LostReason); err != nil {
				return nil, err
			}
		}
		if in.Amount != nil {
			if *in.Amount < 0 {
				return nil, ErrInvalidAmount
			}
			d.Amount = *in.Amount
		}
		if in.Currency != nil {
			if !crm.ValidDealCurrency(*in.Currency) {
				return nil, ErrInvalidCurrency
			}
			d.Currency = *in.Currency
		}
		if in.Probability != nil {
			if *in.Probability < 0 || *in.Probability > 100 {
				return nil, ErrInvalidProbability
			}
			d.Probability = int16(*in.Probability)
		}
		if in.ExpectedCloseDate != nil {
			expected, err := parseOptionalDate(in.ExpectedCloseDate)
			if err != nil {
				return nil, err
			}
			d.ExpectedCloseDate = expected
		}
		if in.AccountID != nil {
			if err := s.validateAccount(ctx, tenantID, userID, in.AccountID); err != nil {
				return nil, err
			}
			d.AccountID = in.AccountID
		}
		if in.LeadID != nil {
			d.LeadID = in.LeadID
		}
		if in.ContactID != nil {
			d.ContactID = in.ContactID
		}
		if in.OwnerID != nil {
			d.OwnerID = in.OwnerID
		}
		if in.Description != nil {
			d.Description = *in.Description
		}
		if in.LostReason != nil {
			d.LostReason = *in.LostReason
		}
		if in.TagsSet {
			d.Tags = domain.StringArray(in.Tags)
		}
	} else {
		if in.Title != nil {
			d.Title = *in.Title
		}
		if in.Stage != nil {
			if err := applyStageChange(d, *in.Stage, in.LostReason); err != nil {
				return nil, err
			}
		}
		if in.Amount != nil {
			if *in.Amount < 0 {
				return nil, ErrInvalidAmount
			}
			d.Amount = *in.Amount
		}
		if in.Currency != nil {
			if !crm.ValidDealCurrency(*in.Currency) {
				return nil, ErrInvalidCurrency
			}
			d.Currency = *in.Currency
		}
		if in.Probability != nil {
			if *in.Probability < 0 || *in.Probability > 100 {
				return nil, ErrInvalidProbability
			}
			d.Probability = int16(*in.Probability)
		}
		if in.ExpectedCloseDate != nil {
			expected, err := parseOptionalDate(in.ExpectedCloseDate)
			if err != nil {
				return nil, err
			}
			d.ExpectedCloseDate = expected
		}
		if in.AccountID != nil {
			if err := s.validateAccount(ctx, tenantID, userID, in.AccountID); err != nil {
				return nil, err
			}
			d.AccountID = in.AccountID
		}
		if in.LeadID != nil {
			d.LeadID = in.LeadID
		}
		if in.ContactID != nil {
			d.ContactID = in.ContactID
		}
		if in.OwnerID != nil {
			d.OwnerID = in.OwnerID
		}
		if in.Description != nil {
			d.Description = *in.Description
		}
		if in.LostReason != nil {
			d.LostReason = *in.LostReason
		}
		if in.TagsSet {
			d.Tags = domain.StringArray(in.Tags)
		}
	}
	d.UpdatedBy = userID

	if err := s.repo.Update(ctx, d); err != nil {
		return nil, err
	}
	dto := toDTO(d)
	return &dto, nil
}

func (s *Service) UpdateStage(ctx context.Context, tenantID, userID, id uuid.UUID, in StageInput) (*DealDTO, error) {
	scope := s.dataScope(ctx, tenantID, userID)
	d, err := s.repo.GetByID(ctx, tenantID, id, scope)
	if err != nil {
		return nil, err
	}
	if err := applyStageChange(d, in.Stage, in.LostReason); err != nil {
		return nil, err
	}
	d.UpdatedBy = userID
	if err := s.repo.Update(ctx, d); err != nil {
		return nil, err
	}
	dto := toDTO(d)
	return &dto, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, userID, id uuid.UUID) error {
	scope := s.dataScope(ctx, tenantID, userID)
	if _, err := s.repo.GetByID(ctx, tenantID, id, scope); err != nil {
		return err
	}
	return s.repo.SoftDelete(ctx, tenantID, id)
}

func (s *Service) Pipeline(ctx context.Context, tenantID, userID uuid.UUID, q PipelineQuery) (*PipelineResult, error) {
	scope := s.dataScope(ctx, tenantID, userID)
	byStage, summary, err := s.repo.Pipeline(ctx, tenantID, repository.DealPipelineFilter{
		OwnerID:   q.OwnerID,
		AccountID: q.AccountID,
		Scope:     scope,
		PerStage:  20,
	})
	if err != nil {
		return nil, err
	}
	stages := make([]PipelineStageDTO, 0, len(crm.DealPipelineStages))
	for _, stage := range crm.DealPipelineStages {
		deals := byStage[stage]
		var count int64
		var amountTotal float64
		for i := range deals {
			count++
			amountTotal += deals[i].Amount
		}
		items := make([]PipelineItemDTO, len(deals))
		for i := range deals {
			items[i] = toPipelineItem(&deals[i])
		}
		stages = append(stages, PipelineStageDTO{
			Stage:       stage,
			Count:       count,
			AmountTotal: amountTotal,
			Items:       items,
		})
	}
	return &PipelineResult{
		Stages: stages,
		Summary: PipelineSummaryDTO{
			OpenCount:    summary.OpenCount,
			OpenAmount:   summary.OpenAmount,
			WonCountMTD:  summary.WonCountMTD,
			WonAmountMTD: summary.WonAmountMTD,
		},
	}, nil
}

// CreateFromLeadInput is used by lead convert (optional create_deal).
type CreateFromLeadInput struct {
	Title  string
	Amount float64
	Stage  string
}

func (s *Service) CreateFromLead(ctx context.Context, tenantID, userID uuid.UUID, leadID, accountID uuid.UUID, in CreateFromLeadInput) (*DealDTO, error) {
	stage := in.Stage
	if stage == "" {
		stage = crm.DealStageQualification
	}
	return s.Create(ctx, tenantID, userID, CreateInput{
		Title:     in.Title,
		Stage:     stage,
		Amount:    in.Amount,
		LeadID:    &leadID,
		AccountID: &accountID,
	})
}

func applyStageChange(d *domain.Deal, to string, lostReason *string) error {
	if !crm.ValidDealStage(to) {
		return ErrInvalidStage
	}
	if crm.IsDealTerminal(d.Stage) && to != d.Stage {
		return ErrDealClosedReadonly
	}
	if to != d.Stage && !crm.CanTransitionDealStage(d.Stage, to) {
		return ErrInvalidStageTransition
	}
	if to == d.Stage {
		return nil
	}
	d.Stage = to
	if to == crm.DealStageWon {
		now := time.Now().UTC()
		d.ClosedAt = &now
		d.Probability = 100
	} else if to == crm.DealStageLost {
		now := time.Now().UTC()
		d.ClosedAt = &now
		d.Probability = 0
		if lostReason != nil {
			d.LostReason = *lostReason
		}
	}
	return nil
}

func (s *Service) validateAccount(ctx context.Context, tenantID, userID uuid.UUID, accountID *uuid.UUID) error {
	if accountID == nil || s.accounts == nil {
		return nil
	}
	_, err := s.accounts.GetByID(ctx, tenantID, *accountID, s.dataScope(ctx, tenantID, userID))
	if errors.Is(err, repository.ErrAccountNotFound) {
		return repository.ErrAccountNotFound
	}
	return err
}

func parseOptionalDate(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, errors.New("invalid expected_close_date")
	}
	return &t, nil
}

func parseDateRange(from, to *string) (*time.Time, *time.Time, error) {
	f, err := parseOptionalDate(from)
	if err != nil {
		return nil, nil, err
	}
	t, err := parseOptionalDate(to)
	if err != nil {
		return nil, nil, err
	}
	return f, t, nil
}

func toDTO(d *domain.Deal) DealDTO {
	tags := []string(d.Tags)
	if tags == nil {
		tags = []string{}
	}
	var expected *string
	if d.ExpectedCloseDate != nil {
		s := d.ExpectedCloseDate.Format("2006-01-02")
		expected = &s
	}
	return DealDTO{
		ID:                d.ID,
		TenantID:          d.TenantID,
		OwnerID:           d.OwnerID,
		Title:             d.Title,
		Stage:             d.Stage,
		Amount:            d.Amount,
		Currency:          d.Currency,
		Probability:       int(d.Probability),
		ExpectedCloseDate: expected,
		AccountID:         d.AccountID,
		LeadID:            d.LeadID,
		ContactID:         d.ContactID,
		Description:       d.Description,
		Tags:              tags,
		LostReason:        d.LostReason,
		ClosedAt:          d.ClosedAt,
		EngagementScore:   int(d.EngagementScore),
		LastActivityAt:    d.LastActivityAt,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
	}
}

func toPipelineItem(d *domain.Deal) PipelineItemDTO {
	var expected *string
	if d.ExpectedCloseDate != nil {
		s := d.ExpectedCloseDate.Format("2006-01-02")
		expected = &s
	}
	return PipelineItemDTO{
		ID:                d.ID,
		Title:             d.Title,
		Amount:            d.Amount,
		Currency:          d.Currency,
		Probability:       int(d.Probability),
		ExpectedCloseDate: expected,
		AccountID:         d.AccountID,
		OwnerID:           d.OwnerID,
	}
}
