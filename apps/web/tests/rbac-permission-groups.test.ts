import { describe, expect, it } from 'vitest'
import { buildPermissionModules } from '~/utils/rbac-permission-groups'

describe('buildPermissionModules', () => {
  it('groups permissions by module and resource', () => {
    const modules = buildPermissionModules([
      { id: '1', resource: 'leads', action: 'view', description: '' },
      { id: '2', resource: 'leads', action: 'create', description: '' },
      { id: '3', resource: 'settings', action: 'update', description: '' },
      { id: '4', resource: 'rbac', action: 'manage', description: '' },
    ])
    expect(modules.map((m) => m.key)).toEqual(['crm', 'settings', 'system'])
    expect(modules[0].resources[0].resource).toBe('leads')
    expect(modules[0].resources[0].items).toHaveLength(2)
  })
})
