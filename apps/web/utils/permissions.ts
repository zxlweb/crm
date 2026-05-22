export type PermissionMap = Record<string, string[]>

export function toPermissionMap(
  data: Array<{ resource: string; actions: string[] }>,
): PermissionMap {
  const map: PermissionMap = {}
  for (const item of data) {
    map[item.resource] = item.actions
  }
  return map
}

export function canAccess(
  permissions: PermissionMap,
  resource: string,
  action: string,
): boolean {
  return permissions[resource]?.includes(action) ?? false
}
