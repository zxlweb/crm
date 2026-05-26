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

	"crm-backend/internal/application/audit"
	dealapp "crm-backend/internal/application/deal"
	"crm-backend/internal/domain"
	httphandler "crm-backend/internal/interfaces/http"
	"crm-backend/internal/interfaces/middleware"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/jwtutil"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type memDealRepo struct {
	mu    sync.Mutex
	items map[uuid.UUID]*domain.Deal
}

func (m *memDealRepo) List(ctx context.Context, tenantID uuid.UUID, f repository.DealListFilter) ([]domain.Deal, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var matched []domain.Deal
	for _, d := range m.items {
		if d.TenantID != tenantID {
			continue
		}
		if !f.ViewAll && d.OwnerID != nil && *d.OwnerID != f.UserID {
			continue
		}
		if f.Stage != "" && d.Stage != f.Stage {
			continue
		}
		if len(f.Stages) > 0 {
			ok := false
			for _, s := range f.Stages {
				if d.Stage == s {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		if f.Search != "" && d.Title != f.Search {
			continue
		}
		matched = append(matched, *d)
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
		return []domain.Deal{}, total, nil
	}
	end := start + pageSize
	if end > len(matched) {
		end = len(matched)
	}
	return matched[start:end], total, nil
}

func (m *memDealRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Deal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.items[id]
	if !ok || d.TenantID != tenantID {
		return nil, repository.ErrDealNotFound
	}
	if !viewAll && d.OwnerID != nil && *d.OwnerID != userID {
		return nil, repository.ErrDealNotFound
	}
	cp := *d
	return &cp, nil
}

func (m *memDealRepo) Create(ctx context.Context, d *domain.Deal) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.items == nil {
		m.items = map[uuid.UUID]*domain.Deal{}
	}
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	if d.Stage == "" {
		d.Stage = crm.DealStageQualification
	}
	cp := *d
	m.items[d.ID] = &cp
	return nil
}

func (m *memDealRepo) Update(ctx context.Context, d *domain.Deal) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	d.UpdatedAt = time.Now()
	cp := *d
	m.items[d.ID] = &cp
	return nil
}

func (m *memDealRepo) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.items[id]
	if !ok || d.TenantID != tenantID {
		return repository.ErrDealNotFound
	}
	delete(m.items, id)
	return nil
}

func (m *memDealRepo) Pipeline(ctx context.Context, tenantID uuid.UUID, f repository.DealPipelineFilter) (map[string][]domain.Deal, repository.DealPipelineSummary, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	byStage := make(map[string][]domain.Deal)
	var summary repository.DealPipelineSummary
	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	for _, d := range m.items {
		if d.TenantID != tenantID {
			continue
		}
		if !f.ViewAll && d.OwnerID != nil && *d.OwnerID != f.UserID {
			continue
		}
		if f.OwnerID != nil && (d.OwnerID == nil || *d.OwnerID != *f.OwnerID) {
			continue
		}
		if f.AccountID != nil && (d.AccountID == nil || *d.AccountID != *f.AccountID) {
			continue
		}
		cp := *d
		list := byStage[d.Stage]
		if len(list) < f.PerStage {
			byStage[d.Stage] = append(list, cp)
		}
		if crm.IsDealOpen(d.Stage) {
			summary.OpenCount++
			summary.OpenAmount += d.Amount
		}
		if d.Stage == crm.DealStageWon && d.ClosedAt != nil && !d.ClosedAt.Before(monthStart) {
			summary.WonCountMTD++
			summary.WonAmountMTD += d.Amount
		}
	}
	return byStage, summary, nil
}

func (m *memDealRepo) CountScoped(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	_, total, err := m.List(ctx, tenantID, repository.DealListFilter{Page: 1, PageSize: 1000, ViewAll: viewAll, UserID: userID})
	return total, err
}

func (m *memDealRepo) StatsByStage(ctx context.Context, tenantID uuid.UUID, f repository.DealStatsFilter, metric string) ([]repository.DealStageStat, int64, error) {
	items, _, err := m.List(ctx, tenantID, repository.DealListFilter{ViewAll: f.ViewAll, UserID: f.UserID, Page: 1, PageSize: 1000})
	if err != nil {
		return nil, 0, err
	}
	counts := map[string]repository.DealStageStat{}
	for _, d := range items {
		row := counts[d.Stage]
		row.Label = d.Stage
		row.Value++
		row.Amount += d.Amount
		counts[d.Stage] = row
	}
	var total int64
	out := make([]repository.DealStageStat, 0, len(counts))
	for _, row := range counts {
		out = append(out, row)
		total += row.Value
	}
	return out, total, nil
}

