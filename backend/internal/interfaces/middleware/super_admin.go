package middleware

import (
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool("is_super_admin") {
			response.Forbidden(c, "需要 Super Admin 权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
