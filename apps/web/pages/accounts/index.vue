<template>
  <PermissionGuard resource="accounts" action="view">
    <div class="accounts-page relative space-y-5" data-testid="accounts-page">
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
            {{ $t('accountsPageDescApi') }}
          </p>
        </div>
        <button
          v-if="canCreate"
          type="button"
          class="group relative inline-flex shrink-0 cursor-pointer items-center gap-1.5 overflow-hidden rounded-xl px-4 py-2 text-sm font-semibold text-ds-on-brand shadow-ds-brand transition-[transform,box-shadow] duration-200 hover:shadow-ds-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ds-brand focus-visible:ring-offset-2 focus-visible:ring-offset-ds-bg disabled:cursor-not-allowed disabled:opacity-60"
          :style="{ background: 'var(--ds-brand-gradient)' }"
          data-testid="account-create-btn"
          :disabled="saving"
          @click="openCreate"
        >
          <span
            class="pointer-events-none absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/30 to-transparent opacity-0 transition-[transform,opacity] duration-500 group-hover:translate-x-full group-hover:opacity-100"
            aria-hidden="true"
          />
          <UIcon name="i-heroicons-plus-20-solid" class="h-4 w-4" aria-hidden="true" />
          <span>{{ $t('accountsCreate') }}</span>
        </button>
      </header>

      <CrmEntityListHero
        v-if="!pending"
        :items="items"
        :total="pagination?.total ?? items.length"
        i18n-prefix="accountsHero"
        at-risk-href="/accounts?health=low"
        test-id="accounts-list-hero"
      />

      <UAlert v-if="loadError" color="red" variant="soft" :title="loadError" />

      <Transition
        mode="out-in"
        enter-active-class="transition-opacity duration-200 ease-out"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition-opacity duration-150 ease-in"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div key="list">
          <div v-if="pending" class="flex justify-center py-24">
            <UIcon name="i-heroicons-arrow-path" class="h-8 w-8 animate-spin text-primary" />
          </div>
          <AccountsListTable
            v-else
            :items="items"
            :can-edit="canUpdate"
            @edit="openEdit"
          >
            <template #toolbar>
              <div class="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center">
                <UiInput
                  v-model="search"
                  search
                  type="search"
                  class="w-full min-w-[12rem] sm:w-56 lg:w-64"
                  :placeholder="$t('accountsSearchPlaceholder')"
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
                  data-testid="accounts-filter-reset"
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
          </AccountsListTable>
        </div>
      </Transition>

      <UiModal v-model:open="formOpen" :title="formTitle">
        <form class="space-y-4" data-testid="account-form" @submit.prevent="submitForm">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="account-name">{{ $t('accountsColName') }}</label>
            <UiInput id="account-name" v-model="formName" required />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="account-industry">{{ $t('accountsColIndustry') }}</label>
            <UiInput id="account-industry" v-model="formIndustry" />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="account-website">{{ $t('accountsFieldWebsite') }}</label>
            <UiInput id="account-website" v-model="formWebsite" type="url" />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="account-lifecycle">{{ $t('leadsColLifecycle') }}</label>
            <UiSelect id="account-lifecycle" v-model="formLifecycle" :items="lifecycleFormItems" />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="account-tags">{{ $t('leadsFieldTags') }}</label>
            <UiInput
              id="account-tags"
              v-model="formTagsText"
              :placeholder="$t('accountsTagsHint')"
            />
          </div>
        </form>
        <template #footer>
          <div class="flex justify-end gap-2">
            <UiButton variant="secondary" @click="formOpen = false">{{ $t('cancel') }}</UiButton>
            <UiButton :loading="saving" data-testid="account-form-submit" @click="submitForm">{{ $t('save') }}</UiButton>
          </div>
        </template>
      </UiModal>
    </div>
  </PermissionGuard>
</template>

<script setup lang="ts">
import type { Account, LifecycleStage, RelationshipHealth } from '~/types/account'
import type { SegmentTemplate } from '~/types/segment'
import { isSegmentCode } from '~/types/segment'

definePageMeta({ layout: 'app', middleware: 'auth' })

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const permission = usePermission()
const accountsApi = useAccounts()
const segmentsApi = useSegments()

const lifecycleOptions: LifecycleStage[] = ['acquire', 'activate', 'grow', 'retain', 'revive']
const healthOptions: RelationshipHealth[] = ['high', 'medium', 'low']

