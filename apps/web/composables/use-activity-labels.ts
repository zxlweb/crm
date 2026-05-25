/** 将 API 返回的 event_type 枚举码映射为当前语言的展示文案 */
export function useActivityLabels() {
  const { t, locale } = useI18n()

  function activityTypeLabel(code: string): string {
    if (!code) return '—'
    const key = `activityType.${code}`
    const translated = t(key)
    return translated === key ? code : translated
  }

  return { activityTypeLabel, locale }
}
