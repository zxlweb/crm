/** Nuxt UI UInput — 对齐 CRM design-system（勿依赖 dark: 或 /opacity，由 .form-input 全局样式兜底） */
export const crmInputUi = {
  rounded: 'rounded-lg',
  placeholder: 'placeholder:text-ds-fg-muted',
  color: {
    white: {
      outline:
        'bg-ds-bg-input text-ds-fg shadow-none ring-1 ring-inset ring-ds-border hover:ring-ds-border focus:ring-2 focus:ring-ds-brand/20',
    },
    gray: {
      outline:
        'bg-ds-bg-input text-ds-fg shadow-none ring-1 ring-inset ring-ds-border hover:ring-ds-border focus:ring-2 focus:ring-ds-brand/20',
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
