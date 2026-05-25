package crm

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestValidSegmentCode(t *testing.T) {
	for _, code := range []string{"high_value", "churn_risk", "new_potential", "needs_activation", "revive_pool"} {
		if !ValidSegmentCode(code) {
			t.Fatalf("expected valid: %s", code)
		}
	}
	if ValidSegmentCode("unknown_segment") {
		t.Fatal("expected invalid code")
	}
}

func TestSegmentEntityForCode(t *testing.T) {
	lead, acct, err := SegmentEntityForCode("high_value")
	if err != nil || !lead || !acct {
		t.Fatalf("high_value: lead=%v account=%v err=%v", lead, acct, err)
	}
	lead, acct, err = SegmentEntityForCode("churn_risk")
	if err != nil || !lead || !acct {
		t.Fatalf("churn_risk: lead=%v account=%v err=%v", lead, acct, err)
	}
}

func TestApplyLeadSegmentFilter_HighValue(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE TABLE leads (
		id TEXT PRIMARY KEY,
		tenant_id TEXT,
		amount REAL,
		status TEXT,
		lifecycle_stage TEXT,
		created_at DATETIME,
		last_activity_at DATETIME,
		deleted_at DATETIME
	)`).Error; err != nil {
		t.Fatal(err)
	}
	highID := uuid.New().String()
	lowID := uuid.New().String()
	tenantID := uuid.New().String()
	now := time.Now().UTC()
	for _, row := range []struct {
		id, title string
		amount    float64
	}{
		{highID, "high", 200000},
		{lowID, "low", 1000},
	} {
		if err := db.Exec(
			`INSERT INTO leads (id, tenant_id, amount, status, lifecycle_stage, created_at) VALUES (?, ?, ?, 'new', 'acquire', ?)`,
			row.id, tenantID, row.amount, now,
		).Error; err != nil {
			t.Fatal(err)
		}
	}

	q := db.Table("leads").Where("tenant_id = ?", tenantID)
	if err := ApplyLeadSegmentFilter(q, "high_value", SegmentApplyOpts{HighValueAmount: 100000}); err != nil {
		t.Fatal(err)
	}
	var ids []string
	if err := q.Pluck("id", &ids).Error; err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != highID {
		t.Fatalf("filtered ids: %v", ids)
	}
}
