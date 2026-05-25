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
import type { PieSeriesOption } from 'echarts/charts'
import type { ChartDonutItem } from '../../../chart/types'
import { useChartTheme } from '../../../chart/use-chart-theme'

const props = withDefaults(
  defineProps<{
    items: ChartDonutItem[]
    height?: number
    loadingText?: string
  }>(),
  { height: 280, loadingText: 'Loading…' },
)

const { colors, baseTooltip } = useChartTheme()

const option = computed(() => {
  const c = colors.value
  const palette = [c.primary, c.primaryEnd, c.axis, c.barTrack, c.glow]

  const series: PieSeriesOption = {
    type: 'pie',
    radius: ['52%', '78%'],
    center: ['50%', '50%'],
    avoidLabelOverlap: true,
    itemStyle: {
      borderRadius: 6,
      borderColor: c.tooltipBg,
      borderWidth: 2,
    },
    label: {
      show: true,
      formatter: '{b}\n{d}%',
      color: c.axis,
      fontSize: 11,
    },
    emphasis: {
      scale: true,
      scaleSize: 6,
      label: { fontWeight: 'bold' },
    },
    data: props.items.map((item, i) => ({
      name: item.name,
      value: item.value,
      itemStyle: { color: palette[i % palette.length] },
    })),
  }

  return {
    animationDuration: 600,
    tooltip: {
      ...baseTooltip(),
      trigger: 'item',
      formatter: (params: { name?: string; value?: number; percent?: number }) => {
        const pct = typeof params.percent === 'number' ? params.percent.toFixed(1) : '0'
        return `<div style="opacity:0.7;font-size:11px">${params.name ?? ''}</div><b>${params.value ?? 0} (${pct}%)</b>`
      },
    },
    series: [series],
  }
})
</script>
