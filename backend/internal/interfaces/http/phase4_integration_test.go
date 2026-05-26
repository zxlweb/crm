package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	auditstatsapp "crm-backend/internal/application/auditstats"
	customfieldapp "crm-backend/internal/application/customfield"
	settingsapp "crm-backend/internal/application/settings"
	"crm-backend/internal/application/audit"
	"crm-backend/internal/application/superadmin"
	"crm-backend/internal/domain"
	httphandler "crm-backend/internal/interfaces/http"
	"crm-backend/internal/interfaces/middleware"
	"crm-backend/internal/pkg/jwtutil"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func TestPhase4_SettingsAndCustomFields(t *testing.T) {
	env := setupPhase4HTTPEnv(t)

	w := env.get("/api/settings/tenant")
	env.assertOK(w)
	data := env.parseData(w)
	if data["default_locale"] != "zh-CN" {
		t.Fatalf("locale: %v", data["default_locale"])
	}

	w = env.patch("/api/settings/tenant", map[string]any{
		"tenant_name":    "Acme China",
		"default_locale": "en-US",
		"business_switches": map[string]any{
			"ai_preview_enabled": true,
		},
	})
	env.assertOK(w)
	env.assertAuditAction("settings.update", 1)

	w = env.get("/api/settings/features")
	env.assertOK(w)

	w = env.post("/api/settings/custom-fields", map[string]any{
		"entity_type": "lead",
		"field_key":   "industry_segment",
		"field_label": map[string]string{"zh-CN": "行业子类", "en-US": "Segment"},
		"field_type":  "select",
		"options": []map[string]any{
			{"value": "saas", "label": map[string]string{"zh-CN": "SaaS", "en-US": "SaaS"}},
		},
		"display_order": 30,
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create field: %d %s", w.Code, w.Body.String())
	}
	fieldID := env.parseDataID(w)
	env.assertAuditAction("custom_field.create", 1)

	w = env.get("/api/settings/custom-fields?entity_type=lead")
	env.assertOK(w)
	body := env.parseBody(w)
	items := body["data"].(map[string]any)["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("fields len %d", len(items))
	}

	w = env.request(http.MethodDelete, "/api/settings/custom-fields/"+fieldID, env.tenantA, nil)
	env.assertOK(w)
}

func TestPhase4_AuditStatsAndExport(t *testing.T) {
	env := setupPhase4HTTPEnv(t)
	tenantID := uuid.MustParse(env.tenantA)
	userID := uuid.MustParse(env.userID)
	env.auditStats.logs = append(env.auditStats.logs, repository.AuditLogExportRow{
		ID: uuid.New(), CreatedAt: time.Now(), Action: "settings.update",
		ResourceType: "settings", UserID: &userID,
	})
	env.auditStats.byAction = []repository.ActionCountRow{{Action: "settings.update", Count: 2}}

	w := env.get("/api/audit/stats/by-action")
	env.assertOK(w)

	w = env.get("/api/audit/export?format=csv")
	if w.Code != http.StatusOK {
		t.Fatalf("export: %d %s", w.Code, w.Body.String())
	}
	if w.Header().Get("Content-Type") == "" {
		t.Fatal("missing content type")
	}

	w = env.get("/api/audit/export?format=csv")
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("rate limit: %d", w.Code)
	}
	_ = tenantID
}

func TestPhase4_SuperAdminInsights(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := superadmin.NewServiceWithInsights(&memInsightsTenantRepo{}, &memInsightsRepo{})
	h := httphandler.NewSuperAdminHandlers(svc, nil)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("is_super_admin", true)
		c.Next()
	})
	r.GET("/api/super-admin/stats/tenant-health", h.TenantHealth)
	r.GET("/api/super-admin/stats/plan-distribution", h.PlanDistribution)
	r.GET("/api/super-admin/stats/top-tenants", h.TopTenants)

	for _, path := range []string{
		"/api/super-admin/stats/tenant-health",
		"/api/super-admin/stats/plan-distribution",
		"/api/super-admin/stats/top-tenants?metric=activity",
	} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
		}
	}
}

