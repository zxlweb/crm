package activity

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"crm-backend/internal/application/appscope"
	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/datascope"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

var (
	ErrNotFound        = repository.ErrActivityNotFound
	ErrSubjectNotFound = errors.New("subject_not_found")
	ErrInvalidEvent    = errors.New("invalid event_type")
	ErrInvalidSubject  = errors.New("invalid subject_type")
	ErrInvalidSentiment = errors.New("invalid sentiment")
	ErrInvalidSource   = errors.New("invalid sentiment_source")
	ErrInvalidDirection = errors.New("invalid direction")
	ErrSentimentAI     = errors.New("sentiment_source_ai_not_allowed")
)

type Service struct {
	repo     repository.ActivityRepository
	leads    repository.LeadRepository
	accounts repository.AccountRepository
	contacts repository.ContactRepository
	tenants  repository.TenantRepository
	enforcer *casbin.Enforcer
	scope    appscope.Provider
}

func NewService(
	repo repository.ActivityRepository,
	leads repository.LeadRepository,
	accounts repository.AccountRepository,
	contacts repository.ContactRepository,
	tenants repository.TenantRepository,
	enforcer *casbin.Enforcer,
	scope appscope.Provider,
) *Service {
	return &Service{repo: repo, leads: leads, accounts: accounts, contacts: contacts, tenants: tenants, enforcer: enforcer, scope: scope}
}

func (s *Service) dataScope(ctx context.Context, tenantID, userID uuid.UUID) datascope.ScopeParams {
	return s.scope.Params(ctx, tenantID, userID)
}

