/**
 * EnterpriseFlow 设计系统 — 主题 token 定义
 * V1：默认浅色紫 SaaS | V2：深色 Optrixx Dashboard（切换主题）
 *
 * 运行时值见 assets/css/design-system.css（CSS 变量）
 * Tailwind 语义类：bg-ds-bg、text-ds-fg、bg-ds-brand 等
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

/** 设计 token 分组（文档与类型参考） */
export const TOKEN_GROUPS = [
  'surface', // 背景层级：bg, elevated, muted, input, sidebar
  'foreground', // 文字：fg, heading, muted, subtle, brand
  'border',
  'brand', // 主色与交互
  'status', // success, danger
  'chart', // 图表 SVG 用
  'shadow',
] as const
