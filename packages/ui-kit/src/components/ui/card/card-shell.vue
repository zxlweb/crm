<template>
  <div class="ds-card overflow-hidden rounded-ds-2xl shadow-ds-sm">
    <div v-if="title || subtitle" class="border-b border-ds-border-muted px-ds-5 pb-ds-4 pt-ds-5">
      <p v-if="title" class="text-ds-sm font-ds-medium text-ds-fg-muted">{{ title }}</p>
      <p v-if="subtitle" class="mt-ds-1 text-ds-xs text-ds-fg-muted">{{ subtitle }}</p>
      <slot name="header-extra" />
    </div>
    <div :class="noPadding ? '' : 'p-ds-5'">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    title?: string
    subtitle?: string
    noPadding?: boolean
  }>(),
  { noPadding: false },
)
</script>
