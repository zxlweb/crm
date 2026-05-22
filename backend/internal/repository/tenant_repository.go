package repository

import (
	"context"
	"errors"

	"crm-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrTenantNotFound = errors.New("tenant not found")

type TenantListFilter struct {
	Page     int
	PageSize int
	Search   string
	IsActive *bool
}

type TenantActivityPoint struct {
	Date       string
	NewTenants int64
	Logins     int64
}

type TenantRepository interface {
	List(ctx context.Context, filter TenantListFilter) ([]domain.Tenant, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
	CountUsers(ctx context.Context, tenantID uuid.UUID) (int64, error)
	CountAll(ctx context.Context) (total int64, active int64, err error)
	CountAllUsers(ctx context.Context) (int64, error)
	TenantActivityTrend(ctx context.Context, days int) ([]TenantActivityPoint, error)
}

type GormTenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &GormTenantRepository{db: db}
}

func (r *GormTenantRepository) List(ctx context.Context, filter TenantListFilter) ([]domain.Tenant, int64, error) {
	q := r.db.WithContext(ctx).Model(&domain.Tenant{})
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		q = q.Where("name ILIKE ? OR domain ILIKE ?", like, like)
	}
	if filter.IsActive != nil {
		q = q.Where("is_active = ?", *filter.IsActive)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var rows []domain.Tenant
	err := q.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error
	return rows, total, err
}

func (r *GormTenantRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tenant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTenantNotFound
	}
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *GormTenantRepository) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	res := r.db.WithContext(ctx).Model(&domain.Tenant{}).
		Where("id = ?", id).
		Update("is_active", active)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrTenantNotFound
	}
	return nil
}

func (r *GormTenantRepository) CountUsers(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("user_tenants").
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	return count, err
}

func (r *GormTenantRepository) CountAll(ctx context.Context) (int64, int64, error) {
	var total, active int64
	if err := r.db.WithContext(ctx).Model(&domain.Tenant{}).Count(&total).Error; err != nil {
		return 0, 0, err
	}
	if err := r.db.WithContext(ctx).Model(&domain.Tenant{}).Where("is_active = ?", true).Count(&active).Error; err != nil {
		return 0, 0, err
	}
	return total, active, nil
}

func (r *GormTenantRepository) CountAllUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.User{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *GormTenantRepository) TenantActivityTrend(ctx context.Context, days int) ([]TenantActivityPoint, error) {
	if days < 1 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	type row struct {
		Day        string `gorm:"column:day"`
		NewTenants int64  `gorm:"column:new_tenants"`
		Logins     int64  `gorm:"column:logins"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		WITH days AS (
			SELECT generate_series(
				(CURRENT_DATE - (?::int - 1) * INTERVAL '1 day')::date,
				CURRENT_DATE::date,
				INTERVAL '1 day'
			)::date AS day
		)
		SELECT
			to_char(d.day, 'MM-DD') AS day,
			COALESCE(nt.c, 0) AS new_tenants,
			COALESCE(lg.c, 0) AS logins
		FROM days d
		LEFT JOIN (
			SELECT created_at::date AS day, COUNT(*)::bigint AS c
			FROM tenants
			WHERE created_at::date >= (CURRENT_DATE - (?::int - 1) * INTERVAL '1 day')::date
			GROUP BY 1
		) nt ON nt.day = d.day
		LEFT JOIN (
			SELECT created_at::date AS day, COUNT(*)::bigint AS c
			FROM audit_logs
			WHERE action = 'auth.login'
			  AND created_at::date >= (CURRENT_DATE - (?::int - 1) * INTERVAL '1 day')::date
			GROUP BY 1
		) lg ON lg.day = d.day
		ORDER BY d.day
	`, days, days, days).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]TenantActivityPoint, 0, len(rows))
	for _, r := range rows {
		out = append(out, TenantActivityPoint{
			Date: r.Day, NewTenants: r.NewTenants, Logins: r.Logins,
		})
	}
	return out, nil
}
