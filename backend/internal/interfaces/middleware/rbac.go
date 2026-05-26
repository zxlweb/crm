package middleware

import (
	"net/http"
	"strings"

	"crm-backend/internal/pkg/rbac"
	"crm-backend/internal/pkg/rbacutil"
	"crm-backend/internal/pkg/response"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool("is_super_admin") {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		// 登录后拉取当前用户权限/角色，或自检权限；不依赖 Casbin 路由策略（避免普通角色被误拦）
		if c.Request.Method == http.MethodGet &&
			(strings.HasSuffix(path, "/rbac/my-permissions") || strings.HasSuffix(path, "/rbac/my-roles")) {
			c.Next()
			return
		}
		if c.Request.Method == http.MethodPost && strings.HasSuffix(path, "/rbac/check") {
			c.Next()
			return
		}

		userID := c.GetString("user_id")
		tenantID := c.GetString("tenant_id")
		resource, action := rbac.RouteToPermission(c.Request.Method, c.Request.URL.Path)

		ok, err := rbacutil.Enforce(c.Request.Context(), enforcer, userID, tenantID, resource, action)
		if err != nil || !ok {
			response.Forbidden(c, "无权限访问该资源")
			c.Abort()
			return
		}

		c.Next()
	}
}
