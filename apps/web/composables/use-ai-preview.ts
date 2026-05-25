import { DEMO_TENANT_ID } from '~/constants/demo'
import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'
import {
  AI_PREVIEW_CHURN_RISK_SCORE,
  AI_PREVIEW_COPILOT_KEYS,
  AI_PREVIEW_INSIGHT_FIXTURES,
} from '~/fixtures/ai-preview'
import type { AiPreviewCopilotScene } from '~/types/ai-preview'
import type { InsightPanelItem } from '~/types/insight'

export function useAiPreview(subjectId?: Ref<string | null | undefined>) {
  const route = useRoute()
  const tenantCookie = useCookie('crm.tenant_id')
  const { t } = useI18n()

  const isPreviewMode = computed(() => {
    if (route.query.preview === '1') return true
    if (tenantCookie.value === DEMO_TENANT_ID) return true
    const id = subjectId?.value ?? (route.params.id as string | undefined)
    return id === DEMO_LEAD_ID
  })

  const previewInsights = computed<InsightPanelItem[]>(() => {
    if (!isPreviewMode.value) return []
    return AI_PREVIEW_INSIGHT_FIXTURES.map((item) => ({
      id: item.id,
      title: t(item.titleKey),
      body: t(item.bodyKey),
    }))
  })

  const churnRiskScore = computed(() => (isPreviewMode.value ? AI_PREVIEW_CHURN_RISK_SCORE : null))

  const disclaimer = computed(() => (isPreviewMode.value ? t('aiPreviewDisclaimer') : ''))

  function generateCopilot(scene: AiPreviewCopilotScene): string {
    return t(AI_PREVIEW_COPILOT_KEYS[scene])
  }

  function mapCopilotKind(kind: 'followup' | 'email'): AiPreviewCopilotScene {
    return kind === 'followup' ? 'followup_script' : 'email_draft'
  }

  return {
    isPreviewMode,
    previewInsights,
    churnRiskScore,
    disclaimer,
    generateCopilot,
    mapCopilotKind,
  }
}
