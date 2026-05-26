<template>
  <section class="space-y-4" data-testid="dashboard-analytics-panel">
    <div class="flex flex-wrap items-end justify-between gap-2 px-1">
      <div class="min-w-0">
        <h2 class="text-base font-semibold tracking-tight text-ds-fg-heading sm:text-lg">
          {{ $t('dashboardAnalyticsTitle') }}
        </h2>
        <p class="mt-0.5 text-xs text-ds-fg-muted">{{ $t('dashboardAnalyticsHint') }}</p>
      </div>
      <NuxtLink
        to="/deals"
        class="inline-flex cursor-pointer items-center gap-1 text-xs font-semibold text-ds-fg-brand transition-colors duration-200 hover:text-ds-brand-hover"
      >
        {{ $t('dashboardAnalyticsDealsLink') }}
        <UIcon name="i-heroicons-arrow-right-20-solid" class="h-3 w-3" aria-hidden="true" />
      </NuxtLink>
    </div>

    <UAlert v-if="loadError" color="red" variant="soft" :title="loadError" />

    <div
      class="grid gap-4"
      :class="showTeamRanking ? 'xl:grid-cols-3' : 'xl:grid-cols-2'"
    >
      <article
        class="ds-analytics-card relative overflow-hidden rounded-2xl border border-ds-border-muted bg-ds-bg-elevated shadow-ds-sm"
        :data-featured="true"
      >
        <span
          class="pointer-events-none absolute inset-x-0 top-0 h-0.5"
          :style="{ background: 'var(--ds-brand-gradient)' }"
          aria-hidden="true"
        />
        <span
          class="pointer-events-none absolute -right-10 -top-10 h-36 w-36 rounded-full opacity-15 blur-3xl"
          :style="{ background: 'var(--ds-brand-gradient)' }"
          aria-hidden="true"
        />
        <header class="flex items-start justify-between gap-3 px-4 pb-2 pt-4 sm:px-5">
          <div class="min-w-0">
            <h3 class="text-sm font-semibold text-ds-fg-heading">
              {{ $t('dashboardQuotaTitle') }}
            </h3>
            <p class="mt-0.5 text-xs text-ds-fg-muted">{{ quotaSubtitle }}</p>
          </div>
          <span
            v-if="quotaHasTarget"
            class="inline-flex shrink-0 items-center gap-1 rounded-full bg-ds-brand-subtle px-2 py-0.5 text-[11px] font-semibold text-ds-fg-brand ring-1 ring-inset ring-ds-brand/25"
          >
            {{ quotaPercent }}%
          </span>
        </header>
        <div class="px-3 pb-3">
          <LeadsReportChartSlot
            :pending="pending"
            :empty="quota == null"
            :height="220"
            :empty-text="$t('dashboardChartEmpty')"
          >
            <div v-if="quotaHasTarget" data-testid="dashboard-quota-gauge" class="flex justify-center">
              <ChartGauge
                :value="quotaPercent"
                :label="$t('chartsGaugeLabel')"
                :height="220"
              />
            </div>
            <div
              v-else-if="quota"
              data-testid="dashboard-quota-gauge"
              class="flex h-[220px] flex-col items-center justify-center gap-1.5 text-center"
            >
              <p class="text-3xl font-bold tabular-nums text-ds-fg-heading">
                <span class="bg-clip-text" :style="brandGradientText">
                  {{ formatDealAmount(quota.won_amount_mtd) }}
                </span>
              </p>
              <p class="text-sm text-ds-fg-muted">{{ $t('dashboardQuotaWonMtdLabel') }}</p>
              <p class="text-xs text-ds-fg-subtle">{{ $t('dashboardQuotaNoTargetHint') }}</p>
            </div>
          </LeadsReportChartSlot>
        </div>
      </article>

      <article
        class="ds-analytics-card relative overflow-hidden rounded-2xl border border-ds-border-muted bg-ds-bg-elevated shadow-ds-sm"
      >
        <span
          class="pointer-events-none absolute inset-x-0 top-0 h-0.5 bg-ds-info opacity-80"
          aria-hidden="true"
        />
        <header class="flex items-start justify-between gap-3 px-4 pb-2 pt-4 sm:px-5">
          <div class="min-w-0">
            <h3 class="text-sm font-semibold text-ds-fg-heading">
              {{ $t('dashboardDealsFunnelTitle') }}
            </h3>
            <p class="mt-0.5 text-xs text-ds-fg-muted">
              {{ $t('dashboardDealsFunnelHint') }}
            </p>
          </div>
          <span
            v-if="funnelItems.length > 0"
            class="inline-flex shrink-0 items-center gap-1 rounded-full bg-ds-info-subtle px-2 py-0.5 text-[11px] font-semibold text-ds-info ring-1 ring-inset ring-ds-info/25"
          >
            {{ funnelItems.length }}
          </span>
        </header>
        <div class="px-3 pb-3">
          <LeadsReportChartSlot
            :pending="pending"
            :empty="!funnelItems.length"
            :height="220"
            :empty-text="$t('dashboardChartEmpty')"
          >
            <div data-testid="dashboard-deals-funnel">
              <ChartFunnel :key="chartLocale" :items="funnelItems" :height="220" />
            </div>
          </LeadsReportChartSlot>
        </div>
      </article>

      <article
        v-if="showTeamRanking"
        class="ds-analytics-card relative overflow-hidden rounded-2xl border border-ds-border-muted bg-ds-bg-elevated shadow-ds-sm"
      >
        <span
          class="pointer-events-none absolute inset-x-0 top-0 h-0.5 bg-ds-success opacity-80"
          aria-hidden="true"
        />
        <header class="flex items-start justify-between gap-3 px-4 pb-2 pt-4 sm:px-5">
          <div class="min-w-0">
            <h3 class="text-sm font-semibold text-ds-fg-heading">
              {{ $t('dashboardTeamRankingTitle') }}
            </h3>
            <p class="mt-0.5 text-xs text-ds-fg-muted">
              {{ $t('dashboardTeamRankingHint') }}
            </p>
          </div>
          <span
            v-if="rankingBars.length > 0"
            class="inline-flex shrink-0 items-center gap-1 rounded-full bg-ds-success-subtle px-2 py-0.5 text-[11px] font-semibold text-ds-success ring-1 ring-inset ring-ds-success/25"
          >
            {{ rankingBars.length }}
          </span>
        </header>
        <div class="px-3 pb-3">
          <LeadsReportChartSlot
            :pending="pending"
            :empty="!rankingBars.length"
            :height="220"
            :empty-text="$t('dashboardChartEmpty')"
          >
            <div data-testid="dashboard-team-ranking">
              <ChartBar :key="`rank-${chartLocale}`" :items="rankingBars" horizontal :height="220" />
            </div>
          </LeadsReportChartSlot>
        </div>
      </article>
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

const brandGradientText = {
  background: 'var(--ds-brand-gradient)',
  '-webkit-background-clip': 'text',
  'background-clip': 'text',
  color: 'transparent',
}

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
