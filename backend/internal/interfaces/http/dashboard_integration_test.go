package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"crm-backend/internal/application/audit"
	dashapp "crm-backend/internal/application/dashboard"
	dealapp "crm-backend/internal/application/deal"
	"crm-backend/internal/domain"
	httphandler "crm-backend/internal/interfaces/http"
	"crm-backend/internal/interfaces/middleware"
	"crm-backend/internal/pkg/jwtutil"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestDashboardHTTP_SummaryAndDataScope(t *testing.T) {
	env := setupDashboardHTTPEnv(t, false)

	w := env.get("/api/dashboard/summary")
	if w.Code != http.StatusOK {
		t.Fatalf("summary status %d body %s", w.Code, w.Body.String())
	}
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	if data["data_scope"] != "self" {
		t.Fatalf("data_scope: %v", data["data_scope"])
	}
	kpis := data["kpis"].(map[string]any)
	if _, ok := kpis["leads_total"]; !ok {
		t.Fatalf("kpis: %v", kpis)
	}
	spark := data["sparklines"].(map[string]any)
	if len(spark["leads"].([]any)) != 7 {
		t.Fatalf("sparklines leads: %v", spark["leads"])
	}
}

func TestDashboardHTTP_TeamRankingRequiresManagerScope(t *testing.T) {
	sales := setupDashboardHTTPEnv(t, false)
	w := sales.get("/api/dashboard/team-ranking")
	if w.Code != http.StatusForbidden {
		t.Fatalf("sales team-ranking: status %d", w.Code)
	}

	manager := setupDashboardHTTPEnv(t, true)
	w = manager.get("/api/dashboard/team-ranking")
	if w.Code != http.StatusOK {
		t.Fatalf("manager team-ranking: status %d body %s", w.Code, w.Body.String())
	}
}

func TestDashboardHTTP_FunnelAndQuota(t *testing.T) {
	env := setupDashboardHTTPEnv(t, true)

	w := env.get("/api/dashboard/funnel?scope=deals")
	if w.Code != http.StatusOK {
		t.Fatalf("funnel: %d %s", w.Code, w.Body.String())
	}

	w = env.get("/api/dashboard/quota")
	if w.Code != http.StatusOK {
		t.Fatalf("quota: %d %s", w.Code, w.Body.String())
	}

	w = env.get("/api/dashboard/todo")
	if w.Code != http.StatusOK {
		t.Fatalf("todo: %d %s", w.Code, w.Body.String())
	}
}

type dashboardHTTPEnv struct {
	router  *gin.Engine
	tenantA string
	token   string
}

func setupDashboardHTTPEnv(t *testing.T, manager bool) *dashboardHTTPEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	enforcer := newDashboardTestEnforcer(t, userA, tenantA, manager)
	leadRepo := &memLeadRepo{items: map[uuid.UUID]*domain.Lead{}}
	accountRepo := &memAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	dealRepo := &memDealRepo{items: map[uuid.UUID]*domain.Deal{}}
	activityRepo := &memActivityRepo{items: map[uuid.UUID]*domain.Activity{}}

	dashSvc := dashapp.NewService(leadRepo, accountRepo, dealRepo, activityRepo, nil, nil, enforcer)
	dashHTTP := httphandler.NewDashboardHandlers(dashSvc, enforcer)
	dealSvc := dealapp.NewService(dealRepo, accountRepo, enforcer)
	dealHTTP := httphandler.NewDealHandlers(dealSvc, audit.NewRecorder(&memAuditRepo{}))

	secret := "dashboard-test-secret"
	token, _, err := jwtutil.GenerateAccess(secret, userA, "sales@test.com", false, &tenantA, nil, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(secret))
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RBACMiddleware(enforcer))
	registerDashboardRoutes(api, dashHTTP, dealHTTP)

	return &dashboardHTTPEnv{router: r, tenantA: tenantA.String(), token: token}
}

func registerDashboardRoutes(api *gin.RouterGroup, dashHTTP *httphandler.DashboardHandlers, dealHTTP *httphandler.DealHandlers) {
	api.GET("/dashboard/summary", dashHTTP.Summary)
	api.GET("/dashboard/funnel", dashHTTP.Funnel)
	api.GET("/dashboard/quota", dashHTTP.Quota)
	api.GET("/dashboard/team-ranking", dashHTTP.TeamRanking)
	api.GET("/dashboard/todo", dashHTTP.Todo)
	api.GET("/deals/stats/by-stage", dealHTTP.StatsByStage)
	api.GET("/deals/stats/win-rate", dealHTTP.StatsWinRate)
}

func newDashboardTestEnforcer(t *testing.T, userID, tenantID uuid.UUID, manager bool) *casbin.Enforcer {
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
	acts := []string{"view"}
	if manager {
		acts = append(acts, "create", "update", "delete")
	}
	for _, act := range acts {
		_, _ = e.AddPolicy(role, dom, "dashboard", act)
		_, _ = e.AddPolicy(role, dom, "deals", act)
		_, _ = e.AddPolicy(role, dom, "leads", act)
	}
	if manager {
		_, _ = e.AddPolicy(role, dom, "rbac", "manage")
	}
	_, _ = e.AddGroupingPolicy(userID.String(), role, dom)
	return e
}

func (e *dashboardHTTPEnv) get(path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, bytes.NewReader(nil))
	req.Header.Set("Authorization", "Bearer "+e.token)
	req.Header.Set("X-Tenant-ID", e.tenantA)
	w := httptest.NewRecorder()
	e.router.ServeHTTP(w, req)
	return w
}
