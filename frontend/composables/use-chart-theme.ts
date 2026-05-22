import type { TooltipComponentOption } from 'echarts'

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

const FALLBACK: ChartThemeColors = {
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

function readCssVar(name: string, fallback: string) {
  if (!import.meta.client) return fallback
  const v = getComputedStyle(document.documentElement).getPropertyValue(name).trim()
  return v || fallback
}

export function useChartTheme() {
  const { id: themeId } = useTheme()

  const colors = computed<ChartThemeColors>(() => {
    void themeId.value
    if (!import.meta.client) return FALLBACK
    return {
      primary: readCssVar('--ds-chart-line-end', FALLBACK.primaryEnd),
      primaryEnd: readCssVar('--ds-chart-line-start', FALLBACK.primary),
      secondary: readCssVar('--ds-chart-secondary', FALLBACK.secondary),
      grid: readCssVar('--ds-chart-grid', FALLBACK.grid),
      axis: readCssVar('--ds-chart-axis', FALLBACK.axis),
      areaStart: readCssVar('--ds-chart-area-start', FALLBACK.areaStart),
      areaMid: readCssVar('--ds-chart-area-mid', FALLBACK.areaMid),
      areaEnd: readCssVar('--ds-chart-area-end', FALLBACK.areaEnd),
      tooltipBg: readCssVar('--ds-chart-tooltip-bg', FALLBACK.tooltipBg),
      tooltipBorder: readCssVar('--ds-chart-tooltip-border', FALLBACK.tooltipBorder),
      tooltipFg: readCssVar('--ds-chart-tooltip-fg', FALLBACK.tooltipFg),
      glow: readCssVar('--ds-chart-glow', FALLBACK.glow),
      dotStroke: readCssVar('--ds-chart-dot-stroke', FALLBACK.dotStroke),
      barTrack: readCssVar('--ds-chart-bar-track', FALLBACK.barTrack),
    }
  })

  /** 柱状图 tooltip：阴影指示线 + 圆角卡片 */
  function barTooltip(horizontal: boolean): TooltipComponentOption {
    const c = colors.value
    return {
      trigger: 'axis',
      backgroundColor: c.tooltipBg,
      borderColor: c.tooltipBorder,
      borderWidth: 1,
      padding: [10, 14],
      textStyle: { color: c.tooltipFg, fontSize: 12 },
      extraCssText: 'border-radius: 12px; box-shadow: 0 8px 28px rgba(124,58,237,0.18);',
      axisPointer: horizontal
        ? {
            type: 'shadow',
            shadowStyle: { color: 'rgba(124, 58, 237, 0.06)' },
          }
        : {
            type: 'line',
            lineStyle: { color: c.primaryEnd, width: 1, type: 'dashed', opacity: 0.5 },
          },
    }
  }

  function baseGrid(padding = { top: 12, right: 12, bottom: 28, left: 44 }) {
    return {
      left: padding.left,
      right: padding.right,
      top: padding.top,
      bottom: padding.bottom,
      containLabel: false,
    }
  }

  function baseTooltip(): TooltipComponentOption {
    const c = colors.value
    return {
      trigger: 'axis',
      backgroundColor: c.tooltipBg,
      borderColor: c.tooltipBorder,
      borderWidth: 1,
      padding: [10, 14],
      textStyle: { color: c.tooltipFg, fontSize: 12 },
      extraCssText: 'border-radius: 12px; box-shadow: 0 8px 24px rgba(0,0,0,0.35);',
      axisPointer: {
        type: 'line',
        lineStyle: { color: c.primary, width: 1, type: 'dashed', opacity: 0.35 },
      },
    }
  }

  function categoryAxis(data: string[]) {
    const c = colors.value
    return {
      type: 'category' as const,
      data,
      boundaryGap: false,
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: c.axis, fontSize: 11, margin: 10 },
      splitLine: { show: false },
    }
  }

  function valueAxis(formatter?: (value: number) => string) {
    const c = colors.value
    return {
      type: 'value' as const,
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: {
        color: c.axis,
        fontSize: 11,
        formatter: formatter ? (v: number) => formatter(v) : undefined,
      },
      splitLine: { show: false },
    }
  }

  function valueAxisWithGrid(formatter?: (value: number) => string) {
    const c = colors.value
    return {
      ...valueAxis(formatter),
      splitLine: {
        show: true,
        lineStyle: { color: c.grid, type: 'dashed' as const, opacity: 0.8 },
      },
    }
  }

  return { colors, baseGrid, baseTooltip, barTooltip, categoryAxis, valueAxis, valueAxisWithGrid }
}
