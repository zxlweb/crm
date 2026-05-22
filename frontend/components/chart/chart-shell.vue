<template>
  <div class="ds-card overflow-hidden rounded-2xl shadow-sm">
    <div v-if="title || metric != null" class="border-b border-ds-border-muted px-5 pb-4 pt-5">
      <div class="flex flex-wrap items-start justify-between gap-3">
        <div>
          <p v-if="title" class="text-sm font-medium text-ds-fg-muted">{{ title }}</p>
          <div v-if="metric != null" class="mt-1 flex flex-wrap items-baseline gap-2">
            <span class="text-2xl font-bold tracking-tight text-ds-fg-heading sm:text-3xl">{{ metric }}</span>
            <span v-if="metricLabel" class="text-xs text-ds-fg-muted">{{ metricLabel }}</span>
          </div>
          <p v-if="subtitle" class="mt-1 text-xs text-ds-fg-muted">{{ subtitle }}</p>
        </div>
        <slot name="header-extra" />
      </div>
      <div v-if="legend?.length" class="mt-3 flex flex-wrap gap-4 text-xs">
        <span
          v-for="item in legend"
          :key="item.label"
          class="flex items-center gap-1.5"
          :class="item.muted ? 'text-ds-fg-muted' : 'text-ds-fg-brand'"
        >
          <span
            class="h-2 w-2 rounded-full"
            :class="item.dashed ? 'border border-ds-fg-muted bg-transparent' : 'bg-ds-brand'"
            :style="item.color ? { background: item.color } : undefined"
          />
          {{ item.label }}
        </span>
      </div>
    </div>
    <div class="px-2 pb-2 pt-1" :style="{ minHeight: `${height}px` }">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ChartLegendItem } from '~/types/chart'

withDefaults(
  defineProps<{
    title?: string
    subtitle?: string
    metric?: string | number | null
    metricLabel?: string
    legend?: ChartLegendItem[]
    height?: number
  }>(),
  { height: 280 },
)
</script>
