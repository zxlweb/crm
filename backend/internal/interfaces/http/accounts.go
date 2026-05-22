package http

import (
	"errors"
	"strconv"
	"time"

	accountapp "crm-backend/internal/application/account"
	"crm-backend/internal/application/audit"
	"crm-backend/internal/application/emotion"
	"crm-backend/internal/application/insights"
	"crm-backend/internal/pkg/response"
	"crm-backend/internal/pkg/tenant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewAccountHandlers(svc *accountapp.Service, rec *audit.Recorder) *AccountHandlers {
	return &AccountHandlers{svc: svc, audit: rec}
}

type AccountHandlers struct {
	svc   *accountapp.Service
	audit *audit.Recorder
}

func (h *AccountHandlers) List(c *gin.Context) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	q := accountapp.ListQuery{
		Page:               queryInt(c, "page", 1),
		PageSize:           queryInt(c, "page_size", 20),
		Search:             c.Query("search"),
		LifecycleStage:     c.Query("lifecycle_stage"),
		RelationshipHealth: c.Query("relationship_health"),
	}
	if oid := c.Query("owner_id"); oid != "" {
		id, err := uuid.Parse(oid)
		if err != nil {
			response.BadRequest(c, "owner_id 格式无效")
			return
		}
		q.OwnerID = &id
	}
	result, err := h.svc.List(c.Request.Context(), tenantID, userID, q)
	if err != nil {
		response.InternalError(c, "获取公司列表失败")
		return
	}
	response.SuccessPage(c, gin.H{"items": result.Items}, response.Pagination{
		Page:     result.Page,
		PageSize: result.Size,
		Total:    result.Total,
	})
}

func (h *AccountHandlers) Create(c *gin.Context) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	var req accountCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	dto, err := h.svc.Create(c.Request.Context(), tenantID, userID, accountapp.CreateInput{
		Name:           req.Name,
		Industry:       req.Industry,
		Website:        req.Website,
		OwnerID:        req.OwnerID,
		LifecycleStage: req.LifecycleStage,
		Tags:           req.Tags,
	})
	if err != nil {
		if errors.Is(err, accountapp.ErrInvalidLifecycle) {
			response.BadRequest(c, "invalid lifecycle_stage")
			return
		}
		response.InternalError(c, "创建公司失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "account.create", "account", &dto.ID, dto, nil)
	response.Created(c, dto)
}

func (h *AccountHandlers) Get(c *gin.Context) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	dto, err := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	if err != nil {
		if errors.Is(err, accountapp.ErrNotFound) {
			response.NotFound(c, "公司不存在")
			return
		}
		response.InternalError(c, "获取公司详情失败")
		return
	}
	response.Success(c, dto)
}

func (h *AccountHandlers) Put(c *gin.Context) {
	h.update(c, true)
}

func (h *AccountHandlers) Patch(c *gin.Context) {
	h.update(c, false)
}

func (h *AccountHandlers) update(c *gin.Context, full bool) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	var req accountUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	in := accountapp.UpdateInput{
		Name:           req.Name,
		Industry:       req.Industry,
		Website:        req.Website,
		OwnerID:        req.OwnerID,
		LifecycleStage: req.LifecycleStage,
	}
	if req.Tags != nil {
		in.Tags = *req.Tags
		in.TagsSet = true
	}
	old, _ := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	dto, err := h.svc.Update(c.Request.Context(), tenantID, userID, id, in, full)
	if err != nil {
		if errors.Is(err, accountapp.ErrNotFound) {
			response.NotFound(c, "公司不存在")
			return
		}
		if errors.Is(err, accountapp.ErrInvalidLifecycle) {
			response.BadRequest(c, "invalid lifecycle_stage")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	action := "account.update"
	if old != nil && dto.LifecycleStage != old.LifecycleStage {
		action = "lifecycle.change"
	}
	recordAudit(c, h.audit, tenantID, action, "account", &id, dto, old)
	response.Success(c, dto)
}

func (h *AccountHandlers) Delete(c *gin.Context) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	old, _ := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	if err := h.svc.Delete(c.Request.Context(), tenantID, userID, id); err != nil {
		if errors.Is(err, accountapp.ErrNotFound) {
			response.NotFound(c, "公司不存在")
			return
		}
		response.InternalError(c, "删除公司失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "account.delete", "account", &id, nil, old)
	response.Success(c, gin.H{"deleted": true})
}

func (h *AccountHandlers) EmotionJourney(c *gin.Context) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	dto, err := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	if err != nil {
		if errors.Is(err, accountapp.ErrNotFound) {
			response.NotFound(c, "公司不存在")
			return
		}
		response.InternalError(c, "获取情绪旅程失败")
		return
	}
	_ = c.Query("range")
	data := emotion.EmptyJourney("account", id, dto.LifecycleStage)
	response.Success(c, data)
}

func (h *AccountHandlers) EvaluateInsights(c *gin.Context) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	dto, err := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	if err != nil {
		if errors.Is(err, accountapp.ErrNotFound) {
			response.NotFound(c, "公司不存在")
			return
		}
		response.InternalError(c, "洞察求值失败")
		return
	}
	var lastAt *time.Time
	if dto.LastActivityAt != nil {
		lastAt = dto.LastActivityAt
	}
	data := insights.EvaluateAccount(id, dto.LifecycleStage, int16(dto.EngagementScore), lastAt)
	response.Success(c, data)
}

type accountCreateRequest struct {
	Name           string     `json:"name" binding:"required,min=1,max=255"`
	Industry       string     `json:"industry"`
	Website        string     `json:"website"`
	OwnerID        *uuid.UUID `json:"owner_id"`
	LifecycleStage string     `json:"lifecycle_stage"`
	Tags           []string   `json:"tags"`
}

type accountUpdateRequest struct {
	Name           *string    `json:"name"`
	Industry       *string    `json:"industry"`
	Website        *string    `json:"website"`
	OwnerID        *uuid.UUID `json:"owner_id"`
	LifecycleStage *string    `json:"lifecycle_stage"`
	Tags           *[]string  `json:"tags"`
}

func accountContext(c *gin.Context) (tenantID, userID uuid.UUID, ok bool) {
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

func queryInt(c *gin.Context, key string, def int) int {
	v := c.Query(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
