package datascope

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	LevelAll        = "all"
	LevelDepartment = "department"
	LevelSelf       = "self"
)

// ScopeParams is the unified owner filter for CRM list/stats queries.
type ScopeParams struct {
	Level      string
	TenantID   uuid.UUID
	UserID     uuid.UUID
	Department string
}

func (p ScopeParams) ViewAll() bool { return p.Level == LevelAll }

func (p ScopeParams) APIScope() string {
	switch p.Level {
	case LevelAll:
		return LevelAll
	case LevelDepartment:
		return LevelDepartment
	default:
		return LevelSelf
	}
}

// ApplyOwnerScope filters by tenant admin (all), sales manager (department), or self.
func ApplyOwnerScope(db *gorm.DB, p ScopeParams) *gorm.DB {
	switch p.Level {
	case LevelAll:
		return db
	case LevelDepartment:
		if p.Department == "" {
			return db.Where("owner_id = ?", p.UserID)
		}
		return db.Where(`owner_id IN (
			SELECT ut.user_id FROM user_tenants ut
			WHERE ut.tenant_id = ? AND ut.department = ?
		)`, p.TenantID, p.Department)
	default:
		return db.Where("owner_id = ?", p.UserID)
	}
}
