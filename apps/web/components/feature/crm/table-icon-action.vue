<template>
  <NuxtLink
    v-if="to"
    :to="to"
    class="inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-ds-fg-subtle transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg"
    :aria-label="label"
  >
    <UIcon :name="icon" class="h-4 w-4" />
  </NuxtLink>
  <button
    v-else
    type="button"
    class="inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-ds-fg-subtle transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg"
    :aria-label="label"
    @click="$emit('click')"
  >
    <UIcon :name="icon" class="h-4 w-4" />
  </button>
</template>

<script setup lang="ts">
defineProps<{
  icon: string
  label: string
  to?: string
}>()

defineEmits<{
  click: []
}>()
</script>
