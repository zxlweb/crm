export type AiPreviewCopilotScene = 'followup_script' | 'email_draft'

export type AiPreviewInsightFixture = {
  id: string
  titleKey: string
  bodyKey: string
}

export type AiPreviewPack = {
  insights: Array<{ id: string; title: string; body: string }>
  churnRiskScore: number
  disclaimer: string
}
