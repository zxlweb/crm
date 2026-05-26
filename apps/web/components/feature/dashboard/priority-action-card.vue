<template>
  <NuxtLink
    :to="item.followHref"
    class="group relative flex cursor-pointer items-center gap-3 overflow-hidden rounded-xl border bg-ds-bg-elevated/80 px-3 py-2.5 pl-4 transition-[border-color,background-color,box-shadow] duration-200 hover:border-ds-brand-muted hover:bg-ds-bg-muted/50 hover:shadow-ds-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ds-brand"
    :class="surfaceClass"
    data-testid="priority-action-card"
  >
    <span
      class="absolute inset-y-2 left-0 w-0.5 rounded-full"
      :class="borderAccentClass"
      aria-hidden="true"
    />

    <div
      class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg"
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
        <span class="inline-flex rounded-md px-1.5 py-0.5 text-[10px] font-semibold" :class="typeBadgeClass">
          {{ item.entityType === 'lead' ? $t('dashboardTypeLead') : $t('dashboardTypeAccount') }}
        </span>
        <span class="inline-flex rounded-md px-1.5 py-0.5 text-[10px] font-semibold" :class="healthBadgeClass">
          {{ healthLabelText }}
        </span>
        <AiPreviewBadge v-if="item.isPreview" />
      </div>

      <p class="mt-0.5 line-clamp-1 text-xs text-ds-fg-muted">
        {{ item.reasons.join(' · ') }}
      </p>

      <div class="mt-1 flex flex-wrap items-center gap-x-2.5 gap-y-0.5 text-[11px]">
        <span class="line-clamp-1 font-medium" :class="insightStripClass">
          {{ $t('dashboardSuggestionPrefix') }}{{ item.suggestion }}
        </span>
        <span v-if="item.contactName" class="inline-flex shrink-0 items-center gap-0.5 text-ds-fg-muted">
          <UIcon name="i-heroicons-user-20-solid" class="h-3 w-3" aria-hidden="true" />
          {{ item.contactName }}
        </span>
        <span
          v-if="item.daysSinceActivity != null"
          class="inline-flex shrink-0 items-center gap-0.5 text-ds-fg-muted"
        >
          <UIcon name="i-heroicons-clock-20-solid" class="h-3 w-3" aria-hidden="true" />
          {{ $t('dashboardDaysAgo', { days: item.daysSinceActivity }) }}
        </span>
      </div>
    </div>

    <div class="hidden shrink-0 items-center gap-2 sm:flex">
      <div
        class="flex h-6 w-14 items-center justify-center overflow-hidden rounded-md bg-ds-bg-muted/40"
        aria-hidden="true"
      >
        <DashboardSparkline
          :values="item.sparkline"
          :tone="item.healthLabel"
          :width="52"
          :height="22"
        />
      </div>
      <span class="min-w-[1.75rem] text-right text-sm font-bold tabular-nums" :class="scoreTextClass">
        {{ item.engagementScore }}
      </span>
    </div>

    <span
      class="inline-flex shrink-0 items-center gap-0.5 text-xs font-semibold text-ds-fg-brand transition-transform duration-200 group-hover:translate-x-0.5"
    >
      <span class="hidden sm:inline">{{ $t('dashboardFollowCta') }}</span>
      <span class="tabular-nums sm:hidden" :class="scoreTextClass">{{ item.engagementScore }}</span>
      <UIcon
        name="i-heroicons-chevron-right-20-solid"
        class="h-4 w-4"
        aria-hidden="true"
      />
    </span>
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
      return 'bg-ds-danger'
    case 'watch':
      return 'bg-ds-warning'
    default:
      return 'bg-ds-success'
  }
})

const surfaceClass = computed(() => {
  switch (props.item.healthLabel) {
    case 'alert':
      return 'border-ds-danger/25'
    case 'watch':
      return 'border-ds-warning/25'
    default:
      return 'border-ds-border-muted'
  }
})

const avatarClass = computed(() => {
  switch (props.item.healthLabel) {
    case 'alert':
      return 'bg-ds-danger-subtle text-ds-danger'
    case 'watch':
      return 'bg-ds-warning-subtle text-ds-warning'
    default:
      return 'bg-ds-success-subtle text-ds-success'
  }
})

const typeBadgeClass = computed(() =>
  props.item.entityType === 'lead'
    ? 'bg-ds-warning-subtle text-ds-warning ring-1 ring-inset ring-ds-warning/20'
    : 'bg-ds-bg-muted text-ds-fg-muted ring-1 ring-inset ring-ds-border-muted',
)

const healthBadgeClass = computed(() => {
  switch (props.item.healthLabel) {
    case 'alert':
      return 'bg-ds-danger-subtle text-ds-danger ring-1 ring-inset ring-ds-danger/20'
    case 'watch':
      return 'bg-ds-warning-subtle text-ds-warning ring-1 ring-inset ring-ds-warning/20'
    default:
      return 'bg-ds-success-subtle text-ds-success ring-1 ring-inset ring-ds-success/20'
  }
})

const insightStripClass = computed(() => {
  switch (props.item.healthLabel) {
    case 'alert':
      return 'text-ds-danger'
    case 'watch':
      return 'text-ds-warning'
    default:
      return 'text-ds-success'
  }
})

const scoreTextClass = computed(() => {
  const band = scoreBand(props.item.engagementScore)
  if (band === 'excellent') return 'text-ds-success'
  if (band === 'watch') return 'text-ds-warning'
  return 'text-ds-danger'
})
</script>
