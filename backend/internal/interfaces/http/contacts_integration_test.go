package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	accountapp "crm-backend/internal/application/account"
	contactapp "crm-backend/internal/application/contact"
	emotionapp "crm-backend/internal/application/emotion"
	"crm-backend/internal/application/audit"
	"crm-backend/internal/domain"
	httphandler "crm-backend/internal/interfaces/http"
	"crm-backend/internal/interfaces/middleware"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/jwtutil"
	"crm-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestContactsHTTP_CRUDAndAccountLink(t *testing.T) {
	env := setupContactsHTTPEnv(t)

	accW := env.post("/api/accounts", map[string]any{"name": "联系人测试公司"})
	accountID := env.parseDataID(accW)

	w := env.post("/api/contacts", map[string]any{
		"account_id":  accountID,
		"first_name":  "张",
		"last_name":   "三",
		"email":       "zhang@example.com",
		"is_primary":  true,
		"tags":        []string{"decision-maker"},
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create contact: %d %s", w.Code, w.Body.String())
	}
	id := env.parseDataID(w)
	data := env.parseData(w)
	if data["display_name"] != "张 三" {
		t.Fatalf("display_name: %v", data["display_name"])
	}
	if data["relationship_health"] != "low" {
		t.Fatalf("relationship_health: %v", data["relationship_health"])
	}

	w = env.get("/api/contacts/" + id)
	env.assertOK(w)

	w = env.get("/api/accounts/" + accountID + "/contacts")
	env.assertOK(w)
	body := env.parseBody(w)
	items := body["data"].(map[string]any)["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("account contacts: %d", len(items))
	}

	w = env.patch(id, map[string]any{"last_name": "叁", "lifecycle_stage": "grow"})
	env.assertOK(w)

	w = env.delete(id)
	env.assertOK(w)

	w = env.get("/api/contacts/" + id)
	if w.Code != http.StatusNotFound {
		t.Fatalf("after delete: %d", w.Code)
	}
}

func TestContactsHTTP_InsightsEvaluate_INS001(t *testing.T) {
	env := setupContactsHTTPEnv(t)

	w := env.post("/api/contacts", map[string]any{
		"first_name": "洞察",
		"last_name":  "测试",
		"email":      "insight-eval@example.com",
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create contact: %d %s", w.Code, w.Body.String())
	}
	id := env.parseDataID(w)

	w = env.post("/api/contacts/"+id+"/insights/evaluate", map[string]any{})
	env.assertOK(w)
	data := env.parseData(w)
	items, ok := data["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatalf("expected INS-001 hit, items: %v", data["items"])
	}
	first, _ := items[0].(map[string]any)
	if first["rule_id"] != "INS-001" {
		t.Fatalf("rule_id: %v", first["rule_id"])
	}
}

func TestContactsHTTP_ListFilterByAccount(t *testing.T) {
	env := setupContactsHTTPEnv(t)
	accW := env.post("/api/accounts", map[string]any{"name": "A Co"})
	accountID := env.parseDataID(accW)

	env.post("/api/contacts", map[string]any{"account_id": accountID, "first_name": "A", "email": "a@ex.com"})
	env.post("/api/contacts", map[string]any{"first_name": "B", "email": "b@ex.com"})

	w := env.get("/api/contacts?account_id=" + accountID)
	env.assertOK(w)
	body := env.parseBody(w)
	items := body["data"].(map[string]any)["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("filter account_id: %d", len(items))
	}
}

type contactsHTTPEnv struct {
	t           *testing.T
	router      *gin.Engine
	contactRepo *memContactRepo
	accountRepo *memAccountRepo
	tenantA     string
	tenantB     string
	userID      string
	token       string
}

func setupContactsHTTPEnv(t *testing.T) *contactsHTTPEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	tenantB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	enforcer := newResourceTestEnforcer(t, userA, tenantA, tenantB, "contacts")
	role := userA.String()
	for _, act := range []string{"view", "create"} {
		_, _ = enforcer.AddPolicy(role, tenantA.String(), "accounts", act)
	}

	accountRepo := &memAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	contactRepo := &memContactRepo{items: map[uuid.UUID]*domain.Contact{}}
	auditRec := audit.NewRecorder(&memAuditRepo{})

	contactSvc := contactapp.NewService(contactRepo, accountRepo, enforcer)
	activityRepo := &memActivityRepo{items: map[uuid.UUID]*domain.Activity{}}
	emotionSvc := emotionapp.NewService(activityRepo)
	contactHTTP := httphandler.NewContactHandlers(contactSvc, auditRec, emotionSvc)

	accountHTTP := httphandler.NewAccountHandlers(
		accountapp.NewService(accountRepo, nil, enforcer),
		auditRec,
		emotionSvc,
	)

	secret := "contacts-test-secret"
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
	api.GET("/accounts/:id/contacts", contactHTTP.ListByAccount)
	api.GET("/contacts", contactHTTP.List)
	api.POST("/contacts", contactHTTP.Create)
	api.GET("/contacts/:id", contactHTTP.Get)
	api.PATCH("/contacts/:id", contactHTTP.Patch)
	api.DELETE("/contacts/:id", contactHTTP.Delete)
	api.POST("/contacts/:id/insights/evaluate", contactHTTP.EvaluateInsights)

	return &contactsHTTPEnv{
		t:           t,
		router:      r,
		contactRepo: contactRepo,
		accountRepo: accountRepo,
		tenantA:     tenantA.String(),
		tenantB:     tenantB.String(),
		userID:      userA.String(),
		token:       token,
	}
}

