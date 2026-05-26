package rbac

import (
	"strings"
)

var methodToAction = map[string]string{
	"GET":    "view",
	"POST":   "create",
	"PUT":    "update",
	"PATCH":  "update",
	"DELETE": "delete",
}

// RouteToPermission 将 HTTP 方法与路径映射为 resource + action
func RouteToPermission(method, path string) (resource, action string) {
	path = strings.TrimPrefix(path, "/api/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 && parts[0] != "" {
		resource = parts[0]
	}
	if resource == "" {
		resource = "unknown"
	}

	// Dashboard (phase-3-deals-dashboard-api.md §4)
	if parts[0] == "dashboard" {
		return "dashboard", "view"
	}

	// RBAC admin (roles, members, user role assignment)
	if parts[0] == "rbac" {
		action = methodToAction[method]
		if action == "create" || action == "update" {
			action = "manage"
		}
		if action == "" {
			action = "view"
		}
		return "rbac", action
	}

	// Deals pipeline / stage / stats (phase-3-deals-dashboard-api.md)
	if len(parts) >= 2 && parts[0] == "deals" {
		if parts[1] == "pipeline" {
			return "deals", "view"
		}
		if parts[1] == "stats" {
			return "deals", "view"
		}
		if len(parts) >= 3 && parts[2] == "stage" {
			return "deals", "update"
		}
	}

	// Leads stats (phase-2-crm-ai.md §10)
	if len(parts) >= 3 && parts[1] == "stats" {
		return "leads", "view"
	}

	// Activities summary (phase-2-crm-ai.md §8)
	if len(parts) >= 2 && parts[0] == "activities" && parts[1] == "summary" {
		return "activities", "view"
	}

	// Segments (phase-2-crm-ai.md §9)
	if parts[0] == "segments" {
		if method == "PATCH" {
			return "segments", "manage"
		}
		if len(parts) >= 3 && parts[2] == "count" {
			return "segments", "view"
		}
		return "segments", "view"
	}

	// Sub-resource routes (phase-2-crm-ai.md)
	if len(parts) >= 3 {
		switch parts[2] {
		case "emotion-journey":
			return resource, "view"
		case "insights":
			if len(parts) >= 4 && parts[3] == "evaluate" {
				return resource, "view"
			}
		case "convert", "assign":
			if resource == "leads" {
				if parts[2] == "assign" {
					return "leads", "assign"
				}
				return "leads", "update"
			}
		case "contacts":
			return resource, "view"
		case "lifecycle":
			return resource, "view"
		}
	}

	action = methodToAction[method]
	if action == "" {
		action = "view"
	}
	return resource, action
}
