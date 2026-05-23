<template>
  <section
    class="overflow-hidden rounded-2xl border border-ds-border/60 bg-ds-bg-elevated shadow-ds-sm"
    data-testid="dashboard-team-heatmap"
  >
    <div class="flex flex-wrap items-center justify-between gap-3 border-b border-ds-border-muted px-5 py-4 sm:px-6">
      <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('dashboardHeatmapTitle') }}</h3>
      <NuxtLink
        to="/leads?tab=reports"
        class="cursor-pointer text-xs font-semibold text-ds-fg-brand transition-colors hover:text-ds-brand-hover"
      >
        {{ $t('dashboardHeatmapViewDetail') }}
      </NuxtLink>
    </div>

    <div class="px-5 py-3 sm:px-6">
      <div
        class="grid grid-cols-[minmax(0,1fr)_4.5rem_4.5rem] gap-x-3 border-b border-ds-border-muted pb-2 text-[11px] font-medium text-ds-fg-muted"
      >
        <span>{{ $t('dashboardHeatmapColMember') }}</span>
        <span class="text-center">{{ $t('dashboardHeatmapColHealth') }}</span>
        <span class="text-center">{{ $t('dashboardHeatmapColEmotion') }}</span>
      </div>

      <ul class="divide-y divide-ds-border-muted">
        <li
          v-for="row in members"
          :key="row.memberId"
          class="grid grid-cols-[minmax(0,1fr)_4.5rem_4.5rem] items-center gap-x-3 py-3"
        >
          <div class="flex min-w-0 items-center gap-2.5">
            <div
              class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-ds-bg-muted text-xs font-bold text-ds-fg-muted"
            >
              {{ row.memberName.slice(0, 1) }}
            </div>
            <span class="truncate text-sm font-medium text-ds-fg-heading">{{ row.memberName }}</span>
          </div>
          <div class="flex justify-center">
            <span
              class="inline-flex min-w-[2.75rem] justify-center rounded-md px-2 py-1 text-xs font-bold tabular-nums"
              :class="scoreBandPillClass(scoreBand(row.healthScore))"
            >
              {{ row.healthScore }}
            </span>
          </div>
          <div class="flex justify-center">
            <span
              class="inline-flex min-w-[2.75rem] justify-center rounded-md px-2 py-1 text-xs font-bold tabular-nums"
              :class="scoreBandPillClass(scoreBand(row.emotionScore))"
            >
              {{ row.emotionScore }}
            </span>
          </div>
        </li>
      </ul>

      <div class="mt-4 flex flex-wrap items-center gap-x-4 gap-y-2 text-[11px] text-ds-fg-muted">
        <span class="inline-flex items-center gap-1.5">
          <span class="h-2 w-2 rounded-sm" :class="scoreBandLegendClass('excellent')" />
          {{ $t('dashboardHeatmapLegendExcellent') }}
        </span>
        <span class="inline-flex items-center gap-1.5">
          <span class="h-2 w-2 rounded-sm" :class="scoreBandLegendClass('watch')" />
          {{ $t('dashboardHeatmapLegendWatch') }}
        </span>
        <span class="inline-flex items-center gap-1.5">
          <span class="h-2 w-2 rounded-sm" :class="scoreBandLegendClass('alert')" />
          {{ $t('dashboardHeatmapLegendAlert') }}
        </span>
      </div>

      <div class="mt-3 flex flex-wrap items-center justify-between gap-2 border-t border-ds-border-muted pt-3">
        <p class="text-xs text-ds-fg-muted">
          {{ $t('dashboardHeatmapAvg', { health: avgHealth, emotion: avgEmotion }) }}
        </p>
        <NuxtLink
          to="/leads?tab=reports"
          class="cursor-pointer text-xs font-semibold text-ds-fg-brand transition-colors hover:text-ds-brand-hover"
        >
          {{ $t('dashboardHeatmapViewDetail') }}
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { DASHBOARD_TEAM_HEATMAP_MEMBERS } from '~/fixtures/dashboard-preview'
import { scoreBand, scoreBandLegendClass, scoreBandPillClass } from '~/utils/dashboard-score-band'

const members = DASHBOARD_TEAM_HEATMAP_MEMBERS

const avgHealth = computed(() => {
  if (!members.length) return 0
  return Math.round(members.reduce((s, m) => s + m.healthScore, 0) / members.length)
})

const avgEmotion = computed(() => {
  if (!members.length) return 0
  return Math.round(members.reduce((s, m) => s + m.emotionScore, 0) / members.length)
})
</script>
