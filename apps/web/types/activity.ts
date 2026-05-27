export type ActivitySubjectType = 'lead' | 'account' | 'contact'

export type ActivityEventType =
  | 'note'
  | 'call'
  | 'email'
  | 'meeting'
  | 'wechat'
  | 'visit'
  | 'system'

export type ActivityDirection = 'inbound' | 'outbound'

/** 人工情绪标注标准 5 档（与 API §2.4.1、后端 ValidActivitySentiments 一致） */
export const ACTIVITY_SENTIMENTS = [
  'positive',
  'neutral',
  'hesitant',
  'negative',
  'unknown',
] as const

export type ActivitySentiment = (typeof ACTIVITY_SENTIMENTS)[number]

export const DEFAULT_ACTIVITY_SENTIMENT: ActivitySentiment = 'unknown'

export type Activity = {
  id: string
  subject_type: ActivitySubjectType
  subject_id: string
  event_type: ActivityEventType
  direction?: string
  body?: string
  metadata?: Record<string, unknown>
  label?: string
  sentiment?: ActivitySentiment | null
  sentiment_source?: string
  occurred_at: string
  created_by?: string
  created_at: string
  updated_at: string
}

export type ActivityCreateInput = {
  subject_type: ActivitySubjectType
  subject_id: string
  event_type: ActivityEventType
  direction?: ActivityDirection
  body?: string
  sentiment?: ActivitySentiment
  sentiment_source?: 'manual' | 'rule'
  occurred_at?: string
  label?: string
  metadata?: Record<string, unknown>
}

export type ActivityUpdateInput = {
  body?: string
  direction?: ActivityDirection
  sentiment?: ActivitySentiment
  sentiment_source?: 'manual' | 'rule'
  occurred_at?: string
  label?: string
  metadata?: Record<string, unknown>
  clear_sentiment?: boolean
}

export type ActivitySummary = {
  items: Array<{ label: string; value: number; percentage?: number }>
  total: number
}

export type Pagination = {
  page: number
  page_size: number
  total: number
}
