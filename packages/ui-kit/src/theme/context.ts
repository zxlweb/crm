import type { ComputedRef, InjectionKey, Ref } from 'vue'
import type { ThemeId } from '../tokens'

/** 应用层 provide，供图表/主题切换组件消费（与 Nuxt useCookie 解耦） */
export interface UiKitThemeContext {
  id: Ref<ThemeId> | ComputedRef<ThemeId>
  setTheme: (theme: ThemeId) => void
  isDark: ComputedRef<boolean>
}

/** 字符串 key，避免 app 与 ui-kit 分包时 Symbol 实例不一致 */
export const UI_KIT_THEME_KEY: InjectionKey<UiKitThemeContext> = 'crm-ui-kit-theme'
