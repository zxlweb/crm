import type { ThemeId } from '~/design-system/tokens'

/** 支持 ?theme=v1|v2 预览切换 */
export default defineNuxtRouteMiddleware((to) => {
  const q = to.query.theme
  if (q !== 'v1' && q !== 'v2') return

  const theme = useTheme()
  theme.setTheme(q as ThemeId)
})
