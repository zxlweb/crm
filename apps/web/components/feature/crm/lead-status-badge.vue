<template>
  <span
    class="inline-flex items-center text-xs font-medium transition-colors duration-200"
    :class="[
      variant === 'plain' ? plainClass : filledClass,
      variant === 'filled' ? 'gap-1.5 rounded-full border px-2.5 py-0.5' : 'gap-1.5',
    ]"
  >
    <span class="h-1.5 w-1.5 shrink-0 rounded-full" :class="dotClass" aria-hidden="true" />
    {{ label }}
    <span
      v-if="status === 'unqualified'"
      class="group/hint relative -mr-0.5 ml-0.5 inline-flex"
      data-testid="lead-status-unqualified-hint"
    >
      <button
        type="button"
        class="inline-flex h-4 w-4 cursor-help items-center justify-center rounded-full text-ds-danger/75 transition-colors duration-200 hover:bg-ds-danger-subtle hover:text-ds-danger focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-ds-brand"
        :aria-label="$t('leadsUnqualifiedHintAria')"
        @click.stop
      >
        <svg class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </button>
      <span
        role="tooltip"
        class="pointer-events-none invisible absolute left-1/2 top-[calc(100%+6px)] z-[var(--ds-z-tooltip)] w-max max-w-[16rem] -translate-x-1/2 rounded-lg border border-ds-border bg-ds-bg-elevated px-3 py-2 text-left text-xs font-normal leading-snug text-ds-fg shadow-ds-md opacity-0 transition-opacity duration-150 group-hover/hint:visible group-hover/hint:opacity-100 group-focus-within/hint:visible group-focus-within/hint:opacity-100"
      >
        {{ $t('leadsOverviewHintUnqualified') }}
      </span>
    </span>
  </span>
</template>

<script setup lang="ts">
import type { LeadStatus } from '~/types/lead'

const props = withDefaults(
  defineProps<{
    status: LeadStatus
    /** filled：详情等彩色 pill；plain：列表等小圆点 + 灰字 */
    variant?: 'filled' | 'plain'
    /** plain 列表场景：圆点默认中性，仅风险态保留语义色 */
    mutedDot?: boolean
  }>(),
  { variant: 'filled', mutedDot: false },
)

const { t } = useI18n()

const label = computed(() => t(`leadStatus.${props.status}`))

const filledClass = computed(() => {
  const map: Record<LeadStatus, string> = {
    new: 'border-ds-border bg-ds-bg-muted text-ds-fg-muted',
    contacted: 'border-ds-info/25 bg-ds-info-subtle text-ds-info',
    qualified: 'border-ds-brand-muted bg-ds-brand-subtle text-ds-fg-brand',
    unqualified: 'border-ds-danger/20 bg-ds-danger-subtle text-ds-danger',
    converted: 'border-ds-success/25 bg-ds-success-subtle text-ds-success',
  }
  return map[props.status]
})

const plainClass = 'text-ds-fg-muted'

const dotClass = computed(() => {
  if (props.variant === 'plain' && props.mutedDot) {
    if (props.status === 'unqualified') return 'bg-ds-danger'
    return 'bg-ds-fg-subtle'
  }
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
