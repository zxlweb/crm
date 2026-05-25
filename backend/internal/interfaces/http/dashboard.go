package http

import (
	"errors"

	dashapp "crm-backend/internal/application/dashboard"
	"crm-backend/internal/pkg/datascope"
	"crm-backend/internal/pkg/response"
	"crm-backend/internal/pkg/tenant"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewDashboardHandlers(svc *dashapp.Service, enforcer *casbin.Enforcer) *DashboardHandlers {
	return &DashboardHandlers{svc: svc, enforcer: enforcer}
}

type DashboardHandlers struct {
	svc      *dashapp.Service
	enforcer *casbin.Enforcer
}

func (h *DashboardHandlers) Summary(c *gin.Context) {
	if !h.allowDashboard(c) {
		return
	}
	tenantID, userID, ok := dashboardContext(c)
	if !ok {
		return
	}
	preview := c.Query("preview") == "1"
	data, err := h.svc.Summary(c.Request.Context(), tenantID, userID, preview)
	if err != nil {
		response.InternalError(c, "获取仪表盘汇总失败")
		return
	}
	response.Success(c, data)
}

func (h *DashboardHandlers) Funnel(c *gin.Context) {
	if !h.allowDashboard(c) {
		return
	}
	tenantID, userID, ok := dashboardContext(c)
	if !ok {
		return
	}
	scope := c.DefaultQuery("scope", "deals")
	data, err := h.svc.Funnel(c.Request.Context(), tenantID, userID, scope)
	if err != nil {
		response.InternalError(c, "获取漏斗数据失败")
		return
	}
	response.Success(c, data)
}

func (h *DashboardHandlers) Quota(c *gin.Context) {
	if !h.allowDashboard(c) {
		return
	}
	tenantID, userID, ok := dashboardContext(c)
	if !ok {
		return
	}
	data, err := h.svc.Quota(c.Request.Context(), tenantID, userID)
	if err != nil {
		response.InternalError(c, "获取配额数据失败")
		return
	}
	response.Success(c, data)
}

func (h *DashboardHandlers) TeamRanking(c *gin.Context) {
	if !h.allowDashboard(c) {
		return
	}
	tenantID, userID, ok := dashboardContext(c)
	if !ok {
		return
	}
	limit := queryInt(c, "limit", 10)
	data, err := h.svc.TeamRanking(c.Request.Context(), tenantID, userID, c.DefaultQuery("metric", "won_amount"), limit)
	if err != nil {
		if errors.Is(err, dashapp.ErrTeamRankingDenied) {
			response.Forbidden(c, "dashboard_scope_denied")
			return
		}
		response.InternalError(c, "获取经理排行失败")
		return
	}
	response.Success(c, data)
}

func (h *DashboardHandlers) Todo(c *gin.Context) {
	if !h.allowDashboard(c) {
		return
	}
	tenantID, userID, ok := dashboardContext(c)
	if !ok {
		return
	}
	date := c.DefaultQuery("date", "")
	data, err := h.svc.Todo(c.Request.Context(), tenantID, userID, date)
	if err != nil {
		response.InternalError(c, "获取待办失败")
		return
	}
	response.Success(c, data)
}

func (h *DashboardHandlers) allowDashboard(c *gin.Context) bool {
	if c.GetBool("is_super_admin") {
		return true
	}
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")
	if datascope.CanAccessDashboard(h.enforcer, userID, tenantID) {
		return true
	}
	response.Forbidden(c, "无权限访问仪表盘")
	return false
}

func dashboardContext(c *gin.Context) (tenantID, userID uuid.UUID, ok bool) {
	tid, okT := tenant.IDFromContext(c.Request.Context())
	if !okT {
		response.BadRequest(c, "缺少租户上下文")
		return uuid.Nil, uuid.Nil, false
	}
	uid, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		response.Unauthorized(c, "用户未认证")
		return uuid.Nil, uuid.Nil, false
	}
	return tid, uid, true
}
