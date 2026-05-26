import { mockAuditByAction, mockAuditTopActors, mockAuditTrend } from '~/fixtures/audit-stats.mock'
import type {
  AuditByAction,
  AuditStatsQuery,
  AuditTopActors,
  AuditTrend,
} from '~/types/audit-stats'
import {
  normalizeAuditByAction,
  normalizeAuditTopActors,
  normalizeAuditTrend,
} from '~/utils/phase4-api-normalize'

export function defaultAuditRange(): Pick<AuditStatsQuery, 'from' | 'to'> {
  const to = new Date()
  const from = new Date(to)
  from.setDate(from.getDate() - 7)
  return {
    from: from.toISOString().slice(0, 10),
    to: to.toISOString().slice(0, 10),
  }
}

function buildQuery(q: AuditStatsQuery): string {
  const params = new URLSearchParams()
  if (q.from) params.set('from', q.from)
  if (q.to) params.set('to', q.to)
  if (q.module) params.set('module', q.module)
  if (q.actor_role) params.set('actor_role', q.actor_role)
  if (q.action) params.set('action', q.action)
  if (q.granularity) params.set('granularity', q.granularity)
  if (q.limit) params.set('limit', String(q.limit))
  const s = params.toString()
  return s ? `?${s}` : ''
}

export function useAuditStats() {
  const api = useApi()
  const config = useRuntimeConfig()
  const forceMock = computed(() => config.public.useSettingsMock === true || config.public.useSettingsMock === 'true')

  async function fetchByAction(q: AuditStatsQuery = defaultAuditRange()): Promise<AuditByAction> {
    if (forceMock.value) return mockAuditByAction()
    const data = await api.request<unknown>(`/api/audit/stats/by-action${buildQuery(q)}`)
    return normalizeAuditByAction(data)
  }

  async function fetchTrend(q: AuditStatsQuery = defaultAuditRange()): Promise<AuditTrend> {
    if (forceMock.value) return mockAuditTrend()
    const data = await api.request<unknown>(
      `/api/audit/stats/trend${buildQuery({ ...q, granularity: q.granularity ?? 'day' })}`,
    )
    return normalizeAuditTrend(data)
  }

  async function fetchTopActors(q: AuditStatsQuery = { ...defaultAuditRange(), limit: 10 }): Promise<AuditTopActors> {
    if (forceMock.value) return mockAuditTopActors()
    const data = await api.request<unknown>(`/api/audit/stats/top-actors${buildQuery(q)}`)
    return normalizeAuditTopActors(data)
  }

  async function exportCsv(q: AuditStatsQuery = defaultAuditRange()): Promise<Blob> {
    if (forceMock.value) return new Blob(['action,count\nsettings.update,18\n'], { type: 'text/csv' })
    const config2 = useRuntimeConfig()
    const headers = new Headers()
    const token = useAuth().accessToken.value
    if (token) headers.set('Authorization', `Bearer ${token}`)
    const tenantId = useCookie<string | null>('crm.tenant_id').value
    if (tenantId) headers.set('X-Tenant-ID', tenantId)

    const params = buildQuery({ ...q, granularity: undefined })
    const sep = params.includes('?') ? '&' : '?'
    const res = await fetch(`${config2.public.apiBase}/api/audit/export${params}${sep}format=csv`, { headers })
    if (!res.ok) {
      const text = await res.text()
      throw new Error(text.trim().slice(0, 120) || `Export failed: ${res.status}`)
    }
    return res.blob()
  }

  return { forceMock, fetchByAction, fetchTrend, fetchTopActors, exportCsv }
}
