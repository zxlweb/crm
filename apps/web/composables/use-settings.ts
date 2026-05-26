import { mockFeatureFlags, mockTenantSettings } from '~/fixtures/settings.mock'
import type {
  FeatureFlag,
  TenantSettings,
  TenantSettingsUpdateInput,
} from '~/types/settings'

export function useSettings() {
  const api = useApi()
  const config = useRuntimeConfig()
  const forceMock = computed(() => config.public.useSettingsMock === true || config.public.useSettingsMock === 'true')

  function normalizeTenantSettings(raw: TenantSettings): TenantSettings {
    const updatedBy = raw.updated_by
    return {
      ...raw,
      updated_by: updatedBy != null ? String(updatedBy) : '',
    }
  }

  async function fetchTenantSettings(): Promise<TenantSettings> {
    if (forceMock.value) return mockTenantSettings()
    const data = await api.request<TenantSettings>('/api/settings/tenant')
    return normalizeTenantSettings(data)
  }

  async function updateTenantSettings(input: TenantSettingsUpdateInput): Promise<TenantSettings> {
    if (forceMock.value) return { ...mockTenantSettings(), ...input }
    const data = await api.request<TenantSettings>('/api/settings/tenant', {
      method: 'PATCH',
      body: JSON.stringify(input),
    })
    return normalizeTenantSettings(data)
  }

  async function fetchFeatures(): Promise<FeatureFlag[]> {
    if (forceMock.value) return mockFeatureFlags()
    return api.request<FeatureFlag[]>('/api/settings/features')
  }

  return { forceMock, fetchTenantSettings, updateTenantSettings, fetchFeatures }
}
