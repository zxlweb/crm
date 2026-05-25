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

	"crm-backend/internal/application/audit"
	emotionapp "crm-backend/internal/application/emotion"
	leadapp "crm-backend/internal/application/lead"
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

type memLeadRepo struct {
	mu    sync.Mutex
	items map[uuid.UUID]*domain.Lead
}

func (m *memLeadRepo) List(ctx context.Context, tenantID uuid.UUID, f repository.LeadListFilter) ([]domain.Lead, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []domain.Lead
	for _, l := range m.items {
		if l.TenantID != tenantID {
			continue
		}
		if !f.ViewAll && l.OwnerID != nil && *l.OwnerID != f.UserID {
			continue
		}
		if f.Segment != "" && !leadMatchesSegment(l, f.Segment, f.SegmentOpts) {
			continue
		}
		out = append(out, *l)
	}
	return out, int64(len(out)), nil
}

func (m *memLeadRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Lead, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.items[id]
	if !ok || l.TenantID != tenantID {
		return nil, repository.ErrLeadNotFound
	}
	if !viewAll && l.OwnerID != nil && *l.OwnerID != userID {
		return nil, repository.ErrLeadNotFound
	}
	cp := *l
	return &cp, nil
}

func (m *memLeadRepo) Create(ctx context.Context, l *domain.Lead) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.items == nil {
		m.items = map[uuid.UUID]*domain.Lead{}
	}
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	now := time.Now()
	l.CreatedAt = now
	l.UpdatedAt = now
	if l.LifecycleStage == "" {
		l.LifecycleStage = "acquire"
	}
	if l.Status == "" {
		l.Status = "new"
	}
	cp := *l
	m.items[l.ID] = &cp
	return nil
}

func (m *memLeadRepo) Update(ctx context.Context, l *domain.Lead) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l.UpdatedAt = time.Now()
	cp := *l
	m.items[l.ID] = &cp
	return nil
}

func (m *memLeadRepo) UpdateEngagementFromActivity(ctx context.Context, tenantID, id, updatedBy uuid.UUID, last *time.Time, score int16) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.items[id]
	if !ok || l.TenantID != tenantID {
		return repository.ErrLeadNotFound
	}
	l.LastActivityAt = last
	l.EngagementScore = score
	l.RelationshipHealth = crm.RelationshipHealthFromScore(score)
	l.UpdatedBy = updatedBy
	l.UpdatedAt = time.Now()
	cp := *l
	m.items[id] = &cp
	return nil
}

func (m *memLeadRepo) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.items[id]
	if !ok || l.TenantID != tenantID {
		return repository.ErrLeadNotFound
	}
	delete(m.items, id)
	return nil
}

func (m *memLeadRepo) filterStats(tenantID uuid.UUID, f repository.LeadStatsFilter) []*domain.Lead {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []*domain.Lead
	for _, l := range m.items {
		if l.TenantID != tenantID {
			continue
		}
		if !f.ViewAll && l.OwnerID != nil && *l.OwnerID != f.UserID {
			continue
		}
		if f.From != nil && l.CreatedAt.Before(*f.From) {
			continue
		}
		if f.To != nil && !l.CreatedAt.Before(*f.To) {
			continue
		}
		cp := *l
		out = append(out, &cp)
	}
	return out
}

func (m *memLeadRepo) StatsBySource(ctx context.Context, tenantID uuid.UUID, f repository.LeadStatsFilter) ([]repository.LabelCount, int64, error) {
	return m.countBy(ctx, tenantID, f, func(l *domain.Lead) string {
		if l.Source == "" {
			return "unknown"
		}
		return l.Source
	})
}

func (m *memLeadRepo) StatsByStatus(ctx context.Context, tenantID uuid.UUID, f repository.LeadStatsFilter) ([]repository.LabelCount, int64, error) {
	return m.countBy(ctx, tenantID, f, func(l *domain.Lead) string { return l.Status })
}

func (m *memLeadRepo) countBy(ctx context.Context, tenantID uuid.UUID, f repository.LeadStatsFilter, keyFn func(*domain.Lead) string) ([]repository.LabelCount, int64, error) {
	_ = ctx
	leads := m.filterStats(tenantID, f)
	counts := map[string]int64{}
	for _, l := range leads {
		counts[keyFn(l)]++
	}
	var total int64
	out := make([]repository.LabelCount, 0, len(counts))
	for label, count := range counts {
		out = append(out, repository.LabelCount{Label: label, Count: count})
		total += count
	}
	return out, total, nil
}

