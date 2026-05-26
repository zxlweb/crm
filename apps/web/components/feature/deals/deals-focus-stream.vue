<template>
  <section
    v-if="cards.length > 0"
    class="ds-focus-stream relative"
    data-testid="deals-focus-stream"
  >
    <header class="mb-2 flex items-end justify-between gap-3 px-1">
      <div class="flex min-w-0 items-center gap-2">
        <span class="relative flex h-2 w-2" aria-hidden="true">
          <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-ds-brand opacity-60" />
          <span class="relative inline-flex h-2 w-2 rounded-full bg-ds-brand" />
        </span>
        <h2 class="truncate text-sm font-semibold tracking-tight text-ds-fg-heading">
          {{ $t('dealsFocusTitle') }}
        </h2>
        <span class="rounded-full bg-ds-brand-subtle px-2 py-0.5 text-[11px] font-semibold text-ds-fg-brand ring-1 ring-inset ring-ds-brand/20">
          {{ cards.length }}
        </span>
      </div>
      <p class="hidden text-xs text-ds-fg-muted sm:block">{{ $t('dealsFocusHint') }}</p>
    </header>

    <div class="relative">
      <!-- Edge fades (only show on overflow / wider screens) -->
      <div
        class="pointer-events-none absolute inset-y-0 left-0 z-10 w-12 bg-gradient-to-r from-ds-bg to-transparent"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute inset-y-0 right-0 z-10 w-12 bg-gradient-to-l from-ds-bg to-transparent"
        aria-hidden="true"
      />

      <div
        class="ds-focus-stream__track flex gap-3 overflow-x-auto scroll-px-4 snap-x snap-mandatory px-2 pb-3 pt-1"
      >
        <NuxtLink
          v-for="card in cards"
          :key="`${card.kind}-${card.item.id}`"
          :to="`/deals/${card.item.id}`"
          class="ds-focus-stream__card group relative flex w-[16.5rem] shrink-0 snap-start cursor-pointer flex-col overflow-hidden rounded-2xl border bg-ds-bg-elevated p-3 shadow-ds-sm transition-[border-color,box-shadow] duration-200 hover:border-ds-brand-muted hover:shadow-ds-md"
          :class="card.cardClass"
          :data-testid="`deals-focus-card-${card.kind}-${card.item.id}`"
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
              <UIcon :name="card.icon" class="h-3 w-3" aria-hidden="true" />
              {{ card.tagLabel }}
            </span>
            <DealsDealStageBadge :stage="card.stage" variant="plain" />
          </div>

          <p class="mt-2 line-clamp-2 text-sm font-semibold leading-snug text-ds-fg-heading transition-colors duration-200 group-hover:text-ds-fg-brand">
            {{ card.item.title }}
          </p>

          <div class="mt-2.5 flex items-end justify-between gap-2">
            <p class="text-base font-bold tabular-nums tracking-tight text-ds-fg-heading">
              {{ formatDealAmount(card.item.amount, card.item.currency) }}
            </p>
            <p class="text-[11px] tabular-nums text-ds-fg-muted">
              {{ $t('dealsCardProbability', { value: card.item.probability }) }}
            </p>
          </div>

          <div class="mt-2 flex items-center justify-between gap-2 text-[11px]">
            <span class="inline-flex items-center gap-1 text-ds-fg-muted">
              <UIcon name="i-heroicons-calendar-days" class="h-3 w-3" aria-hidden="true" />
              {{ card.dateLabel }}
            </span>
            <span
              class="inline-flex items-center gap-1 font-medium"
              :class="card.signalClass"
            >
              <UIcon :name="card.signalIcon" class="h-3 w-3" aria-hidden="true" />
              {{ card.signalLabel }}
            </span>
          </div>
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { DealPipelineItem, DealPipelineStage, DealStage } from '~/types/deal'

const props = defineProps<{
  stages: DealPipelineStage[]
}>()

const { t, locale } = useI18n()
const { formatDealAmount } = useDealLabels()

type FocusKind = 'overdue' | 'closing' | 'hot'

interface FocusCard {
  kind: FocusKind
  stage: DealStage
  item: DealPipelineItem
  daysToClose: number | null
  score: number
  tagLabel: string
  icon: string
  tagClass: string
  cardClass: string
  barClass: string
  dateLabel: string
  signalLabel: string
  signalIcon: string
  signalClass: string
}

function startOfDay(d: Date): number {
  const x = new Date(d)
  x.setHours(0, 0, 0, 0)
  return x.getTime()
}

function daysBetween(iso: string | null): number | null {
  if (!iso) return null
  const close = new Date(`${iso}T00:00:00`)
  if (Number.isNaN(close.getTime())) return null
  const todayStart = startOfDay(new Date())
  const closeStart = startOfDay(close)
  return Math.round((closeStart - todayStart) / 86_400_000)
}