func (m *memDealRepo) StatsWinRate(ctx context.Context, tenantID uuid.UUID, f repository.DealStatsFilter, granularity string) ([]repository.DealWinRatePoint, error) {
	_ = granularity
	return []repository.DealWinRatePoint{}, nil
}

func (m *memDealRepo) DailyCreatedCounts(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID, days int) ([]int64, error) {
	if days < 1 {
		days = 7
	}
	return make([]int64, days), nil
}

func (m *memDealRepo) CountByStage(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) ([]repository.LabelCount, error) {
	rows, _, err := m.StatsByStage(ctx, tenantID, repository.DealStatsFilter{ViewAll: viewAll, UserID: userID}, "count")
	if err != nil {
		return nil, err
	}
	out := make([]repository.LabelCount, len(rows))
	for i, row := range rows {
		out[i] = repository.LabelCount{Label: row.Label, Count: row.Value}
	}
	return out, nil
}

func (m *memDealRepo) TeamRanking(ctx context.Context, tenantID uuid.UUID, metric string, limit int) ([]repository.DealOwnerMetric, error) {
	_ = metric
	_ = limit
	return []repository.DealOwnerMetric{}, nil
}

func TestDealsHTTP_CRUDAndStageMachine(t *testing.T) {
	env := setupDealsHTTPEnv(t)

	w := env.post("/api/deals", map[string]any{"title": "云帆年度订阅", "amount": 280000})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create status %d body %s", w.Code, w.Body.String())
	}
	dealID := env.parseDataID(w)

	w = env.putStage(dealID, map[string]any{"stage": "proposal"})
	env.assertOK(w)
	env.assertStage(w, "proposal")

	w = env.putStage(dealID, map[string]any{"stage": "qualification"})
	env.assertOK(w)

	w = env.putStage(dealID, map[string]any{"stage": "proposal"})
	env.assertOK(w)

	w = env.putStage(dealID, map[string]any{"stage": "negotiation"})
	env.assertOK(w)

	w = env.putStage(dealID, map[string]any{"stage": "won"})
	env.assertOK(w)
	env.assertStage(w, "won")

	w = env.putStage(dealID, map[string]any{"stage": "proposal"})
	env.assertBadMessage(w, "deal_closed_readonly")

	w = env.post("/api/deals", map[string]any{"title": "非法跳转", "stage": "qualification"})
	id2 := env.parseDataID(w)
	w = env.patch(id2, map[string]any{"stage": "won"})
	env.assertBadMessage(w, "invalid_stage_transition")
}

func TestDealsHTTP_Pipeline(t *testing.T) {
	env := setupDealsHTTPEnv(t)

	env.post("/api/deals", map[string]any{"title": "A", "stage": "qualification", "amount": 100})
	env.post("/api/deals", map[string]any{"title": "B", "stage": "proposal", "amount": 200})

	w := env.get("/api/deals/pipeline")
	env.assertOK(w)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	stages, ok := data["stages"].([]any)
	if !ok || len(stages) < 2 {
		t.Fatalf("pipeline stages: %v", data["stages"])
	}
	summary := data["summary"].(map[string]any)
	if summary["open_count"].(float64) < 2 {
		t.Fatalf("open_count: %v", summary["open_count"])
	}
}

func TestDealsHTTP_TenantIsolation(t *testing.T) {
	env := setupDealsHTTPEnv(t)

	w := env.post("/api/deals", map[string]any{"title": "Tenant A Deal"})
	dealID := env.parseDataID(w)

	w = env.getWithTenant("/api/deals/"+dealID, env.tenantB)
	if w.Code != http.StatusNotFound && w.Code != http.StatusForbidden {
		t.Fatalf("cross-tenant get: status %d body %s", w.Code, w.Body.String())
	}

	w = env.getWithTenant("/api/deals", env.tenantB)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	items := data["items"].([]any)
	if len(items) != 0 {
		t.Fatalf("tenant B list should be empty, got %d", len(items))
	}
}

type dealsHTTPEnv struct {
	t        *testing.T
	router   *gin.Engine
	dealRepo *memDealRepo
	tenantA  string
	tenantB  string
	token    string
}

