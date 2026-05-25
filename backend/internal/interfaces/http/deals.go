package http

import (
	"errors"
	"strconv"

	dealapp "crm-backend/internal/application/deal"
	"crm-backend/internal/application/audit"
	"crm-backend/internal/pkg/response"
	"crm-backend/internal/pkg/tenant"
	"crm-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewDealHandlers(svc *dealapp.Service, rec *audit.Recorder) *DealHandlers {
	return &DealHandlers{svc: svc, audit: rec}
}

type DealHandlers struct {
	svc   *dealapp.Service
	audit *audit.Recorder
}

func (h *DealHandlers) List(c *gin.Context) {
	tenantID, userID, ok := dealContext(c)
	if !ok {
		return
	}
	q := dealapp.ListQuery{
		Page:     queryInt(c, "page", 1),
		PageSize: queryInt(c, "page_size", 20),
		Search:   c.Query("search"),
		Stage:    c.Query("stage"),
		Stages:   repository.ParseDealStages(c.Query("stages")),
	}
	if oid := c.Query("owner_id"); oid != "" {
		id, err := uuid.Parse(oid)
		if err != nil {
			response.BadRequest(c, "owner_id 格式无效")
			return
		}
		q.OwnerID = &id
	}
	if aid := c.Query("account_id"); aid != "" {
		id, err := uuid.Parse(aid)
		if err != nil {
			response.BadRequest(c, "account_id 格式无效")
			return
		}
		q.AccountID = &id
	}
	if lid := c.Query("lead_id"); lid != "" {
		id, err := uuid.Parse(lid)
		if err != nil {
			response.BadRequest(c, "lead_id 格式无效")
			return
		}
		q.LeadID = &id
	}
	if v := c.Query("expected_close_from"); v != "" {
		q.ExpectedCloseFrom = &v
	}
	if v := c.Query("expected_close_to"); v != "" {
		q.ExpectedCloseTo = &v
	}
	if v := c.Query("min_amount"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			response.BadRequest(c, "min_amount 格式无效")
			return
		}
		q.MinAmount = &f
	}
	if v := c.Query("max_amount"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			response.BadRequest(c, "max_amount 格式无效")
			return
		}
		q.MaxAmount = &f
	}
	result, err := h.svc.List(c.Request.Context(), tenantID, userID, q)
	if err != nil {
		response.InternalError(c, "获取商机列表失败")
		return
	}
	response.SuccessPage(c, gin.H{"items": result.Items}, response.Pagination{
		Page:     result.Page,
		PageSize: result.Size,
		Total:    result.Total,
	})
}

