import { describe, expect, it } from 'vitest'
import { mockTenantSettings, mockFeatureFlags, mockCustomFields } from '~/fixtures/settings.mock'

describe('settings mock fixtures', () => {
  it('returns valid tenant settings', () => {
    const s = mockTenantSettings()
    expect(s.tenant_id).toBeTruthy()
    expect(s.tenant_name).toBe('Acme China')
    expect(s.default_locale).toBe('zh-CN')
    expect(s.timezone).toBe('Asia/Shanghai')
    expect(s.business_switches).toBeDefined()
    expect(typeof s.business_switches.ai_preview_enabled).toBe('boolean')
    expect(s.business_switches.lead_import_mode).toMatch(/^(manual_review|auto_merge)$/)
    expect(s.sales_quota.amount).toBeGreaterThanOrEqual(0)
    expect(s.sales_quota.currency).toBeTruthy()
  })

  it('returns feature flags with correct types', () => {
    const flags = mockFeatureFlags()
    expect(flags.length).toBeGreaterThan(0)
    for (const f of flags) {
      expect(f.key).toBeTruthy()
      expect(['boolean', 'string', 'number']).toContain(f.type)
    }
  })

  it('returns custom fields with required fields', () => {
    const fields = mockCustomFields()
    expect(fields.length).toBeGreaterThan(0)
    for (const f of fields) {
      expect(f.id).toBeTruthy()
      expect(['account', 'contact', 'lead', 'deal']).toContain(f.entity_type)
      expect(f.field_key).toBeTruthy()
      expect(f.field_label['zh-CN']).toBeTruthy()
      expect(f.field_label['en-US']).toBeTruthy()
      expect(['text', 'select', 'date']).toContain(f.field_type)
    }
  })

  it('select type fields have options', () => {
    const fields = mockCustomFields()
    const selectFields = fields.filter((f) => f.field_type === 'select')
    for (const f of selectFields) {
      expect(f.options.length).toBeGreaterThan(0)
      for (const opt of f.options) {
        expect(opt.value).toBeTruthy()
        expect(opt.label['zh-CN']).toBeTruthy()
        expect(opt.label['en-US']).toBeTruthy()
      }
    }
  })
})
