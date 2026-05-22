package middleware

import (
	"net/http"
	"strings"

	"crm-backend/internal/pkg/rbac"
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

		// 登录/注册后拉取当前用户权限集合，不依赖 Casbin 路由策略（避免新租户管理员被误拦）
		if c.Request.Method == http.MethodGet && strings.HasSuffix(c.Request.URL.Path, "/rbac/my-permissions") {
			c.Next()
			return
		}

		userID := c.GetString("user_id")
		tenantID := c.GetString("tenant_id")
		resource, action := rbac.RouteToPermission(c.Request.Method, c.Request.URL.Path)

		ok, err := enforcer.Enforce(userID, tenantID, resource, action)
		if err != nil || !ok {
			response.Forbidden(c, "无权限访问该资源")
			c.Abort()
			return
		}

		c.Next()
	}
}
