<template>
  <ChartShell :title="$t('adminChartTitle')" :subtitle="$t('adminChartSubtitle')" :height="300">
    <template #header-extra>
      <span class="rounded-lg bg-ds-brand-subtle px-2.5 py-1 text-xs font-medium text-ds-fg-brand">{{ $t('adminChartPeriod') }}</span>
    </template>
    <div v-if="pending" class="flex h-[260px] items-center justify-center">
      <div class="h-8 w-8 animate-spin rounded-full border-2 border-ds-brand-muted border-t-ds-brand" />
    </div>
    <p v-else-if="loadError" class="flex h-[260px] items-center justify-center text-sm text-ds-danger">{{ loadError }}</p>
    <ChartLine
      v-else
      :categories="categories"
      :series="chartSeries"
      :height="260"
      :show-area="true"
    />
  </ChartShell>
</template>

<script setup lang="ts">
import type { ChartSeries } from '@crm/ui-kit'
import type { TenantActivityTrend } from '~/composables/use-super-admin'

const superAdmin = useSuperAdmin()
const { t } = useI18n()

const pending = ref(true)
const loadError = ref('')
const trend = ref<TenantActivityTrend | null>(null)

const categories = computed(() => trend.value?.categories ?? [])

const chartSeries = computed<ChartSeries[]>(() => {
  if (!trend.value?.series.length) return []
  return trend.value.series.map((s) => ({
    name: s.name === 'logins' ? t('adminChartLogins') : t('adminChartNewTenants'),
    data: s.data,
    primary: s.primary,
    compare: !s.primary,
  }))
})

async function load() {
  pending.value = true
  loadError.value = ''
  try {
    trend.value = await superAdmin.fetchTenantActivityTrend(7)
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

onMounted(load)
defineExpose({ reload: load })
</script>
