package http

import (
	dealapp "crm-backend/internal/application/deal"
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *DealHandlers) StatsByStage(c *gin.Context) {
	tenantID, userID, ok := dealContext(c)
	if !ok {
		return
	}
	from, to, okRange := parseStatsDateRange(c)
	if !okRange {
		return
	}
	metric := c.DefaultQuery("metric", "count")
	data, err := h.svc.StatsByStage(c.Request.Context(), tenantID, userID, dealapp.StatsQuery{
		From: from, To: to, Metric: metric,
	})
	if err != nil {
		response.InternalError(c, "获取商机阶段统计失败")
		return
	}
	response.Success(c, data)
}

func (h *DealHandlers) StatsWinRate(c *gin.Context) {
	tenantID, userID, ok := dealContext(c)
	if !ok {
		return
	}
	from, to, okRange := parseStatsDateRange(c)
	if !okRange {
		return
	}
	data, err := h.svc.StatsWinRate(c.Request.Context(), tenantID, userID, dealapp.StatsQuery{
		From: from, To: to, Granularity: c.DefaultQuery("granularity", "week"),
	})
	if err != nil {
		response.InternalError(c, "获取赢单率统计失败")
		return
	}
	response.Success(c, data)
}
