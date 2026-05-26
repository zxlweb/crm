<template>
  <section
    class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4"
    data-testid="lead-detail-metrics"
    :aria-label="$t('leadsMetricsAria')"
  >
    <article
      v-for="metric in metrics"
      :key="metric.key"
      class="ds-lead-metric group relative overflow-hidden rounded-2xl border border-ds-border/80 bg-ds-bg-elevated p-4 shadow-ds-sm transition-[border-color,box-shadow,transform] duration-200 hover:-translate-y-0.5 hover:border-ds-brand-muted/60 hover:shadow-ds-md sm:p-5"
      :class="metric.featured ? 'ds-lead-metric--featured' : ''"
    >
      <span
        v-if="metric.featured"
        class="pointer-events-none absolute inset-x-0 top-0 h-0.5"
        :style="{ background: 'var(--ds-brand-gradient)' }"
        aria-hidden="true"
      />
      <span
        v-if="metric.featured"
        class="pointer-events-none absolute -right-12 -top-10 h-32 w-32 rounded-full opacity-20 blur-3xl"
        :style="{ background: 'var(--ds-brand-gradient)' }"
        aria-hidden="true"
      />
      <span
        v-else-if="metric.accentClass"
        class="pointer-events-none absolute inset-x-0 top-0 h-0.5 opacity-80"
        :class="metric.accentClass"
        aria-hidden="true"
      />

      <div class="relative flex items-start gap-3">
        <span
          class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl transition-colors duration-200"
          :class="metric.iconWrapClass"
          :style="metric.iconStyle"
        >
          <UIcon :name="metric.icon" class="h-5 w-5" aria-hidden="true" />
        </span>
        <div class="min-w-0 flex-1">
          <div class="flex items-center justify-between gap-2">
            <p class="truncate text-xs font-medium uppercase tracking-wider text-ds-fg-muted">
              {{ metric.label }}
            </p>
            <span
              v-if="metric.trend"
              class="inline-flex items-center gap-1 rounded-full px-1.5 py-0.5 text-[10px] font-semibold leading-none ring-1 ring-inset"
              :class="metric.trendClass"
            >
              <UIcon :name="metric.trendIcon" class="h-2.5 w-2.5" aria-hidden="true" />
              <span class="tabular-nums">{{ metric.trend }}</span>
            </span>
          </div>
          <p class="mt-1 truncate text-2xl font-extrabold tabular-nums tracking-tight text-ds-fg-heading sm:text-[1.65rem]">
            <span v-if="metric.featured" class="bg-clip-text" :style="brandGradientText">
              {{ metric.value }}
            </span>
            <template v-else>{{ metric.value }}</template>
          </p>
          <p v-if="metric.hint" class="mt-1 text-xs text-ds-fg-subtle">{{ metric.hint }}</p>

          <!-- Progress bar for featured engagement metric -->
          <div
            v-if="metric.bar"
            class="mt-2 h-1 overflow-hidden rounded-full bg-ds-bg-muted"
          >
            <div
              class="h-full rounded-full transition-all duration-700 ease-out"
              :class="metric.bar.fillClass"
              :style="{ width: `${metric.bar.value}%` }"
            />
          </div>
        </div>
      </div>
    </article>
  </section>
</template>

<script setup lang="ts">
import type { Lead } from '~/types/lead'

type TrendDirection = 'up' | 'down' | 'flat'

interface LeadMetric {
  key: string
  label: string
  value: string
  hint?: string
  featured: boolean
  icon: string
  iconWrapClass: string
  iconStyle?: Record<string, string>
  accentClass?: string
  trend?: string
  trendIcon: string
  trendClass: string
  bar?: { value: number; fillClass: string }
}

const props = defineProps<{
  lead: Lead
}>()

const { t, locale } = useI18n()

const brandGradientText = {
  background: 'var(--ds-brand-gradient)',
  '-webkit-background-clip': 'text',
  'background-clip': 'text',
  color: 'transparent',
}

