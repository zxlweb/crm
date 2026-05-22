/**
 * 刷新后 tenant.list 会丢失；若 cookie 无租户则拉列表并自动选第一个。
 */
export default defineNuxtPlugin(async () => {
  const auth = useAuth()
  if (!auth.isAuthenticated.value) return

  const tenant = useTenant()
  if (tenant.currentTenantId.value && tenant.tenants.value.length > 0) return

  try {
    const list = await tenant.fetchTenants()
    if (!tenant.currentTenantId.value && list.length > 0) {
      await tenant.switchTenant(list[0].id)
    }
  } catch {
    // 顶栏切换器会提示用户手动选择
  }
})
