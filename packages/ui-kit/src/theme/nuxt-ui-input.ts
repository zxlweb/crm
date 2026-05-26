/** Nuxt UI UInput — 对齐 CRM design-system；柔和聚焦 + 平滑过渡 */
export const crmInputUi = {
  wrapper: 'relative',
  base:
    'relative block w-full disabled:cursor-not-allowed disabled:opacity-60 focus:outline-none transition-[box-shadow,background-color,ring] duration-200',
  form: 'form-input',
  rounded: 'rounded-xl',
  placeholder: 'placeholder:text-ds-fg-subtle',
  file: {
    base: 'file:mr-1.5 file:font-medium file:text-ds-fg-muted file:bg-transparent file:border-0 file:p-0 file:outline-none',
  },
  size: {
    '2xs': 'text-xs',
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
    lg: 'gap-x-2.5',
    xl: 'gap-x-2.5',
  },
  padding: {
    '2xs': 'px-2 py-0.5',
    xs: 'px-2.5 py-1',
    sm: 'px-2.5 py-1.5',
    md: 'px-3 py-1.5',
    lg: 'px-3.5 py-2',
    xl: 'px-4 py-2.5',
  },
  leading: {
    padding: {
      '2xs': 'ps-7',
      xs: 'ps-8',
      sm: 'ps-9',
      md: 'ps-10',
      lg: 'ps-11',
      xl: 'ps-12',
    },
  },
  trailing: {
    padding: {
      '2xs': 'pe-7',
      xs: 'pe-8',
      sm: 'pe-9',
      md: 'pe-10',
      lg: 'pe-11',
      xl: 'pe-12',
    },
  },
  color: {
    white: {
      outline:
        'bg-ds-bg-input text-ds-fg ring-1 ring-inset ring-ds-border hover:ring-ds-brand-muted focus:ring-2 focus:ring-ds-brand/35 focus:ring-inset',
    },
    gray: {
      outline:
        'bg-ds-bg-input text-ds-fg ring-1 ring-inset ring-ds-border hover:ring-ds-brand-muted focus:ring-2 focus:ring-ds-brand/35 focus:ring-inset',
    },
  },
  variant: {
    outline:
      'shadow-none bg-ds-bg-input text-ds-fg ring-1 ring-inset ring-ds-border hover:ring-ds-brand-muted focus:ring-2 focus:ring-ds-brand/35 focus:ring-inset',
    none: 'bg-transparent focus:ring-0 focus:shadow-none',
  },
  icon: {
    base: 'flex-shrink-0 text-ds-fg-subtle',
    color: 'text-ds-fg-subtle',
    loading: 'animate-spin',
    size: {
      '2xs': 'h-3 w-3',
      xs: 'h-3.5 w-3.5',
      sm: 'h-4 w-4',
      md: 'h-4 w-4',
      lg: 'h-5 w-5',
      xl: 'h-5 w-5',
    },
    leading: {
      wrapper: 'absolute inset-y-0 start-0 flex items-center',
      pointer: 'pointer-events-none',
      padding: {
        '2xs': 'px-2',
        xs: 'px-2.5',
        sm: 'px-2.5',
        md: 'px-3',
        lg: 'px-3.5',
        xl: 'px-3.5',
      },
    },
    trailing: {
      wrapper: 'absolute inset-y-0 end-0 flex items-center',
      pointer: 'pointer-events-none',
      padding: {
        '2xs': 'px-2',
        xs: 'px-2.5',
        sm: 'px-2.5',
        md: 'px-3',
        lg: 'px-3.5',
        xl: 'px-3.5',
      },
    },
  },
  default: {
    size: 'md' as const,
    color: 'gray' as const,
    variant: 'outline' as const,
    loadingIcon: 'i-heroicons-arrow-path-20-solid',
  },
} as const
