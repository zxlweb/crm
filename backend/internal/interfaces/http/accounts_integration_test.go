package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	accountapp "crm-backend/internal/application/account"
	emotionapp "crm-backend/internal/application/emotion"
	"crm-backend/internal/application/audit"
	"crm-backend/internal/domain"
	httphandler "crm-backend/internal/interfaces/http"
	"crm-backend/internal/interfaces/middleware"
	"crm-backend/internal/pkg/jwtutil"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestAccountsHTTP_CRUDLifecycleAndAudit(t *testing.T) {
	env := setupAccountsHTTPEnv(t)

	w := env.post("/api/accounts", map[string]any{
		"name":     "华创科技有限公司",
		"industry": "software",
		"tags":     []string{"vip"},
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create: %d %s", w.Code, w.Body.String())
	}
	id := env.parseDataID(w)
	data := env.parseData(w)
	if data["lifecycle_stage"] != "acquire" {
		t.Fatalf("lifecycle_stage: %v", data["lifecycle_stage"])
	}
	if data["relationship_health"] != "low" {
		t.Fatalf("relationship_health: %v", data["relationship_health"])
	}
	owner, _ := data["owner_id"].(string)
	if owner != env.userID {
		t.Fatalf("owner_id default: %v", data["owner_id"])
	}

	w = env.get("/api/accounts/" + id)
	env.assertOK(w)

	w = env.patch(id, map[string]any{"name": "华创科技（更新）"})
	env.assertOK(w)

	w = env.patch(id, map[string]any{"lifecycle_stage": "grow"})
	env.assertOK(w)

	env.assertAuditAction("lifecycle.change", 1)

	w = env.get("/api/accounts/" + id)
	data = env.parseData(w)
	if data["lifecycle_stage"] != "grow" {
		t.Fatalf("lifecycle after patch: %v", data["lifecycle_stage"])
	}

	w = env.delete(id)
	env.assertOK(w)
	env.assertAuditAction("account.delete", 1)

	w = env.get("/api/accounts/" + id)
	if w.Code != http.StatusNotFound {
		t.Fatalf("after delete: %d %s", w.Code, w.Body.String())
	}
}

func TestAccountsHTTP_ListSearchAndFilter(t *testing.T) {
	env := setupAccountsHTTPEnv(t)

	env.post("/api/accounts", map[string]any{"name": "Alpha Corp", "lifecycle_stage": "acquire"})
	env.post("/api/accounts", map[string]any{"name": "Beta 华创", "lifecycle_stage": "grow"})

	w := env.get("/api/accounts?search=华创")
	env.assertOK(w)
	body := env.parseBody(w)
	data := body["data"].(map[string]any)
	items := data["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("search items: %d", len(items))
	}

	w = env.get("/api/accounts?lifecycle_stage=grow")
	env.assertOK(w)
	body = env.parseBody(w)
	data = body["data"].(map[string]any)
	items = data["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("lifecycle filter: %d", len(items))
	}

	w = env.get("/api/accounts?page=1&page_size=1")
	env.assertOK(w)
	body = env.parseBody(w)
	pag, ok := body["pagination"].(map[string]any)
	if !ok || pag["total"] == nil {
		t.Fatalf("missing pagination: %v", body["pagination"])
	}
}

func TestAccountsHTTP_ReadOnlyEngagementScore(t *testing.T) {
	env := setupAccountsHTTPEnv(t)

	w := env.post("/api/accounts", map[string]any{"name": "Score Co"})
	id := env.parseDataID(w)

	uid, _ := uuid.Parse(id)
	env.accountRepo.mu.Lock()
	if a, ok := env.accountRepo.items[uid]; ok {
		a.EngagementScore = 85
	}
	env.accountRepo.mu.Unlock()

	w = env.patch(id, map[string]any{"name": "Score Co Renamed"})
	env.assertOK(w)
	data := env.parseData(w)
	if int(data["engagement_score"].(float64)) != 85 {
		t.Fatalf("engagement_score changed: %v", data["engagement_score"])
	}
	if data["relationship_health"] != "high" {
		t.Fatalf("relationship_health: %v", data["relationship_health"])
	}
}

func TestAccountsHTTP_InvalidLifecycle(t *testing.T) {
	env := setupAccountsHTTPEnv(t)
	w := env.post("/api/accounts", map[string]any{"name": "X", "lifecycle_stage": "invalid"})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d", w.Code)
	}
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "invalid lifecycle_stage" {
		t.Fatalf("message: %v", body["message"])
	}
}

func TestAccountsHTTP_TenantIsolation(t *testing.T) {
	env := setupAccountsHTTPEnv(t)

	w := env.post("/api/accounts", map[string]any{"name": "Tenant A Only"})
	id := env.parseDataID(w)

	w = env.getWithTenant("/api/accounts/"+id, env.tenantB)
	if w.Code != http.StatusNotFound && w.Code != http.StatusForbidden {
		t.Fatalf("cross-tenant get: %d %s", w.Code, w.Body.String())
	}

	w = env.getWithTenant("/api/accounts", env.tenantB)
	body := env.parseBody(w)
	data := body["data"].(map[string]any)
	items, _ := data["items"].([]any)
	if len(items) != 0 {
		t.Fatalf("tenant B list len %d", len(items))
	}
}

