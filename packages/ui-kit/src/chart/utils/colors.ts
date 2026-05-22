/** 图表主题色（来自 --ds-chart-* CSS 变量） */
export interface ChartThemeColors {
  primary: string
  primaryEnd: string
  secondary: string
  grid: string
  axis: string
  areaStart: string
  areaMid: string
  areaEnd: string
  tooltipBg: string
  tooltipBorder: string
  tooltipFg: string
  glow: string
  dotStroke: string
  barTrack: string
}

export const CHART_COLOR_FALLBACK: ChartThemeColors = {
  primary: '#7c3aed',
  primaryEnd: '#c4b5fd',
  secondary: '#94a3b8',
  grid: 'transparent',
  axis: '#71717a',
  areaStart: 'rgba(167, 139, 250, 0.28)',
  areaMid: 'rgba(196, 181, 253, 0.14)',
  areaEnd: 'transparent',
  tooltipBg: '#1a1a1a',
  tooltipBorder: 'rgba(255,255,255,0.08)',
  tooltipFg: '#fafafa',
  glow: 'rgba(167,139,250,0.55)',
  dotStroke: '#141414',
  barTrack: 'rgba(139, 92, 246, 0.1)',
}

const CHART_CSS_VARS: Record<keyof ChartThemeColors, string> = {
  primary: '--ds-chart-line-end',
  primaryEnd: '--ds-chart-line-start',
  secondary: '--ds-chart-secondary',
  grid: '--ds-chart-grid',
  axis: '--ds-chart-axis',
  areaStart: '--ds-chart-area-start',
  areaMid: '--ds-chart-area-mid',
  areaEnd: '--ds-chart-area-end',
  tooltipBg: '--ds-chart-tooltip-bg',
  tooltipBorder: '--ds-chart-tooltip-border',
  tooltipFg: '--ds-chart-tooltip-fg',
  glow: '--ds-chart-glow',
  dotStroke: '--ds-chart-dot-stroke',
  barTrack: '--ds-chart-bar-track',
}

/** 读取单条 CSS 变量（SSR / 单测可传入 root） */
export function readCssVar(name: string, fallback: string, root?: HTMLElement | null): string {
  if (!root && typeof document === 'undefined') return fallback
  const el = root ?? (typeof document !== 'undefined' ? document.documentElement : null)
  if (!el) return fallback
  const v = getComputedStyle(el).getPropertyValue(name).trim()
  return v || fallback
}

/** 从 document 解析当前图表色板 */
export function resolveChartColors(
  fallback: ChartThemeColors = CHART_COLOR_FALLBACK,
  root?: HTMLElement | null,
): ChartThemeColors {
  const resolved = { ...fallback }
  for (const key of Object.keys(CHART_CSS_VARS) as (keyof ChartThemeColors)[]) {
    resolved[key] = readCssVar(CHART_CSS_VARS[key], fallback[key], root)
  }
  return resolved
}
