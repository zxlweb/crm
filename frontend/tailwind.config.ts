import type { Config } from 'tailwindcss'

export default {
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Plus Jakarta Sans"', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
      colors: {
        ds: {
          bg: 'var(--ds-bg)',
          'bg-elevated': 'var(--ds-bg-elevated)',
          'bg-muted': 'var(--ds-bg-muted)',
          'bg-input': 'var(--ds-bg-input)',
          'bg-sidebar': 'var(--ds-bg-sidebar)',
          'bg-topbar': 'var(--ds-bg-topbar)',
          'brand-subtle': 'var(--ds-bg-brand-subtle)',
          'success-subtle': 'var(--ds-bg-success-subtle)',
          'danger-subtle': 'var(--ds-bg-danger-subtle)',
          fg: 'var(--ds-fg)',
          'fg-heading': 'var(--ds-fg-heading)',
          'fg-muted': 'var(--ds-fg-muted)',
          'fg-subtle': 'var(--ds-fg-subtle)',
          'fg-brand': 'var(--ds-fg-brand)',
          'fg-nav': 'var(--ds-fg-nav)',
          'fg-nav-active': 'var(--ds-fg-nav-active)',
          'on-brand': 'var(--ds-fg-on-brand)',
          border: 'var(--ds-border)',
          'border-muted': 'var(--ds-border-muted)',
          brand: 'var(--ds-brand)',
          'brand-hover': 'var(--ds-brand-hover)',
          'brand-muted': 'var(--ds-brand-muted)',
          success: 'var(--ds-success)',
          danger: 'var(--ds-danger)',
        },
      },
      boxShadow: {
        'ds-brand': 'var(--ds-brand-shadow)',
      },
    },
  },
} satisfies Config
