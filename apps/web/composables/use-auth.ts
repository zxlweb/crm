export type AuthUser = {
  id: string
  email: string
  name: string
  is_super_admin: boolean
}

export type AuthTenant = {
  id: string
  name: string
  domain: string
}

export type AuthRole = {
  id: string
  name: string
  description: string
}

export type LoginResponse = {
  access_token: string
  refresh_token: string
  expires_in: number
  user: AuthUser
  tenants: AuthTenant[]
  current_tenant?: AuthTenant
  roles?: AuthRole[]
  current_role?: AuthRole
}

const AUTH_COOKIE_OPTS = {
  maxAge: 60 * 60 * 24 * 7,
  sameSite: 'lax' as const,
  path: '/',
}

export function useAuth() {
  const accessToken = useCookie<string | null>('crm.access_token', AUTH_COOKIE_OPTS)
  const refreshToken = useCookie<string | null>('crm.refresh_token', AUTH_COOKIE_OPTS)
  const user = useCookie<AuthUser | null>('crm.user', {
    ...AUTH_COOKIE_OPTS,
    default: () => null,
  })

  function setSession(access: string, refresh: string, profile: AuthUser) {
    accessToken.value = access
    refreshToken.value = refresh
    user.value = profile
  }

  function clearTokens() {
    accessToken.value = null
    refreshToken.value = null
    user.value = null
  }

  const isAuthenticated = computed(() => !!accessToken.value)
  const isSuperAdmin = computed(() => !!user.value?.is_super_admin)

  return {
    accessToken,
    refreshToken,
    user,
    setSession,
    clearTokens,
    isAuthenticated,
    isSuperAdmin,
  }
}
