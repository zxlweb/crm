<template>
  <PermissionGuard resource="deals" action="view">
    <div class="deals-page relative space-y-5" data-testid="deals-page">
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
            {{ isPreviewMode ? $t('dealsPageDescMock') : $t('dealsPageDescApi') }}
          </p>
        </div>
        <button
          v-if="canCreate"
          type="button"
          class="group relative inline-flex shrink-0 cursor-pointer items-center gap-1.5 overflow-hidden rounded-xl px-4 py-2 text-sm font-semibold text-ds-on-brand shadow-ds-brand transition-[transform,box-shadow] duration-200 hover:shadow-ds-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ds-brand focus-visible:ring-offset-2 focus-visible:ring-offset-ds-bg disabled:cursor-not-allowed disabled:opacity-60"
          :style="{ background: 'var(--ds-brand-gradient)' }"
          data-testid="deal-create-btn"
          :disabled="creating"
          @click="createOpen = true"
        >
          <span
            class="pointer-events-none absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/30 to-transparent opacity-0 transition-[transform,opacity] duration-500 group-hover:translate-x-full group-hover:opacity-100"
            aria-hidden="true"
          />
          <UIcon name="i-heroicons-plus-20-solid" class="h-4 w-4" aria-hidden="true" />
          <span>{{ $t('dealsCreate') }}</span>
        </button>
      </header>

      <DealsPipelineHero
        v-if="!pending && pipeline"
        :summary="pipeline.summary"
        :stages="pipeline.stages"
      />

      <UiTabs v-model="activeTab" :items="mainTabs" class="max-w-xs" data-testid="deals-main-tabs" />

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
        <div v-if="activeTab === 'pipeline'" key="pipeline">
          <div v-if="pending && !pipeline" class="flex justify-center py-24">
            <UIcon name="i-heroicons-arrow-path" class="h-8 w-8 animate-spin text-primary" />
          </div>
          <template v-else-if="pipeline">
            <DealsFocusStream :stages="pipeline.stages" />

            <DealsPipelineBoard
              :stages="pipeline.stages"
              :pending="pending"
              :can-update="canUpdate"
              :moving-deal-id="movingDealId"
              @move-stage="onMoveStage"
            >
              <template v-if="canCreate" #empty-action>
                <button
                  type="button"
                  class="mt-4 inline-flex cursor-pointer items-center gap-1.5 rounded-xl border border-ds-border-muted bg-ds-bg-elevated px-3 py-2 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand-muted hover:text-ds-fg-brand"
                  @click="createOpen = true"
                >
                  <UIcon name="i-heroicons-plus-20-solid" class="h-3.5 w-3.5" aria-hidden="true" />
                  {{ $t('dealsCreate') }}
                </button>
              </template>
            </DealsPipelineBoard>
          </template>
        </div>

        <DealsReportsPanel v-else key="reports" data-testid="deals-tab-reports" />
      </Transition>

      <UiModal v-model:open="createOpen" :title="$t('dealsCreateTitle')">
        <form class="space-y-4" @submit.prevent="submitCreate">
          <UiInput v-model="createForm.title" :label="$t('dealsFieldTitle')" required />
          <UiInput v-model="createForm.amount" type="number" min="0" :label="$t('dealsFieldAmount')" />
          <UiInput v-model="createForm.expected_close_date" type="date" :label="$t('dealsFieldCloseDate')" />
          <p v-if="createError" class="text-sm text-red-600">{{ createError }}</p>
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
import type { DealPipelineData, DealStage } from '~/types/deal'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const permission = usePermission()
const dealsApi = useDeals()

const activeTab = ref('pipeline')
const pipeline = ref<DealPipelineData | null>(null)
const pending = ref(true)
const loadError = ref('')
const creating = ref(false)
const createOpen = ref(false)
const createError = ref('')
const tabReady = ref(false)
const movingDealId = ref<string | null>(null)

const createForm = reactive({
  title: '',
  amount: '',
  expected_close_date: '',
})

const mainTabs = computed(() => [
  { id: 'pipeline', label: t('dealsTabPipeline'), icon: 'i-heroicons-queue-list' },
  { id: 'reports', label: t('dealsTabReports'), icon: 'i-heroicons-chart-bar' },
])

const canCreate = computed(() => permission.can('deals', 'create'))
const canUpdate = computed(() => permission.can('deals', 'update'))
const isPreviewMode = computed(() => dealsApi.useMock.value)

async function loadPipeline() {
  pending.value = true
  loadError.value = ''
  try {
    pipeline.value = await dealsApi.fetchPipeline()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
    pipeline.value = null
  } finally {
    pending.value = false
  }
}

async function onMoveStage(dealId: string, stage: DealStage) {
  if (!canUpdate.value) return
  movingDealId.value = dealId
  loadError.value = ''
  try {
    await dealsApi.updateStage(dealId, { stage })
    await loadPipeline()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    movingDealId.value = null
  }
}

async function submitCreate() {
  createError.value = ''
  const title = createForm.title.trim()
  if (!title) {
    createError.value = t('dealsTitleRequired')
    return
  }
  creating.value = true
  try {
    await dealsApi.create({
      title,
      amount: createForm.amount ? Number(createForm.amount) : undefined,
      expected_close_date: createForm.expected_close_date || null,
    })
    createOpen.value = false
    createForm.title = ''
    createForm.amount = ''
    createForm.expected_close_date = ''
    await loadPipeline()
  } catch (e) {
    createError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    creating.value = false
  }
}

function applyRouteQuery() {
  const tab = route.query.tab
  if (tab === 'reports') activeTab.value = 'reports'
  else if (tab === 'pipeline') activeTab.value = 'pipeline'
  if (route.query.create === '1' && canCreate.value) {
    createOpen.value = true
  }
}

watch(activeTab, (tab) => {
  if (!tabReady.value) return
  const query = { ...route.query }
  if (tab === 'reports') query.tab = 'reports'
  else delete query.tab
  router.replace({ query })
})

function onKeydown(e: KeyboardEvent) {
  const target = e.target as HTMLElement | null
  if (target && ['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName)) return
  if (target?.isContentEditable) return
  if ((e.key === 'n' || e.key === 'N') && !e.metaKey && !e.ctrlKey && !e.altKey) {
    if (canCreate.value && !createOpen.value) {
      e.preventDefault()
      createOpen.value = true
    }
  }
}

onMounted(async () => {
  applyRouteQuery()
  tabReady.value = true
  if (typeof globalThis.window !== 'undefined') {
    globalThis.window.addEventListener('keydown', onKeydown)
  }
  await loadPipeline()
})

onUnmounted(() => {
  if (typeof globalThis.window !== 'undefined') {
    globalThis.window.removeEventListener('keydown', onKeydown)
  }
})
</script>