func TestAccountsHTTP_EmotionJourneyAndInsights(t *testing.T) {
	env := setupAccountsHTTPEnv(t)
	w := env.post("/api/accounts", map[string]any{"name": "Journey Co"})
	id := env.parseDataID(w)

	w = env.get("/api/accounts/" + id + "/emotion-journey")
	env.assertOK(w)
	data := env.parseData(w)
	if data["subject_type"] != "account" {
		t.Fatalf("subject_type: %v", data["subject_type"])
	}

	w = env.post("/api/accounts/"+id+"/insights/evaluate", map[string]any{})
	env.assertOK(w)
	data = env.parseData(w)
	items, ok := data["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatalf("insights missing items: %v", data["items"])
	}
	first, _ := items[0].(map[string]any)
	if first["rule_id"] != "INS-001" {
		t.Fatalf("expected INS-001, got %v", first["rule_id"])
	}
}

type accountsHTTPEnv struct {
	t           *testing.T
	router      *gin.Engine
	auditRepo   *memAuditRepo
	accountRepo *memAccountRepo
	tenantA     string
	tenantB     string
	userID      string
	token       string
}

func setupAccountsHTTPEnv(t *testing.T) *accountsHTTPEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	tenantB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	enforcer := newResourceTestEnforcer(t, userA, tenantA, tenantB, "accounts")
	accountRepo := &memAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	auditRepo := &memAuditRepo{}
	auditRec := audit.NewRecorder(auditRepo)
	accountSvc := accountapp.NewService(accountRepo, nil, enforcer)
	activityRepo := &memActivityRepo{items: map[uuid.UUID]*domain.Activity{}}
	emotionSvc := emotionapp.NewService(activityRepo)
	accountHTTP := httphandler.NewAccountHandlers(accountSvc, auditRec, emotionSvc)

	secret := "accounts-test-secret"
	token, _, err := jwtutil.GenerateAccess(secret, userA, "sales@test.com", false, &tenantA, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(secret))
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RBACMiddleware(enforcer))
	api.GET("/accounts", accountHTTP.List)
	api.POST("/accounts", accountHTTP.Create)
	api.GET("/accounts/:id", accountHTTP.Get)
	api.PATCH("/accounts/:id", accountHTTP.Patch)
	api.DELETE("/accounts/:id", accountHTTP.Delete)
	api.GET("/accounts/:id/emotion-journey", accountHTTP.EmotionJourney)
	api.POST("/accounts/:id/insights/evaluate", accountHTTP.EvaluateInsights)

	return &accountsHTTPEnv{
		t:           t,
		router:      r,
		auditRepo:   auditRepo,
		accountRepo: accountRepo,
		tenantA:     tenantA.String(),
		tenantB:     tenantB.String(),
		userID:      userA.String(),
		token:       token,
	}
}

func newResourceTestEnforcer(t *testing.T, userID, tenantA, tenantB uuid.UUID, resource string) *casbin.Enforcer {
	t.Helper()
	e := newLeadsTestEnforcer(t, userID, tenantA)
	role := userID.String()
	domB := tenantB.String()
	for _, act := range []string{"view", "create", "update", "delete"} {
		_, _ = e.AddPolicy(role, domB, resource, act)
	}
	_, _ = e.AddGroupingPolicy(userID.String(), role, domB)
	domA := tenantA.String()
	for _, act := range []string{"view", "create", "update", "delete"} {
		_, _ = e.AddPolicy(role, domA, resource, act)
	}
	return e
}

func (e *accountsHTTPEnv) request(method, path, tenantID string, body any) *httptest.ResponseRecorder {
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

func (e *accountsHTTPEnv) post(path string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPost, path, e.tenantA, body)
}

func (e *accountsHTTPEnv) get(path string) *httptest.ResponseRecorder {
	return e.request(http.MethodGet, path, e.tenantA, nil)
}

func (e *accountsHTTPEnv) getWithTenant(path, tenantID string) *httptest.ResponseRecorder {
	return e.request(http.MethodGet, path, tenantID, nil)
}

func (e *accountsHTTPEnv) patch(id string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPatch, "/api/accounts/"+id, e.tenantA, body)
}

func (e *accountsHTTPEnv) delete(id string) *httptest.ResponseRecorder {
	return e.request(http.MethodDelete, "/api/accounts/"+id, e.tenantA, nil)
}

func (e *accountsHTTPEnv) parseBody(w *httptest.ResponseRecorder) map[string]any {
	e.t.Helper()
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		e.t.Fatalf("parse body: %v", err)
	}
	return body
}

func (e *accountsHTTPEnv) parseData(w *httptest.ResponseRecorder) map[string]any {
	e.t.Helper()
	body := e.parseBody(w)
	data, ok := body["data"].(map[string]any)
	if !ok {
		e.t.Fatalf("no data: %s", w.Body.String())
	}
	return data
}

func (e *accountsHTTPEnv) parseDataID(w *httptest.ResponseRecorder) string {
	data := e.parseData(w)
	id, _ := data["id"].(string)
	if id == "" {
		e.t.Fatalf("no id: %v", data)
	}
	return id
}

func (e *accountsHTTPEnv) assertOK(w *httptest.ResponseRecorder) {
	e.t.Helper()
	if w.Code != http.StatusOK {
		e.t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func (e *accountsHTTPEnv) assertAuditAction(action string, min int) {
	e.t.Helper()
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
		e.t.Fatalf("audit %s: got %d want >= %d", action, n, min)
	}
}
