import type { AuthRole, LoginResponse } from '~/composables/use-auth'

const ROLE_COOKIE_OPTS = {
  maxAge: 60 * 60 * 24 * 7,
  sameSite: 'lax' as const,
  path: '/',
}

export function useActiveRole() {
  const auth = useAuth()
  const currentRoleId = useCookie<string | null>('crm.role_id', ROLE_COOKIE_OPTS)
  const roles = useState<AuthRole[]>('role.list', () => [])

  const currentRole = computed(() =>
    roles.value.find((r) => r.id === currentRoleId.value) ?? null,
  )

  const canSwitchRole = computed(
    () => !auth.isSuperAdmin.value && roles.value.length > 1,
  )

  function applyFromLogin(data: Pick<LoginResponse, 'roles' | 'current_role'>) {
    if (data.roles?.length) {
      roles.value = data.roles
    }
    if (data.current_role) {
      currentRoleId.value = data.current_role.id
    } else if (data.roles?.length) {
      currentRoleId.value = data.roles[0].id
    } else {
      currentRoleId.value = null
    }
  }

  function clearRoles() {
    roles.value = []
    currentRoleId.value = null
  }

  async function fetchRoles() {
    const list = await useApi().request<AuthRole[]>('/api/rbac/my-roles')
    roles.value = list
    if (!currentRoleId.value && list.length > 0) {
      currentRoleId.value = list[0].id
    }
    return list
  }

  async function switchRole(roleId: string) {
    const data = await useApi().request<LoginResponse>('/api/auth/switch-role', {
      method: 'POST',
      body: JSON.stringify({ role_id: roleId }),
    })
    auth.setSession(data.access_token, data.refresh_token, data.user)
    useTenant().applyDepartmentFromLogin(data)
    applyFromLogin(data)
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
    currentRoleId,
    currentRole,
    roles,
    canSwitchRole,
    applyFromLogin,
    clearRoles,
    fetchRoles,
    switchRole,
  }
}
