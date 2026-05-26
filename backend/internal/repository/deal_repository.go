package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/infrastructure/persistence"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/datascope"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrDealNotFound = errors.New("deal not found")

type DealListFilter struct {
	Page               int
	PageSize           int
	Search             string
	Stage              string
	Stages             []string
	OwnerID            *uuid.UUID
	AccountID          *uuid.UUID
	LeadID             *uuid.UUID
	ExpectedCloseFrom  *time.Time
	ExpectedCloseTo    *time.Time
	MinAmount          *float64
	MaxAmount          *float64
	Scope              datascope.ScopeParams
}

type DealPipelineFilter struct {
	OwnerID   *uuid.UUID
	AccountID *uuid.UUID
	Scope     datascope.ScopeParams
	PerStage  int
}

type DealPipelineSummary struct {
	OpenCount    int64
	OpenAmount   float64
	WonCountMTD  int64
	WonAmountMTD float64
}

type DealRepository interface {
	List(ctx context.Context, tenantID uuid.UUID, f DealListFilter) ([]domain.Deal, int64, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID, scope datascope.ScopeParams) (*domain.Deal, error)
	Create(ctx context.Context, d *domain.Deal) error
	Update(ctx context.Context, d *domain.Deal) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	Pipeline(ctx context.Context, tenantID uuid.UUID, f DealPipelineFilter) (map[string][]domain.Deal, DealPipelineSummary, error)
	CountScoped(ctx context.Context, tenantID uuid.UUID, scope datascope.ScopeParams) (int64, error)
	StatsByStage(ctx context.Context, tenantID uuid.UUID, f DealStatsFilter, metric string) ([]DealStageStat, int64, error)
	StatsWinRate(ctx context.Context, tenantID uuid.UUID, f DealStatsFilter, granularity string) ([]DealWinRatePoint, error)
	DailyCreatedCounts(ctx context.Context, tenantID uuid.UUID, scope datascope.ScopeParams, days int) ([]int64, error)
	CountByStage(ctx context.Context, tenantID uuid.UUID, scope datascope.ScopeParams) ([]LabelCount, error)
	TeamRanking(ctx context.Context, tenantID uuid.UUID, metric string, limit int, scope datascope.ScopeParams) ([]DealOwnerMetric, error)
	TeamRankingByDepartment(ctx context.Context, tenantID uuid.UUID, metric string, limit int) ([]DealDepartmentMetric, error)
}

type GormDealRepository struct {
	db *gorm.DB
}

func NewDealRepository(db *gorm.DB) DealRepository {
	return &GormDealRepository{db: db}
}

func (r *GormDealRepository) base(ctx context.Context, tenantID uuid.UUID) *gorm.DB {
	return persistence.DBFromContext(r.db, ctx).Model(&domain.Deal{}).Where("tenant_id = ?", tenantID)
}

func (r *GormDealRepository) scoped(ctx context.Context, tenantID uuid.UUID, scope datascope.ScopeParams) *gorm.DB {
	if scope.TenantID == uuid.Nil {
		scope.TenantID = tenantID
	}
	return datascope.ApplyOwnerScope(r.base(ctx, tenantID), scope)
}

func (r *GormDealRepository) applyListFilters(q *gorm.DB, f DealListFilter) *gorm.DB {
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("title ILIKE ?", like)
	}
	if f.Stage != "" {
		q = q.Where("stage = ?", f.Stage)
	}
	if len(f.Stages) > 0 {
		q = q.Where("stage IN ?", f.Stages)
	}
	if f.OwnerID != nil {
		q = q.Where("owner_id = ?", *f.OwnerID)
	}
	if f.AccountID != nil {
		q = q.Where("account_id = ?", *f.AccountID)
	}
	if f.LeadID != nil {
		q = q.Where("lead_id = ?", *f.LeadID)
	}
	if f.ExpectedCloseFrom != nil {
		q = q.Where("expected_close_date >= ?", *f.ExpectedCloseFrom)
	}
	if f.ExpectedCloseTo != nil {
		q = q.Where("expected_close_date <= ?", *f.ExpectedCloseTo)
	}
	if f.MinAmount != nil {
		q = q.Where("amount >= ?", *f.MinAmount)
	}
	if f.MaxAmount != nil {
		q = q.Where("amount <= ?", *f.MaxAmount)
	}
	return q
}

