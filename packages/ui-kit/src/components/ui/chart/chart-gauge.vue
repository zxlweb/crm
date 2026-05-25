<template>
  <ClientOnly>
    <VChart
      class="w-full"
      :style="{ height: `${height}px` }"
      :option="option"
      autoresize
    />
    <template #fallback>
      <MotionlessGauge :value="value" :label="label" :height="height" />
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import { computed, defineComponent, h } from 'vue'
import type { GaugeSeriesOption } from 'echarts/charts'
import { clampPercent } from '../../../chart/utils/chart-value-utils'
import { useChartTheme } from '../../../chart/use-chart-theme'

const props = withDefaults(
  defineProps<{
    /** Completion percentage 0–100 */
    value: number
    label?: string
    height?: number
    loadingText?: string
  }>(),
  {
    label: '',
    height: 220,
    loadingText: 'Loading…',
  },
)

const { colors, baseTooltip } = useChartTheme()

const clamped = computed(() => clampPercent(props.value))

const option = computed(() => {
  const c = colors.value
  const pct = clamped.value

  const series: GaugeSeriesOption = {
    type: 'gauge',
    startAngle: 210,
    endAngle: -30,
    min: 0,
    max: 100,
    splitNumber: 5,
    radius: '88%',
    center: ['50%', '58%'],
    pointer: {
      show: true,
      length: '58%',
      width: 6,
      itemStyle: { color: c.primaryEnd },
    },
    progress: {
      show: true,
      width: 14,
      roundCap: true,
      itemStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 1,
          y2: 0,
          colorStops: [
            { offset: 0, color: c.primary },
            { offset: 1, color: c.primaryEnd },
          ],
        },
      },
    },
    axisLine: {
      lineStyle: {
        width: 14,
        color: [[1, c.barTrack]],
      },
    },
    axisTick: { show: false },
    splitLine: { show: false },
    axisLabel: {
      distance: 18,
      color: c.axis,
      fontSize: 10,
    },
    anchor: {
      show: true,
      size: 8,
      itemStyle: { color: c.primaryEnd, borderWidth: 2, borderColor: c.dotStroke },
    },
    title: props.label
      ? {
          show: true,
          offsetCenter: [0, '78%'],
          fontSize: 12,
          color: c.axis,
        }
      : { show: false },
    detail: {
      valueAnimation: true,
      fontSize: 28,
      fontWeight: 'bold',
      color: c.primaryEnd,
      offsetCenter: [0, '18%'],
      formatter: '{value}%',
    },
    data: [{ value: pct, name: props.label }],
  }

  return {
    animationDuration: 600,
    tooltip: {
      ...baseTooltip(),
      formatter: () => `${pct}%`,
    },
    series: [series],
  }
})

const MotionlessGauge = defineComponent({
  name: 'MotionlessGauge',
  props: {
    value: { type: Number, required: true },
    label: { type: String, default: '' },
    height: { type: Number, default: 220 },
  },
  setup(p) {
    const pct = computed(() => clampPercent(p.value))
    return () =>
      h(
        'div',
        {
          class: 'flex flex-col items-center justify-center text-ds-fg-muted',
          style: { height: `${p.height}px` },
          role: 'img',
          'aria-label': p.label ? `${p.label}: ${pct.value}%` : `${pct.value}%`,
        },
        [
          h('span', { class: 'text-3xl font-bold tabular-nums text-ds-fg-heading' }, `${pct.value}%`),
          p.label ? h('span', { class: 'mt-1 text-xs' }, p.label) : null,
        ],
      )
  },
})
</script>
