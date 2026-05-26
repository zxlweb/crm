<template>
  <nav
    class="ds-settings-nav flex flex-col gap-1"
    :aria-label="$t('settingsNavAria')"
    data-testid="settings-nav"
  >
    <button
      v-for="item in items"
      :key="item.key"
      type="button"
      class="ds-settings-nav__item group flex w-full cursor-pointer items-start gap-3 rounded-xl border border-transparent px-3 py-3 text-left transition-colors duration-200"
      :class="modelValue === item.key
        ? 'border-ds-border bg-ds-brand-subtle text-ds-fg-brand shadow-ds-sm'
        : 'text-ds-fg hover:bg-ds-bg-muted hover:text-ds-fg-heading'"
      :data-testid="`settings-tab-${item.key}`"
      :aria-current="modelValue === item.key ? 'page' : undefined"
      @click="$emit('update:modelValue', item.key)"
    >
      <span
        class="ds-settings-nav__icon flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border transition-colors duration-200"
        :class="modelValue === item.key
          ? 'border-ds-brand/30 bg-ds-bg-elevated text-ds-fg-brand'
          : 'border-ds-border-muted bg-ds-bg-muted text-ds-fg-muted group-hover:border-ds-border group-hover:text-ds-fg-heading'"
        aria-hidden="true"
      >
        <UIcon :name="item.icon" class="h-5 w-5" />
      </span>
      <span class="min-w-0 flex-1">
        <span class="flex items-center gap-2">
          <span class="truncate text-sm font-semibold">{{ item.label }}</span>
          <span
            v-if="item.badge"
            class="rounded-full bg-ds-bg-elevated px-2 py-0.5 text-[10px] font-semibold text-ds-fg-muted ring-1 ring-inset ring-ds-border-muted"
          >
            {{ item.badge }}
          </span>
        </span>
        <span class="mt-0.5 block text-xs leading-relaxed text-ds-fg-muted">
          {{ item.description }}
        </span>
      </span>
    </button>
  </nav>
</template>

<script setup lang="ts">
export interface SettingsNavItem {
  key: string
  label: string
  description: string
  icon: string
  badge?: string | number
}

defineProps<{
  modelValue: string
  items: SettingsNavItem[]
}>()

defineEmits<{
  'update:modelValue': [key: string]
}>()
</script>