func (m *memLeadRepo) StatsTrend(ctx context.Context, tenantID uuid.UUID, f repository.LeadStatsFilter, granularity string) ([]repository.TrendPoint, error) {
	_ = granularity
	_ = ctx
	leads := m.filterStats(tenantID, f)
	byDay := map[string]int64{}
	for _, l := range leads {
		d := l.CreatedAt.UTC().Format("2006-01-02")
		byDay[d]++
	}
	out := make([]repository.TrendPoint, 0, len(byDay))
	for d, count := range byDay {
		out = append(out, repository.TrendPoint{Date: d, Count: count})
	}
	return out, nil
}

func (m *memLeadRepo) CountScoped(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	_, total, err := m.List(ctx, tenantID, repository.LeadListFilter{Page: 1, PageSize: 1000, ViewAll: viewAll, UserID: userID})
	return total, err
}

func (m *memLeadRepo) DailyCreatedCounts(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID, days int) ([]int64, error) {
	_ = ctx
	if days < 1 {
		days = 7
	}
	return make([]int64, days), nil
}

func (m *memLeadRepo) CountLowEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	_ = ctx
	return 0, nil
}

func (m *memLeadRepo) AvgEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (float64, error) {
	_ = ctx
	return 0, nil
}

func (m *memLeadRepo) ListPriorityCandidates(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID, limit int) ([]domain.Lead, error) {
	items, _, err := m.List(ctx, tenantID, repository.LeadListFilter{Page: 1, PageSize: limit, ViewAll: viewAll, UserID: userID})
	return items, err
}

func (m *memLeadRepo) StatsFunnel(ctx context.Context, tenantID uuid.UUID, f repository.LeadStatsFilter) ([]repository.LabelCount, error) {
	rows, _, err := m.StatsByStatus(ctx, tenantID, f)
	if err != nil {
		return nil, err
	}
	byStatus := map[string]int64{}
	for _, row := range rows {
		byStatus[row.Label] = row.Count
	}
	order := []string{"new", "contacted", "qualified", "unqualified", "converted"}
	out := make([]repository.LabelCount, 0, len(order))
	for _, status := range order {
		out = append(out, repository.LabelCount{Label: status, Count: byStatus[status]})
	}
	return out, nil
}

type memAccountRepo struct {
	mu    sync.Mutex
	items map[uuid.UUID]*domain.Account
}

func (m *memAccountRepo) List(ctx context.Context, tenantID uuid.UUID, f repository.AccountListFilter) ([]domain.Account, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var matched []domain.Account
	for _, a := range m.items {
		if a.TenantID != tenantID {
			continue
		}
		if !f.ViewAll && a.OwnerID != nil && *a.OwnerID != f.UserID {
			continue
		}
		if f.Search != "" && !strings.Contains(strings.ToLower(a.Name), strings.ToLower(f.Search)) {
			continue
		}
		if f.LifecycleStage != "" && a.LifecycleStage != f.LifecycleStage {
			continue
		}
		if f.RelationshipHealth != "" && crm.RelationshipHealthFromScore(a.EngagementScore) != f.RelationshipHealth {
			continue
		}
		if f.OwnerID != nil && (a.OwnerID == nil || *a.OwnerID != *f.OwnerID) {
			continue
		}
		if f.Segment != "" && !accountMatchesSegment(a, f.Segment, f.SegmentOpts) {
			continue
		}
		matched = append(matched, *a)
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
	if pageSize > 100 {
		pageSize = 100
	}
	start := (page - 1) * pageSize
	if start >= len(matched) {
		return []domain.Account{}, total, nil
	}
	end := start + pageSize
	if end > len(matched) {
		end = len(matched)
	}
	return matched[start:end], total, nil
}

func (m *memAccountRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.items[id]
	if !ok || a.TenantID != tenantID {
		return nil, repository.ErrAccountNotFound
	}
	cp := *a
	return &cp, nil
}

func (m *memAccountRepo) Create(ctx context.Context, a *domain.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.items == nil {
		m.items = map[uuid.UUID]*domain.Account{}
	}
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	cp := *a
	m.items[a.ID] = &cp
	return nil
}

func (m *memAccountRepo) Update(ctx context.Context, a *domain.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *a
	m.items[a.ID] = &cp
	return nil
}

func (m *memAccountRepo) UpdateEngagementFromActivity(ctx context.Context, tenantID, id, updatedBy uuid.UUID, last *time.Time, score int16) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.items[id]
	if !ok || a.TenantID != tenantID {
		return repository.ErrAccountNotFound
	}
	a.LastActivityAt = last
	a.EngagementScore = score
	a.UpdatedBy = updatedBy
	cp := *a
	m.items[id] = &cp
	return nil
}

