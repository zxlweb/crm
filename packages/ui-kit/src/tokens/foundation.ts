/**
 * Foundation tokens — 与主题无关的尺度（CSS :root）
 * 运行时值见 styles/design-system.css；此处供文档、/design 预览、类型约束。
 */

export const DS_FONT_FAMILY = {
  sans: { token: '--ds-font-sans', value: '"Plus Jakarta Sans", ui-sans-serif, system-ui, sans-serif' },
  mono: { token: '--ds-font-mono', value: 'ui-monospace, "Cascadia Code", monospace' },
} as const

export const DS_FONT_SIZE = {
  xs: { token: '--ds-text-xs', value: '0.75rem', lineHeight: '--ds-leading-xs' },
  sm: { token: '--ds-text-sm', value: '0.875rem', lineHeight: '--ds-leading-sm' },
  base: { token: '--ds-text-base', value: '1rem', lineHeight: '--ds-leading-base' },
  lg: { token: '--ds-text-lg', value: '1.125rem', lineHeight: '--ds-leading-lg' },
  xl: { token: '--ds-text-xl', value: '1.25rem', lineHeight: '--ds-leading-xl' },
  '2xl': { token: '--ds-text-2xl', value: '1.5rem', lineHeight: '--ds-leading-2xl' },
  '3xl': { token: '--ds-text-3xl', value: '1.875rem', lineHeight: '--ds-leading-3xl' },
  '4xl': { token: '--ds-text-4xl', value: '2.25rem', lineHeight: '--ds-leading-4xl' },
} as const

export const DS_FONT_WEIGHT = {
  normal: { token: '--ds-font-normal', value: '400' },
  medium: { token: '--ds-font-medium', value: '500' },
  semibold: { token: '--ds-font-semibold', value: '600' },
  bold: { token: '--ds-font-bold', value: '700' },
} as const

export const DS_SPACE = {
  0: { token: '--ds-space-0', value: '0' },
  1: { token: '--ds-space-1', value: '0.25rem' },
  2: { token: '--ds-space-2', value: '0.5rem' },
  3: { token: '--ds-space-3', value: '0.75rem' },
  4: { token: '--ds-space-4', value: '1rem' },
  5: { token: '--ds-space-5', value: '1.25rem' },
  6: { token: '--ds-space-6', value: '1.5rem' },
  8: { token: '--ds-space-8', value: '2rem' },
  10: { token: '--ds-space-10', value: '2.5rem' },
  12: { token: '--ds-space-12', value: '3rem' },
  16: { token: '--ds-space-16', value: '4rem' },
} as const

export const DS_RADIUS = {
  none: { token: '--ds-radius-none', value: '0' },
  sm: { token: '--ds-radius-sm', value: '0.375rem' },
  md: { token: '--ds-radius-md', value: '0.5rem' },
  lg: { token: '--ds-radius-lg', value: '0.75rem' },
  xl: { token: '--ds-radius-xl', value: '1rem' },
  '2xl': { token: '--ds-radius-2xl', value: '1.25rem' },
  full: { token: '--ds-radius-full', value: '9999px' },
} as const

export const DS_BORDER_WIDTH = {
  DEFAULT: { token: '--ds-border-width', value: '1px' },
  0: { token: '--ds-border-width-0', value: '0' },
  2: { token: '--ds-border-width-2', value: '2px' },
} as const

export const DS_MOTION = {
  fast: { token: '--ds-duration-fast', value: '150ms' },
  normal: { token: '--ds-duration-normal', value: '200ms' },
  slow: { token: '--ds-duration-slow', value: '300ms' },
  slower: { token: '--ds-duration-slower', value: '500ms' },
} as const

export const DS_EASING = {
  default: { token: '--ds-ease-default', value: 'cubic-bezier(0.4, 0, 0.2, 1)' },
  in: { token: '--ds-ease-in', value: 'cubic-bezier(0.4, 0, 1, 1)' },
  out: { token: '--ds-ease-out', value: 'cubic-bezier(0, 0, 0.2, 1)' },
  inOut: { token: '--ds-ease-in-out', value: 'cubic-bezier(0.4, 0, 0.2, 1)' },
} as const

export const DS_Z_INDEX = {
  base: { token: '--ds-z-base', value: '0' },
  dropdown: { token: '--ds-z-dropdown', value: '1000' },
  sticky: { token: '--ds-z-sticky', value: '1020' },
  fixed: { token: '--ds-z-fixed', value: '1030' },
  modalBackdrop: { token: '--ds-z-modal-backdrop', value: '1040' },
  modal: { token: '--ds-z-modal', value: '1050' },
  popover: { token: '--ds-z-popover', value: '1060' },
  tooltip: { token: '--ds-z-tooltip', value: '1070' },
  toast: { token: '--ds-z-toast', value: '1080' },
} as const

export const DS_LAYOUT = {
  sidebarWidth: { token: '--ds-sidebar-width', value: '16.25rem' },
  topbarHeight: { token: '--ds-topbar-height', value: '4rem' },
  contentMax: { token: '--ds-content-max', value: '80rem' },
  pagePaddingX: { token: '--ds-page-px', value: 'var(--ds-space-6)' },
  pagePaddingY: { token: '--ds-page-py', value: 'var(--ds-space-8)' },
} as const

export const DS_SIZE = {
  inputHeight: { token: '--ds-input-height', value: '2.5rem' },
  buttonHeightSm: { token: '--ds-button-height-sm', value: '2rem' },
  buttonHeightMd: { token: '--ds-button-height-md', value: '2.5rem' },
  buttonHeightLg: { token: '--ds-button-height-lg', value: '3rem' },
  iconSm: { token: '--ds-icon-sm', value: '1rem' },
  iconMd: { token: '--ds-icon-md', value: '1.25rem' },
  iconLg: { token: '--ds-icon-lg', value: '1.5rem' },
} as const
