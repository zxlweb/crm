<template>
  <UiThemeToggle
    variant="icon"
    class="shrink-0"
    :aria-label="$t('themeSkinLabel')"
    :format-label="skinLabel"
  />
</template>

<script setup lang="ts">
import type { ThemeId } from '@crm/ui-kit'

const { t } = useI18n()

function skinLabel(id: ThemeId) {
  return id === 'v1' ? t('themeSkinLight') : t('themeSkinDark')
}
</script>
