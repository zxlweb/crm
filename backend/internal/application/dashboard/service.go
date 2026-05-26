package dashboard

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"crm-backend/internal/application/appscope"
	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/crm"
	"crm-backend/internal/pkg/datascope"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
)

var ErrTeamRankingDenied = errors.New("dashboard_scope_denied")

type Service struct {
	leads     repository.LeadRepository
	accounts  repository.AccountRepository
	deals     repository.DealRepository
	activities repository.ActivityRepository
	tenants   repository.TenantRepository
	users     repository.UserRepository
	enforcer  *casbin.Enforcer
	scope     appscope.Provider
}

func NewService(
	leads repository.LeadRepository,
	accounts repository.AccountRepository,
	deals repository.DealRepository,
	activities repository.ActivityRepository,
	tenants repository.TenantRepository,
	users repository.UserRepository,
	enforcer *casbin.Enforcer,
	scope appscope.Provider,
) *Service {
	return &Service{
		leads: leads, accounts: accounts, deals: deals, activities: activities,
		tenants: tenants, users: users, enforcer: enforcer, scope: scope,
	}
}

func (s *Service) dataScope(ctx context.Context, tenantID, userID uuid.UUID) datascope.ScopeParams {
	return s.scope.Params(ctx, tenantID, userID)
}

type KPIsDTO struct {
	LeadsTotal      int64   `json:"leads_total"`
	AccountsTotal   int64   `json:"accounts_total"`
	DealsTotal      int64   `json:"deals_total"`
	DealsOpenCount  int64   `json:"deals_open_count"`
	DealsOpenAmount float64 `json:"deals_open_amount"`
	AtRiskTotal     int64   `json:"at_risk_total"`
	AvgEngagement   int     `json:"avg_engagement"`
	WeeklyFollowUps int64   `json:"weekly_follow_ups"`
}

type KPITrendsDTO struct {
	LeadsWeeklyTouch    int64  `json:"leads_weekly_touch"`
	AccountsWeeklyTouch int64  `json:"accounts_weekly_touch"`
	DealsWeeklyNew      int64  `json:"deals_weekly_new"`
	EngagementDelta     int    `json:"engagement_delta"`
	EngagementDirection string `json:"engagement_direction"`
}

type SparklinesDTO struct {
	Leads []int64 `json:"leads"`
	Deals []int64 `json:"deals"`
}

type PriorityDTO struct {
	EntityType       string   `json:"entity_type"`
	EntityID         uuid.UUID `json:"entity_id"`
	Title            string   `json:"title"`
	Reasons          []string `json:"reasons"`
	Suggestion       string   `json:"suggestion"`
	Score            int      `json:"score"`
	EngagementScore  int      `json:"engagement_score"`
	IsPreview        bool     `json:"is_preview"`
}

type SummaryDTO struct {
	DataScope            string         `json:"data_scope"`
	CanViewTeamRanking   bool           `json:"can_view_team_ranking"`
	KPIs                 KPIsDTO        `json:"kpis"`
	KPITrends  KPITrendsDTO   `json:"kpi_trends"`
	Sparklines SparklinesDTO  `json:"sparklines"`
	Priorities []PriorityDTO  `json:"priorities"`
}

type FunnelStageDTO struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type FunnelDTO struct {
	Stages []FunnelStageDTO `json:"stages"`
}

type QuotaDTO struct {
	TargetAmount     float64 `json:"target_amount"`
	WonAmountMTD     float64 `json:"won_amount_mtd"`
	CompletionRate   float64 `json:"completion_rate"`
	Period           string  `json:"period"`
}

type TeamRankingItemDTO struct {
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Department string     `json:"department,omitempty"`
	Name       string     `json:"name"`
	Value      float64    `json:"value"`
	Rank       int        `json:"rank"`
}

type TeamRankingDTO struct {
	GroupBy string               `json:"group_by"`
	Items   []TeamRankingItemDTO `json:"items"`
}

type TodoItemDTO struct {
	ID             string    `json:"id"`
	Time           string    `json:"time"`
	Title          string    `json:"title"`
	Subtitle       string    `json:"subtitle"`
	HrefEntityType string    `json:"href_entity_type"`
	HrefEntityID   uuid.UUID `json:"href_entity_id"`
}

type TodoDTO struct {
	Items []TodoItemDTO `json:"items"`
}

