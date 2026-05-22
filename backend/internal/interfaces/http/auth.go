package http

import (
	"errors"

	"crm-backend/internal/application/auth"
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type switchTenantRequest struct {
	TenantID string `json:"tenant_id" binding:"required,uuid"`
}

func NewAuthHandlers(svc *auth.Service) *AuthHandlers {
	return &AuthHandlers{svc: svc}
}

type AuthHandlers struct {
	svc *auth.Service
}

func (h *AuthHandlers) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			response.Unauthorized(c, "邮箱或密码错误")
			return
		}
		response.InternalError(c, "登录失败")
		return
	}
	response.Success(c, result)
}

func (h *AuthHandlers) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			response.Unauthorized(c, "Refresh Token 无效或已过期")
			return
		}
		response.InternalError(c, "刷新 Token 失败")
		return
	}
	response.Success(c, result)
}

func (h *AuthHandlers) Profile(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		response.Unauthorized(c, "无效用户")
		return
	}

	profile, err := h.svc.Profile(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, "获取用户信息失败")
		return
	}
	response.Success(c, profile)
}

func (h *AuthHandlers) ListTenants(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		response.Unauthorized(c, "无效用户")
		return
	}

	tenants, err := h.svc.ListTenants(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, "获取租户列表失败")
		return
	}
	response.Success(c, tenants)
}

func (h *AuthHandlers) SwitchTenant(c *gin.Context) {
	var req switchTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		response.BadRequest(c, "tenant_id 格式无效")
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		response.Unauthorized(c, "无效用户")
		return
	}

	result, err := h.svc.SwitchTenant(
		c.Request.Context(),
		userID,
		c.GetString("email"),
		c.GetBool("is_super_admin"),
		tenantID,
	)
	if err != nil {
		if errors.Is(err, auth.ErrTenantForbidden) {
			response.Forbidden(c, "无权访问该租户")
			return
		}
		response.InternalError(c, "切换租户失败")
		return
	}
	response.Success(c, result)
}

func CurrentUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Success(c, gin.H{
			"user_id":        c.GetString("user_id"),
			"email":          c.GetString("email"),
			"is_super_admin": c.GetBool("is_super_admin"),
			"tenant_id":      c.GetString("tenant_id"),
		})
	}
}
