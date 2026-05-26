<template>
  <PermissionGuard resource="contacts" action="view">
    <div class="contacts-page relative space-y-5" data-testid="contacts-page">
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
            {{ $t('contactsPageDescApi') }}
          </p>
        </div>
        <button
          v-if="canCreate"
          type="button"
          class="group relative inline-flex shrink-0 cursor-pointer items-center gap-1.5 overflow-hidden rounded-xl px-4 py-2 text-sm font-semibold text-ds-on-brand shadow-ds-brand transition-[transform,box-shadow] duration-200 hover:shadow-ds-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ds-brand focus-visible:ring-offset-2 focus-visible:ring-offset-ds-bg disabled:cursor-not-allowed disabled:opacity-60"
          :style="{ background: 'var(--ds-brand-gradient)' }"
          data-testid="contact-create-btn"
          :disabled="saving"
          @click="openCreate"
        >
          <span
            class="pointer-events-none absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/30 to-transparent opacity-0 transition-[transform,opacity] duration-500 group-hover:translate-x-full group-hover:opacity-100"
            aria-hidden="true"
          />
          <UIcon name="i-heroicons-plus-20-solid" class="h-4 w-4" aria-hidden="true" />
          <span>{{ $t('contactsCreate') }}</span>
        </button>
      </header>

      <CrmEntityListHero
        v-if="!pending"
        :items="items"
        :total="pagination?.total ?? items.length"
        i18n-prefix="contactsHero"
        at-risk-href="/contacts?health=low"
        test-id="contacts-list-hero"
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
          <ContactsListTable
            v-else
            :items="items"
            :can-edit="canUpdate"
            :account-names="accountNames"
            @edit="openEdit"
          >
            <template #toolbar>
              <div class="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center">
                <UiInput
                  v-model="search"
                  search
                  type="search"
                  class="w-full min-w-[12rem] sm:w-56 lg:w-64"
                  :placeholder="$t('contactsSearchPlaceholder')"
                  @keyup.enter="onSearch"
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
                <UiSelect
                  v-model="accountFilter"
                  class="w-full sm:w-40"
                  :items="accountSelectItems"
                  :placeholder="$t('contactsFilterAllAccounts')"
                />
                <button
                  v-if="hasActiveFilter"
                  type="button"
                  class="inline-flex shrink-0 cursor-pointer items-center gap-1 rounded-lg border border-ds-border-muted bg-ds-bg-elevated px-2.5 py-1.5 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand-muted hover:text-ds-fg-brand"
                  data-testid="contacts-filter-reset"
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
          </ContactsListTable>
        </div>
      </Transition>

      <UiModal v-model:open="formOpen" :title="formTitle">
        <form class="space-y-4" data-testid="contact-form" @submit.prevent="submitForm">
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('contactsFieldFirstName') }}</label>
              <UiInput v-model="formFirstName" data-testid="contact-form-first-name" />
            </div>
            <div>
              <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('contactsFieldLastName') }}</label>
              <UiInput v-model="formLastName" data-testid="contact-form-last-name" />
            </div>
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('contactsColEmail') }}</label>
            <UiInput v-model="formEmail" type="email" data-testid="contact-form-email" />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('contactsFieldPhone') }}</label>
            <UiInput v-model="formPhone" />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('contactsColAccount') }}</label>
            <UiSelect v-model="formAccountId" :items="accountFormItems" />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('leadsColLifecycle') }}</label>
            <UiSelect v-model="formLifecycle" :items="lifecycleFormItems" />
          </div>
          <label class="flex items-center gap-2 text-sm text-ds-fg">
            <input v-model="formIsPrimary" type="checkbox" class="rounded border-ds-border" />
            {{ $t('contactsFieldPrimary') }}
          </label>
        </form>
        <template #footer>
          <div class="flex justify-end gap-2">
            <UiButton variant="secondary" @click="formOpen = false">{{ $t('cancel') }}</UiButton>
            <UiButton :loading="saving" data-testid="contact-form-submit" @click="submitForm">{{ $t('save') }}</UiButton>
          </div>
        </template>
      </UiModal>
    </div>
  </PermissionGuard>
</template>

<script setup lang="ts">
import type { Contact, LifecycleStage, RelationshipHealth } from '~/types/contact'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const permission = usePermission()
const tenant = useTenant()
const contactsApi = useContacts()
const accountsApi = useAccounts()

const lifecycleOptions: LifecycleStage[] = ['acquire', 'activate', 'grow', 'retain', 'revive']
const healthOptions: RelationshipHealth[] = ['high', 'medium', 'low']

const LIST_PAGE_SIZE = 10

const page = ref(1)
const items = ref<Contact[]>([])
const pagination = ref<{ page: number; page_size: number; total: number } | null>(null)
const pending = ref(true)
const loadError = ref('')
const search = ref('')
const lifecycleFilter = ref('')
const healthFilter = ref('')
const accountFilter = ref('')
const filtersReady = ref(false)
const saving = ref(false)
const formOpen = ref(false)
const editingId = ref<string | null>(null)
const accountNames = ref<Record<string, string>>({})
const accountOptions = ref<Array<{ id: string; name: string }>>([])

const formFirstName = ref('')
const formLastName = ref('')
const formEmail = ref('')
const formPhone = ref('')
const formAccountId = ref('')
const formLifecycle = ref<LifecycleStage>('acquire')
const formIsPrimary = ref(false)

const canCreate = computed(() => permission.can('contacts', 'create'))
const canUpdate = computed(() => permission.can('contacts', 'update'))

