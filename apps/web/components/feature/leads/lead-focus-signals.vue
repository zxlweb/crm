<template>
  <section
    v-if="signals.length > 0"
    class="ds-lead-signals relative"
    data-testid="lead-focus-signals"
  >
    <header class="mb-2 flex items-end justify-between gap-3 px-1">
      <div class="flex min-w-0 items-center gap-2">
        <span class="relative flex h-2 w-2" aria-hidden="true">
          <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-ds-brand opacity-60" />
          <span class="relative inline-flex h-2 w-2 rounded-full bg-ds-brand" />
        </span>
        <h2 class="truncate text-sm font-semibold tracking-tight text-ds-fg-heading">
          {{ $t('leadsSignalsTitle') }}
        </h2>
        <span
          class="rounded-full bg-ds-brand-subtle px-2 py-0.5 text-[11px] font-semibold text-ds-fg-brand ring-1 ring-inset ring-ds-brand/20"
        >
          {{ signals.length }}
        </span>
      </div>
      <p class="hidden text-xs text-ds-fg-muted sm:block">{{ $t('leadsSignalsHint') }}</p>
    </header>

    <div class="relative">
      <div
        class="pointer-events-none absolute inset-y-0 left-0 z-10 w-12 bg-gradient-to-r from-ds-bg to-transparent"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute inset-y-0 right-0 z-10 w-12 bg-gradient-to-l from-ds-bg to-transparent"
        aria-hidden="true"
      />

      <div
        class="ds-lead-signals__track flex gap-3 overflow-x-auto scroll-px-4 snap-x snap-mandatory px-2 pb-3 pt-1"
      >
        <component
          :is="signal.href ? 'NuxtLink' : 'div'"
          v-for="signal in signals"
          :key="signal.id"
          :to="signal.href"
          class="ds-lead-signals__card group relative flex w-[17rem] shrink-0 snap-start flex-col overflow-hidden rounded-2xl border bg-ds-bg-elevated p-3.5 shadow-ds-sm transition-[border-color,box-shadow,transform] duration-200"
          :class="[signal.cardClass, signal.href ? 'cursor-pointer hover:-translate-y-0.5 hover:shadow-ds-md' : '']"
          :data-testid="`lead-signal-${signal.id}`"
        >
          <span
            class="pointer-events-none absolute inset-x-0 top-0 h-0.5 opacity-80"
            :class="signal.barClass"
            aria-hidden="true"
          />

          <div class="flex items-center justify-between gap-2">
            <span
              class="inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-[11px] font-semibold ring-1 ring-inset"
              :class="signal.tagClass"
            >
              <UIcon :name="signal.tagIcon" class="h-3 w-3" aria-hidden="true" />
              {{ signal.tagLabel }}
            </span>
            <span
              v-if="signal.metaLabel"
              class="text-[11px] font-medium tabular-nums text-ds-fg-subtle"
            >
              {{ signal.metaLabel }}
            </span>
          </div>

          <p
            class="mt-2 line-clamp-2 text-sm font-semibold leading-snug text-ds-fg-heading transition-colors duration-200 group-hover:text-ds-fg-brand"
          >
            {{ signal.title }}
          </p>
          <p class="mt-1.5 line-clamp-2 text-xs leading-relaxed text-ds-fg-muted">
            {{ signal.body }}
          </p>

          <div class="mt-auto flex items-center justify-between gap-2 pt-3 text-[11px]">
            <span class="inline-flex items-center gap-1 font-medium" :class="signal.signalClass">
              <UIcon :name="signal.signalIcon" class="h-3 w-3" aria-hidden="true" />
              {{ signal.signalLabel }}
            </span>
            <span
              v-if="signal.cta"
              class="inline-flex items-center gap-1 font-semibold text-ds-fg-brand transition-transform duration-200 group-hover:translate-x-0.5"
            >
              {{ signal.cta }}
              <UIcon name="i-heroicons-arrow-right-20-solid" class="h-3 w-3" aria-hidden="true" />
            </span>
          </div>
        </component>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { Lead } from '~/types/lead'

const LOW_INTENT_TAG = '低意向'

interface SignalCard {
  id: string
  title: string
  body: string
  tagLabel: string
  tagIcon: string
  tagClass: string
  cardClass: string
  barClass: string
  signalLabel: string
  signalIcon: string
  signalClass: string
  metaLabel?: string
  cta?: string
  href?: string
  score: number
}

const props = defineProps<{
  lead: Lead
  churnScore?: number | null
}>()

const { t } = useI18n()

