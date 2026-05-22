import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'
import type { EmotionJourney } from '~/types/emotion-journey'

const now = Date.now()
const at = (daysAgo: number, hour = 10) =>
  new Date(now - daysAgo * 86400000 + hour * 3600000).toISOString()

/** PRD §15.2 演示路径 — 华创科技情绪触点 */
export const DEMO_LEAD_EMOTION_JOURNEY: EmotionJourney = {
  subject_type: 'lead',
  subject_id: DEMO_LEAD_ID,
  lifecycle_current: 'grow',
  lifecycle_bands: [
    { stage: 'acquire', from: at(45), to: at(30) },
    { stage: 'activate', from: at(30), to: at(18) },
    { stage: 'grow', from: at(18), to: at(0) },
  ],
  points: [
    {
      activity_id: 'act-001',
      at: at(28, 11),
      event_type: 'email',
      direction: 'outbound',
      sentiment: 'positive',
      sentiment_score: 2,
      sentiment_source: 'manual',
      label: '邮件：产品方案发送',
      lifecycle_stage_at_time: 'activate',
    },
    {
      activity_id: 'act-002',
      at: at(21, 15),
      event_type: 'call',
      direction: 'inbound',
      sentiment: 'neutral',
      sentiment_score: 0,
      sentiment_source: 'manual',
      label: '电话：预算范围确认',
      lifecycle_stage_at_time: 'grow',
    },
    {
      activity_id: 'act-003',
      at: at(14, 10),
      event_type: 'meeting',
      direction: 'outbound',
      sentiment: 'hesitant',
      sentiment_score: -1,
      sentiment_source: 'rule',
      label: '会议：价格与交付周期',
      lifecycle_stage_at_time: 'grow',
    },
    {
      activity_id: 'act-004',
      at: at(7, 16),
      event_type: 'wechat',
      direction: 'inbound',
      sentiment: 'negative',
      sentiment_score: -2,
      sentiment_source: 'manual',
      label: '微信：竞品对比顾虑',
      lifecycle_stage_at_time: 'grow',
    },
    {
      activity_id: 'act-005',
      at: at(3, 9),
      event_type: 'call',
      direction: 'outbound',
      sentiment: 'hesitant',
      sentiment_score: -1,
      sentiment_source: 'manual',
      label: '电话：ROI 材料跟进',
      lifecycle_stage_at_time: 'grow',
    },
  ],
  milestones: [
    { type: 'qualified', at: at(25), label: '标记为合格线索' },
    { type: 'insight', at: at(14), label: '洞察：价格顾虑' },
  ],
  summary: {
    current_sentiment: 'hesitant',
    trend: 'down',
    days_since_positive: 21,
  },
}

export const DEMO_LEAD_IDS = {
  huachuang: DEMO_LEAD_ID,
  yunfan: 'a1000000-0000-4000-8000-000000000002',
  beichen: 'a1000000-0000-4000-8000-000000000003',
  xinghai: 'a1000000-0000-4000-8000-000000000004',
  mingde: 'a1000000-0000-4000-8000-000000000005',
} as const

const LEAD_002 = DEMO_LEAD_IDS.yunfan
const LEAD_003 = DEMO_LEAD_IDS.beichen
const LEAD_004 = DEMO_LEAD_IDS.xinghai
const LEAD_005 = DEMO_LEAD_IDS.mingde

const LEAD_002_JOURNEY: EmotionJourney = {
  subject_type: 'lead',
  subject_id: LEAD_002,
  lifecycle_current: 'activate',
  lifecycle_bands: [
    { stage: 'acquire', from: at(20), to: at(10) },
    { stage: 'activate', from: at(10), to: at(0) },
  ],
  points: [
    {
      activity_id: 'f1000000-0000-4000-8000-000000000010',
      at: at(5, 10),
      event_type: 'email',
      direction: 'outbound',
      sentiment: 'positive',
      sentiment_score: 2,
      label: '邮件：白皮书',
      lifecycle_stage_at_time: 'activate',
    },
    {
      activity_id: 'f1000000-0000-4000-8000-000000000011',
      at: at(1, 14),
      event_type: 'call',
      direction: 'inbound',
      sentiment: 'positive',
      sentiment_score: 2,
      label: '电话：预约演示',
      lifecycle_stage_at_time: 'activate',
    },
  ],
  milestones: [],
  summary: { current_sentiment: 'positive', trend: 'up', days_since_positive: 1 },
}

