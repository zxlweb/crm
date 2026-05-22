package httputil

import "github.com/gin-gonic/gin"

// ClientIP 从代理头或连接地址取客户端 IP
func ClientIP(c *gin.Context) string {
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	return c.ClientIP()
}
