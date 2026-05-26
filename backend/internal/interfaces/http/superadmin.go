package http

import (
	"errors"
	"strconv"
	"time"

	"crm-backend/internal/application/audit"
	"crm-backend/internal/application/superadmin"
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewSuperAdminHandlers(svc *superadmin.Service, rec *audit.Recorder) *SuperAdminHandlers {
	return &SuperAdminHandlers{svc: svc, audit: rec}
}

type SuperAdminHandlers struct {
	svc   *superadmin.Service
	audit *audit.Recorder
}

func (h *SuperAdminHandlers) Overview(c *gin.Context) {
	data, err := h.svc.Overview(c.Request.Context())
	if err != nil {
		response.InternalError(c, "获取概览失败")
		return
	}
	response.Success(c, data)
}

func (h *SuperAdminHandlers) TenantActivityTrend(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	data, err := h.svc.TenantActivityTrend(c.Request.Context(), days)
	if err != nil {
		response.InternalError(c, "获取租户趋势失败")
		return
	}
	response.Success(c, data)
}

func (h *SuperAdminHandlers) ListTenants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	var isActive *bool
	if v := c.Query("is_active"); v != "" {
		active := v == "true" || v == "1"
		isActive = &active
	}

	result, total, err := h.svc.ListTenants(c.Request.Context(), page, pageSize, search, isActive)
	if err != nil {
		response.InternalError(c, "获取租户列表失败")
		return
	}
	response.SuccessPage(c, result, response.Pagination{
		Page: page, PageSize: pageSize, Total: total,
	})
}

func (h *SuperAdminHandlers) GetTenant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "租户 ID 格式无效")
		return
	}

	tenant, err := h.svc.GetTenant(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, superadmin.ErrTenantNotFound) {
			response.NotFound(c, "租户不存在")
			return
		}
		response.InternalError(c, "获取租户详情失败")
		return
	}
	response.Success(c, tenant)
}

type patchSuperAdminTenantRequest struct {
	IsActive bool `json:"is_active"`
}

func (h *SuperAdminHandlers) PatchTenant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "租户 ID 格式无效")
		return
	}

	var req patchSuperAdminTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tenant, err := h.svc.SetTenantActive(c.Request.Context(), id, req.IsActive)
	if err != nil {
		if errors.Is(err, superadmin.ErrTenantNotFound) {
			response.NotFound(c, "租户不存在")
			return
		}
		response.InternalError(c, "更新租户状态失败")
		return
	}
	recordAudit(c, h.audit, id, "tenant.set_active", "tenant", &id, map[string]bool{"is_active": req.IsActive}, nil)
	response.Success(c, tenant)
}

func (h *SuperAdminHandlers) TenantHealth(c *gin.Context) {
	data, err := h.svc.TenantHealth(c.Request.Context())
	if err != nil {
		response.InternalError(c, "获取租户健康度失败")
		return
	}
	response.Success(c, data)
}

func (h *SuperAdminHandlers) PlanDistribution(c *gin.Context) {
	from, to := parseInsightDateRange(c)
	data, err := h.svc.PlanDistribution(c.Request.Context(), from, to)
	if err != nil {
		response.InternalError(c, "获取套餐分布失败")
		return
	}
	response.Success(c, data)
}

func (h *SuperAdminHandlers) TopTenants(c *gin.Context) {
	metric := c.DefaultQuery("metric", "activity")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := h.svc.TopTenants(c.Request.Context(), metric, limit)
	if err != nil {
		response.InternalError(c, "获取 TOP 租户失败")
		return
	}
	response.Success(c, data)
}

func parseInsightDateRange(c *gin.Context) (time.Time, time.Time) {
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -30)
	if s := c.Query("from"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			from = t
		}
	}
	if s := c.Query("to"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			to = t
		}
	}
	return from, to
}
