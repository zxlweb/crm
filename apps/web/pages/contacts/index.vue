<template>
  <PermissionGuard resource="contacts" action="view">
    <div class="space-y-4" data-testid="contacts-page">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p class="max-w-2xl text-sm text-ds-fg-muted">{{ $t('contactsPageDescApi') }}</p>
        <UiButton
          v-if="canCreate"
          variant="secondary"
          size="sm"
          class="shrink-0"
          icon="i-heroicons-plus-20-solid"
          data-testid="contact-create-btn"
          :loading="saving"
          @click="openCreate"
        >
          {{ $t('contactsCreate') }}
        </UiButton>
      </div>

      <UAlert v-if="loadError" color="red" variant="soft" :title="loadError" class="mb-2" />

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
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
            <UiInput
              v-model="search"
              search
              type="search"
              class="flex-1"
              :placeholder="$t('contactsSearchPlaceholder')"
              @keyup.enter="onSearch"
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
            <UiSelect
              v-model="accountFilter"
              class="sm:w-52"
              :items="accountSelectItems"
              :placeholder="$t('contactsFilterAllAccounts')"
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
      </ContactsListTable>

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
const pending = ref(false)
const loadError = ref('')
const search = ref('')
const lifecycleFilter = ref('')
const healthFilter = ref('')
const accountFilter = ref('')
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

function onPageChange(next: number) {
  page.value = next
  reload()
}

watch([lifecycleFilter, healthFilter, accountFilter], () => {
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
  if (route.query.create === '1' && canCreate.value) {
    openCreate()
  }
  if (route.query.account_id) {
    accountFilter.value = String(route.query.account_id)
  }
  await reload()
})
</script>
