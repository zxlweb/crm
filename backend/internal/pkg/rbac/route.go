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
