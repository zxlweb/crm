import type { TooltipComponentOption } from 'echarts'
import type { ChartThemeColors } from './colors'
import { CHART_COLOR_FALLBACK } from './colors'

/** 柱状图 tooltip：阴影指示线 + 圆角卡片 */
export function barTooltip(colors: ChartThemeColors, horizontal: boolean): TooltipComponentOption {
  return {
    trigger: 'axis',
    backgroundColor: colors.tooltipBg,
    borderColor: colors.tooltipBorder,
    borderWidth: 1,
    padding: [10, 14],
    textStyle: { color: colors.tooltipFg, fontSize: 12 },
    extraCssText: 'border-radius: 12px; box-shadow: 0 8px 28px rgba(124,58,237,0.18);',
    axisPointer: horizontal
      ? {
          type: 'shadow',
          shadowStyle: { color: 'rgba(124, 58, 237, 0.06)' },
        }
      : {
          type: 'line',
          lineStyle: { color: colors.primaryEnd, width: 1, type: 'dashed', opacity: 0.5 },
        },
  }
}

export function baseGrid(
  _colors: ChartThemeColors = CHART_COLOR_FALLBACK,
  padding = { top: 12, right: 12, bottom: 28, left: 44 },
) {
  return {
    left: padding.left,
    right: padding.right,
    top: padding.top,
    bottom: padding.bottom,
    containLabel: false,
  }
}

export function baseTooltip(colors: ChartThemeColors): TooltipComponentOption {
  return {
    trigger: 'axis',
    backgroundColor: colors.tooltipBg,
    borderColor: colors.tooltipBorder,
    borderWidth: 1,
    padding: [10, 14],
    textStyle: { color: colors.tooltipFg, fontSize: 12 },
    extraCssText: 'border-radius: 12px; box-shadow: 0 8px 24px rgba(0,0,0,0.35);',
    axisPointer: {
      type: 'line',
      lineStyle: { color: colors.primary, width: 1, type: 'dashed', opacity: 0.35 },
    },
  }
}

export function categoryAxis(colors: ChartThemeColors, data: string[]) {
  return {
    type: 'category' as const,
    data,
    boundaryGap: false,
    axisLine: { show: false },
    axisTick: { show: false },
    axisLabel: { color: colors.axis, fontSize: 11, margin: 10 },
    splitLine: { show: false },
  }
}

export type ValueAxisOptions = {
  min?: number
  max?: number
  interval?: number
}

export function valueAxis(
  colors: ChartThemeColors,
  formatter?: (value: number) => string,
  scale?: ValueAxisOptions,
) {
  return {
    type: 'value' as const,
    min: scale?.min,
    max: scale?.max,
    interval: scale?.interval,
    axisLine: { show: false },
    axisTick: { show: false },
    axisLabel: {
      color: colors.axis,
      fontSize: 11,
      formatter: formatter ? (v: number) => formatter(v) : undefined,
    },
    splitLine: { show: false },
  }
}

export function valueAxisWithGrid(
  colors: ChartThemeColors,
  formatter?: (value: number) => string,
  scale?: ValueAxisOptions,
) {
  return {
    ...valueAxis(colors, formatter, scale),
    splitLine: {
      show: true,
      lineStyle: { color: colors.grid, type: 'dashed' as const, opacity: 0.8 },
    },
  }
}
