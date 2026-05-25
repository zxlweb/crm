import {
  mockCreateDeal,
  mockFetchPipeline,
  mockGetDeal,
  mockListDeals,
  mockUpdateDeal,
  mockUpdateDealStage,
} from '~/utils/deals-mock-store'
import type {
  Deal,
  DealCreateInput,
  DealListQuery,
  DealPipelineData,
  DealStageUpdateInput,
  DealUpdateInput,
  DealsListData,
  Pagination,
} from '~/types/deal'

type DealsPagePayload = DealsListData | Deal[]

function normalizeListPayload(data: DealsPagePayload | undefined): Deal[] {
  if (!data) return []
  if (Array.isArray(data)) return data
  return data.items ?? []
}

function buildQueryString(query: DealListQuery): string {
  const params = new URLSearchParams()
  if (query.page) params.set('page', String(query.page))
  if (query.page_size) params.set('page_size', String(query.page_size))
  if (query.search) params.set('search', query.search)
  if (query.stage) params.set('stage', query.stage)
  if (query.stages) params.set('stages', query.stages)
  if (query.owner_id) params.set('owner_id', query.owner_id)
  if (query.account_id) params.set('account_id', query.account_id)
  if (query.lead_id) params.set('lead_id', query.lead_id)
  const qs = params.toString()
  return qs ? `?${qs}` : ''
}

export function useDeals() {
  const config = useRuntimeConfig()
  const api = useApi()
  const tenant = useTenant()
  const auth = useAuth()

  const useMock = computed(
    () => config.public.useDealsMock === true || config.public.useDealsMock === 'true',
  )

  async function fetchList(
    query: DealListQuery = {},
  ): Promise<{ data: DealsListData; pagination: Pagination }> {
    if (useMock.value) {
      return mockListDeals(query)
    }
    const path = `/api/deals${buildQueryString(query)}`
    const res = await api.requestPage<DealsPagePayload>(path)
    return {
      data: { items: normalizeListPayload(res.data) },
      pagination: {
        page: res.pagination.page,
        page_size: res.pagination.page_size,
        total: Number(res.pagination.total),
      },
    }
  }

  async function fetchById(id: string): Promise<Deal | null> {
    if (useMock.value) {
      return mockGetDeal(id)
    }
    try {
      return await api.request<Deal>(`/api/deals/${id}`)
    } catch {
      return null
    }
  }

  async function fetchPipeline(): Promise<DealPipelineData> {
    if (useMock.value) {
      return mockFetchPipeline()
    }
    return api.request<DealPipelineData>('/api/deals/pipeline')
  }

  async function create(input: DealCreateInput): Promise<Deal> {
    if (useMock.value) {
      return mockCreateDeal(input, tenant.currentTenantId.value ?? 'demo', auth.user.value?.id ?? null)
    }
    return api.request<Deal>('/api/deals', { method: 'POST', body: JSON.stringify(input) })
  }

  async function update(id: string, input: DealUpdateInput): Promise<Deal> {
    if (useMock.value) {
      const row = mockUpdateDeal(id, input)
      if (!row) throw new Error('Deal not found')
      return row
    }
    return api.request<Deal>(`/api/deals/${id}`, { method: 'PUT', body: JSON.stringify(input) })
  }

  async function updateStage(id: string, input: DealStageUpdateInput): Promise<Deal> {
    if (useMock.value) {
      const row = mockUpdateDealStage(id, input)
      if (!row) throw new Error('Deal not found')
      return row
    }
    return api.request<Deal>(`/api/deals/${id}/stage`, {
      method: 'PUT',
      body: JSON.stringify(input),
    })
  }

  async function remove(id: string): Promise<void> {
    if (useMock.value) {
      mockUpdateDeal(id, { stage: 'lost', lost_reason: 'deleted' })
      return
    }
    await api.request<void>(`/api/deals/${id}`, { method: 'DELETE' })
  }

  return {
    useMock,
    fetchList,
    fetchById,
    fetchPipeline,
    create,
    update,
    updateStage,
    remove,
  }
}
