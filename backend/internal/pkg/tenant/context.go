package tenant

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const tenantIDKey contextKey = "tenant_id"

// WithID 将租户 ID 写入 context
func WithID(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// IDFromContext 从 context 读取租户 ID
func IDFromContext(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(tenantIDKey)
	if v == nil {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok && id != uuid.Nil
}

// MustIDFromContext 读取租户 ID，不存在则 panic（仅用于已校验的中间件之后）
func MustIDFromContext(ctx context.Context) uuid.UUID {
	id, ok := IDFromContext(ctx)
	if !ok {
		panic("tenant_id missing in context")
	}
	return id
}
