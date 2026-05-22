import { canAccess, toPermissionMap, type PermissionMap } from '~/utils/permissions'

export function usePermission() {
  const permissions = useState<PermissionMap>('rbac.permissions', () => ({}))

  function setPermissions(data: Array<{ resource: string; actions: string[] }>) {
    permissions.value = toPermissionMap(data)
  }

  function can(resource: string, action: string): boolean {
    if (useAuth().isSuperAdmin.value) return true
    return canAccess(permissions.value, resource, action)
  }

  return { permissions, setPermissions, can }
}
