import type {
  Account,
  AccountCreateInput,
  AccountListQuery,
  AccountUpdateInput,
  AccountsListData,
  Pagination,
} from '~/types/account'

type AccountsPagePayload = AccountsListData | Account[]

function normalizeListPayload(data: AccountsPagePayload | undefined): Account[] {
  if (!data) return []
  if (Array.isArray(data)) return data
  return data.items ?? []
}

function buildQueryString(query: AccountListQuery): string {
  const params = new URLSearchParams()
  if (query.page) params.set('page', String(query.page))
  if (query.page_size) params.set('page_size', String(query.page_size))
  if (query.search) params.set('search', query.search)
  if (query.lifecycle_stage) params.set('lifecycle_stage', query.lifecycle_stage)
  if (query.relationship_health) params.set('relationship_health', query.relationship_health)
  if (query.segment) params.set('segment', query.segment)
  if (query.owner_id) params.set('owner_id', query.owner_id)
  const qs = params.toString()
  return qs ? `?${qs}` : ''
}

export function useAccounts() {
  const api = useApi()

  async function fetchList(
    query: AccountListQuery = {},
  ): Promise<{ data: AccountsListData; pagination: Pagination }> {
    const path = `/api/accounts${buildQueryString(query)}`
    const res = await api.requestPage<AccountsPagePayload>(path)
    return {
      data: { items: normalizeListPayload(res.data) },
      pagination: {
        page: res.pagination.page,
        page_size: res.pagination.page_size,
        total: Number(res.pagination.total),
      },
    }
  }

  async function fetchById(id: string): Promise<Account | null> {
    try {
      return await api.request<Account>(`/api/accounts/${id}`)
    } catch {
      return null
    }
  }

  async function create(input: AccountCreateInput): Promise<Account> {
    return api.request<Account>('/api/accounts', {
      method: 'POST',
      body: JSON.stringify(input),
    })
  }

  async function update(id: string, input: AccountUpdateInput): Promise<Account> {
    return api.request<Account>(`/api/accounts/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(input),
    })
  }

  async function remove(id: string): Promise<void> {
    await api.request<void>(`/api/accounts/${id}`, { method: 'DELETE' })
  }

  return {
    fetchList,
    fetchById,
    create,
    update,
    remove,
  }
}
