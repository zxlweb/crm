package repository

import (
	"context"
	"errors"

	"crm-backend/internal/domain"
	"crm-backend/internal/infrastructure/persistence"
	"crm-backend/internal/pkg/crm"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrAccountNotFound = errors.New("account not found")

type AccountListFilter struct {
	Page               int
	PageSize           int
	Search             string
	LifecycleStage     string
	RelationshipHealth string
	OwnerID            *uuid.UUID
	ViewAll            bool
	UserID             uuid.UUID
}

type AccountRepository interface {
	List(ctx context.Context, tenantID uuid.UUID, f AccountListFilter) ([]domain.Account, int64, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Account, error)
	Create(ctx context.Context, a *domain.Account) error
	Update(ctx context.Context, a *domain.Account) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
}

type GormAccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &GormAccountRepository{db: db}
}

func healthSQLExpr() string {
	return `CASE
		WHEN engagement_score >= 70 THEN 'high'
		WHEN engagement_score >= 40 THEN 'medium'
		ELSE 'low'
	END`
}

func (r *GormAccountRepository) base(ctx context.Context, tenantID uuid.UUID) *gorm.DB {
	return persistence.DBFromContext(r.db, ctx).Model(&domain.Account{})
}

func (r *GormAccountRepository) List(ctx context.Context, tenantID uuid.UUID, f AccountListFilter) ([]domain.Account, int64, error) {
	q := r.base(ctx, tenantID)
	if !f.ViewAll {
		q = q.Where("(owner_id = ? OR owner_id IS NULL)", f.UserID)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("name ILIKE ?", like)
	}
	if f.LifecycleStage != "" {
		q = q.Where("lifecycle_stage = ?", f.LifecycleStage)
	}
	if f.RelationshipHealth != "" {
		q = q.Where(healthSQLExpr()+" = ?", f.RelationshipHealth)
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

	var items []domain.Account
	err := q.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *GormAccountRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Account, error) {
	q := r.base(ctx, tenantID).Where("id = ?", id)
	if !viewAll {
		q = q.Where("(owner_id = ? OR owner_id IS NULL)", userID)
	}
	var a domain.Account
	err := q.First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *GormAccountRepository) Create(ctx context.Context, a *domain.Account) error {
	if a.LifecycleStage == "" {
		a.LifecycleStage = "acquire"
	}
	if !crm.ValidLifecycleStage(a.LifecycleStage) {
		return errors.New("invalid lifecycle_stage")
	}
	return persistence.DBFromContext(r.db, ctx).Create(a).Error
}

func (r *GormAccountRepository) Update(ctx context.Context, a *domain.Account) error {
	if a.LifecycleStage != "" && !crm.ValidLifecycleStage(a.LifecycleStage) {
		return errors.New("invalid lifecycle_stage")
	}
	return persistence.DBFromContext(r.db, ctx).Save(a).Error
}

func (r *GormAccountRepository) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	res := persistence.DBFromContext(r.db, ctx).
		Where("id = ?", id).
		Delete(&domain.Account{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrAccountNotFound
	}
	return nil
}
