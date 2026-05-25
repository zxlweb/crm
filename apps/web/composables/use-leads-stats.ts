import type {
  LeadStatsDistribution,
  LeadStatsFunnel,
  LeadStatsQuery,
  LeadStatsTrend,
} from '~/types/lead-stats'

function buildStatsQueryString(query: LeadStatsQuery = {}): string {
  const params = new URLSearchParams()
  if (query.from) params.set('from', query.from)
  if (query.to) params.set('to', query.to)
  if (query.granularity) params.set('granularity', query.granularity)
  const qs = params.toString()
  return qs ? `?${qs}` : ''
}

export function useLeadsStats() {
  const api = useApi()

  function fetchBySource(query: LeadStatsQuery = {}) {
    return api.request<LeadStatsDistribution>(
      `/api/leads/stats/by-source${buildStatsQueryString(query)}`,
    )
  }

  function fetchByStatus(query: LeadStatsQuery = {}) {
    return api.request<LeadStatsDistribution>(
      `/api/leads/stats/by-status${buildStatsQueryString(query)}`,
    )
  }

  function fetchTrend(query: LeadStatsQuery = {}) {
    return api.request<LeadStatsTrend>(`/api/leads/stats/trend${buildStatsQueryString(query)}`)
  }

  function fetchFunnel(query: LeadStatsQuery = {}) {
    return api.request<LeadStatsFunnel>(`/api/leads/stats/funnel${buildStatsQueryString(query)}`)
  }

  return { fetchBySource, fetchByStatus, fetchTrend, fetchFunnel }
}

/** 默认近 30 日（UTC 日期字符串） */
export function defaultLeadStatsRange(): { from: string; to: string } {
  const end = new Date()
  const start = new Date(end)
  start.setUTCDate(start.getUTCDate() - 29)
  const fmt = (d: Date) => d.toISOString().slice(0, 10)
  return { from: fmt(start), to: fmt(end) }
}
