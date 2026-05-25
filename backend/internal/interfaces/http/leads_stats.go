package http

import (
	"context"
	"time"

	leadapp "crm-backend/internal/application/lead"
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *LeadHandlers) StatsBySource(c *gin.Context) {
	h.writeDistribution(c, h.svc.StatsBySource)
}

func (h *LeadHandlers) StatsByStatus(c *gin.Context) {
	h.writeDistribution(c, h.svc.StatsByStatus)
}

func (h *LeadHandlers) StatsTrend(c *gin.Context) {
	tenantID, userID, ok := leadContext(c)
	if !ok {
		return
	}
	q, ok := parseLeadStatsQuery(c)
	if !ok {
		return
	}
	data, err := h.svc.StatsTrend(c.Request.Context(), tenantID, userID, q)
	if err != nil {
		response.InternalError(c, "获取线索趋势失败")
		return
	}
	response.Success(c, data)
}

func (h *LeadHandlers) StatsFunnel(c *gin.Context) {
	tenantID, userID, ok := leadContext(c)
	if !ok {
		return
	}
	q, ok := parseLeadStatsQuery(c)
	if !ok {
		return
	}
	data, err := h.svc.StatsFunnel(c.Request.Context(), tenantID, userID, q)
	if err != nil {
		response.InternalError(c, "获取转化漏斗失败")
		return
	}
	response.Success(c, data)
}

type distributionFn func(context.Context, uuid.UUID, uuid.UUID, leadapp.StatsQuery) (*leadapp.DistributionDTO, error)

func (h *LeadHandlers) writeDistribution(c *gin.Context, fn distributionFn) {
	tenantID, userID, ok := leadContext(c)
	if !ok {
		return
	}
	q, ok := parseLeadStatsQuery(c)
	if !ok {
		return
	}
	data, err := fn(c.Request.Context(), tenantID, userID, q)
	if err != nil {
		response.InternalError(c, "获取线索统计失败")
		return
	}
	response.Success(c, data)
}

func parseLeadStatsQuery(c *gin.Context) (leadapp.StatsQuery, bool) {
	from, to, ok := parseStatsDateRange(c)
	if !ok {
		return leadapp.StatsQuery{}, false
	}
	return leadapp.StatsQuery{
		From:        from,
		To:          to,
		Granularity: c.DefaultQuery("granularity", "day"),
	}, true
}

func parseStatsDateRange(c *gin.Context) (from, to *time.Time, ok bool) {
	now := time.Now().UTC()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(24 * time.Hour)
	start := end.AddDate(0, 0, -30)

	if raw := c.Query("from"); raw != "" {
		t, err := time.Parse("2006-01-02", raw)
		if err != nil {
			response.BadRequest(c, "from 日期格式应为 YYYY-MM-DD")
			return nil, nil, false
		}
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	}
	if raw := c.Query("to"); raw != "" {
		t, err := time.Parse("2006-01-02", raw)
		if err != nil {
			response.BadRequest(c, "to 日期格式应为 YYYY-MM-DD")
			return nil, nil, false
		}
		end = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Add(24 * time.Hour)
	}
	if !end.After(start) {
		response.BadRequest(c, "to 必须晚于 from")
		return nil, nil, false
	}
	return &start, &end, true
}
