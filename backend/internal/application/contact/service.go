package contact

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/datascope"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
)

var (
	ErrNotFound          = repository.ErrContactNotFound
	ErrInvalidLifecycle  = errors.New("invalid lifecycle_stage")
	ErrNameRequired      = errors.New("name_or_email_required")
)

type Service struct {
	repo     repository.ContactRepository
	accounts repository.AccountRepository
	enforcer *casbin.Enforcer
}

func NewService(repo repository.ContactRepository, accounts repository.AccountRepository, enforcer *casbin.Enforcer) *Service {
	return &Service{repo: repo, accounts: accounts, enforcer: enforcer}
}

type ContactDTO struct {
	ID                 uuid.UUID  `json:"id"`
	AccountID          *uuid.UUID `json:"account_id,omitempty"`
	FirstName          string     `json:"first_name,omitempty"`
	LastName           string     `json:"last_name,omitempty"`
	DisplayName        string     `json:"display_name"`
	Email              string     `json:"email,omitempty"`
	Phone              string     `json:"phone,omitempty"`
	IsPrimary          bool       `json:"is_primary"`
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
	AccountID          *uuid.UUID
	OwnerID            *uuid.UUID
}

type ListResult struct {
	Items []ContactDTO
	Total int64
	Page  int
	Size  int
}

type CreateInput struct {
	AccountID      *uuid.UUID
	FirstName      string
	LastName       string
	Email          string
	Phone          string
	IsPrimary      bool
	OwnerID        *uuid.UUID
	LifecycleStage string
	Tags           []string
}

type UpdateInput struct {
	AccountID      *uuid.UUID
	AccountIDSet   bool
	FirstName      *string
	LastName       *string
	Email          *string
	Phone          *string
	IsPrimary      *bool
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
	items, total, err := s.repo.List(ctx, tenantID, repository.ContactListFilter{
		Page:               page,
		PageSize:           size,
		Search:             q.Search,
		LifecycleStage:     q.LifecycleStage,
		RelationshipHealth: q.RelationshipHealth,
		AccountID:          q.AccountID,
		OwnerID:            q.OwnerID,
		ViewAll:            viewAll,
		UserID:             userID,
	})
	if err != nil {
		return nil, err
	}
	dtos := make([]ContactDTO, len(items))
	for i := range items {
		dtos[i] = toDTO(&items[i])
	}
	return &ListResult{Items: dtos, Total: total, Page: page, Size: size}, nil
}

func (s *Service) ListByAccount(ctx context.Context, tenantID, userID, accountID uuid.UUID, q ListQuery) (*ListResult, error) {
	viewAll := s.viewAll(ctx, userID.String(), tenantID.String())
	if _, err := s.accounts.GetByID(ctx, tenantID, accountID, viewAll, userID); err != nil {
		return nil, err
	}
	q.AccountID = &accountID
	return s.List(ctx, tenantID, userID, q)
}

func (s *Service) Get(ctx context.Context, tenantID, userID, id uuid.UUID) (*ContactDTO, error) {
	c, err := s.repo.GetByID(ctx, tenantID, id, s.viewAll(ctx, userID.String(), tenantID.String()), userID)
	if err != nil {
		return nil, err
	}
	dto := toDTO(c)
	return &dto, nil
}

