import { describe, expect, it, beforeEach } from 'vitest'
import { MOCK_LEADS_SEED } from '~/fixtures/leads.mock'
import { mockConvertLead, mockGetLead, resetLeadsMockStore } from '~/utils/leads-mock-store'

describe('leads-mock-store convert', () => {
  beforeEach(() => resetLeadsMockStore())

  it('converts qualified lead', () => {
    const qualified = MOCK_LEADS_SEED.find((r) => r.status === 'qualified')!
    const row = mockConvertLead(qualified.id, { create_account: { name: '华创科技（客户）' } })
    expect(row?.status).toBe('converted')
    expect(row?.converted_account_id).toBeTruthy()
    expect(mockGetLead(qualified.id)?.status).toBe('converted')
  })

  it('rejects unqualified lead', () => {
    const bad = MOCK_LEADS_SEED.find((r) => r.status === 'unqualified')!
    expect(() => mockConvertLead(bad.id, { create_account: { name: 'X' } })).toThrow('convert_not_allowed')
  })
})
