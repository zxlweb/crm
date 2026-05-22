<template>
  <ClientOnly>
    <VChart
      class="w-full"
      :style="{ height: `${height}px` }"
      :option="option"
      autoresize
    />
    <template #fallback>
      <div class="flex items-center justify-center text-sm text-ds-fg-muted" :style="{ height: `${height}px` }">
        {{ loadingText }}
      </div>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ComposeOption } from 'echarts/core'
import type { LineSeriesOption } from 'echarts/charts'
import type { GridComponentOption, TooltipComponentOption, XAXisComponentOption, YAXisComponentOption } from 'echarts'
import type { ChartSeries } from '../../../chart/types'
import { useChartTheme } from '../../../chart/use-chart-theme'

const props = withDefaults(
  defineProps<{
    categories: string[]
    series: ChartSeries[]
    height?: number
    yFormatter?: (value: number) => string
    /** 与 categories 等长；每项为 ECharts symbol（如 path://） */
    pointSymbols?: string[]
    pointSymbolSize?: number
    /** 与 categories 等长；在数据点上方显示 emoji */
    pointEmojis?: string[]
    /** 与 categories 等长；用于 tooltip 展示（可含 emoji） */
    pointLabels?: string[]
    /** 固定 Y 轴刻度（如情绪分 -2～2） */
    yMin?: number
    yMax?: number
    yInterval?: number
    showArea?: boolean
    loadingText?: string
  }>(),
  { height: 260, showArea: true, pointSymbolSize: 22, loadingText: 'Loading…' },
)

const { colors, baseGrid, baseTooltip, categoryAxis, valueAxis } = useChartTheme()

const option = computed(() => {
  const c = colors.value

  const hasPointSymbols = (props.pointSymbols?.length ?? 0) > 0
  const hasPointEmojis = (props.pointEmojis?.length ?? 0) > 0
  const hasPointMarkers = hasPointSymbols || hasPointEmojis

  const echartsSeries: LineSeriesOption[] = props.series.map((s) => {
    const isCompare = s.compare === true
    const isPrimary = s.primary !== false && !isCompare
    const usePathSymbols = isPrimary && hasPointSymbols
    const useEmojis = isPrimary && hasPointEmojis

    const seriesData = usePathSymbols
      ? s.data.map((v, i) => ({
          value: v,
          symbol: props.pointSymbols![i] ?? 'circle',
          symbolSize: props.pointSymbolSize,
        }))
      : s.data

    return {
      name: s.name,
      type: 'line',
      smooth: true,
      /** 折线默认仅 hover 显示拐点；有 emoji/path 标记时需常显 */
      showSymbol: usePathSymbols || useEmojis,
      showAllSymbol: usePathSymbols || useEmojis ? true : 'auto',
      symbol: 'circle',
      symbolSize: usePathSymbols
        ? props.pointSymbolSize
        : useEmojis
          ? 6
          : 10,
      label: useEmojis
        ? {
            show: true,
            position: 'top',
            distance: 10,
            fontSize: 22,
            formatter: (p: { dataIndex?: number }) =>
              props.pointEmojis![p.dataIndex ?? 0] ?? '',
          }
        : undefined,
      labelLayout: useEmojis ? { hideOverlap: false } : undefined,
      emphasis: {
        focus: 'series',
        label: useEmojis ? { show: true } : undefined,
        itemStyle: {
          color: c.primaryEnd,
          borderColor: c.dotStroke,
          borderWidth: 3,
          shadowBlur: 12,
          shadowColor: c.glow,
        },
      },
      blur: useEmojis ? { label: { show: true } } : undefined,
      lineStyle: {
        width: isPrimary ? 2.5 : 1.5,
        color: isCompare ? c.secondary : c.primaryEnd,
        type: isCompare ? 'dashed' : 'solid',
        shadowBlur: isPrimary ? 14 : 0,
        shadowColor: isPrimary ? c.glow : 'transparent',
      },
      itemStyle: {
        color: isPrimary ? c.primaryEnd : c.secondary,
        borderColor: c.dotStroke,
        borderWidth: usePathSymbols || useEmojis ? 1.5 : 0,
        opacity: useEmojis && !usePathSymbols ? 0.35 : 1,
      },
      areaStyle: isPrimary && props.showArea
        ? {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: c.areaStart },
                { offset: 0.55, color: c.areaMid },
                { offset: 1, color: c.areaEnd || 'transparent' },
              ],
            },
            opacity: 1,
          }
        : undefined,
      data: seriesData,
    }
  })

  const gridPadding = hasPointMarkers
    ? { top: hasPointEmojis ? 44 : 28, right: 12, bottom: 28, left: 48 }
    : undefined

  return {
    animationDuration: 600,
    animationEasing: 'cubicOut',
    grid: baseGrid(colors.value, gridPadding),
    tooltip: {
      ...baseTooltip(),
      formatter: (params: unknown) => {
        const list = Array.isArray(params) ? params : [params]
        const first = list[0] as { axisValue?: string; axisValueLabel?: string }
        const title = first?.axisValue ?? first?.axisValueLabel ?? ''
        const lines = list.map((p) => {
          const item = p as { seriesName?: string; value?: number; dataIndex?: number }
          const val =
            typeof item.value === 'number'
              ? item.value
              : typeof (item as { value?: { value?: number } }).value === 'object'
                ? (item as { value: { value: number } }).value?.value ?? 0
                : 0
          const idx = item.dataIndex ?? 0
          const pointLabel = props.pointLabels?.[idx]
          const formatted = pointLabel ?? (props.yFormatter ? props.yFormatter(val) : String(val))
          return `<span style="color:${c.tooltipFg};opacity:0.7">${item.seriesName}</span> <b>${formatted}</b>`
        })
        return `<div style="color:${c.tooltipFg};font-size:11px;opacity:0.65;margin-bottom:4px">${title}</div>${lines.join('<br/>')}`
      },
    } as TooltipComponentOption,
    xAxis: categoryAxis(props.categories) as XAXisComponentOption,
    yAxis: valueAxis(props.yFormatter, {
      min: props.yMin,
      max: props.yMax,
      interval: props.yInterval,
    }) as YAXisComponentOption,
    series: echartsSeries,
  }
})
</script>