const LEAD_003_JOURNEY: EmotionJourney = {
  subject_type: 'lead',
  subject_id: LEAD_003,
  lifecycle_current: 'acquire',
  lifecycle_bands: [{ stage: 'acquire', from: at(8), to: at(0) }],
  points: [
    {
      activity_id: 'act-nb-1',
      at: at(12, 9),
      event_type: 'email',
      direction: 'outbound',
      sentiment: 'neutral',
      sentiment_score: 0,
      label: '展会后首次触达',
      lifecycle_stage_at_time: 'acquire',
    },
    {
      activity_id: 'act-nb-2',
      at: at(8, 14),
      event_type: 'call',
      direction: 'outbound',
      sentiment: 'hesitant',
      sentiment_score: -1,
      label: '电话：需求尚不明确',
      lifecycle_stage_at_time: 'acquire',
    },
  ],
  milestones: [],
  summary: { current_sentiment: 'hesitant', trend: 'flat', days_since_positive: 30 },
}

/** 星海教育 — 低意向、不合格路径 */
const LEAD_004_JOURNEY: EmotionJourney = {
  subject_type: 'lead',
  subject_id: LEAD_004,
  lifecycle_current: 'acquire',
  lifecycle_bands: [{ stage: 'acquire', from: at(60), to: at(0) }],
  points: [
    {
      activity_id: 'act-xh-1',
      at: at(55, 10),
      event_type: 'call',
      direction: 'outbound',
      sentiment: 'neutral',
      sentiment_score: 0,
      label: '冷呼：初步介绍产品',
      lifecycle_stage_at_time: 'acquire',
    },
    {
      activity_id: 'act-xh-2',
      at: at(45, 15),
      event_type: 'email',
      direction: 'outbound',
      sentiment: 'neutral',
      sentiment_score: 0,
      label: '邮件：资料包发送',
      lifecycle_stage_at_time: 'acquire',
    },
    {
      activity_id: 'act-xh-3',
      at: at(35, 11),
      event_type: 'call',
      direction: 'inbound',
      sentiment: 'hesitant',
      sentiment_score: -1,
      label: '电话：表示暂无采购计划',
      lifecycle_stage_at_time: 'acquire',
    },
    {
      activity_id: 'act-xh-4',
      at: at(30, 9),
      event_type: 'note',
      direction: 'outbound',
      sentiment: 'negative',
      sentiment_score: -2,
      label: '备注：标记为低意向',
      lifecycle_stage_at_time: 'acquire',
    },
  ],
  milestones: [
    { type: 'status', at: at(25), label: '状态更新为不合格' },
  ],
  summary: {
    current_sentiment: 'negative',
    trend: 'down',
    days_since_positive: 55,
  },
}

const LEAD_005_JOURNEY: EmotionJourney = {
  subject_type: 'lead',
  subject_id: LEAD_005,
  lifecycle_current: 'retain',
  lifecycle_bands: [
    { stage: 'grow', from: at(90), to: at(30) },
    { stage: 'retain', from: at(30), to: at(0) },
  ],
  points: [
    {
      activity_id: 'act-md-2',
      at: at(15, 15),
      event_type: 'call',
      direction: 'inbound',
      sentiment: 'positive',
      sentiment_score: 2,
      label: '合同续签意向确认',
      lifecycle_stage_at_time: 'retain',
    },
    {
      activity_id: 'act-md-1',
      at: at(2, 11),
      event_type: 'meeting',
      direction: 'outbound',
      sentiment: 'positive',
      sentiment_score: 2,
      label: '季度业务回顾',
      lifecycle_stage_at_time: 'retain',
    },
  ],
  milestones: [{ type: 'converted', at: at(60), label: '转化为客户公司' }],
  summary: { current_sentiment: 'positive', trend: 'up', days_since_positive: 2 },
}

const FIXTURES: Record<string, EmotionJourney> = {
  [DEMO_LEAD_ID]: DEMO_LEAD_EMOTION_JOURNEY,
  [LEAD_002]: LEAD_002_JOURNEY,
  [LEAD_003]: LEAD_003_JOURNEY,
  [LEAD_004]: LEAD_004_JOURNEY,
  [LEAD_005]: LEAD_005_JOURNEY,
}

/** Demo Corp 租户下 API 空数据时回退 */
export function mockLeadEmotionJourney(leadId: string): EmotionJourney | null {
  return FIXTURES[leadId] ?? null
}
