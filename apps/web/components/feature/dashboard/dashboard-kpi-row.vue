<template>
  <section
    data-testid="dashboard-kpi-row"
    :aria-labelledby="variant === 'hero' ? undefined : 'dashboard-pulse-heading'"
  >
    <h3
      v-if="variant !== 'hero'"
      id="dashboard-pulse-heading"
      class="mb-3 text-sm font-semibold text-ds-fg-heading"
    >
      {{ $t('dashboardPulseTitle') }}
    </h3>

    <div :class="gridClass">
      <NuxtLink
        v-for="kpi in kpis"
        :key="kpi.key"
        :to="kpi.href"
        class="group relative block min-w-0 cursor-pointer rounded-ds-xl focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ds-brand"
      >
        <CardMetric
          class="h-full transition-[box-shadow,border-color,background-color] duration-200 group-hover:border-ds-brand-muted/50 group-hover:bg-ds-bg-muted/30 group-hover:shadow-ds-md"
          :class="variant === 'hero' ? 'p-4 sm:p-5' : 'pr-7'"
          density="compact"
          :icon-tone="kpi.iconTone"
          :label="kpi.label"
          :value="kpi.value"
          :compare-label="kpi.hint"
          :trend="kpi.trend"
          :trend-direction="kpi.trendDirection"
        >
          <template #icon>
            <component :is="kpi.icon" class="h-5 w-5" aria-hidden="true" />
          </template>
        </CardMetric>
        <UIcon
          v-if="variant !== 'hero'"
          name="i-heroicons-chevron-right-20-solid"
          class="pointer-events-none absolute right-3 top-4 h-4 w-4 shrink-0 text-ds-fg-subtle opacity-0 transition-opacity duration-200 group-hover:opacity-100"
          aria-hidden="true"
        />
      </NuxtLink>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { DashboardKpiTrends } from '~/types/dashboard'
import { h } from 'vue'

type TrendDirection = 'up' | 'down' | 'flat'

const props = withDefaults(
  defineProps<{
    leadsTotal: number
    accountsTotal: number
    avgEngagement: number
    atRiskTotal: number
    kpiTrends: DashboardKpiTrends
    variant?: 'default' | 'hero'
  }>(),
  {
    variant: 'default',
  },
)

const { t } = useI18n()

const gridClass = computed(() => {
  if (props.variant === 'hero') {
    return 'grid grid-cols-2 gap-3 sm:grid-cols-4 xl:gap-4'
  }
  return 'grid min-w-[36rem] grid-cols-4 gap-3 overflow-x-auto pb-0.5 sm:min-w-0'
})

const kpis = computed(() => {
  const atRiskTrend = props.atRiskTotal > 0 ? t('metricNeedsAttention') : t('metricAllClear')
  const atRiskTrendDirection: TrendDirection =
    props.atRiskTotal > 0 ? 'down' : 'flat'

  const leadsTrend =
    props.kpiTrends.leadsWeeklyTouch > 0
      ? t('dashboardKpiTrendWeeklyActive', { count: props.kpiTrends.leadsWeeklyTouch })
      : undefined
  const accountsTrend =
    props.kpiTrends.accountsWeeklyTouch > 0
      ? t('dashboardKpiTrendWeeklyActive', { count: props.kpiTrends.accountsWeeklyTouch })
      : undefined

  let engagementTrend: string | undefined
  let engagementTrendDirection: TrendDirection = 'flat'
  if (props.kpiTrends.engagementDirection === 'up') {
    engagementTrend = t('dashboardKpiTrendEngagementUp', { delta: props.kpiTrends.engagementDelta })
    engagementTrendDirection = 'up'
  } else if (props.kpiTrends.engagementDirection === 'down') {
    engagementTrend = t('dashboardKpiTrendEngagementDown', { delta: props.kpiTrends.engagementDelta })
    engagementTrendDirection = 'down'
  } else if (props.avgEngagement > 0) {
    engagementTrend = t('dashboardKpiTrendEngagementFlat')
    engagementTrendDirection = 'flat'
  }

  return [
    {
      key: 'leads',
      label: t('dashboardKpiLeads'),
      value: props.leadsTotal,
      hint: t('dashboardKpiLeadsHint'),
      href: '/leads',
      trend: leadsTrend,
      trendDirection: 'flat' as TrendDirection,
      iconTone: 'info' as const,
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z' }),
      ]),
    },
    {
      key: 'accounts',
      label: t('dashboardKpiAccounts'),
      value: props.accountsTotal,
      hint: t('dashboardKpiAccountsHint'),
      href: '/accounts',
      trend: accountsTrend,
      trendDirection: 'flat' as TrendDirection,
      iconTone: 'calendar' as const,
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4' }),
      ]),
    },
    {
      key: 'engagement',
      label: t('dashboardKpiEngagement'),
      value: props.avgEngagement,
      hint: t('dashboardKpiEngagementHint'),
      href: '/leads?tab=reports',
      trend: engagementTrend,
      trendDirection: engagementTrendDirection,
      iconTone: 'brand' as const,
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M13 7h8m0 0v8m0-8l-8 8-4-4-6 6' }),
      ]),
    },
    {
      key: 'at-risk',
      label: t('dashboardKpiAtRisk'),
      value: props.atRiskTotal,
      hint: t('dashboardKpiAtRiskHint'),
      href: '/leads?health=low',
      trend: atRiskTrend,
      trendDirection: atRiskTrendDirection,
      iconTone: 'accent' as const,
      icon: h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.75' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z' }),
      ]),
    },
  ]
})
</script>
