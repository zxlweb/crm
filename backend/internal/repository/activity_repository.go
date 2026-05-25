package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/infrastructure/persistence"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrActivityNotFound = errors.New("activity not found")

type ActivityListFilter struct {
	SubjectType string
	SubjectID   uuid.UUID
	Page        int
	PageSize    int
}

type ActivityRepository interface {
	List(ctx context.Context, tenantID uuid.UUID, f ActivityListFilter) ([]domain.Activity, int64, error)
	ListForJourney(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID, since *time.Time) ([]domain.Activity, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Activity, error)
	Create(ctx context.Context, a *domain.Activity) error
	Update(ctx context.Context, a *domain.Activity) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	SummaryByEventType(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID) ([]LabelCount, int64, error)
	LatestOccurredAt(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID) (*time.Time, error)
	CountBySubject(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID) (int64, error)
	CountSince(ctx context.Context, tenantID uuid.UUID, since time.Time) (int64, error)
	CountLeadTouchesSince(ctx context.Context, tenantID uuid.UUID, since time.Time, viewAll bool, userID uuid.UUID) (int64, error)
	CountAccountTouchesSince(ctx context.Context, tenantID uuid.UUID, since time.Time, viewAll bool, userID uuid.UUID) (int64, error)
}

type GormActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &GormActivityRepository{db: db}
}

func (r *GormActivityRepository) base(ctx context.Context, tenantID uuid.UUID) *gorm.DB {
	return persistence.DBFromContext(r.db, ctx).Model(&domain.Activity{}).Where("tenant_id = ?", tenantID)
}

func (r *GormActivityRepository) subjectQuery(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID) *gorm.DB {
	q := r.base(ctx, tenantID)
	if subjectType != "" {
		q = q.Where("subject_type = ?", subjectType)
	}
	if subjectID != uuid.Nil {
		q = q.Where("subject_id = ?", subjectID)
	}
	return q
}

func (r *GormActivityRepository) List(ctx context.Context, tenantID uuid.UUID, f ActivityListFilter) ([]domain.Activity, int64, error) {
	q := r.subjectQuery(ctx, tenantID, f.SubjectType, f.SubjectID)
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
	var items []domain.Activity
	err := q.Order("occurred_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *GormActivityRepository) ListForJourney(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID, since *time.Time) ([]domain.Activity, error) {
	q := r.subjectQuery(ctx, tenantID, subjectType, subjectID)
	if since != nil {
		q = q.Where("occurred_at >= ?", *since)
	}
	var items []domain.Activity
	err := q.Order("occurred_at ASC").Find(&items).Error
	return items, err
}

func (r *GormActivityRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Activity, error) {
	var a domain.Activity
	err := r.base(ctx, tenantID).Where("id = ?", id).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrActivityNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *GormActivityRepository) Create(ctx context.Context, a *domain.Activity) error {
	return persistence.DBFromContext(r.db, ctx).Create(a).Error
}

func (r *GormActivityRepository) Update(ctx context.Context, a *domain.Activity) error {
	return persistence.DBFromContext(r.db, ctx).Save(a).Error
}

func (r *GormActivityRepository) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	res := persistence.DBFromContext(r.db, ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&domain.Activity{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrActivityNotFound
	}
	return nil
}

func (r *GormActivityRepository) SummaryByEventType(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID) ([]LabelCount, int64, error) {
	type row struct {
		Label string
		Count int64
	}
	var rows []row
	err := r.subjectQuery(ctx, tenantID, subjectType, subjectID).
		Select("event_type AS label, COUNT(*) AS count").
		Group("event_type").
		Order("count DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	out := make([]LabelCount, len(rows))
	for i, row := range rows {
		out[i] = LabelCount{Label: row.Label, Count: row.Count}
		total += row.Count
	}
	return out, total, nil
}

func (r *GormActivityRepository) LatestOccurredAt(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID) (*time.Time, error) {
	var max sql.NullTime
	err := r.subjectQuery(ctx, tenantID, subjectType, subjectID).
		Select("MAX(occurred_at)").
		Scan(&max).Error
	if err != nil || !max.Valid {
		return nil, err
	}
	t := max.Time
	return &t, nil
}

func (r *GormActivityRepository) CountBySubject(ctx context.Context, tenantID uuid.UUID, subjectType string, subjectID uuid.UUID) (int64, error) {
	var n int64
	err := r.subjectQuery(ctx, tenantID, subjectType, subjectID).Count(&n).Error
	return n, err
}
