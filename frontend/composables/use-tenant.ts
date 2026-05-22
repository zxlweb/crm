export function useTenant() {
  const currentTenantId = useState<string | null>('tenant.id', () => null)
  const tenants = useState<Array<{ id: string; name: string; domain: string }>>('tenant.list', () => [])

  function setTenant(id: string) {
    currentTenantId.value = id
  }

  function setTenantList(list: Array<{ id: string; name: string; domain: string }>) {
    tenants.value = list
  }

  return { currentTenantId, tenants, setTenant, setTenantList }
}
