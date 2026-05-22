<template>
  <UModal v-model="open" :title="title">
    <div class="space-y-4 p-4">
      <slot />
    </div>
    <template v-if="$slots.footer" #footer>
      <div class="flex justify-end gap-2 p-4">
        <slot name="footer" />
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
const open = defineModel<boolean>('open', { default: false })

defineProps<{
  title?: string
}>()
</script>
