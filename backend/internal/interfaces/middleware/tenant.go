package middleware

import (
	"crm-backend/internal/pkg/response"
	"crm-backend/internal/pkg/tenant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantIDStr := c.GetHeader("X-Tenant-ID")
		if tenantIDStr == "" {
			tenantIDStr = c.GetString("jwt_tenant_id")
		}
		if tenantIDStr == "" {
			response.BadRequest(c, "缺少 X-Tenant-ID Header，请先切换租户")
			c.Abort()
			return
		}

		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			response.BadRequest(c, "租户 ID 格式无效")
			c.Abort()
			return
		}

		// Header 与 Token 中的租户不一致时以 Header 为准（显式切换）
		ctx := tenant.WithID(c.Request.Context(), tenantID)
		c.Request = c.Request.WithContext(ctx)
		c.Set("tenant_id", tenantID.String())
		c.Next()
	}
}
