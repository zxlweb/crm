export type LeadStatsQuery = {
  from?: string
  to?: string
  granularity?: 'day' | 'week'
}

export type LeadStatsDistribution = {
  items: Array<{
    label: string
    value: number
    percentage?: number
  }>
  total: number
}

export type LeadStatsTrend = {
  categories: string[]
  series: Array<{
    name: string
    data: number[]
    primary?: boolean
  }>
}

export type LeadStatsFunnel = {
  stages: Array<{ name: string; count: number }>
  conversion_rates: Array<{ from: string; to: string; rate: number }>
}