func TestPhase4_TenantIsolation(t *testing.T) {
	env := setupPhase4HTTPEnv(t)
	w := env.post("/api/settings/custom-fields", map[string]any{
		"entity_type": "lead",
		"field_key":   "tenant_a_only",
		"field_label": map[string]string{"zh-CN": "A"},
		"field_type":  "text",
	})
	fieldID := env.parseDataID(w)

	w = env.request(http.MethodGet, "/api/settings/custom-fields/"+fieldID, env.tenantB, nil)
	// list is filtered; direct get not exposed �?use list
	w = env.request(http.MethodGet, "/api/settings/custom-fields?entity_type=lead", env.tenantB, nil)
	env.assertOK(w)
	body := env.parseBody(w)
	data := body["data"].(map[string]any)
	items, _ := data["items"].([]any)
	if len(items) != 0 {
		t.Fatalf("tenant B should not see tenant A fields")
	}
}

type phase4HTTPEnv struct {
	t          *testing.T
	router     *gin.Engine
	auditRepo  *memAuditRepo
	auditStats *memAuditStatsRepo
	tenantA    string
	tenantB    string
	userID     string
	token      string
}

func setupPhase4HTTPEnv(t *testing.T) *phase4HTTPEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)
	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	tenantB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	enforcer := newPhase4Enforcer(t, userA, tenantA, tenantB)
	settingsRepo := &memSettingsRepo{tenants: map[uuid.UUID]*domain.Tenant{
		tenantA: {
			ID: tenantA, Name: "Demo", Config: datatypes.JSON([]byte(`{"default_locale":"zh-CN","timezone":"Asia/Shanghai"}`)),
			UpdatedAt: time.Now(),
		},
		tenantB: {ID: tenantB, Name: "Other", Config: datatypes.JSON([]byte(`{}`)), UpdatedAt: time.Now()},
	}}
	customRepo := &memCustomFieldRepo{items: map[uuid.UUID]*domain.CustomField{}}
	auditRepo := &memAuditRepo{}
	auditStatsRepo := &memAuditStatsRepo{}

	settingsSvc := settingsapp.NewService(settingsRepo)
	customSvc := customfieldapp.NewService(customRepo)
	auditStatsSvc := auditstatsapp.NewService(auditStatsRepo)
	auditRec := audit.NewRecorder(auditRepo)

	settingsHTTP := httphandler.NewSettingsHandlers(settingsSvc, auditRec)
	customHTTP := httphandler.NewCustomFieldHandlers(customSvc, auditRec)
	auditHTTP := httphandler.NewAuditStatsHandlers(auditStatsSvc, auditRec)

	secret := "phase4-test-secret"
	token, _, err := jwtutil.GenerateAccess(secret, userA, "admin@test.com", false, &tenantA, nil, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(secret))
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RBACMiddleware(enforcer))
	api.GET("/settings/tenant", settingsHTTP.GetTenant)
	api.PATCH("/settings/tenant", settingsHTTP.PatchTenant)
	api.GET("/settings/features", settingsHTTP.ListFeatures)
	api.GET("/settings/custom-fields", customHTTP.List)
	api.POST("/settings/custom-fields", customHTTP.Create)
	api.PATCH("/settings/custom-fields/:id", customHTTP.Patch)
	api.DELETE("/settings/custom-fields/:id", customHTTP.Delete)
	api.GET("/audit/stats/by-action", auditHTTP.ByAction)
	api.GET("/audit/export", auditHTTP.Export)

	return &phase4HTTPEnv{
		t: t, router: r, auditRepo: auditRepo, auditStats: auditStatsRepo,
		tenantA: tenantA.String(), tenantB: tenantB.String(), userID: userA.String(), token: token,
	}
}