func setupDealsHTTPEnv(t *testing.T) *dealsHTTPEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	tenantB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	enforcer := newDealsTestEnforcer(t, userA, tenantA)
	dealRepo := &memDealRepo{items: map[uuid.UUID]*domain.Deal{}}
	dealSvc := dealapp.NewService(dealRepo, nil, enforcer)
	auditRepo := &memAuditRepo{}
	auditRec := audit.NewRecorder(auditRepo)
	dealHTTP := httphandler.NewDealHandlers(dealSvc, auditRec)

	secret := "deals-test-secret"
	token, _, err := jwtutil.GenerateAccess(secret, userA, "sales@test.com", false, &tenantA, nil, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(secret))
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RBACMiddleware(enforcer))
	registerDealsRoutes(api, dealHTTP)

	return &dealsHTTPEnv{
		t:        t,
		router:   r,
		dealRepo: dealRepo,
		tenantA:  tenantA.String(),
		tenantB:  tenantB.String(),
		token:    token,
	}
}

func registerDealsRoutes(api *gin.RouterGroup, dealHTTP *httphandler.DealHandlers) {
	api.GET("/deals/pipeline", dealHTTP.Pipeline)
	api.GET("/deals", dealHTTP.List)
	api.POST("/deals", dealHTTP.Create)
	api.GET("/deals/:id", dealHTTP.Get)
	api.PUT("/deals/:id", dealHTTP.Put)
	api.PATCH("/deals/:id", dealHTTP.Patch)
	api.DELETE("/deals/:id", dealHTTP.Delete)
	api.PUT("/deals/:id/stage", dealHTTP.PutStage)
}

func newDealsTestEnforcer(t *testing.T, userID, tenantID uuid.UUID) *casbin.Enforcer {
	t.Helper()
	const mdef = `
[request_definition]
r = sub, dom, obj, act
[policy_definition]
p = sub, dom, obj, act
[role_definition]
g = _, _, dom
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
`
	m, err := model.NewModelFromString(mdef)
	if err != nil {
		t.Fatal(err)
	}
	e, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatal(err)
	}
	role := userID.String()
	dom := tenantID.String()
	for _, act := range []string{"view", "create", "update", "delete"} {
		_, _ = e.AddPolicy(role, dom, "deals", act)
	}
	_, _ = e.AddGroupingPolicy(userID.String(), role, dom)
	domB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb").String()
	for _, act := range []string{"view", "create", "update", "delete"} {
		_, _ = e.AddPolicy(role, domB, "deals", act)
	}
	_, _ = e.AddGroupingPolicy(userID.String(), role, domB)
	return e
}

func (e *dealsHTTPEnv) request(method, path, tenantID string, body any) *httptest.ResponseRecorder {
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

func (e *dealsHTTPEnv) post(path string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPost, path, e.tenantA, body)
}

func (e *dealsHTTPEnv) get(path string) *httptest.ResponseRecorder {
	return e.request(http.MethodGet, path, e.tenantA, nil)
}

func (e *dealsHTTPEnv) patch(id string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPatch, "/api/deals/"+id, e.tenantA, body)
}

func (e *dealsHTTPEnv) putStage(id string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPut, "/api/deals/"+id+"/stage", e.tenantA, body)
}

func (e *dealsHTTPEnv) getWithTenant(path, tenantID string) *httptest.ResponseRecorder {
	return e.request(http.MethodGet, path, tenantID, nil)
}

func (e *dealsHTTPEnv) parseDataID(w *httptest.ResponseRecorder) string {
	e.t.Helper()
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		e.t.Fatalf("parse: %v body %s", err, w.Body.String())
	}
	data := body["data"].(map[string]any)
	id, _ := data["id"].(string)
	if id == "" {
		e.t.Fatalf("no id: %v", data)
	}
	return id
}

func (e *dealsHTTPEnv) assertOK(w *httptest.ResponseRecorder) {
	e.t.Helper()
	if w.Code != http.StatusOK {
		e.t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func (e *dealsHTTPEnv) assertStage(w *httptest.ResponseRecorder, want string) {
	e.t.Helper()
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	if data["stage"] != want {
		e.t.Fatalf("stage %v want %s", data["stage"], want)
	}
}

func (e *dealsHTTPEnv) assertBadMessage(w *httptest.ResponseRecorder, msg string) {
	e.t.Helper()
	if w.Code != http.StatusBadRequest {
		e.t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != msg {
		e.t.Fatalf("message %v want %s", body["message"], msg)
	}
}
