package http

import (
	"errors"

	segmentapp "crm-backend/internal/application/segment"
	"crm-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func NewSegmentHandlers(svc *segmentapp.Service) *SegmentHandlers {
	return &SegmentHandlers{svc: svc}
}

type SegmentHandlers struct {
	svc *segmentapp.Service
}

func (h *SegmentHandlers) List(c *gin.Context) {
	tenantID, userID, ok := leadContext(c)
	if !ok {
		return
	}
	withCount := c.Query("with_count") == "1" || c.Query("with_count") == "true"
	items, err := h.svc.List(c.Request.Context(), tenantID, userID, withCount)
	if err != nil {
		response.InternalError(c, "获取分群列表失败")
		return
	}
	response.Success(c, gin.H{"items": items})
}

func (h *SegmentHandlers) Count(c *gin.Context) {
	tenantID, userID, ok := leadContext(c)
	if !ok {
		return
	}
	code := c.Param("code")
	entity := c.DefaultQuery("entity", "lead")
	data, err := h.svc.Count(c.Request.Context(), tenantID, userID, code, entity)
	if err != nil {
		writeSegmentError(c, err)
		return
	}
	response.Success(c, data)
}

func writeSegmentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, segmentapp.ErrNotFound):
		response.NotFound(c, "分群不存在")
	case errors.Is(err, segmentapp.ErrInvalidSegment):
		response.BadRequest(c, "invalid_segment_code")
	default:
		response.InternalError(c, "分群处理失败")
	}
}
