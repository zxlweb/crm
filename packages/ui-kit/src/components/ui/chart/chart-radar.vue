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
import type { ChartRadarItem } from '../../../chart/types'
import { useChartTheme } from '../../../chart/use-chart-theme'

const props = withDefaults(
  defineProps<{
    indicators: Array<{ name: string; max?: number }>
    series: ChartRadarItem[]
    height?: number
    loadingText?: string
  }>(),
  { height: 320, loadingText: 'Loading…' },
)

const { colors, baseTooltip } = useChartTheme()

const option = computed(() => {
  const c = colors.value
  const palette = [c.primary, c.primaryEnd, c.axis, c.barTrack, c.glow]

  return {
    animationDuration: 600,
    tooltip: {
      ...baseTooltip(),
      trigger: 'item',
    },
    radar: {
      indicator: props.indicators.map((ind) => ({
        name: ind.name,
        max: ind.max ?? 100,
      })),
      shape: 'polygon' as const,
      axisName: {
        color: c.axis,
        fontSize: 11,
      },
      splitArea: {
        areaStyle: {
          color: ['transparent', `${c.barTrack}`],
        },
      },
      axisLine: {
        lineStyle: { color: c.barTrack },
      },
      splitLine: {
        lineStyle: { color: c.barTrack },
      },
    },
    series: [{
      type: 'radar' as const,
      data: props.series.map((s, i) => ({
        name: s.name,
        value: s.values,
        areaStyle: {
          color: `${palette[i % palette.length]}20`,
        },
        lineStyle: {
          color: palette[i % palette.length],
          width: 2,
        },
        itemStyle: {
          color: palette[i % palette.length],
        },
      })),
    }],
  }
})
</script>
