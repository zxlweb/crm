<template>
  <div class="inline-flex items-center gap-1 rounded-xl border border-ds-border bg-ds-bg-elevated p-1" role="group" :aria-label="$t('themeSwitchVersion')">
    <button
      v-for="t in themes"
      :key="t.id"
      type="button"
      class="cursor-pointer rounded-lg px-3 py-1.5 text-xs font-medium transition-colors"
      :class="id === t.id
        ? 'bg-ds-brand text-ds-on-brand shadow-sm'
        : 'text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg'"
      :aria-pressed="id === t.id"
      @click="setTheme(t.id)"
    >
      {{ t.id.toUpperCase() }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { THEMES } from '~/design-system/tokens'

const { id, setTheme } = useTheme()
const themes = THEMES
</script>
