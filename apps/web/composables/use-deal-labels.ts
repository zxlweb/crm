import type { DealCurrency, DealStage } from '~/types/deal'

export function useDealLabels() {
  const { t, locale } = useI18n()

  function dealStageLabel(stage: DealStage | string): string {
    if (!stage) return '—'
    const key = `dealStage.${stage}`
    const translated = t(key)
    return translated === key ? stage : translated
  }

  function formatDealAmount(amount: number, currency: DealCurrency = 'CNY'): string {
    const code = currency === 'USD' ? 'USD' : 'CNY'
    return new Intl.NumberFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
      style: 'currency',
      currency: code,
      maximumFractionDigits: 0,
    }).format(amount)
  }

  function formatCloseDate(iso: string | null | undefined): string {
    if (!iso) return '—'
    const d = new Date(`${iso}T00:00:00`)
    if (Number.isNaN(d.getTime())) return iso
    return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
      month: 'short',
      day: 'numeric',
    }).format(d)
  }

  return { dealStageLabel, formatDealAmount, formatCloseDate, locale }
}
