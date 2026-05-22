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
        {{ $t('loading') }}
      </div>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import type { ComposeOption } from 'echarts/core'
import type { LineSeriesOption } from 'echarts/charts'
import type { GridComponentOption, TooltipComponentOption, XAXisComponentOption, YAXisComponentOption } from 'echarts'
import type { ChartSeries } from '~/types/chart'

const props = withDefaults(
  defineProps<{
    categories: string[]
    series: ChartSeries[]
    height?: number
    yFormatter?: (value: number) => string
    showArea?: boolean
  }>(),
  { height: 260, showArea: true },
)

const { colors, baseGrid, baseTooltip, categoryAxis, valueAxis } = useChartTheme()
const { id: themeId } = useTheme()

const option = computed(() => {
  void themeId.value
  const c = colors.value

  const echartsSeries: LineSeriesOption[] = props.series.map((s) => {
    const isCompare = s.compare === true
    const isPrimary = s.primary !== false && !isCompare

    return {
      name: s.name,
      type: 'line',
      smooth: true,
      showSymbol: false,
      symbol: 'circle',
      symbolSize: 10,
      emphasis: {
        focus: 'series',
        itemStyle: {
          color: c.primaryEnd,
          borderColor: c.dotStroke,
          borderWidth: 3,
          shadowBlur: 12,
          shadowColor: c.glow,
        },
      },
      lineStyle: {
        width: isPrimary ? 2.5 : 1.5,
        color: isCompare ? c.secondary : c.primaryEnd,
        type: isCompare ? 'dashed' : 'solid',
        shadowBlur: isPrimary ? 14 : 0,
        shadowColor: isPrimary ? c.glow : 'transparent',
      },
      itemStyle: {
        color: isPrimary ? c.primaryEnd : c.secondary,
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
      data: s.data,
    }
  })

  return {
    animationDuration: 600,
    animationEasing: 'cubicOut',
    grid: baseGrid(),
    tooltip: {
      ...baseTooltip(),
      formatter: (params: unknown) => {
        const list = Array.isArray(params) ? params : [params]
        const first = list[0] as { axisValue?: string; axisValueLabel?: string }
        const title = first?.axisValue ?? first?.axisValueLabel ?? ''
        const lines = list.map((p) => {
          const item = p as { seriesName?: string; value?: number; color?: string }
          const val = typeof item.value === 'number' ? item.value : 0
          const formatted = props.yFormatter ? props.yFormatter(val) : String(val)
          return `<span style="color:${c.tooltipFg};opacity:0.7">${item.seriesName}</span> <b>${formatted}</b>`
        })
        return `<div style="color:${c.tooltipFg};font-size:11px;opacity:0.65;margin-bottom:4px">${title}</div>${lines.join('<br/>')}`
      },
    } as TooltipComponentOption,
    xAxis: categoryAxis(props.categories) as XAXisComponentOption,
    yAxis: valueAxis(props.yFormatter) as YAXisComponentOption,
    series: echartsSeries,
  } satisfies ComposeOption
})
</script>
