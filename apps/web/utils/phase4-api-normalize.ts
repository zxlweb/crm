import type { AuditByAction, AuditTopActors, AuditTrend } from '~/types/audit-stats'
import type { CustomField } from '~/types/settings'

type LabeledRecord = Record<string, unknown>

function pickString(obj: LabeledRecord, ...keys: string[]): string {
  for (const key of keys) {
    const v = obj[key]
    if (typeof v === 'string') return v
  }
  return ''
}

function pickNumber(obj: LabeledRecord, ...keys: string[]): number {
  for (const key of keys) {
    const v = obj[key]
    if (typeof v === 'number') return v
    if (typeof v === 'string' && v !== '') return Number(v)
  }
  return 0
}

function parseJsonObject(raw: unknown): Record<string, string> {
  if (raw && typeof raw === 'object' && !Array.isArray(raw)) {
    const out: Record<string, string> = {}
    for (const [k, v] of Object.entries(raw as Record<string, unknown>)) {
      if (typeof v === 'string') out[k] = v
    }
    return out
  }
  if (typeof raw === 'string' && raw.trim()) {
    try {
      return parseJsonObject(JSON.parse(raw))
    } catch {
      return {}
    }
  }
  return {}
}

function parseOptions(raw: unknown): CustomField['options'] {
  if (Array.isArray(raw)) return raw as CustomField['options']
  if (typeof raw === 'string' && raw.trim()) {
    try {
      const parsed = JSON.parse(raw)
      return Array.isArray(parsed) ? parsed : []
    } catch {
      return []
    }
  }
  return []
}

export function normalizeCustomField(raw: unknown): CustomField {
  const row = (raw ?? {}) as LabeledRecord
  const label = parseJsonObject(row.field_label ?? row.FieldLabel)
  return {
    id: pickString(row, 'id', 'ID'),
    entity_type: pickString(row, 'entity_type', 'EntityType') as CustomField['entity_type'],
    field_key: pickString(row, 'field_key', 'FieldKey'),
    field_label: {
      'zh-CN': label['zh-CN'] ?? label['zh_CN'] ?? '',
      'en-US': label['en-US'] ?? label['en_US'] ?? '',
    },
    field_type: pickString(row, 'field_type', 'FieldType') as CustomField['field_type'],
    required: Boolean(row.required ?? row.Required),
    options: parseOptions(row.options ?? row.Options),
    default_value: (row.default_value ?? row.DefaultValue ?? null) as string | null,
    display_order: pickNumber(row, 'display_order', 'DisplayOrder') || 100,
    is_active: row.is_active !== false && row.IsActive !== false,
    created_at: pickString(row, 'created_at', 'CreatedAt'),
    updated_at: pickString(row, 'updated_at', 'UpdatedAt'),
  }
}

export function normalizeCustomFieldsList(data: unknown): CustomField[] {
  if (Array.isArray(data)) return data.map(normalizeCustomField)
  const wrapped = data as { items?: unknown[] } | null
  if (wrapped?.items && Array.isArray(wrapped.items)) {
    return wrapped.items.map(normalizeCustomField)
  }
  return []
}

export function normalizeAuditByAction(data: unknown): AuditByAction {
  const raw = (data ?? {}) as LabeledRecord
  const itemsRaw = (raw.items ?? raw.Items) as unknown[] | undefined
  const items = (itemsRaw ?? []).map((row) => {
    const r = (row ?? {}) as LabeledRecord
    return {
      action: pickString(r, 'action', 'Action'),
      count: pickNumber(r, 'count', 'Count'),
    }
  })
  const total = pickNumber(raw, 'total', 'Total') || items.reduce((s, i) => s + i.count, 0)
  return { items, total }
}

export function normalizeAuditTrend(data: unknown): AuditTrend {
  const raw = (data ?? {}) as LabeledRecord
  const itemsRaw = (raw.items ?? raw.Items) as unknown[] | undefined
  const items = (itemsRaw ?? []).map((row) => {
    const r = (row ?? {}) as LabeledRecord
    return {
      date: pickString(r, 'date', 'Date', 'bucket', 'Bucket'),
      count: pickNumber(r, 'count', 'Count'),
    }
  })
  const gran = pickString(raw, 'granularity', 'Granularity')
  return {
    items,
    granularity: gran === 'week' ? 'week' : 'day',
  }
}

export function normalizeAuditTopActors(data: unknown): AuditTopActors {
  const raw = (data ?? {}) as LabeledRecord
  const itemsRaw = (raw.items ?? raw.Items) as unknown[] | undefined
  const items = (itemsRaw ?? []).map((row) => {
    const r = (row ?? {}) as LabeledRecord
    const actorId = r.actor_id ?? r.ActorID
    return {
      actor_id: typeof actorId === 'string' ? actorId : '',
      actor_name: pickString(r, 'actor_name', 'ActorName') || '—',
      count: pickNumber(r, 'count', 'Count'),
    }
  })
  return { items }
}
