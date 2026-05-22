import { DEMO_TENANT_ID } from '~/constants/demo'
import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'
import { mockLeadEmotionJourney } from '~/fixtures/emotion-journey.mock'
import type { EmotionJourney, EmotionSubjectType } from '~/types/emotion-journey'

export type EmotionJourneyQuery = {
  range?: '30d' | '90d' | 'all'
  useDemoFixture?: boolean
}

export function useEmotionJourney() {
  const api = useApi()
  const route = useRoute()

  const isPreview = computed(() => route.query.preview === '1')
  const tenantCookie = useCookie('crm.tenant_id')

  function shouldUseFixture(subjectId: string, query: EmotionJourneyQuery) {
    return (
      query.useDemoFixture === true ||
      isPreview.value ||
      subjectId === DEMO_LEAD_ID ||
      tenantCookie.value === DEMO_TENANT_ID
    )
  }

  async function fetchJourney(
    subjectType: EmotionSubjectType,
    subjectId: string,
    query: EmotionJourneyQuery = {},
  ): Promise<{ journey: EmotionJourney; fromFixture: boolean }> {
    const range = query.range ?? '90d'
    const useFixture = shouldUseFixture(subjectId, query)

    const fixture = () => {
      const mock = mockLeadEmotionJourney(subjectId)
      if (mock) return { journey: mock, fromFixture: true }
      return null
    }

    try {
      const data = await api.request<EmotionJourney>(
        `/api/${subjectType}s/${subjectId}/emotion-journey?range=${range}`,
      )
      if (useFixture && data.points.length === 0) {
        const mock = fixture()
        if (mock) return mock
      }
      return { journey: data, fromFixture: false }
    } catch (e) {
      if (useFixture) {
        const mock = fixture()
        if (mock) return mock
      }
      throw e
    }
  }

  return { fetchJourney, isPreview }
}
