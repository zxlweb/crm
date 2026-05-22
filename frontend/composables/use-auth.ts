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

export type LoginResponse = {
  access_token: string
  refresh_token: string
  expires_in: number
  user: AuthUser
  tenants: AuthTenant[]
  current_tenant?: AuthTenant
}

export function useAuth() {
  const accessToken = useState<string | null>('auth.accessToken', () => null)
  const refreshToken = useState<string | null>('auth.refreshToken', () => null)
  const user = useState<AuthUser | null>('auth.user', () => null)

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