func newPhase4Enforcer(t *testing.T, userID, tenantA, tenantB uuid.UUID) *casbin.Enforcer {
	t.Helper()
	e := newLeadsTestEnforcer(t, userID, tenantA)
	role := userID.String()
	dom := tenantA.String()
	for _, res := range []string{"settings", "custom_fields"} {
		_, _ = e.AddPolicy(role, dom, res, "view")
		_, _ = e.AddPolicy(role, dom, res, "update")
	}
	_, _ = e.AddPolicy(role, dom, "audit", "view")
	_, _ = e.AddPolicy(role, dom, "audit", "export")
	domB := tenantB.String()
	_, _ = e.AddGroupingPolicy(userID.String(), role, domB)
	for _, res := range []string{"settings", "custom_fields", "audit"} {
		act := "view"
		if res != "audit" {
			_, _ = e.AddPolicy(role, domB, res, act)
			_, _ = e.AddPolicy(role, domB, res, "update")
			continue
		}
		_, _ = e.AddPolicy(role, domB, res, "view")
	}
	return e
}

type memSettingsRepo struct {
	mu      sync.Mutex
	tenants map[uuid.UUID]*domain.Tenant
}

func (m *memSettingsRepo) GetTenant(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tenants[tenantID]
	if !ok {
		return nil, repository.ErrTenantNotFound
	}
	cp := *t
	return &cp, nil
}

func (m *memSettingsRepo) UpdateTenant(ctx context.Context, tenantID uuid.UUID, name *string, config datatypes.JSON) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tenants[tenantID]
	if !ok {
		return repository.ErrTenantNotFound
	}
	if name != nil {
		t.Name = *name
	}
	if config != nil {
		t.Config = config
	}
	t.UpdatedAt = time.Now()
	return nil
}

type memCustomFieldRepo struct {
	mu    sync.Mutex
	items map[uuid.UUID]*domain.CustomField
}

func (m *memCustomFieldRepo) List(ctx context.Context, tenantID uuid.UUID, f repository.CustomFieldListFilter) ([]domain.CustomField, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []domain.CustomField
	for _, it := range m.items {
		if it.TenantID != tenantID {
			continue
		}
		if f.EntityType != "" && it.EntityType != f.EntityType {
			continue
		}
		if f.ActiveOnly && !it.IsActive {
			continue
		}
		out = append(out, *it)
	}
	return out, nil
}

func (m *memCustomFieldRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomField, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	it, ok := m.items[id]
	if !ok || it.TenantID != tenantID {
		return nil, repository.ErrCustomFieldNotFound
	}
	cp := *it
	return &cp, nil
}

func (m *memCustomFieldRepo) Create(ctx context.Context, f *domain.CustomField) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, it := range m.items {
		if it.TenantID == f.TenantID && it.EntityType == f.EntityType && it.FieldKey == f.FieldKey {
			return errDuplicateKey
		}
	}
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now
	cp := *f
	m.items[f.ID] = &cp
	return nil
}

func (m *memCustomFieldRepo) Update(ctx context.Context, f *domain.CustomField) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.items[f.ID]; !ok {
		return repository.ErrCustomFieldNotFound
	}
	f.UpdatedAt = time.Now()
	cp := *f
	m.items[f.ID] = &cp
	return nil
}

func (m *memCustomFieldRepo) Deactivate(ctx context.Context, tenantID, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	it, ok := m.items[id]
	if !ok || it.TenantID != tenantID {
		return repository.ErrCustomFieldNotFound
	}
	it.IsActive = false
	return nil
}

var errDuplicateKey = &duplicateKeyError{}

type duplicateKeyError struct{}

func (e *duplicateKeyError) Error() string { return "duplicate key" }

type memAuditStatsRepo struct {
	byAction []repository.ActionCountRow
	logs     []repository.AuditLogExportRow
}

func (m *memAuditStatsRepo) CountByAction(ctx context.Context, tenantID uuid.UUID, f repository.AuditStatsFilter) ([]repository.ActionCountRow, int64, error) {
	var total int64
	for _, r := range m.byAction {
		total += r.Count
	}
	return m.byAction, total, nil
}

