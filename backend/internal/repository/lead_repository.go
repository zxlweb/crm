package repository

import (
	"context"
	"errors"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/infrastructure/persistence"
	"crm-backend/internal/pkg/crm"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrLeadNotFound = errors.New("lead not found")

type LeadListFilter struct {
	Page               int
	PageSize           int
	Search             string
	Status             string
	Source             string
	LifecycleStage     string
	RelationshipHealth string
	Segment            string
	SegmentOpts        crm.SegmentApplyOpts
	OwnerID            *uuid.UUID
	ViewAll            bool
	UserID             uuid.UUID
}

type LeadRepository interface {
	List(ctx context.Context, tenantID uuid.UUID, f LeadListFilter) ([]domain.Lead, int64, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Lead, error)
	Create(ctx context.Context, l *domain.Lead) error
	Update(ctx context.Context, l *domain.Lead) error
	UpdateEngagementFromActivity(ctx context.Context, tenantID, id, updatedBy uuid.UUID, last *time.Time, score int16) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	StatsBySource(ctx context.Context, tenantID uuid.UUID, f LeadStatsFilter) ([]LabelCount, int64, error)
	StatsByStatus(ctx context.Context, tenantID uuid.UUID, f LeadStatsFilter) ([]LabelCount, int64, error)
	StatsTrend(ctx context.Context, tenantID uuid.UUID, f LeadStatsFilter, granularity string) ([]TrendPoint, error)
	StatsFunnel(ctx context.Context, tenantID uuid.UUID, f LeadStatsFilter) ([]LabelCount, error)
	CountScoped(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error)
	DailyCreatedCounts(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID, days int) ([]int64, error)
	CountLowEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (int64, error)
	AvgEngagement(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID) (float64, error)
	ListPriorityCandidates(ctx context.Context, tenantID uuid.UUID, viewAll bool, userID uuid.UUID, limit int) ([]domain.Lead, error)
}

type GormLeadRepository struct {
	db *gorm.DB
}

func NewLeadRepository(db *gorm.DB) LeadRepository {
	return &GormLeadRepository{db: db}
}

func (r *GormLeadRepository) base(ctx context.Context, tenantID uuid.UUID) *gorm.DB {
	return persistence.DBFromContext(r.db, ctx).Model(&domain.Lead{}).Where("tenant_id = ?", tenantID)
}

func (r *GormLeadRepository) List(ctx context.Context, tenantID uuid.UUID, f LeadListFilter) ([]domain.Lead, int64, error) {
	q := r.base(ctx, tenantID)
	if !f.ViewAll {
		q = q.Where("(owner_id = ? OR owner_id IS NULL)", f.UserID)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("(title ILIKE ? OR source ILIKE ?)", like, like)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Source != "" {
		q = q.Where("source = ?", f.Source)
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
	if f.Segment != "" {
		if err := crm.ApplyLeadSegmentFilter(q, f.Segment, f.SegmentOpts); err != nil {
			if errors.Is(err, crm.ErrInvalidSegmentCode) {
				return nil, 0, err
			}
			return nil, 0, err
		}
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

	var items []domain.Lead
	err := q.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *GormLeadRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID, viewAll bool, userID uuid.UUID) (*domain.Lead, error) {
	q := r.base(ctx, tenantID).Where("id = ?", id)
	if !viewAll {
		q = q.Where("(owner_id = ? OR owner_id IS NULL)", userID)
	}
	var l domain.Lead
	err := q.First(&l).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrLeadNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *GormLeadRepository) Create(ctx context.Context, l *domain.Lead) error {
	if l.LifecycleStage == "" {
		l.LifecycleStage = "acquire"
	}
	if !crm.ValidLifecycleStage(l.LifecycleStage) {
		return errors.New("invalid lifecycle_stage")
	}
	if l.Status == "" {
		l.Status = "new"
	}
	if !crm.ValidLeadStatus(l.Status) {
		return errors.New("invalid status")
	}
	return persistence.DBFromContext(r.db, ctx).Create(l).Error
}

func (r *GormLeadRepository) Update(ctx context.Context, l *domain.Lead) error {
	if l.LifecycleStage != "" && !crm.ValidLifecycleStage(l.LifecycleStage) {
		return errors.New("invalid lifecycle_stage")
	}
	if l.Status != "" && !crm.ValidLeadStatus(l.Status) {
		return errors.New("invalid status")
	}
	return persistence.DBFromContext(r.db, ctx).Save(l).Error
}

func (r *GormLeadRepository) UpdateEngagementFromActivity(ctx context.Context, tenantID, id, updatedBy uuid.UUID, last *time.Time, score int16) error {
	return r.base(ctx, tenantID).Where("id = ?", id).Updates(map[string]any{
		"last_activity_at":    last,
		"engagement_score":    score,
		"relationship_health": crm.RelationshipHealthFromScore(score),
		"updated_by":          updatedBy,
	}).Error
}

func (r *GormLeadRepository) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	res := persistence.DBFromContext(r.db, ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&domain.Lead{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrLeadNotFound
	}
	return nil
}
