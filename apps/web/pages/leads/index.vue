<template>
  <PermissionGuard resource="leads" action="view">
    <div class="leads-page relative space-y-5" data-testid="leads-page">
      <div
        class="pointer-events-none absolute inset-x-0 top-0 h-64 overflow-hidden"
        aria-hidden="true"
      >
        <div
          class="absolute -left-12 top-0 h-52 w-52 rounded-full blur-3xl opacity-60"
          :style="{ background: 'var(--ds-blur-brand)' }"
        />
        <div
          class="absolute right-0 top-4 h-40 w-40 rounded-full blur-3xl opacity-50"
          :style="{ background: 'var(--ds-blur-accent)' }"
        />
      </div>

      <header
        class="relative flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
      >
        <div class="min-w-0">
          <p class="max-w-2xl text-sm text-ds-fg-muted">
            {{ $t('leadsPageDescApi') }}
          </p>
        </div>
        <div class="flex shrink-0 items-center gap-2">
          <button
            type="button"
            class="inline-flex h-9 cursor-pointer items-center gap-1.5 rounded-xl border border-ds-border-muted bg-ds-bg-elevated px-3 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand-muted hover:text-ds-fg-brand disabled:cursor-not-allowed disabled:opacity-60"
            data-testid="leads-pool-recycle-btn"
            :title="$t('leadsPoolRecycleBtnHint', { days: poolSettings?.inactiveDays ?? 30 })"
            :disabled="recycling || !poolSettings?.enabled"
            @click="runRecycle"
          >
            <UIcon
              :name="recycling ? 'i-heroicons-arrow-path' : 'i-heroicons-arrow-uturn-down'"
              class="h-3.5 w-3.5"
              :class="recycling ? 'animate-spin' : ''"
              aria-hidden="true"
            />
            <span>{{ $t('leadsPoolRecycleBtn') }}</span>
          </button>
          <button
            type="button"
            class="inline-flex h-9 cursor-pointer items-center gap-1.5 rounded-xl border border-ds-border-muted bg-ds-bg-elevated px-3 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand-muted hover:text-ds-fg-brand"
            data-testid="leads-pool-settings-btn"
            @click="poolSettingsOpen = true"
          >
            <UIcon name="i-heroicons-cog-6-tooth" class="h-3.5 w-3.5" aria-hidden="true" />
            <span>{{ $t('leadsPoolSettingsBtn') }}</span>
          </button>
          <button
            v-if="canCreate"
            type="button"
            class="ds-leads-create-btn group relative inline-flex shrink-0 cursor-pointer items-center gap-1.5 overflow-hidden rounded-xl px-4 py-2 text-sm font-semibold text-ds-on-brand shadow-ds-brand transition-[transform,box-shadow] duration-200 hover:shadow-ds-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ds-brand focus-visible:ring-offset-2 focus-visible:ring-offset-ds-bg disabled:cursor-not-allowed disabled:opacity-60"
            :style="{ background: 'var(--ds-brand-gradient)' }"
            data-testid="lead-create-btn"
            :disabled="creating"
            @click="createOpen = true"
          >
            <span
              class="pointer-events-none absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/30 to-transparent opacity-0 transition-[transform,opacity] duration-500 group-hover:translate-x-full group-hover:opacity-100"
              aria-hidden="true"
            />
            <UIcon name="i-heroicons-plus-20-solid" class="h-4 w-4" aria-hidden="true" />
            <span>{{ $t('leadsCreate') }}</span>
          </button>
        </div>
      </header>

      <UAlert
        v-if="recycleNotice"
        :color="recycleNotice.tone"
        variant="soft"
        :title="recycleNotice.title"
        :description="recycleNotice.description"
        :close-button="{ icon: 'i-heroicons-x-mark', variant: 'link', padded: false }"
        @close="recycleNotice = null"
      />

      <LeadsListHero
        v-if="activeTab === 'list' && !pending"
        :items="items"
        :total="pagination?.total ?? items.length"
      />

      <UiTabs v-model="activeTab" :items="mainTabs" class="max-w-xs" data-testid="leads-main-tabs" />

      <div v-if="activeTab === 'list'" class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <UiTabs
          v-model="poolFilter"
          :items="poolTabs"
          :aria-label="$t('leadsPoolGroupAria')"
          data-testid="leads-pool-tabs"
        />
        <div
          v-if="poolStats && poolSettings?.enabled"
          class="flex items-center gap-3 text-[11px] text-ds-fg-subtle"
        >
          <span class="inline-flex items-center gap-1">
            <span class="h-2 w-2 rounded-full bg-ds-warning" aria-hidden="true" />
            {{ $t('leadsPoolRecyclableSoon', { n: poolStats.recyclableSoon }) }}
          </span>
          <span
            v-if="poolSettings.last_recycled_at"
            class="inline-flex items-center gap-1"
          >
            <UIcon name="i-heroicons-clock" class="h-3 w-3" aria-hidden="true" />
            {{ $t('leadsPoolLastRecycledAt', { time: formatTime(poolSettings.last_recycled_at) }) }}
          </span>
        </div>
      </div>

      <UAlert v-if="activeTab === 'list' && loadError" color="red" variant="soft" :title="loadError" />

      <Transition
          mode="out-in"
          enter-active-class="transition-opacity duration-200 ease-out"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="transition-opacity duration-150 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="activeTab === 'list'" key="list">
            <div v-if="pending" class="flex justify-center py-24">
              <UIcon name="i-heroicons-arrow-path" class="h-8 w-8 animate-spin text-primary" />
            </div>
            <LeadsListTable
              v-else
              :items="items"
              :pool="poolFilter"
              :current-user-id="currentUserId"
              :can-claim-pool="canClaim"
              :can-manage-pool="canManagePool"
              :pool-settings="poolSettings"
              :pending-id="actionPendingId"
              @claim="onClaim"
              @release="onRelease"
              @transfer="openTransfer"
            >
              <template #toolbar>
                <div class="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center">
                  <UiInput
                    v-model="search"
                    search
                    type="search"
                    class="w-full min-w-[12rem] sm:w-56 lg:w-64"
                    :placeholder="$t('leadsSearchPlaceholder')"
                    @keyup.enter="onSearch"
                  />
                  <UiSelect
                    v-model="segmentFilter"
                    class="w-full sm:w-40"
                    :items="segmentSelectItems"
                    :placeholder="$t('segmentAll')"
                    data-testid="segment-select"
                  />
                  <UiSelect
                    v-model="statusFilter"
                    class="w-full sm:w-36"
                    :items="statusSelectItems"
                    :placeholder="$t('leadsFilterAllStatus')"
                  />
                  <UiSelect
                    v-model="lifecycleFilter"
                    class="w-full sm:w-36"
                    :items="lifecycleSelectItems"
                    :placeholder="$t('accountsFilterAllLifecycle')"
                  />
                  <UiSelect
                    v-model="healthFilter"
                    class="w-full sm:w-36"
                    :items="healthSelectItems"
                    :placeholder="$t('accountsFilterAllHealth')"
                  />
                  <button
                    v-if="hasActiveFilter"
                    type="button"
                    class="inline-flex shrink-0 cursor-pointer items-center gap-1 rounded-lg border border-ds-border-muted bg-ds-bg-elevated px-2.5 py-1.5 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand-muted hover:text-ds-fg-brand"
                    data-testid="leads-filter-reset"
                    @click="resetFilters"
                  >
                    <UIcon name="i-heroicons-x-mark" class="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{{ $t('leadsFilterReset') }}</span>
                  </button>
                </div>
              </template>
              <template v-if="pagination && pagination.total > 0" #footer>
                <UiTablePagination
                  :page="page"
                  :page-size="pagination.page_size"
                  :total="pagination.total"
                  :range-text="tableRangeLabel"
                  :prev-label="$t('paginationPrev')"
                  :next-label="$t('paginationNext')"
                  @update:page="onPageChange"
                />
              </template>
            </LeadsListTable>
          </div>

          <LeadsReportsPanel v-else key="reports" data-testid="leads-tab-reports" />
        </Transition>

      <UiModal v-model:open="createOpen" :title="$t('leadsCreateTitle')">
        <form class="space-y-4" data-testid="lead-create-form" @submit.prevent="submitCreate">
          <div>
            <label
              class="mb-1.5 block text-sm font-medium text-ds-fg"
              for="leads-create-title"
            >
              {{ $t('leadsColTitle') }}
            </label>
            <UiInput id="leads-create-title" v-model="createTitle" required data-testid="lead-create-title" />
          </div>
          <div>
            <label
              class="mb-1.5 block text-sm font-medium text-ds-fg"
              for="leads-create-amount"
            >
              {{ $t('leadsFieldAmount') }}
            </label>
            <UiInput
              id="leads-create-amount"
              v-model="createAmount"
              type="number"
              min="0"
              step="1000"
              :placeholder="$t('leadsAmountUnset')"
              data-testid="lead-create-amount"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg">
              {{ $t('leadsColSource') }}
            </label>
            <UiSelect
              v-model="createSource"
              :items="createSourceItems"
              :placeholder="$t('leadSource.unknown')"
              data-testid="lead-create-source"
            />
          </div>
          <div>
            <label
              class="mb-1.5 block text-sm font-medium text-ds-fg"
              for="leads-create-close-date"
            >
              {{ $t('leadsFieldExpectedClose') }}
            </label>
            <UiInput
              id="leads-create-close-date"
              v-model="createExpectedCloseDate"
              type="date"
              data-testid="lead-create-close-date"
            />
          </div>
        </form>
        <template #footer>
          <div class="flex justify-end gap-2">
            <UiButton variant="secondary" @click="createOpen = false">{{ $t('cancel') }}</UiButton>
            <UiButton :loading="creating" @click="submitCreate">{{ $t('save') }}</UiButton>
          </div>
        </template>
      </UiModal>

      <UiModal
        v-model:open="poolSettingsOpen"
        :title="$t('leadsPoolSettingsTitle')"
        :subtitle="$t('leadsPoolSettingsSubtitle')"
      >
        <form
          class="space-y-4"
          data-testid="leads-pool-settings-form"
          @submit.prevent="submitPoolSettings"
        >
          <label
            class="flex items-start gap-3 rounded-xl border border-ds-border-muted bg-ds-bg-muted/40 px-3 py-2.5 cursor-pointer"
          >
            <input
              v-model="settingsDraftEnabled"
              type="checkbox"
              class="mt-0.5 h-4 w-4 cursor-pointer rounded border-ds-border accent-ds-brand"
            />
            <span class="flex-1">
              <span class="block text-sm font-medium text-ds-fg-heading">
                {{ $t('leadsPoolSettingsEnableLabel') }}
              </span>
              <span class="mt-0.5 block text-xs text-ds-fg-muted">
                {{ $t('leadsPoolSettingsEnableHint') }}
              </span>
            </span>
          </label>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="pool-inactive-days">
              {{ $t('leadsPoolSettingsDaysLabel') }}
            </label>
            <UiInput
              id="pool-inactive-days"
              v-model="settingsDraftDays"
              type="number"
              min="1"
              max="365"
              :disabled="!settingsDraftEnabled"
            />
            <p class="mt-1.5 text-xs text-ds-fg-subtle">
              {{ $t('leadsPoolSettingsDaysHint') }}
            </p>
          </div>
        </form>
        <template #footer>
          <div class="flex justify-end gap-2">
            <UiButton variant="secondary" @click="poolSettingsOpen = false">{{ $t('cancel') }}</UiButton>
            <UiButton
              :loading="savingSettings"
              data-testid="leads-pool-settings-submit"
              @click="submitPoolSettings"
            >
              {{ $t('save') }}
            </UiButton>
          </div>
        </template>
      </UiModal>

      <UiModal
        v-model:open="transferOpen"
        :title="$t('leadsPoolTransferTitle')"
        :subtitle="transferLead ? transferLead.title : ''"
      >
        <div class="space-y-3" data-testid="leads-pool-transfer-form">
          <p class="text-xs text-ds-fg-muted">{{ $t('leadsPoolTransferHint') }}</p>
          <UiSelect
            v-model="transferTarget"
            :items="transferTargetItems"
            :placeholder="$t('leadsPoolTransferPlaceholder')"
            data-testid="leads-pool-transfer-select"
          />
        </div>
        <template #footer>
          <div class="flex justify-end gap-2">
            <UiButton variant="secondary" @click="transferOpen = false">{{ $t('cancel') }}</UiButton>
            <UiButton
              :loading="transferring"
              :disabled="!transferTarget"
              data-testid="leads-pool-transfer-submit"
              @click="submitTransfer"
            >
              {{ $t('confirm') }}
            </UiButton>
          </div>
        </template>
      </UiModal>
    </div>
  </PermissionGuard>
