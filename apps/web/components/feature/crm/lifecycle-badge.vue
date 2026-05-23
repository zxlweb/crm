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
import type { LifecycleStage } from '~/types/lead'

const props = withDefaults(
  defineProps<{
    stage: LifecycleStage
    /** filled：列表等彩色 pill；plain：详情头中性辅信息 */
    variant?: 'filled' | 'plain'
    mutedDot?: boolean
  }>(),
  { variant: 'filled', mutedDot: false },
)

const { t } = useI18n()

const label = computed(() => t(`lifecycle.${props.stage}`))

const filledClass = computed(() => {
  const map: Record<LifecycleStage, string> = {
    acquire: 'bg-ds-bg-muted text-ds-fg-muted',
    activate: 'bg-ds-info-subtle text-ds-info',
    grow: 'bg-ds-brand-subtle text-ds-fg-brand',
    retain: 'bg-ds-success-subtle text-ds-success',
    revive: 'bg-ds-warning-subtle text-ds-warning',
  }
  return map[props.stage]
})

const plainClass = 'text-ds-fg-muted'

const dotClass = computed(() => {
  if (props.variant === 'plain' && props.mutedDot) {
    return 'bg-ds-fg-subtle'
  }
  const map: Record<LifecycleStage, string> = {
    acquire: 'bg-ds-fg-subtle',
    activate: 'bg-ds-info',
    grow: 'bg-ds-brand',
    retain: 'bg-ds-success',
    revive: 'bg-ds-warning',
  }
  return map[props.stage]
})
</script>
