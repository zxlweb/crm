<template>
  <section class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4" data-testid="deals-kpi-row">
    <article
      v-for="metric in metrics"
      :key="metric.key"
      class="group relative overflow-hidden rounded-ds-xl border border-ds-border/80 bg-ds-bg/80 p-4 shadow-ds-sm backdrop-blur-sm transition-[box-shadow,border-color] duration-200 hover:border-ds-brand-muted/40 hover:shadow-ds-md sm:p-5"
    >
      <div class="flex items-start gap-3">
        <span
          class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl transition-colors duration-200"
          :class="metric.iconWrapClass"
        >
          <component :is="metric.icon" class="h-5 w-5" aria-hidden="true" />
        </span>
        <div class="min-w-0 flex-1">
          <p class="text-xs font-medium text-ds-fg-muted">{{ metric.label }}</p>
          <p class="mt-0.5 truncate text-xl font-bold tabular-nums tracking-tight text-ds-fg-heading sm:text-2xl">
            {{ metric.value }}
          </p>
          <p v-if="metric.hint" class="mt-1 text-xs text-ds-fg-subtle">{{ metric.hint }}</p>
        </div>
      </div>
      <ChartSparkline
        v-if="metric.sparkline?.length"
        :values="metric.sparkline"
        :tone="metric.sparklineTone ?? 'auto'"
        :width="72"
        :height="28"
        class="pointer-events-none absolute bottom-3 right-3 opacity-90"
      />
    </article>
  </section>
</template>

<script setup lang="ts">
import type { DealPipelineSummary, DealPipelineStage } from '~/types/deal'
import { h } from 'vue'

type SparklineTone = 'auto' | 'up' | 'down' | 'flat' | 'brand'

const props = defineProps<{
  summary: DealPipelineSummary | null
  stages?: DealPipelineStage[]
}>()

const { t } = useI18n()
const { formatDealAmount } = useDealLabels()

function sparklineFromStages(stages: DealPipelineStage[], openOnly = true): number[] {
  const filtered = openOnly
    ? stages.filter((s) => s.stage !== 'won' && s.stage !== 'lost')
    : stages
  const counts = filtered.map((s) => s.count)
  if (counts.length >= 2) return counts
  const total = counts.reduce((a, b) => a + b, 0)
  return [Math.max(0, total - 2), Math.max(0, total - 1), total]
}

const metrics = computed(() => {
  const s = props.summary
  if (!s) return []

  const stageList = props.stages ?? []
  const openSparkline = sparklineFromStages(stageList, true)
  const wonSparkline = [
    Math.max(0, s.won_count_mtd - 1),
    s.won_count_mtd,
    s.won_count_mtd,
    s.won_count_mtd,
  ]

  return [
    {
      key: 'open-count',
      label: t('dealsSummaryOpenCount'),
      value: s.open_count,
      hint: t('dealsSummaryOpenCountHint'),
      iconWrapClass: 'bg-ds-info-subtle text-ds-info',
      sparkline: openSparkline,
      sparklineTone: 'brand' as SparklineTone,
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2' }),
      ]),
    },
    {
      key: 'open-amount',
      label: t('dealsSummaryOpenAmount'),
      value: formatDealAmount(s.open_amount),
      hint: t('dealsSummaryOpenAmountHint'),
      iconWrapClass: 'bg-ds-brand-subtle text-ds-fg-brand',
      sparkline: openSparkline.map((n, i) => n * (i + 1) * 10000),
      sparklineTone: 'up' as SparklineTone,
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z' }),
      ]),
    },
    {
      key: 'won-mtd',
      label: t('dealsSummaryWonMtd'),
      value: s.won_count_mtd,
      hint: t('dealsSummaryWonMtdHint'),
      iconWrapClass: 'bg-ds-success-subtle text-ds-success',
      sparkline: wonSparkline,
      sparklineTone: 'up' as SparklineTone,
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z' }),
      ]),
    },
    {
      key: 'won-amount',
      label: t('dealsSummaryWonAmount'),
      value: formatDealAmount(s.won_amount_mtd),
      hint: t('dealsSummaryWonAmountHint'),
      iconWrapClass: 'bg-ds-warning-subtle text-ds-warning',
      sparkline: wonSparkline.map((n) => n * 80000),
      sparklineTone: 'flat' as SparklineTone,
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M13 7h8m0 0v8m0-8l-8 8-4-4-6 6' }),
      ]),
    },
  ]
})
</script>
