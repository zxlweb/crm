export interface AuditStatItem {
  action: string
  count: number
}

export interface AuditByAction {
  items: AuditStatItem[]
  total: number
}

export interface AuditTrendPoint {
  date: string
  count: number
}

export interface AuditTrend {
  items: AuditTrendPoint[]
  granularity: 'day' | 'week'
}

export interface AuditTopActor {
  actor_id: string
  actor_name: string
  count: number
}

export interface AuditTopActors {
  items: AuditTopActor[]
}

export interface AuditStatsQuery {
  from?: string
  to?: string
  module?: string
  actor_role?: string
  action?: string
  granularity?: 'day' | 'week'
  limit?: number
}
