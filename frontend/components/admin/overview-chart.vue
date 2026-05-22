<template>
  <ChartShell :title="$t('adminChartTitle')" :subtitle="$t('adminChartSubtitle')" :height="300">
    <template #header-extra>
      <span class="rounded-lg bg-ds-brand-subtle px-2.5 py-1 text-xs font-medium text-ds-fg-brand">{{ $t('adminChartPeriod') }}</span>
    </template>
    <ChartLine
      :categories="categories"
      :series="series"
      :height="260"
      :show-area="true"
    />
  </ChartShell>
</template>

<script setup lang="ts">
import type { ChartSeries } from '~/types/chart'

const props = defineProps<{
  activeTenants: number
  totalTenants: number
}>()

const categories = ['1', '2', '3', '4', '5', '6', '7']

const series = computed<ChartSeries[]>(() => {
  const total = Math.max(props.totalTenants, 1)
  const ratio = props.activeTenants / total
  const base = [0.45, 0.5, 0.48, 0.55, 0.52, ratio, ratio * 0.92]
  const scaled = base.map((v) => Math.round(v * total))
  const compare = base.map((v) => Math.round(v * total * 0.85))
  return [
    { name: 'Current', data: scaled, primary: true },
    { name: 'Before', data: compare, compare: true },
  ]
})
</script>
