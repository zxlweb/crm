import { describe, expect, it } from 'vitest'
import {
  SENTIMENT_ECHART_SYMBOL,
  TAG_ICONS,
  sentimentEchartSymbol,
} from './tag-icons'
import { svgPathDToEchart, tagIconPathsToEchartPath } from './tag-icon-echart'

describe('tag-icon-echart', () => {
  it('converts quadratic smile path for echarts', () => {
    expect(svgPathDToEchart('M7.5 13.25 Q12 17.25 16.5 13.25')).toBe(
      'M7.5,13.25Q12,17.25,16.5,13.25',
    )
  })

  it('builds path:// from tag icon shapes', () => {
    const path = tagIconPathsToEchartPath(TAG_ICONS.positive)
    expect(path.startsWith('path://')).toBe(true)
    expect(path).toContain('M12,4a8,8')
    expect(path).toContain('Q12,17.25')
  })

  it('keeps sentiment echart symbols in sync with TAG_ICONS', () => {
    for (const key of ['positive', 'neutral', 'hesitant', 'negative'] as const) {
      expect(SENTIMENT_ECHART_SYMBOL[key]).toBe(
        tagIconPathsToEchartPath(TAG_ICONS[key]),
      )
    }
    expect(sentimentEchartSymbol('positive')).toBe(SENTIMENT_ECHART_SYMBOL.positive)
  })
})
