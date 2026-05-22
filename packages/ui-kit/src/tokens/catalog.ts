/**
 * 设计系统 Token 目录 — 三层逻辑 + 分组，供文档与 /design 预览
 */

export type DsTokenLayerId = 'foundation' | 'semantic' | 'component'

export interface DsTokenLayer {
  id: DsTokenLayerId
  titleKey: string
  descKey: string
}

export const DS_TOKEN_LAYERS: DsTokenLayer[] = [
  {
    id: 'foundation',
    titleKey: 'dsLayerFoundationTitle',
    descKey: 'dsLayerFoundationDesc',
  },
  {
    id: 'semantic',
    titleKey: 'dsLayerSemanticTitle',
    descKey: 'dsLayerSemanticDesc',
  },
  {
    id: 'component',
    titleKey: 'dsLayerComponentTitle',
    descKey: 'dsLayerComponentDesc',
  },
]

export interface DsTokenEntry {
  name: string
  token: string
  /** CSS var() 或 preview 用 */
  preview?: 'color' | 'length' | 'shadow' | 'font' | 'text-sample' | 'space' | 'radius'
  tailwind?: string
  noteKey?: string
}

export interface DsTokenGroup {
  id: string
  layer: DsTokenLayerId
  titleKey: string
  tokens: DsTokenEntry[]
}