func (h *DealHandlers) Create(c *gin.Context) {
	tenantID, userID, ok := dealContext(c)
	if !ok {
		return
	}
	var req dealCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	dto, err := h.svc.Create(c.Request.Context(), tenantID, userID, dealapp.CreateInput{
		Title:             req.Title,
		Stage:             req.Stage,
		Amount:            req.Amount,
		Currency:          req.Currency,
		Probability:       req.Probability,
		ExpectedCloseDate: req.ExpectedCloseDate,
		AccountID:         req.AccountID,
		LeadID:            req.LeadID,
		ContactID:         req.ContactID,
		OwnerID:           req.OwnerID,
		Description:       req.Description,
		Tags:              req.Tags,
	})
	if err != nil {
		if errors.Is(err, dealapp.ErrInvalidStage) {
			response.BadRequest(c, "invalid stage")
			return
		}
		if errors.Is(err, dealapp.ErrInvalidCurrency) || errors.Is(err, dealapp.ErrInvalidAmount) || errors.Is(err, dealapp.ErrInvalidProbability) {
			response.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, repository.ErrAccountNotFound) {
			response.BadRequest(c, "account_not_found")
			return
		}
		response.InternalError(c, "创建商机失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "deal.create", "deal", &dto.ID, dto, nil)
	response.Created(c, dto)
}

func (h *DealHandlers) Get(c *gin.Context) {
	tenantID, userID, ok := dealContext(c)
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
		if errors.Is(err, dealapp.ErrNotFound) {
			response.NotFound(c, "商机不存在")
			return
		}
		response.InternalError(c, "获取商机详情失败")
		return
	}
	response.Success(c, dto)
}

func (h *DealHandlers) Put(c *gin.Context) {
	h.update(c, true)
}

func (h *DealHandlers) Patch(c *gin.Context) {
	h.update(c, false)
}

func (h *DealHandlers) update(c *gin.Context, full bool) {
	tenantID, userID, ok := dealContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	var req dealUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	in := dealapp.UpdateInput{
		Title:             req.Title,
		Stage:             req.Stage,
		Amount:            req.Amount,
		Currency:          req.Currency,
		Probability:       req.Probability,
		ExpectedCloseDate: req.ExpectedCloseDate,
		AccountID:         req.AccountID,
		LeadID:            req.LeadID,
		ContactID:         req.ContactID,
		OwnerID:           req.OwnerID,
		Description:       req.Description,
		LostReason:        req.LostReason,
	}
	if req.Tags != nil {
		in.Tags = *req.Tags
		in.TagsSet = true
	}
	old, _ := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	dto, err := h.svc.Update(c.Request.Context(), tenantID, userID, id, in, full)
	if err != nil {
		if errors.Is(err, dealapp.ErrNotFound) {
			response.NotFound(c, "商机不存在")
			return
		}
		if errors.Is(err, dealapp.ErrInvalidStageTransition) {
			response.BadRequest(c, "invalid_stage_transition")
			return
		}
		if errors.Is(err, dealapp.ErrDealClosedReadonly) {
			response.BadRequest(c, "deal_closed_readonly")
			return
		}
		if errors.Is(err, dealapp.ErrInvalidStage) || errors.Is(err, dealapp.ErrInvalidCurrency) ||
			errors.Is(err, dealapp.ErrInvalidAmount) || errors.Is(err, dealapp.ErrInvalidProbability) {
			response.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, repository.ErrAccountNotFound) {
			response.BadRequest(c, "account_not_found")
			return
		}
		response.InternalError(c, "更新商机失败")
		return
	}
	action := "deal.update"
	if req.Stage != nil && old != nil && *req.Stage != old.Stage {
		action = "deal.stage_change"
	}
	recordAudit(c, h.audit, tenantID, action, "deal", &id, dto, old)
	response.Success(c, dto)
}

func (h *DealHandlers) Delete(c *gin.Context) {
	tenantID, userID, ok := dealContext(c)
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
		if errors.Is(err, dealapp.ErrNotFound) {
			response.NotFound(c, "商机不存在")
			return
		}
		response.InternalError(c, "删除商机失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "deal.delete", "deal", &id, nil, old)
	response.Success(c, gin.H{"deleted": true})
}

func (h *DealHandlers) Pipeline(c *gin.Context) {
	tenantID, userID, ok := dealContext(c)
	if !ok {
		return
	}
	q := dealapp.PipelineQuery{}
	if oid := c.Query("owner_id"); oid != "" {
		id, err := uuid.Parse(oid)
		if err != nil {
			response.BadRequest(c, "owner_id 格式无效")
			return
		}
		q.OwnerID = &id
	}
	if aid := c.Query("account_id"); aid != "" {
		id, err := uuid.Parse(aid)
		if err != nil {
			response.BadRequest(c, "account_id 格式无效")
			return
		}
		q.AccountID = &id
	}
	data, err := h.svc.Pipeline(c.Request.Context(), tenantID, userID, q)
	if err != nil {
		response.InternalError(c, "获取 Pipeline 失败")
		return
	}
	response.Success(c, data)
}

func (h *DealHandlers) PutStage(c *gin.Context) {
	tenantID, userID, ok := dealContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	var req dealStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	old, _ := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	dto, err := h.svc.UpdateStage(c.Request.Context(), tenantID, userID, id, dealapp.StageInput{
		Stage:      req.Stage,
		LostReason: req.LostReason,
		Note:       req.Note,
	})
	if err != nil {
		if errors.Is(err, dealapp.ErrNotFound) {
			response.NotFound(c, "商机不存在")
			return
		}
		if errors.Is(err, dealapp.ErrInvalidStageTransition) {
			response.BadRequest(c, "invalid_stage_transition")
			return
		}
		if errors.Is(err, dealapp.ErrDealClosedReadonly) {
			response.BadRequest(c, "deal_closed_readonly")
			return
		}
		if errors.Is(err, dealapp.ErrInvalidStage) {
			response.BadRequest(c, "invalid stage")
			return
		}
		response.InternalError(c, "更新商机阶段失败")
		return
	}
	recordAudit(c, h.audit, tenantID, "deal.stage_change", "deal", &id, dto, old)
	response.Success(c, dto)
}

type dealCreateRequest struct {
	Title             string     `json:"title" binding:"required,min=1,max=255"`
	Stage             string     `json:"stage"`
	Amount            float64    `json:"amount"`
	Currency          string     `json:"currency"`
	Probability       *int       `json:"probability"`
	ExpectedCloseDate *string    `json:"expected_close_date"`
	AccountID         *uuid.UUID `json:"account_id"`
	LeadID            *uuid.UUID `json:"lead_id"`
	ContactID         *uuid.UUID `json:"contact_id"`
	OwnerID           *uuid.UUID `json:"owner_id"`
	Description       string     `json:"description"`
	Tags              []string   `json:"tags"`
}

type dealUpdateRequest struct {
	Title             *string    `json:"title"`
	Stage             *string    `json:"stage"`
	Amount            *float64   `json:"amount"`
	Currency          *string    `json:"currency"`
	Probability       *int       `json:"probability"`
	ExpectedCloseDate *string    `json:"expected_close_date"`
	AccountID         *uuid.UUID `json:"account_id"`
	LeadID            *uuid.UUID `json:"lead_id"`
	ContactID         *uuid.UUID `json:"contact_id"`
	OwnerID           *uuid.UUID `json:"owner_id"`
	Description       *string    `json:"description"`
	LostReason        *string    `json:"lost_reason"`
	Tags              *[]string  `json:"tags"`
}

type dealStageRequest struct {
	Stage      string  `json:"stage" binding:"required"`
	LostReason *string `json:"lost_reason"`
	Note       string  `json:"note"`
}

func dealContext(c *gin.Context) (tenantID, userID uuid.UUID, ok bool) {
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
