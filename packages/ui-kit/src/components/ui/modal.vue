<template>
  <UModal v-model="open">
    <div class="ds-modal flex flex-col overflow-hidden rounded-2xl">
      <span
        v-if="title || $slots.title || !hideClose"
        class="pointer-events-none block h-0.5 w-full shrink-0"
        :style="{ background: 'var(--ds-brand-gradient)' }"
        aria-hidden="true"
      />
      <header
        v-if="title || $slots.title || !hideClose"
        class="ds-modal__header flex items-start justify-between gap-3 border-b border-ds-border-muted px-5 py-4 sm:px-6"
      >
        <div class="min-w-0 flex-1">
          <slot name="title">
            <h2 v-if="title" class="text-base font-semibold tracking-tight text-ds-fg-heading">
              {{ title }}
            </h2>
          </slot>
          <p v-if="subtitle || $slots.subtitle" class="mt-1 text-xs text-ds-fg-muted">
            <slot name="subtitle">{{ subtitle }}</slot>
          </p>
        </div>
        <button
          v-if="!hideClose"
          type="button"
          class="ds-modal__close inline-flex h-8 w-8 shrink-0 cursor-pointer items-center justify-center rounded-lg text-ds-fg-subtle transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-heading focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ds-brand/30"
          :aria-label="closeAriaLabel"
          data-testid="ui-modal-close"
          @click="open = false"
        >
          <UIcon name="i-heroicons-x-mark" class="h-4 w-4" aria-hidden="true" />
        </button>
      </header>

      <div class="ds-modal__body px-5 py-5 sm:px-6">
        <slot />
      </div>

      <footer
        v-if="$slots.footer"
        class="ds-modal__footer flex flex-wrap items-center justify-end gap-2 border-t border-ds-border-muted bg-ds-bg-muted/30 px-5 py-3.5 sm:px-6"
      >
        <slot name="footer" />
      </footer>
    </div>
  </UModal>
</template>

<script setup lang="ts">
const open = defineModel<boolean>('open', { default: false })

withDefaults(
  defineProps<{
    title?: string
    subtitle?: string
    hideClose?: boolean
    closeAriaLabel?: string
  }>(),
  {
    title: '',
    subtitle: '',
    hideClose: false,
    closeAriaLabel: 'Close',
  },
)
</script>