func (m *memAuditStatsRepo) Trend(ctx context.Context, tenantID uuid.UUID, f repository.AuditStatsFilter, granularity string) ([]repository.TrendRow, error) {
	return []repository.TrendRow{{Bucket: "2026-05-26", Count: 2}}, nil
}

func (m *memAuditStatsRepo) TopActors(ctx context.Context, tenantID uuid.UUID, f repository.AuditStatsFilter) ([]repository.ActorCountRow, error) {
	return nil, nil
}

func (m *memAuditStatsRepo) ExportRows(ctx context.Context, tenantID uuid.UUID, f repository.AuditStatsFilter) ([]repository.AuditLogExportRow, error) {
	return m.logs, nil
}

type memInsightsTenantRepo struct{}

func (m *memInsightsTenantRepo) List(ctx context.Context, filter repository.TenantListFilter) ([]domain.Tenant, int64, error) {
	return nil, 0, nil
}
func (m *memInsightsTenantRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	return nil, repository.ErrTenantNotFound
}
func (m *memInsightsTenantRepo) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	return nil
}
func (m *memInsightsTenantRepo) CountUsers(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *memInsightsTenantRepo) CountAll(ctx context.Context) (int64, int64, error) {
	return 0, 0, nil
}
func (m *memInsightsTenantRepo) CountAllUsers(ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *memInsightsTenantRepo) TenantActivityTrend(ctx context.Context, days int) ([]repository.TenantActivityPoint, error) {
	return nil, nil
}

type memInsightsRepo struct{}

func (m *memInsightsRepo) TenantHealth(ctx context.Context) ([]repository.TenantHealthRow, error) {
	return []repository.TenantHealthRow{{
		TenantID: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		TenantName: "Demo", Scores: map[string]int{"activity": 80}, Overall: 80,
	}}, nil
}

func (m *memInsightsRepo) PlanDistribution(ctx context.Context, from, to time.Time) ([]repository.PlanCountRow, error) {
	return []repository.PlanCountRow{{Plan: "professional", Count: 1}}, nil
}

func (m *memInsightsRepo) TopTenants(ctx context.Context, metric string, limit int) ([]repository.TopTenantRow, error) {
	return []repository.TopTenantRow{{TenantID: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), TenantName: "Demo", Value: 10}}, nil
}

func (e *phase4HTTPEnv) request(method, path, tenantID string, body any) *httptest.ResponseRecorder {
	var rdr *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	} else {
		rdr = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, rdr)
	req.Header.Set("Authorization", "Bearer "+e.token)
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.router.ServeHTTP(w, req)
	return w
}

func (e *phase4HTTPEnv) get(path string) *httptest.ResponseRecorder {
	return e.request(http.MethodGet, path, e.tenantA, nil)
}

func (e *phase4HTTPEnv) post(path string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPost, path, e.tenantA, body)
}

func (e *phase4HTTPEnv) patch(path string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPatch, path, e.tenantA, body)
}

func (e *phase4HTTPEnv) parseBody(w *httptest.ResponseRecorder) map[string]any {
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	return body
}

func (e *phase4HTTPEnv) parseData(w *httptest.ResponseRecorder) map[string]any {
	body := e.parseBody(w)
	data, _ := body["data"].(map[string]any)
	return data
}

func (e *phase4HTTPEnv) parseDataID(w *httptest.ResponseRecorder) string {
	return e.parseData(w)["id"].(string)
}

func (e *phase4HTTPEnv) assertOK(w *httptest.ResponseRecorder) {
	if w.Code != http.StatusOK {
		e.t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func (e *phase4HTTPEnv) assertAuditAction(action string, min int) {
	tenantAID := uuid.MustParse(e.tenantA)
	e.auditRepo.mu.Lock()
	defer e.auditRepo.mu.Unlock()
	var n int
	for _, log := range e.auditRepo.logs {
		if log.TenantID == tenantAID && log.Action == action {
			n++
		}
	}
	if n < min {
		e.t.Fatalf("audit %s: got %d", action, n)
	}
}
