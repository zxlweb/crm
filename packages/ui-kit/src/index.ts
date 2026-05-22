import './styles/design-system.css'

// Tokens & theme
export * from './tokens'
export { UI_KIT_THEME_KEY, type UiKitThemeContext } from './theme/context'
export { bridgeUiKitThemeFromApp, useUiKitTheme } from './theme/bridge'

// Chart subsystem
export * from './chart/types'
export * from './chart/utils/colors'
export * from './chart/utils/echarts-parts'
export { useChartTheme } from './chart/use-chart-theme'

// Vue components (named imports for non-Nuxt consumers)
export { default as ChartShell } from './components/ui/chart/chart-shell.vue'
export { default as ChartLine } from './components/ui/chart/chart-line.vue'
export { default as ChartBar } from './components/ui/chart/chart-bar.vue'
export { default as ChartFunnel } from './components/ui/chart/chart-funnel.vue'
export { default as UiThemeToggle } from './components/ui/theme-toggle.vue'
export { default as CardMetric } from './components/ui/card/card-metric.vue'
export { default as CardShell } from './components/ui/card/card-shell.vue'
export * from './card/types'

// Tag / activity SVG icons
export * from './icons/tag-icons'
export * from './icons/sentiment-emoji'
export { default as UiTagIcon } from './components/ui/tag-icon.vue'
export { default as UiSentimentEmoji } from './components/ui/sentiment-emoji.vue'
