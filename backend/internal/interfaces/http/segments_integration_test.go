package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"crm-backend/internal/application/audit"
	emotionapp "crm-backend/internal/application/emotion"
	leadapp "crm-backend/internal/application/lead"
	segmentapp "crm-backend/internal/application/segment"
	"crm-backend/internal/domain"
	httphandler "crm-backend/internal/interfaces/http"
	"crm-backend/internal/interfaces/middleware"
	"crm-backend/internal/pkg/jwtutil"
	"crm-backend/internal/pkg/response"
	"crm-backend/internal/pkg/tenant"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (e *segmentsHTTPEnv) get(path string) *httptest.ResponseRecorder {
	return e.leadsHTTPEnv.request(http.MethodGet, path, e.tenantA, nil)
}

func (e *segmentsHTTPEnv) getWithTenant(path, tenant string) *httptest.ResponseRecorder {
	return e.leadsHTTPEnv.request(http.MethodGet, path, tenant, nil)
}

func TestSegmentsHTTP_ListTemplates(t *testing.T) {
	env := setupSegmentsHTTPEnv(t)
	w := env.get("/api/segments")
	env.assertOK(w)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	items := body["data"].(map[string]any)["items"].([]any)
	if len(items) != 5 {
		t.Fatalf("expected 5 segment templates, got %d", len(items))
	}
}

func TestSegmentsHTTP_CountMatchesLeadsList_LINT04b(t *testing.T) {
	env := setupSegmentsHTTPEnv(t)
	tenantID := uuid.MustParse(env.tenantA)

	for _, l := range []*domain.Lead{
		{TenantID: tenantID, Title: "High", Amount: 150000, Status: "new", LifecycleStage: "acquire"},
		{TenantID: tenantID, Title: "Low", Amount: 5000, Status: "new", LifecycleStage: "acquire"},
	} {
		if err := env.leadRepo.Create(context.Background(), l); err != nil {
			t.Fatal(err)
		}
	}

	w := env.get("/api/segments/high_value/count")
	env.assertOK(w)
	var countBody map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &countBody)
	count := int64(countBody["data"].(map[string]any)["count"].(float64))
	if count != 1 {
		t.Fatalf("high_value count: %d", count)
	}

	w = env.get("/api/leads?segment=high_value")
	env.assertOK(w)
	var listBody map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &listBody)
	total := int64(listBody["pagination"].(map[string]any)["total"].(float64))
	if total != count {
		t.Fatalf("list total %d != segment count %d", total, count)
	}
}

func TestSegmentsHTTP_ChurnRiskStaleLeads(t *testing.T) {
	env := setupSegmentsHTTPEnv(t)
	tenantID := uuid.MustParse(env.tenantA)
	old := time.Now().UTC().Add(-10 * 24 * time.Hour)
	recent := time.Now().UTC().Add(-1 * time.Hour)

	for i, last := range []*time.Time{&old, &old, &recent} {
		if err := env.leadRepo.Create(context.Background(), &domain.Lead{
			TenantID: tenantID, Title: "L" + strconv.Itoa(i), Status: "new",
			LifecycleStage: "acquire", LastActivityAt: last,
		}); err != nil {
			t.Fatal(err)
		}
	}

	w := env.get("/api/leads?segment=churn_risk")
	env.assertOK(w)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	total := int64(body["pagination"].(map[string]any)["total"].(float64))
	if total != 2 {
		t.Fatalf("churn_risk list total: %d", total)
	}
}

func TestSegmentsHTTP_TenantIsolation(t *testing.T) {
	env := setupSegmentsHTTPEnv(t)
	tenantA := uuid.MustParse(env.tenantA)

	_ = env.leadRepo.Create(context.Background(), &domain.Lead{
		TenantID: tenantA, Title: "A revive", LifecycleStage: "revive", Status: "new",
	})

	w := env.getWithTenant("/api/segments/revive_pool/count", env.tenantB)
	env.assertOK(w)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["data"].(map[string]any)["count"].(float64) != 0 {
		t.Fatalf("tenant B count: %v", body["data"].(map[string]any)["count"])
	}

	w = env.getWithTenant("/api/leads?segment=revive_pool", env.tenantB)
	env.assertOK(w)
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["pagination"].(map[string]any)["total"].(float64) != 0 {
		t.Fatalf("tenant B list: %v", body["pagination"].(map[string]any)["total"])
	}
}

func TestSegmentsHTTP_RequiresSegmentsView(t *testing.T) {
	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	enforcer := newLeadsTestEnforcer(t, userA, tenantA)
	segmentRepo := &memSegmentRepo{}
	leadRepo := &memLeadRepo{items: map[uuid.UUID]*domain.Lead{}}
	accountRepo := &memAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	segmentSvc := segmentapp.NewService(segmentRepo, leadRepo, accountRepo, nil, enforcer)
	segmentHTTP := httphandler.NewSegmentHandlers(segmentSvc)

	secret := "segments-no-view"
	token, _, _ := jwtutil.GenerateAccess(secret, userA, "u@test.com", false, &tenantA, nil, time.Hour)
	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(secret))
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RBACMiddleware(enforcer))
	api.GET("/segments", segmentHTTP.List)

	req := httptest.NewRequest(http.MethodGet, "/api/segments", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Tenant-ID", tenantA.String())
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d %s", rec.Code, rec.Body.String())
	}
}

