import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'
import { DASHBOARD_PREVIEW_WEEKLY_FOLLOWUPS } from '~/fixtures/dashboard-preview'
import type { Account } from '~/types/account'
import type { Lead } from '~/types/lead'
import type { DashboardSummaryPriority } from '~/types/dashboard-stats'
import type { DashboardInsightItem, DashboardKpiTrends, DashboardSnapshot, PriorityActionItem } from '~/types/dashboard'
import {
  buildPriorityFromAccount,
  buildPriorityFromLead,
  mergePriorities,
} from '~/utils/dashboard-priority'
import { buildEngagementSparkline } from '~/utils/dashboard-sparkline'

function averageEngagement(items: Array<{ engagement_score: number }>): number {
  if (items.length === 0) return 0
  const sum = items.reduce((acc, row) => acc + row.engagement_score, 0)
  return Math.round(sum / items.length)
}

function countRecentTouch(
  items: Array<{ last_activity_at: string | null }>,
  withinDays = 7,
): number {
  const cutoff = Date.now() - withinDays * 86400000
  return items.filter((row) => {
    if (!row.last_activity_at) return false
    return new Date(row.last_activity_at).getTime() >= cutoff
  }).length
}

function mapSummaryPriority(p: DashboardSummaryPriority): PriorityActionItem {
  const entityType = p.entity_type
  const entityId = String(p.entity_id)
  const urgency =
    p.score >= 55 ? 'coral' : p.score >= 30 ? ('amber' as const) : ('neutral' as const)
  const healthLabel =
    p.score >= 55 ? 'alert' : p.score >= 30 ? ('watch' as const) : ('healthy' as const)

  return {
    id: `${entityType}-${entityId}`,
    entityType,
    entityId,
    title: p.title,
    reasons: p.reasons,
    suggestion: p.suggestion,
    insightTip: p.suggestion,
    followHref: entityType === 'lead' ? `/leads/${p.entity_id}` : `/accounts/${p.entity_id}`,
    urgency,
    healthLabel,
    isPreview: p.is_preview,
    score: p.score,
    engagementScore: p.engagement_score,
    sparkline: buildEngagementSparkline(p.engagement_score),
    daysSinceActivity: null,
  }
}

function buildKpiTrends(
  leadsPool: Lead[],
  accountsPool: Account[],
): DashboardKpiTrends {
  const pool = [...leadsPool, ...accountsPool]
  const cutoff = Date.now() - 7 * 86400000
  const recent = pool.filter(
    (row) => row.last_activity_at && new Date(row.last_activity_at).getTime() >= cutoff,
  )
  const older = pool.filter(
    (row) => !row.last_activity_at || new Date(row.last_activity_at).getTime() < cutoff,
  )

  let engagementDelta = 0
  let engagementDirection: DashboardKpiTrends['engagementDirection'] = 'flat'
  if (recent.length > 0 && older.length > 0) {
    const rawDelta = averageEngagement(recent) - averageEngagement(older)
    engagementDelta = Math.min(15, Math.abs(rawDelta))
    if (rawDelta > 2) engagementDirection = 'up'
    else if (rawDelta < -2) engagementDirection = 'down'
  }

  return {
    leadsWeeklyTouch: countRecentTouch(leadsPool),
    accountsWeeklyTouch: countRecentTouch(accountsPool),
    engagementDelta,
    engagementDirection,
  }
}

function buildWeeklyFollowUpCount(
  leadsPool: Lead[],
  accountsPool: Account[],
  isPreviewMode: boolean,
): number {
  const computed = countRecentTouch(leadsPool) + countRecentTouch(accountsPool)
  if (isPreviewMode) {
    return Math.max(computed, DASHBOARD_PREVIEW_WEEKLY_FOLLOWUPS)
  }
  return computed
}

async function safeList<T extends { pagination: { total: number }; data: { items: unknown[] } }>(
  loader: () => Promise<T>,
  fallback: T,
): Promise<T> {
  try {
    return await loader()
  } catch {
    return fallback
  }
}

