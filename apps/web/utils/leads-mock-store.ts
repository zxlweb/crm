import { MOCK_LEADS_SEED } from '~/fixtures/leads.mock'
import type {
  Lead,
  LeadConvertInput,
  LeadCreateInput,
  LeadListQuery,
  LeadPool,
  LeadPoolSettings,
  LeadPoolStats,
  LeadRecycleSummary,
  LeadUpdateInput,
  LeadsListData,
  Pagination,
} from '~/types/lead'

function cloneSeed(): Lead[] {
  return MOCK_LEADS_SEED.map((row) => ({ ...row, tags: [...row.tags] }))
}

let store: Lead[] | null = null

// ---------- 客户池设置（mock：模块级状态） ----------
const DEFAULT_POOL_SETTINGS: LeadPoolSettings = {
  enabled: true,
  inactiveDays: 30,
  last_recycled_at: null,
}

let poolSettings: LeadPoolSettings = { ...DEFAULT_POOL_SETTINGS }

export function mockGetPoolSettings(): LeadPoolSettings {
  return { ...poolSettings }
}

export function mockUpdatePoolSettings(input: Partial<LeadPoolSettings>): LeadPoolSettings {
  poolSettings = {
    ...poolSettings,
    ...input,
    inactiveDays:
      input.inactiveDays != null
        ? Math.max(1, Math.min(365, Math.trunc(input.inactiveDays)))
        : poolSettings.inactiveDays,
  }
  return { ...poolSettings }
}

function getStore(): Lead[] {
  if (!store) store = cloneSeed()
  return store
}

export function resetLeadsMockStore(): void {
  store = null
  poolSettings = { ...DEFAULT_POOL_SETTINGS }
}

function poolOf(lead: Lead, currentUserId: string | null | undefined): LeadPool {
  if (!lead.owner_id) return 'public'
  if (currentUserId && lead.owner_id === currentUserId) return 'mine'
  return 'others'
}

function matchesPool(lead: Lead, pool: LeadPool | undefined, currentUserId: string | null | undefined): boolean {
  if (!pool || pool === 'all') return true
  if (pool === 'public') return !lead.owner_id
  if (pool === 'mine') return Boolean(currentUserId && lead.owner_id === currentUserId)
  if (pool === 'others') return Boolean(lead.owner_id && lead.owner_id !== currentUserId)
  return true
}

function matchesQuery(
  lead: Lead,
  query: LeadListQuery,
  currentUserId: string | null | undefined,
): boolean {
  if (!matchesPool(lead, query.pool, currentUserId)) return false
  if (query.status && lead.status !== query.status) return false
  if (query.source && lead.source !== query.source) return false
  if (query.owner_id && lead.owner_id !== query.owner_id) return false
  if (query.lifecycle_stage && lead.lifecycle_stage !== query.lifecycle_stage) return false
  if (query.relationship_health && lead.relationship_health !== query.relationship_health) return false
  if (query.search) {
    const q = query.search.toLowerCase()
    if (!lead.title.toLowerCase().includes(q) && !lead.source.toLowerCase().includes(q)) {
      return false
    }
  }
  return true
}

export function mockListLeads(
  query: LeadListQuery = {},
  currentUserId: string | null = null,
): { data: LeadsListData; pagination: Pagination } {
  const page = query.page ?? 1
  const pageSize = Math.min(query.page_size ?? 20, 100)
  const filtered = getStore()
    .filter((row) => matchesQuery(row, query, currentUserId))
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())

  const start = (page - 1) * pageSize
  const items = filtered.slice(start, start + pageSize)

  return {
    data: { items },
    pagination: { page, page_size: pageSize, total: filtered.length },
  }
}

export function mockPoolStats(currentUserId: string | null = null): LeadPoolStats {
  const all = getStore()
  let mine = 0
  let pub = 0
  let others = 0
  let recyclableSoon = 0
  const now = Date.now()
  const days = poolSettings.inactiveDays
  for (const lead of all) {
    const pool = poolOf(lead, currentUserId)
    if (pool === 'mine') mine += 1
    else if (pool === 'public') pub += 1
    else if (pool === 'others') others += 1

    if (pool !== 'public' && poolSettings.enabled) {
      const baseline = lead.last_activity_at
        ? new Date(lead.last_activity_at).getTime()
        : new Date(lead.created_at).getTime()
      const idleDays = (now - baseline) / 86_400_000
      const remaining = days - idleDays
      if (remaining <= 3) recyclableSoon += 1
    }
  }
  return { mine, public: pub, others, all: all.length, recyclableSoon }
}

export function mockClaimLead(id: string, userId: string | null): Lead | null {
  const idx = getStore().findIndex((row) => row.id === id)
  if (idx < 0) return null
  const current = getStore()[idx]
  if (current.owner_id) {
    throw new Error('lead_already_owned')
  }
  if (!userId) {
    throw new Error('claim_requires_login')
  }
  const updated: Lead = {
    ...current,
    owner_id: userId,
    last_activity_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  }
  getStore()[idx] = updated
  return updated
}