func (r *GormDealRepository) List(ctx context.Context, tenantID uuid.UUID, f DealListFilter) ([]domain.Deal, int64, error) {
	q := r.scoped(ctx, tenantID, f.Scope)
	q = r.applyListFilters(q, f)

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

	var items []domain.Deal
	err := q.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *GormDealRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID, scope datascope.ScopeParams) (*domain.Deal, error) {
	q := r.scoped(ctx, tenantID, scope).Where("id = ?", id)
	var d domain.Deal
	err := q.First(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDealNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *GormDealRepository) Create(ctx context.Context, d *domain.Deal) error {
	if d.Stage == "" {
		d.Stage = crm.DealStageQualification
	}
	if !crm.ValidDealStage(d.Stage) {
		return errors.New("invalid stage")
	}
	if d.Currency == "" {
		d.Currency = "CNY"
	}
	if !crm.ValidDealCurrency(d.Currency) {
		return errors.New("invalid currency")
	}
	return persistence.DBFromContext(r.db, ctx).Create(d).Error
}

func (r *GormDealRepository) Update(ctx context.Context, d *domain.Deal) error {
	if d.Stage != "" && !crm.ValidDealStage(d.Stage) {
		return errors.New("invalid stage")
	}
	if d.Currency != "" && !crm.ValidDealCurrency(d.Currency) {
		return errors.New("invalid currency")
	}
	return persistence.DBFromContext(r.db, ctx).Save(d).Error
}

func (r *GormDealRepository) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	res := persistence.DBFromContext(r.db, ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&domain.Deal{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrDealNotFound
	}
	return nil
}

func (r *GormDealRepository) Pipeline(ctx context.Context, tenantID uuid.UUID, f DealPipelineFilter) (map[string][]domain.Deal, DealPipelineSummary, error) {
	perStage := f.PerStage
	if perStage < 1 {
		perStage = 20
	}

	q := r.scoped(ctx, tenantID, f.Scope)
	if f.OwnerID != nil {
		q = q.Where("owner_id = ?", *f.OwnerID)
	}
	if f.AccountID != nil {
		q = q.Where("account_id = ?", *f.AccountID)
	}

	byStage := make(map[string][]domain.Deal, len(crm.DealPipelineStages))
	for _, stage := range crm.DealPipelineStages {
		var items []domain.Deal
		sq := q.Session(&gorm.Session{}).Where("stage = ?", stage).Order("updated_at DESC").Limit(perStage)
		if err := sq.Find(&items).Error; err != nil {
			return nil, DealPipelineSummary{}, err
		}
		byStage[stage] = items
	}

	var summary DealPipelineSummary
	openQ := q.Session(&gorm.Session{}).Where("stage IN ?", []string{
		crm.DealStageQualification, crm.DealStageProposal, crm.DealStageNegotiation,
	})
	if err := openQ.Count(&summary.OpenCount).Error; err != nil {
		return nil, DealPipelineSummary{}, err
	}
	if err := openQ.Select("COALESCE(SUM(amount), 0)").Scan(&summary.OpenAmount).Error; err != nil {
		return nil, DealPipelineSummary{}, err
	}

	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	wonQ := q.Session(&gorm.Session{}).Where("stage = ? AND closed_at >= ?", crm.DealStageWon, monthStart)
	if err := wonQ.Count(&summary.WonCountMTD).Error; err != nil {
		return nil, DealPipelineSummary{}, err
	}
	if err := wonQ.Select("COALESCE(SUM(amount), 0)").Scan(&summary.WonAmountMTD).Error; err != nil {
		return nil, DealPipelineSummary{}, err
	}

	return byStage, summary, nil
}

// ParseDealStages splits comma-separated stage query values.
func ParseDealStages(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" && crm.ValidDealStage(p) {
			out = append(out, p)
		}
	}
	return out
}
