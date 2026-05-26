import { describe, expect, it } from 'vitest'
import {
  normalizeAuditByAction,
  normalizeAuditTopActors,
  normalizeAuditTrend,
  normalizeCustomField,
  normalizeCustomFieldsList,
} from '~/utils/phase4-api-normalize'

describe('phase4-api-normalize', () => {
  it('unwraps custom fields list', () => {
    const list = normalizeCustomFieldsList({
      items: [{
        id: 'cf-1',
        entity_type: 'lead',
        field_key: 'segment',
        field_label: { 'zh-CN': '子类', 'en-US': 'Segment' },
        field_type: 'text',
        required: false,
        is_active: true,
        created_at: '2026-05-01T00:00:00Z',
        updated_at: '2026-05-01T00:00:00Z',
      }],
    })
    expect(list).toHaveLength(1)
    expect(list[0].field_label['zh-CN']).toBe('子类')
  })

  it('normalizes PascalCase audit by-action', () => {
    const data = normalizeAuditByAction({
      Items: [{ Action: 'settings.update', Count: 3 }],
      Total: 3,
    })
    expect(data.items[0]).toEqual({ action: 'settings.update', count: 3 })
    expect(data.total).toBe(3)
  })

  it('maps trend bucket to date', () => {
    const data = normalizeAuditTrend({
      Items: [{ Bucket: '2026-05-20', Count: 12 }],
    })
    expect(data.items[0].date).toBe('2026-05-20')
    expect(data.items[0].count).toBe(12)
  })

  it('normalizes top actors', () => {
    const data = normalizeAuditTopActors({
      Items: [{ ActorName: 'Alice', Count: 5 }],
    })
    expect(data.items[0].actor_name).toBe('Alice')
    expect(data.items[0].count).toBe(5)
  })

  it('parses string field_label on custom field', () => {
    const field = normalizeCustomField({
      id: 'x',
      entity_type: 'lead',
      field_key: 'k',
      field_label: '{"zh-CN":"A","en-US":"B"}',
      field_type: 'text',
      required: false,
      is_active: true,
      created_at: '',
      updated_at: '',
    })
    expect(field.field_label['en-US']).toBe('B')
  })
})
