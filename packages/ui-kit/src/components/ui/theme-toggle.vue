<template>
  <div
    class="inline-flex items-center gap-1 rounded-xl border border-ds-border bg-ds-bg-elevated p-1"
    role="group"
    :aria-label="ariaLabel"
  >
    <button
      v-for="t in themeList"
      :key="t.id"
      type="button"
      class="cursor-pointer rounded-lg px-3 py-1.5 text-xs font-medium transition-colors"
      :class="id === t.id
        ? 'bg-ds-brand text-ds-on-brand shadow-sm'
        : 'text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg'"
      :aria-pressed="id === t.id"
      @click="setTheme(t.id)"
    >
      {{ labelFor(t.id) }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { THEMES, type ThemeId } from '../../tokens'
import { useUiKitTheme } from '../../theme/bridge'

const props = withDefaults(
  defineProps<{
    ariaLabel?: string
    /** 自定义按钮文案，默认 V1 / V2 */
    formatLabel?: (id: ThemeId) => string
  }>(),
  {
    ariaLabel: 'Theme version',
    formatLabel: (id: ThemeId) => id.toUpperCase(),
  },
)

const { id, setTheme } = useUiKitTheme()
const themeList = THEMES
const labelFor = (themeId: ThemeId) => props.formatLabel(themeId)
</script>
