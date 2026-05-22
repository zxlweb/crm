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
import type { FunnelSeriesOption } from 'echarts/charts'
import type { ChartFunnelItem } from '~/types/chart'

const props = withDefaults(
  defineProps<{
    items: ChartFunnelItem[]
    height?: number
  }>(),
  { height: 320 },
)

const { colors, baseTooltip } = useChartTheme()
const { id: themeId } = useTheme()

const option = computed(() => {
  void themeId.value
  const c = colors.value
  const sorted = [...props.items].sort((a, b) => b.value - a.value)
  const palette = [c.primaryEnd, c.primaryEnd, c.primary, c.primary, c.primary]

  const series: FunnelSeriesOption = {
    type: 'funnel',
    left: '8%',
    right: '8%',
    top: 16,
    bottom: 16,
    sort: 'descending',
    gap: 4,
    color: palette,
    label: {
      show: true,
      position: 'inside',
      color: '#fff',
      fontSize: 12,
      formatter: '{b}\n{c}',
    },
    labelLine: { show: false },
    itemStyle: {
      borderColor: c.tooltipBg,
      borderWidth: 2,
    },
    emphasis: {
      label: { fontSize: 13, fontWeight: 'bold' },
    },
    data: sorted.map((item) => ({ name: item.name, value: item.value })),
  }

  return {
    animationDuration: 600,
    tooltip: {
      ...baseTooltip(),
      trigger: 'item',
    },
    series: [series],
  }
})
</script>