</template>

<script setup lang="ts">
import { MOCK_USER_PROFILES } from '~/fixtures/users.mock'
import type {
  Lead,
  LeadCreateInput,
  LeadPool,
  LeadPoolSettings,
  LeadPoolStats,
  LeadStatus,
  LifecycleStage,
  RelationshipHealth,
} from '~/types/lead'
import { isLeadPool } from '~/types/lead'
import type { SegmentTemplate } from '~/types/segment'
import { isSegmentCode } from '~/types/segment'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const permission = usePermission()
const auth = useAuth()
const leadsApi = useLeads()
const segmentsApi = useSegments()

const LIST_PAGE_SIZE = 10

const activeTab = ref('list')
const page = ref(1)
const items = ref<Lead[]>([])
const pagination = ref<{ page: number; page_size: number; total: number } | null>(null)
const pending = ref(true)
const loadError = ref('')
const search = ref('')
const segmentFilter = ref('')
const statusFilter = ref('')
const lifecycleFilter = ref<LifecycleStage | ''>('')
const healthFilter = ref('')
const segments = ref<SegmentTemplate[]>([])
const filtersReady = ref(false)
const creating = ref(false)
const createTitle = ref('')
const createAmount = ref('')
const createSource = ref('')
const createExpectedCloseDate = ref('')
const createOpen = ref(false)