func (m *memAccountRepo) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	delete(m.items, id)
	return nil
}

func (m *memAccountRepo) CountScoped(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	_, total, err := m.List(ctx, tenantID, repository.AccountListFilter{Page: 1, PageSize: 1000, ViewAll: viewAll, UserID: userID})
	return total, err
}

func (m *memAccountRepo) CountLowEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error) {
	_ = ctx
	return 0, nil
}

type memAuditRepo struct {
	mu   sync.Mutex
	logs []domain.AuditLog
}

func (m *memAuditRepo) Create(ctx context.Context, log *domain.AuditLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	m.logs = append(m.logs, *log)
	return nil
}

func TestLeadsHTTP_StatusMachineAndConvert(t *testing.T) {
	env := setupLeadsHTTPEnv(t)

	w := env.post("/api/leads", map[string]any{"title": "华创线索", "source": "web"})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create status %d body %s", w.Code, w.Body.String())
	}
	leadID := env.parseDataID(w)

	w = env.patch(leadID, map[string]any{"status": "contacted"})
	env.assertOK(w)
	env.assertStatus(w, "contacted")

	w = env.patch(leadID, map[string]any{"status": "qualified"})
	env.assertOK(w)

	w = env.patch(leadID, map[string]any{"status": "new"})
	env.assertBadMessage(w, "invalid_status_transition")

	w = env.post("/api/leads/"+leadID+"/convert", map[string]any{
		"create_account": map[string]any{"name": "华创科技有限公司"},
	})
	env.assertOK(w)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	if data["status"] != "converted" || data["converted_account_id"] == nil {
		t.Fatalf("convert data: %v", data)
	}

	env.auditRepo.mu.Lock()
	var convertCount int
	tenantAID := uuid.MustParse(env.tenantA)
	for _, log := range env.auditRepo.logs {
		if log.TenantID == tenantAID && log.Action == "lead.convert" {
			convertCount++
		}
	}
	env.auditRepo.mu.Unlock()
	if convertCount < 1 {
		t.Fatal("expected lead.convert audit log")
	}

	w = env.patch(leadID, map[string]any{"status": "contacted"})
	env.assertBadMessage(w, "invalid_status_transition")
}

func TestLeadsHTTP_ConvertRequiresAccount(t *testing.T) {
	env := setupLeadsHTTPEnv(t)

	w := env.post("/api/leads", map[string]any{"title": "X", "status": "qualified"})
	leadID := env.parseDataID(w)

	w = env.post("/api/leads/"+leadID+"/convert", map[string]any{})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d %s", w.Code, w.Body.String())
	}
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "convert_requires_account" {
		t.Fatalf("message: %v", body["message"])
	}
}

func TestLeadsHTTP_StatsEndpoints(t *testing.T) {
	env := setupLeadsHTTPEnv(t)

	w := env.post("/api/leads", map[string]any{"title": "A", "source": "web", "status": "new"})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create: %d %s", w.Code, w.Body.String())
	}
	w = env.post("/api/leads", map[string]any{"title": "B", "source": "web", "status": "contacted"})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create: %d %s", w.Code, w.Body.String())
	}

	for _, path := range []string{
		"/api/leads/stats/by-source",
		"/api/leads/stats/by-status",
		"/api/leads/stats/trend",
		"/api/leads/stats/funnel",
	} {
		w = env.getWithTenant(path, env.tenantA)
		if w.Code != http.StatusOK {
			t.Fatalf("%s status %d body %s", path, w.Code, w.Body.String())
		}
	}

	w = env.getWithTenant("/api/leads/stats/by-source", env.tenantA)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	total, _ := data["total"].(float64)
	if total < 2 {
		t.Fatalf("expected total >= 2, got %v", data)
	}
}

func TestLeadsHTTP_TenantIsolation(t *testing.T) {
	env := setupLeadsHTTPEnv(t)

	w := env.post("/api/leads", map[string]any{"title": "Tenant A Lead"})
	leadID := env.parseDataID(w)

	w = env.getWithTenant("/api/leads/"+leadID, env.tenantB)
	// 跨租户：无 B 域策略 → 403；或有策略但无数据 → 404（均不泄露 A 租户数据）
	if w.Code != http.StatusNotFound && w.Code != http.StatusForbidden {
		t.Fatalf("cross-tenant get: status %d body %s", w.Code, w.Body.String())
	}

	w = env.getWithTenant("/api/leads", env.tenantB)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data, _ := body["data"].(map[string]any)
	items, _ := data["items"].([]any)
	if len(items) != 0 {
		t.Fatalf("tenant B list should be empty, got %d items", len(items))
	}
}

