import type {
  DashboardFunnel,
  DashboardQuota,
  DashboardSummary,
  DashboardTeamRanking,
} from '~/types/dashboard-stats'
import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'

export const MOCK_DASHBOARD_SUMMARY: DashboardSummary = {
  data_scope: 'self',
  kpis: {
    leads_total: 5,
    accounts_total: 2,
    deals_total: 6,
    deals_open_count: 4,
    deals_open_amount: 1760000,
    at_risk_total: 3,
    avg_engagement: 62,
    weekly_follow_ups: 12,
  },
  kpi_trends: {
    leadsWeeklyTouch: 3,
    accountsWeeklyTouch: 2,
    dealsWeeklyNew: 1,
    engagementDelta: 15,
    engagementDirection: 'up',
  },
  sparklines: {
    leads: [2, 3, 1, 4, 3, 5, 4],
    deals: [0, 1, 0, 2, 1, 1, 2],
  },
  priorities: [
    {
      entity_type: 'lead',
      entity_id: DEMO_LEAD_ID,
      title: '华创科技',
      reasons: ['7 天未跟进'],
      suggestion: '今日电话确认方案',
      score: 38,
      engagement_score: 62,
      is_preview: true,
    },
  ],
}

export const MOCK_DASHBOARD_FUNNEL: DashboardFunnel = {
  stages: [
    { name: 'qualification', count: 4 },
    { name: 'proposal', count: 2 },
    { name: 'negotiation', count: 1 },
    { name: 'won', count: 1 },
  ],
}

export const MOCK_DASHBOARD_QUOTA: DashboardQuota = {
  target_amount: 5000000,
  won_amount_mtd: 3400000,
  completion_rate: 0.68,
  period: '2026-05',
}

export const MOCK_DASHBOARD_TEAM_RANKING: DashboardTeamRanking = {
  items: [
    { user_id: 'user-demo-001', name: '张三', value: 1200000, rank: 1 },
    { user_id: 'user-demo-002', name: '李四', value: 980000, rank: 2 },
    { user_id: 'user-demo-003', name: '王磊', value: 760000, rank: 3 },
    { user_id: 'user-demo-004', name: '赵敏', value: 540000, rank: 4 },
  ],
}

export function mockDashboardSummary(): DashboardSummary {
  return structuredClone(MOCK_DASHBOARD_SUMMARY)
}

export function mockDashboardFunnel(): DashboardFunnel {
  return structuredClone(MOCK_DASHBOARD_FUNNEL)
}

export function mockDashboardQuota(): DashboardQuota {
  return structuredClone(MOCK_DASHBOARD_QUOTA)
}

export function mockDashboardTeamRanking(): DashboardTeamRanking {
  return structuredClone(MOCK_DASHBOARD_TEAM_RANKING)
}
