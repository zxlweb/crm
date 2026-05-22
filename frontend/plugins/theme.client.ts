import { applyThemeToDocument } from '~/composables/use-theme'

export default defineNuxtPlugin(() => {
  const { id } = useTheme()

  applyThemeToDocument(id.value)

  watch(id, (theme) => {
    applyThemeToDocument(theme)
  })
})
