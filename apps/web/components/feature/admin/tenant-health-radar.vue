<template>
  <div class="ds-card rounded-2xl p-5 shadow-sm" data-testid="tenant-health-radar">
    <h3 class="mb-1 text-sm font-semibold text-ds-fg-heading">{{ $t('adminHealthRadarTitle') }}</h3>
    <p class="mb-3 text-xs text-ds-fg-muted">{{ $t('adminHealthRadarHint') }}</p>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="h-6 w-6 animate-spin rounded-full border-2 border-ds-brand-muted border-t-ds-brand" />
    </div>

    <template v-else-if="healthData">
      <ChartRadar
        :indicators="indicators"
        :series="radarSeries"
        :height="320"
        data-testid="admin-radar-chart"
      />

      <div class="mt-4 overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="border-b border-ds-border-muted text-xs uppercase text-ds-fg-brand">
            <tr>
              <th class="px-3 py-2">{{ $t('name') }}</th>
              <th class="px-3 py-2 text-right">{{ $t('adminHealthOverall') }}</th>
              <th v-for="dim in healthData.dimensions" :key="dim" class="px-3 py-2 text-right">
                {{ $t(`adminHealthDim.${dim}`) }}
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-ds-border-muted">
            <tr v-for="item in healthData.items" :key="item.tenant_id" class="transition-colors hover:bg-ds-bg-muted">
              <td class="px-3 py-2 font-medium text-ds-fg-heading">{{ item.tenant_name }}</td>
              <td class="px-3 py-2 text-right">
                <span
                  class="rounded-full px-2 py-0.5 text-xs font-medium"
                  :class="scoreColor(item.overall_score)"
                >
                  {{ item.overall_score }}
                </span>
              </td>
              <td v-for="dim in healthData.dimensions" :key="dim" class="px-3 py-2 text-right text-ds-fg-muted">
                {{ item.scores[dim] }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <p v-else class="py-8 text-center text-sm text-ds-fg-muted">{{ $t('dashboardChartEmpty') }}</p>
  </div>
</template>

<script setup lang="ts">
import type { ChartRadarItem } from '@crm/ui-kit'
import type { TenantHealthData } from '~/types/admin-insights'

const { t } = useI18n()
const insights = useAdminTenantInsights()

const loading = ref(true)
const healthData = ref<TenantHealthData | null>(null)

const indicators = computed(() => {
  if (!healthData.value) return []
  return healthData.value.dimensions.map((dim) => ({
    name: t(`adminHealthDim.${dim}`),
    max: 100,
  }))
})

const radarSeries = computed<ChartRadarItem[]>(() => {
  if (!healthData.value) return []
  return healthData.value.items.map((item) => ({
    name: item.tenant_name,
    values: healthData.value!.dimensions.map((dim) => item.scores[dim]),
  }))
})

function scoreColor(score: number): string {
  if (score >= 70) return 'bg-ds-success-subtle text-ds-success'
  if (score >= 50) return 'bg-yellow-100 text-yellow-700'
  return 'bg-ds-danger-subtle text-ds-danger'
}

async function load() {
  loading.value = true
  try {
    healthData.value = await insights.fetchTenantHealth()
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
