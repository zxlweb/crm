package repository

import (
	"context"
	"errors"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/infrastructure/persistence"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/datascope"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrContactNotFound = errors.New("contact not found")

type ContactListFilter struct {
	Page               int
	PageSize           int
	Search             string
	LifecycleStage     string
	RelationshipHealth string
	AccountID          *uuid.UUID
	OwnerID            *uuid.UUID
	Scope              datascope.ScopeParams
}

type ContactRepository interface {
	List(ctx context.Context, tenantID uuid.UUID, f ContactListFilter) ([]domain.Contact, int64, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID, scope datascope.ScopeParams) (*domain.Contact, error)
	Create(ctx context.Context, c *domain.Contact) error
	Update(ctx context.Context, c *domain.Contact) error
	UpdateEngagementFromActivity(ctx context.Context, tenantID, id, updatedBy uuid.UUID, last *time.Time, score int16) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
}

type GormContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &GormContactRepository{db: db}
}

func (r *GormContactRepository) base(ctx context.Context, tenantID uuid.UUID) *gorm.DB {
	return persistence.DBFromContext(r.db, ctx).Model(&domain.Contact{}).Where("tenant_id = ?", tenantID)
}

func (r *GormContactRepository) List(ctx context.Context, tenantID uuid.UUID, f ContactListFilter) ([]domain.Contact, int64, error) {
	q := datascope.ApplyOwnerScope(r.base(ctx, tenantID), f.Scope)
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where(
			"(first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?)",
			like, like, like, like,
		)
	}
	if f.LifecycleStage != "" {
		q = q.Where("lifecycle_stage = ?", f.LifecycleStage)
	}
	if f.RelationshipHealth != "" {
		q = q.Where(healthSQLExpr()+" = ?", f.RelationshipHealth)
	}
	if f.AccountID != nil {
		q = q.Where("account_id = ?", *f.AccountID)
	}
	if f.OwnerID != nil {
		q = q.Where("owner_id = ?", *f.OwnerID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := f.Page
	if page < 1 {
		page = 1
	}
	pageSize := f.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var items []domain.Contact
	err := q.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *GormContactRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID, scope datascope.ScopeParams) (*domain.Contact, error) {
	q := datascope.ApplyOwnerScope(r.base(ctx, tenantID).Where("id = ?", id), scope)
	var c domain.Contact
	err := q.First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrContactNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *GormContactRepository) Create(ctx context.Context, c *domain.Contact) error {
	if c.LifecycleStage == "" {
		c.LifecycleStage = "acquire"
	}
	if !crm.ValidLifecycleStage(c.LifecycleStage) {
		return errors.New("invalid lifecycle_stage")
	}
	return persistence.DBFromContext(r.db, ctx).Create(c).Error
}

func (r *GormContactRepository) Update(ctx context.Context, c *domain.Contact) error {
	if c.LifecycleStage != "" && !crm.ValidLifecycleStage(c.LifecycleStage) {
		return errors.New("invalid lifecycle_stage")
	}
	return persistence.DBFromContext(r.db, ctx).Save(c).Error
}

func (r *GormContactRepository) UpdateEngagementFromActivity(ctx context.Context, tenantID, id, updatedBy uuid.UUID, last *time.Time, score int16) error {
	return r.base(ctx, tenantID).Where("id = ?", id).Updates(map[string]any{
		"last_activity_at": last,
		"engagement_score": score,
		"updated_by":       updatedBy,
	}).Error
}

func (r *GormContactRepository) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	res := persistence.DBFromContext(r.db, ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&domain.Contact{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrContactNotFound
	}
	return nil
}
