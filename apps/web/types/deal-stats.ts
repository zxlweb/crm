/** Phase 3 — docs/api/phase-3-deals-dashboard-api.md §3.5 */

export type DealStatsQuery = {
  from?: string
  to?: string
  metric?: 'count' | 'amount'
  granularity?: 'week' | 'month'
}

export type DealStatsByStageItem = {
  label: string
  value: number
  amount?: number
}

export type DealStatsByStage = {
  items: DealStatsByStageItem[]
  total: number
}

export type DealStatsWinRateItem = {
  period: string
  won: number
  lost: number
  rate: number
}

export type DealStatsWinRate = {
  items: DealStatsWinRateItem[]
}
