package http

import (
	"crm-backend/internal/infrastructure/persistence"
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HealthHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := persistence.PingDB(db); err != nil {
			response.Error(c, 503, 503, "database unreachable")
			return
		}

		response.Success(c, gin.H{
			"status": "ok",
			"db":     "connected",
		})
	}
}
