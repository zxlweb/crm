import { mockPlanDistribution, mockTenantHealth, mockTopTenants } from '~/fixtures/admin-insights.mock'
import type {
  AdminInsightsQuery,
  PlanDistribution,
  TenantHealthData,
  TopTenants,
} from '~/types/admin-insights'

function buildQuery(q: AdminInsightsQuery): string {
  const params = new URLSearchParams()
  if (q.from) params.set('from', q.from)
  if (q.to) params.set('to', q.to)
  if (q.metric) params.set('metric', q.metric)
  if (q.limit) params.set('limit', String(q.limit))
  const s = params.toString()
  return s ? `?${s}` : ''
}

export function useAdminTenantInsights() {
  const api = useApi()
  const config = useRuntimeConfig()
  const forceMock = computed(() => config.public.useAdminInsightsMock === true || config.public.useAdminInsightsMock === 'true')

  async function fetchTenantHealth(): Promise<TenantHealthData> {
    if (forceMock.value) return mockTenantHealth()
    return api.request<TenantHealthData>('/api/super-admin/stats/tenant-health', { skipTenant: true })
  }

  async function fetchPlanDistribution(q: AdminInsightsQuery = {}): Promise<PlanDistribution> {
    if (forceMock.value) return mockPlanDistribution()
    return api.request<PlanDistribution>(`/api/super-admin/stats/plan-distribution${buildQuery(q)}`, { skipTenant: true })
  }

  async function fetchTopTenants(q: AdminInsightsQuery = { metric: 'activity', limit: 10 }): Promise<TopTenants> {
    if (forceMock.value) return mockTopTenants()
    return api.request<TopTenants>(`/api/super-admin/stats/top-tenants${buildQuery(q)}`, { skipTenant: true })
  }

  return { forceMock, fetchTenantHealth, fetchPlanDistribution, fetchTopTenants }
}