const LIST_PAGE_SIZE = 10

const page = ref(1)
const items = ref<Account[]>([])
const pagination = ref<{ page: number; page_size: number; total: number } | null>(null)
const pending = ref(true)
const loadError = ref('')
const search = ref('')
const segmentFilter = ref('')
const lifecycleFilter = ref('')
const healthFilter = ref('')
const segments = ref<SegmentTemplate[]>([])
const filtersReady = ref(false)
const saving = ref(false)
const formOpen = ref(false)
const editingId = ref<string | null>(null)

const formName = ref('')
const formIndustry = ref('')
const formWebsite = ref('')
const formLifecycle = ref<LifecycleStage>('acquire')
const formTagsText = ref('')

const canCreate = computed(() => permission.can('accounts', 'create'))
const canUpdate = computed(() => permission.can('accounts', 'update'))

const formTitle = computed(() =>
  editingId.value ? t('accountsEditTitle') : t('accountsCreateTitle'),
)

const lifecycleSelectItems = computed(() => [
  { label: t('accountsFilterAllLifecycle'), value: '' },
  ...lifecycleOptions.map((s) => ({ label: t(`lifecycle.${s}`), value: s })),
])

const lifecycleFormItems = computed(() =>
  lifecycleOptions.map((s) => ({ label: t(`lifecycle.${s}`), value: s })),
)

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

const hasActiveFilter = computed(
  () =>
    Boolean(segmentFilter.value) ||
    Boolean(lifecycleFilter.value) ||
    Boolean(healthFilter.value) ||
    Boolean(search.value),
)

function parseTags(text: string): string[] {
  return text
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
}

function resetForm() {
  formName.value = ''
  formIndustry.value = ''
  formWebsite.value = ''
  formLifecycle.value = 'acquire'
  formTagsText.value = ''
  editingId.value = null
}

function openCreate() {
  resetForm()
  formOpen.value = true
}

function openEdit(row: Account) {
  if (!canUpdate.value) return
  editingId.value = row.id
  formName.value = row.name
  formIndustry.value = row.industry ?? ''
  formWebsite.value = row.website ?? ''
  formLifecycle.value = row.lifecycle_stage
  formTagsText.value = (row.tags ?? []).join(', ')
  formOpen.value = true
}

function formatLoadError(e: unknown): string {
  const msg = e instanceof Error ? e.message : t('loadFailed')
  if (/404|page not found/i.test(msg)) {
    return t('accountsApiNotFoundHint')
  }
  return msg
}

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

function resetFilters() {
  search.value = ''
  segmentFilter.value = ''
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
  const tenant = useTenant()
  if (!tenant.currentTenantId.value) {
    pending.value = false
    loadError.value = t('tenantRequiredHint')
    return
  }

  pending.value = true
  loadError.value = ''
  try {
    const { data, pagination: pageMeta } = await accountsApi.fetchList({
      page: page.value,
      page_size: LIST_PAGE_SIZE,
      search: search.value || undefined,
      lifecycle_stage: (lifecycleFilter.value || undefined) as LifecycleStage | undefined,
      relationship_health: (healthFilter.value || undefined) as RelationshipHealth | undefined,
      segment: segmentFilter.value || undefined,
    })
    items.value = data.items
    pagination.value = pageMeta
    page.value = pageMeta.page
  } catch (e) {
    loadError.value = formatLoadError(e)
  } finally {
    pending.value = false
  }
}

async function submitForm() {
  if (!formName.value.trim()) return
  saving.value = true
  loadError.value = ''
  const payload = {
    name: formName.value.trim(),
    industry: formIndustry.value.trim() || undefined,
    website: formWebsite.value.trim() || undefined,
    lifecycle_stage: formLifecycle.value,
    tags: parseTags(formTagsText.value),
  }
  try {
    if (editingId.value) {
      await accountsApi.update(editingId.value, payload)
    } else {
      await accountsApi.create(payload)
    }
    formOpen.value = false
    resetForm()
    await reload()
  } catch (e) {
    loadError.value = formatLoadError(e)
  } finally {
    saving.value = false
  }
}

watch([lifecycleFilter, healthFilter], () => {
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
    openCreate()
  }
}

onMounted(async () => {
  segments.value = await segmentsApi.fetchListWithCounts()
  applyRouteQuery()
  filtersReady.value = true
  await reload()
})
</script>
