import { mockDealStatsByStage, mockDealStatsWinRate } from '~/fixtures/deal-stats.mock'
import type { DealStatsByStage, DealStatsQuery, DealStatsWinRate } from '~/types/deal-stats'

export function defaultDealStatsRange(): Pick<DealStatsQuery, 'from' | 'to'> {
  const to = new Date()
  const from = new Date(to)
  from.setDate(from.getDate() - 90)
  return {
    from: from.toISOString().slice(0, 10),
    to: to.toISOString().slice(0, 10),
  }
}

function buildQuery(q: DealStatsQuery): string {
  const params = new URLSearchParams()
  if (q.from) params.set('from', q.from)
  if (q.to) params.set('to', q.to)
  if (q.metric) params.set('metric', q.metric)
  if (q.granularity) params.set('granularity', q.granularity)
  const s = params.toString()
  return s ? `?${s}` : ''
}

export function useDealStats() {
  const api = useApi()
  const config = useRuntimeConfig()
  const forceMock = computed(() => config.public.useDealsMock === true)

  async function fetchByStage(q: DealStatsQuery = defaultDealStatsRange()): Promise<DealStatsByStage> {
    if (forceMock.value) {
      return mockDealStatsByStage()
    }
    return api.request<DealStatsByStage>(`/api/deals/stats/by-stage${buildQuery(q)}`)
  }

  async function fetchWinRate(q: DealStatsQuery = defaultDealStatsRange()): Promise<DealStatsWinRate> {
    if (forceMock.value) {
      return mockDealStatsWinRate()
    }
    return api.request<DealStatsWinRate>(`/api/deals/stats/win-rate${buildQuery({ ...q, granularity: q.granularity ?? 'month' })}`)
  }

  return { forceMock, fetchByStage, fetchWinRate }
}
