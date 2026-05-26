import type {
  CustomField,
  FeatureFlag,
  TenantSettings,
} from '~/types/settings'

export function mockTenantSettings(): TenantSettings {
  return {
    tenant_id: 'demo-tenant-001',
    tenant_name: 'Acme China',
    default_locale: 'zh-CN',
    timezone: 'Asia/Shanghai',
    business_switches: {
      ai_preview_enabled: false,
      lead_import_mode: 'manual_review',
    },
    sales_quota: {
      amount: 5000000,
      currency: 'CNY',
      period: '2026-05',
    },
    updated_at: '2026-05-25T14:00:00Z',
    updated_by: 'user-001',
  }
}

export function mockFeatureFlags(): FeatureFlag[] {
  return [
    { key: 'ai_preview_enabled', type: 'boolean', value: false, default_value: false, description: 'AI Preview 总开关' },
    { key: 'lead_import_mode', type: 'string', value: 'manual_review', default_value: 'manual_review', description: '线索导入策略' },
  ]
}

export function mockCustomFields(): CustomField[] {
  return [
    {
      id: 'cf-001',
      entity_type: 'lead',
      field_key: 'industry_segment',
      field_label: { 'zh-CN': '行业子类', 'en-US': 'Industry Segment' },
      field_type: 'select',
      required: false,
      options: [
        { value: 'saas', label: { 'zh-CN': 'SaaS', 'en-US': 'SaaS' } },
        { value: 'manufacturing', label: { 'zh-CN': '制造业', 'en-US': 'Manufacturing' } },
        { value: 'fintech', label: { 'zh-CN': '金融科技', 'en-US': 'FinTech' } },
      ],
      default_value: null,
      display_order: 30,
      is_active: true,
      created_at: '2026-05-20T10:00:00Z',
      updated_at: '2026-05-20T10:00:00Z',
    },
    {
      id: 'cf-002',
      entity_type: 'account',
      field_key: 'contract_start',
      field_label: { 'zh-CN': '合同起始日', 'en-US': 'Contract Start' },
      field_type: 'date',
      required: false,
      options: [],
      default_value: null,
      display_order: 40,
      is_active: true,
      created_at: '2026-05-21T10:00:00Z',
      updated_at: '2026-05-21T10:00:00Z',
    },
    {
      id: 'cf-003',
      entity_type: 'deal',
      field_key: 'deal_memo',
      field_label: { 'zh-CN': '备注说明', 'en-US': 'Deal Memo' },
      field_type: 'text',
      required: false,
      options: [],
      default_value: '',
      display_order: 50,
      is_active: true,
      created_at: '2026-05-22T10:00:00Z',
      updated_at: '2026-05-22T10:00:00Z',
    },
  ]
}
