import {
  mockDashboardFunnel,
  mockDashboardQuota,
  mockDashboardSummary,
  mockDashboardTeamRanking,
} from '~/fixtures/dashboard-stats.mock'
import type {
  DashboardFunnel,
  DashboardQuota,
  DashboardSummary,
  DashboardTeamRanking,
} from '~/types/dashboard-stats'
import type { DashboardKpiTrends } from '~/types/dashboard'

type ApiKpiTrends = {
  leads_weekly_touch?: number
  accounts_weekly_touch?: number
  deals_weekly_new?: number
  engagement_delta?: number
  engagement_direction?: string
}

type ApiDashboardSummary = Omit<DashboardSummary, 'kpi_trends'> & {
  kpi_trends?: ApiKpiTrends
}

function mapKpiTrends(raw?: ApiKpiTrends): DashboardKpiTrends {
  const dir = raw?.engagement_direction
  const direction: DashboardKpiTrends['engagementDirection'] =
    dir === 'up' || dir === 'down' || dir === 'flat' ? dir : 'flat'
  return {
    leadsWeeklyTouch: raw?.leads_weekly_touch ?? 0,
    accountsWeeklyTouch: raw?.accounts_weekly_touch ?? 0,
    dealsWeeklyNew: raw?.deals_weekly_new ?? 0,
    engagementDelta: raw?.engagement_delta ?? 0,
    engagementDirection: direction,
  }
}

function normalizeSummary(raw: ApiDashboardSummary): DashboardSummary {
  return {
    data_scope: raw.data_scope,
    can_view_team_ranking: raw.can_view_team_ranking,
    kpis: raw.kpis,
    kpi_trends: mapKpiTrends(raw.kpi_trends),
    sparklines: raw.sparklines ?? { leads: [], deals: [] },
    priorities: (raw.priorities ?? []).map((p) => ({
      ...p,
      entity_id: String(p.entity_id),
    })),
  }
}

export function useDashboardStats() {
  const api = useApi()
  const config = useRuntimeConfig()
  const forceMock = computed(() => config.public.useDashboardMock === true)

  async function fetchSummary(opts?: { preview?: boolean }): Promise<DashboardSummary> {
    if (forceMock.value) {
      return mockDashboardSummary()
    }
    const q = opts?.preview ? '?preview=1' : ''
    const raw = await api.request<ApiDashboardSummary>(`/api/dashboard/summary${q}`)
    return normalizeSummary(raw)
  }

  async function fetchFunnel(scope: 'deals' | 'leads' = 'deals'): Promise<DashboardFunnel> {
    if (forceMock.value) {
      return mockDashboardFunnel()
    }
    return api.request<DashboardFunnel>(`/api/dashboard/funnel?scope=${scope}`)
  }

  async function fetchQuota(): Promise<DashboardQuota> {
    if (forceMock.value) {
      return mockDashboardQuota()
    }
    return api.request<DashboardQuota>('/api/dashboard/quota')
  }

  async function fetchTeamRanking(metric = 'won_amount', limit = 10): Promise<DashboardTeamRanking> {
    if (forceMock.value) {
      return mockDashboardTeamRanking()
    }
    return api.request<DashboardTeamRanking>(
      `/api/dashboard/team-ranking?metric=${encodeURIComponent(metric)}&limit=${limit}`,
    )
  }

  /** 403 dashboard_scope_denied 时返回 null，不阻断其它图表 */
  async function fetchTeamRankingOptional(
    metric = 'won_amount',
    limit = 10,
  ): Promise<DashboardTeamRanking | null> {
    if (forceMock.value) {
      return mockDashboardTeamRanking()
    }
    try {
      return await fetchTeamRanking(metric, limit)
    } catch (e) {
      const msg = e instanceof Error ? e.message : ''
      if (msg.includes('dashboard_scope_denied') || msg.includes('403')) {
        return null
      }
      throw e
    }
  }

  return {
    forceMock,
    fetchSummary,
    fetchFunnel,
    fetchQuota,
    fetchTeamRanking,
    fetchTeamRankingOptional,
  }
}
