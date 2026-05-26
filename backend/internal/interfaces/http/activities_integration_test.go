package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"crm-backend/internal/application/audit"
	actapp "crm-backend/internal/application/activity"
	emotionapp "crm-backend/internal/application/emotion"
	leadapp "crm-backend/internal/application/lead"
	"crm-backend/internal/domain"
	httphandler "crm-backend/internal/interfaces/http"
	"crm-backend/internal/interfaces/middleware"
	"crm-backend/internal/pkg/jwtutil"
	"crm-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setupActivitiesHTTPEnv(t *testing.T) (*leadsHTTPEnv, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	enforcer := newLeadsTestEnforcer(t, userA, tenantA)
	leadRepo := &memLeadRepo{items: map[uuid.UUID]*domain.Lead{}}
	accountRepo := &memAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	activityRepo := &memActivityRepo{items: map[uuid.UUID]*domain.Activity{}}
	auditRec := audit.NewRecorder(&memAuditRepo{})
	leadSvc := leadapp.NewService(leadRepo, accountRepo, activityRepo, nil, enforcer, nil)
	contactRepo := &memContactRepo{items: map[uuid.UUID]*domain.Contact{}}
	actSvc := actapp.NewService(activityRepo, leadRepo, accountRepo, contactRepo, nil, enforcer)
	emotionSvc := emotionapp.NewService(activityRepo)
	leadHTTP := httphandler.NewLeadHandlers(leadSvc, auditRec, emotionSvc)
	actHTTP := httphandler.NewActivityHandlers(actSvc, auditRec)
	secret := "activities-test-secret"
	token, _, err := jwtutil.GenerateAccess(secret, userA, "sales@test.com", false, &tenantA, nil, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(secret))
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RBACMiddleware(enforcer))
	registerLeadsRoutes(api, leadHTTP)
	api.GET("/activities/summary", actHTTP.Summary)
	api.GET("/activities", actHTTP.List)
	api.POST("/activities", actHTTP.Create)
	api.GET("/activities/:id", actHTTP.Get)
	api.PATCH("/activities/:id", actHTTP.Patch)
	api.DELETE("/activities/:id", actHTTP.Delete)
	env := &leadsHTTPEnv{
		t:        t,
		router:   r,
		leadRepo: leadRepo,
		tenantA:  tenantA.String(),
		tenantB:  uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb").String(),
		token:    token,
	}
	return env, tenantA.String()
}

type memActivityRepo struct {
	items map[uuid.UUID]*domain.Activity
}

func (m *memActivityRepo) List(ctx context.Context, tenantID uuid.UUID, f repository.ActivityListFilter) ([]domain.Activity, int64, error) {
	var out []domain.Activity
	for _, a := range m.items {
		if a.TenantID != tenantID || a.SubjectType != f.SubjectType || a.SubjectID != f.SubjectID {
			continue
		}
		out = append(out, *a)
	}
	return out, int64(len(out)), nil
}

func (m *memActivityRepo) ListForJourney(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID, since *time.Time) ([]domain.Activity, error) {
	items, _, err := m.List(ctx, tenantID, repository.ActivityListFilter{
		SubjectType: subjectType,
		SubjectID:   subjectID,
		Page:        1,
		PageSize:    10000,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Activity, 0, len(items))
	for _, a := range items {
		if since != nil && a.OccurredAt.Before(*since) {
			continue
		}
		out = append(out, a)
	}
	// ASC by occurred_at
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].OccurredAt.Before(out[i].OccurredAt) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out, nil
}

func (m *memActivityRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Activity, error) {
	a, ok := m.items[id]
	if !ok || a.TenantID != tenantID {
		return nil, repository.ErrActivityNotFound
	}
	cp := *a
	return &cp, nil
}

