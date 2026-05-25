<template>
  <div class="relative" :style="{ height: `${height}px` }">
    <div
      v-if="pending"
      class="flex h-full items-center justify-center"
    >
      <UIcon name="i-heroicons-arrow-path" class="h-7 w-7 animate-spin text-primary" />
    </div>
    <p
      v-else-if="empty"
      class="flex h-full items-center justify-center text-sm text-ds-fg-muted"
    >
      {{ emptyText || $t('leadsReportEmpty') }}
    </p>
    <slot v-else />
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    pending?: boolean
    empty?: boolean
    height?: number
    emptyText?: string
  }>(),
  { pending: false, empty: false, height: 240, emptyText: '' },
)
</script>
