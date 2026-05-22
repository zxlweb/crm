/**
 * 主题元数据（色彩语义由 design-system.css 按 data-theme 切换）
 */

export type ThemeId = 'v1' | 'v2'

export const THEME_COOKIE = 'crm-theme'
export const DEFAULT_THEME: ThemeId = 'v1'

export interface ThemeMeta {
  id: ThemeId
  nameKey: string
  descKey: string
}

export const THEMES: ThemeMeta[] = [
  { id: 'v1', nameKey: 'themeV1Name', descKey: 'themeV1Desc' },
  { id: 'v2', nameKey: 'themeV2Name', descKey: 'themeV2Desc' },
]

export function applyThemeToDocument(theme: ThemeId) {
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('data-theme', theme)
  }
}
