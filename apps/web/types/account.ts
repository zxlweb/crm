/** Phase 2 API — docs/api/phase-2-crm-ai.md §4 */

import type { LifecycleStage, RelationshipHealth } from '~/types/lead'

export type { LifecycleStage, RelationshipHealth }

export type Account = {
  id: string
  name: string
  industry: string
  website: string
  owner_id: string | null
  lifecycle_stage: LifecycleStage
  relationship_health: RelationshipHealth
  engagement_score: number
  last_activity_at: string | null
  tags: string[]
  created_at: string
  updated_at: string
}

export type AccountListQuery = {
  page?: number
  page_size?: number
  search?: string
  lifecycle_stage?: LifecycleStage
  relationship_health?: RelationshipHealth
  segment?: string
  owner_id?: string
}

export type Pagination = {
  page: number
  page_size: number
  total: number
}

export type AccountsListData = {
  items: Account[]
}

export type AccountCreateInput = {
  name: string
  industry?: string
  website?: string
  owner_id?: string | null
  lifecycle_stage?: LifecycleStage
  tags?: string[]
}

export type AccountUpdateInput = Partial<AccountCreateInput>

export type AccountDetailTab = 'overview' | 'timeline' | 'emotion'
