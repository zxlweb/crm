/** Phase 3 API — docs/api/phase-3-deals-dashboard-api.md §2–§3 */

export type DealStage =
  | 'qualification'
  | 'proposal'
  | 'negotiation'
  | 'won'
  | 'lost'

export type DealCurrency = 'CNY' | 'USD'

export type Deal = {
  id: string
  tenant_id: string
  title: string
  stage: DealStage
  amount: number
  currency: DealCurrency
  probability: number
  expected_close_date: string | null
  account_id: string | null
  lead_id: string | null
  contact_id: string | null
  owner_id: string | null
  description: string
  tags: string[]
  lost_reason: string | null
  closed_at: string | null
  engagement_score: number
  last_activity_at: string | null
  created_at: string
  updated_at: string
}

export type DealListQuery = {
  page?: number
  page_size?: number
  search?: string
  stage?: DealStage
  stages?: string
  owner_id?: string
  account_id?: string
  lead_id?: string
}

export type DealCreateInput = {
  title: string
  stage?: DealStage
  amount?: number
  currency?: DealCurrency
  probability?: number
  expected_close_date?: string | null
  account_id?: string | null
  lead_id?: string | null
  contact_id?: string | null
  owner_id?: string | null
  description?: string
  tags?: string[]
}

export type DealUpdateInput = Partial<DealCreateInput> & {
  stage?: DealStage
  lost_reason?: string | null
}

export type DealStageUpdateInput = {
  stage: DealStage
  lost_reason?: string | null
  note?: string
}

export type DealPipelineItem = Pick<
  Deal,
  | 'id'
  | 'title'
  | 'amount'
  | 'currency'
  | 'probability'
  | 'expected_close_date'
  | 'account_id'
  | 'owner_id'
>

export type DealPipelineStage = {
  stage: DealStage
  count: number
  amount_total: number
  items: DealPipelineItem[]
}

export type DealPipelineSummary = {
  open_count: number
  open_amount: number
  won_count_mtd: number
  won_amount_mtd: number
}

export type DealPipelineData = {
  stages: DealPipelineStage[]
  summary: DealPipelineSummary
}

export type DealsListData = {
  items: Deal[]
}

export type Pagination = {
  page: number
  page_size: number
  total: number
}

export const DEAL_PIPELINE_STAGES: DealStage[] = [
  'qualification',
  'proposal',
  'negotiation',
  'won',
  'lost',
]

export const OPEN_DEAL_STAGES: DealStage[] = ['qualification', 'proposal', 'negotiation']