function formatRelativeDate(iso: string | null, days: number | null): string {
  if (!iso || days === null) return t('dealsFocusNoDate')
  const fmt = new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    month: 'short',
    day: 'numeric',
  }).format(new Date(`${iso}T00:00:00`))
  return fmt
}

function signalForDays(days: number | null): { label: string; icon: string; cls: string } {
  if (days === null) return { label: t('dealsFocusNoDate'), icon: 'i-heroicons-clock', cls: 'text-ds-fg-muted' }
  if (days < 0) return { label: t('dealsFocusOverdueBy', { n: Math.abs(days) }), icon: 'i-heroicons-exclamation-triangle', cls: 'text-ds-danger' }
  if (days === 0) return { label: t('dealsFocusDueToday'), icon: 'i-heroicons-bell-alert', cls: 'text-ds-warning' }
  if (days <= 7) return { label: t('dealsFocusInDays', { n: days }), icon: 'i-heroicons-clock', cls: 'text-ds-warning' }
  return { label: t('dealsFocusInDays', { n: days }), icon: 'i-heroicons-clock', cls: 'text-ds-fg-muted' }
}

const cards = computed<FocusCard[]>(() => {
  const openStages = props.stages.filter((s) => s.stage !== 'won' && s.stage !== 'lost')
  const collected: FocusCard[] = []
  const seen = new Set<string>()

  const overdue: FocusCard[] = []
  const closing: FocusCard[] = []
  const hot: FocusCard[] = []

  for (const col of openStages) {
    for (const item of col.items) {
      const days = daysBetween(item.expected_close_date)
      const signal = signalForDays(days)
      const base: Omit<FocusCard, 'kind' | 'tagLabel' | 'icon' | 'tagClass' | 'cardClass' | 'barClass'> & {
        amountScore: number
      } = {
        stage: col.stage,
        item,
        daysToClose: days,
        score: 0,
        dateLabel: formatRelativeDate(item.expected_close_date, days),
        signalLabel: signal.label,
        signalIcon: signal.icon,
        signalClass: signal.cls,
        amountScore: Math.log10(Math.max(item.amount, 1)) * 10,
      }

      if (days !== null && days < 0) {
        overdue.push({
          ...base,
          kind: 'overdue',
          score: 1000 - days + base.amountScore,
          tagLabel: t('dealsFocusTagOverdue'),
          icon: 'i-heroicons-exclamation-triangle',
          tagClass: 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/25',
          cardClass: 'border-ds-danger/30',
          barClass: 'bg-ds-danger',
        })
      } else if (days !== null && days >= 0 && days <= 7) {
        closing.push({
          ...base,
          kind: 'closing',
          score: 500 - days * 5 + base.amountScore,
          tagLabel: days === 0 ? t('dealsFocusTagToday') : t('dealsFocusTagSoon'),
          icon: 'i-heroicons-bell-alert',
          tagClass: 'bg-ds-warning-subtle text-ds-warning ring-ds-warning/25',
          cardClass: 'border-ds-warning/30',
          barClass: 'bg-ds-warning',
        })
      }

      hot.push({
        ...base,
        kind: 'hot',
        score: item.probability * Math.log10(Math.max(item.amount, 1)),
        tagLabel: t('dealsFocusTagHot'),
        icon: 'i-heroicons-fire',
        tagClass: 'bg-ds-brand-subtle text-ds-fg-brand ring-ds-brand/25',
        cardClass: 'border-ds-brand-muted',
        barClass: 'bg-ds-brand-gradient',
      })
    }
  }

  const push = (card: FocusCard) => {
    if (seen.has(card.item.id)) return
    seen.add(card.item.id)
    collected.push(card)
  }

  overdue.sort((a, b) => b.score - a.score).forEach(push)
  closing.sort((a, b) => b.score - a.score).forEach(push)
  hot
    .sort((a, b) => b.score - a.score)
    .slice(0, 4)
    .forEach(push)

  return collected.slice(0, 8)
})
</script>

<style scoped>
.ds-focus-stream__track {
  scrollbar-width: thin;
  scrollbar-color: var(--ds-border) transparent;
}
.ds-focus-stream__track::-webkit-scrollbar {
  height: 6px;
}
.ds-focus-stream__track::-webkit-scrollbar-thumb {
  background: var(--ds-border);
  border-radius: 9999px;
}
.ds-focus-stream__track::-webkit-scrollbar-track {
  background: transparent;
}

.ds-focus-stream__card {
  background-image: linear-gradient(
    180deg,
    color-mix(in srgb, var(--ds-bg-elevated) 88%, var(--ds-bg-muted) 12%) 0%,
    var(--ds-bg-elevated) 60%
  );
}

.bg-ds-brand-gradient {
  background: var(--ds-brand-gradient);
}
</style>
