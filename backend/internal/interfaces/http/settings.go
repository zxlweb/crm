package http

import (
	"encoding/json"
	"errors"

	settingsapp "crm-backend/internal/application/settings"
	"crm-backend/internal/application/audit"
	customapp "crm-backend/internal/application/customfield"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewSettingsHandlers(svc *settingsapp.Service, rec *audit.Recorder) *SettingsHandlers {
	return &SettingsHandlers{svc: svc, audit: rec}
}

type SettingsHandlers struct {
	svc   *settingsapp.Service
	audit *audit.Recorder
}

// GetTenant GET /api/settings/tenant
func (h *SettingsHandlers) GetTenant(c *gin.Context) {
	tenantID, ok := settingsTenantID(c)
	if !ok {
		return
	}
	dto, err := h.svc.GetTenant(c.Request.Context(), tenantID)
	if err != nil {
		if errors.Is(err, settingsapp.ErrTenantNotFound) {
			response.NotFound(c, "租户不存在")
			return
		}
		response.InternalError(c, "获取租户配置失败")
		return
	}
	response.Success(c, dto)
}

type patchTenantRequest struct {
	TenantName       *string                     `json:"tenant_name"`
	DefaultLocale    *string                     `json:"default_locale"`
	Timezone         *string                     `json:"timezone"`
	BusinessSwitches *patchBusinessSwitchesJSON  `json:"business_switches"`
	SalesQuota       *crm.SalesQuota             `json:"sales_quota"`
}

type patchBusinessSwitchesJSON struct {
	AIPreviewEnabled *bool   `json:"ai_preview_enabled"`
	LeadImportMode   *string `json:"lead_import_mode"`
}

// PatchTenant PATCH /api/settings/tenant
func (h *SettingsHandlers) PatchTenant(c *gin.Context) {
	tenantID, ok := settingsTenantID(c)
	if !ok {
		return
	}
	var req patchTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := uuid.Parse(c.GetString("user_id"))
	var bs *crm.BusinessSwitchesPatch
	if req.BusinessSwitches != nil {
		bs = &crm.BusinessSwitchesPatch{
			AIPreviewEnabled: req.BusinessSwitches.AIPreviewEnabled,
			LeadImportMode:   req.BusinessSwitches.LeadImportMode,
		}
	}
	dto, err := h.svc.PatchTenant(c.Request.Context(), tenantID, settingsapp.PatchInput{
		TenantName:       req.TenantName,
		DefaultLocale:    req.DefaultLocale,
		Timezone:         req.Timezone,
		BusinessSwitches: bs,
		SalesQuota:       req.SalesQuota,
		UpdatedBy:        userID,
	})
	if err != nil {
		if errors.Is(err, crm.ErrSettingKeyInvalid) {
			response.BadRequest(c, "setting_key_invalid")
			return
		}
		if errors.Is(err, settingsapp.ErrTenantNotFound) {
			response.NotFound(c, "租户不存在")
			return
		}
		response.InternalError(c, "更新租户配置失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "settings.update", "settings", &tenantID, dto, nil)
	response.Success(c, dto)
}

// ListFeatures GET /api/settings/features
func (h *SettingsHandlers) ListFeatures(c *gin.Context) {
	tenantID, ok := settingsTenantID(c)
	if !ok {
		return
	}
	items, err := h.svc.ListFeatures(c.Request.Context(), tenantID)
	if err != nil {
		if errors.Is(err, settingsapp.ErrTenantNotFound) {
			response.NotFound(c, "租户不存在")
			return
		}
		response.InternalError(c, "获取功能开关失败")
		return
	}
	response.Success(c, gin.H{"items": items})
}

func settingsTenantID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.GetString("tenant_id"))
	if err != nil {
		response.BadRequest(c, "租户上下文无效")
		return uuid.Nil, false
	}
	return id, true
}

func NewCustomFieldHandlers(svc *customapp.Service, rec *audit.Recorder) *CustomFieldHandlers {
	return &CustomFieldHandlers{svc: svc, audit: rec}
}

type CustomFieldHandlers struct {
	svc   *customapp.Service
	audit *audit.Recorder
}

func (h *CustomFieldHandlers) List(c *gin.Context) {
	tenantID, ok := settingsTenantID(c)
	if !ok {
		return
	}
	items, err := h.svc.List(c.Request.Context(), tenantID, c.Query("entity_type"))
	if err != nil {
		response.InternalError(c, "获取自定义字段失败")
		return
	}
	response.Success(c, gin.H{"items": items})
}

type createCustomFieldRequest struct {
	EntityType   string          `json:"entity_type" binding:"required"`
	FieldKey     string          `json:"field_key" binding:"required"`
	FieldLabel   json.RawMessage `json:"field_label" binding:"required"`
	FieldType    string          `json:"field_type" binding:"required"`
	Required     bool            `json:"required"`
	Options      json.RawMessage `json:"options"`
	DefaultValue json.RawMessage `json:"default_value"`
	DisplayOrder int             `json:"display_order"`
}

func (h *CustomFieldHandlers) Create(c *gin.Context) {
	tenantID, ok := settingsTenantID(c)
	if !ok {
		return
	}
	var req createCustomFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	dto, err := h.svc.Create(c.Request.Context(), tenantID, customapp.CreateInput{
		EntityType: req.EntityType, FieldKey: req.FieldKey, FieldLabel: req.FieldLabel,
		FieldType: req.FieldType, Required: req.Required, Options: req.Options,
		DefaultValue: req.DefaultValue, DisplayOrder: req.DisplayOrder,
	})
	if err != nil {
		mapCustomFieldErr(c, err)
		return
	}
	recordAudit(c, h.audit, tenantID, "custom_field.create", "custom_field", &dto.ID, dto, nil)
	response.Created(c, dto)
}

func (h *CustomFieldHandlers) Patch(c *gin.Context) {
	tenantID, ok := settingsTenantID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "字段 ID 格式无效")
		return
	}
	var req struct {
		FieldLabel   *json.RawMessage `json:"field_label"`
		FieldType    *string          `json:"field_type"`
		Required     *bool            `json:"required"`
		Options      json.RawMessage  `json:"options"`
		DefaultValue json.RawMessage  `json:"default_value"`
		DisplayOrder *int             `json:"display_order"`
		IsActive     *bool            `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	dto, err := h.svc.Update(c.Request.Context(), tenantID, id, customapp.UpdateInput{
		FieldLabel: req.FieldLabel, FieldType: req.FieldType, Required: req.Required,
		Options: req.Options, DefaultValue: req.DefaultValue, DisplayOrder: req.DisplayOrder,
		IsActive: req.IsActive,
	})
	if err != nil {
		mapCustomFieldErr(c, err)
		return
	}
	recordAudit(c, h.audit, tenantID, "custom_field.update", "custom_field", &id, dto, nil)
	response.Success(c, dto)
}

func (h *CustomFieldHandlers) Delete(c *gin.Context) {
	tenantID, ok := settingsTenantID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "字段 ID 格式无效")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), tenantID, id); err != nil {
		mapCustomFieldErr(c, err)
		return
	}
	recordAudit(c, h.audit, tenantID, "custom_field.delete", "custom_field", &id, nil, nil)
	response.Success(c, gin.H{"id": id.String()})
}

func mapCustomFieldErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, customapp.ErrNotFound):
		response.NotFound(c, "自定义字段不存在")
	case errors.Is(err, customapp.ErrKeyConflict):
		response.Conflict(c, "custom_field_key_conflict")
	case errors.Is(err, customapp.ErrTypeInvalid):
		response.BadRequest(c, "custom_field_type_invalid")
	default:
		response.InternalError(c, "自定义字段操作失败")
	}
}