func (m *memActivityRepo) Create(ctx context.Context, a *domain.Activity) error {
	if m.items == nil {
		m.items = map[uuid.UUID]*domain.Activity{}
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

func (m *memActivityRepo) Update(ctx context.Context, a *domain.Activity) error {
	cp := *a
	m.items[a.ID] = &cp
	return nil
}

func (m *memActivityRepo) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	if _, ok := m.items[id]; !ok {
		return repository.ErrActivityNotFound
	}
	delete(m.items, id)
	return nil
}

func (m *memActivityRepo) SummaryByEventType(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID) ([]repository.LabelCount, int64, error) {
	items, total, err := m.List(ctx, tenantID, repository.ActivityListFilter{SubjectType: subjectType, SubjectID: subjectID, Page: 1, PageSize: 1000})
	if err != nil {
		return nil, 0, err
	}
	counts := map[string]int64{}
	for _, a := range items {
		counts[a.EventType]++
	}
	out := make([]repository.LabelCount, 0, len(counts))
	for k, v := range counts {
		out = append(out, repository.LabelCount{Label: k, Count: v})
	}
	return out, total, nil
}

func (m *memActivityRepo) LatestOccurredAt(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID) (*time.Time, error) {
	items, _, _ := m.List(ctx, tenantID, repository.ActivityListFilter{SubjectType: subjectType, SubjectID: subjectID, Page: 1, PageSize: 1000})
	if len(items) == 0 {
		return nil, nil
	}
	latest := items[0].OccurredAt
	for _, a := range items[1:] {
		if a.OccurredAt.After(latest) {
			latest = a.OccurredAt
		}
	}
	return &latest, nil
}

func (m *memActivityRepo) CountBySubject(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID) (int64, error) {
	_, total, err := m.List(ctx, tenantID, repository.ActivityListFilter{SubjectType: subjectType, SubjectID: subjectID, Page: 1, PageSize: 1})
	return total, err
}

func (m *memActivityRepo) CountSince(ctx context.Context, tenantID uuid.UUID, since time.Time) (int64, error) {
	var n int64
	for _, a := range m.items {
		if a.TenantID == tenantID && !a.OccurredAt.Before(since) {
			n++
		}
	}
	return n, nil
}

func (m *memActivityRepo) CountLeadTouchesSince(ctx context.Context, tenantID uuid.UUID, since time.Time, viewAll bool, userID uuid.UUID) (int64, error) {
	_ = viewAll
	_ = userID
	var n int64
	for _, a := range m.items {
		if a.TenantID == tenantID && a.SubjectType == "lead" && !a.OccurredAt.Before(since) {
			n++
		}
	}
	return n, nil
}

func (m *memActivityRepo) CountAccountTouchesSince(ctx context.Context, tenantID uuid.UUID, since time.Time, viewAll bool, userID uuid.UUID) (int64, error) {
	_ = viewAll
	_ = userID
	var n int64
	for _, a := range m.items {
		if a.TenantID == tenantID && a.SubjectType == "account" && !a.OccurredAt.Before(since) {
			n++
		}
	}
	return n, nil
}

func TestActivitiesHTTP_CRUDAndSummary(t *testing.T) {
	env, tenant := setupActivitiesHTTPEnv(t)
	w := env.post("/api/leads", map[string]any{"title": "Activity Lead"})
	leadID := env.parseDataID(w)

	w = env.post("/api/activities", map[string]any{
		"subject_type": "lead",
		"subject_id":   leadID,
		"event_type":   "call",
		"direction":    "outbound",
		"body":         "首次电话沟通",
		"sentiment":    "positive",
		"label":        "电话：首次沟通",
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create activity: %d %s", w.Code, w.Body.String())
	}
	actID := env.parseDataID(w)

	w = env.getWithTenant("/api/activities?subject_type=lead&subject_id="+leadID, tenant)
	env.assertOK(w)
	var listBody map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &listBody)
	data := listBody["data"].(map[string]any)
	items := data["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("list items: %d", len(items))
	}

	w = env.getWithTenant("/api/activities/summary?subject_type=lead&subject_id="+leadID, tenant)
	env.assertOK(w)
	var sumBody map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &sumBody)
	sum := sumBody["data"].(map[string]any)
	if sum["total"].(float64) != 1 {
		t.Fatalf("summary total: %v", sum["total"])
	}

	w = env.request(http.MethodPatch, "/api/activities/"+actID, tenant, map[string]any{
		"sentiment": "hesitant",
	})
	env.assertOK(w)

	w = env.request(http.MethodDelete, "/api/activities/"+actID, tenant, nil)
	env.assertOK(w)

	w = env.getWithTenant("/api/leads/"+leadID, tenant)
	env.assertOK(w)
}

func TestActivitiesHTTP_SentimentKeywordRuleInference(t *testing.T) {
	env, tenant := setupActivitiesHTTPEnv(t)
	w := env.post("/api/leads", map[string]any{"title": "Keyword Lead"})
	leadID := env.parseDataID(w)

	w = env.post("/api/activities", map[string]any{
		"subject_type": "lead",
		"subject_id":   leadID,
		"event_type":   "call",
		"direction":    "outbound",
		"body":         "客户非常满意，愿意推进合同",
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create with keyword body: %d %s", w.Code, w.Body.String())
	}

	w = env.getWithTenant("/api/activities?subject_type=lead&subject_id="+leadID, tenant)
	env.assertOK(w)
	var listBody map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &listBody)
	data := listBody["data"].(map[string]any)
	items := data["items"].([]any)
	row, _ := items[0].(map[string]any)
	if row["sentiment"] != "positive" {
		t.Fatalf("inferred sentiment: %v", row["sentiment"])
	}
	if row["sentiment_source"] != "rule" {
		t.Fatalf("sentiment_source: %v", row["sentiment_source"])
	}
}