func (s *Service) Create(ctx context.Context, tenantID, userID uuid.UUID, in CreateInput) (*ContactDTO, error) {
	if !hasIdentity(in.FirstName, in.LastName, in.Email) {
		return nil, ErrNameRequired
	}
	stage := in.LifecycleStage
	if stage == "" {
		stage = "acquire"
	}
	if !crm.ValidLifecycleStage(stage) {
		return nil, ErrInvalidLifecycle
	}
	viewAll := s.viewAll(ctx, userID.String(), tenantID.String())
	if in.AccountID != nil {
		if _, err := s.accounts.GetByID(ctx, tenantID, *in.AccountID, viewAll, userID); err != nil {
			return nil, err
		}
	}
	owner := in.OwnerID
	if owner == nil {
		owner = &userID
	}
	c := &domain.Contact{
		TenantID:        tenantID,
		AccountID:       in.AccountID,
		OwnerID:         owner,
		FirstName:       strings.TrimSpace(in.FirstName),
		LastName:        strings.TrimSpace(in.LastName),
		Email:           strings.TrimSpace(in.Email),
		Phone:           strings.TrimSpace(in.Phone),
		IsPrimary:       in.IsPrimary,
		LifecycleStage:  stage,
		EngagementScore: 0,
		Tags:            domain.StringArray(in.Tags),
		AuditFields: domain.AuditFields{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
	}
	if err := s.repo.Create(ctx, c); err != nil {
		if err.Error() == "invalid lifecycle_stage" {
			return nil, ErrInvalidLifecycle
		}
		return nil, err
	}
	dto := toDTO(c)
	return &dto, nil
}

func (s *Service) Update(ctx context.Context, tenantID, userID, id uuid.UUID, in UpdateInput, full bool) (*ContactDTO, error) {
	viewAll := s.viewAll(ctx, userID.String(), tenantID.String())
	c, err := s.repo.GetByID(ctx, tenantID, id, viewAll, userID)
	if err != nil {
		return nil, err
	}

	if in.AccountIDSet {
		if in.AccountID != nil {
			if _, err := s.accounts.GetByID(ctx, tenantID, *in.AccountID, viewAll, userID); err != nil {
				return nil, err
			}
			c.AccountID = in.AccountID
		} else {
			c.AccountID = nil
		}
	}
	if in.FirstName != nil {
		c.FirstName = strings.TrimSpace(*in.FirstName)
	}
	if in.LastName != nil {
		c.LastName = strings.TrimSpace(*in.LastName)
	}
	if in.Email != nil {
		c.Email = strings.TrimSpace(*in.Email)
	}
	if in.Phone != nil {
		c.Phone = strings.TrimSpace(*in.Phone)
	}
	if in.IsPrimary != nil {
		c.IsPrimary = *in.IsPrimary
	}
	if in.OwnerID != nil {
		c.OwnerID = in.OwnerID
	}
	if in.LifecycleStage != nil {
		if !crm.ValidLifecycleStage(*in.LifecycleStage) {
			return nil, ErrInvalidLifecycle
		}
		c.LifecycleStage = *in.LifecycleStage
	}
	if in.TagsSet {
		c.Tags = domain.StringArray(in.Tags)
	}

	if full {
		if !hasIdentity(c.FirstName, c.LastName, c.Email) {
			return nil, ErrNameRequired
		}
	} else if !hasIdentity(c.FirstName, c.LastName, c.Email) {
		return nil, ErrNameRequired
	}

	c.UpdatedBy = userID
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	dto := toDTO(c)
	return &dto, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, userID, id uuid.UUID) error {
	viewAll := s.viewAll(ctx, userID.String(), tenantID.String())
	if _, err := s.repo.GetByID(ctx, tenantID, id, viewAll, userID); err != nil {
		return err
	}
	return s.repo.SoftDelete(ctx, tenantID, id)
}

func (s *Service) viewAll(ctx context.Context, userID, tenantID string) bool {
	return datascope.CanViewAllTenantData(ctx, s.enforcer, userID, tenantID)
}

func hasIdentity(first, last, email string) bool {
	return strings.TrimSpace(first) != "" || strings.TrimSpace(last) != "" || strings.TrimSpace(email) != ""
}

func displayName(c *domain.Contact) string {
	name := strings.TrimSpace(strings.TrimSpace(c.FirstName) + " " + strings.TrimSpace(c.LastName))
	if name != "" {
		return name
	}
	if c.Email != "" {
		return c.Email
	}
	if c.Phone != "" {
		return c.Phone
	}
	return "-"
}

func toDTO(c *domain.Contact) ContactDTO {
	tags := []string(c.Tags)
	if tags == nil {
		tags = []string{}
	}
	return ContactDTO{
		ID:                 c.ID,
		AccountID:          c.AccountID,
		FirstName:          c.FirstName,
		LastName:           c.LastName,
		DisplayName:        displayName(c),
		Email:              c.Email,
		Phone:              c.Phone,
		IsPrimary:          c.IsPrimary,
		OwnerID:            c.OwnerID,
		LifecycleStage:     c.LifecycleStage,
		RelationshipHealth: crm.RelationshipHealthFromScore(c.EngagementScore),
		EngagementScore:    int(c.EngagementScore),
		LastActivityAt:     c.LastActivityAt,
		Tags:               tags,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}
