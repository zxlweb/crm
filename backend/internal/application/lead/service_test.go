package lead

import (
	"context"
	"testing"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

type mockLeadRepo struct {
	items map[uuid.UUID]*domain.Lead
}

func (m *mockLeadRepo) List(ctx context.Context, tenantID uuid.UUID, f repository.LeadListFilter) ([]domain.Lead, int64, error) {
	var out []domain.Lead
	for _, l := range m.items {
		if l.TenantID == tenantID {
			out = append(out, *l)
		}
	}
	return out, int64(len(out)), nil
}

func (m *mockLeadRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Lead, error) {
	l, ok := m.items[id]
	if !ok || l.TenantID != tenantID {
		return nil, repository.ErrLeadNotFound
	}
	return l, nil
}

func (m *mockLeadRepo) Create(ctx context.Context, l *domain.Lead) error {
	if m.items == nil {
		m.items = map[uuid.UUID]*domain.Lead{}
	}
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	m.items[l.ID] = l
	return nil
}

func (m *mockLeadRepo) Update(ctx context.Context, l *domain.Lead) error {
	m.items[l.ID] = l
	return nil
}

func (m *mockLeadRepo) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	delete(m.items, id)
	return nil
}

func (m *mockLeadRepo) UpdateEngagementFromActivity(ctx context.Context, tenantID, id, updatedBy uuid.UUID, last *time.Time, score int16) error {
	l, ok := m.items[id]
	if !ok || l.TenantID != tenantID {
		return repository.ErrLeadNotFound
	}
	l.LastActivityAt = last
	l.EngagementScore = score
	return nil
}

func (m *mockLeadRepo) StatsBySource(ctx context.Context, tenantID uuid.UUID, f repository.LeadStatsFilter) ([]repository.LabelCount, int64, error) {
	return nil, 0, nil
}

func (m *mockLeadRepo) StatsByStatus(ctx context.Context, tenantID uuid.UUID, f repository.LeadStatsFilter) ([]repository.LabelCount, int64, error) {
	return nil, 0, nil
}

func (m *mockLeadRepo) StatsTrend(ctx context.Context, tenantID uuid.UUID, f repository.LeadStatsFilter, granularity string) ([]repository.TrendPoint, error) {
	return nil, nil
}

func (m *mockLeadRepo) StatsFunnel(ctx context.Context, tenantID uuid.UUID, f repository.LeadStatsFilter) ([]repository.LabelCount, error) {
	return nil, nil
}

func (m *mockLeadRepo) CountScoped(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	_ = viewAll
	_ = userID
	var n int64
	for _, l := range m.items {
		if l.TenantID == tenantID {
			n++
		}
	}
	return n, nil
}

func (m *mockLeadRepo) DailyCreatedCounts(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID, days int) ([]int64, error) {
	_ = ctx
	_ = tenantID
	_ = viewAll
	_ = userID
	if days < 1 {
		days = 7
	}
	return make([]int64, days), nil
}

func (m *mockLeadRepo) CountLowEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	_ = ctx
	_ = tenantID
	_ = viewAll
	_ = userID
	return 0, nil
}

func (m *mockLeadRepo) AvgEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (float64, error) {
	_ = ctx
	_ = tenantID
	_ = viewAll
	_ = userID
	return 0, nil
}

func (m *mockLeadRepo) ListPriorityCandidates(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID, limit int) ([]domain.Lead, error) {
	_ = ctx
	_ = limit
	items, _, err := m.List(ctx, tenantID, repository.LeadListFilter{Page: 1, PageSize: 20, ViewAll: viewAll, UserID: userID})
	return items, err
}

type mockAccountRepo struct {
	items map[uuid.UUID]*domain.Account
}

func (m *mockAccountRepo) List(ctx context.Context, tenantID uuid.UUID, f repository.AccountListFilter) ([]domain.Account, int64, error) {
	return nil, 0, nil
}

func (m *mockAccountRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Account, error) {
	a, ok := m.items[id]
	if !ok || a.TenantID != tenantID {
		return nil, repository.ErrAccountNotFound
	}
	return a, nil
}

func (m *mockAccountRepo) Create(ctx context.Context, a *domain.Account) error {
	if m.items == nil {
		m.items = map[uuid.UUID]*domain.Account{}
	}
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
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

func (m *mockAccountRepo) UpdateEngagementFromActivity(ctx context.Context, tenantID, id, updatedBy uuid.UUID, last *time.Time, score int16) error {
	a, ok := m.items[id]
	if !ok || a.TenantID != tenantID {
		return repository.ErrAccountNotFound
	}
	a.LastActivityAt = last
	a.EngagementScore = score
	return nil
}

func (m *mockAccountRepo) CountScoped(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	_ = ctx
	_ = viewAll
	_ = userID
	return int64(len(m.items)), nil
}

func (m *mockAccountRepo) CountLowEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	_ = ctx
	_ = tenantID
	_ = viewAll
	_ = userID
	return 0, nil
}

func TestService_StatusTransition(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	leadID := uuid.New()
	repo := &mockLeadRepo{items: map[uuid.UUID]*domain.Lead{
		leadID: {
			ID: leadID, TenantID: tenantID, OwnerID: &userID,
			Title: "L", Status: "new", LifecycleStage: "acquire",
		},
	}}
	svc := NewService(repo, &mockAccountRepo{}, nil, nil, nil, nil)

	status := "qualified"
	_, err := svc.Update(context.Background(), tenantID, userID, leadID, UpdateInput{Status: &status}, false)
	if err != ErrInvalidStatusTransition {
		t.Fatalf("expected invalid_status_transition, got %v", err)
	}

	status = "contacted"
	if _, err := svc.Update(context.Background(), tenantID, userID, leadID, UpdateInput{Status: &status}, false); err != nil {
		t.Fatalf("new->contacted: %v", err)
	}
}

func TestService_ConvertCreatesAccount(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	leadID := uuid.New()
	leadRepo := &mockLeadRepo{items: map[uuid.UUID]*domain.Lead{
		leadID: {
			ID: leadID, TenantID: tenantID, OwnerID: &userID,
			Title: "Deal", Status: "qualified", LifecycleStage: "grow",
		},
	}}
	accRepo := &mockAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	svc := NewService(leadRepo, accRepo, nil, nil, nil, nil)

	dto, err := svc.Convert(context.Background(), tenantID, userID, leadID, ConvertInput{
		CreateAccount: &ConvertAccountInput{Name: "Acme Corp"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if dto.LeadDTO.Status != "converted" || dto.ConvertedAccountID == nil {
		t.Fatalf("convert result: %+v", dto)
	}
	if len(accRepo.items) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accRepo.items))
	}
}