func TestActivitiesHTTP_SentimentValidation(t *testing.T) {
	env, tenant := setupActivitiesHTTPEnv(t)
	w := env.post("/api/leads", map[string]any{"title": "Sentiment Lead"})
	leadID := env.parseDataID(w)

	w = env.post("/api/activities", map[string]any{
		"subject_type": "lead",
		"subject_id":   leadID,
		"event_type":   "call",
		"direction":    "outbound",
		"sentiment":    "not-a-mood",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid sentiment: %d %s", w.Code, w.Body.String())
	}

	w = env.post("/api/activities", map[string]any{
		"subject_type":     "lead",
		"subject_id":       leadID,
		"event_type":       "call",
		"direction":        "outbound",
		"sentiment":        "positive",
		"sentiment_source": "ai",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("ai sentiment_source: %d %s", w.Code, w.Body.String())
	}

	w = env.post("/api/activities", map[string]any{
		"subject_type":     "lead",
		"subject_id":       leadID,
		"event_type":       "call",
		"direction":        "outbound",
		"body":             "情绪标注",
		"sentiment":        "negative",
		"sentiment_source": "manual",
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create with sentiment: %d %s", w.Code, w.Body.String())
	}
	actID := env.parseDataID(w)

	w = env.getWithTenant("/api/activities?subject_type=lead&subject_id="+leadID, tenant)
	env.assertOK(w)
	var listBody map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &listBody)
	data := listBody["data"].(map[string]any)
	items := data["items"].([]any)
	row, _ := items[0].(map[string]any)
	if row["sentiment"] != "negative" {
		t.Fatalf("sentiment: %v", row["sentiment"])
	}

	w = env.request(http.MethodPatch, "/api/activities/"+actID, tenant, map[string]any{
		"sentiment": "positive",
	})
	env.assertOK(w)

	w = env.getWithTenant("/api/activities/"+actID, tenant)
	env.assertOK(w)
	var getBody map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &getBody)
	got := getBody["data"].(map[string]any)
	if got["sentiment"] != "positive" {
		t.Fatalf("patched sentiment: %v", got["sentiment"])
	}
}

func TestActivitiesHTTP_SentimentKeywordRule_TenantIsolation(t *testing.T) {
	env, tenantA := setupActivitiesHTTPEnv(t)
	tenantB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb").String()

	w := env.post("/api/leads", map[string]any{"title": "Sentiment Lead"})
	leadID := env.parseDataID(w)

	w = env.post("/api/activities", map[string]any{
		"subject_type": "lead",
		"subject_id":   leadID,
		"event_type":   "call",
		"direction":    "outbound",
		"body":         "客户觉得太贵了，要再考虑一下",
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create activity: %d %s", w.Code, w.Body.String())
	}
	actID := env.parseDataID(w)

	var createBody map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &createBody)
	data := createBody["data"].(map[string]any)
	if data["sentiment"] != "hesitant" {
		t.Fatalf("sentiment want hesitant, got %v", data["sentiment"])
	}
	if data["sentiment_source"] != "rule" {
		t.Fatalf("sentiment_source want rule, got %v", data["sentiment_source"])
	}

	w = env.getWithTenant("/api/activities/"+actID, tenantB)
	if w.Code != http.StatusNotFound {
		t.Fatalf("tenant B expected 404 on activity get, got %d", w.Code)
	}
	_ = tenantA
}

func TestActivitiesHTTP_RequiresViewPermission(t *testing.T) {
	tenantA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userA := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	enforcer := newLeadsTestEnforcerActs(t, userA, tenantA, "create")
	leadRepo := &memLeadRepo{items: map[uuid.UUID]*domain.Lead{}}
	accountRepo := &memAccountRepo{items: map[uuid.UUID]*domain.Account{}}
	activityRepo := &memActivityRepo{items: map[uuid.UUID]*domain.Activity{}}
	contactRepo := &memContactRepo{items: map[uuid.UUID]*domain.Contact{}}
	actSvc := actapp.NewService(activityRepo, leadRepo, accountRepo, contactRepo, nil, enforcer)
	actHTTP := httphandler.NewActivityHandlers(actSvc, audit.NewRecorder(&memAuditRepo{}))
	secret := "act-no-view"
	token, _, _ := jwtutil.GenerateAccess(secret, userA, "u@test.com", false, &tenantA, nil, time.Hour)
	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(secret))
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RBACMiddleware(enforcer))
	api.GET("/activities", actHTTP.List)
	w := httptest.NewRequest(http.MethodGet, "/api/activities?subject_type=lead&subject_id="+uuid.New().String(), nil)
	w.Header.Set("Authorization", "Bearer "+token)
	w.Header.Set("X-Tenant-ID", tenantA.String())
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, w)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}
