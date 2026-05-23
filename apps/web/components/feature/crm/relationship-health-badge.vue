<template>
  <span
    class="inline-flex items-center text-xs font-medium transition-colors duration-200"
    :class="[variant === 'plain' ? plainClass : filledClass, variant === 'filled' ? 'rounded-full px-2.5 py-0.5' : 'gap-1.5']"
  >
    <span
      v-if="variant === 'plain'"
      class="h-1.5 w-1.5 shrink-0 rounded-full"
      :class="dotClass"
      aria-hidden="true"
    />
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import type { RelationshipHealth } from '~/types/lead'

const props = withDefaults(
  defineProps<{
    health: RelationshipHealth
    variant?: 'filled' | 'plain'
    mutedDot?: boolean
  }>(),
  { variant: 'filled', mutedDot: false },
)

const { t } = useI18n()

const label = computed(() => t(`relationshipHealth.${props.health}`))

const filledClass = computed(() => {
  const map: Record<RelationshipHealth, string> = {
    high: 'bg-ds-success-subtle text-ds-success',
    medium: 'bg-ds-warning-subtle text-ds-warning',
    low: 'bg-ds-danger-subtle text-ds-danger',
  }
  return map[props.health]
})

const plainClass = 'text-ds-fg-muted'

const dotClass = computed(() => {
  if (props.variant === 'plain' && props.mutedDot) {
    if (props.health === 'low') return 'bg-ds-danger'
    return 'bg-ds-fg-subtle'
  }
  const map: Record<RelationshipHealth, string> = {
    high: 'bg-ds-success',
    medium: 'bg-ds-warning',
    low: 'bg-ds-danger',
  }
  return map[props.health]
})
</script>
