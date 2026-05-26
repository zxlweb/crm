package http

import (
	"errors"
	"strconv"
	"time"

	auditstatsapp "crm-backend/internal/application/auditstats"
	"crm-backend/internal/application/audit"
	"crm-backend/internal/repository"
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewAuditStatsHandlers(svc *auditstatsapp.Service, rec *audit.Recorder) *AuditStatsHandlers {
	return &AuditStatsHandlers{svc: svc, audit: rec}
}

type AuditStatsHandlers struct {
	svc   *auditstatsapp.Service
	audit *audit.Recorder
}

func (h *AuditStatsHandlers) ByAction(c *gin.Context) {
	tenantID, f, ok := h.parseFilter(c)
	if !ok {
		return
	}
	data, err := h.svc.ByAction(c.Request.Context(), tenantID, f)
	if err != nil {
		response.InternalError(c, "获取审计统计失败")
		return
	}
	response.Success(c, data)
}

func (h *AuditStatsHandlers) Trend(c *gin.Context) {
	tenantID, f, ok := h.parseFilter(c)
	if !ok {
		return
	}
	granularity := c.DefaultQuery("granularity", "day")
	data, err := h.svc.Trend(c.Request.Context(), tenantID, f, granularity)
	if err != nil {
		response.InternalError(c, "获取审计趋势失败")
		return
	}
	response.Success(c, data)
}

func (h *AuditStatsHandlers) TopActors(c *gin.Context) {
	tenantID, f, ok := h.parseFilter(c)
	if !ok {
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	f.Limit = limit
	data, err := h.svc.TopActors(c.Request.Context(), tenantID, f)
	if err != nil {
		response.InternalError(c, "获取操作人统计失败")
		return
	}
	response.Success(c, data)
}

func (h *AuditStatsHandlers) Export(c *gin.Context) {
	tenantID, f, ok := h.parseFilter(c)
	if !ok {
		return
	}
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		response.Unauthorized(c, "用户无效")
		return
	}
	if c.Query("format") != "" && c.Query("format") != "csv" {
		response.BadRequest(c, "仅支�?format=csv")
		return
	}
	csv, err := h.svc.ExportCSV(c.Request.Context(), tenantID, userID, f)
	if err != nil {
		if errors.Is(err, auditstatsapp.ErrExportRateLimited) {
			response.TooManyRequests(c, "audit_export_rate_limited")
			return
		}
		response.InternalError(c, "导出审计日志失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "audit.export", "audit", nil, gin.H{"format": "csv"}, nil)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=audit-export.csv")
	c.String(200, csv)
}

func (h *AuditStatsHandlers) parseFilter(c *gin.Context) (uuid.UUID, repository.AuditStatsFilter, bool) {
	tenantID, ok := settingsTenantID(c)
	if !ok {
		return uuid.Nil, repository.AuditStatsFilter{}, false
	}
	from, to := parseTimeRange(c)
	f := repository.AuditStatsFilter{
		From:      from,
		To:        to,
		Module:    c.Query("module"),
		ActorRole: c.Query("actor_role"),
		Action:    c.Query("action"),
	}
	if aid := c.Query("actor_id"); aid != "" {
		id, err := uuid.Parse(aid)
		if err != nil {
			response.BadRequest(c, "actor_id 格式无效")
			return uuid.Nil, f, false
		}
		f.ActorID = &id
	}
	return tenantID, f, true
}

func parseTimeRange(c *gin.Context) (time.Time, time.Time) {
	var from, to time.Time
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			from = t
		} else if t, err := time.Parse("2006-01-02", v); err == nil {
			from = t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			to = t
		} else if t, err := time.Parse("2006-01-02", v); err == nil {
			to = t.Add(24*time.Hour - time.Nanosecond)
		}
	}
	if from.IsZero() && to.IsZero() {
		to = time.Now()
		from = to.AddDate(0, 0, -7)
	}
	return from, to
}
