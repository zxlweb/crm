package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"crm-backend/internal/application/auth"
	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/password"
	"crm-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	hash, _ := password.Hash("password123")
	uid := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	repo := &mockAuthRepo{
		user: &domain.User{ID: uid, Email: "admin@demo.com", PasswordHash: hash, IsSuperAdmin: true},
		tenants: []repository.TenantBrief{
			{ID: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), Name: "Demo", Domain: "demo"},
		},
	}
	svc := auth.NewService(repo, "test-secret", time.Hour, 24*time.Hour, nil)
	h := NewAuthHandlers(svc)

	r := gin.New()
	r.POST("/api/auth/login", h.Login)

	body, _ := json.Marshal(map[string]string{
		"email":    "admin@demo.com",
		"password": "password123",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

type mockAuthRepo struct {
	user    *domain.User
	tenants []repository.TenantBrief
}

func (m *mockAuthRepo) FindByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	if m.user != nil && m.user.ID == userID {
		return m.user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockAuthRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.user == nil || m.user.Email != email {
		return nil, repository.ErrUserNotFound
	}
	return m.user, nil
}

func (m *mockAuthRepo) ListActiveTenantsForUser(ctx context.Context, userID uuid.UUID) ([]repository.TenantBrief, error) {
	return m.tenants, nil
}

func (m *mockAuthRepo) UserBelongsToTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error) {
	for _, t := range m.tenants {
		if t.ID == tenantID {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockAuthRepo) RegisterWithTenant(ctx context.Context, in repository.RegisterInput) (*domain.User, uuid.UUID, error) {
	return nil, uuid.Nil, repository.ErrEmailExists
}

func TestSwitchTenantHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	hash, _ := password.Hash("password123")
	uid := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	tid := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	repo := &mockAuthRepo{
		user: &domain.User{ID: uid, Email: "admin@demo.com", PasswordHash: hash, IsSuperAdmin: true},
		tenants: []repository.TenantBrief{
			{ID: tid, Name: "Demo", Domain: "demo"},
		},
	}
	svc := auth.NewService(repo, "test-secret", time.Hour, 24*time.Hour, nil)
	h := NewAuthHandlers(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uid.String())
		c.Set("email", "admin@demo.com")
		c.Set("is_super_admin", true)
		c.Next()
	})
	r.POST("/api/auth/switch-tenant", h.SwitchTenant)

	body, _ := json.Marshal(map[string]string{"tenant_id": tid.String()})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/switch-tenant", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			CurrentTenant struct {
				ID string `json:"id"`
			} `json:"current_tenant"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Data.CurrentTenant.ID != tid.String() {
		t.Fatalf("expected current tenant %s, got %+v", tid, resp.Data)
	}
}
