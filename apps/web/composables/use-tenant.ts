import type { AuthTenant, LoginResponse } from '~/composables/use-auth'

const TENANT_COOKIE_OPTS = {
  maxAge: 60 * 60 * 24 * 7,
  sameSite: 'lax' as const,
  path: '/',
}

export function useTenant() {
  const auth = useAuth()
  const currentTenantId = useCookie<string | null>('crm.tenant_id', TENANT_COOKIE_OPTS)
  const tenants = useState<AuthTenant[]>('tenant.list', () => [])

  const currentTenant = computed(() =>
    tenants.value.find((t) => t.id === currentTenantId.value) ?? null,
  )

  function setTenant(id: string) {
    currentTenantId.value = id || null
  }

  function setTenantList(list: AuthTenant[]) {
    tenants.value = list
    if (!currentTenantId.value && list.length > 0) {
      currentTenantId.value = list[0].id
    }
  }

  function clearTenant() {
    currentTenantId.value = null
    tenants.value = []
  }

  /** 拉取当前用户可访问的租户（无需 X-Tenant-ID） */
  async function fetchTenants() {
    const list = await useApi().request<AuthTenant[]>('/api/auth/tenants', { skipTenant: true })
    setTenantList(list)
    return list
  }

  /** 切换租户：更新 token + cookie，并刷新 RBAC */
  async function switchTenant(tenantId: string) {
    const data = await useApi().request<LoginResponse>('/api/auth/switch-tenant', {
      method: 'POST',
      skipTenant: true,
      body: JSON.stringify({ tenant_id: tenantId }),
    })
    auth.setSession(data.access_token, data.refresh_token, data.user)
    setTenantList(data.tenants)
    setTenant(tenantId)
    if (!auth.isSuperAdmin.value) {
      try {
        await useRbac().loadMyPermissions()
      } catch {
        // 切换成功即可继续
      }
    }
    return data
  }

  return {
    currentTenantId,
    currentTenant,
    tenants,
    setTenant,
    setTenantList,
    clearTenant,
    fetchTenants,
    switchTenant,
  }
}
