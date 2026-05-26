<template>
  <section
    class="ds-deals-pipeline-hero relative overflow-hidden rounded-2xl border border-ds-border bg-ds-bg-elevated/85 p-4 shadow-ds-sm backdrop-blur-sm sm:p-5"
    data-testid="deals-pipeline-hero"
  >
    <span
      class="pointer-events-none absolute -right-12 -top-12 h-44 w-44 rounded-full opacity-25 blur-3xl"
      :style="{ background: 'var(--ds-brand-gradient)' }"
      aria-hidden="true"
    />

    <div class="relative grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
      <article
        v-for="(stat, idx) in stats"
        :key="stat.key"
        class="ds-deals-pipeline-hero__tile group relative overflow-hidden rounded-xl border border-ds-border-muted bg-ds-bg-elevated/70 p-3.5 transition-[border-color,box-shadow,transform] duration-200 hover:-translate-y-0.5 hover:border-ds-brand-muted/60 hover:shadow-ds-md"
        :class="stat.featured ? 'ds-deals-pipeline-hero__tile--featured' : ''"
      >
        <span
          class="pointer-events-none absolute inset-x-0 top-0 h-0.5"
          :class="stat.accentClass"
          aria-hidden="true"
        />
        <span
          v-if="stat.featured"
          class="pointer-events-none absolute -right-8 -top-8 h-24 w-24 rounded-full opacity-25 blur-2xl"
          :style="{ background: 'var(--ds-brand-gradient)' }"
          aria-hidden="true"
        />

        <div class="relative flex items-start gap-3">
          <span
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl transition-colors duration-200"
            :class="stat.iconWrapClass"
            :style="stat.iconStyle"
          >
            <UIcon :name="stat.icon" class="h-4 w-4" aria-hidden="true" />
          </span>
          <div class="min-w-0 flex-1">
            <div class="flex items-center justify-between gap-2">
              <p class="truncate text-[11px] font-medium uppercase tracking-wider text-ds-fg-muted">
                {{ stat.label }}
              </p>
              <span
                v-if="stat.tagLabel"
                class="inline-flex items-center gap-1 rounded-full px-1.5 py-0.5 text-[10px] font-semibold leading-none ring-1 ring-inset"
                :class="stat.tagClass"
              >
                <UIcon
                  v-if="stat.tagIcon"
                  :name="stat.tagIcon"
                  class="h-2.5 w-2.5"
                  aria-hidden="true"
                />
                {{ stat.tagLabel }}
              </span>
            </div>
            <p class="mt-1 truncate text-[1.5rem] font-extrabold tabular-nums tracking-tight text-ds-fg-heading sm:text-2xl">
              <span
                v-if="stat.featured"
                class="bg-clip-text"
                :style="brandGradientText"
              >
                {{ stat.value }}
              </span>
              <span
                v-else
                class="text-ds-fg-heading"
              >
                {{ stat.value }}
              </span>
            </p>
            <p class="mt-0.5 text-xs text-ds-fg-subtle">{{ stat.hint }}</p>
          </div>
        </div>

        <button
          v-if="stat.scrollTarget"
          type="button"
          class="absolute inset-0 cursor-pointer rounded-xl"
          :aria-label="stat.label"
          :tabindex="idx === 0 ? 0 : -1"
          @click="scrollToFocus"
        />
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { DealPipelineStage, DealPipelineSummary } from '~/types/deal'

interface HeroStat {
  key: string
  label: string
  value: string | number
  hint: string
  icon: string
  iconWrapClass: string
  iconStyle?: Record<string, string>
  accentClass: string
  featured?: boolean
  tagLabel?: string
  tagIcon?: string
  tagClass?: string
  scrollTarget?: boolean
}

const props = defineProps<{
  summary: DealPipelineSummary
  stages: DealPipelineStage[]
}>()

const { t } = useI18n()
const { formatDealAmount } = useDealLabels()

const brandGradientText = {
  background: 'var(--ds-brand-gradient)',
  '-webkit-background-clip': 'text',
  'background-clip': 'text',
  color: 'transparent',
}

function startOfDay(d: Date): number {
  const x = new Date(d)
  x.setHours(0, 0, 0, 0)
  return x.getTime()
}

const closingThisWeekCount = computed(() => {
  const today = startOfDay(new Date())
  let n = 0
  for (const col of props.stages) {
    if (col.stage === 'won' || col.stage === 'lost') continue
    for (const item of col.items) {
      if (!item.expected_close_date) continue
      const close = startOfDay(new Date(`${item.expected_close_date}T00:00:00`))
      const days = Math.round((close - today) / 86_400_000)
      if (days >= 0 && days <= 7) n++
    }
  }
  return n
})

const overdueCount = computed(() => {
  const today = startOfDay(new Date())
  let n = 0
  for (const col of props.stages) {
    if (col.stage === 'won' || col.stage === 'lost') continue
    for (const item of col.items) {
      if (!item.expected_close_date) continue
      const close = startOfDay(new Date(`${item.expected_close_date}T00:00:00`))
      if (close < today) n++
    }
  }
  return n
})

