/** 将 API 返回的 status/source 枚举码映射为当前语言的展示文案 */
const SOURCE_ALIASES: Record<string, string> = {
  web: 'website',
}

export function useLeadLabels() {
  const { t, locale } = useI18n()

  function leadStatusLabel(code: string): string {
    if (!code) return '—'
    const key = `leadStatus.${code}`
    const translated = t(key)
    return translated === key ? code : translated
  }

  function leadSourceLabel(code: string): string {
    const raw = (code || '').trim()
    if (!raw || raw === 'unknown') {
      return t('leadSource.unknown')
    }
    const normalized = SOURCE_ALIASES[raw] ?? raw
    const key = `leadSource.${normalized}`
    const translated = t(key)
    return translated === key ? raw : translated
  }

  function formatStatsDate(isoDate: string): string {
    const d = new Date(`${isoDate}T00:00:00`)
    if (Number.isNaN(d.getTime())) return isoDate
    return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
      month: 'numeric',
      day: 'numeric',
    }).format(d)
  }

  return { leadStatusLabel, leadSourceLabel, formatStatsDate, locale }
}
