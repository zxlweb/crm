import type { DealStatsByStage, DealStatsWinRate } from '~/types/deal-stats'

export const MOCK_DEAL_STATS_BY_STAGE: DealStatsByStage = {
  items: [
    { label: 'qualification', value: 4, amount: 580000 },
    { label: 'proposal', value: 2, amount: 520000 },
    { label: 'negotiation', value: 1, amount: 960000 },
    { label: 'won', value: 1, amount: 340000 },
    { label: 'lost', value: 1, amount: 120000 },
  ],
  total: 9,
}

export const MOCK_DEAL_STATS_WIN_RATE: DealStatsWinRate = {
  items: [
    { period: '2026-W18', won: 1, lost: 0, rate: 1 },
    { period: '2026-W19', won: 0, lost: 1, rate: 0 },
    { period: '2026-W20', won: 2, lost: 1, rate: 0.67 },
    { period: '2026-W21', won: 1, lost: 0, rate: 1 },
  ],
}

export function mockDealStatsByStage(): DealStatsByStage {
  return structuredClone(MOCK_DEAL_STATS_BY_STAGE)
}

export function mockDealStatsWinRate(): DealStatsWinRate {
  return structuredClone(MOCK_DEAL_STATS_WIN_RATE)
}
