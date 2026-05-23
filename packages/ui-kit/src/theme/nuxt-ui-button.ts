/** Nuxt UI UButton — 列表页等场景：primary 略克制，secondary 更中性 */
export const crmButtonUi = {
  rounded: 'rounded-lg',
  font: 'font-medium',
  default: {
    size: 'md' as const,
    color: 'primary' as const,
    variant: 'solid' as const,
  },
  color: {
    primary: {
      solid:
        'shadow-none bg-ds-brand text-ds-on-brand hover:bg-ds-brand-hover active:bg-ds-brand-hover disabled:bg-ds-brand/60',
      outline:
        'ring-1 ring-inset ring-ds-border bg-ds-bg-elevated text-ds-fg hover:bg-ds-bg-muted shadow-none',
      soft: 'bg-ds-brand-subtle/60 text-ds-fg-brand hover:bg-ds-brand-subtle shadow-none',
      ghost: 'text-ds-fg-brand hover:bg-ds-brand-subtle/50 shadow-none',
    },
    gray: {
      solid: 'shadow-none bg-ds-bg-muted text-ds-fg hover:bg-ds-border-muted',
      outline:
        'ring-1 ring-inset ring-ds-border bg-ds-bg-elevated text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg shadow-none',
      soft: 'bg-ds-bg-muted/80 text-ds-fg-muted hover:bg-ds-bg-muted shadow-none',
      ghost: 'text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg shadow-none',
    },
  },
} as const
