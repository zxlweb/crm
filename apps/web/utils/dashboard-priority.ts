import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'
import type { Account } from '~/types/account'
import type { Lead, RelationshipHealth } from '~/types/lead'
import type { PriorityActionItem, PriorityHealthLabel } from '~/types/dashboard'

function daysSince(iso: string | null): number | null {
  if (!iso) return null
  const diff = Date.now() - new Date(iso).getTime()
  return Math.max(0, Math.floor(diff / 86400000))
}

function buildSparkline(score: number): number[] {
  const v = Math.min(100, Math.max(0, score))
  return [Math.max(0, v - 24), Math.max(0, v - 16), Math.max(0, v - 7), v]
}

function healthLabelFrom(
  health: RelationshipHealth,
  urgency: PriorityActionItem['urgency'],
): PriorityHealthLabel {
  if (health === 'low' || urgency === 'coral') return 'alert'
  if (health === 'medium' || urgency === 'amber') return 'watch'
  return 'healthy'
}

function pickContactName(tags: string[]): string | undefined {
  const tag = tags.find((t) => t.trim().length > 0)
  if (!tag) return undefined
  if (/总监|经理|总|先生|女士|经理/.test(tag) || tag.length <= 6) return tag
  return undefined
}

function scoreRecord(input: {
  relationship_health: RelationshipHealth
  last_activity_at: string | null
  engagement_score: number
  status?: Lead['status']
}): number {
  let score = 0
  const inactive = daysSince(input.last_activity_at)

  if (input.relationship_health === 'low') score += 40
  else if (input.relationship_health === 'medium') score += 12

  if (inactive != null) {
    if (inactive >= 14) score += 35
    else if (inactive >= 7) score += 28
    else if (inactive >= 3) score += 8
  } else {
    score += 20
  }

  if (input.engagement_score < 30) score += 12
  if (input.status === 'qualified' && inactive != null && inactive >= 5) score += 18

  return score
}

function urgencyFromScore(score: number): PriorityActionItem['urgency'] {
  if (score >= 55) return 'coral'
  if (score >= 30) return 'amber'
  return 'neutral'
}

function pickLeadSuggestion(
  lead: Lead,
  labels: { suggestReengage: string; suggestCall: string; suggestFollowup: string },
): string {
  const inactive = daysSince(lead.last_activity_at)
  if (inactive != null && inactive >= 21) return labels.suggestReengage
  if (lead.relationship_health === 'low') return labels.suggestCall
  if (lead.status === 'qualified') return labels.suggestFollowup
  return labels.suggestFollowup
}

function pickAccountSuggestion(
  account: Account,
  labels: { suggestReengage: string; suggestVisit: string; suggestFollowup: string },
): string {
  const inactive = daysSince(account.last_activity_at)
  if (inactive != null && inactive >= 21) return labels.suggestReengage
  if (account.relationship_health === 'low') return labels.suggestVisit
  return labels.suggestFollowup
}

function pickLeadInsightTip(
  lead: Lead,
  labels: {
    insightTipStale: string
    insightTipLowEngagement: string
    insightTipCall: string
    insightTipFollowup: string
  },
): string {
  const inactive = daysSince(lead.last_activity_at)
  if (inactive != null && inactive >= 21) return labels.insightTipStale
  if (lead.engagement_score < 30) return labels.insightTipLowEngagement
  if (lead.relationship_health === 'low') return labels.insightTipCall
  return labels.insightTipFollowup
}

function pickAccountInsightTip(
  account: Account,
  labels: {
    insightTipStale: string
    insightTipLowEngagement: string
    insightTipVisit: string
    insightTipFollowup: string
  },
): string {
  const inactive = daysSince(account.last_activity_at)
  if (inactive != null && inactive >= 21) return labels.insightTipStale
  if (account.engagement_score < 30) return labels.insightTipLowEngagement
  if (account.relationship_health === 'low') return labels.insightTipVisit
  return labels.insightTipFollowup
}

