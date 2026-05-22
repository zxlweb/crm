package http

import (
	"crm-backend/internal/application/audit"
	"crm-backend/internal/pkg/httputil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func recordAudit(c *gin.Context, rec *audit.Recorder, tenantID uuid.UUID, action, resourceType string, resourceID *uuid.UUID, newValue, oldValue any) {
	if rec == nil || tenantID == uuid.Nil {
		return
	}
	var userID *uuid.UUID
	if id, err := uuid.Parse(c.GetString("user_id")); err == nil {
		userID = &id
	}
	rec.Record(c.Request.Context(), audit.Entry{
		TenantID:     tenantID,
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		NewValue:     newValue,
		OldValue:     oldValue,
		IPAddress:    httputil.ClientIP(c),
	})
}
