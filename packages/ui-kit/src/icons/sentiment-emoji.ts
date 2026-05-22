/** 情绪 → Emoji（系统原生表情，跨平台一致展示） */

export type SentimentKey = 'positive' | 'neutral' | 'hesitant' | 'negative' | 'unknown'

export const SENTIMENT_EMOJI: Record<SentimentKey, string> = {
  positive: '😊',
  neutral: '😐',
  hesitant: '😕',
  negative: '😞',
  unknown: '❔',
}

const SENTIMENT_ALIAS: Record<string, SentimentKey> = {
  positive: 'positive',
  neutral: 'neutral',
  hesitant: 'hesitant',
  negative: 'negative',
  unknown: 'unknown',
}

/** 情绪分 → 情绪键（与情绪曲线 Y 轴一致） */
export const SENTIMENT_SCORE_KEY: Record<number, SentimentKey> = {
  2: 'positive',
  0: 'neutral',
  [-1]: 'hesitant',
  [-2]: 'negative',
}

export function resolveSentimentKey(sentiment: string): SentimentKey {
  return SENTIMENT_ALIAS[sentiment] ?? 'unknown'
}

export function sentimentEmoji(sentiment: string): string {
  return SENTIMENT_EMOJI[resolveSentimentKey(sentiment)]
}

export function sentimentEmojiByScore(score: number): string {
  return SENTIMENT_EMOJI[SENTIMENT_SCORE_KEY[score] ?? 'unknown']
}
