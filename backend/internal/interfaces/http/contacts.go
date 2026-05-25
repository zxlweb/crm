package http

import (
	"errors"

	contactapp "crm-backend/internal/application/contact"
	"crm-backend/internal/application/audit"
	"crm-backend/internal/application/emotion"
	"crm-backend/internal/application/insights"
	"crm-backend/internal/pkg/response"
	"crm-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewContactHandlers(svc *contactapp.Service, rec *audit.Recorder, emotionSvc *emotion.Service) *ContactHandlers {
	return &ContactHandlers{svc: svc, audit: rec, emotion: emotionSvc}
}

type ContactHandlers struct {
	svc     *contactapp.Service
	audit   *audit.Recorder
	emotion *emotion.Service
}

func (h *ContactHandlers) List(c *gin.Context) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	q, ok := parseContactListQuery(c)
	if !ok {
		return
	}
	result, err := h.svc.List(c.Request.Context(), tenantID, userID, q)
	if err != nil {
		response.InternalError(c, "获取联系人列表失败")
		return
	}
	response.SuccessPage(c, gin.H{"items": result.Items}, response.Pagination{
		Page:     result.Page,
		PageSize: result.Size,
		Total:    result.Total,
	})
}

func (h *ContactHandlers) ListByAccount(c *gin.Context) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	q, ok := parseContactListQuery(c)
	if !ok {
		return
	}
	result, err := h.svc.ListByAccount(c.Request.Context(), tenantID, userID, accountID, q)
	if err != nil {
		writeContactError(c, err)
		return
	}
	response.SuccessPage(c, gin.H{"items": result.Items}, response.Pagination{
		Page:     result.Page,
		PageSize: result.Size,
		Total:    result.Total,
	})
}

func (h *ContactHandlers) Create(c *gin.Context) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	var req contactCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	dto, err := h.svc.Create(c.Request.Context(), tenantID, userID, contactapp.CreateInput{
		AccountID:      req.AccountID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		IsPrimary:      req.IsPrimary,
		OwnerID:        req.OwnerID,
		LifecycleStage: req.LifecycleStage,
		Tags:           req.Tags,
	})
	if err != nil {
		writeContactError(c, err)
		return
	}
	recordAudit(c, h.audit, tenantID, "contact.create", "contact", &dto.ID, dto, nil)
	response.Created(c, dto)
}

func (h *ContactHandlers) Get(c *gin.Context) {
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
		writeContactError(c, err)
		return
	}
	response.Success(c, dto)
}

func (h *ContactHandlers) Put(c *gin.Context) {
	h.update(c, true)
}

func (h *ContactHandlers) Patch(c *gin.Context) {
	h.update(c, false)
}

func (h *ContactHandlers) update(c *gin.Context, full bool) {
	tenantID, userID, ok := accountContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	var req contactUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	in := contactapp.UpdateInput{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		IsPrimary:      req.IsPrimary,
		OwnerID:        req.OwnerID,
		LifecycleStage: req.LifecycleStage,
	}
	if req.AccountID != nil || req.ClearAccount {
		in.AccountIDSet = true
		if req.ClearAccount {
			in.AccountID = nil
		} else {
			in.AccountID = req.AccountID
		}
	}
	if req.Tags != nil {
		in.Tags = *req.Tags
		in.TagsSet = true
	}
	old, _ := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	dto, err := h.svc.Update(c.Request.Context(), tenantID, userID, id, in, full)
	if err != nil {
		writeContactError(c, err)
		return
	}
	recordAudit(c, h.audit, tenantID, "contact.update", "contact", &id, dto, old)
	response.Success(c, dto)
}

func (h *ContactHandlers) Delete(c *gin.Context) {
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
		writeContactError(c, err)
		return
	}
	recordAudit(c, h.audit, tenantID, "contact.delete", "contact", &id, nil, old)
	response.Success(c, gin.H{"deleted": true})
}

func (h *ContactHandlers) EmotionJourney(c *gin.Context) {
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
		writeContactError(c, err)
		return
	}
	rangeQ := c.DefaultQuery("range", "90d")
	if h.emotion == nil {
		response.Success(c, emotion.EmptyJourney("contact", id, dto.LifecycleStage))
		return
	}
	data, err := h.emotion.BuildJourney(c.Request.Context(), tenantID, emotion.SubjectInput{
		SubjectType:      "contact",
		SubjectID:        id,
		LifecycleCurrent: dto.LifecycleStage,
		CreatedAt:        dto.CreatedAt,
	}, rangeQ)
	if err != nil {
		if errors.Is(err, emotion.ErrInvalidRange) {
			response.BadRequest(c, emotion.FormatRangeError())
			return
		}
		response.InternalError(c, "获取情绪旅程失败")
		return
	}
	response.Success(c, data)
}

func (h *ContactHandlers) EvaluateInsights(c *gin.Context) {
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
		writeContactError(c, err)
		return
	}
	data := insights.EvaluateContact(id, dto.LifecycleStage, int16(dto.EngagementScore), dto.LastActivityAt)
	response.Success(c, data)
}

type contactCreateRequest struct {
	AccountID      *uuid.UUID `json:"account_id"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	IsPrimary      bool       `json:"is_primary"`
	OwnerID        *uuid.UUID `json:"owner_id"`
	LifecycleStage string     `json:"lifecycle_stage"`
	Tags           []string   `json:"tags"`
}

type contactUpdateRequest struct {
	AccountID      *uuid.UUID `json:"account_id"`
	ClearAccount   bool       `json:"clear_account"`
	FirstName      *string    `json:"first_name"`
	LastName       *string    `json:"last_name"`
	Email          *string    `json:"email"`
	Phone          *string    `json:"phone"`
	IsPrimary      *bool      `json:"is_primary"`
	OwnerID        *uuid.UUID `json:"owner_id"`
	LifecycleStage *string    `json:"lifecycle_stage"`
	Tags           *[]string  `json:"tags"`
}

func parseContactListQuery(c *gin.Context) (contactapp.ListQuery, bool) {
	q := contactapp.ListQuery{
		Page:               queryInt(c, "page", 1),
		PageSize:           queryInt(c, "page_size", 20),
		Search:             c.Query("search"),
		LifecycleStage:     c.Query("lifecycle_stage"),
		RelationshipHealth: c.Query("relationship_health"),
	}
	if aid := c.Query("account_id"); aid != "" {
		id, err := uuid.Parse(aid)
		if err != nil {
			response.BadRequest(c, "account_id 格式无效")
			return q, false
		}
		q.AccountID = &id
	}
	if oid := c.Query("owner_id"); oid != "" {
		id, err := uuid.Parse(oid)
		if err != nil {
			response.BadRequest(c, "owner_id 格式无效")
			return q, false
		}
		q.OwnerID = &id
	}
	return q, true
}

func writeContactError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, contactapp.ErrNotFound):
		response.NotFound(c, "联系人不存在")
	case errors.Is(err, repository.ErrAccountNotFound):
		response.NotFound(c, "关联公司不存在")
	case errors.Is(err, contactapp.ErrInvalidLifecycle):
		response.BadRequest(c, "invalid lifecycle_stage")
	case errors.Is(err, contactapp.ErrNameRequired):
		response.BadRequest(c, "name_or_email_required")
	default:
		response.InternalError(c, "处理联系人失败")
	}
}