const LEAD_SOURCE_CODES = ['website', 'referral', 'exhibition', 'cold_call', 'partner'] as const

// ===== 客户池（公海 / 私海） =====
const poolFilter = ref<LeadPool>('mine')
const poolStats = ref<LeadPoolStats | null>(null)
const poolSettings = ref<LeadPoolSettings | null>(null)
const actionPendingId = ref<string | null>(null)
const recycling = ref(false)
const recycleNotice = ref<
  | null
  | { tone: 'green' | 'amber'; title: string; description?: string }
>(null)

const poolSettingsOpen = ref(false)
const settingsDraftEnabled = ref(true)
const settingsDraftDays = ref('30')
const savingSettings = ref(false)

const transferOpen = ref(false)
const transferLead = ref<Lead | null>(null)
const transferTarget = ref('')
const transferring = ref(false)

const statusOptions: LeadStatus[] = ['new', 'contacted', 'qualified', 'unqualified', 'converted']
const lifecycleOptions: LifecycleStage[] = ['acquire', 'activate', 'grow', 'retain', 'revive']
const healthOptions: RelationshipHealth[] = ['high', 'medium', 'low']

const canCreate = computed(() => permission.can('leads', 'create'))
const canClaim = computed(() => permission.can('leads', 'update') || permission.can('leads', 'create'))
const canManagePool = computed(() => permission.can('leads', 'update'))
const currentUserId = computed(() => auth.user.value?.id ?? null)