type ActivityDTO struct {
	ID              uuid.UUID      `json:"id"`
	SubjectType     string         `json:"subject_type"`
	SubjectID       uuid.UUID      `json:"subject_id"`
	EventType       string         `json:"event_type"`
	Direction       string         `json:"direction,omitempty"`
	Body            string         `json:"body,omitempty"`
	Metadata        map[string]any `json:"metadata"`
	Label           string         `json:"label,omitempty"`
	Sentiment       *string        `json:"sentiment,omitempty"`
	SentimentSource *string        `json:"sentiment_source,omitempty"`
	OccurredAt      time.Time      `json:"occurred_at"`
	CreatedBy       *uuid.UUID     `json:"created_by,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type ListQuery struct {
	SubjectType string
	SubjectID   uuid.UUID
	Page        int
	PageSize    int
}

type ListResult struct {
	Items []ActivityDTO
	Total int64
	Page  int
	Size  int
}

type CreateInput struct {
	SubjectType     string
	SubjectID       uuid.UUID
	EventType       string
	Direction       string
	Body            string
	Metadata        map[string]any
	Sentiment       *string
	SentimentSource string
	OccurredAt      *time.Time
	Label           string
}

type UpdateInput struct {
	Body            *string
	Direction       *string
	Metadata        map[string]any
	MetadataSet     bool
	Sentiment       *string
	SentimentClear  bool
	SentimentSource *string
	OccurredAt      *time.Time
	Label           *string
}

type SummaryDTO struct {
	Items []LabelCountDTO `json:"items"`
	Total int64           `json:"total"`
}

type LabelCountDTO struct {
	Label      string  `json:"label"`
	Value      int64   `json:"value"`
	Percentage float64 `json:"percentage,omitempty"`
}

func (s *Service) List(ctx context.Context, tenantID, userID uuid.UUID, q ListQuery) (*ListResult, error) {
	if !crm.ValidActivitySubjectType(q.SubjectType) || q.SubjectID == uuid.Nil {
		return nil, ErrInvalidSubject
	}
	if err := s.assertSubjectAccess(ctx, tenantID, userID, q.SubjectType, q.SubjectID); err != nil {
		return nil, err
	}
	page := q.Page
	if page < 1 {
		page = 1
	}
	size := q.PageSize
	if size < 1 {
		size = 20
	}
	items, total, err := s.repo.List(ctx, tenantID, repository.ActivityListFilter{
		SubjectType: q.SubjectType,
		SubjectID:   q.SubjectID,
		Page:        page,
		PageSize:    size,
	})
	if err != nil {
		return nil, err
	}
	dtos := make([]ActivityDTO, len(items))
	for i := range items {
		dtos[i] = toDTO(&items[i])
	}
	return &ListResult{Items: dtos, Total: total, Page: page, Size: size}, nil
}

func (s *Service) Get(ctx context.Context, tenantID, userID, id uuid.UUID) (*ActivityDTO, error) {
	a, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if err := s.assertSubjectAccess(ctx, tenantID, userID, a.SubjectType, a.SubjectID); err != nil {
		return nil, err
	}
	dto := toDTO(a)
	return &dto, nil
}

func (s *Service) Create(ctx context.Context, tenantID, userID uuid.UUID, in CreateInput) (*ActivityDTO, error) {
	if err := validateCreateBase(in); err != nil {
		return nil, err
	}
	if err := s.assertSubjectAccess(ctx, tenantID, userID, in.SubjectType, in.SubjectID); err != nil {
		return nil, err
	}
	occurred := time.Now().UTC()
	if in.OccurredAt != nil {
		occurred = in.OccurredAt.UTC()
	}
	meta := in.Metadata
	if meta == nil {
		meta = map[string]any{}
	}
	if in.Label != "" {
		meta["label"] = in.Label
	}
	src := in.SentimentSource
	sentiment := in.Sentiment
	if sentiment == nil && in.Body != "" {
		if s, ruleSrc, ok := s.inferSentiment(ctx, tenantID, in.Body); ok {
			sentiment = &s
			src = ruleSrc
		}
	}
	if src == "" && sentiment != nil {
		src = "manual"
	}
	if err := validateSentimentFields(sentiment, src); err != nil {
		return nil, err
	}
	a := &domain.Activity{
		TenantID:        tenantID,
		SubjectType:     in.SubjectType,
		SubjectID:       in.SubjectID,
		EventType:       in.EventType,
		Direction:       in.Direction,
		Body:            in.Body,
		Metadata:        mustJSON(meta),
		Sentiment:       sentiment,
		OccurredAt:      occurred,
		CreatedBy:       &userID,
	}
	if src != "" {
		a.SentimentSource = &src
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	if err := s.refreshSubject(ctx, tenantID, userID, in.SubjectType, in.SubjectID); err != nil {
		return nil, err
	}
	dto := toDTO(a)
	return &dto, nil
}

func (s *Service) Update(ctx context.Context, tenantID, userID, id uuid.UUID, in UpdateInput) (*ActivityDTO, error) {
	if err := validateUpdate(in); err != nil {
		return nil, err
	}
	a, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if err := s.assertSubjectAccess(ctx, tenantID, userID, a.SubjectType, a.SubjectID); err != nil {
		return nil, err
	}
	if in.Body != nil {
		a.Body = *in.Body
	}
	if in.Direction != nil {
		if !crm.ValidActivityDirection(*in.Direction) {
			return nil, ErrInvalidDirection
		}
		a.Direction = *in.Direction
	}
	if in.MetadataSet {
		meta := in.Metadata
		if meta == nil {
			meta = map[string]any{}
		}
		if in.Label != nil {
			meta["label"] = *in.Label
		}
		a.Metadata = mustJSON(meta)
	} else if in.Label != nil {
		meta := parseMetadata(a.Metadata)
		meta["label"] = *in.Label
		a.Metadata = mustJSON(meta)
	}
	if in.SentimentClear {
		a.Sentiment = nil
		a.SentimentSource = nil
	} else if in.Sentiment != nil {
		if !crm.ValidActivitySentiment(*in.Sentiment) {
			return nil, ErrInvalidSentiment
		}
		a.Sentiment = in.Sentiment
		if in.SentimentSource != nil {
			if *in.SentimentSource == "ai" {
				return nil, ErrSentimentAI
			}
			if !crm.ValidSentimentSource(*in.SentimentSource) {
				return nil, ErrInvalidSource
			}
			a.SentimentSource = in.SentimentSource
		} else if a.SentimentSource == nil || (a.SentimentSource != nil && *a.SentimentSource == "") {
			src := "manual"
			a.SentimentSource = &src
		}
	} else if in.Body != nil && *in.Body != "" && !isManualSentiment(a.SentimentSource) {
		if s, ruleSrc, ok := s.inferSentiment(ctx, tenantID, *in.Body); ok {
			a.Sentiment = &s
			a.SentimentSource = &ruleSrc
		}
	}
	if in.SentimentSource != nil && !in.SentimentClear && in.Sentiment == nil {
		if *in.SentimentSource == "ai" {
			return nil, ErrSentimentAI
		}
		if !crm.ValidSentimentSource(*in.SentimentSource) {
			return nil, ErrInvalidSource
		}
		a.SentimentSource = in.SentimentSource
	}
	if in.OccurredAt != nil {
		a.OccurredAt = in.OccurredAt.UTC()
	}
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	if err := s.refreshSubject(ctx, tenantID, userID, a.SubjectType, a.SubjectID); err != nil {
		return nil, err
	}
	dto := toDTO(a)
	return &dto, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, userID, id uuid.UUID) error {
	a, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if err := s.assertSubjectAccess(ctx, tenantID, userID, a.SubjectType, a.SubjectID); err != nil {
		return err
	}
	subjectType, subjectID := a.SubjectType, a.SubjectID
	if err := s.repo.SoftDelete(ctx, tenantID, id); err != nil {
		return err
	}
	return s.refreshSubject(ctx, tenantID, userID, subjectType, subjectID)
}

func (s *Service) Summary(ctx context.Context, tenantID, userID uuid.UUID, subjectType string, subjectID uuid.UUID) (*SummaryDTO, error) {
	if subjectType != "" && subjectID != uuid.Nil {
		if err := s.assertSubjectAccess(ctx, tenantID, userID, subjectType, subjectID); err != nil {
			return nil, err
		}
	}
	rows, total, err := s.repo.SummaryByEventType(ctx, tenantID, subjectType, subjectID)
	if err != nil {
		return nil, err
	}
	items := make([]LabelCountDTO, len(rows))
	for i, row := range rows {
		pct := 0.0
		if total > 0 {
			pct = float64(row.Count) / float64(total)
		}
		items[i] = LabelCountDTO{Label: row.Label, Value: row.Count, Percentage: pct}
	}
	return &SummaryDTO{Items: items, Total: total}, nil
}

func isManualSentiment(src *string) bool {
	return src != nil && *src == "manual"
}

func (s *Service) inferSentiment(ctx context.Context, tenantID uuid.UUID, body string) (sentiment, source string, ok bool) {
	rules := crm.DefaultSentimentKeywordRules()
	if s.tenants != nil {
		if t, err := s.tenants.FindByID(ctx, tenantID); err == nil && t != nil {
			rules = crm.ParseTenantCRMConfig(t.Config).SentimentKeywordRules
		}
	}
	sentiment, matched := crm.InferSentimentFromBody(body, rules)
	if !matched {
		return "", "", false
	}
	return sentiment, "rule", true
}

func validateCreateBase(in CreateInput) error {
	if !crm.ValidActivitySubjectType(in.SubjectType) || in.SubjectID == uuid.Nil {
		return ErrInvalidSubject
	}
	if !crm.ValidActivityEventType(in.EventType) {
		return ErrInvalidEvent
	}
	if !crm.ValidActivityDirection(in.Direction) {
		return ErrInvalidDirection
	}
	return nil
}

func validateSentimentFields(sentiment *string, src string) error {
	if sentiment != nil && !crm.ValidActivitySentiment(*sentiment) {
		return ErrInvalidSentiment
	}
	if src == "ai" {
		return ErrSentimentAI
	}
	if src != "" && !crm.ValidSentimentSource(src) {
		return ErrInvalidSource
	}
	if src == "rule" && sentiment == nil {
		return ErrInvalidSentiment
	}
	return nil
}

func validateUpdate(in UpdateInput) error {
	if in.Sentiment != nil && !crm.ValidActivitySentiment(*in.Sentiment) {
		return ErrInvalidSentiment
	}
	if in.SentimentSource != nil {
		if *in.SentimentSource == "ai" {
			return ErrSentimentAI
		}
		if !crm.ValidSentimentSource(*in.SentimentSource) {
			return ErrInvalidSource
		}
		if *in.SentimentSource == "rule" && in.Sentiment == nil && !in.SentimentClear {
			return ErrInvalidSentiment
		}
	}
	return nil
}

func (s *Service) assertSubjectAccess(ctx context.Context, tenantID, userID uuid.UUID, subjectType string, subjectID uuid.UUID) error {
	scope := s.dataScope(ctx, tenantID, userID)
	switch subjectType {
	case "lead":
		_, err := s.leads.GetByID(ctx, tenantID, subjectID, scope)
		if errors.Is(err, repository.ErrLeadNotFound) {
			return ErrSubjectNotFound
		}
		return err
	case "account":
		_, err := s.accounts.GetByID(ctx, tenantID, subjectID, scope)
		if errors.Is(err, repository.ErrAccountNotFound) {
			return ErrSubjectNotFound
		}
		return err
	case "contact":
		_, err := s.contacts.GetByID(ctx, tenantID, subjectID, scope)
		if errors.Is(err, repository.ErrContactNotFound) {
			return ErrSubjectNotFound
		}
		return err
	default:
		return ErrInvalidSubject
	}
}

func (s *Service) refreshSubject(ctx context.Context, tenantID, userID uuid.UUID, subjectType string, subjectID uuid.UUID) error {
	latest, err := s.repo.LatestOccurredAt(ctx, tenantID, subjectType, subjectID)
	if err != nil {
		return err
	}
	count, err := s.repo.CountBySubject(ctx, tenantID, subjectType, subjectID)
	if err != nil {
		return err
	}
	score := int16(count * 3)
	if score > 100 {
		score = 100
	}
	scope := s.dataScope(ctx, tenantID, userID)
	switch subjectType {
	case "lead":
		if _, err := s.leads.GetByID(ctx, tenantID, subjectID, scope); err != nil {
			return err
		}
		return s.leads.UpdateEngagementFromActivity(ctx, tenantID, subjectID, userID, latest, score)
	case "account":
		if _, err := s.accounts.GetByID(ctx, tenantID, subjectID, scope); err != nil {
			return err
		}
		return s.accounts.UpdateEngagementFromActivity(ctx, tenantID, subjectID, userID, latest, score)
	case "contact":
		if _, err := s.contacts.GetByID(ctx, tenantID, subjectID, scope); err != nil {
			return err
		}
		return s.contacts.UpdateEngagementFromActivity(ctx, tenantID, subjectID, userID, latest, score)
	default:
		return nil
	}
}

func toDTO(a *domain.Activity) ActivityDTO {
	meta := parseMetadata(a.Metadata)
	label, _ := meta["label"].(string)
	return ActivityDTO{
		ID:              a.ID,
		SubjectType:     a.SubjectType,
		SubjectID:       a.SubjectID,
		EventType:       a.EventType,
		Direction:       a.Direction,
		Body:            a.Body,
		Metadata:        meta,
		Label:           label,
		Sentiment:       a.Sentiment,
		SentimentSource: a.SentimentSource,
		OccurredAt:      a.OccurredAt,
		CreatedBy:       a.CreatedBy,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

func parseMetadata(raw datatypes.JSON) map[string]any {
	if len(raw) == 0 {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil || m == nil {
		return map[string]any{}
	}
	return m
}

func mustJSON(m map[string]any) datatypes.JSON {
	b, err := json.Marshal(m)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return datatypes.JSON(b)
}
