import type { LifecycleStage, RelationshipHealth } from '~/types/lead'
import type { InsightPanelItem, InsightSubjectType } from '~/types/insight'

type DetailInsightsOptions = {
  subjectType: InsightSubjectType
  subjectId: Ref<string | null | undefined>
  previewInsights?: Ref<InsightPanelItem[]>
  fallbackEngagement?: Ref<number | null | undefined>
}

export function useDetailInsights(opts: DetailInsightsOptions) {
  const insightsApi = useInsights()
  const { t } = useI18n()

  const insights = ref<InsightPanelItem[]>([])
  const engagementScore = ref<number | null>(null)
  const lifecycleStage = ref<LifecycleStage | null>(null)
  const relationshipHealth = ref<RelationshipHealth | null>(null)
  const pending = ref(false)

  async function reload() {
    const id = opts.subjectId.value
    if (!id) {
      insights.value = []
      engagementScore.value = opts.fallbackEngagement?.value ?? null
      return
    }

    const preview = opts.previewInsights?.value ?? []
    if (preview.length > 0) {
      insights.value = preview
      engagementScore.value = opts.fallbackEngagement?.value ?? null
      return
    }

    pending.value = true
    try {
      const data = await insightsApi.evaluate(opts.subjectType, id)
      insights.value = insightsApi.toPanelItems(data.items, t)
      engagementScore.value = data.engagement_score
      lifecycleStage.value = data.lifecycle_stage
      relationshipHealth.value = data.relationship_health
    } catch {
      insights.value = []
      engagementScore.value = opts.fallbackEngagement?.value ?? null
    } finally {
      pending.value = false
    }
  }

  watch(
    () => [opts.subjectId.value, opts.previewInsights?.value] as const,
    reload,
    { immediate: true, deep: true },
  )

  return { insights, engagementScore, lifecycleStage, relationshipHealth, pending, reload }
}
