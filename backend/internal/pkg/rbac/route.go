package rbac

import "strings"

var methodToAction = map[string]string{
	"GET":    "view",
	"POST":   "create",
	"PUT":    "update",
	"PATCH":  "update",
	"DELETE": "delete",
}

// RouteToPermission 将 HTTP 方法与路径映射为 resource + action
func RouteToPermission(method, path string) (resource, action string) {
	action = methodToAction[method]
	if action == "" {
		action = "view"
	}

	path = strings.TrimPrefix(path, "/api/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 && parts[0] != "" {
		resource = parts[0]
	}
	if resource == "" {
		resource = "unknown"
	}
	return resource, action
}
