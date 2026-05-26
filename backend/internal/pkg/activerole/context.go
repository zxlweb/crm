package activerole

import "context"

type ctxKey struct{}

// WithID stores the active role for the current request (from JWT).
func WithID(ctx context.Context, roleID string) context.Context {
	if roleID == "" {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, roleID)
}

// FromContext returns the active role ID when set.
func FromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKey{}).(string)
	return v, ok && v != ""
}