func (s *Service) Summary(ctx context.Context, tenantID, userID uuid.UUID, preview bool) (*SummaryDTO, error) {
	scope := s.dataScope(ctx, tenantID, userID)
	opts := s.segmentOpts(ctx, tenantID)

	leadsTotal, _ := s.leads.CountScoped(ctx, tenantID, scope)
	accountsTotal, _ := s.accounts.CountScoped(ctx, tenantID, scope)
	dealsTotal, _ := s.deals.CountScoped(ctx, tenantID, scope)

	_, pipeSummary, _ := s.deals.Pipeline(ctx, tenantID, repository.DealPipelineFilter{
		Scope: scope, PerStage: 1,
	})

	leadAtRisk, _ := s.leads.CountLowEngagement(ctx, tenantID, scope)
	accAtRisk, _ := s.accounts.CountLowEngagement(ctx, tenantID, scope)

	avgLead, _ := s.leads.AvgEngagement(ctx, tenantID, scope)
	avgEng := int(avgLead)

	weekStart := time.Now().UTC().AddDate(0, 0, -7)
	var weeklyFollow int64
	if scope.Level == datascope.LevelAll {
		weeklyFollow, _ = s.activities.CountSince(ctx, tenantID, weekStart)
	} else {
		leadTouches, _ := s.activities.CountLeadTouchesSince(ctx, tenantID, weekStart, scope)
		accTouches, _ := s.activities.CountAccountTouchesSince(ctx, tenantID, weekStart, scope)
		weeklyFollow = leadTouches + accTouches
	}

	leadSpark, _ := s.leads.DailyCreatedCounts(ctx, tenantID, scope, 7)
	dealSpark, _ := s.deals.DailyCreatedCounts(ctx, tenantID, scope, 7)
	if leadSpark == nil {
		leadSpark = make([]int64, 7)
	}
	if dealSpark == nil {
		dealSpark = make([]int64, 7)
	}

	leadTouches, _ := s.activities.CountLeadTouchesSince(ctx, tenantID, weekStart, scope)
	accTouches, _ := s.activities.CountAccountTouchesSince(ctx, tenantID, weekStart, scope)
	dealsNew := int64(0)
	if len(dealSpark) > 0 {
		for _, c := range dealSpark {
			dealsNew += c
		}
	}

	priorities := s.buildPriorities(ctx, tenantID, userID, scope, opts, preview)

	_ = opts
	return &SummaryDTO{
		DataScope:          scope.APIScope(),
		CanViewTeamRanking: s.scope.CanViewTeamRanking(ctx, tenantID, userID, scope),
		KPIs: KPIsDTO{
			LeadsTotal: leadsTotal, AccountsTotal: accountsTotal, DealsTotal: dealsTotal,
			DealsOpenCount: pipeSummary.OpenCount, DealsOpenAmount: pipeSummary.OpenAmount,
			AtRiskTotal: leadAtRisk + accAtRisk, AvgEngagement: avgEng, WeeklyFollowUps: weeklyFollow,
		},
		KPITrends: KPITrendsDTO{
			LeadsWeeklyTouch: leadTouches, AccountsWeeklyTouch: accTouches, DealsWeeklyNew: dealsNew,
			EngagementDelta: 0, EngagementDirection: "flat",
		},
		Sparklines: SparklinesDTO{Leads: leadSpark, Deals: dealSpark},
		Priorities: priorities,
	}, nil
}

func (s *Service) Funnel(ctx context.Context, tenantID, userID uuid.UUID, scopeParam string) (*FunnelDTO, error) {
	scope := s.dataScope(ctx, tenantID, userID)
	if scopeParam == "leads" {
		rows, err := s.leads.StatsFunnel(ctx, tenantID, repository.LeadStatsFilter{Scope: scope})
		if err != nil {
			return nil, err
		}
		stages := make([]FunnelStageDTO, len(rows))
		for i, row := range rows {
			stages[i] = FunnelStageDTO{Name: row.Label, Count: row.Count}
		}
		return &FunnelDTO{Stages: stages}, nil
	}
	rows, err := s.deals.CountByStage(ctx, tenantID, scope)
	if err != nil {
		return nil, err
	}
	stages := make([]FunnelStageDTO, len(rows))
	for i, row := range rows {
		stages[i] = FunnelStageDTO{Name: row.Label, Count: row.Count}
	}
	return &FunnelDTO{Stages: stages}, nil
}

func (s *Service) Quota(ctx context.Context, tenantID, userID uuid.UUID) (*QuotaDTO, error) {
	scope := s.dataScope(ctx, tenantID, userID)
	target := 0.0
	period := time.Now().UTC().Format("2006-01")
	if s.tenants != nil {
		if t, err := s.tenants.FindByID(ctx, tenantID); err == nil && t != nil {
			cfg := crm.ParseTenantCRMConfig(t.Config)
			target = cfg.SalesQuota.Amount
			if cfg.SalesQuota.Period != "" {
				period = cfg.SalesQuota.Period
			}
		}
	}
	_, summary, err := s.deals.Pipeline(ctx, tenantID, repository.DealPipelineFilter{
		Scope: scope, PerStage: 1,
	})
	if err != nil {
		return nil, err
	}
	won := summary.WonAmountMTD
	rate := 0.0
	if target > 0 {
		rate = won / target
	}
	return &QuotaDTO{
		TargetAmount: target, WonAmountMTD: won, CompletionRate: rate, Period: period,
	}, nil
}

