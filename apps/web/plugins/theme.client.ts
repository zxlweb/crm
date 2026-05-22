import { applyThemeToDocument } from '@crm/ui-kit'

export default defineNuxtPlugin(() => {
  const { id } = useTheme()

  applyThemeToDocument(id.value)

  watch(id, (theme) => {
    applyThemeToDocument(theme)
  })
})