function scrollToFocus() {
  if (import.meta.server) return
  document.querySelector('[data-testid="deals-focus-stream"]')?.scrollIntoView({
    behavior: 'smooth',
    block: 'nearest',
  })
}

const stats = computed<HeroStat[]>(() => {
  const s = props.summary
  return [
    {
      key: 'open-amount',
      label: t('dealsHeroStatOpenAmount'),
      value: formatDealAmount(s.open_amount),
      hint: t('dealsHeroStatOpenAmountHint', { n: s.open_count }),
      icon: 'i-heroicons-banknotes',
      iconWrapClass: 'text-ds-on-brand shadow-ds-brand',
      iconStyle: { background: 'var(--ds-brand-gradient)' },
      accentClass: 'bg-ds-brand opacity-70',
      featured: true,
      tagLabel: t('dealsHeroStatLive'),
      tagIcon: 'i-heroicons-bolt',
      tagClass: 'bg-ds-brand-subtle text-ds-fg-brand ring-ds-brand/25',
    },
    {
      key: 'closing',
      label: t('dealsHeroStatClosing'),
      value: closingThisWeekCount.value,
      hint: t('dealsHeroStatClosingHint'),
      icon: 'i-heroicons-calendar-days',
      iconWrapClass:
        closingThisWeekCount.value > 0
          ? 'bg-ds-success-subtle text-ds-success'
          : 'bg-ds-bg-muted text-ds-fg-muted',
      accentClass:
        closingThisWeekCount.value > 0 ? 'bg-ds-success opacity-70' : 'bg-ds-border opacity-50',
      tagLabel:
        closingThisWeekCount.value > 0
          ? t('dealsHeroStatClosingTag')
          : t('dealsHeroStatQuietTag'),
      tagIcon:
        closingThisWeekCount.value > 0
          ? 'i-heroicons-arrow-trending-up'
          : 'i-heroicons-minus',
      tagClass:
        closingThisWeekCount.value > 0
          ? 'bg-ds-success-subtle text-ds-success ring-ds-success/25'
          : 'bg-ds-bg-muted text-ds-fg-muted ring-ds-border-muted',
      scrollTarget: closingThisWeekCount.value > 0,
    },
    {
      key: 'overdue',
      label: t('dealsHeroStatOverdue'),
      value: overdueCount.value,
      hint: t('dealsHeroStatOverdueHint'),
      icon: 'i-heroicons-exclamation-triangle',
      iconWrapClass:
        overdueCount.value > 0
          ? 'bg-ds-danger-subtle text-ds-danger'
          : 'bg-ds-bg-muted text-ds-fg-muted',
      accentClass:
        overdueCount.value > 0 ? 'bg-ds-danger opacity-70' : 'bg-ds-border opacity-50',
      tagLabel:
        overdueCount.value > 0 ? t('dealsHeroStatActNow') : t('dealsHeroStatOnTrack'),
      tagIcon:
        overdueCount.value > 0
          ? 'i-heroicons-exclamation-triangle'
          : 'i-heroicons-check-circle',
      tagClass:
        overdueCount.value > 0
          ? 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/25'
          : 'bg-ds-success-subtle text-ds-success ring-ds-success/25',
      scrollTarget: overdueCount.value > 0,
    },
    {
      key: 'won-mtd',
      label: t('dealsHeroStatWonMtd'),
      value: s.won_count_mtd,
      hint: t('dealsHeroStatWonMtdHint', { amount: formatDealAmount(s.won_amount_mtd) }),
      icon: 'i-heroicons-trophy',
      iconWrapClass:
        s.won_count_mtd > 0
          ? 'bg-ds-warning-subtle text-ds-warning'
          : 'bg-ds-bg-muted text-ds-fg-muted',
      accentClass:
        s.won_count_mtd > 0 ? 'bg-ds-warning opacity-70' : 'bg-ds-border opacity-50',
      tagLabel:
        s.won_count_mtd > 0 ? t('dealsHeroStatMomentum') : t('dealsHeroStatQuietTag'),
      tagIcon: s.won_count_mtd > 0 ? 'i-heroicons-fire' : 'i-heroicons-minus',
      tagClass:
        s.won_count_mtd > 0
          ? 'bg-ds-warning-subtle text-ds-warning ring-ds-warning/25'
          : 'bg-ds-bg-muted text-ds-fg-muted ring-ds-border-muted',
    },
  ]
})
</script>

<style scoped>
.ds-deals-pipeline-hero__tile--featured {
  background-image: linear-gradient(
    180deg,
    color-mix(in srgb, var(--ds-bg-elevated) 92%, var(--ds-brand) 8%) 0%,
    var(--ds-bg-elevated) 70%
  );
}
</style>
