<template>
  <div class="min-h-screen bg-ds-bg">
    <header class="sticky top-0 z-10 border-b border-ds-border bg-ds-bg-topbar px-6 py-4 backdrop-blur-md">
      <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-4">
        <div>
          <p class="text-xs font-medium text-ds-fg-brand">{{ $t('chartsPageLabel') }}</p>
          <h1 class="text-xl font-bold text-ds-fg-heading">{{ $t('chartsPageTitle') }}</h1>
          <p class="mt-1 text-sm text-ds-fg-muted">{{ $t('chartsPageDesc') }}</p>
        </div>
        <div class="flex items-center gap-3">
          <UiThemeToggle />
          <NuxtLink to="/cards" class="text-sm text-ds-fg-muted transition-colors hover:text-ds-fg-brand">
            {{ $t('cardsPageTitle') }}
          </NuxtLink>
          <NuxtLink to="/design" class="text-sm text-ds-fg-muted transition-colors hover:text-ds-fg-brand">
            {{ $t('themePreviewTitle') }}
          </NuxtLink>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-6xl space-y-8 px-6 py-8">
      <!-- 折线图 · Sales Overview 风格 -->
      <section>
        <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">{{ $t('chartsSectionLine') }}</h2>
        <ChartShell
          :title="$t('chartsSalesOverview')"
          :metric="salesMetric"
          :metric-label="$t('chartsAverageSales')"
          :legend="lineLegend"
          :height="300"
        >
          <ChartLine
            :categories="salesDays"
            :series="salesSeries"
            :y-formatter="formatCurrencyK"
            :height="280"
          />
        </ChartShell>
      </section>

      <section class="grid gap-6 lg:grid-cols-2">
        <!-- 柱状图 -->
        <div>
          <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">{{ $t('chartsSectionBar') }}</h2>
          <ChartShell
            :title="$t('chartsProductPerformance')"
            :subtitle="$t('chartsProductPerformanceDesc')"
            :height="320"
          >
            <ChartBar
              :items="productBars"
              horizontal
              :height="300"
            />
          </ChartShell>
        </div>

        <!-- 月度柱状 -->
        <div>
          <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">{{ $t('chartsSectionBarVertical') }}</h2>
          <ChartShell
            :title="$t('chartsMonthlyRevenue')"
            :metric="monthlyTotal"
            :height="320"
          >
            <ChartBar
              :items="monthlyBars"
              :horizontal="false"
              :value-formatter="formatCurrencyK"
              :height="300"
            />
          </ChartShell>
        </div>
      </section>

      <!-- Activity 跟进类型（纵向柱，与线索详情摘要一致） -->
      <section>
        <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">{{ $t('chartsSectionActivity') }}</h2>
        <ChartShell
          :title="$t('chartsActivityTypeTitle')"
          :subtitle="$t('chartsActivityTypeDesc')"
          :height="300"
        >
          <ChartBar :items="activityTypeBars" :horizontal="false" :height="260" />
        </ChartShell>
      </section>

      <!-- 环形图 -->
      <section>
        <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">{{ $t('chartsSectionDonut') }}</h2>
        <ChartShell
          :title="$t('chartsDonutTitle')"
          :subtitle="$t('chartsDonutDesc')"
          :height="320"
        >
          <ChartDonut :items="sourceDonut" :height="300" />
        </ChartShell>
      </section>

      <!-- 漏斗图 -->
      <section>
        <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">{{ $t('chartsSectionFunnel') }}</h2>
        <ChartShell
          :title="$t('chartsPipelineFunnel')"
          :subtitle="$t('chartsPipelineFunnelDesc')"
          :metric="funnelTop"
          :metric-label="$t('chartsLeadsTotal')"
          :height="360"
        >
          <ChartFunnel :items="pipelineFunnel" :height="340" />
        </ChartShell>
      </section>

      <!-- 双折线对比 -->
      <section>
        <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">{{ $t('chartsSectionCompare') }}</h2>
        <ChartShell
          :title="$t('chartsTenantTrend')"
          :legend="tenantLegend"
          :height="300"
        >
          <ChartLine
            :categories="tenantDays"
            :series="tenantSeries"
            :show-area="false"
            :height="280"
          />
        </ChartShell>
      </section>

      <!-- Sparkline · KPI 内嵌趋势 -->
      <section>
        <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">{{ $t('chartsSectionSparkline') }}</h2>
        <div class="grid gap-6 lg:grid-cols-2">
          <ChartShell :title="$t('chartsSparklineTitle')" :subtitle="$t('chartsSparklineDesc')" :height="120">
            <div class="flex items-center justify-between gap-4 px-2">
              <div>
                <p class="text-xs text-ds-fg-muted">{{ $t('chartsSparklineMetricLabel') }}</p>
                <p class="text-2xl font-bold tabular-nums text-ds-fg-heading">¥2.4M</p>
              </div>
              <ChartSparkline :values="sparklineUp" tone="up" :width="120" :height="40" />
            </div>
          </ChartShell>
          <ChartShell :title="$t('chartsSparklineFlatTitle')" :height="120">
            <div class="flex items-center justify-end gap-4 px-2">
              <ChartSparkline :values="sparklineMixed" tone="auto" :width="140" :height="40" />
            </div>
          </ChartShell>
        </div>
      </section>

      <!-- Gauge · 配额完成率 -->
      <section>
        <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">{{ $t('chartsSectionGauge') }}</h2>
        <ChartShell
          :title="$t('chartsGaugeTitle')"
          :subtitle="$t('chartsGaugeDesc')"
          :height="280"
          data-testid="dashboard-quota-gauge"
        >
          <ChartGauge :value="72" :label="$t('chartsGaugeLabel')" :height="240" />
        </ChartShell>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import type { ChartBarItem, ChartDonutItem, ChartFunnelItem, ChartLegendItem, ChartSeries } from '@crm/ui-kit'

