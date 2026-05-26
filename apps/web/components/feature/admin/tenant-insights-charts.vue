<template>
  <div class="grid gap-6 lg:grid-cols-2" data-testid="tenant-insights-charts">
    <!-- Plan Distribution Donut -->
    <div class="ds-card rounded-2xl p-5 shadow-sm">
      <h3 class="mb-1 text-sm font-semibold text-ds-fg-heading">{{ $t('adminPlanDistTitle') }}</h3>
      <p class="mb-3 text-xs text-ds-fg-muted">{{ $t('adminPlanDistHint') }}</p>

      <div v-if="loading" class="flex justify-center py-8">
        <div class="h-6 w-6 animate-spin rounded-full border-2 border-ds-brand-muted border-t-ds-brand" />
      </div>
      <ChartDonut
        v-else-if="planItems.length > 0"
        :items="planItems"
        :height="260"
        data-testid="admin-plan-donut"
      />
      <p v-else class="py-8 text-center text-sm text-ds-fg-muted">{{ $t('dashboardChartEmpty') }}</p>
    </div>

    <!-- Top Tenants Bar -->
    <div class="ds-card rounded-2xl p-5 shadow-sm">
      <div class="mb-3 flex items-center justify-between">
        <div>
          <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('adminTopTenantsTitle') }}</h3>
          <p class="text-xs text-ds-fg-muted">{{ $t('adminTopTenantsHint') }}</p>
        </div>
        <select
          v-model="topMetric"
          class="ds-input rounded-lg px-2 py-1 text-xs"
          data-testid="admin-top-metric"
          @change="refreshTop"
        >
          <option value="activity">{{ $t('adminMetricActivity') }}</option>
          <option value="revenue">{{ $t('adminMetricRevenue') }}</option>
          <option value="risk">{{ $t('adminMetricRisk') }}</option>
        </select>
      </div>

      <div v-if="loading" class="flex justify-center py-8">
        <div class="h-6 w-6 animate-spin rounded-full border-2 border-ds-brand-muted border-t-ds-brand" />
      </div>
      <ChartBar
        v-else-if="topItems.length > 0"
        :items="topItems"
        :height="260"
        :horizontal="true"
        :value-formatter="(v: number) => String(v)"
        data-testid="admin-top-tenants-bar"
      />
      <p v-else class="py-8 text-center text-sm text-ds-fg-muted">{{ $t('dashboardChartEmpty') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ChartBarItem, ChartDonutItem } from '@crm/ui-kit'

const insights = useAdminTenantInsights()

const loading = ref(true)
const planItems = ref<ChartDonutItem[]>([])
const topItems = ref<ChartBarItem[]>([])
const topMetric = ref<'activity' | 'revenue' | 'risk'>('activity')

async function refresh() {
  loading.value = true
  try {
    const [plans, tops] = await Promise.all([
      insights.fetchPlanDistribution(),
      insights.fetchTopTenants({ metric: topMetric.value, limit: 10 }),
    ])
    planItems.value = plans.items.map((p) => ({ name: p.plan, value: Number(p.count) }))
    topItems.value = tops.items.map((t) => ({ name: t.tenant_name, value: Math.round(t.value) }))
  } finally {
    loading.value = false
  }
}

async function refreshTop() {
  try {
    const tops = await insights.fetchTopTenants({ metric: topMetric.value, limit: 10 })
    topItems.value = tops.items.map((t) => ({ name: t.tenant_name, value: Math.round(t.value) }))
  } catch {
    // keep existing data
  }
}

onMounted(refresh)
</script>
