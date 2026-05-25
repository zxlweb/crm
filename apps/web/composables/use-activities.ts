import type {
  Activity,
  ActivityCreateInput,
  ActivitySubjectType,
  ActivitySummary,
  ActivityUpdateInput,
  Pagination,
} from '~/types/activity'

function buildQuery(params: Record<string, string | number | undefined>): string {
  const q = new URLSearchParams()
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== '') q.set(k, String(v))
  }
  const s = q.toString()
  return s ? `?${s}` : ''
}

export function useActivities() {
  const api = useApi()

  async function fetchList(opts: {
    subjectType: ActivitySubjectType
    subjectId: string
    page?: number
    pageSize?: number
  }): Promise<{ items: Activity[]; pagination: Pagination }> {
    const qs = buildQuery({
      subject_type: opts.subjectType,
      subject_id: opts.subjectId,
      page: opts.page ?? 1,
      page_size: opts.pageSize ?? 50,
    })
    const res = await api.requestPage<{ items: Activity[] }>(`/api/activities${qs}`)
    return {
      items: res.data?.items ?? [],
      pagination: {
        page: res.pagination.page,
        page_size: res.pagination.page_size,
        total: Number(res.pagination.total),
      },
    }
  }

  async function fetchById(id: string): Promise<Activity | null> {
    try {
      return await api.request<Activity>(`/api/activities/${id}`)
    } catch {
      return null
    }
  }

  async function create(input: ActivityCreateInput): Promise<Activity> {
    return api.request<Activity>('/api/activities', {
      method: 'POST',
      body: JSON.stringify(input),
    })
  }

  async function update(id: string, input: ActivityUpdateInput): Promise<Activity> {
    return api.request<Activity>(`/api/activities/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(input),
    })
  }

  async function remove(id: string): Promise<void> {
    await api.request<void>(`/api/activities/${id}`, { method: 'DELETE' })
  }

  async function fetchSummary(opts: {
    subjectType?: ActivitySubjectType
    subjectId?: string
  }): Promise<ActivitySummary> {
    const qs = buildQuery({
      subject_type: opts.subjectType,
      subject_id: opts.subjectId,
    })
    return api.request<ActivitySummary>(`/api/activities/summary${qs}`)
  }

  return { fetchList, fetchById, create, update, remove, fetchSummary }
}
