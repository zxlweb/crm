import type {
  Contact,
  ContactCreateInput,
  ContactListQuery,
  ContactUpdateInput,
  ContactsListData,
  Pagination,
} from '~/types/contact'

type ContactsPagePayload = ContactsListData | Contact[]

function normalizeListPayload(data: ContactsPagePayload | undefined): Contact[] {
  if (!data) return []
  if (Array.isArray(data)) return data
  return data.items ?? []
}

function buildQueryString(query: ContactListQuery): string {
  const params = new URLSearchParams()
  if (query.page) params.set('page', String(query.page))
  if (query.page_size) params.set('page_size', String(query.page_size))
  if (query.search) params.set('search', query.search)
  if (query.lifecycle_stage) params.set('lifecycle_stage', query.lifecycle_stage)
  if (query.relationship_health) params.set('relationship_health', query.relationship_health)
  if (query.account_id) params.set('account_id', query.account_id)
  if (query.owner_id) params.set('owner_id', query.owner_id)
  const qs = params.toString()
  return qs ? `?${qs}` : ''
}

export function useContacts() {
  const api = useApi()

  async function fetchList(
    query: ContactListQuery = {},
  ): Promise<{ data: ContactsListData; pagination: Pagination }> {
    const path = `/api/contacts${buildQueryString(query)}`
    const res = await api.requestPage<ContactsPagePayload>(path)
    return {
      data: { items: normalizeListPayload(res.data) },
      pagination: {
        page: res.pagination.page,
        page_size: res.pagination.page_size,
        total: Number(res.pagination.total),
      },
    }
  }

  async function fetchByAccount(
    accountId: string,
    query: Omit<ContactListQuery, 'account_id'> = {},
  ): Promise<{ data: ContactsListData; pagination: Pagination }> {
    const params = new URLSearchParams()
    if (query.page) params.set('page', String(query.page))
    if (query.page_size) params.set('page_size', String(query.page_size))
    const qs = params.toString()
    const path = `/api/accounts/${accountId}/contacts${qs ? `?${qs}` : ''}`
    const res = await api.requestPage<ContactsPagePayload>(path)
    return {
      data: { items: normalizeListPayload(res.data) },
      pagination: {
        page: res.pagination.page,
        page_size: res.pagination.page_size,
        total: Number(res.pagination.total),
      },
    }
  }

  async function fetchById(id: string): Promise<Contact | null> {
    try {
      return await api.request<Contact>(`/api/contacts/${id}`)
    } catch {
      return null
    }
  }

  async function create(input: ContactCreateInput): Promise<Contact> {
    return api.request<Contact>('/api/contacts', {
      method: 'POST',
      body: JSON.stringify(input),
    })
  }

  async function update(id: string, input: ContactUpdateInput): Promise<Contact> {
    return api.request<Contact>(`/api/contacts/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(input),
    })
  }

  async function remove(id: string): Promise<void> {
    await api.request<void>(`/api/contacts/${id}`, { method: 'DELETE' })
  }

  return {
    fetchList,
    fetchByAccount,
    fetchById,
    create,
    update,
    remove,
  }
}
