/**
 * Card 场景注册表 — 与 Chart 子系统并列。
 * 新增固定设计风格时：加场景 id + 对应 Card* 组件，在 05-component-scenarios.md 补一行即可。
 */

export type CardScenarioId = 'dashboard' | 'content'

export interface CardScenarioMeta {
  id: CardScenarioId
  /** Nuxt 自动导入组件名 */
  component: 'CardMetric' | 'CardShell'
  /** i18n / 文档用简述 */
  descKey: string
  /** 典型页面 */
  pages: string[]
}

/** 看板 KPI：左方标 + 主数值 + 底栏同比 */
export const CARD_SCENARIOS: Record<CardScenarioId, CardScenarioMeta> = {
  dashboard: {
    id: 'dashboard',
    component: 'CardMetric',
    descKey: 'cardScenarioDashboard',
    pages: ['/admin', '/dashboard', '仪表盘顶栏 KPI'],
  },
  content: {
    id: 'content',
    component: 'CardShell',
    descKey: 'cardScenarioContent',
    pages: ['表格容器', '表单区块', '详情页模块'],
  },
}

export type CardMetricIconTone = 'info' | 'calendar' | 'accent' | 'brand' | 'neutral'
export type CardMetricTrendDirection = 'up' | 'down' | 'flat'