export function mockReleaseLead(id: string): Lead | null {
  const idx = getStore().findIndex((row) => row.id === id)
  if (idx < 0) return null
  const current = getStore()[idx]
  if (!current.owner_id) {
    return current
  }
  const updated: Lead = {
    ...current,
    owner_id: null,
    updated_at: new Date().toISOString(),
  }
  getStore()[idx] = updated
  return updated
}

export function mockTransferLead(id: string, toUserId: string): Lead | null {
  if (!toUserId) {
    throw new Error('transfer_requires_target')
  }
  const idx = getStore().findIndex((row) => row.id === id)
  if (idx < 0) return null
  const current = getStore()[idx]
  const updated: Lead = {
    ...current,
    owner_id: toUserId,
    updated_at: new Date().toISOString(),
  }
  getStore()[idx] = updated
  return updated
}

/**
 * 扫描私海中沉默超过 inactiveDays 的线索并回收至公海。
 * 已转化（converted）的线索不会被回收，避免误伤客户档案。
 */
export function mockRecycleStaleLeads(now: number = Date.now()): LeadRecycleSummary {
  const summary: LeadRecycleSummary = {
    recycled: 0,
    scanned: 0,
    threshold_days: poolSettings.inactiveDays,
    recycled_ids: [],
  }
  if (!poolSettings.enabled) {
    poolSettings.last_recycled_at = new Date(now).toISOString()
    return summary
  }
  const limitMs = poolSettings.inactiveDays * 86_400_000
  const list = getStore()
  for (let i = 0; i < list.length; i += 1) {
    const lead = list[i]
    if (!lead.owner_id) continue
    if (lead.status === 'converted') continue
    summary.scanned += 1
    const baseline = lead.last_activity_at
      ? new Date(lead.last_activity_at).getTime()
      : new Date(lead.created_at).getTime()
    if (now - baseline >= limitMs) {
      const tags = lead.tags.includes('已回收') ? lead.tags : [...lead.tags, '已回收']
      list[i] = {
        ...lead,
        owner_id: null,
        tags,
        updated_at: new Date(now).toISOString(),
      }
      summary.recycled += 1
      summary.recycled_ids.push(lead.id)
    }
  }
  poolSettings.last_recycled_at = new Date(now).toISOString()
  return summary
}

export function mockGetLead(id: string): Lead | null {
  return getStore().find((row) => row.id === id) ?? null
}

export function mockCreateLead(input: LeadCreateInput, tenantId: string, ownerId: string | null): Lead {
  const ts = new Date().toISOString()
  const lead: Lead = {
    id: `lead-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`,
    tenant_id: tenantId,
    owner_id: input.owner_id ?? ownerId,
    title: input.title,
    status: input.status ?? 'new',
    source: input.source ?? '',
    amount: input.amount ?? 0,
    expected_close_date: input.expected_close_date ?? null,
    lifecycle_stage: input.lifecycle_stage ?? 'acquire',
    relationship_health: 'medium',
    engagement_score: 0,
    last_activity_at: null,
    tags: input.tags ?? [],
    converted_account_id: null,
    converted_contact_id: null,
    created_at: ts,
    updated_at: ts,
  }
  getStore().unshift(lead)
  return lead
}

export function mockUpdateLead(id: string, input: LeadUpdateInput): Lead | null {
  const idx = getStore().findIndex((row) => row.id === id)
  if (idx < 0) return null
  const current = getStore()[idx]
  const updated: Lead = {
    ...current,
    ...input,
    tags: input.tags ?? current.tags,
    updated_at: new Date().toISOString(),
  }
  getStore()[idx] = updated
  return updated
}

export function mockDeleteLead(id: string): boolean {
  const before = getStore().length
  store = getStore().filter((row) => row.id !== id)
  return store.length < before
}

export function mockConvertLead(id: string, input: LeadConvertInput): Lead | null {
  const idx = getStore().findIndex((row) => row.id === id)
  if (idx < 0) return null
  const current = getStore()[idx]
  if (current.status !== 'qualified') {
    throw new Error('convert_not_allowed')
  }
  const accountId =
    input.account_id ??
    `acc-${Date.now().toString(36)}`
  const accountName = input.create_account?.name?.trim() || current.title
  const updated: Lead = {
    ...current,
    status: 'converted',
    lifecycle_stage: 'retain',
    converted_account_id: accountId,
    converted_contact_id: input.contact_id ?? null,
    tags: [...current.tags, accountName !== current.title ? `公司:${accountName}` : '已转化'].filter(
      (tag, i, arr) => arr.indexOf(tag) === i,
    ),
    updated_at: new Date().toISOString(),
  }
  getStore()[idx] = updated
  return updated
}
