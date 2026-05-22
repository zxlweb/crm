import { describe, expect, it } from 'vitest'
import { CHART_COLOR_FALLBACK, readCssVar, resolveChartColors } from './colors'

describe('chart colors utils', () => {
  it('readCssVar returns fallback without DOM', () => {
    expect(readCssVar('--ds-chart-axis', '#111', null)).toBe('#111')
  })

  it('resolveChartColors uses fallback keys when root is null', () => {
    const colors = resolveChartColors(CHART_COLOR_FALLBACK, null)
    expect(colors.primary).toBe(CHART_COLOR_FALLBACK.primary)
    expect(colors.axis).toBe(CHART_COLOR_FALLBACK.axis)
  })

  it.skipIf(typeof document === 'undefined')('resolveChartColors reads CSS variables from root', () => {
    const root = document.createElement('div')
    root.style.setProperty('--ds-chart-axis', '#abcdef')
    const colors = resolveChartColors(CHART_COLOR_FALLBACK, root)
    expect(colors.axis).toBe('#abcdef')
  })
})
