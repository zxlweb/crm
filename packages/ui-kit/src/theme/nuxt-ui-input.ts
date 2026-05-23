/** Nuxt UI UInput — 对齐 CRM design-system，弱化默认 gray ring + primary focus */
export const crmInputUi = {
  rounded: 'rounded-lg',
  placeholder: 'placeholder:text-ds-fg-muted',
  color: {
    white: {
      outline:
        'bg-ds-bg-elevated text-ds-fg shadow-none ring-1 ring-inset ring-ds-border focus:ring-2 focus:ring-ds-brand/12',
    },
    gray: {
      outline:
        'bg-ds-bg-muted/40 text-ds-fg shadow-none ring-1 ring-inset ring-transparent hover:bg-ds-bg-muted/70 focus:bg-ds-bg-elevated focus:ring-2 focus:ring-ds-brand/10',
    },
  },
  icon: {
    base: 'text-ds-fg-subtle',
  },
  default: {
    size: 'md' as const,
    color: 'gray' as const,
    variant: 'outline' as const,
    loadingIcon: 'i-heroicons-arrow-path-20-solid',
  },
} as const