const mainTabs = computed(() => [
  { id: 'list', label: t('leadsTabList'), icon: 'i-heroicons-queue-list' },
  { id: 'reports', label: t('leadsTabReports'), icon: 'i-heroicons-chart-bar' },
])

const poolTabs = computed(() => [
  {
    id: 'mine',
    label: t('leadsPoolTabMine'),
    icon: 'i-heroicons-user-circle',
    count: poolStats.value?.mine,
  },
  {
    id: 'public',
    label: t('leadsPoolTabPublic'),
    icon: 'i-heroicons-globe-asia-australia',
    count: poolStats.value?.public,
  },
  {
    id: 'others',
    label: t('leadsPoolTabOthers'),
    icon: 'i-heroicons-user-group',
    count: poolStats.value?.others,
  },
  {
    id: 'all',
    label: t('leadsPoolTabAll'),
    icon: 'i-heroicons-rectangle-stack',
    count: poolStats.value?.all,
  },
])

const transferTargetItems = computed(() => {
  const list = Object.values(MOCK_USER_PROFILES).filter(
    (u) => u.id !== transferLead.value?.owner_id,
  )
  return list.map((u) => ({
    label: `${u.name} · ${u.role || t('leadsPoolMember')}`,
    value: u.id,
  }))
})

const statusSelectItems = computed(() => [
  { label: t('leadsFilterAllStatus'), value: '' },
  ...statusOptions.map((s) => ({ label: t(`leadStatus.${s}`), value: s })),
])