func TestSegmentsHTTP_InvalidSegmentCode(t *testing.T) {
	env := setupSegmentsHTTPEnv(t)
	w := env.get("/api/segments/not-a-segment/count")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("count: expected 400, got %d", w.Code)
	}
	w = env.get("/api/leads?segment=not-a-segment")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("list: expected 400, got %d", w.Code)
	}
}

type segmentsHTTPEnv struct {
	*leadsHTTPEnv
	segmentRepo *memSegmentRepo
}

func setupSegmentsHTTPEnv(t *testing.T) *segmentsHTTPEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	tenantB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	enforcer := newSegmentsTestEnforcer(t, userA, tenantA)
	leadRepo := &memLeadRepo{items: map[uuid.UUID]*domain.Lead{}}
	accountRepo := &memAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	segmentRepo := &memSegmentRepo{}
	activityRepo := &memActivityRepo{items: map[uuid.UUID]*domain.Activity{}}
	auditRec := audit.NewRecorder(&memAuditRepo{})

	leadSvc := leadapp.NewService(leadRepo, accountRepo, activityRepo, nil, enforcer, nil)
	segmentSvc := segmentapp.NewService(segmentRepo, leadRepo, accountRepo, nil, enforcer)
	emotionSvc := emotionapp.NewService(activityRepo)
	leadHTTP := httphandler.NewLeadHandlers(leadSvc, auditRec, emotionSvc)
	segmentHTTP := httphandler.NewSegmentHandlers(segmentSvc)

	secret := "segments-test-secret"
	token, _, err := jwtutil.GenerateAccess(secret, userA, "sales@test.com", false, &tenantA, nil, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(secret))
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RBACMiddleware(enforcer))
	registerLeadsRoutesWithSegment(api, leadHTTP, leadSvc)
	api.GET("/segments", segmentHTTP.List)
	api.GET("/segments/:code/count", segmentHTTP.Count)

	base := &leadsHTTPEnv{
		t:        t,
		router:   r,
		leadRepo: leadRepo,
		tenantA:  tenantA.String(),
		tenantB:  tenantB.String(),
		token:    token,
	}
	return &segmentsHTTPEnv{leadsHTTPEnv: base, segmentRepo: segmentRepo}
}

func registerLeadsRoutesWithSegment(api *gin.RouterGroup, leadHTTP *httphandler.LeadHandlers, leadSvc *leadapp.Service) {
	api.GET("/leads/stats/by-source", leadHTTP.StatsBySource)
	api.GET("/leads/stats/by-status", leadHTTP.StatsByStatus)
	api.GET("/leads/stats/trend", leadHTTP.StatsTrend)
	api.GET("/leads/stats/funnel", leadHTTP.StatsFunnel)
	api.GET("/leads", func(c *gin.Context) {
		tid, okT := tenant.IDFromContext(c.Request.Context())
		if !okT {
			response.BadRequest(c, "缺少租户上下文")
			return
		}
		uid, err := uuid.Parse(c.GetString("user_id"))
		if err != nil {
			response.Unauthorized(c, "用户未认证")
			return
		}
		result, err := leadSvc.List(c.Request.Context(), tid, uid, leadapp.ListQuery{
			Page:               segmentQueryInt(c, "page", 1),
			PageSize:           segmentQueryInt(c, "page_size", 20),
			Search:             c.Query("search"),
			Status:             c.Query("status"),
			Source:             c.Query("source"),
			LifecycleStage:     c.Query("lifecycle_stage"),
			RelationshipHealth: c.Query("relationship_health"),
			Segment:            c.Query("segment"),
		})
		if err != nil {
			if errors.Is(err, leadapp.ErrInvalidSegment) {
				response.BadRequest(c, "invalid_segment_code")
				return
			}
			response.InternalError(c, "获取线索列表失败")
			return
		}
		response.SuccessPage(c, gin.H{"items": result.Items}, response.Pagination{
			Page: result.Page, PageSize: result.Size, Total: result.Total,
		})
	})
	api.POST("/leads", leadHTTP.Create)
	api.GET("/leads/:id", leadHTTP.Get)
	api.PATCH("/leads/:id", leadHTTP.Patch)
	api.POST("/leads/:id/convert", leadHTTP.Convert)
}

func newSegmentsTestEnforcer(t *testing.T, userID, tenantID uuid.UUID) *casbin.Enforcer {
	t.Helper()
	e := newLeadsTestEnforcer(t, userID, tenantID)
	role := userID.String()
	dom := tenantID.String()
	_, _ = e.AddPolicy(role, dom, "segments", "view")
	domB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb").String()
	_, _ = e.AddPolicy(role, domB, "segments", "view")
	return e
}

func segmentQueryInt(c *gin.Context, key string, def int) int {
	if v := c.Query(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