function trendClassFor(direction: TrendDirection): string {
  switch (direction) {
    case 'up':
      return 'bg-ds-success-subtle text-ds-success ring-ds-success/25'
    case 'down':
      return 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/25'
    default:
      return 'bg-ds-bg-muted text-ds-fg-muted ring-ds-border-muted'
  }
}

function trendIconFor(direction: TrendDirection): string {
  switch (direction) {
    case 'up':
      return 'i-heroicons-arrow-trending-up'
    case 'down':
      return 'i-heroicons-arrow-trending-down'
    default:
      return 'i-heroicons-minus'
  }
}

const amountDisplay = computed(() => {
  if (props.lead.amount <= 0) return t('leadsAmountUnset')
  return new Intl.NumberFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    style: 'currency',
    currency: 'CNY',
    maximumFractionDigits: 0,
  }).format(props.lead.amount)
})

const engagementDir = computed<TrendDirection>(() => {
  if (props.lead.engagement_score >= 60) return 'up'
  if (props.lead.engagement_score < 35) return 'down'
  return 'flat'
})

const engagementTrendLabel = computed(() => {
  if (props.lead.engagement_score >= 60) return t('leadsEngagementHigh')
  if (props.lead.engagement_score < 35) return t('leadsEngagementLow')
  return t('leadsEngagementMid')
})

const daysIdle = computed<number | null>(() => {
  if (!props.lead.last_activity_at) return null
  return Math.max(
    0,
    Math.floor((Date.now() - new Date(props.lead.last_activity_at).getTime()) / 86_400_000),
  )
})

const daysIdleDir = computed<TrendDirection>(() => {
  if (daysIdle.value == null) return 'flat'
  if (daysIdle.value > 14) return 'down'
  if (daysIdle.value <= 3) return 'up'
  return 'flat'
})

const lastActivityShort = computed(() => {
  if (!props.lead.last_activity_at) return t('leadsMetricNoActivity')
  return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    month: 'short',
    day: 'numeric',
  }).format(new Date(props.lead.last_activity_at))
})

const daysToClose = computed<number | null>(() => {
  if (!props.lead.expected_close_date) return null
  const close = new Date(`${props.lead.expected_close_date}T00:00:00`).getTime()
  if (Number.isNaN(close)) return null
  return Math.round((close - Date.now()) / 86_400_000)
})

const closeValueShort = computed(() => {
  if (!props.lead.expected_close_date) return t('leadsCloseUnset')
  const d = new Date(props.lead.expected_close_date)
  if (locale.value === 'zh') {
    return `${d.getMonth() + 1}/${d.getDate()}`
  }
  return new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric' }).format(d)
})

const closeHint = computed(() => {
  if (!props.lead.expected_close_date) return t('leadsColLifecycle')
  return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  }).format(new Date(props.lead.expected_close_date))
})

const closeDir = computed<TrendDirection>(() => {
  if (daysToClose.value == null) return 'flat'
  if (daysToClose.value < 0) return 'down'
  if (daysToClose.value <= 7) return 'up'
  return 'flat'
})

const closeTrendLabel = computed(() => {
  if (daysToClose.value == null) return ''
  if (daysToClose.value < 0) return t('leadsHeroOverdueBy', { n: Math.abs(daysToClose.value) })
  if (daysToClose.value === 0) return t('leadsHeroDueToday')
  return t('leadsHeroDueInDays', { n: daysToClose.value })
})

const engagementBarFill = computed(() => {
  if (props.lead.engagement_score >= 60) return 'bg-ds-success'
  if (props.lead.engagement_score < 35) return 'bg-ds-danger'
  return 'bg-ds-warning'
})

const idleAccent = computed(() => {
  if (daysIdle.value == null) return 'bg-ds-fg-subtle'
  if (daysIdle.value > 14) return 'bg-ds-danger'
  if (daysIdle.value > 7) return 'bg-ds-warning'
  return 'bg-ds-success'
})

