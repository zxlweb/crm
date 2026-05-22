<template>
  <PermissionGuard resource="leads" action="view">
    <div class="space-y-6" data-testid="leads-page">
      <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p class="text-xs font-medium text-ds-fg-brand">{{ $t('leadsPageLabel') }}</p>
          <h2 class="text-xl font-bold text-ds-fg-heading">{{ $t('leadsPageTitle') }}</h2>
          <p class="mt-1 text-sm text-ds-fg-muted">{{ $t('leadsPageDescApi') }}</p>
        </div>
        <UiButton
          v-if="canCreate"
          data-testid="lead-create-btn"
          :loading="creating"
          @click="createOpen = true"
        >
          {{ $t('leadsCreate') }}
        </UiButton>
      </div>

      <UiTabs v-model="activeTab" :items="mainTabs" />

      <div v-if="pending" class="flex justify-center py-24">
        <UIcon name="i-heroicons-arrow-path" class="h-8 w-8 animate-spin text-primary" />
      </div>

      <UAlert v-else-if="loadError" color="red" variant="soft" :title="loadError" />

      <template v-else>
        <Transition
          mode="out-in"
          enter-active-class="transition-opacity duration-200 ease-out"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="transition-opacity duration-150 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="activeTab === 'list'" key="list" class="space-y-4">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
              <UiInput
                v-model="search"
                type="search"
                class="flex-1"
                :placeholder="$t('leadsSearchPlaceholder')"
                @keyup.enter="reload"
              />
              <UiSelect
                v-model="statusFilter"
                class="sm:w-48"
                :items="statusSelectItems"
                :placeholder="$t('leadsFilterAllStatus')"
              />
            </div>

            <LeadsListTable :items="items" />

            <p v-if="pagination" class="text-xs text-ds-fg-muted">
              {{ $t('adminShowing', { count: items.length, total: pagination.total }) }}
            </p>
          </div>

          <LeadsReportsPanel v-else key="reports" />
        </Transition>
      </template>

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
import type { LeadStatus } from '~/types/lead'
import type { Lead } from '~/types/lead'

definePageMeta({ layout: 'app', middleware: 'auth' })

const { t } = useI18n()
const permission = usePermission()
const leadsApi = useLeads()

const activeTab = ref('list')
const items = ref<Lead[]>([])
const pagination = ref<{ total: number } | null>(null)
const pending = ref(true)
const loadError = ref('')
const search = ref('')
const statusFilter = ref('')
const creating = ref(false)
const createTitle = ref('')
const createOpen = ref(false)

const statusOptions: LeadStatus[] = ['new', 'contacted', 'qualified', 'unqualified', 'converted']

const canCreate = computed(() => permission.can('leads', 'create'))

const mainTabs = computed(() => [
  { id: 'list', label: t('leadsTabList') },
  { id: 'reports', label: t('leadsTabReports') },
])

const statusSelectItems = computed(() => [
  { label: t('leadsFilterAllStatus'), value: '' },
  ...statusOptions.map((s) => ({ label: t(`leadStatus.${s}`), value: s })),
])

async function reload() {
  pending.value = true
  loadError.value = ''
  try {
    const { data, pagination: page } = await leadsApi.fetchList({
      search: search.value || undefined,
      status: (statusFilter.value || undefined) as LeadStatus | undefined,
    })
    items.value = data.items
    pagination.value = page
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

watch(statusFilter, () => reload())

onMounted(reload)
</script>
