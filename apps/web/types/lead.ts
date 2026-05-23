/** Phase 2 API contract — docs/api/phase-2-crm-ai.md §2.3, §6 */

export type LeadStatus =
  | 'new'
  | 'contacted'
  | 'qualified'
  | 'unqualified'
  | 'converted'

export type LifecycleStage =
  | 'acquire'
  | 'activate'
  | 'grow'
  | 'retain'
  | 'revive'

export type RelationshipHealth = 'high' | 'medium' | 'low'

export type Lead = {
  id: string
  tenant_id: string
  owner_id: string | null
  title: string
  status: LeadStatus
  source: string
  amount: number
  expected_close_date: string | null
  lifecycle_stage: LifecycleStage
  relationship_health: RelationshipHealth
  engagement_score: number
  last_activity_at: string | null
  tags: string[]
  converted_account_id: string | null
  converted_contact_id: string | null
  created_at: string
  updated_at: string
}

export type LeadListQuery = {
  page?: number
  page_size?: number
  status?: LeadStatus
  source?: string
  owner_id?: string
  lifecycle_stage?: LifecycleStage
  relationship_health?: RelationshipHealth
  segment?: string
  search?: string
}

export type Pagination = {
  page: number
  page_size: number
  total: number
}

export type LeadsListData = {
  items: Lead[]
}

export type LeadCreateInput = {
  title: string
  status?: LeadStatus
  source?: string
  amount?: number
  expected_close_date?: string | null
  owner_id?: string | null
  lifecycle_stage?: LifecycleStage
  tags?: string[]
}

export type LeadUpdateInput = Partial<LeadCreateInput> & {
  status?: LeadStatus
}

export type LeadConvertInput = {
  account_id?: string
  contact_id?: string
  create_account?: { name: string }
}

/** 线索详情 Tab（无概览，信息集中在头部 + KPI） */
export type LeadDetailTab = 'timeline' | 'emotion'
