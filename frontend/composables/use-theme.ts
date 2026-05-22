import { DEFAULT_THEME, THEME_COOKIE, type ThemeId } from '~/design-system/tokens'

export function useTheme() {
  const cookie = useCookie<ThemeId>(THEME_COOKIE, {
    default: () => DEFAULT_THEME,
    maxAge: 60 * 60 * 24 * 365,
    sameSite: 'lax',
  })

  const id = computed({
    get: () => (cookie.value === 'v2' ? 'v2' : 'v1') as ThemeId,
    set: (value: ThemeId) => {
      cookie.value = value
    },
  })

  const isDark = computed(() => id.value === 'v2')

  function setTheme(theme: ThemeId) {
    id.value = theme
    applyThemeToDocument(theme)
  }

  function toggleTheme() {
    setTheme(id.value === 'v1' ? 'v2' : 'v1')
  }

  return { id, isDark, setTheme, toggleTheme }
}

export function applyThemeToDocument(theme: ThemeId) {
  if (import.meta.client) {
    document.documentElement.setAttribute('data-theme', theme)
  }
}
