import { describe, expect, it } from 'vitest'
import { CHART_COLOR_FALLBACK } from './colors'
import { categoryAxis, valueAxisWithGrid } from './echarts-parts'

describe('chart echarts-parts utils', () => {
  it('categoryAxis returns category type with axis labels', () => {
    const axis = categoryAxis(CHART_COLOR_FALLBACK, ['Mon', 'Tue'])
    expect(axis.type).toBe('category')
    expect(axis.data).toEqual(['Mon', 'Tue'])
    expect(axis.axisLabel.color).toBe(CHART_COLOR_FALLBACK.axis)
  })

  it('valueAxisWithGrid enables dashed splitLine', () => {
    const axis = valueAxisWithGrid(CHART_COLOR_FALLBACK)
    expect(axis.splitLine?.show).toBe(true)
    expect(axis.splitLine?.lineStyle?.color).toBe(CHART_COLOR_FALLBACK.grid)
  })
})