const lifecycleSelectItems = computed(() => [
  { label: t('accountsFilterAllLifecycle'), value: '' },
  ...lifecycleOptions.map((s) => ({ label: t(`lifecycle.${s}`), value: s })),
])

const healthSelectItems = computed(() => [
  { label: t('accountsFilterAllHealth'), value: '' },
  ...healthOptions.map((h) => ({ label: t(`relationshipHealth.${h}`), value: h })),
])

const createSourceItems = computed(() => [
  { label: t('leadSource.unknown'), value: '' },
  ...LEAD_SOURCE_CODES.map((code) => ({
    label: t(`leadSource.${code}`),
    value: code,
  })),
])

const segmentSelectItems = computed(() => [
  { label: t('segmentAll'), value: '' },
  ...segments.value.map((segment) => ({
    label:
      segment.count !== null && segment.count !== undefined
        ? t('segmentOptionCount', { name: t(segment.name_key), count: segment.count })
        : t(segment.name_key),
    value: segment.code,
  })),
])

const tableRangeLabel = computed(() => {
  if (!pagination.value?.total) return ''
  const { page: p, page_size: size, total } = pagination.value
  const start = (p - 1) * size + 1
  const end = Math.min(p * size, total)
  return t('tableShowingRange', { start, end, total })
})

const hasActiveFilter = computed(
  () =>
    Boolean(segmentFilter.value) ||
    Boolean(statusFilter.value) ||
    Boolean(lifecycleFilter.value) ||
    Boolean(healthFilter.value) ||
    Boolean(search.value),
)

function onSearch() {
  page.value = 1
  reload()
}

function resetFilters() {
  search.value = ''
  segmentFilter.value = ''
  statusFilter.value = ''
  lifecycleFilter.value = ''
  healthFilter.value = ''
  page.value = 1
  reload()
}

function onPageChange(next: number) {
  page.value = next
  reload()
}

