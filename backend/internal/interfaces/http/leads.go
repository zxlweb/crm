package http

import (
	"errors"

	leadapp "crm-backend/internal/application/lead"
	"crm-backend/internal/application/audit"
	"crm-backend/internal/application/emotion"
	"crm-backend/internal/pkg/response"
	"crm-backend/internal/pkg/tenant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewLeadHandlers(svc *leadapp.Service, rec *audit.Recorder) *LeadHandlers {
	return &LeadHandlers{svc: svc, audit: rec}
}

type LeadHandlers struct {
	svc   *leadapp.Service
	audit *audit.Recorder
}

func (h *LeadHandlers) List(c *gin.Context) {
	tenantID, userID, ok := leadContext(c)
	if !ok {
		return
	}
	q := leadapp.ListQuery{
		Page:               queryInt(c, "page", 1),
		PageSize:           queryInt(c, "page_size", 20),
		Search:             c.Query("search"),
		Status:             c.Query("status"),
		Source:             c.Query("source"),
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
		response.InternalError(c, "获取线索列表失败")
		return
	}
	response.SuccessPage(c, gin.H{"items": result.Items}, response.Pagination{
		Page:     result.Page,
		PageSize: result.Size,
		Total:    result.Total,
	})
}

func (h *LeadHandlers) Create(c *gin.Context) {
	tenantID, userID, ok := leadContext(c)
	if !ok {
		return
	}
	var req leadCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	dto, err := h.svc.Create(c.Request.Context(), tenantID, userID, leadapp.CreateInput{
		Title:             req.Title,
		Status:            req.Status,
		Source:            req.Source,
		Amount:            req.Amount,
		ExpectedCloseDate: req.ExpectedCloseDate,
		OwnerID:           req.OwnerID,
		LifecycleStage:    req.LifecycleStage,
		Tags:              req.Tags,
	})
	if err != nil {
		if errors.Is(err, leadapp.ErrInvalidLifecycle) || errors.Is(err, leadapp.ErrInvalidStatus) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "创建线索失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "lead.create", "lead", &dto.ID, dto, nil)
	response.Created(c, dto)
}

func (h *LeadHandlers) Get(c *gin.Context) {
	tenantID, userID, ok := leadContext(c)
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
		if errors.Is(err, leadapp.ErrNotFound) {
			response.NotFound(c, "线索不存在")
			return
		}
		response.InternalError(c, "获取线索详情失败")
		return
	}
	response.Success(c, dto)
}

func (h *LeadHandlers) Put(c *gin.Context) {
	h.update(c, true)
}

func (h *LeadHandlers) Patch(c *gin.Context) {
	h.update(c, false)
}

func (h *LeadHandlers) update(c *gin.Context, full bool) {
	tenantID, userID, ok := leadContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	var req leadUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	in := leadapp.UpdateInput{
		Title:             req.Title,
		Status:            req.Status,
		Source:            req.Source,
		Amount:            req.Amount,
		ExpectedCloseDate: req.ExpectedCloseDate,
		OwnerID:           req.OwnerID,
		LifecycleStage:    req.LifecycleStage,
	}
	if req.Tags != nil {
		in.Tags = *req.Tags
		in.TagsSet = true
	}
	old, _ := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	dto, err := h.svc.Update(c.Request.Context(), tenantID, userID, id, in, full)
	if err != nil {
		if errors.Is(err, leadapp.ErrNotFound) {
			response.NotFound(c, "线索不存在")
			return
		}
		if errors.Is(err, leadapp.ErrInvalidStatusTransition) {
			response.BadRequest(c, "invalid_status_transition")
			return
		}
		if errors.Is(err, leadapp.ErrInvalidLifecycle) || errors.Is(err, leadapp.ErrInvalidStatus) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "更新线索失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "lead.update", "lead", &id, dto, old)
	response.Success(c, dto)
}

func (h *LeadHandlers) Delete(c *gin.Context) {
	tenantID, userID, ok := leadContext(c)
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
		if errors.Is(err, leadapp.ErrNotFound) {
			response.NotFound(c, "线索不存在")
			return
		}
		response.InternalError(c, "删除线索失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "lead.delete", "lead", &id, nil, old)
	response.Success(c, gin.H{"deleted": true})
}

func (h *LeadHandlers) Convert(c *gin.Context) {
	tenantID, userID, ok := leadContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	var req leadConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	var createAcc *leadapp.ConvertAccountInput
	if req.CreateAccount != nil && req.CreateAccount.Name != "" {
		createAcc = &leadapp.ConvertAccountInput{Name: req.CreateAccount.Name}
	}
	old, _ := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	dto, err := h.svc.Convert(c.Request.Context(), tenantID, userID, id, leadapp.ConvertInput{
		AccountID:     req.AccountID,
		ContactID:     req.ContactID,
		CreateAccount: createAcc,
	})
	if err != nil {
		if errors.Is(err, leadapp.ErrNotFound) {
			response.NotFound(c, "线索不存在")
			return
		}
		if errors.Is(err, leadapp.ErrInvalidStatusTransition) ||
			errors.Is(err, leadapp.ErrConvertNotAllowed) ||
			errors.Is(err, leadapp.ErrAlreadyConverted) {
			response.BadRequest(c, "invalid_status_transition")
			return
		}
		if errors.Is(err, leadapp.ErrConvertMissingAccount) {
			response.BadRequest(c, "convert_requires_account")
			return
		}
		response.InternalError(c, "转化线索失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "lead.convert", "lead", &id, dto, old)
	response.Success(c, dto)
}

func (h *LeadHandlers) EmotionJourney(c *gin.Context) {
	tenantID, userID, ok := leadContext(c)
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
		if errors.Is(err, leadapp.ErrNotFound) {
			response.NotFound(c, "线索不存在")
			return
		}
		response.InternalError(c, "获取情绪旅程失败")
		return
	}
	_ = c.Query("range")
	data := emotion.EmptyLeadJourney(id, dto.LifecycleStage)
	response.Success(c, data)
}

type leadConvertRequest struct {
	AccountID     *uuid.UUID              `json:"account_id"`
	ContactID     *uuid.UUID              `json:"contact_id"`
	CreateAccount *leadConvertAccountBody `json:"create_account"`
}

type leadConvertAccountBody struct {
	Name string `json:"name"`
}

type leadCreateRequest struct {
	Title             string     `json:"title" binding:"required,min=1,max=255"`
	Status            string     `json:"status"`
	Source            string     `json:"source"`
	Amount            float64    `json:"amount"`
	ExpectedCloseDate *string    `json:"expected_close_date"`
	OwnerID           *uuid.UUID `json:"owner_id"`
	LifecycleStage    string     `json:"lifecycle_stage"`
	Tags              []string   `json:"tags"`
}

type leadUpdateRequest struct {
	Title             *string    `json:"title"`
	Status            *string    `json:"status"`
	Source            *string    `json:"source"`
	Amount            *float64   `json:"amount"`
	ExpectedCloseDate *string    `json:"expected_close_date"`
	OwnerID           *uuid.UUID `json:"owner_id"`
	LifecycleStage    *string    `json:"lifecycle_stage"`
	Tags              *[]string  `json:"tags"`
}

func leadContext(c *gin.Context) (tenantID, userID uuid.UUID, ok bool) {
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
