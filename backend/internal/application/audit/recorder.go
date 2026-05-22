package audit

import (
	"context"
	"encoding/json"

	"crm-backend/internal/domain"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Recorder 异步友好地写入审计日志（失败不阻断主流程）
type Recorder struct {
	repo repository.AuditRepository
}

func NewRecorder(repo repository.AuditRepository) *Recorder {
	return &Recorder{repo: repo}
}

type Entry struct {
	TenantID     uuid.UUID
	UserID       *uuid.UUID
	Action       string
	ResourceType string
	ResourceID   *uuid.UUID
	OldValue     any
	NewValue     any
	IPAddress    string
}

func (r *Recorder) Record(ctx context.Context, e Entry) {
	if r == nil || r.repo == nil {
		return
	}
	log := &domain.AuditLog{
		TenantID:     e.TenantID,
		UserID:       e.UserID,
		Action:       e.Action,
		ResourceType: e.ResourceType,
		ResourceID:   e.ResourceID,
		IPAddress:    e.IPAddress,
		OldValue:     toJSON(e.OldValue),
		NewValue:     toJSON(e.NewValue),
	}
	_ = r.repo.Create(ctx, log)
}

func toJSON(v any) datatypes.JSON {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return datatypes.JSON(b)
}
