/** 全角色统一的退出登录（清空会话、租户、权限缓存） */
export function useLogout() {
  const auth = useAuth()
  const tenant = useTenant()
  const activeRole = useActiveRole()
  const { setPermissions } = usePermission()

  async function logout() {
    auth.clearTokens()
    tenant.clearTenant()
    activeRole.clearRoles()
    setPermissions([])
    await navigateTo('/login')
  }

  return { logout }
}
