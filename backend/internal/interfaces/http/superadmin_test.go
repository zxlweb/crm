package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"crm-backend/internal/application/superadmin"
	"crm-backend/internal/domain"
	"crm-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockSuperTenantRepo struct{}

func (m *mockSuperTenantRepo) List(ctx context.Context, filter repository.TenantListFilter) ([]domain.Tenant, int64, error) {
	return []domain.Tenant{
		{ID: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), Name: "Demo", Domain: "demo", IsActive: true, CreatedAt: time.Now()},
	}, 1, nil
}

func (m *mockSuperTenantRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	return &domain.Tenant{ID: id, Name: "Demo", Domain: "demo", IsActive: true, CreatedAt: time.Now()}, nil
}

func (m *mockSuperTenantRepo) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	return nil
}

func (m *mockSuperTenantRepo) CountUsers(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return 1, nil
}

func (m *mockSuperTenantRepo) CountAll(ctx context.Context) (int64, int64, error) {
	return 1, 1, nil
}

func (m *mockSuperTenantRepo) CountAllUsers(ctx context.Context) (int64, error) {
	return 1, nil
}

func TestSuperAdminOverview_ForbiddenForNonAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := superadmin.NewService(&mockSuperTenantRepo{})
	h := NewSuperAdminHandlers(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("is_super_admin", false)
		c.Next()
	})
	r.GET("/api/super-admin/overview", superAdminGate(), h.Overview)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/super-admin/overview", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func TestSuperAdminOverview_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := superadmin.NewService(&mockSuperTenantRepo{})
	h := NewSuperAdminHandlers(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("is_super_admin", true)
		c.Next()
	})
	r.GET("/api/super-admin/overview", h.Overview)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/super-admin/overview", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func superAdminGate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool("is_super_admin") {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "需要 Super Admin 权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
