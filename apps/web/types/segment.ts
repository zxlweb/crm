/** Phase 2 — docs/api/phase-2-crm-ai.md §2.5, §9 */

export const SEGMENT_CODES = [
  'high_value',
  'churn_risk',
  'new_potential',
  'needs_activation',
  'revive_pool',
] as const

export type SegmentCode = (typeof SEGMENT_CODES)[number]

export type SegmentTemplate = {
  code: SegmentCode
  name_key: string
  description_key: string
  filter?: Record<string, unknown>
  count?: number
}

export type SegmentCountResult = {
  count: number
}

export function isSegmentCode(value: string): value is SegmentCode {
  return (SEGMENT_CODES as readonly string[]).includes(value)
}
