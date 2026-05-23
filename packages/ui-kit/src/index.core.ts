/**
 * Server-safe entry — 不含 .vue 导出，供 Nuxt Nitro / app.config 等在 dev 下解析。
 * 组件由 @crm/ui-kit/nuxt 模块 auto-import；库构建仍走 index.ts。
 */
import './styles/design-system.css'

export * from './tokens'
export { UI_KIT_THEME_KEY, type UiKitThemeContext } from './theme/context'
export { bridgeUiKitThemeFromApp, useUiKitTheme } from './theme/bridge'

export * from './chart/types'
export * from './chart/utils/colors'
export * from './chart/utils/echarts-parts'
export { useChartTheme } from './chart/use-chart-theme'

export * from './icons/tag-icons'
export * from './icons/sentiment-emoji'
export * from './card/types'

export { crmTableUi } from './theme/nuxt-ui-table'
export { crmPaginationUi } from './theme/nuxt-ui-pagination'
export { crmInputUi } from './theme/nuxt-ui-input'
export { crmSelectUi } from './theme/nuxt-ui-select'
export { crmButtonUi } from './theme/nuxt-ui-button'
export type { UiTableColumn } from './components/ui/table-types'
