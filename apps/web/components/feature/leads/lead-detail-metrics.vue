<template>
  <section
    class="overflow-x-auto pb-0.5"
    data-testid="lead-detail-metrics"
    aria-label="lead metrics"
  >
    <div class="grid min-w-[36rem] grid-cols-4 gap-3 sm:min-w-0">
    <CardMetric
      class="min-w-0"
      density="compact"
      :label="$t('leadsFieldAmount')"
      :value="amountDisplay"
      value-display="full"
      :compare-label="$t('leadsMetricAmountHint')"
      :trend="amountTrend"
      :trend-direction="lead.amount > 0 ? 'up' : 'flat'"
      icon-tone="neutral"
    >
      <template #icon>
        <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </template>
    </CardMetric>

    <CardMetric
      class="min-w-0"
      density="compact"
      :label="$t('leadsColEngagement')"
      :value="lead.engagement_score"
      :compare-label="$t('leadsMetricEngagementHint')"
      :trend="engagementTrend"
      :trend-direction="engagementDirection"
      icon-tone="neutral"
    >
      <template #icon>
        <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
        </svg>
      </template>
    </CardMetric>

    <CardMetric
      class="min-w-0"
      density="compact"
      :label="$t('leadsMetricDaysIdle')"
      :value="daysIdle"
      :compare-label="$t('leadsFieldLastActivity')"
      :trend="lastActivityShort"
      :trend-direction="daysIdleDirection"
      icon-tone="neutral"
    >
      <template #icon>
        <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
      </template>
    </CardMetric>

    <CardMetric
      class="min-w-0"
      density="compact"
      :label="$t('leadsFieldExpectedClose')"
      :value="closeValueShort"
      value-display="full"
      :compare-label="closeCompareLabel"
      :trend="$t(`lifecycle.${lead.lifecycle_stage}`)"
      trend-direction="flat"
      icon-tone="neutral"
    >
      <template #icon>
        <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </template>
    </CardMetric>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { Lead } from '~/types/lead'

type TrendDirection = 'up' | 'down' | 'flat'

const props = defineProps<{
  lead: Lead
}>()

const { t, locale } = useI18n()

const amountDisplay = computed(() => {
  if (props.lead.amount <= 0) return t('leadsAmountUnset')
  return new Intl.NumberFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    style: 'currency',
    currency: 'CNY',
    maximumFractionDigits: 0,
  }).format(props.lead.amount)
})

const amountTrend = computed(() =>
  props.lead.amount > 0 ? t('leadsMetricAmountSet') : t('leadsMetricAmountUnset'),
)

const engagementDirection = computed<TrendDirection>(() => {
  if (props.lead.engagement_score >= 60) return 'up'
  if (props.lead.engagement_score < 35) return 'down'
  return 'flat'
})

const engagementTrend = computed(() => {
  if (props.lead.engagement_score >= 60) return t('leadsEngagementHigh')
  if (props.lead.engagement_score < 35) return t('leadsEngagementLow')
  return t('leadsEngagementMid')
})

const daysIdle = computed(() => {
  if (!props.lead.last_activity_at) return '—'
  const diff = Math.max(
    0,
    Math.floor((Date.now() - new Date(props.lead.last_activity_at).getTime()) / 86400000),
  )
  return t('leadsMetricDaysValue', { days: diff })
})

const daysIdleDirection = computed<TrendDirection>(() => {
  if (!props.lead.last_activity_at) return 'flat'
  const diff = Math.floor((Date.now() - new Date(props.lead.last_activity_at).getTime()) / 86400000)
  if (diff > 14) return 'down'
  if (diff <= 3) return 'up'
  return 'flat'
})

const lastActivityShort = computed(() => {
  if (!props.lead.last_activity_at) return t('leadsMetricNoActivity')
  return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    month: 'short',
    day: 'numeric',
  }).format(new Date(props.lead.last_activity_at))
})

/** 单行四列：日期用短格式 */
const closeValueShort = computed(() => {
  if (!props.lead.expected_close_date) return t('leadsCloseUnset')
  const d = new Date(props.lead.expected_close_date)
  if (locale.value === 'zh') {
    return `${d.getMonth() + 1}/${d.getDate()}`
  }
  return new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric' }).format(d)
})

const closeCompareLabel = computed(() => {
  if (!props.lead.expected_close_date) return t('leadsColLifecycle')
  const d = new Date(props.lead.expected_close_date)
  const full = new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  }).format(d)
  return full
})
</script>
