<template>
  <section
    v-if="cards.length > 0"
    class="ds-dash-focus relative"
    data-testid="dashboard-focus-stream"
  >
    <header class="mb-2 flex items-end justify-between gap-3 px-1">
      <div class="flex min-w-0 items-center gap-2">
        <span class="relative flex h-2 w-2" aria-hidden="true">
          <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-ds-brand opacity-60" />
          <span class="relative inline-flex h-2 w-2 rounded-full bg-ds-brand" />
        </span>
        <h2 class="truncate text-sm font-semibold tracking-tight text-ds-fg-heading">
          {{ $t('dashboardFocusTitle') }}
        </h2>
        <span
          class="rounded-full bg-ds-brand-subtle px-2 py-0.5 text-[11px] font-semibold text-ds-fg-brand ring-1 ring-inset ring-ds-brand/20"
        >
          {{ cards.length }}
        </span>
      </div>
      <p class="hidden text-xs text-ds-fg-muted sm:block">{{ $t('dashboardFocusHint') }}</p>
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
        class="ds-dash-focus__track flex gap-3 overflow-x-auto scroll-px-4 snap-x snap-mandatory px-2 pb-3 pt-1"
      >
        <NuxtLink
          v-for="card in cards"
          :key="card.id"
          :to="card.href"
          class="ds-dash-focus__card group relative flex w-[17.5rem] shrink-0 snap-start cursor-pointer flex-col overflow-hidden rounded-2xl border bg-ds-bg-elevated p-3.5 shadow-ds-sm transition-[border-color,box-shadow,transform] duration-200 hover:-translate-y-0.5 hover:shadow-ds-md"
          :class="card.cardClass"
          :data-testid="`dashboard-focus-card-${card.kind}`"
        >
          <span
            class="pointer-events-none absolute inset-x-0 top-0 h-0.5 opacity-80"
            :class="card.barClass"
            aria-hidden="true"
          />

          <div class="flex items-center justify-between gap-2">
            <span
              class="inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-[11px] font-semibold ring-1 ring-inset"
              :class="card.tagClass"
            >
              <UIcon :name="card.tagIcon" class="h-3 w-3" aria-hidden="true" />
              {{ card.tagLabel }}
            </span>
            <span
              v-if="card.metaLabel"
              class="text-[11px] font-medium tabular-nums text-ds-fg-subtle"
            >
              {{ card.metaLabel }}
            </span>
          </div>

          <p
            class="mt-2 line-clamp-2 text-sm font-semibold leading-snug text-ds-fg-heading transition-colors duration-200 group-hover:text-ds-fg-brand"
          >
            {{ card.title }}
          </p>

          <p
            v-if="card.body"
            class="mt-1.5 line-clamp-2 text-xs leading-relaxed text-ds-fg-muted"
          >
            {{ card.body }}
          </p>

          <div class="mt-auto flex items-center justify-between gap-2 pt-3 text-[11px]">
            <span
              class="inline-flex items-center gap-1 font-medium"
              :class="card.signalClass"
            >
              <UIcon :name="card.signalIcon" class="h-3 w-3" aria-hidden="true" />
              {{ card.signalLabel }}
            </span>
            <span
              class="inline-flex items-center gap-1 font-semibold text-ds-fg-brand transition-transform duration-200 group-hover:translate-x-0.5"
            >
              {{ card.ctaLabel }}
              <UIcon name="i-heroicons-arrow-right-20-solid" class="h-3 w-3" aria-hidden="true" />
            </span>
          </div>
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type {
  DashboardInsightItem,
  DashboardInsightVariant,
  PriorityActionItem,
} from '~/types/dashboard'

const props = withDefaults(
  defineProps<{
    priorities: PriorityActionItem[]
    insights: DashboardInsightItem[]
    insightDetailHref?: string
    maxCards?: number
    maxPriorities?: number
  }>(),
  { insightDetailHref: undefined, maxCards: 6, maxPriorities: 2 },
)

const { t } = useI18n()

type FocusKind = 'alert' | 'watch' | 'opportunity' | 'churn' | 'insight'

interface FocusCard {
  id: string
  kind: FocusKind
  href: string
  title: string
  body?: string
  metaLabel?: string
  tagLabel: string
  tagIcon: string
  tagClass: string
  cardClass: string
  barClass: string
  signalLabel: string
  signalIcon: string
  signalClass: string
  ctaLabel: string
  score: number
}

const followCta = computed(() => t('dashboardFocusCtaFollow'))
const reviewCta = computed(() => t('dashboardFocusCtaReview'))

const insightVariantMeta: Record<
  DashboardInsightVariant,
  {
    kind: FocusKind
    tagKey: string
    tagIcon: string
    tagClass: string
    cardClass: string
    barClass: string
    signalIcon: string
    signalClass: string
  }
> = {
  churn: {
    kind: 'churn',
    tagKey: 'dashboardFocusTagChurn',
    tagIcon: 'i-heroicons-exclamation-triangle',
    tagClass: 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/25',
    cardClass: 'border-ds-danger/30',
    barClass: 'bg-ds-danger',
    signalIcon: 'i-heroicons-exclamation-triangle',
    signalClass: 'text-ds-danger',
  },
  opportunity: {
    kind: 'opportunity',
    tagKey: 'dashboardFocusTagOpportunity',
    tagIcon: 'i-heroicons-sparkles',
    tagClass: 'bg-ds-success-subtle text-ds-success ring-ds-success/25',
    cardClass: 'border-ds-success/30',
    barClass: 'bg-ds-success',
    signalIcon: 'i-heroicons-arrow-trending-up',
    signalClass: 'text-ds-success',
  },
  rule: {
    kind: 'insight',
    tagKey: 'dashboardFocusTagInsight',
    tagIcon: 'i-heroicons-light-bulb',
    tagClass: 'bg-ds-brand-subtle text-ds-fg-brand ring-ds-brand/25',
    cardClass: 'border-ds-brand-muted',
    barClass: 'bg-ds-brand',
    signalIcon: 'i-heroicons-cpu-chip',
    signalClass: 'text-ds-fg-brand',
  },
}

