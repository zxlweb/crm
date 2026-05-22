import { computed } from 'vue'
import type { ChartThemeColors } from './utils/colors'
import { CHART_COLOR_FALLBACK, resolveChartColors } from './utils/colors'
import {
  barTooltip,
  baseGrid,
  baseTooltip,
  categoryAxis,
  valueAxis,
  valueAxisWithGrid,
  type ValueAxisOptions,
} from './utils/echarts-parts'
import { useUiKitTheme } from '../theme/bridge'

export type { ChartThemeColors } from './utils/colors'

/**
 * 图表主题：读取 DS CSS 变量；随 UI_KIT_THEME_KEY 中的 themeId 更新。
 * 仅 ui-kit 图表组件使用。
 */
export function useChartTheme() {
  const { id: themeId } = useUiKitTheme()

  const colors = computed<ChartThemeColors>(() => {
    void themeId.value
    if (typeof document === 'undefined') return CHART_COLOR_FALLBACK
    return resolveChartColors()
  })

  return {
    colors,
    baseGrid: (padding?: { top: number; right: number; bottom: number; left: number }) =>
      baseGrid(colors.value, padding),
    baseTooltip: () => baseTooltip(colors.value),
    barTooltip: (horizontal: boolean) => barTooltip(colors.value, horizontal),
    categoryAxis: (data: string[]) => categoryAxis(colors.value, data),
    valueAxis: (formatter?: (value: number) => string, scale?: ValueAxisOptions) =>
      valueAxis(colors.value, formatter, scale),
    valueAxisWithGrid: (formatter?: (value: number) => string, scale?: ValueAxisOptions) =>
      valueAxisWithGrid(colors.value, formatter, scale),
  }
}
