/** Clamp percentage-style chart values to 0–100. */
export function clampPercent(value: number): number {
  if (!Number.isFinite(value)) return 0
  return Math.min(100, Math.max(0, value))
}

/** Normalize sparkline input; empty → [0]. */
export function normalizeSparklineValues(values: number[]): number[] {
  if (!values.length) return [0]
  return values.map((v) => (Number.isFinite(v) ? v : 0))
}
