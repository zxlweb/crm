<template>
  <article
    class="ds-card ds-card-metric rounded-ds-xl shadow-ds-sm"
    :class="density === 'compact' ? 'ds-card-metric--compact p-ds-4' : 'p-ds-5'"
  >
    <div class="ds-card-metric__body">
      <div class="flex items-start" :class="density === 'compact' ? 'gap-ds-2' : 'gap-ds-3'">
        <div
          class="ds-card-metric__icon shrink-0"
          :class="density === 'compact' ? 'ds-card-metric__icon--sm' : ''"
          :data-tone="iconTone"
          aria-hidden="true"
        >
          <slot name="icon" />
        </div>
        <div class="min-w-0 flex-1">
          <p
            class="font-ds-medium text-ds-fg-muted"
            :class="density === 'compact' ? 'text-ds-xs leading-tight' : 'text-ds-sm'"
          >
            {{ label }}
          </p>
          <p
            class="mt-0.5 font-ds-bold leading-none tracking-tight text-ds-fg-heading tabular-nums"
            :class="valueClass"
            :title="typeof value === 'string' ? String(value) : undefined"
          >
            {{ formattedValue }}
          </p>
        </div>
      </div>
    </div>

    <div
      v-if="compareLabel || trend != null || $slots['footer-trailing']"
      class="ds-card-metric__footer"
    >
      <span
        v-if="compareLabel"
        class="ds-card-metric__compare min-w-0 flex-1 leading-tight text-ds-fg-muted"
        :class="footerCompareClass"
      >
        {{ compareLabel }}
      </span>
      <div
        v-if="$slots['footer-trailing'] || trend != null"
        class="ds-card-metric__footer-trailing ml-auto flex shrink-0 items-center gap-1.5"
      >
        <slot name="footer-trailing" />
        <span
          v-if="trend != null"
          class="ds-card-metric__trend max-w-[9rem] truncate"
          :class="[trendClass, density === 'compact' ? 'ds-card-metric__trend--compact' : '']"
          :title="trendDisplay"
        >
          {{ trendDisplay }}
          <span v-if="trendArrow" class="ml-0.5" aria-hidden="true">{{ trendArrow }}</span>
        </span>
      </div>
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
    /** full：长文案完整展示；compact：短数字，超出省略 */
    valueDisplay?: 'full' | 'truncate'
    /** compact：单行四列等窄卡；default：Admin 看板 */
    density?: 'default' | 'compact'
  }>(),
  {
    iconTone: 'brand',
    trendDirection: 'flat',
    valueDisplay: 'truncate',
    density: 'default',
  },
)

const valueClass = computed(() => {
  if (props.density === 'compact') {
    return props.valueDisplay === 'full'
      ? 'text-base xl:text-lg'
      : 'text-xl xl:text-2xl'
  }
  return props.valueDisplay === 'full'
    ? 'text-ds-2xl sm:text-ds-3xl'
    : 'truncate text-ds-3xl'
})

const footerCompareClass = computed(() =>
  props.density === 'compact'
    ? 'line-clamp-2 text-[10px] xl:text-ds-xs'
    : 'truncate text-ds-xs',
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
