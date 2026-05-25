package http_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"crm-backend/internal/application/audit"
	emotionapp "crm-backend/internal/application/emotion"
	leadapp "crm-backend/internal/application/lead"
	"crm-backend/internal/domain"
	httphandler "crm-backend/internal/interfaces/http"
	"crm-backend/internal/interfaces/middleware"
	"crm-backend/internal/pkg/jwtutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestLeadsHTTP_InsightsEvaluate_INS001(t *testing.T) {
	env := setupLeadsInsightsHTTPEnv(t)

	old := time.Now().UTC().Add(-10 * 24 * time.Hour)
	leadID := uuid.New()
	env.leadRepo.items[leadID] = &domain.Lead{
		ID:             leadID,
		TenantID:       uuid.MustParse(env.tenantA),
		Title:          "Silent Lead",
		Status:         "new",
		LifecycleStage: "acquire",
		LastActivityAt: &old,
		Timestamps: domain.Timestamps{
			CreatedAt: time.Now().UTC().Add(-30 * 24 * time.Hour),
			UpdatedAt: time.Now().UTC(),
		},
	}

	w := env.post("/api/leads/"+leadID.String()+"/insights/evaluate", map[string]any{})
	env.assertOK(w)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	items := data["items"].([]any)
	if len(items) == 0 {
		t.Fatalf("expected insights, got %v", data)
	}
	first := items[0].(map[string]any)
	if first["rule_id"] != "INS-001" {
		t.Fatalf("expected INS-001, got %v", first["rule_id"])
	}

	w = env.request(http.MethodPost, "/api/leads/"+leadID.String()+"/insights/evaluate", env.tenantB, map[string]any{})
	if w.Code != http.StatusNotFound {
		t.Fatalf("tenant B expected 404, got %d %s", w.Code, w.Body.String())
	}
}

func setupLeadsInsightsHTTPEnv(t *testing.T) *leadsHTTPEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	tenantB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	enforcer := newLeadsTestEnforcer(t, userA, tenantA)
	leadRepo := &memLeadRepo{items: map[uuid.UUID]*domain.Lead{}}
	accountRepo := &memAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	activityRepo := &memActivityRepo{items: map[uuid.UUID]*domain.Activity{}}
	auditRec := audit.NewRecorder(&memAuditRepo{})

	leadSvc := leadapp.NewService(leadRepo, accountRepo, activityRepo, nil, enforcer, nil)
	emotionSvc := emotionapp.NewService(activityRepo)
	leadHTTP := httphandler.NewLeadHandlers(leadSvc, auditRec, emotionSvc)

	secret := "insights-test-secret"
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
	api.POST("/leads/:id/insights/evaluate", leadHTTP.EvaluateInsights)

	return &leadsHTTPEnv{
		t:        t,
		router:   r,
		leadRepo: leadRepo,
		tenantA:  tenantA.String(),
		tenantB:  tenantB.String(),
		token:    token,
	}
}
