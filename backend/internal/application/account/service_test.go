package account

import (
	"context"
	"testing"

	"crm-backend/internal/domain"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

type mockAccountRepo struct {
	items map[uuid.UUID]*domain.Account
}

func (m *mockAccountRepo) List(ctx context.Context, tenantID uuid.UUID, f repository.AccountListFilter) ([]domain.Account, int64, error) {
	var out []domain.Account
	for _, a := range m.items {
		if a.TenantID != tenantID {
			continue
		}
		if !f.ViewAll && a.OwnerID != nil && *a.OwnerID != f.UserID {
			continue
		}
		out = append(out, *a)
	}
	return out, int64(len(out)), nil
}

func (m *mockAccountRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Account, error) {
	a, ok := m.items[id]
	if !ok || a.TenantID != tenantID {
		return nil, repository.ErrAccountNotFound
	}
	if !viewAll && a.OwnerID != nil && *a.OwnerID != userID {
		return nil, repository.ErrAccountNotFound
	}
	return a, nil
}

func (m *mockAccountRepo) Create(ctx context.Context, a *domain.Account) error {
	if m.items == nil {
		m.items = map[uuid.UUID]*domain.Account{}
	}
	m.items[a.ID] = a
	return nil
}

func (m *mockAccountRepo) Update(ctx context.Context, a *domain.Account) error {
	m.items[a.ID] = a
	return nil
}

func (m *mockAccountRepo) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	delete(m.items, id)
	return nil
}

func TestService_CreateDefaultsOwnerAndStage(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	repo := &mockAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	svc := NewService(repo, nil)

	dto, err := svc.Create(context.Background(), tenantID, userID, CreateInput{Name: "Acme"})
	if err != nil {
		t.Fatal(err)
	}
	if dto.LifecycleStage != "acquire" {
		t.Fatalf("lifecycle: %s", dto.LifecycleStage)
	}
	if dto.OwnerID == nil || *dto.OwnerID != userID {
		t.Fatalf("owner: %+v", dto.OwnerID)
	}
	if dto.RelationshipHealth != "low" {
		t.Fatalf("health: %s", dto.RelationshipHealth)
	}
}
