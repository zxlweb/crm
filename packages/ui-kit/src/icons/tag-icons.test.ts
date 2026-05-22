import { describe, expect, it } from 'vitest'
import {
  ACTIVITY_TYPE_ICON,
  TAG_ICONS,
  resolveActivityTypeIcon,
  resolveSentimentIcon,
} from './tag-icons'

describe('tag-icons', () => {
  it('defines all activity type icons', () => {
    for (const key of Object.keys(ACTIVITY_TYPE_ICON)) {
      expect(TAG_ICONS[ACTIVITY_TYPE_ICON[key]!].length).toBeGreaterThan(0)
    }
  })

  it('resolves unknown activity to note', () => {
    expect(resolveActivityTypeIcon('unknown')).toBe('note')
  })

  it('resolves sentiment icons', () => {
    expect(resolveSentimentIcon('positive')).toBe('positive')
    expect(resolveSentimentIcon('unknown')).toBe('neutral')
  })
})
