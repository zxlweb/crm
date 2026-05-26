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

/**
 * 客户池视图：
 * - `mine`   私海：当前登录用户负责的线索
 * - `public` 公海：owner_id 为 null，任何销售可领取
 * - `others` 他人私海：分配给其他销售（仅 admin / 主管视角）
 * - `all`    全部
 */
export type LeadPool = 'mine' | 'public' | 'others' | 'all'

export const LEAD_POOLS: LeadPool[] = ['mine', 'public', 'others', 'all']

export function isLeadPool(value: unknown): value is LeadPool {
  return typeof value === 'string' && (LEAD_POOLS as string[]).includes(value)
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
  pool?: LeadPool
}

/** 客户池设置（mock：仅存在于前端模块作用域） */
export type LeadPoolSettings = {
  /** 自动回收开关 */
  enabled: boolean
  /** 私海中无活动达到该天数后回收至公海 */
  inactiveDays: number
  /** 上次自动扫描回收时间（ISO） */
  last_recycled_at: string | null
}

export type LeadPoolStats = {
  mine: number
  public: number
  others: number
  all: number
  /** 距离自动回收 ≤ 3 天的私海线索数 */
  recyclableSoon: number
}

export type LeadClaimResult = {
  lead: Lead
  recycledCount?: number
}

export type LeadRecycleSummary = {
  recycled: number
  scanned: number
  threshold_days: number
  recycled_ids: string[]
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
