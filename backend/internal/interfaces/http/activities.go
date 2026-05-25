package http

import (
	"errors"
	"time"

	actapp "crm-backend/internal/application/activity"
	"crm-backend/internal/application/audit"
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewActivityHandlers(svc *actapp.Service, rec *audit.Recorder) *ActivityHandlers {
	return &ActivityHandlers{svc: svc, audit: rec}
}

type ActivityHandlers struct {
	svc   *actapp.Service
	audit *audit.Recorder
}

func (h *ActivityHandlers) List(c *gin.Context) {
	tenantID, userID, ok := activityContext(c)
	if !ok {
		return
	}
	subjectType := c.Query("subject_type")
	subjectID, err := uuid.Parse(c.Query("subject_id"))
	if err != nil || subjectType == "" {
		response.BadRequest(c, "subject_type 与 subject_id 必填")
		return
	}
	result, err := h.svc.List(c.Request.Context(), tenantID, userID, actapp.ListQuery{
		SubjectType: subjectType,
		SubjectID:   subjectID,
		Page:        queryInt(c, "page", 1),
		PageSize:    queryInt(c, "page_size", 20),
	})
	if err != nil {
		writeActivityError(c, err)
		return
	}
	response.SuccessPage(c, gin.H{"items": result.Items}, response.Pagination{
		Page:     result.Page,
		PageSize: result.Size,
		Total:    result.Total,
	})
}

func (h *ActivityHandlers) Summary(c *gin.Context) {
	tenantID, userID, ok := activityContext(c)
	if !ok {
		return
	}
	subjectType := c.Query("subject_type")
	var subjectID uuid.UUID
	if sid := c.Query("subject_id"); sid != "" {
		var err error
		subjectID, err = uuid.Parse(sid)
		if err != nil {
			response.BadRequest(c, "subject_id 格式无效")
			return
		}
	}
	data, err := h.svc.Summary(c.Request.Context(), tenantID, userID, subjectType, subjectID)
	if err != nil {
		writeActivityError(c, err)
		return
	}
	response.Success(c, data)
}

func (h *ActivityHandlers) Create(c *gin.Context) {
	tenantID, userID, ok := activityContext(c)
	if !ok {
		return
	}
	var req activityCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		response.BadRequest(c, "subject_id 格式无效")
		return
	}
	var occurred *time.Time
	if req.OccurredAt != nil && *req.OccurredAt != "" {
		t, err := time.Parse(time.RFC3339, *req.OccurredAt)
		if err != nil {
			response.BadRequest(c, "occurred_at 格式无效")
			return
		}
		occurred = &t
	}
	dto, err := h.svc.Create(c.Request.Context(), tenantID, userID, actapp.CreateInput{
		SubjectType:     req.SubjectType,
		SubjectID:       subjectID,
		EventType:       req.EventType,
		Direction:       req.Direction,
		Body:            req.Body,
		Metadata:        req.Metadata,
		Sentiment:       req.Sentiment,
		SentimentSource: req.SentimentSource,
		OccurredAt:      occurred,
		Label:           req.Label,
	})
	if err != nil {
		writeActivityError(c, err)
		return
	}
	recordAudit(c, h.audit, tenantID, "activity.create", "activity", &dto.ID, dto, nil)
	response.Created(c, dto)
}

func (h *ActivityHandlers) Get(c *gin.Context) {
	tenantID, userID, ok := activityContext(c)
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
		writeActivityError(c, err)
		return
	}
	response.Success(c, dto)
}

func (h *ActivityHandlers) Patch(c *gin.Context) {
	tenantID, userID, ok := activityContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID 格式无效")
		return
	}
	var req activityUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	in := actapp.UpdateInput{
		Body:       req.Body,
		Direction:  req.Direction,
		Sentiment:  req.Sentiment,
		Label:      req.Label,
	}
	if req.Metadata != nil {
		in.Metadata = *req.Metadata
		in.MetadataSet = true
	}
	if req.SentimentSource != nil {
		in.SentimentSource = req.SentimentSource
	}
	if req.OccurredAt != nil && *req.OccurredAt != "" {
		t, err := time.Parse(time.RFC3339, *req.OccurredAt)
		if err != nil {
			response.BadRequest(c, "occurred_at 格式无效")
			return
		}
		in.OccurredAt = &t
	}
	if req.ClearSentiment {
		in.SentimentClear = true
	}
	old, _ := h.svc.Get(c.Request.Context(), tenantID, userID, id)
	dto, err := h.svc.Update(c.Request.Context(), tenantID, userID, id, in)
	if err != nil {
		writeActivityError(c, err)
		return
	}
	recordAudit(c, h.audit, tenantID, "activity.update", "activity", &id, dto, old)
	response.Success(c, dto)
}

func (h *ActivityHandlers) Delete(c *gin.Context) {
	tenantID, userID, ok := activityContext(c)
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
		writeActivityError(c, err)
		return
	}
	recordAudit(c, h.audit, tenantID, "activity.delete", "activity", &id, nil, old)
	response.Success(c, gin.H{"deleted": true})
}

type activityCreateRequest struct {
	SubjectType     string         `json:"subject_type" binding:"required"`
	SubjectID       string         `json:"subject_id" binding:"required"`
	EventType       string         `json:"event_type" binding:"required"`
	Direction       string         `json:"direction"`
	Body            string         `json:"body"`
	Metadata        map[string]any `json:"metadata"`
	Sentiment       *string        `json:"sentiment"`
	SentimentSource string         `json:"sentiment_source"`
	OccurredAt      *string        `json:"occurred_at"`
	Label           string         `json:"label"`
}

type activityUpdateRequest struct {
	Body             *string         `json:"body"`
	Direction        *string         `json:"direction"`
	Metadata         *map[string]any `json:"metadata"`
	Sentiment        *string         `json:"sentiment"`
	SentimentSource  *string         `json:"sentiment_source"`
	OccurredAt       *string         `json:"occurred_at"`
	Label            *string         `json:"label"`
	ClearSentiment   bool            `json:"clear_sentiment"`
}

func activityContext(c *gin.Context) (tenantID, userID uuid.UUID, ok bool) {
	return leadContext(c)
}

func writeActivityError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, actapp.ErrNotFound):
		response.NotFound(c, "跟进记录不存在")
	case errors.Is(err, actapp.ErrSubjectNotFound):
		response.NotFound(c, "主体不存在")
	case errors.Is(err, actapp.ErrInvalidEvent),
		errors.Is(err, actapp.ErrInvalidSubject),
		errors.Is(err, actapp.ErrInvalidSentiment),
		errors.Is(err, actapp.ErrInvalidSource),
		errors.Is(err, actapp.ErrInvalidDirection),
		errors.Is(err, actapp.ErrSentimentAI):
		response.BadRequest(c, err.Error())
	default:
		response.InternalError(c, "处理跟进记录失败")
	}
}
