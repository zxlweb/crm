<template>
  <div class="space-y-6" data-testid="leads-reports-panel">
    <p class="text-sm text-ds-fg-muted">{{ $t('leadsReportsHint') }}</p>

    <UAlert v-if="loadError" color="red" variant="soft" :title="loadError" />

    <section class="grid gap-6 lg:grid-cols-2">
      <CardShell :title="$t('leadsReportStatus')" :subtitle="$t('leadsReportStatusHint')" :height="280">
        <LeadsReportChartSlot :pending="pending" :empty="!statusBars.length" :height="240">
          <ChartBar :key="`status-${chartLocale}`" :items="statusBars" :height="240" :horizontal="false" />
        </LeadsReportChartSlot>
      </CardShell>
      <CardShell :title="$t('leadsReportSource')" :subtitle="$t('leadsReportSourceHint')" :height="280">
        <LeadsReportChartSlot :pending="pending" :empty="!sourceDonut.length" :height="240">
          <ChartDonut :key="`source-${chartLocale}`" :items="sourceDonut" :height="240" />
        </LeadsReportChartSlot>
      </CardShell>
    </section>

    <section class="grid gap-6 lg:grid-cols-2">
      <CardShell :title="$t('leadsReportTrend')" :subtitle="$t('leadsReportTrendHint')" :height="300">
        <LeadsReportChartSlot :pending="pending" :empty="!trendCategories.length" :height="260">
          <ChartLine
            :key="`trend-${chartLocale}`"
            :categories="trendCategories"
            :series="trendSeries"
            :height="260"
            :show-area="true"
          />
        </LeadsReportChartSlot>
      </CardShell>
      <CardShell :title="$t('leadsReportFunnel')" :subtitle="$t('leadsReportFunnelHint')" :height="300">
        <LeadsReportChartSlot :pending="pending" :empty="!funnelItems.length" :height="260">
          <ChartFunnel :key="`funnel-${chartLocale}`" :items="funnelItems" :height="260" />
        </LeadsReportChartSlot>
      </CardShell>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { ChartBarItem, ChartDonutItem, ChartFunnelItem, ChartSeries } from '@crm/ui-kit'
import type { LeadStatsDistribution, LeadStatsFunnel, LeadStatsTrend } from '~/types/lead-stats'
import { defaultLeadStatsRange } from '~/composables/use-leads-stats'

const { t, locale } = useI18n()
const { leadStatusLabel, leadSourceLabel, formatStatsDate } = useLeadLabels()
const statsApi = useLeadsStats()

const chartLocale = computed(() => locale.value)

const pending = ref(true)
const loadError = ref('')
const byStatus = ref<LeadStatsDistribution | null>(null)
const bySource = ref<LeadStatsDistribution | null>(null)
const trend = ref<LeadStatsTrend | null>(null)
const funnel = ref<LeadStatsFunnel | null>(null)

const statusBars = computed<ChartBarItem[]>(() => {
  void locale.value
  return (byStatus.value?.items ?? []).map((item) => ({
    name: leadStatusLabel(item.label),
    value: item.value,
  }))
})

const sourceDonut = computed<ChartDonutItem[]>(() => {
  void locale.value
  return (bySource.value?.items ?? []).map((item) => ({
    name: leadSourceLabel(item.label),
    value: item.value,
  }))
})

const trendCategories = computed(() => {
  void locale.value
  return (trend.value?.categories ?? []).map((d) => formatStatsDate(d))
})

const trendSeries = computed<ChartSeries[]>(() => {
  void locale.value
  if (!trend.value?.series.length) return []
  return trend.value.series.map((s) => ({
    name: t('leadsReportTrendSeries'),
    data: s.data,
    primary: s.primary ?? true,
  }))
})

const funnelItems = computed<ChartFunnelItem[]>(() => {
  void locale.value
  return (funnel.value?.stages ?? []).map((stage) => ({
    name: leadStatusLabel(stage.name),
    value: stage.count,
  }))
})

async function load() {
  pending.value = true
  loadError.value = ''
  const range = defaultLeadStatsRange()
  try {
    const [statusRes, sourceRes, trendRes, funnelRes] = await Promise.all([
      statsApi.fetchByStatus(range),
      statsApi.fetchBySource(range),
      statsApi.fetchTrend({ ...range, granularity: 'day' }),
      statsApi.fetchFunnel(range),
    ])
    byStatus.value = statusRes
    bySource.value = sourceRes
    trend.value = trendRes
    funnel.value = funnelRes
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

onMounted(load)
defineExpose({ reload: load })
</script>
