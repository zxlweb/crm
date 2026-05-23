/** Nuxt UI UPagination — 对齐订单表参考：独立圆角块 + Previous/Next 文案 */
export const crmPaginationUi = {
  wrapper: 'flex items-center gap-1.5',
  base: '',
  rounded: '',
  default: {
    size: 'sm' as const,
    activeButton: {
      color: 'primary' as const,
      variant: 'solid' as const,
      class:
        'min-w-[2.25rem] h-9 justify-center rounded-lg font-medium shadow-none bg-ds-brand text-ds-on-brand hover:bg-ds-brand-hover',
    },
    inactiveButton: {
      color: 'gray' as const,
      variant: 'outline' as const,
      class:
        'min-w-[2.25rem] h-9 justify-center rounded-lg font-medium shadow-none ring-1 ring-inset ring-ds-border bg-ds-bg-elevated text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg',
    },
    prevButton: {
      color: 'gray' as const,
      variant: 'outline' as const,
      class:
        'h-9 rounded-lg px-3 font-medium shadow-none ring-1 ring-inset ring-ds-border bg-ds-bg-elevated text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg',
      icon: 'i-heroicons-chevron-left-20-solid',
    },
    nextButton: {
      color: 'gray' as const,
      variant: 'outline' as const,
      class:
        'h-9 rounded-lg px-3 font-medium shadow-none ring-1 ring-inset ring-ds-border bg-ds-bg-elevated text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg',
      icon: 'i-heroicons-chevron-right-20-solid',
    },
  },
} as const
