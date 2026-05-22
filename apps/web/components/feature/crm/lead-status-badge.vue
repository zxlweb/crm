<template>
  <span
    class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors duration-200"
    :class="toneClass"
  >
    <span class="h-1.5 w-1.5 rounded-full" :class="dotClass" />
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import type { LeadStatus } from '~/types/lead'

const props = defineProps<{
  status: LeadStatus
}>()

const { t } = useI18n()

const label = computed(() => t(`leadStatus.${props.status}`))

const toneClass = computed(() => {
  const map: Record<LeadStatus, string> = {
    new: 'bg-ds-bg-muted text-ds-fg-muted',
    contacted: 'bg-ds-info-subtle text-ds-info',
    qualified: 'bg-ds-brand-subtle text-ds-fg-brand',
    unqualified: 'bg-ds-danger-subtle text-ds-danger',
    converted: 'bg-ds-success-subtle text-ds-success',
  }
  return map[props.status]
})

const dotClass = computed(() => {
  const map: Record<LeadStatus, string> = {
    new: 'bg-ds-fg-subtle',
    contacted: 'bg-ds-info',
    qualified: 'bg-ds-brand',
    unqualified: 'bg-ds-danger',
    converted: 'bg-ds-success',
  }
  return map[props.status]
})
</script>
