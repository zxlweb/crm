package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"crm-backend/internal/application/audit"
	emotionapp "crm-backend/internal/application/emotion"
	leadapp "crm-backend/internal/application/lead"
	"crm-backend/internal/domain"
	httphandler "crm-backend/internal/interfaces/http"
	"crm-backend/internal/interfaces/middleware"
	"crm-backend/internal/pkg/jwtutil"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type statsDistribution struct {
	total int64
	items []struct {
		label      string
		value      int64
		percentage float64
	}
}

type statsTrend struct {
	categories []string
	series     []struct {
		name string
		data []float64
	}
}

type statsFunnel struct {
	stages []struct {
		name  string
		count int64
	}
}

// L-STAT-01：统计与列表同源权限与数据范围
func TestLeadsHTTP_Stats_LSTAT01_BySourceMatchesListScope(t *testing.T) {
	env := setupLeadsHTTPEnv(t)
	env.post("/api/leads", map[string]any{"title": "A", "source": "website", "status": "new"})
	env.post("/api/leads", map[string]any{"title": "B", "source": "referral", "status": "contacted"})

	listTotal := env.parseListTotal(env.getWithTenant("/api/leads", env.tenantA))
	src := env.parseDistribution(env.getWithTenant("/api/leads/stats/by-source", env.tenantA))
	if src.total != listTotal {
		t.Fatalf("by-source total %d want list total %d", src.total, listTotal)
	}
	if src.total < 2 {
		t.Fatalf("expected at least 2 leads in stats, got %d", src.total)
	}
}

// L-STAT-02：趋势响应结构
func TestLeadsHTTP_Stats_LSTAT02_TrendShape(t *testing.T) {
	env := setupLeadsHTTPEnv(t)
	env.post("/api/leads", map[string]any{"title": "Trend Lead"})

	w := env.getWithTenant("/api/leads/stats/trend?granularity=day", env.tenantA)
	env.assertOK(w)
	tr := env.parseTrend(w)
	if len(tr.categories) == 0 {
		t.Fatal("expected non-empty categories")
	}
	if len(tr.series) == 0 || len(tr.series[0].data) != len(tr.categories) {
		t.Fatalf("series length mismatch: categories=%d series=%v", len(tr.categories), tr.series)
	}
}

// L-STAT-03：漏斗各阶段 count 之和与列表 total 一致
func TestLeadsHTTP_Stats_LSTAT03_FunnelSumMatchesList(t *testing.T) {
	env := setupLeadsHTTPEnv(t)
	env.post("/api/leads", map[string]any{"title": "F1", "status": "new"})
	env.post("/api/leads", map[string]any{"title": "F2", "status": "qualified"})

	listTotal := env.parseListTotal(env.getWithTenant("/api/leads", env.tenantA))
	funnel := env.parseFunnel(env.getWithTenant("/api/leads/stats/funnel", env.tenantA))
	var stageSum int64
	for _, s := range funnel.stages {
		stageSum += s.count
	}
	if stageSum != listTotal {
		t.Fatalf("funnel stage sum %d want list total %d", stageSum, listTotal)
	}
}

func TestLeadsHTTP_Stats_ByStatusMatchesList(t *testing.T) {
	env := setupLeadsHTTPEnv(t)
	env.post("/api/leads", map[string]any{"title": "S1", "status": "new"})
	env.post("/api/leads", map[string]any{"title": "S2", "status": "new"})

	listTotal := env.parseListTotal(env.getWithTenant("/api/leads", env.tenantA))
	byStatus := env.parseDistribution(env.getWithTenant("/api/leads/stats/by-status", env.tenantA))
	if byStatus.total != listTotal {
		t.Fatalf("by-status total %d want %d", byStatus.total, listTotal)
	}
}

func TestLeadsHTTP_Stats_TenantIsolation(t *testing.T) {
	env := setupLeadsHTTPEnv(t)
	env.post("/api/leads", map[string]any{"title": "Tenant A only"})

	w := env.getWithTenant("/api/leads/stats/by-source", env.tenantB)
	env.assertOK(w)
	src := env.parseDistribution(w)
	if src.total != 0 {
		t.Fatalf("tenant B stats total want 0, got %d", src.total)
	}
}

func TestLeadsHTTP_Stats_InvalidDateRange(t *testing.T) {
	env := setupLeadsHTTPEnv(t)
	w := env.getWithTenant("/api/leads/stats/by-source?from=2026-05-10&to=2026-05-01", env.tenantA)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d %s", w.Code, w.Body.String())
	}
}

func TestLeadsHTTP_Stats_RequiresLeadsView(t *testing.T) {
	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	enforcer := newLeadsTestEnforcerActs(t, userA, tenantA, "create")

	env := setupLeadsRouterWithEnforcer(t, enforcer, userA, tenantA)
	w := env.getWithTenant("/api/leads/stats/by-source", env.tenantA)
	if w.Code != http.StatusForbidden {
		t.Fatalf("stats expected 403, got %d %s", w.Code, w.Body.String())
	}
	w = env.getWithTenant("/api/leads", env.tenantA)
	if w.Code != http.StatusForbidden {
		t.Fatalf("list expected 403, got %d", w.Code)
	}
}

func TestLeadsHTTP_Stats_ViewerCanAccess(t *testing.T) {
	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	enforcer := newLeadsTestEnforcerActs(t, userA, tenantA, "view")

	env := setupLeadsRouterWithEnforcer(t, enforcer, userA, tenantA)
	env.assertOK(env.getWithTenant("/api/leads/stats/trend", env.tenantA))
}

