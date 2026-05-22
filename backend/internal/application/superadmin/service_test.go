package superadmin

import (
	"context"
	"testing"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

type mockTenantRepo struct {
	tenants []domain.Tenant
}

func (m *mockTenantRepo) List(ctx context.Context, filter repository.TenantListFilter) ([]domain.Tenant, int64, error) {
	return m.tenants, int64(len(m.tenants)), nil
}

func (m *mockTenantRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	for i := range m.tenants {
		if m.tenants[i].ID == id {
			return &m.tenants[i], nil
		}
	}
	return nil, repository.ErrTenantNotFound
}

func (m *mockTenantRepo) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	for i := range m.tenants {
		if m.tenants[i].ID == id {
			m.tenants[i].IsActive = active
			return nil
		}
	}
	return repository.ErrTenantNotFound
}

func (m *mockTenantRepo) CountUsers(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return 2, nil
}

func (m *mockTenantRepo) CountAll(ctx context.Context) (int64, int64, error) {
	var active int64
	for _, t := range m.tenants {
		if t.IsActive {
			active++
		}
	}
	return int64(len(m.tenants)), active, nil
}

func (m *mockTenantRepo) CountAllUsers(ctx context.Context) (int64, error) {
	return 5, nil
}

func (m *mockTenantRepo) TenantActivityTrend(ctx context.Context, days int) ([]repository.TenantActivityPoint, error) {
	return []repository.TenantActivityPoint{
		{Date: "05-20", NewTenants: 1, Logins: 2},
	}, nil
}

func TestService_Overview(t *testing.T) {
	repo := &mockTenantRepo{
		tenants: []domain.Tenant{
			{ID: uuid.New(), Name: "A", Domain: "a", IsActive: true, CreatedAt: time.Now()},
			{ID: uuid.New(), Name: "B", Domain: "b", IsActive: false, CreatedAt: time.Now()},
		},
	}
	svc := NewService(repo)

	out, err := svc.Overview(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if out.TenantCount != 2 || out.ActiveTenantCount != 1 || out.UserCount != 5 {
		t.Fatalf("overview: %+v", out)
	}
}

func TestService_TenantActivityTrend(t *testing.T) {
	svc := NewService(&mockTenantRepo{})
	out, err := svc.TenantActivityTrend(context.Background(), 7)
	if err != nil || len(out.Categories) != 1 || len(out.Series) != 2 {
		t.Fatalf("trend: %+v err=%v", out, err)
	}
}

func TestService_SetTenantActive(t *testing.T) {
	id := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	repo := &mockTenantRepo{
		tenants: []domain.Tenant{
			{ID: id, Name: "Demo", Domain: "demo", IsActive: true, CreatedAt: time.Now()},
		},
	}
	svc := NewService(repo)

	dto, err := svc.SetTenantActive(context.Background(), id, false)
	if err != nil || dto.IsActive {
		t.Fatalf("dto: %+v err=%v", dto, err)
	}
}
