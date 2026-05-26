import { mockCustomFields } from '~/fixtures/settings.mock'
import type {
  CustomField,
  CustomFieldCreateInput,
  CustomFieldUpdateInput,
  EntityType,
} from '~/types/settings'
import { normalizeCustomField, normalizeCustomFieldsList } from '~/utils/phase4-api-normalize'

export function useCustomFields() {
  const api = useApi()
  const config = useRuntimeConfig()
  const forceMock = computed(() => config.public.useSettingsMock === true || config.public.useSettingsMock === 'true')

  let _mockStore: CustomField[] | null = null
  function getMockStore(): CustomField[] {
    if (!_mockStore) _mockStore = [...mockCustomFields()]
    return _mockStore
  }

  async function fetchList(entityType?: EntityType): Promise<CustomField[]> {
    if (forceMock.value) {
      const store = getMockStore()
      return entityType ? store.filter((f) => f.entity_type === entityType && f.is_active) : store.filter((f) => f.is_active)
    }
    const q = entityType ? `?entity_type=${entityType}` : ''
    const data = await api.request<unknown>(`/api/settings/custom-fields${q}`)
    return normalizeCustomFieldsList(data)
  }

  async function create(input: CustomFieldCreateInput): Promise<CustomField> {
    if (forceMock.value) {
      const store = getMockStore()
      const existing = store.find((f) => f.entity_type === input.entity_type && f.field_key === input.field_key && f.is_active)
      if (existing) throw new Error('custom_field_key_conflict')
      const now = new Date().toISOString()
      const field: CustomField = {
        id: `cf-${Date.now()}`,
        entity_type: input.entity_type,
        field_key: input.field_key,
        field_label: input.field_label,
        field_type: input.field_type,
        required: input.required ?? false,
        options: input.options ?? [],
        default_value: input.default_value ?? null,
        display_order: input.display_order ?? 100,
        is_active: true,
        created_at: now,
        updated_at: now,
      }
      store.push(field)
      return field
    }
    const data = await api.request<unknown>('/api/settings/custom-fields', {
      method: 'POST',
      body: JSON.stringify(input),
    })
    return normalizeCustomField(data)
  }

  async function update(id: string, input: CustomFieldUpdateInput): Promise<CustomField> {
    if (forceMock.value) {
      const store = getMockStore()
      const idx = store.findIndex((f) => f.id === id)
      if (idx < 0) throw new Error('Not found')
      store[idx] = { ...store[idx], ...input, updated_at: new Date().toISOString() }
      return store[idx]
    }
    const data = await api.request<unknown>(`/api/settings/custom-fields/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(input),
    })
    return normalizeCustomField(data)
  }

  async function remove(id: string): Promise<void> {
    if (forceMock.value) {
      const store = getMockStore()
      const idx = store.findIndex((f) => f.id === id)
      if (idx >= 0) store[idx] = { ...store[idx], is_active: false }
      return
    }
    await api.request<void>(`/api/settings/custom-fields/${id}`, { method: 'DELETE' })
  }

  return { forceMock, fetchList, create, update, remove }
}
