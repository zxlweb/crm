import type { RelationshipHealth } from '~/types/lead'

export type PriorityEntityType = 'lead' | 'account'

export type PriorityHealthLabel = 'alert' | 'watch' | 'healthy'

export type PriorityActionItem = {
  id: string
  entityType: PriorityEntityType
  entityId: string
  title: string
  reasons: string[]
  suggestion: string
  /** Rich copy for tinted insight strip (Figma) */
  insightTip: string
  followHref: string
  urgency: 'amber' | 'coral' | 'neutral'
  healthLabel: PriorityHealthLabel
  isPreview?: boolean
  score: number
  engagementScore: number
  sparkline: number[]
  contactName?: string
  daysSinceActivity: number | null
}

export type DashboardKpiTrends = {
  leadsWeeklyTouch: number
  accountsWeeklyTouch: number
  engagementDelta: number
  engagementDirection: 'up' | 'down' | 'flat'
}

export type DashboardSnapshot = {
  leadsTotal: number
  accountsTotal: number
  atRiskTotal: number
  avgEngagement: number
  priorityCount: number
  priorities: PriorityActionItem[]
  pipelineLeads: import('~/types/lead').Lead[]
  pipelineAccounts: import('~/types/account').Account[]
  weeklyFollowUpCount: number
  kpiTrends: DashboardKpiTrends
}

export type DashboardInsightVariant = 'churn' | 'opportunity' | 'rule'

export type DashboardInsightItem = {
  id: string
  variant: DashboardInsightVariant
  title: string
  body: string
  isPreview?: boolean
  urgent?: boolean
}

export type DashboardCalendarEvent = {
  id: string
  date: string
  time: string
  title: string
  subtitle: string
  href?: string
}

export type TeamHeatmapMemberRow = {
  memberId: string
  memberName: string
  healthScore: number
  emotionScore: number
}

/** @deprecated Grid heatmap cells — use TeamHeatmapMemberRow */
export type TeamHeatmapCell = {
  memberId: string
  memberName: string
  health: RelationshipHealth
  emotion: 'positive' | 'neutral' | 'negative'
  count: number
}
