import type { AiPreviewCopilotScene, AiPreviewInsightFixture } from '~/types/ai-preview'

/** §15.2 演示洞察 — title/body 走 i18n key */
export const AI_PREVIEW_INSIGHT_FIXTURES: AiPreviewInsightFixture[] = [
  { id: 'INS-P1', titleKey: 'aiStubInsight1Title', bodyKey: 'aiStubInsight1Body' },
  { id: 'INS-P2', titleKey: 'aiStubInsight2Title', bodyKey: 'aiStubInsight2Body' },
  { id: 'INS-P3', titleKey: 'aiPreviewChurnTitle', bodyKey: 'aiPreviewChurnBody' },
]

/** §15.2 步骤 5 — 流失风险样例分数 */
export const AI_PREVIEW_CHURN_RISK_SCORE = 72

export const AI_PREVIEW_COPILOT_KEYS: Record<AiPreviewCopilotScene, string> = {
  followup_script: 'aiCopilotMockFollowup',
  email_draft: 'aiCopilotMockEmail',
}
