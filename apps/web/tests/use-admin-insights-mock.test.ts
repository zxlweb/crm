import { describe, expect, it } from 'vitest'
import { mockTenantHealth, mockPlanDistribution, mockTopTenants } from '~/fixtures/admin-insights.mock'

describe('admin insights mock fixtures', () => {
  it('tenant health has 5 dimensions', () => {
    const data = mockTenantHealth()
    expect(data.dimensions).toHaveLength(5)
    expect(data.dimensions).toContain('activity')
    expect(data.dimensions).toContain('audit_risk')
  })

  it('each tenant has all dimension scores 0-100', () => {
    const data = mockTenantHealth()
    expect(data.items.length).toBeGreaterThan(0)
    for (const tenant of data.items) {
      expect(tenant.tenant_id).toBeTruthy()
      expect(tenant.tenant_name).toBeTruthy()
      for (const dim of data.dimensions) {
        const score = tenant.scores[dim]
        expect(score).toBeGreaterThanOrEqual(0)
        expect(score).toBeLessThanOrEqual(100)
      }
      expect(tenant.overall_score).toBeGreaterThanOrEqual(0)
      expect(tenant.overall_score).toBeLessThanOrEqual(100)
    }
  })

  it('plan distribution returns non-empty items', () => {
    const data = mockPlanDistribution()
    expect(data.items.length).toBeGreaterThan(0)
    for (const item of data.items) {
      expect(item.plan).toBeTruthy()
      expect(item.count).toBeGreaterThan(0)
    }
  })

  it('top tenants returns ranked items with metric', () => {
    const data = mockTopTenants()
    expect(data.metric).toBeTruthy()
    expect(data.items.length).toBeGreaterThan(0)
    const values = data.items.map((t) => t.value)
    expect(values).toEqual([...values].sort((a, b) => b - a))
  })
})
