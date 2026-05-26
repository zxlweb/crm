package customfield

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

var (
	ErrNotFound      = errors.New("custom field not found")
	ErrKeyConflict   = crm.ErrCustomFieldKeyConflict
	ErrTypeInvalid   = crm.ErrCustomFieldTypeInvalid
)

type Service struct {
	repo repository.CustomFieldRepository
}

func NewService(repo repository.CustomFieldRepository) *Service {
	return &Service{repo: repo}
}

type FieldDTO struct {
	ID           uuid.UUID       `json:"id"`
	EntityType   string          `json:"entity_type"`
	FieldKey     string          `json:"field_key"`
	FieldLabel   json.RawMessage `json:"field_label"`
	FieldType    string          `json:"field_type"`
	Required     bool            `json:"required"`
	Options      json.RawMessage `json:"options,omitempty"`
	DefaultValue json.RawMessage `json:"default_value,omitempty"`
	DisplayOrder int             `json:"display_order"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

type CreateInput struct {
	EntityType   string
	FieldKey     string
	FieldLabel   json.RawMessage
	FieldType    string
	Required     bool
	Options      json.RawMessage
	DefaultValue json.RawMessage
	DisplayOrder int
}

type UpdateInput struct {
	FieldLabel   *json.RawMessage
	FieldType    *string
	Required     *bool
	Options      json.RawMessage
	DefaultValue json.RawMessage
	DisplayOrder *int
	IsActive     *bool
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, entityType string) ([]FieldDTO, error) {
	rows, err := s.repo.List(ctx, tenantID, repository.CustomFieldListFilter{
		EntityType: entityType,
		ActiveOnly: false,
	})
	if err != nil {
		return nil, err
	}
	out := make([]FieldDTO, 0, len(rows))
	for i := range rows {
		out = append(out, toDTO(&rows[i]))
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, tenantID uuid.UUID, in CreateInput) (*FieldDTO, error) {
	if !crm.ValidEntityType(in.EntityType) || !crm.ValidFieldKey(in.FieldKey) || !crm.ValidFieldType(in.FieldType) {
		return nil, ErrTypeInvalid
	}
	if err := crm.ValidateCustomFieldInput(in.FieldType, in.Required, in.Options, in.DefaultValue); err != nil {
		return nil, err
	}
	if in.DisplayOrder == 0 {
		in.DisplayOrder = 100
	}
	f := &domain.CustomField{
		TenantID:     tenantID,
		EntityType:   in.EntityType,
		FieldKey:     in.FieldKey,
		FieldLabel:   datatypes.JSON(in.FieldLabel),
		FieldType:    in.FieldType,
		Required:     in.Required,
		Options:      datatypes.JSON(in.Options),
		DefaultValue: datatypes.JSON(in.DefaultValue),
		DisplayOrder: in.DisplayOrder,
		IsActive:     true,
	}
	if err := s.repo.Create(ctx, f); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrKeyConflict
		}
		return nil, err
	}
	return toDTOPtr(f), nil
}

func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, in UpdateInput) (*FieldDTO, error) {
	f, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, repository.ErrCustomFieldNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	fieldType := f.FieldType
	if in.FieldType != nil {
		fieldType = *in.FieldType
	}
	if in.FieldLabel != nil {
		f.FieldLabel = datatypes.JSON(*in.FieldLabel)
	}
	if in.Required != nil {
		f.Required = *in.Required
	}
	if in.DisplayOrder != nil {
		f.DisplayOrder = *in.DisplayOrder
	}
	if in.IsActive != nil {
		f.IsActive = *in.IsActive
	}
	opts := []byte(f.Options)
	if len(in.Options) > 0 {
		opts = in.Options
		f.Options = datatypes.JSON(in.Options)
	}
	def := []byte(f.DefaultValue)
	if len(in.DefaultValue) > 0 {
		def = in.DefaultValue
		f.DefaultValue = datatypes.JSON(in.DefaultValue)
	}
	f.FieldType = fieldType
	if err := crm.ValidateCustomFieldInput(fieldType, f.Required, opts, def); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, f); err != nil {
		if errors.Is(err, repository.ErrCustomFieldNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	updated, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toDTOPtr(updated), nil
}

func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	err := s.repo.Deactivate(ctx, tenantID, id)
	if errors.Is(err, repository.ErrCustomFieldNotFound) {
		return ErrNotFound
	}
	return err
}

func toDTO(f *domain.CustomField) FieldDTO {
	return FieldDTO{
		ID:           f.ID,
		EntityType:   f.EntityType,
		FieldKey:     f.FieldKey,
		FieldLabel:   json.RawMessage(f.FieldLabel),
		FieldType:    f.FieldType,
		Required:     f.Required,
		Options:      json.RawMessage(f.Options),
		DefaultValue: json.RawMessage(f.DefaultValue),
		DisplayOrder: f.DisplayOrder,
		IsActive:     f.IsActive,
		CreatedAt:    f.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    f.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toDTOPtr(f *domain.CustomField) *FieldDTO {
	d := toDTO(f)
	return &d
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique")
}
