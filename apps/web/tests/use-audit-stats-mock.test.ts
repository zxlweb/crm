import { describe, expect, it } from 'vitest'
import { mockAuditByAction, mockAuditTrend, mockAuditTopActors } from '~/fixtures/audit-stats.mock'

describe('audit stats mock fixtures', () => {
  it('by-action returns items with valid counts', () => {
    const data = mockAuditByAction()
    expect(data.items.length).toBeGreaterThan(0)
    expect(data.total).toBeGreaterThan(0)
    for (const item of data.items) {
      expect(item.action).toBeTruthy()
      expect(item.count).toBeGreaterThan(0)
    }
  })

  it('trend returns daily data points', () => {
    const data = mockAuditTrend()
    expect(data.granularity).toBe('day')
    expect(data.items.length).toBeGreaterThan(0)
    for (const point of data.items) {
      expect(point.date).toMatch(/^\d{4}-\d{2}-\d{2}$/)
      expect(typeof point.count).toBe('number')
    }
  })

  it('top actors returns ranked list', () => {
    const data = mockAuditTopActors()
    expect(data.items.length).toBeGreaterThan(0)
    for (const actor of data.items) {
      expect(actor.actor_id).toBeTruthy()
      expect(actor.actor_name).toBeTruthy()
      expect(actor.count).toBeGreaterThan(0)
    }
    const counts = data.items.map((a) => a.count)
    expect(counts).toEqual([...counts].sort((a, b) => b - a))
  })
})
