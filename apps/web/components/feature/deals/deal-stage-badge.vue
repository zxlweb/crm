<template>
  <span
    class="inline-flex items-center text-xs font-medium transition-colors duration-200"
    :class="[
      variant === 'plain' ? 'gap-1.5 text-ds-fg-muted' : 'gap-1.5 rounded-full border px-2.5 py-0.5',
      variant === 'filled' ? filledClass : '',
    ]"
  >
    <span class="h-1.5 w-1.5 shrink-0 rounded-full" :class="dotClass" aria-hidden="true" />
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import type { DealStage } from '~/types/deal'

const props = withDefaults(
  defineProps<{
    stage: DealStage
    variant?: 'filled' | 'plain'
  }>(),
  { variant: 'filled' },
)

const { dealStageLabel } = useDealLabels()

const label = computed(() => dealStageLabel(props.stage))

const filledClass = computed(() => {
  const map: Record<DealStage, string> = {
    qualification: 'border-ds-info/25 bg-ds-info-subtle text-ds-info',
    proposal: 'border-ds-brand-muted bg-ds-brand-subtle text-ds-fg-brand',
    negotiation: 'border-ds-warning/25 bg-ds-warning-subtle text-ds-warning',
    won: 'border-ds-success/25 bg-ds-success-subtle text-ds-success',
    lost: 'border-ds-border bg-ds-bg-muted text-ds-fg-muted',
  }
  return map[props.stage]
})

const dotClass = computed(() => {
  const map: Record<DealStage, string> = {
    qualification: 'bg-ds-info',
    proposal: 'bg-ds-brand',
    negotiation: 'bg-ds-warning',
    won: 'bg-ds-success',
    lost: 'bg-ds-fg-subtle',
  }
  return map[props.stage]
})
</script>
