package auth

import (
	"context"
	"testing"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/password"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

type mockUserRepo struct {
	user    *domain.User
	tenants []repository.TenantBrief
}

func (m *mockUserRepo) FindByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	if m.user != nil && m.user.ID == userID {
		return m.user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.user == nil || m.user.Email != email {
		return nil, repository.ErrUserNotFound
	}
	return m.user, nil
}

func (m *mockUserRepo) ListActiveTenantsForUser(ctx context.Context, userID uuid.UUID) ([]repository.TenantBrief, error) {
	return m.tenants, nil
}

func (m *mockUserRepo) UserBelongsToTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error) {
	for _, t := range m.tenants {
		if t.ID == tenantID {
			return true, nil
		}
	}
	return false, nil
}

func TestService_LoginSuccess(t *testing.T) {
	hash, _ := password.Hash("password123")
	uid := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	tid := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	repo := &mockUserRepo{
		user: &domain.User{
			ID:           uid,
			Email:        "admin@demo.com",
			PasswordHash: hash,
			IsSuperAdmin: true,
		},
		tenants: []repository.TenantBrief{
			{ID: tid, Name: "Demo Corp", Domain: "demo"},
		},
	}
	svc := NewService(repo, "secret", time.Hour, 24*time.Hour)

	result, err := svc.Login(context.Background(), "admin@demo.com", "password123")
	if err != nil {
		t.Fatal(err)
	}
	if result.AccessToken == "" || result.RefreshToken == "" || len(result.Tenants) != 1 {
		t.Fatalf("result: %+v", result)
	}
	if !result.User.IsSuperAdmin {
		t.Fatalf("user: %+v", result.User)
	}
}

func TestService_Profile(t *testing.T) {
	uid := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	repo := &mockUserRepo{
		user: &domain.User{ID: uid, Email: "admin@demo.com", IsSuperAdmin: true},
	}
	svc := NewService(repo, "secret", time.Hour, 24*time.Hour)

	profile, err := svc.Profile(context.Background(), uid)
	if err != nil || !profile.IsSuperAdmin {
		t.Fatalf("profile: %+v err=%v", profile, err)
	}
}

func TestService_LoginInvalidPassword(t *testing.T) {
	hash, _ := password.Hash("password123")
	repo := &mockUserRepo{
		user: &domain.User{ID: uuid.New(), Email: "a@b.com", PasswordHash: hash},
	}
	svc := NewService(repo, "secret", time.Hour, 24*time.Hour)

	_, err := svc.Login(context.Background(), "a@b.com", "wrong")
	if err != ErrInvalidCredentials {
		t.Fatalf("got %v", err)
	}
}

func TestService_SwitchTenantSuccess(t *testing.T) {
	uid := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	tid := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	repo := &mockUserRepo{
		user:    &domain.User{ID: uid, Email: "admin@demo.com", IsSuperAdmin: true},
		tenants: []repository.TenantBrief{{ID: tid, Name: "Demo", Domain: "demo"}},
	}
	svc := NewService(repo, "secret", time.Hour, 24*time.Hour)

	result, err := svc.SwitchTenant(context.Background(), uid, "admin@demo.com", true, tid)
	if err != nil {
		t.Fatal(err)
	}
	if result.CurrentTenant == nil || result.CurrentTenant.ID != tid.String() {
		t.Fatalf("current tenant: %+v", result.CurrentTenant)
	}
}
