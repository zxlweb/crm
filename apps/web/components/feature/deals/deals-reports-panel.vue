<template>
  <div class="space-y-6" data-testid="deals-reports-panel">
    <p class="text-sm text-ds-fg-muted">{{ $t('dealsReportsHint') }}</p>

    <UAlert v-if="loadError" color="red" variant="soft" :title="loadError" />

    <section class="grid gap-6 lg:grid-cols-2">
      <CardShell :title="$t('dealsReportByStage')" :subtitle="$t('dealsReportByStageHint')" :height="300">
        <LeadsReportChartSlot :pending="pending" :empty="!stageBars.length" :height="260" :empty-text="$t('dealsPipelineEmpty')">
          <ChartBar :key="`stage-${chartLocale}`" :items="stageBars" :horizontal="false" :height="240" />
        </LeadsReportChartSlot>
      </CardShell>

      <CardShell :title="$t('dealsReportWinRate')" :subtitle="$t('dealsReportWinRateHint')" :height="300">
        <LeadsReportChartSlot :pending="pending" :empty="!winRateCategories.length" :height="260" :empty-text="$t('dealsPipelineEmpty')">
          <ChartLine
            :key="`win-${chartLocale}`"
            :categories="winRateCategories"
            :series="winRateSeries"
            :height="240"
            :y-formatter="formatWinRate"
            :show-area="false"
            :always-show-symbols="winRateSparse"
          />
        </LeadsReportChartSlot>
        <p v-if="winRateSparse && winRateCategories.length" class="mt-2 text-center text-xs text-ds-fg-subtle">
          {{ $t('dealsReportWinRateSparseHint') }}
        </p>
      </CardShell>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { ChartBarItem, ChartSeries } from '@crm/ui-kit'
import type { DealStatsByStage, DealStatsWinRate } from '~/types/deal-stats'
import { defaultDealStatsRange } from '~/composables/use-deal-stats'

const { t, locale } = useI18n()
const { dealStageLabel } = useDealLabels()
const statsApi = useDealStats()

const chartLocale = computed(() => locale.value)

const pending = ref(true)
const loadError = ref('')
const byStage = ref<DealStatsByStage | null>(null)
const winRate = ref<DealStatsWinRate | null>(null)

const stageBars = computed<ChartBarItem[]>(() => {
  void locale.value
  return (byStage.value?.items ?? []).map((item) => ({
    name: dealStageLabel(item.label),
    value: item.value,
  }))
})

const winRateCategories = computed(() => winRate.value?.items.map((i) => i.period) ?? [])

const winRateSparse = computed(() => (winRate.value?.items.length ?? 0) <= 2)

const winRateSeries = computed<ChartSeries[]>(() => {
  void locale.value
  if (!winRate.value?.items.length) return []
  return [
    {
      name: t('dealsReportWinRateSeries'),
      data: winRate.value.items.map((i) => Math.round(i.rate * 100)),
      primary: true,
    },
  ]
})

function formatWinRate(value: number) {
  return `${value}%`
}

async function load() {
  pending.value = true
  loadError.value = ''
  const range = defaultDealStatsRange()
  try {
    const [stageRes, rateRes] = await Promise.all([
      statsApi.fetchByStage({ ...range, metric: 'count' }),
      statsApi.fetchWinRate({ ...range, granularity: 'week' }),
    ])
    byStage.value = stageRes
    winRate.value = rateRes
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

onMounted(load)
defineExpose({ reload: load })
</script>
