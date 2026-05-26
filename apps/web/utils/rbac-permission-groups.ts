import type { PermissionDictItem } from '~/composables/use-rbac'

/** Display order and module grouping for settings → role permissions tree */
export const PERMISSION_MODULE_ORDER = [
  'crm',
  'insights',
  'settings',
  'platform',
  'system',
] as const

export type PermissionModuleKey = (typeof PERMISSION_MODULE_ORDER)[number]

const RESOURCE_TO_MODULE: Record<string, PermissionModuleKey> = {
  leads: 'crm',
  accounts: 'crm',
  contacts: 'crm',
  deals: 'crm',
  dashboard: 'insights',
  settings: 'settings',
  custom_fields: 'settings',
  audit: 'settings',
  admin_tenant_insights: 'platform',
  rbac: 'system',
}

export type PermissionResourceNode = {
  resource: string
  items: PermissionDictItem[]
}

export type PermissionModuleNode = {
  key: PermissionModuleKey
  resources: PermissionResourceNode[]
}

export function buildPermissionModules(items: PermissionDictItem[]): PermissionModuleNode[] {
  const byModule = new Map<PermissionModuleKey, Map<string, PermissionDictItem[]>>()

  for (const item of items) {
    const mod = RESOURCE_TO_MODULE[item.resource] ?? 'system'
    if (!byModule.has(mod)) byModule.set(mod, new Map())
    const resMap = byModule.get(mod)!
    if (!resMap.has(item.resource)) resMap.set(item.resource, [])
    resMap.get(item.resource)!.push(item)
  }

  const modules: PermissionModuleNode[] = []
  for (const key of PERMISSION_MODULE_ORDER) {
    const resMap = byModule.get(key)
    if (!resMap) continue
    const resources: PermissionResourceNode[] = []
    for (const [resource, perms] of resMap) {
      perms.sort((a, b) => a.action.localeCompare(b.action))
      resources.push({ resource, items: perms })
    }
    resources.sort((a, b) => a.resource.localeCompare(b.resource))
    modules.push({ key, resources })
  }
  return modules
}

export function resourcePermissionIds(node: PermissionResourceNode): string[] {
  return node.items.map((p) => p.id)
}

export function modulePermissionIds(node: PermissionModuleNode): string[] {
  return node.resources.flatMap(resourcePermissionIds)
}
