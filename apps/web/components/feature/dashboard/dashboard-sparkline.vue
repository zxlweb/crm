<template>
  <svg
    :width="width"
    :height="height"
    :viewBox="`0 0 ${width} ${height}`"
    class="block shrink-0 overflow-hidden"
    aria-hidden="true"
  >
    <polyline
      :points="polylinePoints"
      fill="none"
      :stroke="strokeColor"
      stroke-width="1.75"
      stroke-linecap="round"
      stroke-linejoin="round"
      vector-effect="non-scaling-stroke"
    />
  </svg>
</template>

<script setup lang="ts">
import type { PriorityHealthLabel } from '~/types/dashboard'
import { scoreBand } from '~/utils/dashboard-score-band'

const props = withDefaults(
  defineProps<{
    values: number[]
    width?: number
    height?: number
    /** 与跟进卡健康度色一致；未传则按末点活跃分推断 */
    tone?: PriorityHealthLabel
  }>(),
  {
    width: 56,
    height: 24,
    tone: undefined,
  },
)

const latest = computed(() => props.values[props.values.length - 1] ?? 0)

const strokeColor = computed(() => {
  if (props.tone === 'alert') return 'var(--ds-danger, #dc2626)'
  if (props.tone === 'watch') return 'var(--ds-warning, #f59e0b)'
  if (props.tone === 'healthy') return 'var(--ds-success, #059669)'

  const band = scoreBand(latest.value)
  if (band === 'excellent') return 'var(--ds-success, #059669)'
  if (band === 'watch') return 'var(--ds-warning, #f59e0b)'
  return 'var(--ds-danger, #dc2626)'
})

const polylinePoints = computed(() => {
  const vals = props.values.length ? props.values : [0]
  const max = Math.max(...vals, 1)
  const min = Math.min(...vals, 0)
  const range = max - min || 1
  const padX = 1
  const padY = 2
  const innerW = props.width - padX * 2
  const innerH = props.height - padY * 2

  return vals
    .map((v, i) => {
      const x =
        vals.length === 1
          ? props.width / 2
          : padX + (i / (vals.length - 1)) * innerW
      const y = padY + innerH - ((v - min) / range) * innerH
      return `${x},${y}`
    })
    .join(' ')
})
</script>
