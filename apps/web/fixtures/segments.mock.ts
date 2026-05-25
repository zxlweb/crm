import type { SegmentTemplate } from '~/types/segment'

/** 契约预置五群 — 与 phase-2-crm-ai.md §2.5 一致；API 未就绪时 fallback */
export const SEGMENT_TEMPLATES_MOCK: SegmentTemplate[] = [
  {
    code: 'high_value',
    name_key: 'segment.high_value',
    description_key: 'segment.high_value.desc',
  },
  {
    code: 'churn_risk',
    name_key: 'segment.churn_risk',
    description_key: 'segment.churn_risk.desc',
  },
  {
    code: 'new_potential',
    name_key: 'segment.new_potential',
    description_key: 'segment.new_potential.desc',
  },
  {
    code: 'needs_activation',
    name_key: 'segment.needs_activation',
    description_key: 'segment.needs_activation.desc',
  },
  {
    code: 'revive_pool',
    name_key: 'segment.revive_pool',
    description_key: 'segment.revive_pool.desc',
  },
]
