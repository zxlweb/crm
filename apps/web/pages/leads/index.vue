<template>
  <PermissionGuard resource="leads" action="view">
    <div class="space-y-4" data-testid="leads-page">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p class="max-w-2xl text-sm text-ds-fg-muted">
          {{ $t('leadsPageDescApi') }}
        </p>
        <UiButton
          v-if="canCreate"
          variant="secondary"
          size="sm"
          class="shrink-0"
          icon="i-heroicons-plus-20-solid"
          data-testid="lead-create-btn"
          :loading="creating"
          @click="createOpen = true"
        >
          {{ $t('leadsCreate') }}
        </UiButton>
      </div>

      <UiTabs v-model="activeTab" :items="mainTabs" class="max-w-xs" data-testid="leads-main-tabs" />

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
            <LeadsListTable v-else :items="items">
              <template #toolbar>
                <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
                  <UiInput
                    v-model="search"
                    search
                    type="search"
                    class="flex-1"
                    :placeholder="$t('leadsSearchPlaceholder')"
                    @keyup.enter="onSearch"
                  />
                  <UiSelect
                    v-model="segmentFilter"
                    class="sm:w-52"
                    :items="segmentSelectItems"
                    :placeholder="$t('segmentAll')"
                    data-testid="segment-select"
                  />
                  <UiSelect
                    v-model="statusFilter"
                    class="sm:w-48"
                    :items="statusSelectItems"
                    :placeholder="$t('leadsFilterAllStatus')"
                  />
                  <UiSelect
                    v-model="healthFilter"
                    class="sm:w-48"
                    :items="healthSelectItems"
                    :placeholder="$t('accountsFilterAllHealth')"
                  />
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
        <form class="space-y-4" @submit.prevent="submitCreate">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('leadsColTitle') }}</label>
            <UiInput v-model="createTitle" required />
          </div>
        </form>
        <template #footer>
          <div class="flex justify-end gap-2">
            <UiButton variant="secondary" @click="createOpen = false">{{ $t('cancel') }}</UiButton>
            <UiButton :loading="creating" @click="submitCreate">{{ $t('save') }}</UiButton>
          </div>
        </template>
      </UiModal>
    </div>
  </PermissionGuard>
</template>

<script setup lang="ts">
import type { Lead, LeadStatus, RelationshipHealth } from '~/types/lead'
import type { SegmentTemplate } from '~/types/segment'
import { isSegmentCode } from '~/types/segment'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const permission = usePermission()
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
const healthFilter = ref('')
const segments = ref<SegmentTemplate[]>([])
const filtersReady = ref(false)
const creating = ref(false)
const createTitle = ref('')
const createOpen = ref(false)

const statusOptions: LeadStatus[] = ['new', 'contacted', 'qualified', 'unqualified', 'converted']
const healthOptions: RelationshipHealth[] = ['high', 'medium', 'low']

const canCreate = computed(() => permission.can('leads', 'create'))

const mainTabs = computed(() => [
  { id: 'list', label: t('leadsTabList') },
  { id: 'reports', label: t('leadsTabReports') },
])

const statusSelectItems = computed(() => [
  { label: t('leadsFilterAllStatus'), value: '' },
  ...statusOptions.map((s) => ({ label: t(`leadStatus.${s}`), value: s })),
])

const healthSelectItems = computed(() => [
  { label: t('accountsFilterAllHealth'), value: '' },
  ...healthOptions.map((h) => ({ label: t(`relationshipHealth.${h}`), value: h })),
])

const segmentSelectItems = computed(() => [
  { label: t('segmentAll'), value: '' },
  ...segments.value.map((segment) => ({
    label:
      segment.count != null
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

function onSearch() {
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
      relationship_health: (healthFilter.value || undefined) as RelationshipHealth | undefined,
      segment: segmentFilter.value || undefined,
    })
    items.value = data.items
    pagination.value = pageMeta
    page.value = pageMeta.page
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

async function submitCreate() {
  if (!createTitle.value.trim()) return
  creating.value = true
  try {
    await leadsApi.create({ title: createTitle.value.trim() })
    createOpen.value = false
    createTitle.value = ''
    await reload()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    creating.value = false
  }
}

watch([statusFilter, healthFilter], () => {
  if (!filtersReady.value) return
  page.value = 1
  reload()
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

function applyRouteQuery() {
  const tab = route.query.tab
  if (tab === 'reports') activeTab.value = 'reports'
  const health = route.query.health
  if (typeof health === 'string' && healthOptions.includes(health as RelationshipHealth)) {
    healthFilter.value = health
  }
  const seg = route.query.segment
  if (typeof seg === 'string' && isSegmentCode(seg)) {
    segmentFilter.value = seg
  } else {
    segmentFilter.value = ''
  }
  if (route.query.create === '1' && canCreate.value) {
    createOpen.value = true
  }
}

onMounted(async () => {
  segments.value = await segmentsApi.fetchListWithCounts()
  applyRouteQuery()
  filtersReady.value = true
  await reload()
})
</script>