func TestLeadsHTTP_Stats_OwnerScopeMatchesList(t *testing.T) {
	env := setupLeadsHTTPEnv(t)
	userB := uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")

	env.post("/api/leads", map[string]any{"title": "Owned by A"})
	w := env.post("/api/leads", map[string]any{"title": "Owned by B"})
	leadID := uuid.MustParse(env.parseDataID(w))

	env.leadRepo.mu.Lock()
	if l, ok := env.leadRepo.items[leadID]; ok {
		l.OwnerID = &userB
		cp := *l
		env.leadRepo.items[leadID] = &cp
	}
	env.leadRepo.mu.Unlock()

	listTotal := env.parseListTotal(env.getWithTenant("/api/leads", env.tenantA))
	if listTotal != 1 {
		t.Fatalf("scoped list total want 1, got %d", listTotal)
	}
	byStatus := env.parseDistribution(env.getWithTenant("/api/leads/stats/by-status", env.tenantA))
	if byStatus.total != 1 {
		t.Fatalf("scoped stats total want 1, got %d", byStatus.total)
	}
}

func setupLeadsRouterWithEnforcer(t *testing.T, enforcer *casbin.Enforcer, userID, tenantID uuid.UUID) *leadsHTTPEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)
	leadRepo := &memLeadRepo{items: map[uuid.UUID]*domain.Lead{}}
	accountRepo := &memAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	activityRepo := &memActivityRepo{items: map[uuid.UUID]*domain.Activity{}}
	leadSvc := leadapp.NewService(leadRepo, accountRepo, activityRepo, nil, enforcer, nil)
	emotionSvc := emotionapp.NewService(activityRepo)
	leadHTTP := httphandler.NewLeadHandlers(leadSvc, audit.NewRecorder(&memAuditRepo{}), emotionSvc)
	secret := "leads-test-custom"
	token, _, err := jwtutil.GenerateAccess(secret, userID, "u@test.com", false, &tenantID, nil, time.Hour)
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
		tenantA:   tenantID.String(),
		tenantB:   uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb").String(),
		token:     token,
	}
}

func registerLeadsRoutes(api *gin.RouterGroup, leadHTTP *httphandler.LeadHandlers) {
	api.GET("/leads/stats/by-source", leadHTTP.StatsBySource)
	api.GET("/leads/stats/by-status", leadHTTP.StatsByStatus)
	api.GET("/leads/stats/trend", leadHTTP.StatsTrend)
	api.GET("/leads/stats/funnel", leadHTTP.StatsFunnel)
	api.GET("/leads", leadHTTP.List)
	api.POST("/leads", leadHTTP.Create)
	api.GET("/leads/:id", leadHTTP.Get)
	api.PATCH("/leads/:id", leadHTTP.Patch)
	api.POST("/leads/:id/convert", leadHTTP.Convert)
	api.GET("/leads/:id/emotion-journey", leadHTTP.EmotionJourney)
}

func newLeadsTestEnforcerActs(t *testing.T, userID, tenantID uuid.UUID, acts ...string) *casbin.Enforcer {
	t.Helper()
	e := newLeadsTestEnforcer(t, userID, tenantID)
	// 清空默认全权限，仅保留指定 action
	role := userID.String()
	dom := tenantID.String()
	_, _ = e.RemoveFilteredPolicy(0, role, dom)
	for _, act := range acts {
		_, _ = e.AddPolicy(role, dom, "leads", act)
	}
	return e
}

func (e *leadsHTTPEnv) parseListTotal(w *httptest.ResponseRecorder) int64 {
	e.t.Helper()
	e.assertOK(w)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	pag, _ := body["pagination"].(map[string]any)
	total, _ := pag["total"].(float64)
	return int64(total)
}

func (e *leadsHTTPEnv) parseDistribution(w *httptest.ResponseRecorder) statsDistribution {
	e.t.Helper()
	e.assertOK(w)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	out := statsDistribution{}
	out.total = int64(data["total"].(float64))
	raw, _ := data["items"].([]any)
	for _, it := range raw {
		m := it.(map[string]any)
		out.items = append(out.items, struct {
			label      string
			value      int64
			percentage float64
		}{
			label:      m["label"].(string),
			value:      int64(m["value"].(float64)),
			percentage: m["percentage"].(float64),
		})
	}
	return out
}

func (e *leadsHTTPEnv) parseTrend(w *httptest.ResponseRecorder) statsTrend {
	e.t.Helper()
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	out := statsTrend{}
	for _, c := range data["categories"].([]any) {
		out.categories = append(out.categories, c.(string))
	}
	for _, s := range data["series"].([]any) {
		m := s.(map[string]any)
		var pts []float64
		for _, d := range m["data"].([]any) {
			pts = append(pts, d.(float64))
		}
		out.series = append(out.series, struct {
			name string
			data []float64
		}{name: m["name"].(string), data: pts})
	}
	return out
}

func (e *leadsHTTPEnv) parseFunnel(w *httptest.ResponseRecorder) statsFunnel {
	e.t.Helper()
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	out := statsFunnel{}
	for _, s := range data["stages"].([]any) {
		m := s.(map[string]any)
		out.stages = append(out.stages, struct {
			name  string
			count int64
		}{name: m["name"].(string), count: int64(m["count"].(float64))})
	}
	return out
}