func (s *Service) TeamRanking(ctx context.Context, tenantID, userID uuid.UUID, metric string, limit int) (*TeamRankingDTO, error) {
	scope := s.dataScope(ctx, tenantID, userID)
	if !s.scope.CanViewTeamRanking(ctx, tenantID, userID, scope) {
		return nil, ErrTeamRankingDenied
	}
	if metric == "" {
		metric = "won_amount"
	}
	if scope.Level == datascope.LevelAll {
		rows, err := s.deals.TeamRankingByDepartment(ctx, tenantID, metric, limit)
		if err != nil {
			return nil, err
		}
		items := make([]TeamRankingItemDTO, 0, len(rows))
		for i, row := range rows {
			items = append(items, TeamRankingItemDTO{
				Department: row.Department,
				Name:       row.Department,
				Value:      row.Value,
				Rank:       i + 1,
			})
		}
		return &TeamRankingDTO{GroupBy: "department", Items: items}, nil
	}
	rows, err := s.deals.TeamRanking(ctx, tenantID, metric, limit, scope)
	if err != nil {
		return nil, err
	}
	items := make([]TeamRankingItemDTO, 0, len(rows))
	for i, row := range rows {
		name := row.OwnerID.String()
		if s.users != nil {
			if u, err := s.users.FindByID(ctx, row.OwnerID); err == nil && u != nil && u.Name != "" {
				name = u.Name
			}
		}
		uid := row.OwnerID
		items = append(items, TeamRankingItemDTO{
			UserID: &uid, Name: name, Value: row.Value, Rank: i + 1,
		})
	}
	return &TeamRankingDTO{GroupBy: "user", Items: items}, nil
}

func (s *Service) Todo(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ string) (*TodoDTO, error) {
	return &TodoDTO{Items: []TodoItemDTO{}}, nil
}

func (s *Service) buildPriorities(ctx context.Context, tenantID, userID uuid.UUID, scope datascope.ScopeParams, opts crm.SegmentApplyOpts, preview bool) []PriorityDTO {
	_ = preview
	candidates, err := s.leads.ListPriorityCandidates(ctx, tenantID, scope, 20)
	if err != nil || len(candidates) == 0 {
		return []PriorityDTO{}
	}
	type scored struct {
		lead  domain.Lead
		score int
	}
	var rows []scored
	daysSilent := opts.DaysSilent
	for _, l := range candidates {
		days := crm.DaysSince(l.LastActivityAt)
		score := int(l.EngagementScore)
		if days > daysSilent {
			score = 30
		}
		rows = append(rows, scored{lead: l, score: score})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].score < rows[j].score })
	if len(rows) > 5 {
		rows = rows[:5]
	}
	out := make([]PriorityDTO, 0, len(rows))
	for _, row := range rows {
		l := row.lead
		reasons := []string{}
		suggestion := "????????"
		if crm.DaysSince(l.LastActivityAt) > daysSilent {
			reasons = append(reasons, fmt.Sprintf("%d ????", daysSilent))
			suggestion = "????????"
		}
		if crm.RelationshipHealthFromScore(l.EngagementScore) == "low" {
			reasons = append(reasons, "?????")
		}
		if len(reasons) == 0 {
			reasons = append(reasons, "???")
		}
		out = append(out, PriorityDTO{
			EntityType: "lead", EntityID: l.ID, Title: l.Title,
			Reasons: reasons, Suggestion: suggestion, Score: row.score,
			EngagementScore: int(l.EngagementScore), IsPreview: false,
		})
	}
	return out
}

func (s *Service) segmentOpts(ctx context.Context, tenantID uuid.UUID) crm.SegmentApplyOpts {
	opts := crm.SegmentApplyOpts{DaysSilent: 7, HighValueAmount: 100000}
	if s.tenants == nil {
		return opts
	}
	t, err := s.tenants.FindByID(ctx, tenantID)
	if err != nil || t == nil {
		return opts
	}
	cfg := crm.ParseTenantCRMConfig(t.Config)
	opts.DaysSilent = cfg.InsightThresholds.DaysSilent
	opts.HighValueAmount = cfg.InsightThresholds.HighValueAmount
	return opts
}