export function buildPriorityFromLead(
  lead: Lead,
  labels: {
    reasonNoFollowup: (days: number) => string
    reasonLowHealth: string
    reasonMediumHealth: string
    reasonHealthDeclining: string
    suggestCall: string
    suggestFollowup: string
    suggestReengage: string
    insightTipCall: string
    insightTipFollowup: string
    insightTipStale: string
    insightTipLowEngagement: string
  },
): PriorityActionItem | null {
  const score = scoreRecord(lead)
  if (score < 20) return null

  const reasons: string[] = []
  const inactive = daysSince(lead.last_activity_at)
  if (inactive != null && inactive >= 7) {
    reasons.push(labels.reasonNoFollowup(inactive))
  }
  if (lead.relationship_health === 'low') {
    reasons.push(labels.reasonLowHealth)
    if (inactive != null && inactive >= 3) {
      reasons.push(labels.reasonHealthDeclining)
    }
  } else if (lead.relationship_health === 'medium' && reasons.length === 0) {
    reasons.push(labels.reasonMediumHealth)
  }
  if (reasons.length === 0 && inactive != null && inactive >= 3) {
    reasons.push(labels.reasonNoFollowup(inactive))
  }
  if (reasons.length === 0) return null

  const urgency = urgencyFromScore(score)
  const suggestion = pickLeadSuggestion(lead, labels)
  const insightTip = pickLeadInsightTip(lead, labels)

  return {
    id: `lead-${lead.id}`,
    entityType: 'lead',
    entityId: lead.id,
    title: lead.title,
    reasons,
    suggestion,
    insightTip,
    followHref: `/leads/${lead.id}#timeline`,
    urgency,
    healthLabel: healthLabelFrom(lead.relationship_health, urgency),
    isPreview: lead.id === DEMO_LEAD_ID,
    score,
    engagementScore: lead.engagement_score,
    sparkline: buildSparkline(lead.engagement_score),
    contactName: pickContactName(lead.tags),
    daysSinceActivity: inactive,
  }
}

export function buildPriorityFromAccount(
  account: Account,
  labels: {
    reasonNoFollowup: (days: number) => string
    reasonLowHealth: string
    reasonMediumHealth: string
    reasonHealthDeclining: string
    suggestVisit: string
    suggestFollowup: string
    suggestReengage: string
    insightTipVisit: string
    insightTipFollowup: string
    insightTipStale: string
    insightTipLowEngagement: string
  },
): PriorityActionItem | null {
  const score = scoreRecord(account)
  if (score < 20) return null

  const reasons: string[] = []
  const inactive = daysSince(account.last_activity_at)
  if (inactive != null && inactive >= 7) {
    reasons.push(labels.reasonNoFollowup(inactive))
  }
  if (account.relationship_health === 'low') {
    reasons.push(labels.reasonLowHealth)
    if (inactive != null && inactive >= 3) {
      reasons.push(labels.reasonHealthDeclining)
    }
  } else if (account.relationship_health === 'medium' && reasons.length === 0) {
    reasons.push(labels.reasonMediumHealth)
  }
  if (reasons.length === 0) return null

  const urgency = urgencyFromScore(score)
  const suggestion = pickAccountSuggestion(account, labels)
  const insightTip = pickAccountInsightTip(account, labels)

  return {
    id: `account-${account.id}`,
    entityType: 'account',
    entityId: account.id,
    title: account.name,
    reasons,
    suggestion,
    insightTip,
    followHref: `/accounts/${account.id}?tab=timeline`,
    urgency,
    healthLabel: healthLabelFrom(account.relationship_health, urgency),
    score,
    engagementScore: account.engagement_score,
    sparkline: buildSparkline(account.engagement_score),
    contactName: pickContactName(account.tags),
    daysSinceActivity: inactive,
  }
}

export function mergePriorities(items: PriorityActionItem[], limit = 7): PriorityActionItem[] {
  return [...items].sort((a, b) => b.score - a.score).slice(0, limit)
}
