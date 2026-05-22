import { bridgeUiKitThemeFromApp, UI_KIT_THEME_KEY } from '@crm/ui-kit'

/**
 * Bridge app useTheme() into @crm/ui-kit module singleton (inject 在 monorepo 分包下不可靠).
 */
export default defineNuxtPlugin((nuxtApp) => {
  const theme = useTheme()
  bridgeUiKitThemeFromApp(theme)
  nuxtApp.vueApp.provide(UI_KIT_THEME_KEY, theme)
})