type leadsHTTPEnv struct {
	t         *testing.T
	router    *gin.Engine
	leadRepo  *memLeadRepo
	auditRepo *memAuditRepo
	tenantA   string
	tenantB   string
	token     string
}

func setupLeadsHTTPEnv(t *testing.T) *leadsHTTPEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	tenantB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	enforcer := newLeadsTestEnforcer(t, userA, tenantA)
	leadRepo := &memLeadRepo{items: map[uuid.UUID]*domain.Lead{}}
	accountRepo := &memAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	auditRepo := &memAuditRepo{}
	auditRec := audit.NewRecorder(auditRepo)
	activityRepo := &memActivityRepo{items: map[uuid.UUID]*domain.Activity{}}
	leadSvc := leadapp.NewService(leadRepo, accountRepo, activityRepo, nil, enforcer, nil)
	emotionSvc := emotionapp.NewService(activityRepo)
	leadHTTP := httphandler.NewLeadHandlers(leadSvc, auditRec, emotionSvc)

	secret := "leads-test-secret"
	token, _, err := jwtutil.GenerateAccess(secret, userA, "sales@test.com", false, &tenantA, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(secret))
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RBACMiddleware(enforcer))
	registerLeadsRoutes(api, leadHTTP)

	return &leadsHTTPEnv{
		t:         t,
		router:    r,
		leadRepo:  leadRepo,
		auditRepo: auditRepo,
		tenantA:   tenantA.String(),
		tenantB:   tenantB.String(),
		token:     token,
	}
}

func newLeadsTestEnforcer(t *testing.T, userID, tenantID uuid.UUID) *casbin.Enforcer {
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
		_, _ = e.AddPolicy(role, dom, "leads", act)
	}
	for _, act := range []string{"view", "create", "update", "delete"} {
		_, _ = e.AddPolicy(role, dom, "activities", act)
	}
	for _, act := range []string{"view", "create", "update", "delete"} {
		_, _ = e.AddPolicy(role, dom, "contacts", act)
	}
	_, _ = e.AddGroupingPolicy(userID.String(), role, dom)
	// Tenant B: same user may switch tenant; grant policies so RBAC passes and repo returns empty/404
	domB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb").String()
	for _, act := range []string{"view", "create", "update", "delete"} {
		_, _ = e.AddPolicy(role, domB, "leads", act)
	}
	for _, act := range []string{"view", "create", "update", "delete"} {
		_, _ = e.AddPolicy(role, domB, "activities", act)
	}
	for _, act := range []string{"view", "create", "update", "delete"} {
		_, _ = e.AddPolicy(role, domB, "contacts", act)
	}
	_, _ = e.AddGroupingPolicy(userID.String(), role, domB)
	return e
}

func (e *leadsHTTPEnv) request(method, path, tenantID string, body any) *httptest.ResponseRecorder {
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

func (e *leadsHTTPEnv) post(path string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPost, path, e.tenantA, body)
}

func (e *leadsHTTPEnv) patch(leadID string, body any) *httptest.ResponseRecorder {
	return e.request(http.MethodPatch, "/api/leads/"+leadID, e.tenantA, body)
}

func (e *leadsHTTPEnv) getWithTenant(path, tenantID string) *httptest.ResponseRecorder {
	return e.request(http.MethodGet, path, tenantID, nil)
}

func (e *leadsHTTPEnv) parseDataID(w *httptest.ResponseRecorder) string {
	e.t.Helper()
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		e.t.Fatalf("parse: %v body %s", err, w.Body.String())
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		e.t.Fatalf("no data: %s", w.Body.String())
	}
	id, _ := data["id"].(string)
	if id == "" {
		e.t.Fatalf("no id in data: %v", data)
	}
	return id
}

func (e *leadsHTTPEnv) assertOK(w *httptest.ResponseRecorder) {
	e.t.Helper()
	if w.Code != http.StatusOK {
		e.t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func (e *leadsHTTPEnv) assertStatus(w *httptest.ResponseRecorder, want string) {
	e.t.Helper()
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	if data["status"] != want {
		e.t.Fatalf("status %v want %s", data["status"], want)
	}
}

func (e *leadsHTTPEnv) assertBadMessage(w *httptest.ResponseRecorder, msg string) {
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