const daysIdle = computed<number | null>(() => {
  if (!props.lead.last_activity_at) return null
  return Math.max(
    0,
    Math.floor((Date.now() - new Date(props.lead.last_activity_at).getTime()) / 86_400_000),
  )
})

const daysToClose = computed<number | null>(() => {
  if (!props.lead.expected_close_date) return null
  const close = new Date(`${props.lead.expected_close_date}T00:00:00`).getTime()
  if (Number.isNaN(close)) return null
  return Math.round((close - Date.now()) / 86_400_000)
})

const ACTIVITY_TRIGGER = '#timeline'

function buildRelationshipLow(): SignalCard {
  return {
    id: 'relationship-low',
    title: t('leadsSignalRelationshipLowTitle'),
    body: t('leadsSignalRelationshipLowBody'),
    tagLabel: t('leadsSignalTagRisk'),
    tagIcon: 'i-heroicons-shield-exclamation',
    tagClass: 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/25',
    cardClass: 'border-ds-danger/30',
    barClass: 'bg-ds-danger',
    signalLabel: t('leadsSignalSignalUrgent'),
    signalIcon: 'i-heroicons-exclamation-triangle',
    signalClass: 'text-ds-danger',
    cta: t('leadsSignalCtaContact'),
    href: ACTIVITY_TRIGGER,
    score: 1000,
  }
}

function buildChurnHigh(score: number): SignalCard {
  return {
    id: 'churn-high',
    title: t('leadsSignalChurnTitle'),
    body: t('leadsSignalChurnBody', { score }),
    metaLabel: `${score}%`,
    tagLabel: t('leadsSignalTagChurn'),
    tagIcon: 'i-heroicons-arrow-trending-down',
    tagClass: 'bg-ds-warning-subtle text-ds-warning ring-ds-warning/25',
    cardClass: 'border-ds-warning/35',
    barClass: 'bg-ds-warning',
    signalLabel: t('leadsSignalSignalAi'),
    signalIcon: 'i-heroicons-cpu-chip',
    signalClass: 'text-ds-warning',
    cta: t('leadsSignalCtaReview'),
    score: 900,
  }
}

function buildIdle(days: number): SignalCard {
  const isCritical = days > 14
  return {
    id: 'idle',
    title: t('leadsSignalIdleTitle', { days }),
    body: t('leadsSignalIdleBody'),
    metaLabel: t('leadsMetricDaysValue', { days }),
    tagLabel: t('leadsSignalTagIdle'),
    tagIcon: 'i-heroicons-clock',
    tagClass: isCritical
      ? 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/25'
      : 'bg-ds-warning-subtle text-ds-warning ring-ds-warning/25',
    cardClass: isCritical ? 'border-ds-danger/30' : 'border-ds-warning/30',
    barClass: isCritical ? 'bg-ds-danger' : 'bg-ds-warning',
    signalLabel: t('leadsSignalSignalAction'),
    signalIcon: 'i-heroicons-bell-alert',
    signalClass: isCritical ? 'text-ds-danger' : 'text-ds-warning',
    cta: t('leadsSignalCtaLogActivity'),
    href: ACTIVITY_TRIGGER,
    score: 700 + days,
  }
}

function buildOverdueClose(days: number): SignalCard {
  return {
    id: 'overdue-close',
    title: t('leadsSignalOverdueTitle', { n: days }),
    body: t('leadsSignalOverdueBody'),
    metaLabel: t('leadsHeroOverdueBy', { n: days }),
    tagLabel: t('leadsSignalTagOverdue'),
    tagIcon: 'i-heroicons-exclamation-triangle',
    tagClass: 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/25',
    cardClass: 'border-ds-danger/30',
    barClass: 'bg-ds-danger',
    signalLabel: t('leadsSignalSignalUrgent'),
    signalIcon: 'i-heroicons-exclamation-triangle',
    signalClass: 'text-ds-danger',
    cta: t('leadsSignalCtaUpdate'),
    score: 850,
  }
}

function buildClosingSoon(days: number): SignalCard {
  const isToday = days === 0
  return {
    id: 'closing-soon',
    title: isToday
      ? t('leadsSignalClosingTodayTitle')
      : t('leadsSignalClosingSoonTitle', { n: days }),
    body: t('leadsSignalClosingSoonBody'),
    metaLabel: isToday ? t('leadsHeroDueToday') : t('leadsHeroDueInDays', { n: days }),
    tagLabel: t('leadsSignalTagClosing'),
    tagIcon: 'i-heroicons-bell-alert',
    tagClass: 'bg-ds-warning-subtle text-ds-warning ring-ds-warning/25',
    cardClass: 'border-ds-warning/35',
    barClass: 'bg-ds-warning',
    signalLabel: t('leadsSignalSignalAction'),
    signalIcon: 'i-heroicons-fire',
    signalClass: 'text-ds-warning',
    cta: t('leadsSignalCtaPrepare'),
    score: 800,
  }
}

