import {
  mockClaimLead,
  mockConvertLead,
  mockCreateLead,
  mockDeleteLead,
  mockGetLead,
  mockGetPoolSettings,
  mockListLeads,
  mockPoolStats,
  mockRecycleStaleLeads,
  mockReleaseLead,
  mockTransferLead,
  mockUpdateLead,
  mockUpdatePoolSettings,
} from '~/utils/leads-mock-store'
import type {
  Lead,
  LeadConvertInput,
  LeadCreateInput,
  LeadListQuery,
  LeadPoolSettings,
  LeadPoolStats,
  LeadRecycleSummary,
  LeadUpdateInput,
  LeadsListData,
  Pagination,
} from '~/types/lead'

type LeadsPagePayload = LeadsListData | Lead[]

function normalizeListPayload(data: LeadsPagePayload | undefined): Lead[] {
  if (!data) return []
  if (Array.isArray(data)) return data
  return data.items ?? []
}

function buildQueryString(query: LeadListQuery): string {
  const params = new URLSearchParams()
  if (query.page) params.set('page', String(query.page))
  if (query.page_size) params.set('page_size', String(query.page_size))
  if (query.status) params.set('status', query.status)
  if (query.source) params.set('source', query.source)
  if (query.owner_id) params.set('owner_id', query.owner_id)
  if (query.lifecycle_stage) params.set('lifecycle_stage', query.lifecycle_stage)
  if (query.relationship_health) params.set('relationship_health', query.relationship_health)
  if (query.segment) params.set('segment', query.segment)
  if (query.search) params.set('search', query.search)
  if (query.pool) params.set('pool', query.pool)
  const qs = params.toString()
  return qs ? `?${qs}` : ''
}

export function useLeads() {
  const config = useRuntimeConfig()
  const api = useApi()
  const tenant = useTenant()
  const auth = useAuth()

  const useMock = computed(
    () => config.public.useLeadsMock === true || config.public.useLeadsMock === 'true',
  )

  async function fetchList(
    query: LeadListQuery = {},
  ): Promise<{ data: LeadsListData; pagination: Pagination }> {
    if (useMock.value) {
      return mockListLeads(query, auth.user.value?.id ?? null)
    }
    const path = `/api/leads${buildQueryString(query)}`
    const res = await api.requestPage<LeadsPagePayload>(path)
    return {
      data: { items: normalizeListPayload(res.data) },
      pagination: {
        page: res.pagination.page,
        page_size: res.pagination.page_size,
        total: Number(res.pagination.total),
      },
    }
  }

  async function fetchById(id: string): Promise<Lead | null> {
    if (useMock.value) {
      return mockGetLead(id)
    }
    try {
      return await api.request<Lead>(`/api/leads/${id}`)
    } catch {
      return null
    }
  }

  async function create(input: LeadCreateInput): Promise<Lead> {
    if (useMock.value) {
      return mockCreateLead(
        input,
        tenant.currentTenantId.value ?? 'tenant-local',
        auth.user.value?.id ?? null,
      )
    }
    return api.request<Lead>('/api/leads', { method: 'POST', body: JSON.stringify(input) })
  }

  async function update(id: string, input: LeadUpdateInput): Promise<Lead> {
    if (useMock.value) {
      const row = mockUpdateLead(id, input)
      if (!row) throw new Error('lead_not_found')
      return row
    }
    return api.request<Lead>(`/api/leads/${id}`, { method: 'PATCH', body: JSON.stringify(input) })
  }

  async function remove(id: string): Promise<void> {
    if (useMock.value) {
      if (!mockDeleteLead(id)) throw new Error('lead_not_found')
      return
    }
    await api.request<void>(`/api/leads/${id}`, { method: 'DELETE' })
  }

  async function convert(id: string, input: LeadConvertInput): Promise<Lead> {
    if (useMock.value) {
      const row = mockConvertLead(id, input)
      if (!row) throw new Error('lead_not_found')
      return row
    }
    return api.request<Lead>(`/api/leads/${id}/convert`, {
      method: 'POST',
      body: JSON.stringify(input),
    })
  }

  // ===== 客户池（公海/私海）操作 =====

  async function claim(id: string): Promise<Lead> {
    if (useMock.value) {
      const row = mockClaimLead(id, auth.user.value?.id ?? null)
      if (!row) throw new Error('lead_not_found')
      return row
    }
    return api.request<Lead>(`/api/leads/${id}/claim`, { method: 'POST' })
  }

  async function release(id: string): Promise<Lead> {
    if (useMock.value) {
      const row = mockReleaseLead(id)
      if (!row) throw new Error('lead_not_found')
      return row
    }
    return api.request<Lead>(`/api/leads/${id}/release`, { method: 'POST' })
  }

  async function transfer(id: string, toUserId: string): Promise<Lead> {
    if (useMock.value) {
      const row = mockTransferLead(id, toUserId)
      if (!row) throw new Error('lead_not_found')
      return row
    }
    return api.request<Lead>(`/api/leads/${id}/transfer`, {
      method: 'POST',
      body: JSON.stringify({ owner_id: toUserId }),
    })
  }

  async function recycleStale(): Promise<LeadRecycleSummary> {
    if (useMock.value) {
      return mockRecycleStaleLeads()
    }
    return api.request<LeadRecycleSummary>('/api/leads/pool/recycle', { method: 'POST' })
  }

  async function poolStats(): Promise<LeadPoolStats> {
    if (useMock.value) {
      return mockPoolStats(auth.user.value?.id ?? null)
    }
    return api.request<LeadPoolStats>('/api/leads/pool/stats')
  }

  async function getPoolSettings(): Promise<LeadPoolSettings> {
    if (useMock.value) {
      return mockGetPoolSettings()
    }
    return api.request<LeadPoolSettings>('/api/leads/pool/settings')
  }

  async function updatePoolSettings(input: Partial<LeadPoolSettings>): Promise<LeadPoolSettings> {
    if (useMock.value) {
      return mockUpdatePoolSettings(input)
    }
    return api.request<LeadPoolSettings>('/api/leads/pool/settings', {
      method: 'PATCH',
      body: JSON.stringify(input),
    })
  }

  return {
    useMock,
    fetchList,
    fetchById,
    create,
    update,
    remove,
    convert,
    claim,
    release,
    transfer,
    recycleStale,
    poolStats,
    getPoolSettings,
    updatePoolSettings,
  }
}
