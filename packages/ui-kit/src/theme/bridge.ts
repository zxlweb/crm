import { computed, ref, type ComputedRef } from 'vue'
import { applyThemeToDocument, DEFAULT_THEME, type ThemeId } from '../tokens'
import type { UiKitThemeContext } from './context'

/** 应用层在 Nuxt 插件中注册；未注册时用本地 fallback（仅 document，不写 cookie） */
let appTheme: UiKitThemeContext | null = null

const fallbackId = ref<ThemeId>(DEFAULT_THEME)

export function bridgeUiKitThemeFromApp(ctx: UiKitThemeContext) {
  appTheme = ctx
}

/** ui-kit 组件统一消费主题（避免跨包 inject 失效） */
export function useUiKitTheme(): {
  id: ComputedRef<ThemeId>
  setTheme: (theme: ThemeId) => void
  isDark: ComputedRef<boolean>
} {
  if (appTheme) {
    return {
      id: computed(() => appTheme!.id.value as ThemeId),
      setTheme: (t) => appTheme!.setTheme(t),
      isDark: computed(() => appTheme!.isDark.value),
    }
  }

  return {
    id: computed(() => fallbackId.value),
    setTheme: (t) => {
      fallbackId.value = t
      applyThemeToDocument(t)
    },
    isDark: computed(() => fallbackId.value === 'v2'),
  }
}