/** 完整分组目录（色彩项随 data-theme 变化） */
export const DS_TOKEN_GROUPS: DsTokenGroup[] = [
  // ─── Foundation ───
  {
    id: 'typography-family',
    layer: 'foundation',
    titleKey: 'dsGroupTypographyFamily',
    tokens: [
      { name: 'Sans', token: '--ds-font-sans', preview: 'font', tailwind: 'font-sans' },
      { name: 'Mono', token: '--ds-font-mono', preview: 'font', tailwind: 'font-mono' },
    ],
  },
  {
    id: 'typography-size',
    layer: 'foundation',
    titleKey: 'dsGroupTypographySize',
    tokens: [
      { name: 'XS', token: '--ds-text-xs', preview: 'text-sample', tailwind: 'text-ds-xs' },
      { name: 'SM', token: '--ds-text-sm', preview: 'text-sample', tailwind: 'text-ds-sm' },
      { name: 'Base', token: '--ds-text-base', preview: 'text-sample', tailwind: 'text-ds-base' },
      { name: 'LG', token: '--ds-text-lg', preview: 'text-sample', tailwind: 'text-ds-lg' },
      { name: 'XL', token: '--ds-text-xl', preview: 'text-sample', tailwind: 'text-ds-xl' },
      { name: '2XL', token: '--ds-text-2xl', preview: 'text-sample', tailwind: 'text-ds-2xl' },
      { name: '3XL', token: '--ds-text-3xl', preview: 'text-sample', tailwind: 'text-ds-3xl' },
    ],
  },
  {
    id: 'typography-weight',
    layer: 'foundation',
    titleKey: 'dsGroupTypographyWeight',
    tokens: [
      { name: 'Normal', token: '--ds-font-normal', preview: 'text-sample', tailwind: 'font-ds-normal' },
      { name: 'Medium', token: '--ds-font-medium', preview: 'text-sample', tailwind: 'font-ds-medium' },
      { name: 'Semibold', token: '--ds-font-semibold', preview: 'text-sample', tailwind: 'font-ds-semibold' },
      { name: 'Bold', token: '--ds-font-bold', preview: 'text-sample', tailwind: 'font-ds-bold' },
    ],
  },
  {
    id: 'spacing',
    layer: 'foundation',
    titleKey: 'dsGroupSpacing',
    tokens: [1, 2, 3, 4, 5, 6, 8, 10, 12].map((n) => ({
      name: String(n),
      token: `--ds-space-${n}`,
      preview: 'space' as const,
      tailwind: n === 4 ? 'p-ds-4' : `p-ds-${n}`,
    })),
  },
  {
    id: 'radius',
    layer: 'foundation',
    titleKey: 'dsGroupRadius',
    tokens: [
      { name: 'SM', token: '--ds-radius-sm', preview: 'radius', tailwind: 'rounded-ds-sm' },
      { name: 'MD', token: '--ds-radius-md', preview: 'radius', tailwind: 'rounded-ds-md' },
      { name: 'LG', token: '--ds-radius-lg', preview: 'radius', tailwind: 'rounded-ds-lg' },
      { name: 'XL', token: '--ds-radius-xl', preview: 'radius', tailwind: 'rounded-ds-xl' },
      { name: '2XL', token: '--ds-radius-2xl', preview: 'radius', tailwind: 'rounded-ds-2xl' },
      { name: 'Full', token: '--ds-radius-full', preview: 'radius', tailwind: 'rounded-ds-full' },
    ],
  },
  {
    id: 'shadow',
    layer: 'foundation',
    titleKey: 'dsGroupShadow',
    tokens: [
      { name: 'SM', token: '--ds-shadow-sm', preview: 'shadow', tailwind: 'shadow-ds-sm' },
      { name: 'MD', token: '--ds-shadow-md', preview: 'shadow', tailwind: 'shadow-ds-md' },
      { name: 'LG', token: '--ds-shadow-lg', preview: 'shadow', tailwind: 'shadow-ds-lg' },
      { name: 'Brand', token: '--ds-brand-shadow', preview: 'shadow', tailwind: 'shadow-ds-brand' },
    ],
  },
  {
    id: 'motion',
    layer: 'foundation',
    titleKey: 'dsGroupMotion',
    tokens: [
      { name: 'Fast', token: '--ds-duration-fast', preview: 'length', tailwind: 'duration-ds-fast' },
      { name: 'Normal', token: '--ds-duration-normal', preview: 'length', tailwind: 'duration-ds-normal' },
      { name: 'Slow', token: '--ds-duration-slow', preview: 'length', tailwind: 'duration-ds-slow' },
      { name: 'Ease default', token: '--ds-ease-default', preview: 'length' },
    ],
  },
  {
    id: 'z-index',
    layer: 'foundation',
    titleKey: 'dsGroupZIndex',
    tokens: [
      { name: 'Dropdown', token: '--ds-z-dropdown', tailwind: 'z-ds-dropdown' },
      { name: 'Sticky', token: '--ds-z-sticky', tailwind: 'z-ds-sticky' },
      { name: 'Modal', token: '--ds-z-modal', tailwind: 'z-ds-modal' },
      { name: 'Toast', token: '--ds-z-toast', tailwind: 'z-ds-toast' },
    ],
  },
  {
    id: 'layout',
    layer: 'foundation',
    titleKey: 'dsGroupLayout',
    tokens: [
      { name: 'Sidebar width', token: '--ds-sidebar-width', preview: 'length' },
      { name: 'Topbar height', token: '--ds-topbar-height', preview: 'length' },
      { name: 'Content max', token: '--ds-content-max', preview: 'length' },
    ],
  },
  // ─── Semantic (theme) ───
  {
    id: 'surface',
    layer: 'semantic',
    titleKey: 'dsGroupSurface',
    tokens: [
      { name: 'Page BG', token: '--ds-bg', preview: 'color', tailwind: 'bg-ds-bg' },
      { name: 'Elevated', token: '--ds-bg-elevated', preview: 'color', tailwind: 'bg-ds-bg-elevated' },
      { name: 'Muted', token: '--ds-bg-muted', preview: 'color', tailwind: 'bg-ds-bg-muted' },
      { name: 'Input', token: '--ds-bg-input', preview: 'color', tailwind: 'bg-ds-bg-input' },
      { name: 'Sidebar', token: '--ds-bg-sidebar', preview: 'color', tailwind: 'bg-ds-bg-sidebar' },
    ],
  },
  {
    id: 'foreground',
    layer: 'semantic',
    titleKey: 'dsGroupForeground',
    tokens: [
      { name: 'Body', token: '--ds-fg', preview: 'color', tailwind: 'text-ds-fg' },
      { name: 'Heading', token: '--ds-fg-heading', preview: 'color', tailwind: 'text-ds-fg-heading' },
      { name: 'Muted', token: '--ds-fg-muted', preview: 'color', tailwind: 'text-ds-fg-muted' },
      { name: 'Brand', token: '--ds-fg-brand', preview: 'color', tailwind: 'text-ds-fg-brand' },
      { name: 'On brand', token: '--ds-fg-on-brand', preview: 'color', tailwind: 'text-ds-on-brand' },
    ],
  },
  {
    id: 'border',
    layer: 'semantic',
    titleKey: 'dsGroupBorder',
    tokens: [
      { name: 'Default', token: '--ds-border', preview: 'color', tailwind: 'border-ds-border' },
      { name: 'Muted', token: '--ds-border-muted', preview: 'color', tailwind: 'border-ds-border-muted' },
      { name: 'Focus', token: '--ds-border-focus', preview: 'color' },
    ],
  },
  {
    id: 'brand',
    layer: 'semantic',
    titleKey: 'dsGroupBrand',
    tokens: [
      { name: 'Brand', token: '--ds-brand', preview: 'color', tailwind: 'bg-ds-brand' },
      { name: 'Brand hover', token: '--ds-brand-hover', preview: 'color', tailwind: 'bg-ds-brand-hover' },
      { name: 'Brand muted', token: '--ds-brand-muted', preview: 'color', tailwind: 'bg-ds-brand-muted' },
    ],
  },
  {
    id: 'status',
    layer: 'semantic',
    titleKey: 'dsGroupStatus',
    tokens: [
      { name: 'Success', token: '--ds-success', preview: 'color', tailwind: 'text-ds-success' },
      { name: 'Danger', token: '--ds-danger', preview: 'color', tailwind: 'text-ds-danger' },
      { name: 'Success subtle BG', token: '--ds-bg-success-subtle', preview: 'color' },
      { name: 'Danger subtle BG', token: '--ds-bg-danger-subtle', preview: 'color' },
    ],
  },
  {
    id: 'chart',
    layer: 'semantic',
    titleKey: 'dsGroupChart',
    tokens: [
      { name: 'Line end', token: '--ds-chart-line-end', preview: 'color' },
      { name: 'Line start', token: '--ds-chart-line-start', preview: 'color' },
      { name: 'Grid', token: '--ds-chart-grid', preview: 'color' },
      { name: 'Tooltip BG', token: '--ds-chart-tooltip-bg', preview: 'color' },
    ],
  },
  // ─── Component ───
  {
    id: 'component-size',
    layer: 'component',
    titleKey: 'dsGroupComponentSize',
    tokens: [
      { name: 'Input height', token: '--ds-input-height', preview: 'length' },
      { name: 'Button MD', token: '--ds-button-height-md', preview: 'length' },
      { name: 'Icon MD', token: '--ds-icon-md', preview: 'length' },
    ],
  },
  {
    id: 'component-utils',
    layer: 'component',
    titleKey: 'dsGroupComponentUtils',
    tokens: [
      { name: '.ds-card', token: '—', noteKey: 'dsUtilCard' },
      { name: '.ds-input', token: '—', noteKey: 'dsUtilInput' },
      { name: '.ds-btn-primary', token: '—', noteKey: 'dsUtilBtnPrimary' },
      { name: '.ds-nav-active', token: '—', noteKey: 'dsUtilNavActive' },
    ],
  },
]

/** @deprecated 使用 DS_TOKEN_GROUPS 的 layer 字段 */
export const TOKEN_GROUPS = [
  'surface',
  'foreground',
  'border',
  'brand',
  'status',
  'chart',
  'shadow',
  'typography',
  'spacing',
  'radius',
  'motion',
] as const
