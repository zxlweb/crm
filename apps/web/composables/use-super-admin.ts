export type SuperAdminOverview = {
  tenant_count: number
  active_tenant_count: number
  user_count: number
}

export type TenantActivityTrend = {
  categories: string[]
  series: Array<{
    name: string
    data: number[]
    primary?: boolean
  }>
}

export type SuperAdminTenant = {
  id: string
  name: string
  domain: string
  is_active: boolean
  user_count: number
  created_at: string
}

export function useSuperAdmin() {
  const api = useApi()

  function fetchOverview() {
    return api.request<SuperAdminOverview>('/api/super-admin/overview', { skipTenant: true })
  }

  function fetchTenantActivityTrend(days = 7) {
    return api.request<TenantActivityTrend>(
      `/api/super-admin/stats/tenant-activity?days=${days}`,
      { skipTenant: true },
    )
  }

  function fetchTenants(params: { page?: number; page_size?: number; search?: string } = {}) {
    const q = new URLSearchParams()
    if (params.page) q.set('page', String(params.page))
    if (params.page_size) q.set('page_size', String(params.page_size))
    if (params.search) q.set('search', params.search)
    const query = q.toString()
    return api.requestPage<{ items: SuperAdminTenant[] }>(
      `/api/super-admin/tenants${query ? `?${query}` : ''}`,
      { skipTenant: true },
    )
  }

  function patchTenantActive(id: string, isActive: boolean) {
    return api.request<SuperAdminTenant>(`/api/super-admin/tenants/${id}`, {
      method: 'PATCH',
      skipTenant: true,
      body: JSON.stringify({ is_active: isActive }),
    })
  }

  return { fetchOverview, fetchTenantActivityTrend, fetchTenants, patchTenantActive }
}
