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
    </main>
  </div>
</template>

<script setup lang="ts">
import type { ChartBarItem, ChartFunnelItem, ChartLegendItem, ChartSeries } from '@crm/ui-kit'

definePageMeta({ layout: 'auth' })

const { t } = useI18n()

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
</script>