async function reload() {
  pending.value = true
  loadError.value = ''
  try {
    const { data, pagination: pageMeta } = await leadsApi.fetchList({
      page: page.value,
      page_size: LIST_PAGE_SIZE,
      search: search.value || undefined,
      status: (statusFilter.value || undefined) as LeadStatus | undefined,
      lifecycle_stage: (lifecycleFilter.value || undefined) as LifecycleStage | undefined,
      relationship_health: (healthFilter.value || undefined) as RelationshipHealth | undefined,
      segment: segmentFilter.value || undefined,
      pool: poolFilter.value,
    })
    items.value = data.items
    pagination.value = pageMeta
    page.value = pageMeta.page
    await refreshPoolStats()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

async function refreshPoolStats() {
  try {
    poolStats.value = await leadsApi.poolStats()
  } catch {
    // 客户池统计失败不阻塞主列表
  }
}

async function loadPoolSettings() {
  try {
    poolSettings.value = await leadsApi.getPoolSettings()
    if (poolSettings.value) {
      settingsDraftEnabled.value = poolSettings.value.enabled
      settingsDraftDays.value = String(poolSettings.value.inactiveDays)
    }
  } catch {
    // ignore
  }
}

function formatTime(iso: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso
  const now = Date.now()
  const diff = Math.floor((now - date.getTime()) / 60_000)
  if (diff < 1) return t('justNow')
  if (diff < 60) return t('minutesAgo', { n: diff })
  if (diff < 60 * 24) return t('hoursAgo', { n: Math.floor(diff / 60) })
  return t('daysAgo', { n: Math.floor(diff / (60 * 24)) })
}

async function onClaim(row: Lead) {
  if (actionPendingId.value) return
  actionPendingId.value = row.id
  try {
    await leadsApi.claim(row.id)
    recycleNotice.value = {
      tone: 'green',
      title: t('leadsPoolClaimSuccess', { title: row.title }),
    }
    await reload()
  } catch (e) {
    loadError.value = describeClaimError(e)
  } finally {
    actionPendingId.value = null
  }
}

function describeClaimError(e: unknown): string {
  if (e instanceof Error && e.message === 'lead_already_owned') {
    return t('leadsPoolClaimRaceLost')
  }
  if (e instanceof Error) return e.message
  return t('loadFailed')
}

async function onRelease(row: Lead) {
  if (actionPendingId.value) return
  actionPendingId.value = row.id
  try {
    await leadsApi.release(row.id)
    recycleNotice.value = {
      tone: 'amber',
      title: t('leadsPoolReleaseSuccess', { title: row.title }),
    }
    await reload()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    actionPendingId.value = null
  }
}

function openTransfer(row: Lead) {
  transferLead.value = row
  transferTarget.value = ''
  transferOpen.value = true
}

async function submitTransfer() {
  if (!transferLead.value || !transferTarget.value || transferring.value) return
  transferring.value = true
  try {
    await leadsApi.transfer(transferLead.value.id, transferTarget.value)
    recycleNotice.value = {
      tone: 'green',
      title: t('leadsPoolTransferSuccess', {
        title: transferLead.value.title,
        name: MOCK_USER_PROFILES[transferTarget.value]?.name ?? transferTarget.value,
      }),
    }
    transferOpen.value = false
    transferLead.value = null
    transferTarget.value = ''
    await reload()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    transferring.value = false
  }
}

async function runRecycle() {
  if (recycling.value || !poolSettings.value?.enabled) return
  recycling.value = true
  try {
    const summary = await leadsApi.recycleStale()
    poolSettings.value = await leadsApi.getPoolSettings()
    if (summary.recycled === 0) {
      recycleNotice.value = {
        tone: 'green',
        title: t('leadsPoolRecycleNoneTitle'),
        description: t('leadsPoolRecycleNoneHint', { days: summary.threshold_days }),
      }
    } else {
      recycleNotice.value = {
        tone: 'amber',
        title: t('leadsPoolRecycleDoneTitle', { n: summary.recycled }),
        description: t('leadsPoolRecycleDoneHint', { days: summary.threshold_days }),
      }
    }
    await reload()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    recycling.value = false
  }
}

async function submitPoolSettings() {
  if (savingSettings.value) return
  savingSettings.value = true
  try {
    const parsedDays = Number(settingsDraftDays.value)
    const safeDays = Number.isFinite(parsedDays) && parsedDays > 0 ? Math.trunc(parsedDays) : 30
    poolSettings.value = await leadsApi.updatePoolSettings({
      enabled: settingsDraftEnabled.value,
      inactiveDays: Math.max(1, Math.min(365, safeDays)),
    })
    settingsDraftDays.value = String(poolSettings.value.inactiveDays)
    poolSettingsOpen.value = false
    await refreshPoolStats()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    savingSettings.value = false
  }
}

function resetCreateForm() {
  createTitle.value = ''
  createAmount.value = ''
  createSource.value = ''
  createExpectedCloseDate.value = ''
}

async function submitCreate() {
  if (!createTitle.value.trim()) return
  creating.value = true
  try {
    const payload: LeadCreateInput = { title: createTitle.value.trim() }
    const parsedAmount = Number(createAmount.value)
    if (Number.isFinite(parsedAmount) && parsedAmount > 0) {
      payload.amount = parsedAmount
    }
    if (createSource.value) {
      payload.source = createSource.value
    }
    if (createExpectedCloseDate.value) {
      payload.expected_close_date = createExpectedCloseDate.value
    }
    await leadsApi.create(payload)
    createOpen.value = false
    resetCreateForm()
    await reload()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    creating.value = false
  }
}

watch(createOpen, (open) => {
  if (!open) resetCreateForm()
})

watch([statusFilter, healthFilter, lifecycleFilter], () => {
  if (!filtersReady.value) return
  page.value = 1
  reload()
})

watch(poolFilter, (next, prev) => {
  if (!filtersReady.value) return
  if (next === prev) return
  syncPoolQuery(next)
  page.value = 1
  reload()
})

watch(recycleNotice, (val) => {
  if (!val || import.meta.server) return
  setTimeout(() => {
    if (recycleNotice.value === val) recycleNotice.value = null
  }, 5000)
})

watch(segmentFilter, (code) => {
  if (!filtersReady.value) return
  syncSegmentQuery(code)
  page.value = 1
  reload()
})

watch(
  () => route.query.segment,
  (seg) => {
    if (!filtersReady.value) return
    const code = typeof seg === 'string' && isSegmentCode(seg) ? seg : ''
    if (code !== segmentFilter.value) {
      segmentFilter.value = code
    }
  },
)

function syncSegmentQuery(code: string) {
  const current = typeof route.query.segment === 'string' ? route.query.segment : ''
  if (code === current) return
  const query = { ...route.query }
  if (code) query.segment = code
  else delete query.segment
  router.replace({ query })
}

function syncPoolQuery(value: LeadPool) {
  const current = typeof route.query.pool === 'string' ? route.query.pool : ''
  if (value === current) return
  const query = { ...route.query }
  if (value && value !== 'mine') query.pool = value
  else delete query.pool
  router.replace({ query })
}

function applyRouteQuery() {
  const tab = route.query.tab
  if (tab === 'reports') activeTab.value = 'reports'
  const health = route.query.health
  if (typeof health === 'string' && healthOptions.includes(health as RelationshipHealth)) {
    healthFilter.value = health
  }
  const lifecycle = route.query.lifecycle
  if (typeof lifecycle === 'string' && lifecycleOptions.includes(lifecycle as LifecycleStage)) {
    lifecycleFilter.value = lifecycle
  }
  const seg = route.query.segment
  if (typeof seg === 'string' && isSegmentCode(seg)) {
    segmentFilter.value = seg
  } else {
    segmentFilter.value = ''
  }
  const pool = route.query.pool
  if (typeof pool === 'string' && isLeadPool(pool)) {
    poolFilter.value = pool
  }
  if (route.query.create === '1' && canCreate.value) {
    createOpen.value = true
  }
}

onMounted(async () => {
  segments.value = await segmentsApi.fetchListWithCounts()
  applyRouteQuery()
  await loadPoolSettings()
  filtersReady.value = true
  await reload()
})
</script>
