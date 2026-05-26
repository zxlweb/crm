export interface TenantSettings {
  tenant_id: string
  tenant_name: string
  default_locale: string
  timezone: string
  business_switches: BusinessSwitches
  sales_quota: SalesQuota
  updated_at: string
  updated_by: string
}

export interface BusinessSwitches {
  ai_preview_enabled: boolean
  lead_import_mode: 'manual_review' | 'auto_merge'
}

export interface SalesQuota {
  amount: number
  currency: string
  period: string
}

export interface TenantSettingsUpdateInput {
  tenant_name?: string
  default_locale?: string
  timezone?: string
  business_switches?: Partial<BusinessSwitches>
  sales_quota?: Partial<SalesQuota>
}

export interface FeatureFlag {
  key: string
  type: 'boolean' | 'string' | 'number'
  value: boolean | string | number
  default_value: boolean | string | number
  description?: string
}

export type FieldType = 'text' | 'select' | 'date'
export type EntityType = 'account' | 'contact' | 'lead' | 'deal'

export interface CustomFieldOption {
  value: string
  label: { 'zh-CN': string; 'en-US': string }
}

export interface CustomField {
  id: string
  entity_type: EntityType
  field_key: string
  field_label: { 'zh-CN': string; 'en-US': string }
  field_type: FieldType
  required: boolean
  options: CustomFieldOption[]
  default_value: string | null
  display_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CustomFieldCreateInput {
  entity_type: EntityType
  field_key: string
  field_label: { 'zh-CN': string; 'en-US': string }
  field_type: FieldType
  required?: boolean
  options?: CustomFieldOption[]
  default_value?: string | null
  display_order?: number
}

export interface CustomFieldUpdateInput {
  field_label?: { 'zh-CN': string; 'en-US': string }
  required?: boolean
  options?: CustomFieldOption[]
  default_value?: string | null
  display_order?: number
  is_active?: boolean
}
