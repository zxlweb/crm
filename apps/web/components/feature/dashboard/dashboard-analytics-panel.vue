<template>
  <section class="space-y-6" data-testid="dashboard-analytics-panel">
    <div class="flex flex-wrap items-end justify-between gap-2">
      <div>
        <h2 class="text-lg font-semibold text-ds-fg-heading">{{ $t('dashboardAnalyticsTitle') }}</h2>
        <p class="mt-0.5 text-sm text-ds-fg-muted">{{ $t('dashboardAnalyticsHint') }}</p>
      </div>
      <NuxtLink
        to="/deals"
        class="text-xs font-semibold text-ds-fg-brand hover:underline"
      >
        {{ $t('dashboardAnalyticsDealsLink') }}
      </NuxtLink>
    </div>

    <UAlert v-if="loadError" color="red" variant="soft" :title="loadError" />

    <div
      class="grid gap-6"
      :class="showTeamRanking ? 'xl:grid-cols-3' : 'xl:grid-cols-2'"
    >
      <CardShell :title="$t('dashboardQuotaTitle')" :subtitle="quotaSubtitle" :height="300">
        <LeadsReportChartSlot :pending="pending" :empty="quota == null" :height="260" :empty-text="$t('dashboardChartEmpty')">
          <div v-if="quotaHasTarget" data-testid="dashboard-quota-gauge" class="flex justify-center">
            <ChartGauge
              :value="quotaPercent"
              :label="$t('chartsGaugeLabel')"
              :height="240"
            />
          </div>
          <div
            v-else-if="quota"
            data-testid="dashboard-quota-gauge"
            class="flex h-[240px] flex-col items-center justify-center gap-2 text-center"
          >
            <p class="text-3xl font-bold tabular-nums text-ds-fg-heading">
              {{ formatDealAmount(quota.won_amount_mtd) }}
            </p>
            <p class="text-sm text-ds-fg-muted">{{ $t('dashboardQuotaWonMtdLabel') }}</p>
            <p class="text-xs text-ds-fg-subtle">{{ $t('dashboardQuotaNoTargetHint') }}</p>
          </div>
        </LeadsReportChartSlot>
      </CardShell>

      <CardShell :title="$t('dashboardDealsFunnelTitle')" :subtitle="$t('dashboardDealsFunnelHint')" :height="300">
        <LeadsReportChartSlot :pending="pending" :empty="!funnelItems.length" :height="260" :empty-text="$t('dashboardChartEmpty')">
          <div data-testid="dashboard-deals-funnel">
            <ChartFunnel :key="chartLocale" :items="funnelItems" :height="240" />
          </div>
        </LeadsReportChartSlot>
      </CardShell>

      <CardShell
        v-if="showTeamRanking"
        :title="$t('dashboardTeamRankingTitle')"
        :subtitle="$t('dashboardTeamRankingHint')"
        :height="300"
      >
        <LeadsReportChartSlot :pending="pending" :empty="!rankingBars.length" :height="260" :empty-text="$t('dashboardChartEmpty')">
          <div data-testid="dashboard-team-ranking">
            <ChartBar :key="`rank-${chartLocale}`" :items="rankingBars" horizontal :height="240" />
          </div>
        </LeadsReportChartSlot>
      </CardShell>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { ChartBarItem, ChartFunnelItem } from '@crm/ui-kit'
import type { DashboardFunnel, DashboardQuota, DashboardTeamRanking } from '~/types/dashboard-stats'

const props = withDefaults(
  defineProps<{
    /** department / all 数据范围，或租户管理员 */
    showTeamRanking?: boolean
  }>(),
  { showTeamRanking: false },
)

const { t, locale } = useI18n()
const { dealStageLabel, formatDealAmount } = useDealLabels()
const dashboardStats = useDashboardStats()

const pending = ref(true)
const loadError = ref('')
const quota = ref<DashboardQuota | null>(null)
const funnel = ref<DashboardFunnel | null>(null)
const ranking = ref<DashboardTeamRanking | null>(null)

const chartLocale = computed(() => locale.value)
const showTeamRanking = computed(() => props.showTeamRanking)

const quotaHasTarget = computed(() => (quota.value?.target_amount ?? 0) > 0)

const quotaPercent = computed(() =>
  quota.value ? Math.round(quota.value.completion_rate * 100) : 0,
)

const quotaSubtitle = computed(() => {
  if (!quota.value) return t('dashboardQuotaHint')
  if (!quotaHasTarget.value) {
    return t('dashboardQuotaNoTargetSubtitle')
  }
  return t('dashboardQuotaPeriodHint', {
    period: quota.value.period,
    target: formatDealAmount(quota.value.target_amount),
  })
})

const funnelItems = computed<ChartFunnelItem[]>(() => {
  void locale.value
  return (funnel.value?.stages ?? []).map((s) => ({
    name: dealStageLabel(s.name),
    value: s.count,
  }))
})

const rankingBars = computed<ChartBarItem[]>(() => {
  void locale.value
  return (ranking.value?.items ?? []).map((row) => ({
    name: row.name,
    value: row.value,
  }))
})

async function load() {
  pending.value = true
  loadError.value = ''
  try {
    const [quotaRes, funnelRes, rankingRes] = await Promise.all([
      dashboardStats.fetchQuota(),
      dashboardStats.fetchFunnel('deals'),
      showTeamRanking.value
        ? dashboardStats.fetchTeamRankingOptional('won_amount')
        : Promise.resolve(null),
    ])
    quota.value = quotaRes
    funnel.value = funnelRes
    ranking.value = rankingRes
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

onMounted(load)

watch(showTeamRanking, () => {
  void load()
})
</script>
