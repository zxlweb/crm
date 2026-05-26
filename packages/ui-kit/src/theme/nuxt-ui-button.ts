/** Nuxt UI UButton — 与 design-system 对齐：圆角、品牌渐变 primary、克制 secondary */
export const crmButtonUi = {
  base: 'focus:outline-none focus-visible:outline-0 disabled:cursor-not-allowed disabled:opacity-60 flex-shrink-0 transition-[background-color,box-shadow,transform,color] duration-200 active:scale-[0.98] motion-reduce:transform-none',
  font: 'font-semibold',
  rounded: 'rounded-xl',
  size: {
    '2xs': 'text-[11px]',
    xs: 'text-xs',
    sm: 'text-xs',
    md: 'text-sm',
    lg: 'text-sm',
    xl: 'text-base',
  },
  gap: {
    '2xs': 'gap-x-1',
    xs: 'gap-x-1.5',
    sm: 'gap-x-1.5',
    md: 'gap-x-2',
    lg: 'gap-x-2',
    xl: 'gap-x-2.5',
  },
  padding: {
    '2xs': 'px-2 py-1',
    xs: 'px-2.5 py-1.5',
    sm: 'px-3 py-1.5',
    md: 'px-3.5 py-2',
    lg: 'px-4 py-2.5',
    xl: 'px-5 py-3',
  },
  square: {
    '2xs': 'p-1',
    xs: 'p-1.5',
    sm: 'p-1.5',
    md: 'p-2',
    lg: 'p-2.5',
    xl: 'p-3',
  },
  color: {
    primary: {
      solid:
        'shadow-ds-brand text-ds-on-brand bg-ds-brand hover:bg-ds-brand-hover hover:shadow-ds-lg focus-visible:ring-2 focus-visible:ring-ds-brand/40 focus-visible:ring-offset-2 focus-visible:ring-offset-ds-bg disabled:bg-ds-brand/60 disabled:shadow-none',
      outline:
        'shadow-none ring-1 ring-inset ring-ds-border bg-ds-bg-elevated text-ds-fg hover:bg-ds-bg-muted hover:ring-ds-brand-muted hover:text-ds-fg-brand focus-visible:ring-2 focus-visible:ring-ds-brand/30',
      soft: 'shadow-none bg-ds-brand-subtle/70 text-ds-fg-brand hover:bg-ds-brand-subtle focus-visible:ring-2 focus-visible:ring-ds-brand/30',
      ghost:
        'shadow-none text-ds-fg-brand hover:bg-ds-brand-subtle/60 focus-visible:ring-2 focus-visible:ring-ds-brand/25',
      link: 'shadow-none text-ds-fg-brand underline-offset-4 hover:underline focus-visible:ring-2 focus-visible:ring-ds-brand/25',
    },
    gray: {
      solid:
        'shadow-none bg-ds-bg-muted text-ds-fg hover:bg-ds-border-muted focus-visible:ring-2 focus-visible:ring-ds-brand/20',
      outline:
        'shadow-none ring-1 ring-inset ring-ds-border bg-ds-bg-elevated text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg hover:ring-ds-brand-muted focus-visible:ring-2 focus-visible:ring-ds-brand/25',
      soft: 'shadow-none bg-ds-bg-muted/80 text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg focus-visible:ring-2 focus-visible:ring-ds-brand/20',
      ghost:
        'shadow-none text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg focus-visible:ring-2 focus-visible:ring-ds-brand/20',
      link: 'shadow-none text-ds-fg-muted underline-offset-4 hover:underline focus-visible:ring-2 focus-visible:ring-ds-brand/20',
    },
    red: {
      solid:
        'shadow-none bg-ds-danger text-white hover:opacity-90 focus-visible:ring-2 focus-visible:ring-ds-danger/40',
      outline:
        'shadow-none ring-1 ring-inset ring-ds-danger/30 bg-ds-bg-elevated text-ds-danger hover:bg-ds-danger-subtle focus-visible:ring-2 focus-visible:ring-ds-danger/30',
      soft: 'shadow-none bg-ds-danger-subtle text-ds-danger hover:bg-ds-danger-subtle/80 focus-visible:ring-2 focus-visible:ring-ds-danger/25',
      ghost:
        'shadow-none text-ds-danger hover:bg-ds-danger-subtle focus-visible:ring-2 focus-visible:ring-ds-danger/25',
    },
  },
  icon: {
    base: 'flex-shrink-0',
    loading: 'animate-spin',
    size: {
      '2xs': 'h-3 w-3',
      xs: 'h-3.5 w-3.5',
      sm: 'h-4 w-4',
      md: 'h-4 w-4',
      lg: 'h-5 w-5',
      xl: 'h-5 w-5',
    },
  },
  default: {
    size: 'md' as const,
    color: 'primary' as const,
    variant: 'solid' as const,
    loadingIcon: 'i-heroicons-arrow-path-20-solid',
  },
} as const
