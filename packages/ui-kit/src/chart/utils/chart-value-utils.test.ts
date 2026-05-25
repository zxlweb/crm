import { describe, expect, it } from 'vitest'
import { clampPercent, normalizeSparklineValues } from './chart-value-utils'

describe('clampPercent', () => {
  it('clamps to 0–100', () => {
    expect(clampPercent(-5)).toBe(0)
    expect(clampPercent(72.4)).toBe(72.4)
    expect(clampPercent(120)).toBe(100)
    expect(clampPercent(Number.NaN)).toBe(0)
  })
})

describe('normalizeSparklineValues', () => {
  it('returns [0] for empty input', () => {
    expect(normalizeSparklineValues([])).toEqual([0])
  })

  it('replaces non-finite values with 0', () => {
    expect(normalizeSparklineValues([1, Number.NaN, 3])).toEqual([1, 0, 3])
  })
})
