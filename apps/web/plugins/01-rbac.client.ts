/**
 * 刷新后 useState 权限会丢失；登录时 loadMyPermissions 只在 finishSession 调用一次。
 * 客户端启动与路由切换时补拉权限，避免 PermissionGuard 空白页。
 */
export default defineNuxtPlugin(async () => {
  const auth = useAuth()
  const tenant = useTenant()
  const permission = usePermission()
  const rbac = useRbac()
  const hydratedForTenant = useState<string | null>('rbac.hydratedForTenant', () => null)

  async function hydratePermissions(force = false) {
    if (!auth.isAuthenticated.value || auth.isSuperAdmin.value) return
    const tid = tenant.currentTenantId.value
    if (!tid) return
    if (!force && hydratedForTenant.value === tid && Object.keys(permission.permissions.value).length > 0) {
      return
    }
    try {
      await rbac.loadMyPermissions()
      hydratedForTenant.value = tid
    } catch {
      // 不阻断导航；页面内 API 错误自行展示
    }
  }

  await hydratePermissions()

  watch(
    () => tenant.currentTenantId.value,
    (tid, prev) => {
      if (tid && tid !== prev) {
        hydratedForTenant.value = null
        void hydratePermissions(true)
      }
    },
  )

  const router = useRouter()
  router.afterEach(() => {
    void hydratePermissions()
  })
})
