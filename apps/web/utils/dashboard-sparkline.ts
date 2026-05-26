/** 生成跟进卡内嵌趋势折线：单调缓升，避免低分时贴底形成「横线+折角」 */

export function buildEngagementSparkline(score: number): number[] {
  const v = Math.min(100, Math.max(0, Math.round(score)))
  if (v <= 0) return [2, 3, 4, 5]

  const a = Math.max(4, Math.round(v * 0.68))
  const b = Math.max(a + 1, Math.round(v * 0.78))
  const c = Math.max(b + 1, Math.round(v * 0.88))
  return [a, b, c, v]
}
