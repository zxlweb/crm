import { SEGMENT_TEMPLATES_MOCK } from '~/fixtures/segments.mock'
import type { SegmentCountResult, SegmentTemplate } from '~/types/segment'

type SegmentsListPayload = SegmentTemplate[] | { items: SegmentTemplate[] }

function normalizeList(data: SegmentsListPayload | undefined): SegmentTemplate[] {
  if (!data) return []
  if (Array.isArray(data)) return data
  return data.items ?? []
}

export function useSegments() {
  const api = useApi()

  async function fetchList(): Promise<SegmentTemplate[]> {
    try {
      const data = await api.request<SegmentsListPayload>('/api/segments')
      const items = normalizeList(data)
      return items.length > 0 ? items : SEGMENT_TEMPLATES_MOCK.map((s) => ({ ...s }))
    } catch {
      return SEGMENT_TEMPLATES_MOCK.map((s) => ({ ...s }))
    }
  }

  async function fetchCount(code: string): Promise<number> {
    try {
      const data = await api.request<SegmentCountResult>(
        `/api/segments/${encodeURIComponent(code)}/count`,
      )
      return data.count
    } catch {
      return 0
    }
  }

  async function fetchListWithCounts(): Promise<SegmentTemplate[]> {
    const list = await fetchList()
    return Promise.all(
      list.map(async (segment) => ({
        ...segment,
        count: segment.count ?? (await fetchCount(segment.code)),
      })),
    )
  }

  return { fetchList, fetchCount, fetchListWithCounts }
}
