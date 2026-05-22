import { describe, expect, it, beforeEach } from 'vitest'
import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'
import {
  mockCreateLead,
  mockDeleteLead,
  mockGetLead,
  mockListLeads,
  mockUpdateLead,
  resetLeadsMockStore,
} from '~/utils/leads-mock-store'

describe('leads-mock-store', () => {
  beforeEach(() => {
    resetLeadsMockStore()
  })

  it('lists demo lead 华创科技', () => {
    const { data, pagination } = mockListLeads({ search: '华创' })
    expect(pagination.total).toBe(1)
    expect(data.items[0]?.id).toBe(DEMO_LEAD_ID)
    expect(data.items[0]?.title).toContain('华创')
  })

  it('filters by status', () => {
    const { data } = mockListLeads({ status: 'converted' })
    expect(data.items.every((row) => row.status === 'converted')).toBe(true)
  })

  it('creates and updates a lead', () => {
    const created = mockCreateLead({ title: '测试线索' }, 't1', 'u1')
    expect(created.title).toBe('测试线索')
    expect(created.status).toBe('new')

    const updated = mockUpdateLead(created.id, { status: 'contacted' })
    expect(updated?.status).toBe('contacted')
    expect(mockGetLead(created.id)?.status).toBe('contacted')
  })

  it('deletes a lead', () => {
    const row = mockGetLead(DEMO_LEAD_ID)
    expect(row).not.toBeNull()
    expect(mockDeleteLead(DEMO_LEAD_ID)).toBe(true)
    expect(mockGetLead(DEMO_LEAD_ID)).toBeNull()
  })
})
