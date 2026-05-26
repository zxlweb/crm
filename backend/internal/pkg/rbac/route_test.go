package rbac

import "testing"

func TestRouteToPermission(t *testing.T) {
	tests := []struct {
		method   string
		path     string
		resource string
		action   string
	}{
		{"GET", "/api/leads", "leads", "view"},
		{"POST", "/api/leads", "leads", "create"},
		{"PUT", "/api/leads/123", "leads", "update"},
		{"DELETE", "/api/deals/1", "deals", "delete"},
		{"GET", "/api/rbac/permissions", "rbac", "view"},
		{"GET", "/api/rbac/members", "rbac", "view"},
		{"POST", "/api/rbac/roles", "rbac", "manage"},
		{"PUT", "/api/rbac/members/uuid/roles", "rbac", "manage"},
		{"POST", "/api/auth/login", "auth", "create"},
		{"GET", "/api/accounts/uuid/emotion-journey", "accounts", "view"},
		{"POST", "/api/accounts/uuid/insights/evaluate", "accounts", "view"},
		{"POST", "/api/leads/uuid/assign", "leads", "assign"},
		{"GET", "/api/leads/stats/by-source", "leads", "view"},
		{"GET", "/api/leads/stats/trend", "leads", "view"},
		{"GET", "/api/deals/pipeline", "deals", "view"},
		{"PUT", "/api/deals/uuid/stage", "deals", "update"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			resource, action := RouteToPermission(tt.method, tt.path)
			if resource != tt.resource || action != tt.action {
				t.Fatalf("got %s:%s, want %s:%s", resource, action, tt.resource, tt.action)
			}
		})
	}
}

func TestRouteToPermission_UnknownMethodDefaultsView(t *testing.T) {
	resource, action := RouteToPermission("OPTIONS", "/api/leads")
	if resource != "leads" || action != "view" {
		t.Fatalf("got %s:%s", resource, action)
	}
}
