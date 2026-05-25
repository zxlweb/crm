import type { LifecycleStage, RelationshipHealth } from '~/types/lead'

export type InsightSubjectType = 'lead' | 'account' | 'contact'

export type InsightSuggestedAction = {
  activity_event_type: string
  activity_direction: string
  title: string
}

export type InsightItem = {
  id: string
  priority: number
  title_key: string
  body_key: string
  suggested_action?: InsightSuggestedAction
  rule_id: string
  explanation_key: string
}

export type InsightEvaluateResult = {
  items: InsightItem[]
  engagement_score: number
  engagement_explanation_key: string
  lifecycle_stage: LifecycleStage
  relationship_health: RelationshipHealth
}

export type InsightPanelItem = {
  id: string
  title: string
  body: string
}
