import { MOCK_DEALS_SEED } from '~/fixtures/deals.mock'
import {
  DEAL_PIPELINE_STAGES,
  type Deal,
  type DealCreateInput,
  type DealListQuery,
  type DealPipelineData,
  type DealPipelineItem,
  type DealStage,
  type DealStageUpdateInput,
  type DealUpdateInput,
  type DealsListData,
  type Pagination,
} from '~/types/deal'

function cloneSeed(): Deal[] {
  return MOCK_DEALS_SEED.map((row) => ({ ...row, tags: [...row.tags] }))
}

let store: Deal[] | null = null

function getStore(): Deal[] {
  if (!store) store = cloneSeed()
  return store
}

export function resetDealsMockStore(): void {
  store = null
}

function matchesQuery(deal: Deal, query: DealListQuery): boolean {
  if (query.stage && deal.stage !== query.stage) return false
  if (query.stages) {
    const allowed = query.stages.split(',').map((s) => s.trim())
    if (!allowed.includes(deal.stage)) return false
  }
  if (query.owner_id && deal.owner_id !== query.owner_id) return false
  if (query.account_id && deal.account_id !== query.account_id) return false
  if (query.lead_id && deal.lead_id !== query.lead_id) return false
  if (query.search) {
    const q = query.search.toLowerCase()
    if (!deal.title.toLowerCase().includes(q)) return false
  }
  return true
}

function toPipelineItem(deal: Deal): DealPipelineItem {
  return {
    id: deal.id,
    title: deal.title,
    amount: deal.amount,
    currency: deal.currency,
    probability: deal.probability,
    expected_close_date: deal.expected_close_date,
    account_id: deal.account_id,
    owner_id: deal.owner_id,
  }
}

export function mockListDeals(query: DealListQuery = {}): { data: DealsListData; pagination: Pagination } {
  const page = query.page ?? 1
  const pageSize = Math.min(query.page_size ?? 20, 100)
  const filtered = getStore()
    .filter((row) => matchesQuery(row, query))
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())

  const start = (page - 1) * pageSize
  const items = filtered.slice(start, start + pageSize)

  return {
    data: { items },
    pagination: { page, page_size: pageSize, total: filtered.length },
  }
}

export function mockGetDeal(id: string): Deal | null {
  return getStore().find((row) => row.id === id) ?? null
}

export function mockFetchPipeline(): DealPipelineData {
  const rows = getStore()
  const stages = DEAL_PIPELINE_STAGES.map((stage) => {
    const items = rows
      .filter((d) => d.stage === stage)
      .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
      .slice(0, 20)
      .map(toPipelineItem)
    const amount_total = items.reduce((sum, item) => sum + item.amount, 0)
    return { stage, count: items.length, amount_total, items }
  })

  const open = rows.filter((d) => d.stage !== 'won' && d.stage !== 'lost')
  const now = new Date()
  const monthStart = new Date(now.getFullYear(), now.getMonth(), 1)
  const wonMtd = rows.filter(
    (d) => d.stage === 'won' && d.closed_at && new Date(d.closed_at) >= monthStart,
  )

  return {
    stages,
    summary: {
      open_count: open.length,
      open_amount: open.reduce((sum, d) => sum + d.amount, 0),
      won_count_mtd: wonMtd.length,
      won_amount_mtd: wonMtd.reduce((sum, d) => sum + d.amount, 0),
    },
  }
}

export function mockCreateDeal(input: DealCreateInput, tenantId: string, ownerId: string | null): Deal {
  const ts = new Date().toISOString()
  const deal: Deal = {
    id: `deal-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`,
    tenant_id: tenantId,
    owner_id: input.owner_id ?? ownerId,
    title: input.title,
    stage: input.stage ?? 'qualification',
    amount: input.amount ?? 0,
    currency: input.currency ?? 'CNY',
    probability: input.probability ?? 30,
    expected_close_date: input.expected_close_date ?? null,
    account_id: input.account_id ?? null,
    lead_id: input.lead_id ?? null,
    contact_id: input.contact_id ?? null,
    description: input.description ?? '',
    tags: input.tags ?? [],
    lost_reason: null,
    closed_at: null,
    engagement_score: 0,
    last_activity_at: null,
    created_at: ts,
    updated_at: ts,
  }
  getStore().unshift(deal)
  return deal
}

export function mockUpdateDeal(id: string, input: DealUpdateInput): Deal | null {
  const idx = getStore().findIndex((row) => row.id === id)
  if (idx < 0) return null
  const current = getStore()[idx]
  const nextStage = input.stage ?? current.stage
  const updated: Deal = {
    ...current,
    ...input,
    stage: nextStage,
    tags: input.tags ?? current.tags,
    closed_at:
      nextStage === 'won' || nextStage === 'lost'
        ? current.closed_at ?? new Date().toISOString()
        : null,
    probability:
      nextStage === 'won' ? 100 : nextStage === 'lost' ? 0 : (input.probability ?? current.probability),
    updated_at: new Date().toISOString(),
  }
  getStore()[idx] = updated
  return updated
}

export function mockUpdateDealStage(id: string, input: DealStageUpdateInput): Deal | null {
  return mockUpdateDeal(id, {
    stage: input.stage,
    lost_reason: input.lost_reason ?? null,
  })
}
