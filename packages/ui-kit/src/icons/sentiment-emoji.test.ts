import { describe, expect, it } from 'vitest'
import { sentimentEmoji, sentimentEmojiByScore } from './sentiment-emoji'

describe('sentiment-emoji', () => {
  it('maps known sentiments', () => {
    expect(sentimentEmoji('positive')).toBe('😊')
    expect(sentimentEmoji('negative')).toBe('😞')
  })

  it('falls back for unknown', () => {
    expect(sentimentEmoji('other')).toBe('❔')
  })

  it('maps scores', () => {
    expect(sentimentEmojiByScore(2)).toBe('😊')
    expect(sentimentEmojiByScore(-2)).toBe('😞')
  })
})
