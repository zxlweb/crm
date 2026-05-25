<template>
  <PermissionGuard resource="accounts" action="view">
    <div class="space-y-4" data-testid="accounts-page">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p class="max-w-2xl text-sm text-ds-fg-muted">
          {{ $t('accountsPageDescApi') }}
        </p>
        <UiButton
          v-if="canCreate"
          variant="secondary"
          size="sm"
          class="shrink-0"
          icon="i-heroicons-plus-20-solid"
          data-testid="account-create-btn"
          :loading="saving"
          @click="openCreate"
        >
          {{ $t('accountsCreate') }}
        </UiButton>
      </div>

      <div v-if="pending" class="flex justify-center py-24">
        <UIcon name="i-heroicons-arrow-path" class="h-8 w-8 animate-spin text-primary" />
      </div>

      <UAlert v-else-if="loadError" color="red" variant="soft" :title="loadError" />

      <template v-else>
        <AccountsListTable
          :items="items"
          :can-edit="canUpdate"
          @edit="openEdit"
        >
          <template #toolbar>
            <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
              <UiInput
                v-model="search"
                search
                type="search"
                class="flex-1"
                :placeholder="$t('accountsSearchPlaceholder')"
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
                v-model="lifecycleFilter"
                class="sm:w-48"
                :items="lifecycleSelectItems"
                :placeholder="$t('accountsFilterAllLifecycle')"
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
        </AccountsListTable>
      </template>

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