interface PriorityMeta {
  kind: FocusKind
  tagKey: string
  tagIcon: string
  tagClass: string
  cardClass: string
  barClass: string
}

const ALERT_META: PriorityMeta = {
  kind: 'alert',
  tagKey: 'dashboardFocusTagAlert',
  tagIcon: 'i-heroicons-fire',
  tagClass: 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/25',
  cardClass: 'border-ds-danger/35',
  barClass: 'bg-ds-danger',
}

const WATCH_META: PriorityMeta = {
  kind: 'watch',
  tagKey: 'dashboardFocusTagWatch',
  tagIcon: 'i-heroicons-eye',
  tagClass: 'bg-ds-warning-subtle text-ds-warning ring-ds-warning/25',
  cardClass: 'border-ds-warning/30',
  barClass: 'bg-ds-warning',
}

interface SignalDescriptor {
  label: string
  icon: string
  cls: string
}

function priorityActivitySignal(days: number | null | undefined): SignalDescriptor {
  if (days === null || days === undefined) {
    return { label: t('dashboardFocusSignalNew'), icon: 'i-heroicons-sparkles', cls: 'text-ds-fg-muted' }
  }
  if (days <= 0) {
    return { label: t('dashboardFocusSignalToday'), icon: 'i-heroicons-bolt', cls: 'text-ds-fg-brand' }
  }
  const tone = days <= 3 ? 'text-ds-warning' : 'text-ds-fg-muted'
  return { label: t('dashboardDaysAgo', { days }), icon: 'i-heroicons-clock', cls: tone }
}

function priorityCard(item: PriorityActionItem): FocusCard {
  const isAlert = item.healthLabel === 'alert'
  const meta = isAlert ? ALERT_META : WATCH_META
  const signal = priorityActivitySignal(item.daysSinceActivity)

  return {
    id: `priority-${item.id}`,
    kind: meta.kind,
    href: item.followHref,
    title: item.title,
    body: item.suggestion,
    metaLabel: t('dashboardFocusScore', { value: item.engagementScore }),
    tagLabel: t(meta.tagKey),
    tagIcon: meta.tagIcon,
    tagClass: meta.tagClass,
    cardClass: meta.cardClass,
    barClass: meta.barClass,
    signalLabel: signal.label,
    signalIcon: signal.icon,
    signalClass: signal.cls,
    ctaLabel: followCta.value,
    score: (isAlert ? 1000 : 600) + item.score,
  }
}

function insightBaseScore(kind: FocusKind, urgent: boolean): number {
  if (urgent) return 900
  if (kind === 'opportunity') return 500
  return 400
}

function insightCard(item: DashboardInsightItem): FocusCard {
  const meta = insightVariantMeta[item.variant] ?? insightVariantMeta.rule
  const href = props.insightDetailHref ?? '/leads?tab=reports'
  const urgent = Boolean(item.urgent)
  return {
    id: `insight-${item.id}`,
    kind: meta.kind,
    href,
    title: item.title,
    body: item.body,
    tagLabel: t(meta.tagKey),
    tagIcon: meta.tagIcon,
    tagClass: meta.tagClass,
    cardClass: meta.cardClass,
    barClass: meta.barClass,
    signalLabel: urgent ? t('dashboardFocusSignalUrgent') : t('dashboardFocusSignalAi'),
    signalIcon: meta.signalIcon,
    signalClass: meta.signalClass,
    ctaLabel: reviewCta.value,
    score: insightBaseScore(meta.kind, urgent) + (urgent ? 100 : 0),
  }
}

const cards = computed<FocusCard[]>(() => {
  // Limit priorities to the top N urgent ones so the stream doesn't visually
  // echo the structured "follow-up queue" section below. Bias toward insights.
  const priorities = [...props.priorities]
    .sort((a, b) => b.score - a.score)
    .slice(0, props.maxPriorities)
    .map(priorityCard)
  const insights = props.insights.map(insightCard)

  // Interleave alerts → urgent insights → watch → other insights
  const list: FocusCard[] = [
    ...priorities.filter((c) => c.kind === 'alert'),
    ...insights.filter((c) => c.kind === 'churn'),
    ...priorities.filter((c) => c.kind === 'watch'),
    ...insights.filter((c) => c.kind !== 'churn'),
  ]

  const seen = new Set<string>()
  const deduped = list.filter((c) => {
    if (seen.has(c.id)) return false
    seen.add(c.id)
    return true
  })

  return deduped
    .sort((a, b) => b.score - a.score)
    .slice(0, props.maxCards)
})
</script>

<style scoped>
.ds-dash-focus__track {
  scrollbar-width: thin;
  scrollbar-color: var(--ds-border) transparent;
}
.ds-dash-focus__track::-webkit-scrollbar {
  height: 6px;
}
.ds-dash-focus__track::-webkit-scrollbar-thumb {
  background: var(--ds-border);
  border-radius: 9999px;
}
.ds-dash-focus__track::-webkit-scrollbar-track {
  background: transparent;
}

.ds-dash-focus__card {
  background-image: linear-gradient(
    180deg,
    color-mix(in srgb, var(--ds-bg-elevated) 88%, var(--ds-bg-muted) 12%) 0%,
    var(--ds-bg-elevated) 60%
  );
}
</style>
