/** Phase 2 API — docs/api/phase-2-crm-ai.md §5 */

import type { LifecycleStage, RelationshipHealth } from '~/types/lead'

export type { LifecycleStage, RelationshipHealth }

export type Contact = {
  id: string
  account_id: string | null
  first_name: string
  last_name: string
  display_name: string
  email: string
  phone: string
  is_primary: boolean
  owner_id: string | null
  lifecycle_stage: LifecycleStage
  relationship_health: RelationshipHealth
  engagement_score: number
  last_activity_at: string | null
  tags: string[]
  created_at: string
  updated_at: string
}

export type ContactListQuery = {
  page?: number
  page_size?: number
  search?: string
  lifecycle_stage?: LifecycleStage
  relationship_health?: RelationshipHealth
  account_id?: string
  owner_id?: string
}

export type Pagination = {
  page: number
  page_size: number
  total: number
}

export type ContactsListData = {
  items: Contact[]
}

export type ContactCreateInput = {
  account_id?: string | null
  first_name?: string
  last_name?: string
  email?: string
  phone?: string
  is_primary?: boolean
  owner_id?: string | null
  lifecycle_stage?: LifecycleStage
  tags?: string[]
}

export type ContactUpdateInput = Partial<ContactCreateInput> & {
  clear_account?: boolean
}
