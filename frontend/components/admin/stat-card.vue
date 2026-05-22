<template>
  <div class="ds-card rounded-2xl p-5 shadow-sm">
    <div class="flex items-start justify-between gap-3">
      <div>
        <p class="text-xs font-medium text-ds-fg-muted">{{ label }}</p>
        <p class="mt-2 text-2xl font-bold tracking-tight text-ds-fg-heading">{{ value }}</p>
        <p
          v-if="trend"
          class="mt-1 flex items-center gap-1 text-xs font-medium"
          :class="trendPositive ? 'text-ds-success' : 'text-ds-fg-muted'"
        >
          <svg v-if="trendPositive" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 10l7-7m0 0l7 7m-7-7v18" />
          </svg>
          {{ trend }}
        </p>
      </div>
      <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-ds-brand-subtle">
        <slot name="icon" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  value: string | number
  trend?: string
  trendPositive?: boolean
}>()
</script>
