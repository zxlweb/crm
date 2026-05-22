type ApiOptions = RequestInit & {
  skipTenant?: boolean
}

const TENANT_COOKIE_OPTS = {
  maxAge: 60 * 60 * 24 * 7,
  sameSite: 'lax' as const,
  path: '/',
}

function applyAuthHeaders(headers: Headers, options: ApiOptions) {
  headers.set('Content-Type', 'application/json')
  const token = useAuth().accessToken.value
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }
  if (!options.skipTenant) {
    const tenantId = useCookie<string | null>('crm.tenant_id', TENANT_COOKIE_OPTS).value
    if (tenantId) {
      headers.set('X-Tenant-ID', tenantId)
    }
  }
}

export function useApi() {
  const config = useRuntimeConfig()

  async function request<T>(path: string, options: ApiOptions = {}): Promise<T> {
    const headers = new Headers(options.headers)
    applyAuthHeaders(headers, options)

    const res = await fetch(`${config.public.apiBase}${path}`, {
      ...options,
      headers,
    })

    const body = await parseResponseBody(res)
    if (!res.ok) {
      throw new Error(body.message || res.statusText || `HTTP ${res.status}`)
    }
    return body.data as T
  }

  async function requestPage<T>(path: string, options: ApiOptions = {}): Promise<{ data: T; pagination: { page: number; page_size: number; total: number } }> {
    const headers = new Headers(options.headers)
    applyAuthHeaders(headers, options)

    const res = await fetch(`${config.public.apiBase}${path}`, {
      ...options,
      headers,
    })

    const body = await parseResponseBody(res)
    if (!res.ok) {
      throw new Error(body.message || res.statusText || `HTTP ${res.status}`)
    }
    return { data: body.data as T, pagination: body.pagination }
  }

  return { request, requestPage }
}

type ApiBody = {
  code?: number
  message?: string
  data?: unknown
  pagination?: { page: number; page_size: number; total: number }
}

async function parseResponseBody(res: Response): Promise<ApiBody> {
  const text = await res.text()
  if (!text.trim()) {
    return { message: res.statusText || `HTTP ${res.status}` }
  }
  try {
    return JSON.parse(text) as ApiBody
  } catch {
    throw new Error(text.trim().slice(0, 120) || res.statusText || `HTTP ${res.status}`)
  }
}
