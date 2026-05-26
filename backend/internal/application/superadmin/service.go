package superadmin

import (
	"context"
	"errors"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

var ErrTenantNotFound = errors.New("tenant not found")

type OverviewDTO struct {
	TenantCount       int64 `json:"tenant_count"`
	ActiveTenantCount int64 `json:"active_tenant_count"`
	UserCount         int64 `json:"user_count"`
}

type TenantActivityTrendDTO struct {
	Categories []string         `json:"categories"`
	Series     []TrendSeriesDTO `json:"series"`
}

type TrendSeriesDTO struct {
	Name    string  `json:"name"`
	Data    []int64 `json:"data"`
	Primary bool    `json:"primary,omitempty"`
}

type TenantDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	IsActive  bool   `json:"is_active"`
	UserCount int64  `json:"user_count"`
	CreatedAt string `json:"created_at"`
}

type ListResult struct {
	Items []TenantDTO `json:"items"`
}

type Service struct {
	tenants  repository.TenantRepository
	insights repository.TenantInsightsRepository
}

func NewService(tenants repository.TenantRepository, insights repository.TenantInsightsRepository) *Service {
	return &Service{tenants: tenants, insights: insights}
}

func (s *Service) TenantActivityTrend(ctx context.Context, days int) (*TenantActivityTrendDTO, error) {
	points, err := s.tenants.TenantActivityTrend(ctx, days)
	if err != nil {
		return nil, err
	}
	categories := make([]string, 0, len(points))
	newData := make([]int64, 0, len(points))
	loginData := make([]int64, 0, len(points))
	for _, p := range points {
		categories = append(categories, p.Date)
		newData = append(newData, p.NewTenants)
		loginData = append(loginData, p.Logins)
	}
	return &TenantActivityTrendDTO{
		Categories: categories,
		Series: []TrendSeriesDTO{
			{Name: "logins", Data: loginData, Primary: true},
			{Name: "new_tenants", Data: newData},
		},
	}, nil
}

func (s *Service) Overview(ctx context.Context) (*OverviewDTO, error) {
	total, active, err := s.tenants.CountAll(ctx)
	if err != nil {
		return nil, err
	}
	users, err := s.tenants.CountAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return &OverviewDTO{
		TenantCount:       total,
		ActiveTenantCount: active,
		UserCount:         users,
	}, nil
}

func (s *Service) ListTenants(ctx context.Context, page, pageSize int, search string, isActive *bool) (*ListResult, int64, error) {
	rows, total, err := s.tenants.List(ctx, repository.TenantListFilter{
		Page: page, PageSize: pageSize, Search: search, IsActive: isActive,
	})
	if err != nil {
		return nil, 0, err
	}
	items := make([]TenantDTO, 0, len(rows))
	for _, t := range rows {
		dto, err := s.toTenantDTO(ctx, &t)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *dto)
	}
	return &ListResult{Items: items}, total, nil
}

func (s *Service) GetTenant(ctx context.Context, id uuid.UUID) (*TenantDTO, error) {
	t, err := s.tenants.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTenantNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return s.toTenantDTO(ctx, t)
}

func (s *Service) SetTenantActive(ctx context.Context, id uuid.UUID, active bool) (*TenantDTO, error) {
	if err := s.tenants.SetActive(ctx, id, active); err != nil {
		if errors.Is(err, repository.ErrTenantNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return s.GetTenant(ctx, id)
}

func (s *Service) toTenantDTO(ctx context.Context, t *domain.Tenant) (*TenantDTO, error) {
	userCount, err := s.tenants.CountUsers(ctx, t.ID)
	if err != nil {
		return nil, err
	}
	return &TenantDTO{
		ID:        t.ID.String(),
		Name:      t.Name,
		Domain:    t.Domain,
		IsActive:  t.IsActive,
		UserCount: userCount,
		CreatedAt: t.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}