const formTitle = computed(() => (editingId.value ? t('contactsEditTitle') : t('contactsCreateTitle')))

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

const accountSelectItems = computed(() => [
  { label: t('contactsFilterAllAccounts'), value: '' },
  ...accountOptions.value.map((a) => ({ label: a.name, value: a.id })),
])

const accountFormItems = computed(() => [
  { label: t('contactsNoAccount'), value: '' },
  ...accountOptions.value.map((a) => ({ label: a.name, value: a.id })),
])

const hasActiveFilter = computed(
  () =>
    Boolean(lifecycleFilter.value) ||
    Boolean(healthFilter.value) ||
    Boolean(accountFilter.value) ||
    Boolean(search.value),
)

const tableRangeLabel = computed(() => {
  if (!pagination.value?.total) return ''
  const { page: p, page_size: size, total } = pagination.value
  const start = (p - 1) * size + 1
  const end = Math.min(p * size, total)
  return t('tableShowingRange', { start, end, total })
})

function resetForm() {
  formFirstName.value = ''
  formLastName.value = ''
  formEmail.value = ''
  formPhone.value = ''
  formAccountId.value = route.query.account_id ? String(route.query.account_id) : ''
  formLifecycle.value = 'acquire'
  formIsPrimary.value = false
  editingId.value = null
}

function openCreate() {
  resetForm()
  formOpen.value = true
}

function openEdit(row: Contact) {
  if (!canUpdate.value) return
  editingId.value = row.id
  formFirstName.value = row.first_name
  formLastName.value = row.last_name
  formEmail.value = row.email
  formPhone.value = row.phone
  formAccountId.value = row.account_id ?? ''
  formLifecycle.value = row.lifecycle_stage
  formIsPrimary.value = row.is_primary
  formOpen.value = true
}

async function loadAccountOptions() {
  try {
    const { data } = await accountsApi.fetchList({ page: 1, page_size: 100 })
    accountOptions.value = data.items.map((a) => ({ id: a.id, name: a.name }))
    accountNames.value = Object.fromEntries(data.items.map((a) => [a.id, a.name]))
  } catch {
    accountOptions.value = []
  }
}

async function reload() {
  if (!tenant.currentTenantId.value) {
    pending.value = false
    items.value = []
    pagination.value = null
    loadError.value = t('tenantRequiredHint')
    return
  }

  pending.value = true
  loadError.value = ''
  try {
    const { data, pagination: pageMeta } = await contactsApi.fetchList({
      page: page.value,
      page_size: LIST_PAGE_SIZE,
      search: search.value || undefined,
      lifecycle_stage: (lifecycleFilter.value || undefined) as LifecycleStage | undefined,
      relationship_health: (healthFilter.value || undefined) as RelationshipHealth | undefined,
      account_id: accountFilter.value || undefined,
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

async function submitForm() {
  const hasIdentity =
    formFirstName.value.trim() || formLastName.value.trim() || formEmail.value.trim()
  if (!hasIdentity) {
    loadError.value = t('contactsNameRequired')
    return
  }
  saving.value = true
  loadError.value = ''
  const payload = {
    first_name: formFirstName.value.trim(),
    last_name: formLastName.value.trim(),
    email: formEmail.value.trim(),
    phone: formPhone.value.trim(),
    account_id: formAccountId.value || null,
    lifecycle_stage: formLifecycle.value,
    is_primary: formIsPrimary.value,
  }
  try {
    if (editingId.value) {
      await contactsApi.update(editingId.value, payload)
    } else {
      await contactsApi.create(payload)
    }
    formOpen.value = false
    resetForm()
    await reload()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    saving.value = false
  }
}

function onSearch() {
  page.value = 1
  reload()
}

function resetFilters() {
  search.value = ''
  lifecycleFilter.value = ''
  healthFilter.value = ''
  accountFilter.value = ''
  page.value = 1
  syncAccountQuery('')
  reload()
}

function onPageChange(next: number) {
  page.value = next
  reload()
}

function syncAccountQuery(id: string) {
  const current = typeof route.query.account_id === 'string' ? route.query.account_id : ''
  if (id === current) return
  const query = { ...route.query }
  if (id) query.account_id = id
  else delete query.account_id
  router.replace({ query })
}

function applyRouteQuery() {
  const health = route.query.health
  if (typeof health === 'string' && healthOptions.includes(health as RelationshipHealth)) {
    healthFilter.value = health
  }
  if (route.query.account_id) {
    accountFilter.value = String(route.query.account_id)
  }
  if (route.query.create === '1' && canCreate.value) {
    openCreate()
  }
}

watch([lifecycleFilter, healthFilter], () => {
  if (!filtersReady.value) return
  page.value = 1
  reload()
})

watch(accountFilter, (id) => {
  if (!filtersReady.value) return
  syncAccountQuery(id)
  page.value = 1
  reload()
})

watch(
  () => tenant.currentTenantId.value,
  (id, prev) => {
    if (!id || id === prev) return
    page.value = 1
    void reload()
  },
)

async function ensureTenantContext() {
  if (tenant.currentTenantId.value) return
  try {
    const list = await tenant.fetchTenants()
    if (!tenant.currentTenantId.value && list.length > 0) {
      await tenant.switchTenant(list[0].id)
    }
  } catch {
    // reload() 会展示 tenantRequiredHint
  }
}

onMounted(async () => {
  await ensureTenantContext()
  await loadAccountOptions()
  applyRouteQuery()
  filtersReady.value = true
  await reload()
})
</script>
