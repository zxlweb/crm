<template>
  <article class="ds-card ds-card-metric rounded-ds-xl p-ds-5 shadow-ds-sm">
    <div class="flex items-start gap-ds-3">
      <div class="ds-card-metric__icon" :data-tone="iconTone" aria-hidden="true">
        <slot name="icon" />
      </div>
      <div class="min-w-0 flex-1">
        <p class="text-ds-sm font-ds-medium text-ds-fg-muted">
          {{ label }}
        </p>
        <p class="mt-ds-1 text-ds-3xl font-ds-bold leading-none tracking-tight text-ds-fg-heading tabular-nums">
          {{ formattedValue }}
        </p>
      </div>
    </div>

    <div
      v-if="compareLabel || trend != null"
      class="flex items-center justify-between gap-ds-2"
    >
      <span v-if="compareLabel" class="text-ds-xs text-ds-fg-muted">{{ compareLabel }}</span>
      <span
        v-if="trend != null"
        class="ml-auto text-ds-xs font-ds-semibold tabular-nums"
        :class="trendClass"
      >
        {{ trendDisplay }}
        <span class="ml-0.5" aria-hidden="true">{{ trendArrow }}</span>
      </span>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { CardMetricIconTone, CardMetricTrendDirection } from '../../../card/types'

export type { CardMetricIconTone, CardMetricTrendDirection } from '../../../card/types'

const props = withDefaults(
  defineProps<{
    label: string
    value: string | number
    compareLabel?: string
    trend?: string
    trendDirection?: CardMetricTrendDirection
    iconTone?: CardMetricIconTone
  }>(),
  {
    iconTone: 'brand',
    trendDirection: 'flat',
  },
)

const formattedValue = computed(() => {
  if (typeof props.value === 'number') {
    return props.value.toLocaleString()
  }
  return props.value
})

const trendClass = computed(() => {
  if (props.trendDirection === 'up') return 'ds-card-metric__trend--up'
  if (props.trendDirection === 'down') return 'ds-card-metric__trend--down'
  return 'ds-card-metric__trend--flat'
})

const trendArrow = computed(() => {
  if (props.trendDirection === 'up') return '↑'
  if (props.trendDirection === 'down') return '↓'
  return ''
})

const trendDisplay = computed(() => props.trend ?? '')
</script>
