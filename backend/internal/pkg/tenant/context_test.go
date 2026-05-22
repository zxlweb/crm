package tenant

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestWithIDAndIDFromContext(t *testing.T) {
	tid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ctx := WithID(context.Background(), tid)

	got, ok := IDFromContext(ctx)
	if !ok {
		t.Fatal("expected tenant id in context")
	}
	if got != tid {
		t.Fatalf("got %v, want %v", got, tid)
	}
}

func TestIDFromContext_Missing(t *testing.T) {
	_, ok := IDFromContext(context.Background())
	if ok {
		t.Fatal("expected false when tenant id not set")
	}
}

func TestIDFromContext_NilUUID(t *testing.T) {
	ctx := context.WithValue(context.Background(), tenantIDKey, uuid.Nil)
	_, ok := IDFromContext(ctx)
	if ok {
		t.Fatal("expected false for nil uuid")
	}
}

func TestMustIDFromContext_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	MustIDFromContext(context.Background())
}
