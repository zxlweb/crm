<template>
  <span
    class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors duration-200"
    :class="toneClass"
  >
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import type { RelationshipHealth } from '~/types/lead'

const props = defineProps<{
  health: RelationshipHealth
}>()

const { t } = useI18n()

const label = computed(() => t(`relationshipHealth.${props.health}`))

const toneClass = computed(() => {
  const map: Record<RelationshipHealth, string> = {
    high: 'bg-ds-success-subtle text-ds-success',
    medium: 'bg-ds-warning-subtle text-ds-warning',
    low: 'bg-ds-danger-subtle text-ds-danger',
  }
  return map[props.health]
})
</script>
