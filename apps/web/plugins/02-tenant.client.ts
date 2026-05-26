/**
 * 刷新后 tenant.list 会丢失；若 cookie 无租户则拉列表并自动选第一个。
 * 已登录但无 department cookie 时，用 refresh 补全事业部名称。
 */
import type { LoginResponse } from '~/composables/use-auth'

export default defineNuxtPlugin(async () => {
  const auth = useAuth()
  if (!auth.isAuthenticated.value) return

  const tenant = useTenant()
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase as string

  if (!tenant.currentTenantId.value || tenant.tenants.value.length === 0) {
    try {
      const list = await tenant.fetchTenants()
      if (list.length > 0) {
        const cookieValid = list.some((t) => t.id === tenant.currentTenantId.value)
        const targetId = cookieValid ? tenant.currentTenantId.value! : list[0].id
        await tenant.switchTenant(targetId)
      }
    } catch {
      // 顶栏切换器会提示用户手动选择
    }
  } else if (tenant.currentTenantId.value) {
    const cookieValid = tenant.tenants.value.some((t) => t.id === tenant.currentTenantId.value)
    if (!cookieValid && tenant.tenants.value.length > 0) {
      try {
        await tenant.switchTenant(tenant.tenants.value[0].id)
      } catch {
        // ignore
      }
    }
  }

  if (tenant.currentTenantId.value && !tenant.currentDepartment.value && auth.refreshToken.value) {
    try {
      const data = await $fetch<LoginResponse>(`${apiBase}/api/auth/refresh`, {
        method: 'POST',
        body: { refresh_token: auth.refreshToken.value },
      })
      auth.setSession(data.access_token, data.refresh_token, data.user)
      tenant.applyDepartmentFromLogin(data)
      useActiveRole().applyFromLogin(data)
    } catch {
      // 无 department 时不阻断；用户可重新登录
    }
  }
})
