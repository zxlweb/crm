<template>
  <ClientOnly>
    <VChart
      class="w-full"
      :style="chartStyle"
      :option="option"
      autoresize
    />
    <template #fallback>
      <MotionlessSparkline
        :values="values"
        :width="width"
        :height="height"
        :stroke="strokeColor"
      />
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import { computed, defineComponent, h } from 'vue'
import type { LineSeriesOption } from 'echarts/charts'
import { normalizeSparklineValues } from '../../../chart/utils/chart-value-utils'
import { useChartTheme } from '../../../chart/use-chart-theme'

const props = withDefaults(
  defineProps<{
    values: number[]
    width?: number
    height?: number
    /** Semantic stroke: auto (trend), up, down, flat, brand */
    tone?: 'auto' | 'up' | 'down' | 'flat' | 'brand'
    showArea?: boolean
    loadingText?: string
  }>(),
  {
    width: 96,
    height: 32,
    tone: 'auto',
    showArea: true,
    loadingText: 'Loading…',
  },
)

const { colors } = useChartTheme()

const normalized = computed(() => normalizeSparklineValues(props.values))

const strokeColor = computed(() => {
  const c = colors.value
  if (props.tone === 'brand') return c.primary
  if (props.tone === 'up') return '#059669'
  if (props.tone === 'down') return '#dc2626'
  if (props.tone === 'flat') return c.axis
  const vals = normalized.value
  const first = vals[0] ?? 0
  const last = vals[vals.length - 1] ?? 0
  if (last > first) return '#059669'
  if (last < first) return '#dc2626'
  return c.primary
})

const chartStyle = computed(() => ({
  width: `${props.width}px`,
  height: `${props.height}px`,
}))

const option = computed(() => {
  const data = normalized.value

  const series: LineSeriesOption = {
    type: 'line',
    data,
    smooth: true,
    symbol: 'none',
    lineStyle: {
      width: 2,
      color: strokeColor.value,
    },
    areaStyle: props.showArea
      ? {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: `${strokeColor.value}33` },
              { offset: 1, color: `${strokeColor.value}00` },
            ],
          },
        }
      : undefined,
  }

  return {
    animationDuration: 400,
    grid: { left: 0, right: 0, top: 2, bottom: 2 },
    xAxis: { type: 'category', show: false, boundaryGap: false, data: data.map((_, i) => i) },
    yAxis: { type: 'value', show: false, scale: true },
    series: [series],
    tooltip: { show: false },
  }
})

const MotionlessSparkline = defineComponent({
  name: 'MotionlessSparkline',
  props: {
    values: { type: Array as () => number[], required: true },
    width: { type: Number, default: 96 },
    height: { type: Number, default: 32 },
    stroke: { type: String, default: 'currentColor' },
  },
  setup(p) {
    const points = computed(() => {
      const vals = normalizeSparklineValues(p.values)
      const max = Math.max(...vals, 1)
      const min = Math.min(...vals, 0)
      const range = max - min || 1
      const padY = 2
      const innerH = p.height - padY * 2
      return vals
        .map((v, i) => {
          const x = vals.length === 1 ? p.width / 2 : (i / (vals.length - 1)) * p.width
          const y = padY + innerH - ((v - min) / range) * innerH
          return `${x},${y}`
        })
        .join(' ')
    })
    return () =>
      h(
        'svg',
        {
          width: p.width,
          height: p.height,
          viewBox: `0 0 ${p.width} ${p.height}`,
          class: 'shrink-0 overflow-visible',
          'aria-hidden': 'true',
        },
        [
          h('polyline', {
            points: points.value,
            fill: 'none',
            stroke: p.stroke,
            'stroke-width': 2,
            'stroke-linecap': 'round',
            'stroke-linejoin': 'round',
          }),
        ],
      )
  },
})
</script>
