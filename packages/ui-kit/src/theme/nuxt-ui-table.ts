/**
 * Nuxt UI UTable 主题覆盖 — 对齐 CRM 订单表参考风格（浅灰表头、柔和分隔、行 hover）
 * 在 apps/web/app.config.ts 中挂到 ui.table，或由 UiTable 组件 merge。
 */
export const crmTableUi = {
  wrapper: 'relative overflow-x-auto',
  base: 'min-w-full',
  divide: 'divide-y divide-ds-border/80',
  thead: 'bg-ds-bg-muted/70',
  tbody: 'divide-y divide-ds-border/80',
  tr: {
    base: 'transition-colors duration-200',
    selected: 'bg-ds-bg-muted',
    expanded: 'bg-ds-bg-muted/80',
    active: 'hover:bg-ds-bg-muted/60',
  },
  th: {
    base: 'text-left rtl:text-right',
    padding: 'px-4 py-3 first:ps-5 last:pe-5',
    color: 'text-ds-fg-muted',
    font: 'font-medium',
    size: 'text-xs',
  },
  td: {
    base: 'whitespace-nowrap',
    padding: 'px-4 py-3.5 first:ps-5 last:pe-5',
    color: 'text-ds-fg-muted',
    font: '',
    size: 'text-sm',
  },
  checkbox: {
    padding: 'ps-5',
  },
  loadingState: {
    wrapper: 'flex flex-col items-center justify-center px-6 py-14',
    label: 'text-sm text-center text-ds-fg-muted',
    icon: 'mb-4 h-6 w-6 animate-spin text-ds-fg-subtle',
  },
  emptyState: {
    wrapper: 'flex flex-col items-center justify-center px-6 py-14',
    label: 'text-sm text-center text-ds-fg-muted',
    icon: 'mb-4 h-6 w-6 text-ds-fg-subtle',
  },
  default: {
    sortAscIcon: 'i-heroicons-chevron-up-20-solid',
    sortDescIcon: 'i-heroicons-chevron-down-20-solid',
    sortButton: {
      icon: 'i-heroicons-arrows-up-down-20-solid',
      trailing: true,
      square: true,
      color: 'gray',
      variant: 'ghost',
      size: 'xs',
      class: '-m-0.5 ms-1 text-ds-fg-subtle/80 hover:text-ds-fg-muted',
    },
    checkbox: {
      color: 'primary',
      rounded: 'rounded-md',
    },
    loadingState: {
      icon: 'i-heroicons-arrow-path-20-solid',
      label: 'Loading...',
    },
    emptyState: {
      icon: 'i-heroicons-circle-stack-20-solid',
      label: 'No items.',
    },
  },
} as const
