<template>
  <NuxtLink
    :to="item.followHref"
    class="group relative flex cursor-pointer flex-col gap-2 overflow-hidden rounded-xl border bg-ds-bg-elevated p-3.5 pl-5 transition-colors duration-200 hover:bg-ds-bg-muted/40 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ds-brand sm:flex-row sm:items-center sm:gap-3"
    :class="surfaceClass"
    data-testid="priority-action-card"
  >
    <div
      class="absolute inset-y-0 left-0 w-1"
      :class="borderAccentClass"
      aria-hidden="true"
    />

    <div class="flex min-w-0 flex-1 items-start gap-3">
      <div
        class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full"
        :class="avatarClass"
        aria-hidden="true"
      >
        <UIcon :name="entityIcon" class="h-4 w-4" />
      </div>

      <div class="min-w-0 flex-1">
        <div class="flex flex-wrap items-center gap-1.5">
          <h3 class="truncate text-sm font-semibold text-ds-fg-heading group-hover:text-ds-fg-brand">
            {{ item.title }}
          </h3>
          <span class="inline-flex rounded-md px-1.5 py-0.5 text-[11px] font-semibold" :class="typeBadgeClass">
            {{ item.entityType === 'lead' ? $t('dashboardTypeLead') : $t('dashboardTypeAccount') }}
          </span>
          <span class="inline-flex rounded-md px-1.5 py-0.5 text-[11px] font-semibold" :class="healthBadgeClass">
            {{ healthLabelText }}
          </span>
          <AiPreviewBadge v-if="item.isPreview" />
        </div>

        <p class="mt-0.5 line-clamp-1 text-xs text-ds-fg-muted">
          {{ item.reasons.join(' · ') }}
        </p>

        <p class="mt-1 line-clamp-1 text-xs" :class="insightStripClass">
          {{ $t('dashboardSuggestionPrefix') }}{{ item.suggestion }}
        </p>

        <div class="mt-1.5 flex flex-wrap items-center gap-x-3 gap-y-0.5 text-xs text-ds-fg-muted">
          <span v-if="item.contactName" class="inline-flex items-center gap-1">
            <UIcon name="i-heroicons-user-20-solid" class="h-3.5 w-3.5" />
            {{ item.contactName }}
          </span>
          <span v-if="item.daysSinceActivity != null" class="inline-flex items-center gap-1">
            <UIcon name="i-heroicons-clock-20-solid" class="h-3.5 w-3.5" />
            {{ $t('dashboardDaysAgo', { days: item.daysSinceActivity }) }}
          </span>
        </div>
      </div>
    </div>

    <div class="flex shrink-0 items-center justify-between gap-2 sm:flex-col sm:items-end sm:justify-center sm:gap-1.5">
      <div class="flex items-center gap-2">
        <DashboardSparkline :values="item.sparkline" />
        <span class="min-w-[2rem] text-right text-sm font-bold tabular-nums" :class="scoreTextClass">
          {{ item.engagementScore }}
        </span>
      </div>
      <span
        class="inline-flex items-center gap-1 text-xs font-semibold text-ds-fg-brand transition-colors group-hover:text-ds-brand-hover"
      >
        {{ $t('dashboardFollowCta') }}
        <UIcon
          name="i-heroicons-chevron-right-20-solid"
          class="h-4 w-4 transition-transform group-hover:translate-x-0.5"
          aria-hidden="true"
        />
      </span>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
import type { PriorityActionItem } from '~/types/dashboard'
import { scoreBand } from '~/utils/dashboard-score-band'

const props = defineProps<{
  item: PriorityActionItem
}>()

const { t } = useI18n()

const entityIcon = computed(() =>
  props.item.entityType === 'lead'
    ? 'i-heroicons-light-bulb-20-solid'
    : 'i-heroicons-building-office-2-20-solid',
)

const healthLabelText = computed(() => {
  const map = {
    alert: t('dashboardHealthAlert'),
    watch: t('dashboardHealthWatch'),
    healthy: t('dashboardHealthHealthy'),
  }
  return map[props.item.healthLabel]
})

const borderAccentClass = computed(() => {
  switch (props.item.healthLabel) {
    case 'alert':
      return 'bg-red-500'
    case 'watch':
      return 'bg-amber-400'
    default:
      return 'bg-emerald-500'
  }
})

const surfaceClass = computed(() => {
  switch (props.item.healthLabel) {
    case 'alert':
      return 'border-red-200/80 dark:border-red-900/40'
    case 'watch':
      return 'border-amber-200/80 dark:border-amber-900/40'
    default:
      return 'border-emerald-200/70 dark:border-emerald-900/35'
  }
})

const avatarClass = computed(() => {
  switch (props.item.healthLabel) {
    case 'alert':
      return 'bg-red-50 text-red-600 dark:bg-red-950/40 dark:text-red-400'
    case 'watch':
      return 'bg-amber-50 text-amber-600 dark:bg-amber-950/40 dark:text-amber-400'
    default:
      return 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/40 dark:text-emerald-400'
  }
})

const typeBadgeClass = computed(() =>
  props.item.entityType === 'lead'
    ? 'bg-orange-100 text-orange-800 dark:bg-orange-950/50 dark:text-orange-300'
    : 'bg-stone-100 text-stone-700 dark:bg-stone-800/50 dark:text-stone-300',
)

const healthBadgeClass = computed(() => {
  switch (props.item.healthLabel) {
    case 'alert':
      return 'bg-red-100 text-red-700 dark:bg-red-950/50 dark:text-red-300'
    case 'watch':
      return 'bg-amber-100 text-amber-800 dark:bg-amber-950/50 dark:text-amber-300'
    default:
      return 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/50 dark:text-emerald-300'
  }
})

const insightStripClass = computed(() => {
  switch (props.item.healthLabel) {
    case 'alert':
      return 'text-red-800/90 dark:text-red-200/90'
    case 'watch':
      return 'text-amber-800/90 dark:text-amber-200/90'
    default:
      return 'text-emerald-800/90 dark:text-emerald-200/90'
  }
})

const scoreTextClass = computed(() => {
  const band = scoreBand(props.item.engagementScore)
  if (band === 'excellent') return 'text-emerald-600 dark:text-emerald-400'
  if (band === 'watch') return 'text-amber-600 dark:text-amber-400'
  return 'text-red-600 dark:text-red-400'
})
</script>
