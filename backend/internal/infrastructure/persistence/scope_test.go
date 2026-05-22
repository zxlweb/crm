package persistence

import (
	"context"
	"testing"

	"crm-backend/internal/pkg/tenant"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testLead 仅用于 SQLite 内存测试（避免 Postgres uuid 方言）
type testLead struct {
	ID       string `gorm:"primaryKey"`
	TenantID string `gorm:"index;not null"`
	Title    string
}

func TestTenantScope_FiltersByTenantID(t *testing.T) {
	db := openTestDB(t)

	tenantA := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	tenantB := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"

	db.Create(&testLead{ID: "1", TenantID: tenantA, Title: "Lead A"})
	db.Create(&testLead{ID: "2", TenantID: tenantB, Title: "Lead B"})

	tid := uuid.MustParse(tenantA)
	ctx := tenant.WithID(context.Background(), tid)

	var leads []testLead
	if err := DBFromContext(db, ctx).Find(&leads).Error; err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(leads) != 1 || leads[0].Title != "Lead A" {
		t.Fatalf("expected 1 lead for tenant A, got %+v", leads)
	}
}

func TestTenantScope_NilTenantReturnsAll(t *testing.T) {
	db := openTestDB(t)

	tenantA := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	db.Create(&testLead{ID: "1", TenantID: tenantA, Title: "Lead A"})
	db.Create(&testLead{ID: "2", TenantID: tenantA, Title: "Lead A2"})

	var leads []testLead
	if err := db.Scopes(TenantScope(uuid.Nil)).Find(&leads).Error; err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(leads) != 2 {
		t.Fatalf("expected 2 leads without scope, got %d", len(leads))
	}
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&testLead{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}
