package middleware

import (
	"strings"

	"crm-backend/internal/pkg/activerole"
	"crm-backend/internal/pkg/jwtutil"
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供 Token")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "Token 格式错误")
			c.Abort()
			return
		}

		claims, err := jwtutil.Parse(jwtSecret, parts[1], jwtutil.TokenTypeAccess)
		if err != nil {
			response.Unauthorized(c, "Token 无效或已过期")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("is_super_admin", claims.IsSuperAdmin)
		if claims.TenantID != "" {
			c.Set("jwt_tenant_id", claims.TenantID)
		}
		if claims.ActiveRoleID != "" {
			c.Set("active_role_id", claims.ActiveRoleID)
		}

		ctx := activerole.WithID(c.Request.Context(), claims.ActiveRoleID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
