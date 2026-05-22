/** docs/api/phase-2-crm-ai.md §7 */

export type EmotionSubjectType = 'lead' | 'contact' | 'account'

export type ActivitySentiment = 'positive' | 'neutral' | 'hesitant' | 'negative' | 'unknown'

export type EmotionTrend = 'up' | 'down' | 'flat'

export type EmotionJourneyPoint = {
  activity_id: string
  at: string
  event_type: string
  direction?: string
  sentiment: ActivitySentiment | null
  sentiment_score: number | null
  sentiment_source?: string
  label?: string
  lifecycle_stage_at_time?: string
}

export type EmotionJourneyBand = {
  stage: string
  from: string
  to: string
}

export type EmotionJourneyMilestone = {
  type: string
  at: string
  label: string
}

export type EmotionJourneySummary = {
  current_sentiment: ActivitySentiment | null
  trend: EmotionTrend
  days_since_positive?: number
}

export type EmotionJourney = {
  subject_type: EmotionSubjectType
  subject_id: string
  lifecycle_current: string
  lifecycle_bands: EmotionJourneyBand[]
  points: EmotionJourneyPoint[]
  milestones: EmotionJourneyMilestone[]
  summary: EmotionJourneySummary
}