type memContactRepo struct {
	mu    sync.Mutex
	items map[uuid.UUID]*domain.Contact
}

func (m *memContactRepo) List(ctx context.Context, tenantID uuid.UUID, f repository.ContactListFilter) ([]domain.Contact, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var matched []domain.Contact
	for _, c := range m.items {
		if c.TenantID != tenantID {
			continue
		}
		if !f.ViewAll && c.OwnerID != nil && *c.OwnerID != f.UserID {
			continue
		}
		if f.Search != "" {
			needle := strings.ToLower(f.Search)
			hay := strings.ToLower(c.FirstName + " " + c.LastName + " " + c.Email + " " + c.Phone)
			if !strings.Contains(hay, needle) {
				continue
			}
		}
		if f.LifecycleStage != "" && c.LifecycleStage != f.LifecycleStage {
			continue
		}
		if f.RelationshipHealth != "" && crm.RelationshipHealthFromScore(c.EngagementScore) != f.RelationshipHealth {
			continue
		}
		if f.AccountID != nil && (c.AccountID == nil || *c.AccountID != *f.AccountID) {
			continue
		}
		if f.OwnerID != nil && (c.OwnerID == nil || *c.OwnerID != *f.OwnerID) {
			continue
		}
		matched = append(matched, *c)
	}
	total := int64(len(matched))
	page := f.Page
	if page < 1 {
		page = 1
	}
	pageSize := f.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	start := (page - 1) * pageSize
	if start >= len(matched) {
		return []domain.Contact{}, total, nil
	}
	end := start + pageSize
	if end > len(matched) {
		end = len(matched)
	}
	return matched[start:end], total, nil
}

func (m *memContactRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Contact, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.items[id]
	if !ok || c.TenantID != tenantID {
		return nil, repository.ErrContactNotFound
	}
	cp := *c
	return &cp, nil
}

func (m *memContactRepo) Create(ctx context.Context, c *domain.Contact) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.items == nil {
		m.items = map[uuid.UUID]*domain.Contact{}
	}
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	cp := *c
	m.items[c.ID] = &cp
	return nil
}

func (m *memContactRepo) Update(ctx context.Context, c *domain.Contact) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c.UpdatedAt = time.Now()
	cp := *c
	m.items[c.ID] = &cp
	return nil
}

func (m *memContactRepo) UpdateEngagementFromActivity(ctx context.Context, tenantID, id, updatedBy uuid.UUID, last *time.Time, score int16) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.items[id]
	if !ok || c.TenantID != tenantID {
		return repository.ErrContactNotFound
	}
	c.LastActivityAt = last
	c.EngagementScore = score
	c.UpdatedBy = updatedBy
	cp := *c
	m.items[id] = &cp
	return nil
}

func (m *memContactRepo) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.items[id]
	if !ok || c.TenantID != tenantID {
		return repository.ErrContactNotFound
	}
	delete(m.items, id)
	return nil
}

func (e *contactsHTTPEnv) request(method, path, tenantID string, body any) *httptest.ResponseRecorder {
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

func (e *contactsHTTPEnv) post(path string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPost, path, e.tenantA, body)
}

func (e *contactsHTTPEnv) get(path string) *httptest.ResponseRecorder {
	return e.request(http.MethodGet, path, e.tenantA, nil)
}

func (e *contactsHTTPEnv) patch(id string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPatch, "/api/contacts/"+id, e.tenantA, body)
}

func (e *contactsHTTPEnv) delete(id string) *httptest.ResponseRecorder {
	return e.request(http.MethodDelete, "/api/contacts/"+id, e.tenantA, nil)
}

func (e *contactsHTTPEnv) parseBody(w *httptest.ResponseRecorder) map[string]any {
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	return body
}

func (e *contactsHTTPEnv) parseData(w *httptest.ResponseRecorder) map[string]any {
	body := e.parseBody(w)
	data, _ := body["data"].(map[string]any)
	return data
}

func (e *contactsHTTPEnv) parseDataID(w *httptest.ResponseRecorder) string {
	data := e.parseData(w)
	id, _ := data["id"].(string)
	return id
}

func (e *contactsHTTPEnv) assertOK(w *httptest.ResponseRecorder) {
	e.t.Helper()
	if w.Code != http.StatusOK {
		e.t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
