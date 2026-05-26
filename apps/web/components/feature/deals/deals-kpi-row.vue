<template>
  <section class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4" data-testid="deals-kpi-row">
    <article
      v-for="metric in metrics"
      :key="metric.key"
      class="ds-kpi-card group relative overflow-hidden rounded-2xl border border-ds-border/80 bg-ds-bg-elevated p-4 shadow-ds-sm transition-[border-color,box-shadow,transform] duration-200 hover:-translate-y-0.5 hover:border-ds-brand-muted/60 hover:shadow-ds-md sm:p-5"
      :class="metric.featured ? 'ds-kpi-card--featured' : ''"
    >
      <!-- Featured top gradient bar -->
      <span
        v-if="metric.featured"
        class="pointer-events-none absolute inset-x-0 top-0 h-0.5"
        :style="{ background: 'var(--ds-brand-gradient)' }"
        aria-hidden="true"
      />
      <!-- Ambient glow for featured -->
      <span
        v-if="metric.featured"
        class="pointer-events-none absolute -right-12 -top-10 h-32 w-32 rounded-full opacity-20 blur-3xl"
        :style="{ background: 'var(--ds-brand-gradient)' }"
        aria-hidden="true"
      />

      <div class="relative flex items-start gap-3">
        <span
          class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl transition-colors duration-200"
          :class="metric.iconWrapClass"
          :style="metric.iconStyle"
        >
          <component :is="metric.icon" class="h-5 w-5" aria-hidden="true" />
        </span>
        <div class="min-w-0 flex-1">
          <div class="flex items-center justify-between gap-2">
            <p class="text-xs font-medium uppercase tracking-wider text-ds-fg-muted">
              {{ metric.label }}
            </p>
            <span
              v-if="metric.trend"
              class="inline-flex items-center gap-1 rounded-full px-1.5 py-0.5 text-[10px] font-semibold tabular-nums leading-none ring-1 ring-inset"
              :class="metric.trendClass"
            >
              <UIcon :name="metric.trendIcon" class="h-2.5 w-2.5" aria-hidden="true" />
              {{ metric.trend }}
            </span>
          </div>
          <p class="mt-1 truncate text-2xl font-extrabold tabular-nums tracking-tight text-ds-fg-heading sm:text-[1.65rem]">
            {{ metric.value }}
          </p>
          <p v-if="metric.hint" class="mt-1 text-xs text-ds-fg-subtle">{{ metric.hint }}</p>
        </div>
      </div>

      <ChartSparkline
        v-if="metric.sparkline?.length"
        :values="metric.sparkline"
        :tone="metric.sparklineTone ?? 'auto'"
        :width="84"
        :height="32"
        class="pointer-events-none absolute bottom-3 right-3 opacity-90 transition-opacity duration-200 group-hover:opacity-100"
      />
    </article>
  </section>
</template>

<script setup lang="ts">
import type { DealPipelineSummary, DealPipelineStage } from '~/types/deal'
import { h, type VNode } from 'vue'

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

function trendClassFor(direction: 'up' | 'down' | 'flat'): string {
  switch (direction) {
    case 'up':
      return 'bg-ds-success-subtle text-ds-success ring-ds-success/25'
    case 'down':
      return 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/25'
    case 'flat':
    default:
      return 'bg-ds-bg-muted text-ds-fg-muted ring-ds-border-muted'
  }
}

function wonAmountTrendClass(amount: number): string {
  if (amount > 0) return 'bg-ds-warning-subtle text-ds-warning ring-ds-warning/25'
  return trendClassFor('flat')
}

function trendIconFor(direction: 'up' | 'down' | 'flat'): string {
  switch (direction) {
    case 'up':
      return 'i-heroicons-arrow-trending-up'
    case 'down':
      return 'i-heroicons-arrow-trending-down'
    case 'flat':
    default:
      return 'i-heroicons-minus'
  }
}

