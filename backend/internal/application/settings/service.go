package settings

import (
	"context"
	"errors"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

var ErrTenantNotFound = errors.New("tenant not found")

type Service struct {
	repo repository.SettingsRepository
}

func NewService(repo repository.SettingsRepository) *Service {
	return &Service{repo: repo}
}

type TenantDTO struct {
	TenantID         uuid.UUID            `json:"tenant_id"`
	TenantName       string               `json:"tenant_name"`
	DefaultLocale    string               `json:"default_locale"`
	Timezone         string               `json:"timezone"`
	BusinessSwitches crm.BusinessSwitches `json:"business_switches"`
	SalesQuota       crm.SalesQuota       `json:"sales_quota"`
	UpdatedAt        string               `json:"updated_at"`
	UpdatedBy        *uuid.UUID           `json:"updated_by,omitempty"`
}

type PatchInput struct {
	TenantName       *string
	DefaultLocale    *string
	Timezone         *string
	BusinessSwitches *crm.BusinessSwitchesPatch
	SalesQuota       *crm.SalesQuota
	UpdatedBy        uuid.UUID
}

func (s *Service) GetTenant(ctx context.Context, tenantID uuid.UUID) (*TenantDTO, error) {
	t, err := s.repo.GetTenant(ctx, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrTenantNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return toDTO(t, nil), nil
}

func (s *Service) PatchTenant(ctx context.Context, tenantID uuid.UUID, in PatchInput) (*TenantDTO, error) {
	t, err := s.repo.GetTenant(ctx, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrTenantNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	patch := crm.TenantSettingsPatch{
		DefaultLocale:    in.DefaultLocale,
		Timezone:         in.Timezone,
		BusinessSwitches: in.BusinessSwitches,
		SalesQuota:       in.SalesQuota,
	}
	cfg, err := crm.MergeTenantConfig(t.Config, in.TenantName, patch)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateTenant(ctx, tenantID, in.TenantName, cfg); err != nil {
		if errors.Is(err, repository.ErrTenantNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	updated, err := s.repo.GetTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return toDTO(updated, &in.UpdatedBy), nil
}

func (s *Service) ListFeatures(ctx context.Context, tenantID uuid.UUID) ([]map[string]any, error) {
	t, err := s.repo.GetTenant(ctx, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrTenantNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	settings := crm.ParseTenantSettings(t.Config)
	return crm.SettingsFeatureCatalog(settings), nil
}

func toDTO(t *domain.Tenant, updatedBy *uuid.UUID) *TenantDTO {
	settings := crm.ParseTenantSettings(t.Config)
	period := settings.SalesQuota.Period
	if period == "" {
		period = time.Now().Format("2006-01")
		settings.SalesQuota.Period = period
	}
	return &TenantDTO{
		TenantID:         t.ID,
		TenantName:       t.Name,
		DefaultLocale:    settings.DefaultLocale,
		Timezone:         settings.Timezone,
		BusinessSwitches: settings.BusinessSwitches,
		SalesQuota:       settings.SalesQuota,
		UpdatedAt:        t.UpdatedAt.UTC().Format(time.RFC3339),
		UpdatedBy:        updatedBy,
	}
}
