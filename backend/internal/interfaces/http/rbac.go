package http

import (
	"errors"

	"crm-backend/internal/application/audit"
	rbacapp "crm-backend/internal/application/rbac"
	"crm-backend/internal/pkg/response"
	"crm-backend/internal/pkg/tenant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewRBACHandlers(svc *rbacapp.Service, rec *audit.Recorder) *RBACHandlers {
	return &RBACHandlers{svc: svc, audit: rec}
}

type RBACHandlers struct {
	svc   *rbacapp.Service
	audit *audit.Recorder
}

func (h *RBACHandlers) ListPermissions(c *gin.Context) {
	data, err := h.svc.ListPermissionDictionary(c.Request.Context())
	if err != nil {
		response.InternalError(c, "读取权限字典失败")
		return
	}
	response.Success(c, data)
}

func (h *RBACHandlers) ListPermissionItems(c *gin.Context) {
	data, err := h.svc.ListPermissionItems(c.Request.Context())
	if err != nil {
		response.InternalError(c, "读取权限列表失败")
		return
	}
	response.Success(c, data)
}

func (h *RBACHandlers) MyPermissions(c *gin.Context) {
	tenantID, userID, ok := rbacContext(c)
	if !ok {
		return
	}
	if c.GetBool("is_super_admin") {
		data, err := h.svc.ListPermissionDictionary(c.Request.Context())
		if err != nil {
			response.InternalError(c, "读取权限失败")
			return
		}
		response.Success(c, data)
		return
	}
	data, err := h.svc.MyPermissions(c.Request.Context(), tenantID, userID, activeRoleIDFromGin(c))
	if err != nil {
		if errors.Is(err, rbacapp.ErrRoleForbidden) {
			response.Forbidden(c, "当前活跃角色无效")
			return
		}
		response.InternalError(c, "读取权限失败")
		return
	}
	response.Success(c, data)
}

func (h *RBACHandlers) MyRoles(c *gin.Context) {
	tenantID, userID, ok := rbacContext(c)
	if !ok {
		return
	}
	data, err := h.svc.ListUserRoles(c.Request.Context(), tenantID, userID)
	if err != nil {
		response.InternalError(c, "获取角色列表失败")
		return
	}
	response.Success(c, data)
}

func (h *RBACHandlers) ListRoles(c *gin.Context) {
	tenantID, _, ok := rbacContext(c)
	if !ok {
		return
	}
	data, err := h.svc.ListRoles(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalError(c, "获取角色列表失败")
		return
	}
	response.Success(c, data)
}

type createRoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

func (h *RBACHandlers) CreateRole(c *gin.Context) {
	tenantID, _, ok := rbacContext(c)
	if !ok {
		return
	}
	var req createRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	role, err := h.svc.CreateRole(c.Request.Context(), tenantID, req.Name, req.Description)
	if err != nil {
		response.InternalError(c, "创建角色失败")
		return
	}
	if rid, err := uuid.Parse(role.ID); err == nil {
		recordAudit(c, h.audit, tenantID, "rbac.role.create", "role", &rid, role, nil)
	}
	response.Created(c, role)
}

func (h *RBACHandlers) UpdateRole(c *gin.Context) {
	tenantID, _, ok := rbacContext(c)
	if !ok {
		return
	}
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "角色 ID 格式无效")
		return
	}
	var req createRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	role, err := h.svc.UpdateRole(c.Request.Context(), tenantID, roleID, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, rbacapp.ErrRoleNotFound) {
			response.NotFound(c, "角色不存在")
			return
		}
		response.InternalError(c, "更新角色失败")
		return
	}
	response.Success(c, role)
}

type assignPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}

func (h *RBACHandlers) AssignRolePermissions(c *gin.Context) {
	tenantID, _, ok := rbacContext(c)
	if !ok {
		return
	}
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "角色 ID 格式无效")
		return
	}
	var req assignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	permIDs, err := parseUUIDs(req.PermissionIDs)
	if err != nil {
		response.BadRequest(c, "permission_ids 格式无效")
		return
	}
	role, err := h.svc.SetRolePermissions(c.Request.Context(), tenantID, roleID, permIDs)
	if err != nil {
		if errors.Is(err, rbacapp.ErrRoleNotFound) {
			response.NotFound(c, "角色不存在")
			return
		}
		if errors.Is(err, rbacapp.ErrInvalidPermissions) {
			response.BadRequest(c, "存在无效的权限 ID")
			return
		}
		response.InternalError(c, "分配权限失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "rbac.role.set_permissions", "role", &roleID, role, nil)
	response.Success(c, role)
}

func (h *RBACHandlers) ListUserRoles(c *gin.Context) {
	tenantID, _, ok := rbacContext(c)
	if !ok {
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "用户 ID 格式无效")
		return
	}
	data, err := h.svc.ListUserRoles(c.Request.Context(), tenantID, userID)
	if err != nil {
		response.InternalError(c, "获取用户角色失败")
		return
	}
	response.Success(c, data)
}

type assignUserRolesRequest struct {
	RoleIDs []string `json:"role_ids" binding:"required"`
}

func (h *RBACHandlers) AssignUserRoles(c *gin.Context) {
	tenantID, _, ok := rbacContext(c)
	if !ok {
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "用户 ID 格式无效")
		return
	}
	var req assignUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	roleIDs, err := parseUUIDs(req.RoleIDs)
	if err != nil {
		response.BadRequest(c, "role_ids 格式无效")
		return
	}
	data, err := h.svc.SetUserRoles(c.Request.Context(), tenantID, userID, roleIDs)
	if err != nil {
		if errors.Is(err, rbacapp.ErrInvalidRoles) {
			response.BadRequest(c, "存在无效的角色 ID")
			return
		}
		response.InternalError(c, "分配角色失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "rbac.user.set_roles", "user", &userID, data, nil)
	response.Success(c, data)
}

func (h *RBACHandlers) Check(c *gin.Context) {
	tenantID, userID, ok := rbacContext(c)
	if !ok {
		return
	}
	var req rbacapp.CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	allowed, err := h.svc.Check(
		c.Request.Context(),
		tenantID,
		userID.String(),
		c.GetBool("is_super_admin"),
		c.GetString("active_role_id"),
		req.Resource,
		req.Action,
	)
	if err != nil {
		response.InternalError(c, "权限检查失败")
		return
	}
	response.Success(c, rbacapp.CheckResult{Allowed: allowed})
}

func activeRoleIDFromGin(c *gin.Context) *uuid.UUID {
	s := c.GetString("active_role_id")
	if s == "" {
		return nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &id
}

func rbacContext(c *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := tenant.IDFromContext(c.Request.Context())
	if !ok {
		response.BadRequest(c, "缺少租户上下文")
		return uuid.Nil, uuid.Nil, false
	}
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		response.Unauthorized(c, "无效用户")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, userID, true
}

func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(ids))
	for _, s := range ids {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}