const closeAccent = computed(() => {
  if (daysToClose.value == null) return 'bg-ds-fg-subtle'
  if (daysToClose.value < 0) return 'bg-ds-danger'
  if (daysToClose.value <= 7) return 'bg-ds-warning'
  return 'bg-ds-info'
})

function idleIconWrap(days: number | null): string {
  if (days == null) return 'bg-ds-bg-muted text-ds-fg-muted'
  if (days > 14) return 'bg-ds-danger-subtle text-ds-danger'
  if (days > 7) return 'bg-ds-warning-subtle text-ds-warning'
  return 'bg-ds-success-subtle text-ds-success'
}

function closeIconWrap(days: number | null): string {
  if (days == null) return 'bg-ds-bg-muted text-ds-fg-muted'
  if (days < 0) return 'bg-ds-danger-subtle text-ds-danger'
  if (days <= 7) return 'bg-ds-warning-subtle text-ds-warning'
  return 'bg-ds-info-subtle text-ds-info'
}

const metrics = computed<LeadMetric[]>(() => [
  {
    key: 'amount',
    label: t('leadsFieldAmount'),
    value: amountDisplay.value,
    hint: t('leadsMetricAmountHint'),
    featured: false,
    icon: 'i-heroicons-banknotes',
    iconWrapClass: 'bg-ds-bg-muted text-ds-fg-muted',
    accentClass: 'bg-ds-fg-subtle',
    trend: props.lead.amount > 0 ? t('leadsMetricAmountSet') : t('leadsMetricAmountUnset'),
    trendIcon: props.lead.amount > 0 ? 'i-heroicons-check-circle' : 'i-heroicons-pencil-square',
    trendClass: trendClassFor(props.lead.amount > 0 ? 'flat' : 'down'),
  },
  {
    key: 'engagement',
    label: t('leadsColEngagement'),
    value: String(props.lead.engagement_score),
    hint: t('leadsMetricEngagementHint'),
    featured: true,
    icon: 'i-heroicons-bolt',
    iconWrapClass: 'text-ds-on-brand shadow-ds-brand',
    iconStyle: { background: 'var(--ds-brand-gradient)' },
    trend: engagementTrendLabel.value,
    trendIcon: trendIconFor(engagementDir.value),
    trendClass: trendClassFor(engagementDir.value),
    bar: {
      value: Math.max(0, Math.min(100, props.lead.engagement_score)),
      fillClass: engagementBarFill.value,
    },
  },
  {
    key: 'idle',
    label: t('leadsMetricDaysIdle'),
    value:
      daysIdle.value == null
        ? '—'
        : t('leadsMetricDaysValue', { days: daysIdle.value }),
    hint: lastActivityShort.value,
    featured: false,
    icon: 'i-heroicons-clock',
    iconWrapClass: idleIconWrap(daysIdle.value),
    accentClass: idleAccent.value,
    trend: daysIdle.value == null ? t('leadsMetricNoActivity') : lastActivityShort.value,
    trendIcon: trendIconFor(daysIdleDir.value),
    trendClass: trendClassFor(daysIdleDir.value),
  },
  {
    key: 'close',
    label: t('leadsFieldExpectedClose'),
    value: closeValueShort.value,
    hint: closeHint.value,
    featured: false,
    icon: 'i-heroicons-calendar-days',
    iconWrapClass: closeIconWrap(daysToClose.value),
    accentClass: closeAccent.value,
    trend: closeTrendLabel.value || t(`lifecycle.${props.lead.lifecycle_stage}`),
    trendIcon: trendIconFor(closeDir.value),
    trendClass: trendClassFor(closeDir.value),
  },
])
</script>

<style scoped>
.ds-lead-metric--featured {
  background-image: linear-gradient(
    180deg,
    color-mix(in srgb, var(--ds-bg-elevated) 92%, var(--ds-brand) 8%) 0%,
    var(--ds-bg-elevated) 70%
  );
}
</style>
