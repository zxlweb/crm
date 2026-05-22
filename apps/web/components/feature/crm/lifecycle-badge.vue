<template>
  <span
    class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors duration-200"
    :class="toneClass"
  >
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import type { LifecycleStage } from '~/types/lead'

const props = defineProps<{
  stage: LifecycleStage
}>()

const { t } = useI18n()

const label = computed(() => t(`lifecycle.${props.stage}`))

const toneClass = computed(() => {
  const map: Record<LifecycleStage, string> = {
    acquire: 'bg-ds-bg-muted text-ds-fg-muted',
    activate: 'bg-ds-info-subtle text-ds-info',
    grow: 'bg-ds-brand-subtle text-ds-fg-brand',
    retain: 'bg-ds-success-subtle text-ds-success',
    revive: 'bg-ds-warning-subtle text-ds-warning',
  }
  return map[props.stage]
})
</script>
