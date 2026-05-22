import { describe, expect, it } from 'vitest'
import { canAccess, toPermissionMap } from './permissions'

describe('permissions utils', () => {
  it('toPermissionMap groups actions by resource', () => {
    const map = toPermissionMap([
      { resource: 'leads', actions: ['view', 'create'] },
      { resource: 'deals', actions: ['view'] },
    ])
    expect(map.leads).toEqual(['view', 'create'])
    expect(map.deals).toEqual(['view'])
  })

  it('canAccess returns true when action allowed', () => {
    const map = toPermissionMap([
      { resource: 'leads', actions: ['view'] },
    ])
    expect(canAccess(map, 'leads', 'view')).toBe(true)
    expect(canAccess(map, 'leads', 'delete')).toBe(false)
    expect(canAccess(map, 'deals', 'view')).toBe(false)
  })
})
