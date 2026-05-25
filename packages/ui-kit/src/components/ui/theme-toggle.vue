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
      class="cursor-pointer rounded-lg transition-colors duration-200"
      :class="[
        id === t.id
          ? 'bg-ds-brand text-ds-on-brand shadow-sm'
          : 'text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg',
        variant === 'icon' ? 'flex h-8 w-8 items-center justify-center' : 'px-3 py-1.5 text-xs font-medium',
      ]"
      :aria-pressed="id === t.id"
      :aria-label="labelFor(t.id)"
      :title="labelFor(t.id)"
      @click="setTheme(t.id)"
    >
      <template v-if="variant === 'icon'">
        <svg
          v-if="t.id === 'v1'"
          class="h-4 w-4 shrink-0"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.75"
          aria-hidden="true"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"
          />
        </svg>
        <svg
          v-else
          class="h-4 w-4 shrink-0"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.75"
          aria-hidden="true"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z"
          />
        </svg>
      </template>
      <template v-else>
        {{ labelFor(t.id) }}
      </template>
    </button>
  </div>
</template>

<script setup lang="ts">
import { THEMES, type ThemeId } from '../../tokens'
import { useUiKitTheme } from '../../theme/bridge'

const props = withDefaults(
  defineProps<{
    ariaLabel?: string
    variant?: 'text' | 'icon'
    /** 自定义按钮文案 / aria-label，默认 V1 / V2 */
    formatLabel?: (id: ThemeId) => string
  }>(),
  {
    ariaLabel: 'Theme version',
    variant: 'text',
    formatLabel: (id: ThemeId) => id.toUpperCase(),
  },
)

const { id, setTheme } = useUiKitTheme()
const themeList = THEMES
const labelFor = (themeId: ThemeId) => props.formatLabel(themeId)
</script>
