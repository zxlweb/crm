<template>
  <svg
    :width="width"
    :height="height"
    :viewBox="`0 0 ${width} ${height}`"
    class="shrink-0 overflow-visible"
    aria-hidden="true"
  >
    <polyline
      :points="polylinePoints"
      fill="none"
      :stroke="strokeColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
</template>

<script setup lang="ts">
import { scoreBand } from '~/utils/dashboard-score-band'

const props = withDefaults(
  defineProps<{
    values: number[]
    width?: number
    height?: number
  }>(),
  {
    width: 56,
    height: 24,
  },
)

const latest = computed(() => props.values[props.values.length - 1] ?? 0)

const strokeColor = computed(() => {
  const band = scoreBand(latest.value)
  if (band === 'excellent') return 'var(--ds-success, #059669)'
  if (band === 'watch') return '#f59e0b'
  return 'var(--ds-warning, #dc2626)'
})

const polylinePoints = computed(() => {
  const vals = props.values.length ? props.values : [0]
  const max = Math.max(...vals, 1)
  const min = Math.min(...vals, 0)
  const range = max - min || 1
  const padY = 2
  const innerH = props.height - padY * 2

  return vals
    .map((v, i) => {
      const x = vals.length === 1 ? props.width / 2 : (i / (vals.length - 1)) * props.width
      const y = padY + innerH - ((v - min) / range) * innerH
      return `${x},${y}`
    })
    .join(' ')
})
</script>
