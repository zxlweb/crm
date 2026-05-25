import type {
  InsightEvaluateResult,
  InsightItem,
  InsightPanelItem,
  InsightSubjectType,
} from '~/types/insight'

const RESOURCE_PATH: Record<InsightSubjectType, string> = {
  lead: 'leads',
  account: 'accounts',
  contact: 'contacts',
}

export function useInsights() {
  const api = useApi()

  async function evaluate(subjectType: InsightSubjectType, id: string): Promise<InsightEvaluateResult> {
    return api.request<InsightEvaluateResult>(
      `/api/${RESOURCE_PATH[subjectType]}/${id}/insights/evaluate`,
      { method: 'POST', body: JSON.stringify({}) },
    )
  }

  function toPanelItems(items: InsightItem[], te: (key: string) => string): InsightPanelItem[] {
    return items.map((item) => ({
      id: item.id,
      title: te(item.title_key),
      body: te(item.body_key),
    }))
  }

  return { evaluate, toPanelItems }
}
