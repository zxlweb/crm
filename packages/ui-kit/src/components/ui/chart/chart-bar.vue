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
import type { BarSeriesOption } from 'echarts/charts'
import type { ChartBarItem } from '../../../chart/types'
import { useChartTheme } from '../../../chart/use-chart-theme'

const props = withDefaults(
  defineProps<{
    items: ChartBarItem[]
    height?: number
    horizontal?: boolean
    valueFormatter?: (value: number) => string
    animated?: boolean
    loadingText?: string
  }>(),
  { height: 280, horizontal: true, animated: true, loadingText: 'Loading…' },
)

const { colors, baseGrid, barTooltip, categoryAxis, valueAxisWithGrid } = useChartTheme()

function barGradient(horizontal: boolean) {
  const c = colors.value
  return {
    type: 'linear' as const,
    x: 0,
    y: 0,
    x2: horizontal ? 1 : 0,
    y2: horizontal ? 0 : 1,
    global: false,
    colorStops: [
      { offset: 0, color: c.primary },
      { offset: 1, color: c.primaryEnd },
    ],
  }
}

function buildSeries(): BarSeriesOption {
  const c = colors.value
  const horizontal = props.horizontal
  const values = props.items.map((i) => i.value)
  const maxVal = Math.max(...values, 1)
  const radius = horizontal ? [0, 8, 8, 0] : [8, 8, 0, 0]

  return {
    name: horizontal ? 'share' : 'revenue',
    type: 'bar',
    barMaxWidth: horizontal ? 16 : 40,
    showBackground: true,
    backgroundStyle: {
      color: c.barTrack,
      borderRadius: radius,
    },
    animation: props.animated,
    animationDuration: 900,
    animationEasing: 'cubicOut',
    animationDelay: (idx: number) => (props.animated ? idx * 90 : 0),
    animationDurationUpdate: 300,
    emphasis: {
      focus: 'series',
      itemStyle: {
        shadowBlur: 16,
        shadowColor: c.glow,
        borderRadius: radius,
        color: barGradient(horizontal),
      },
      label: horizontal
        ? {
            show: true,
            fontWeight: 'bold',
            color: c.primary,
          }
        : undefined,
    },
    itemStyle: {
      borderRadius: radius,
      color: barGradient(horizontal),
    },
    data: values,
    label: horizontal
      ? {
          show: true,
          position: 'right',
          distance: 8,
          formatter: (params) => {
            const val = typeof params.value === 'number' ? params.value : 0
            const pct = Math.round((val / maxVal) * 100)
            return props.valueFormatter ? props.valueFormatter(val) : `${pct}%`
          },
          color: c.axis,
          fontSize: 11,
        }
      : undefined,
  }
}

const option = computed(() => {
  const names = props.items.map((i) => i.name)
  const values = props.items.map((i) => i.value)
  const maxVal = Math.max(...values, 1)
  const series = buildSeries()

  if (props.horizontal) {
    return {
      animation: props.animated,
      animationDuration: 900,
      animationEasing: 'cubicOut' as const,
      grid: baseGrid({ top: 8, right: 52, bottom: 8, left: 96 }),
      tooltip: barTooltip(true),
      xAxis: {
        type: 'value' as const,
        min: 0,
        max: 100,
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: { color: colors.value.axis, fontSize: 11 },
        splitLine: { show: false },
      },
      yAxis: { ...categoryAxis(names), inverse: true },
      series: [series],
    }
  }

  const yMax = Math.ceil(maxVal * 1.25)

  return {
    animation: props.animated,
    animationDuration: 900,
    animationEasing: 'cubicOut' as const,
    grid: baseGrid({ top: 12, right: 16, bottom: 28, left: 48 }),
    tooltip: {
      ...barTooltip(false),
      formatter: (params: unknown) => {
        const list = Array.isArray(params) ? params : [params]
        const item = list[0] as { axisValue?: string, value?: number }
        const val = typeof item.value === 'number' ? item.value : 0
        const formatted = props.valueFormatter ? props.valueFormatter(val) : String(val)
        return `<div style="opacity:0.7;font-size:11px">${item.axisValue ?? ''}</div><b>${formatted}</b>`
      },
    },
    xAxis: { ...categoryAxis(names), boundaryGap: true },
    yAxis: { ...valueAxisWithGrid(props.valueFormatter), max: yMax },
    series: [series],
  }
})
</script>
