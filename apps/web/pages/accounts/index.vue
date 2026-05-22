<template>
  <PermissionGuard resource="accounts" action="view">
    <div class="space-y-6" data-testid="accounts-page">
      <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p class="text-xs font-medium text-ds-fg-brand">{{ $t('accountsPageLabel') }}</p>
          <h2 class="text-xl font-bold text-ds-fg-heading">{{ $t('accountsPageTitle') }}</h2>
          <p class="mt-1 text-sm text-ds-fg-muted">{{ $t('accountsPageDescApi') }}</p>
        </div>
        <UiButton
          v-if="canCreate"
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
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
          <UiInput
            v-model="search"
            type="search"
            class="flex-1"
            :placeholder="$t('accountsSearchPlaceholder')"
            @keyup.enter="reload"
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

        <AccountsListTable
          :items="items"
          :can-edit="canUpdate"
          @edit="openEdit"
        />

        <p v-if="pagination" class="text-xs text-ds-fg-muted">
          {{ $t('adminShowing', { count: items.length, total: pagination.total }) }}
        </p>
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

definePageMeta({ layout: 'app', middleware: 'auth' })

const { t } = useI18n()
const permission = usePermission()
const accountsApi = useAccounts()

const lifecycleOptions: LifecycleStage[] = ['acquire', 'activate', 'grow', 'retain', 'revive']
const healthOptions: RelationshipHealth[] = ['high', 'medium', 'low']

const items = ref<Account[]>([])
const pagination = ref<{ total: number } | null>(null)
const pending = ref(true)
const loadError = ref('')
const search = ref('')
const lifecycleFilter = ref('')
const healthFilter = ref('')
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
    const { data, pagination: page } = await accountsApi.fetchList({
      search: search.value || undefined,
      lifecycle_stage: (lifecycleFilter.value || undefined) as LifecycleStage | undefined,
      relationship_health: (healthFilter.value || undefined) as RelationshipHealth | undefined,
    })
    items.value = data.items
    pagination.value = page
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

watch([lifecycleFilter, healthFilter], () => reload())

onMounted(reload)
</script>
