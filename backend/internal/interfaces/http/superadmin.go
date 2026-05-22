package http

import (
	"errors"
	"strconv"

	"crm-backend/internal/application/superadmin"
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewSuperAdminHandlers(svc *superadmin.Service) *SuperAdminHandlers {
	return &SuperAdminHandlers{svc: svc}
}

type SuperAdminHandlers struct {
	svc *superadmin.Service
}

func (h *SuperAdminHandlers) Overview(c *gin.Context) {
	data, err := h.svc.Overview(c.Request.Context())
	if err != nil {
		response.InternalError(c, "获取概览失败")
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

type patchTenantRequest struct {
	IsActive bool `json:"is_active"`
}

func (h *SuperAdminHandlers) PatchTenant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "租户 ID 格式无效")
		return
	}

	var req patchTenantRequest
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
	response.Success(c, tenant)
}
