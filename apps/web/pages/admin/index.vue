<template>
  <div class="space-y-6">
    <div v-if="pending" class="flex items-center justify-center py-24">
      <div class="h-8 w-8 animate-spin rounded-full border-2 border-ds-brand-muted border-t-ds-brand" />
    </div>

    <p v-else-if="loadError" class="rounded-xl bg-ds-danger-subtle px-4 py-3 text-sm text-ds-danger">{{ loadError }}</p>

    <template v-else>
      <section class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <CardMetric
          :label="$t('tenantTotal')"
          :value="overview?.tenant_count ?? 0"
          :compare-label="$t('metricCompareTotal')"
          :trend="`+${overview?.tenant_count ?? 0}`"
          trend-direction="up"
          icon-tone="brand"
        >
          <template #icon>
            <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5" />
            </svg>
          </template>
        </CardMetric>

        <CardMetric
          :label="$t('tenantActive')"
          :value="overview?.active_tenant_count ?? 0"
          :compare-label="$t('metricCompareActiveRate')"
          :trend="activeRate"
          trend-direction="up"
          icon-tone="calendar"
        >
          <template #icon>
            <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </template>
        </CardMetric>

        <CardMetric
          :label="$t('userTotal')"
          :value="overview?.user_count ?? 0"
          :compare-label="$t('metricCompareScope')"
          icon-tone="info"
        >
          <template #icon>
            <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </template>
        </CardMetric>

        <CardMetric
          :label="$t('adminInactiveTenants')"
          :value="inactiveCount"
          :compare-label="$t('metricCompareStatus')"
          :trend="inactiveCount > 0 ? $t('metricNeedsAttention') : $t('metricAllClear')"
          :trend-direction="inactiveCount > 0 ? 'down' : 'flat'"
          icon-tone="neutral"
        >
          <template #icon>
            <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
            </svg>
          </template>
        </CardMetric>
      </section>

      <section class="grid gap-6 lg:grid-cols-3">
        <div class="lg:col-span-2">
          <AdminOverviewChart />
        </div>
        <div class="ds-bg-insight rounded-2xl border border-ds-border p-6 text-ds-on-brand shadow-ds-brand">
          <p class="text-sm font-medium opacity-80">{{ $t('adminInsightLabel') }}</p>
          <p class="mt-3 text-2xl font-bold leading-snug">{{ $t('adminInsightTitle') }}</p>
          <p class="mt-2 text-sm opacity-90">{{ $t('adminInsightDesc') }}</p>
          <div class="mt-6 flex items-end gap-2">
            <span class="text-4xl font-bold">{{ activePercent }}%</span>
            <span class="pb-1 text-sm opacity-80">{{ $t('adminInsightMetric') }}</span>
          </div>
        </div>
      </section>

      <section id="tenants" class="ds-card overflow-hidden rounded-2xl shadow-sm">
        <div class="flex flex-col gap-4 border-b border-ds-border-muted p-5 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h3 class="font-semibold text-ds-fg-heading">{{ $t('tenantList') }}</h3>
            <p class="text-xs text-ds-fg-muted">{{ $t('adminTableHint') }}</p>
          </div>
          <div class="relative">
            <svg class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ds-fg-subtle" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              v-model="search"
              type="search"
              class="ds-input w-full rounded-xl py-2.5 pl-10 pr-4 text-sm sm:w-64"
              :placeholder="$t('searchTenant')"
              @keyup.enter="reloadTenants"
            >
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full min-w-[640px] text-left text-sm">
            <thead>
              <tr class="border-b border-ds-border-muted bg-ds-bg-muted text-xs font-medium uppercase tracking-wide text-ds-fg-brand">
                <th class="px-5 py-3">{{ $t('name') }}</th>
                <th class="px-5 py-3">{{ $t('domain') }}</th>
                <th class="px-5 py-3">{{ $t('members') }}</th>
                <th class="px-5 py-3">{{ $t('status') }}</th>
                <th class="px-5 py-3 text-right">{{ $t('actions') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-ds-border-muted">
              <tr v-for="row in tenants" :key="row.id" class="transition-colors hover:bg-ds-bg-muted">
                <td class="px-5 py-4">
                  <div class="flex items-center gap-3">
                    <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-ds-brand-subtle text-xs font-bold text-ds-fg-brand">
                      {{ row.name.charAt(0).toUpperCase() }}
                    </div>
                    <span class="font-medium text-ds-fg-heading">{{ row.name }}</span>
                  </div>
                </td>
                <td class="px-5 py-4 text-ds-fg-muted">{{ row.domain }}</td>
                <td class="px-5 py-4 text-ds-fg-muted">{{ row.user_count }}</td>
                <td class="px-5 py-4">
                  <span
                    class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium"
                    :class="row.is_active ? 'bg-ds-success-subtle text-ds-success' : 'bg-ds-bg-muted text-ds-fg-muted'"
                  >
                    <span class="h-1.5 w-1.5 rounded-full" :class="row.is_active ? 'bg-ds-success' : 'bg-ds-fg-subtle'" />
                    {{ row.is_active ? $t('active') : $t('inactive') }}
                  </span>
                </td>
                <td class="px-5 py-4 text-right">
                  <button
                    type="button"
                    class="cursor-pointer rounded-lg px-3 py-1.5 text-xs font-medium transition-colors disabled:opacity-50"
                    :class="row.is_active
                      ? 'border border-ds-border text-ds-fg-brand hover:border-ds-brand'
                      : 'bg-ds-brand text-ds-on-brand hover:bg-ds-brand-hover'"
                    :disabled="togglingId === row.id"
                    @click="toggleActive(row)"
                  >
                    {{ row.is_active ? $t('deactivate') : $t('activate') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <p v-if="tenants.length === 0" class="py-12 text-center text-sm text-ds-fg-muted">{{ $t('noTenants') }}</p>

        <div v-if="pagination" class="flex items-center justify-between border-t border-ds-border-muted px-5 py-3 text-xs text-ds-fg-muted">
          <span>{{ $t('adminShowing', { count: tenants.length, total: pagination.total }) }}</span>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import type { SuperAdminOverview, SuperAdminTenant } from '~/composables/use-super-admin'

definePageMeta({ layout: 'admin', middleware: 'super-admin' })

const { t } = useI18n()
const superAdmin = useSuperAdmin()

const overview = ref<SuperAdminOverview | null>(null)
const tenants = ref<SuperAdminTenant[]>([])
const pagination = ref<{ total: number } | null>(null)
const pending = ref(true)
const loadError = ref('')
const search = ref('')
const togglingId = ref<string | null>(null)

const inactiveCount = computed(() => {
  const total = overview.value?.tenant_count ?? 0
  const active = overview.value?.active_tenant_count ?? 0
  return Math.max(0, total - active)
})

const activePercent = computed(() => {
  const total = overview.value?.tenant_count ?? 0
  if (total === 0) return 0
  return Math.round(((overview.value?.active_tenant_count ?? 0) / total) * 100)
})

const activeRate = computed(() => `${activePercent.value}% ${t('adminActiveRate')}`)

async function reloadTenants() {
  const { data, pagination: page } = await superAdmin.fetchTenants({ search: search.value || undefined })
  tenants.value = data.items
  pagination.value = page
}

async function load() {
  pending.value = true
  loadError.value = ''
  try {
    overview.value = await superAdmin.fetchOverview()
    await reloadTenants()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

async function toggleActive(row: SuperAdminTenant) {
  togglingId.value = row.id
  try {
    const updated = await superAdmin.patchTenantActive(row.id, !row.is_active)
    const idx = tenants.value.findIndex((item) => item.id === row.id)
    if (idx >= 0) tenants.value[idx] = updated
    overview.value = await superAdmin.fetchOverview()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    togglingId.value = null
  }
}

onMounted(load)
</script>
