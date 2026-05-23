import { MOCK_LEADS_SEED } from '~/fixtures/leads.mock'
import type {
  Lead,
  LeadConvertInput,
  LeadCreateInput,
  LeadListQuery,
  LeadUpdateInput,
  LeadsListData,
  Pagination,
} from '~/types/lead'

function cloneSeed(): Lead[] {
  return MOCK_LEADS_SEED.map((row) => ({ ...row, tags: [...row.tags] }))
}

let store: Lead[] | null = null

function getStore(): Lead[] {
  if (!store) store = cloneSeed()
  return store
}

export function resetLeadsMockStore(): void {
  store = null
}

function matchesQuery(lead: Lead, query: LeadListQuery): boolean {
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

export function mockListLeads(query: LeadListQuery = {}): { data: LeadsListData; pagination: Pagination } {
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
