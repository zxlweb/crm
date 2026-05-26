import type { AuditByAction, AuditTopActors, AuditTrend } from '~/types/audit-stats'

export function mockAuditByAction(): AuditByAction {
  return {
    items: [
      { action: 'settings.update', count: 18 },
      { action: 'custom_field.create', count: 7 },
      { action: 'rbac.role.assign', count: 5 },
      { action: 'lead.create', count: 42 },
      { action: 'lead.update', count: 38 },
      { action: 'deal.create', count: 21 },
      { action: 'account.update', count: 15 },
      { action: 'audit.export', count: 3 },
    ],
    total: 149,
  }
}

export function mockAuditTrend(): AuditTrend {
  const today = new Date()
  const items = Array.from({ length: 14 }, (_, i) => {
    const d = new Date(today)
    d.setDate(d.getDate() - (13 - i))
    return {
      date: d.toISOString().slice(0, 10),
      count: Math.floor(Math.random() * 30) + 5,
    }
  })
  return { items, granularity: 'day' }
}

export function mockAuditTopActors(): AuditTopActors {
  return {
    items: [
      { actor_id: 'u-001', actor_name: 'Alice Wang', count: 56 },
      { actor_id: 'u-002', actor_name: 'Bob Li', count: 41 },
      { actor_id: 'u-003', actor_name: 'Charlie Chen', count: 28 },
      { actor_id: 'u-004', actor_name: 'Diana Zhang', count: 15 },
      { actor_id: 'u-005', actor_name: 'Eric Liu', count: 9 },
    ],
  }
}
