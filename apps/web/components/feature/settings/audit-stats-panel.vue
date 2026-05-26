<template>
  <div class="space-y-6" data-testid="audit-stats-panel">
    <!-- Filters -->
    <div class="flex flex-wrap items-end gap-4">
      <div>
        <label class="mb-1 block text-xs font-medium text-ds-fg-muted">{{ $t('auditFrom') }}</label>
        <input v-model="filters.from" type="date" class="ds-input rounded-xl px-3 py-2 text-sm" data-testid="audit-filter-from">
      </div>
      <div>
        <label class="mb-1 block text-xs font-medium text-ds-fg-muted">{{ $t('auditTo') }}</label>
        <input v-model="filters.to" type="date" class="ds-input rounded-xl px-3 py-2 text-sm" data-testid="audit-filter-to">
      </div>
      <div>
        <label class="mb-1 block text-xs font-medium text-ds-fg-muted">{{ $t('auditGranularity') }}</label>
        <select v-model="filters.granularity" class="ds-input rounded-xl px-3 py-2 text-sm" data-testid="audit-filter-granularity">
          <option value="day">{{ $t('auditGranDay') }}</option>
          <option value="week">{{ $t('auditGranWeek') }}</option>
        </select>
      </div>
      <button
        type="button"
        class="ds-btn-primary rounded-xl px-4 py-2 text-sm font-medium"
        data-testid="audit-refresh-btn"
        @click="refresh"
      >
        {{ $t('auditRefresh') }}
      </button>
      <button
        v-if="canExport"
        type="button"
        class="rounded-xl border border-ds-border px-4 py-2 text-sm font-medium text-ds-fg-brand transition-colors hover:border-ds-brand"
        data-testid="audit-export-btn"
        @click="handleExport"
      >
        {{ $t('adminExport') }} CSV
      </button>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="h-6 w-6 animate-spin rounded-full border-2 border-ds-brand-muted border-t-ds-brand" />
    </div>

    <template v-else>
      <!-- Charts row -->
      <div class="grid gap-6 lg:grid-cols-2">
        <!-- By Action Donut -->
        <div class="ds-card rounded-2xl p-5 shadow-sm">
          <h3 class="mb-1 text-sm font-semibold text-ds-fg-heading">{{ $t('auditByActionTitle') }}</h3>
          <p class="mb-3 text-xs text-ds-fg-muted">{{ $t('auditByActionHint') }}</p>
          <ChartDonut
            v-if="donutItems.length > 0"
            :items="donutItems"
            :height="260"
            data-testid="audit-donut"
          />
          <p v-else class="py-8 text-center text-sm text-ds-fg-muted">{{ $t('dashboardChartEmpty') }}</p>
        </div>

        <!-- Trend Bar -->
        <div class="ds-card rounded-2xl p-5 shadow-sm">
          <h3 class="mb-1 text-sm font-semibold text-ds-fg-heading">{{ $t('auditTrendTitle') }}</h3>
          <p class="mb-3 text-xs text-ds-fg-muted">{{ $t('auditTrendHint') }}</p>
          <ChartBar
            v-if="trendItems.length > 0"
            :items="trendItems"
            :height="260"
            :horizontal="false"
            :value-formatter="(v: number) => String(v)"
            data-testid="audit-trend-bar"
          />
          <p v-else class="py-8 text-center text-sm text-ds-fg-muted">{{ $t('dashboardChartEmpty') }}</p>
        </div>
      </div>

      <!-- Top Actors Bar -->
      <div class="ds-card rounded-2xl p-5 shadow-sm">
        <h3 class="mb-1 text-sm font-semibold text-ds-fg-heading">{{ $t('auditTopActorsTitle') }}</h3>
        <p class="mb-3 text-xs text-ds-fg-muted">{{ $t('auditTopActorsHint') }}</p>
        <ChartBar
          v-if="topActorItems.length > 0"
          :items="topActorItems"
          :height="200"
          :horizontal="true"
          :value-formatter="(v: number) => String(v)"
          data-testid="audit-top-actors-bar"
        />
        <p v-else class="py-8 text-center text-sm text-ds-fg-muted">{{ $t('dashboardChartEmpty') }}</p>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import type { ChartBarItem, ChartDonutItem } from '@crm/ui-kit'
import { defaultAuditRange } from '~/composables/use-audit-stats'

const { can } = usePermission()
const auditStats = useAuditStats()

const canExport = computed(() => can('audit', 'export'))

const loading = ref(true)
const donutItems = ref<ChartDonutItem[]>([])
const trendItems = ref<ChartBarItem[]>([])
const topActorItems = ref<ChartBarItem[]>([])

const range = defaultAuditRange()
const filters = reactive({
  from: range.from,
  to: range.to,
  granularity: 'day' as 'day' | 'week',
})

async function refresh() {
  loading.value = true
  try {
    const [byAction, trend, topActors] = await Promise.all([
      auditStats.fetchByAction({ from: filters.from, to: filters.to }),
      auditStats.fetchTrend({ from: filters.from, to: filters.to, granularity: filters.granularity }),
      auditStats.fetchTopActors({ from: filters.from, to: filters.to, limit: 10 }),
    ])

    donutItems.value = byAction.items.map((i) => ({ name: i.action, value: i.count }))
    trendItems.value = trend.items.map((i) => ({ name: i.date, value: i.count }))
    topActorItems.value = topActors.items.map((i) => ({ name: i.actor_name, value: i.count }))
  } finally {
    loading.value = false
  }
}

async function handleExport() {
  try {
    const blob = await auditStats.exportCsv({ from: filters.from, to: filters.to })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `audit-export-${filters.from}-${filters.to}.csv`
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    // silently handle
  }
}

onMounted(refresh)
</script>