function buildEngagementHigh(score: number): SignalCard {
  return {
    id: 'engagement-high',
    title: t('leadsSignalEngagementHighTitle', { score }),
    body: t('leadsSignalEngagementHighBody'),
    metaLabel: t('leadsHeroEngagement', { score }),
    tagLabel: t('leadsSignalTagOpportunity'),
    tagIcon: 'i-heroicons-sparkles',
    tagClass: 'bg-ds-success-subtle text-ds-success ring-ds-success/25',
    cardClass: 'border-ds-success/30',
    barClass: 'bg-ds-success',
    signalLabel: t('leadsSignalSignalOpportunity'),
    signalIcon: 'i-heroicons-arrow-trending-up',
    signalClass: 'text-ds-success',
    cta: t('leadsSignalCtaAdvance'),
    href: ACTIVITY_TRIGGER,
    score: 600,
  }
}

function buildIntentLow(): SignalCard {
  return {
    id: 'intent-low',
    title: t('leadsSignalIntentLowTitle'),
    body: t('leadsSignalIntentLowBody'),
    tagLabel: t('leadsSignalTagIntent'),
    tagIcon: 'i-heroicons-eye-slash',
    tagClass: 'bg-ds-bg-muted text-ds-fg-muted ring-ds-border-muted',
    cardClass: 'border-ds-border-muted',
    barClass: 'bg-ds-fg-subtle',
    signalLabel: t('leadsSignalSignalAi'),
    signalIcon: 'i-heroicons-light-bulb',
    signalClass: 'text-ds-fg-muted',
    cta: t('leadsSignalCtaQualify'),
    score: 500,
  }
}

function buildAmountMissing(): SignalCard {
  return {
    id: 'amount-missing',
    title: t('leadsSignalAmountMissingTitle'),
    body: t('leadsSignalAmountMissingBody'),
    tagLabel: t('leadsSignalTagData'),
    tagIcon: 'i-heroicons-banknotes',
    tagClass: 'bg-ds-info-subtle text-ds-info ring-ds-info/25',
    cardClass: 'border-ds-info/25',
    barClass: 'bg-ds-info',
    signalLabel: t('leadsSignalSignalAction'),
    signalIcon: 'i-heroicons-pencil-square',
    signalClass: 'text-ds-info',
    cta: t('leadsSignalCtaEdit'),
    score: 450,
  }
}

function hasLowIntentTag(tags: readonly string[]): boolean {
  return tags.some(
    (tag) => tag === LOW_INTENT_TAG || tag.toLowerCase() === 'low intent',
  )
}

const signals = computed<SignalCard[]>(() => {
  const out: SignalCard[] = []
  const lead = props.lead
  const churn = props.churnScore
  const idle = daysIdle.value
  const close = daysToClose.value
  const isRiskLow = lead.relationship_health === 'low'

  if (isRiskLow) out.push(buildRelationshipLow())
  if (churn != null && churn >= 55 && !isRiskLow) out.push(buildChurnHigh(churn))
  if (idle != null && idle > 7) out.push(buildIdle(idle))
  if (close != null && close < 0) out.push(buildOverdueClose(Math.abs(close)))
  if (close != null && close >= 0 && close <= 7) out.push(buildClosingSoon(close))
  if (lead.engagement_score >= 70 && !isRiskLow) out.push(buildEngagementHigh(lead.engagement_score))
  if (hasLowIntentTag(lead.tags ?? [])) out.push(buildIntentLow())
  if (lead.amount === 0 && lead.status === 'qualified') out.push(buildAmountMissing())

  return out.sort((a, b) => b.score - a.score).slice(0, 6)
})
</script>

<style scoped>
.ds-lead-signals__track {
  scrollbar-width: thin;
  scrollbar-color: var(--ds-border) transparent;
}
.ds-lead-signals__track::-webkit-scrollbar {
  height: 6px;
}
.ds-lead-signals__track::-webkit-scrollbar-thumb {
  background: var(--ds-border);
  border-radius: 9999px;
}
.ds-lead-signals__track::-webkit-scrollbar-track {
  background: transparent;
}

.ds-lead-signals__card {
  background-image: linear-gradient(
    180deg,
    color-mix(in srgb, var(--ds-bg-elevated) 88%, var(--ds-bg-muted) 12%) 0%,
    var(--ds-bg-elevated) 60%
  );
}
</style>