interface KpiMetric {
  key: string
  label: string
  value: number | string
  hint: string
  featured: boolean
  iconWrapClass: string
  iconStyle?: Record<string, string>
  sparkline: number[]
  sparklineTone: SparklineTone
  trend?: string
  trendIcon: string
  trendClass: string
  icon: VNode
}

const metrics = computed<KpiMetric[]>(() => {
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

  const lateStageCount = stageList
    .filter((c) => c.stage === 'proposal' || c.stage === 'negotiation')
    .reduce((sum, c) => sum + c.count, 0)
  const earlyStageCount = stageList
    .filter((c) => c.stage === 'qualification')
    .reduce((sum, c) => sum + c.count, 0)

  const openTrendDir: 'up' | 'down' | 'flat' = lateStageCount > earlyStageCount ? 'up' : lateStageCount === earlyStageCount ? 'flat' : 'down'

  return [
    {
      key: 'open-count',
      label: t('dealsSummaryOpenCount'),
      value: s.open_count,
      hint: t('dealsSummaryOpenCountHint'),
      featured: false,
      iconWrapClass: 'bg-ds-info-subtle text-ds-info',
      sparkline: openSparkline,
      sparklineTone: 'brand' as SparklineTone,
      trend: t('dealsKpiTrendActive', { n: lateStageCount }),
      trendIcon: trendIconFor(openTrendDir),
      trendClass: trendClassFor(openTrendDir),
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2' }),
      ]),
    },
    {
      key: 'open-amount',
      label: t('dealsSummaryOpenAmount'),
      value: formatDealAmount(s.open_amount),
      hint: t('dealsSummaryOpenAmountHint'),
      featured: true,
      iconWrapClass: 'text-ds-on-brand shadow-ds-brand',
      iconStyle: { background: 'var(--ds-brand-gradient)' },
      sparkline: openSparkline.map((n, i) => n * (i + 1) * 10000),
      sparklineTone: 'up' as SparklineTone,
      trend: t('dealsKpiTrendPipeline'),
      trendIcon: 'i-heroicons-bolt',
      trendClass: 'bg-ds-brand-subtle text-ds-fg-brand ring-ds-brand/25',
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z' }),
      ]),
    },
    {
      key: 'won-mtd',
      label: t('dealsSummaryWonMtd'),
      value: s.won_count_mtd,
      hint: t('dealsSummaryWonMtdHint'),
      featured: false,
      iconWrapClass: 'bg-ds-success-subtle text-ds-success',
      sparkline: wonSparkline,
      sparklineTone: 'up' as SparklineTone,
      trend: s.won_count_mtd > 0 ? t('dealsKpiTrendMomentum') : t('dealsKpiTrendQuiet'),
      trendIcon: s.won_count_mtd > 0 ? 'i-heroicons-arrow-trending-up' : 'i-heroicons-minus',
      trendClass: trendClassFor(s.won_count_mtd > 0 ? 'up' : 'flat'),
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z' }),
      ]),
    },
    {
      key: 'won-amount',
      label: t('dealsSummaryWonAmount'),
      value: formatDealAmount(s.won_amount_mtd),
      hint: t('dealsSummaryWonAmountHint'),
      featured: false,
      iconWrapClass: 'bg-ds-warning-subtle text-ds-warning',
      sparkline: wonSparkline.map((n) => n * 80000),
      sparklineTone: 'flat' as SparklineTone,
      trend: s.won_amount_mtd > 0 ? t('dealsKpiTrendRevenue') : t('dealsKpiTrendQuiet'),
      trendIcon: s.won_amount_mtd > 0 ? 'i-heroicons-trophy' : 'i-heroicons-minus',
      trendClass: wonAmountTrendClass(s.won_amount_mtd),
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M13 7h8m0 0v8m0-8l-8 8-4-4-6 6' }),
      ]),
    },
  ]
})
</script>

<style scoped>
.ds-kpi-card--featured {
  background-image: linear-gradient(
    180deg,
    color-mix(in srgb, var(--ds-bg-elevated) 92%, var(--ds-brand) 8%) 0%,
    var(--ds-bg-elevated) 70%
  );
}
</style>