definePageMeta({ layout: 'auth' })

const { t } = useI18n()
const { activityTypeLabel } = useActivityLabels()

const salesDays = Array.from({ length: 15 }, (_, i) => String(i + 1))
const salesCurrent = [8, 10, 9, 12, 11, 14, 16, 18, 17, 20, 19, 22, 21, 20, 24]
const salesBefore = salesCurrent.map((v) => Math.round(v * 0.72))

const salesSeries = computed<ChartSeries[]>(() => [
  { name: t('chartsSeriesCurrent'), data: salesCurrent, primary: true },
  { name: t('chartsSeriesBefore'), data: salesBefore, compare: true },
])

const lineLegend = computed<ChartLegendItem[]>(() => [
  { label: t('chartsSeriesCurrent') },
  { label: t('chartsSeriesBefore'), muted: true, dashed: true },
])

const salesMetric = '$88,692.00'

function formatCurrencyK(value: number) {
  return `$${value}k`
}

const sourceDonut = computed<ChartDonutItem[]>(() => [
  { name: t('leadSource.website'), value: 120 },
  { name: t('leadSource.exhibition'), value: 80 },
  { name: t('leadSource.referral'), value: 45 },
  { name: t('leadSource.partner'), value: 30 },
])

const activityTypeBars = computed<ChartBarItem[]>(() => [
  { name: activityTypeLabel('call'), value: 18 },
  { name: activityTypeLabel('email'), value: 12 },
  { name: activityTypeLabel('meeting'), value: 8 },
  { name: activityTypeLabel('wechat'), value: 6 },
  { name: activityTypeLabel('note'), value: 4 },
])

const productBars: ChartBarItem[] = [
  { name: 'Baby & Kids', value: 86 },
  { name: 'Personal Care', value: 72 },
  { name: 'Electronics', value: 58 },
  { name: 'Home & Living', value: 45 },
  { name: 'Fashion', value: 38 },
]

const monthlyBars: ChartBarItem[] = [
  { name: 'Jan', value: 12 },
  { name: 'Feb', value: 15 },
  { name: 'Mar', value: 14 },
  { name: 'Apr', value: 18 },
  { name: 'May', value: 20 },
  { name: 'Jun', value: 22 },
]

const monthlyTotal = '$156k'

const pipelineFunnel = computed<ChartFunnelItem[]>(() => [
  { name: t('chartsFunnelLead'), value: 1200 },
  { name: t('chartsFunnelMql'), value: 680 },
  { name: t('chartsFunnelSql'), value: 320 },
  { name: t('chartsFunnelOpp'), value: 140 },
  { name: t('chartsFunnelWon'), value: 48 },
])

const funnelTop = '1,200'

const tenantDays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
const tenantSeries: ChartSeries[] = [
  { name: 'Active', data: [4, 5, 4, 6, 5, 7, 6], primary: true },
  { name: 'New', data: [1, 2, 1, 3, 2, 2, 3], compare: true },
]

const tenantLegend: ChartLegendItem[] = [
  { label: 'Active' },
  { label: 'New', muted: true, dashed: true },
]

const sparklineUp = [12, 14, 13, 16, 18, 17, 22, 24]
const sparklineMixed = [8, 12, 9, 11, 10, 13, 11, 15]
</script>
