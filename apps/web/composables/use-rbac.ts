import type { PermissionGroup } from '~/utils/permissions'

export type RoleItem = {
  id: string
  name: string
  description: string
  is_system: boolean
  permission_ids: string[]
  user_count: number
}

export type MemberRoleItem = {
  id: string
  name: string
  is_system: boolean
}

export type MemberItem = {
  id: string
  email: string
  name: string
  avatar_url: string
  department?: string
  roles: MemberRoleItem[]
  joined_at: string
}

export type PermissionDictItem = {
  id: string
  resource: string
  action: string
  description: string
}

export function useRbac() {
  const api = useApi()
  const permission = usePermission()

  async function loadMyPermissions() {
    const data = await api.request<PermissionGroup[]>('/api/rbac/my-permissions', { skipTenant: false })
    permission.setPermissions(data)
    return data
  }

  function fetchPermissionDictionary() {
    return api.request<PermissionGroup[]>('/api/rbac/permissions')
  }

  function fetchPermissionItems() {
    return api.request<PermissionDictItem[]>('/api/rbac/permission-items')
  }

  function fetchRoles() {
    return api.request<RoleItem[]>('/api/rbac/roles')
  }

  function fetchMembers() {
    return api.request<MemberItem[]>('/api/rbac/members')
  }

  function assignMemberRoles(userId: string, roleIds: string[]) {
    return api.request<RoleItem[]>(`/api/rbac/members/${userId}/roles`, {
      method: 'PUT',
      body: JSON.stringify({ role_ids: roleIds }),
    })
  }

  function createRole(payload: { name: string; description?: string }) {
    return api.request<RoleItem>('/api/rbac/roles', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  }

  function updateRole(id: string, payload: { name: string; description?: string }) {
    return api.request<RoleItem>(`/api/rbac/roles/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  }

  function assignRolePermissions(id: string, permissionIds: string[]) {
    return api.request<RoleItem>(`/api/rbac/roles/${id}/permissions`, {
      method: 'POST',
      body: JSON.stringify({ permission_ids: permissionIds }),
    })
  }

  function checkPermission(resource: string, action: string) {
    return api.request<{ allowed: boolean }>('/api/rbac/check', {
      method: 'POST',
      body: JSON.stringify({ resource, action }),
    })
  }

  return {
    loadMyPermissions,
    fetchPermissionDictionary,
    fetchPermissionItems,
    fetchRoles,
    fetchMembers,
    assignMemberRoles,
    createRole,
    updateRole,
    assignRolePermissions,
    checkPermission,
  }
}