export function useDashboard() {
  const leadsApi = useLeads()
  const accountsApi = useAccounts()
  const dashboardStats = useDashboardStats()
  const { t } = useI18n()

  function priorityLabels() {
    return {
      reasonNoFollowup: (days: number) => t('dashboardReasonNoFollowup', { days }),
      reasonLowHealth: t('dashboardReasonLowHealth'),
      reasonMediumHealth: t('dashboardReasonMediumHealth'),
      reasonHealthDeclining: t('dashboardReasonHealthDeclining'),
      suggestCall: t('dashboardSuggestCall'),
      suggestFollowup: t('dashboardSuggestFollowup'),
      suggestReengage: t('dashboardSuggestReengage'),
      insightTipCall: t('dashboardInsightTipCall'),
      insightTipFollowup: t('dashboardInsightTipFollowup'),
      insightTipStale: t('dashboardInsightTipStale'),
      insightTipLowEngagement: t('dashboardInsightTipLowEngagement'),
    }
  }

  async function loadSnapshotLegacy(isPreviewMode = false): Promise<DashboardSnapshot> {
    const emptyLeads = { data: { items: [] as Lead[] }, pagination: { page: 1, page_size: 0, total: 0 } }
    const emptyAccounts = { data: { items: [] as Account[] }, pagination: { page: 1, page_size: 0, total: 0 } }

    const labels = priorityLabels()
    const accountLabels = {
      reasonNoFollowup: labels.reasonNoFollowup,
      reasonLowHealth: labels.reasonLowHealth,
      reasonMediumHealth: labels.reasonMediumHealth,
      reasonHealthDeclining: labels.reasonHealthDeclining,
      suggestVisit: t('dashboardSuggestVisit'),
      suggestFollowup: labels.suggestFollowup,
      suggestReengage: labels.suggestReengage,
      insightTipVisit: t('dashboardInsightTipVisit'),
      insightTipFollowup: labels.insightTipFollowup,
      insightTipStale: labels.insightTipStale,
      insightTipLowEngagement: labels.insightTipLowEngagement,
    }

    const [
      pipelineLeadsRes,
      pipelineAccountsRes,
      leadsPoolRes,
      accountsPoolRes,
      atRiskLeadsRes,
      atRiskAccountsRes,
    ] = await Promise.all([
      safeList(() => leadsApi.fetchList({ page: 1, page_size: 3 }), emptyLeads),
      safeList(() => accountsApi.fetchList({ page: 1, page_size: 3 }), emptyAccounts),
      safeList(() => leadsApi.fetchList({ page: 1, page_size: 50 }), emptyLeads),
      safeList(() => accountsApi.fetchList({ page: 1, page_size: 50 }), emptyAccounts),
      safeList(
        () => leadsApi.fetchList({ page: 1, page_size: 1, relationship_health: 'low' }),
        emptyLeads,
      ),
      safeList(
        () => accountsApi.fetchList({ page: 1, page_size: 1, relationship_health: 'low' }),
        emptyAccounts,
      ),
    ])

    const priorityCandidates = [
      ...leadsPoolRes.data.items
        .map((lead) => buildPriorityFromLead(lead, labels))
        .filter(Boolean),
      ...accountsPoolRes.data.items
        .map((account) => buildPriorityFromAccount(account, accountLabels))
        .filter(Boolean),
    ] as NonNullable<ReturnType<typeof buildPriorityFromLead>>[]

    const priorities = ensurePreviewPriorities(
      mergePriorities(priorityCandidates),
      isPreviewMode,
    )

    const engagementPool = [...leadsPoolRes.data.items, ...accountsPoolRes.data.items]
    const kpiTrends = buildKpiTrends(leadsPoolRes.data.items, accountsPoolRes.data.items)
    const weeklyFollowUpCount = buildWeeklyFollowUpCount(
      leadsPoolRes.data.items,
      accountsPoolRes.data.items,
      isPreviewMode,
    )

    return {
      leadsTotal: pipelineLeadsRes.pagination.total,
      accountsTotal: pipelineAccountsRes.pagination.total,
      dealsTotal: 0,
      dealsOpenCount: 0,
      dealsOpenAmount: 0,
      atRiskTotal: atRiskLeadsRes.pagination.total + atRiskAccountsRes.pagination.total,
      avgEngagement: averageEngagement(engagementPool),
      priorityCount: priorities.length,
      priorities,
      pipelineLeads: pipelineLeadsRes.data.items,
      pipelineAccounts: pipelineAccountsRes.data.items,
      weeklyFollowUpCount,
      kpiTrends,
      sparklines: { leads: [], deals: [] },
      summaryFromApi: false,
    }
  }

  async function loadSnapshot(isPreviewMode = false): Promise<DashboardSnapshot> {
    const fromMock = dashboardStats.forceMock.value

    const emptyLeads = { data: { items: [] as Lead[] }, pagination: { page: 1, page_size: 0, total: 0 } }
    const emptyAccounts = { data: { items: [] as Account[] }, pagination: { page: 1, page_size: 0, total: 0 } }

    const [summary, pipelineLeadsRes, pipelineAccountsRes] = await Promise.all([
      dashboardStats.fetchSummary({ preview: isPreviewMode }),
      safeList(() => leadsApi.fetchList({ page: 1, page_size: 3 }), emptyLeads),
      safeList(() => accountsApi.fetchList({ page: 1, page_size: 3 }), emptyAccounts),
    ])

    let priorities = mergePriorities(summary.priorities.map(mapSummaryPriority))
    if (priorities.length === 0) {
      const legacy = await loadSnapshotLegacy(isPreviewMode)
      priorities = legacy.priorities
    } else {
      priorities = ensurePreviewPriorities(priorities, isPreviewMode)
    }

    return {
      leadsTotal: summary.kpis.leads_total,
      accountsTotal: summary.kpis.accounts_total,
      dealsTotal: summary.kpis.deals_total,
      dealsOpenCount: summary.kpis.deals_open_count,
      dealsOpenAmount: summary.kpis.deals_open_amount,
      atRiskTotal: summary.kpis.at_risk_total,
      avgEngagement: summary.kpis.avg_engagement,
      priorityCount: priorities.length,
      priorities,
      pipelineLeads: pipelineLeadsRes.data.items,
      pipelineAccounts: pipelineAccountsRes.data.items,
      weeklyFollowUpCount: summary.kpis.weekly_follow_ups,
      kpiTrends: summary.kpi_trends,
      sparklines: summary.sparklines,
      dataScope: summary.data_scope,
      canViewTeamRanking: summary.can_view_team_ranking ?? summary.data_scope !== 'self',
      summaryFromApi: !fromMock,
    }
  }

  function loadInsightItems(isPreview: boolean): DashboardInsightItem[] {
    const items: DashboardInsightItem[] = [
      {
        id: 'rule-1',
        variant: 'rule',
        title: t('dashboardInsightRuleTitle'),
        body: t('dashboardInsightRuleBody'),
      },
    ]
    if (isPreview) {
      items.unshift(
        {
          id: 'preview-opportunity',
          variant: 'opportunity',
          title: t('dashboardInsightOpportunityTitle'),
          body: t('dashboardInsightOpportunityBody'),
          isPreview: true,
        },
        {
          id: 'preview-churn',
          variant: 'churn',
          title: t('dashboardInsightChurnTitle'),
          body: t('dashboardInsightChurnBody'),
          isPreview: true,
          urgent: true,
        },
      )
    }
    return items
  }

  function ensurePreviewPriorities(
    priorities: PriorityActionItem[],
    isPreviewMode: boolean,
  ): PriorityActionItem[] {
    if (!isPreviewMode) return priorities
    const hasDemo = priorities.some(
      (item) => item.isPreview || item.entityId === DEMO_LEAD_ID,
    )
    if (hasDemo) return priorities

    const demoCard: PriorityActionItem = {
      id: 'preview-huachuang',
      entityType: 'lead',
      entityId: DEMO_LEAD_ID,
      title: '华创科技',
      reasons: [t('dashboardReasonNoFollowup', { days: 7 }), t('dashboardReasonHealthDeclining')],
      suggestion: t('dashboardSuggestCall'),
      insightTip: t('dashboardInsightTipHuachuang'),
      followHref: `/leads/${DEMO_LEAD_ID}#timeline`,
      urgency: 'coral',
      healthLabel: 'alert',
      isPreview: true,
      score: 99,
      engagementScore: 38,
      sparkline: [58, 48, 42, 38],
      contactName: '王总监',
      daysSinceActivity: 7,
    }
    return mergePriorities([demoCard, ...priorities])
  }

  return { loadSnapshot, loadInsightItems }
}
