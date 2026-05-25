package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEmotionJourneyHTTP_WithActivities(t *testing.T) {
	env, _ := setupActivitiesHTTPEnv(t)

	w := env.post("/api/leads", map[string]any{"title": "Journey Lead"})
	leadID := env.parseDataID(w)

	occurred := time.Now().UTC().Add(-2 * 24 * time.Hour).Format(time.RFC3339)
	w = env.post("/api/activities", map[string]any{
		"subject_type": "lead",
		"subject_id":   leadID,
		"event_type":   "call",
		"direction":    "outbound",
		"body":         "客户犹豫价格",
		"sentiment":    "hesitant",
		"occurred_at":  occurred,
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create activity: %d %s", w.Code, w.Body.String())
	}

	w = env.get("/api/leads/" + leadID + "/emotion-journey")
	env.assertOK(w)
	data := parseJourneyData(t, w)
	if len(data["points"].([]any)) == 0 {
		t.Fatalf("expected points, got %v", data["points"])
	}
	if data["summary"].(map[string]any)["trend"] == "" {
		t.Fatal("missing trend")
	}
}

func TestEmotionJourneyHTTP_EmptyStill200(t *testing.T) {
	env, _ := setupActivitiesHTTPEnv(t)
	w := env.post("/api/leads", map[string]any{"title": "Empty Journey"})
	leadID := env.parseDataID(w)

	w = env.get("/api/leads/" + leadID + "/emotion-journey")
	env.assertOK(w)
	data := parseJourneyData(t, w)
	if len(data["points"].([]any)) != 0 {
		t.Fatalf("expected empty points: %v", data["points"])
	}
	if data["summary"].(map[string]any)["trend"] != "flat" {
		t.Fatalf("trend: %v", data["summary"])
	}
}

func TestEmotionJourneyHTTP_Range30dFiltersOld(t *testing.T) {
	env, _ := setupActivitiesHTTPEnv(t)
	w := env.post("/api/leads", map[string]any{"title": "Range Lead"})
	leadID := env.parseDataID(w)

	old := time.Now().UTC().Add(-40 * 24 * time.Hour).Format(time.RFC3339)
	w = env.post("/api/activities", map[string]any{
		"subject_type": "lead", "subject_id": leadID, "event_type": "note",
		"sentiment": "negative", "occurred_at": old,
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create old activity: %d", w.Code)
	}

	recent := time.Now().UTC().Add(-2 * 24 * time.Hour).Format(time.RFC3339)
	w = env.post("/api/activities", map[string]any{
		"subject_type": "lead", "subject_id": leadID, "event_type": "call",
		"sentiment": "neutral", "occurred_at": recent,
	})
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("create recent activity: %d", w.Code)
	}

	w = env.get("/api/leads/" + leadID + "/emotion-journey?range=30d")
	env.assertOK(w)
	data := parseJourneyData(t, w)
	if len(data["points"].([]any)) != 1 {
		t.Fatalf("30d points: %v", data["points"])
	}

	w = env.get("/api/leads/" + leadID + "/emotion-journey?range=all")
	env.assertOK(w)
	data = parseJourneyData(t, w)
	if len(data["points"].([]any)) != 2 {
		t.Fatalf("all points: %v", data["points"])
	}
}

func TestEmotionJourneyHTTP_TenantIsolation(t *testing.T) {
	env, _ := setupActivitiesHTTPEnv(t)
	w := env.post("/api/leads", map[string]any{"title": "Tenant A Journey"})
	leadID := env.parseDataID(w)

	w = env.getWithTenant("/api/leads/"+leadID+"/emotion-journey", env.tenantB)
	if w.Code != http.StatusNotFound && w.Code != http.StatusForbidden {
		t.Fatalf("cross-tenant status %d %s", w.Code, w.Body.String())
	}
}

func TestEmotionJourneyHTTP_InvalidRange(t *testing.T) {
	env, _ := setupActivitiesHTTPEnv(t)
	w := env.post("/api/leads", map[string]any{"title": "Bad Range"})
	leadID := env.parseDataID(w)

	w = env.get("/api/leads/" + leadID + "/emotion-journey?range=1y")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestEmotionJourneyHTTP_ConvertedMilestone(t *testing.T) {
	env, _ := setupActivitiesHTTPEnv(t)
	w := env.post("/api/leads", map[string]any{"title": "Convert Lead", "status": "qualified"})
	leadID := env.parseDataID(w)

	w = env.post("/api/leads/"+leadID+"/convert", map[string]any{
		"create_account": map[string]any{"name": "Converted Co"},
	})
	env.assertOK(w)

	w = env.get("/api/leads/" + leadID + "/emotion-journey")
	env.assertOK(w)
	data := parseJourneyData(t, w)
	milestones, ok := data["milestones"].([]any)
	if !ok || len(milestones) == 0 {
		t.Fatalf("expected converted milestone: %v", data["milestones"])
	}
	first := milestones[0].(map[string]any)
	if first["type"] != "converted" {
		t.Fatalf("milestone type: %v", first["type"])
	}
}

func parseJourneyData(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("no data: %s", w.Body.String())
	}
	return data
}

func (e *leadsHTTPEnv) get(path string) *httptest.ResponseRecorder {
	return e.request(http.MethodGet, path, e.tenantA, nil)
}
