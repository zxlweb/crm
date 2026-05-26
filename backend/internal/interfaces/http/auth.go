package http

import (
	"errors"

	"crm-backend/internal/application/auth"
	"crm-backend/internal/pkg/httputil"
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

type switchRoleRequest struct {
	RoleID string `json:"role_id" binding:"required,uuid"`
}

type registerRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
	CompanyName string `json:"company_name" binding:"required,min=2,max=255"`
	Domain      string `json:"domain" binding:"omitempty,min=2,max=50"`
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

	result, err := h.svc.Login(c.Request.Context(), req.Email, req.Password, httputil.ClientIP(c))
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

func (h *AuthHandlers) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.Register(c.Request.Context(), auth.RegisterInput{
		Email:       req.Email,
		Password:    req.Password,
		Name:        req.Name,
		CompanyName: req.CompanyName,
		Domain:      req.Domain,
	}, httputil.ClientIP(c))
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrEmailExists):
			response.Error(c, 409, 409, "该邮箱已注册")
		case errors.Is(err, auth.ErrDomainExists):
			response.Error(c, 409, 409, "该域名标识已被占用")
		case errors.Is(err, auth.ErrInvalidDomain):
			response.BadRequest(c, "域名标识格式无效（小写字母、数字、连字符）")
		default:
			response.InternalError(c, "注册失败")
		}
		return
	}
	response.Created(c, result)
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
		httputil.ClientIP(c),
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

func (h *AuthHandlers) SwitchRole(c *gin.Context) {
	var req switchRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		response.BadRequest(c, "role_id 格式无效")
		return
	}

	tenantID, err := uuid.Parse(c.GetString("tenant_id"))
	if err != nil {
		response.BadRequest(c, "缺少租户上下文")
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		response.Unauthorized(c, "无效用户")
		return
	}

	result, err := h.svc.SwitchRole(
		c.Request.Context(),
		userID,
		c.GetString("email"),
		c.GetBool("is_super_admin"),
		tenantID,
		roleID,
		httputil.ClientIP(c),
	)
	if err != nil {
		if errors.Is(err, auth.ErrRoleForbidden) {
			response.Forbidden(c, "无权使用该角色")
			return
		}
		response.InternalError(c, "切换角色失败")
		return
	}
	response.Success(c, result)
}

func CurrentUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := gin.H{
			"user_id":        c.GetString("user_id"),
			"email":          c.GetString("email"),
			"is_super_admin": c.GetBool("is_super_admin"),
			"tenant_id":      c.GetString("tenant_id"),
		}
		if rid := c.GetString("active_role_id"); rid != "" {
			payload["active_role_id"] = rid
		}
		response.Success(c, payload)
	}
}
