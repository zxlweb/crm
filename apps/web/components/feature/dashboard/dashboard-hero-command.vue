<template>
  <header
    class="ds-dash-hero relative overflow-hidden rounded-3xl border border-ds-border bg-ds-bg-elevated/85 px-5 py-5 shadow-ds-md backdrop-blur-md sm:px-7 sm:py-7"
    data-testid="dashboard-hero"
  >
    <div
      class="pointer-events-none absolute -right-12 -top-12 h-56 w-56 rounded-full opacity-30 blur-3xl"
      :style="{ background: 'var(--ds-brand-gradient)' }"
      aria-hidden="true"
    />
    <div
      class="pointer-events-none absolute -bottom-16 left-1/3 h-40 w-40 rounded-full opacity-20 blur-3xl"
      :style="{ background: 'var(--ds-blur-brand)' }"
      aria-hidden="true"
    />

    <div class="relative">
      <div class="min-w-0">
        <div
          class="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-[0.18em] text-ds-fg-brand"
        >
          <span class="relative flex h-2 w-2" aria-hidden="true">
            <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-ds-brand opacity-50" />
            <span class="relative inline-flex h-2 w-2 rounded-full bg-ds-brand" />
          </span>
          <span>{{ greeting }}</span>
          <span aria-hidden="true" class="text-ds-fg-subtle">·</span>
          <span class="text-ds-fg-muted">
            {{ isPreviewMode ? $t('dashboardHeroDemoBadge') : $t('dashboardHeroLiveBadge') }}
          </span>
        </div>

        <h1
          id="dashboard-hero-heading"
          class="mt-1.5 text-2xl font-bold tracking-tight text-ds-fg-heading sm:text-3xl"
        >
          {{ headline }}
        </h1>

        <div class="mt-4 flex flex-wrap items-end gap-x-4 gap-y-2">
          <div>
            <p
              class="text-[11px] font-medium uppercase tracking-wider text-ds-fg-muted"
            >
              {{ $t('dashboardHeroTotalLabel') }}
            </p>
            <p
              class="mt-0.5 text-3xl font-extrabold tabular-nums tracking-tight text-ds-fg-heading sm:text-4xl"
            >
              <span class="bg-clip-text" :style="brandGradientText">
                {{ pipelineValueLabel }}
              </span>
            </p>
          </div>
          <div class="flex flex-wrap items-center gap-1.5">
            <NuxtLink
              to="/deals"
              class="ds-dash-hero__chip inline-flex cursor-pointer items-center gap-1.5 rounded-full border border-ds-info/25 bg-ds-info-subtle px-2.5 py-1 text-xs font-medium text-ds-info transition-colors duration-200 hover:bg-ds-info-subtle/70"
              data-testid="dashboard-hero-chip-open"
            >
              <UIcon name="i-heroicons-briefcase" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ $t('dashboardHeroOpenDeals', { n: snapshot.dealsOpenCount }) }}
            </NuxtLink>
            <NuxtLink
              v-if="snapshot.priorityCount > 0"
              to="#dashboard-priority-section"
              class="ds-dash-hero__chip inline-flex cursor-pointer items-center gap-1.5 rounded-full border border-ds-warning/25 bg-ds-warning-subtle px-2.5 py-1 text-xs font-medium text-ds-warning transition-colors duration-200 hover:bg-ds-warning-subtle/70"
              data-testid="dashboard-hero-chip-priority"
            >
              <UIcon name="i-heroicons-bell-alert" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ $t('dashboardHeroPriority', { n: snapshot.priorityCount }) }}
            </NuxtLink>
            <NuxtLink
              v-if="snapshot.atRiskTotal > 0"
              to="/leads?health=low"
              class="ds-dash-hero__chip inline-flex cursor-pointer items-center gap-1.5 rounded-full border border-ds-danger/25 bg-ds-danger-subtle px-2.5 py-1 text-xs font-medium text-ds-danger transition-colors duration-200 hover:bg-ds-danger-subtle/70"
              data-testid="dashboard-hero-chip-at-risk"
            >
              <UIcon name="i-heroicons-exclamation-triangle" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ $t('dashboardHeroAtRisk', { n: snapshot.atRiskTotal }) }}
            </NuxtLink>
            <span
              v-if="weeklyFollowUpNote"
              class="inline-flex items-center gap-1.5 rounded-full border border-ds-success/25 bg-ds-success-subtle px-2.5 py-1 text-xs font-medium text-ds-success"
              data-testid="dashboard-hero-chip-follow"
            >
              <UIcon name="i-heroicons-calendar-days" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ weeklyFollowUpNote }}
            </span>
          </div>
        </div>

        <p class="mt-3 max-w-xl text-sm text-ds-fg-muted">
          {{ subtitleText }}
        </p>
        <p class="mt-2 text-xs text-ds-fg-subtle">
          <span>{{ $t('dashboardHeroLastSync') }}</span>
          <span class="ml-1 font-medium tabular-nums text-ds-fg-muted">{{ lastSyncLabel }}</span>
        </p>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import type { DashboardSnapshot } from '~/types/dashboard'

const props = withDefaults(
  defineProps<{
    snapshot: DashboardSnapshot
    greeting: string
    headline: string
    weeklyFollowUpNote?: string
    isPreviewMode?: boolean
  }>(),
  {
    weeklyFollowUpNote: '',
    isPreviewMode: false,
  },
)

const { t, locale } = useI18n()
const { formatDealAmount } = useDealLabels()

const brandGradientText = {
  background: 'var(--ds-brand-gradient)',
  '-webkit-background-clip': 'text',
  'background-clip': 'text',
  color: 'transparent',
}

const pipelineValueLabel = computed(() => formatDealAmount(props.snapshot.dealsOpenAmount))

const subtitleText = computed(() => {
  const count = props.snapshot.priorityCount
  if (props.isPreviewMode) return t('dashboardHeroSubtitleDemo')
  if (count > 0) return t('dashboardHeroSubtitleBusy', { count })
  return t('dashboardHeroSubtitleQuiet')
})

const lastSyncLabel = computed(() =>
  new Intl.DateTimeFormat(locale.value, {
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date()),
)
</script>
