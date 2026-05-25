<template>
  <!-- Nuxt UI v2 UModal 仅默认插槽，footer 需放在 panel 内 -->
  <UModal v-model="open" :title="title">
    <div class="space-y-4 p-4">
      <slot />
      <div
        v-if="$slots.footer"
        class="flex justify-end gap-2 border-t border-ds-border pt-4"
      >
        <slot name="footer" />
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
const open = defineModel<boolean>('open', { default: false })

defineProps<{
  title?: string
}>()
</script>
