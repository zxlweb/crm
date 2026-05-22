type ApiOptions = RequestInit & {
  skipTenant?: boolean
}

export function useApi() {
  const config = useRuntimeConfig()
  const tenant = useTenant()

  async function request<T>(path: string, options: ApiOptions = {}): Promise<T> {
    const headers = new Headers(options.headers)
    headers.set('Content-Type', 'application/json')

    const token = useAuth().accessToken.value
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
    if (!options.skipTenant && tenant.currentTenantId.value) {
      headers.set('X-Tenant-ID', tenant.currentTenantId.value)
    }

    const res = await fetch(`${config.public.apiBase}${path}`, {
      ...options,
      headers,
    })

    const body = await res.json()
    if (!res.ok) {
      throw new Error(body.message || 'Request failed')
    }
    return body.data as T
  }

  async function requestPage<T>(path: string, options: ApiOptions = {}): Promise<{ data: T; pagination: { page: number; page_size: number; total: number } }> {
    const headers = new Headers(options.headers)
    headers.set('Content-Type', 'application/json')

    const token = useAuth().accessToken.value
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }

    const res = await fetch(`${config.public.apiBase}${path}`, {
      ...options,
      headers,
    })

    const body = await res.json()
    if (!res.ok) {
      throw new Error(body.message || 'Request failed')
    }
    return { data: body.data as T, pagination: body.pagination }
  }

  return { request, requestPage }
}
