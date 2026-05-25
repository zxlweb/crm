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

export type ActivitySentiment = 'positive' | 'neutral' | 'hesitant' | 'negative' | 'unknown'

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
