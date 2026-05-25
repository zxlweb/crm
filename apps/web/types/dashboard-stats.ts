/** Phase 3 — docs/api/phase-3-deals-dashboard-api.md §4 */

import type { DashboardKpiTrends } from '~/types/dashboard'

export type DashboardDataScope = 'self' | 'department' | 'all'

export type DashboardSummaryKpis = {
  leads_total: number
  accounts_total: number
  deals_total: number
  deals_open_count: number
  deals_open_amount: number
  at_risk_total: number
  avg_engagement: number
  weekly_follow_ups: number
}

export type DashboardSummaryPriority = {
  entity_type: 'lead' | 'account'
  entity_id: string
  title: string
  reasons: string[]
  suggestion: string
  score: number
  engagement_score: number
  is_preview?: boolean
}

export type DashboardSummary = {
  data_scope: DashboardDataScope
  kpis: DashboardSummaryKpis
  kpi_trends: DashboardKpiTrends & { deals_weekly_new?: number }
  sparklines: {
    leads: number[]
    deals: number[]
  }
  priorities: DashboardSummaryPriority[]
}

export type DashboardFunnelStage = {
  name: string
  count: number
}

export type DashboardFunnel = {
  stages: DashboardFunnelStage[]
}

export type DashboardQuota = {
  target_amount: number
  won_amount_mtd: number
  completion_rate: number
  period: string
}

export type DashboardTeamRankingItem = {
  user_id: string
  name: string
  value: number
  rank: number
}

export type DashboardTeamRanking = {
  items: DashboardTeamRankingItem[]
}
