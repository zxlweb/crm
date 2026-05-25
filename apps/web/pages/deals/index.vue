<template>
  <PermissionGuard resource="deals" action="view">
    <div class="deals-page relative" data-testid="deals-page">
      <div
        class="pointer-events-none absolute inset-x-0 top-0 h-56 overflow-hidden"
        aria-hidden="true"
      >
        <div
          class="absolute -left-12 top-0 h-44 w-44 rounded-full blur-3xl opacity-60"
          :style="{ background: 'var(--ds-blur-brand)' }"
        />
        <div
          class="absolute right-0 top-4 h-36 w-36 rounded-full blur-3xl opacity-50"
          :style="{ background: 'var(--ds-blur-accent)' }"
        />
      </div>

      <div class="relative mx-auto max-w-[1400px] space-y-6">
        <header class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
          <div class="min-w-0">
            <p class="text-xs font-medium uppercase tracking-wide text-ds-fg-brand">
              {{ $t('dealsPageLabel') }}
            </p>
            <h1 class="mt-0.5 text-xl font-bold tracking-tight text-ds-fg-heading sm:text-2xl">
              {{ $t('dealsPageTitle') }}
            </h1>
            <p class="mt-1.5 max-w-2xl text-sm text-ds-fg-muted">
              {{ isPreviewMode ? $t('dealsPageDescMock') : $t('dealsPageDesc') }}
            </p>
          </div>
          <UiButton
            v-if="canCreate"
            variant="secondary"
            size="sm"
            class="shrink-0"
            icon="i-heroicons-plus-20-solid"
            data-testid="deal-create-btn"
            :loading="creating"
            @click="createOpen = true"
          >
            {{ $t('dealsCreate') }}
          </UiButton>
        </header>

        <DealsPipelineSkeleton v-if="pending && !pipeline" />

        <template v-else-if="pipeline">
          <DealsKpiRow :summary="pipeline.summary" :stages="pipeline.stages" />

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
              <DealsPipelineBoard
                :stages="pipeline.stages"
                :pending="pending"
                :can-update="canUpdate"
                :moving-deal-id="movingDealId"
                @move-stage="onMoveStage"
              >
                <template v-if="canCreate" #empty-action>
                  <UiButton
                    variant="secondary"
                    size="sm"
                    class="mt-4"
                    icon="i-heroicons-plus-20-solid"
                    @click="createOpen = true"
                  >
                    {{ $t('dealsCreate') }}
                  </UiButton>
                </template>
              </DealsPipelineBoard>
            </div>
            <DealsReportsPanel v-else key="reports" />
          </Transition>
        </template>

        <UAlert v-else-if="loadError" color="red" variant="soft" :title="loadError" />
      </div>

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
  { id: 'pipeline', label: t('dealsTabPipeline') },
  { id: 'reports', label: t('dealsTabReports') },
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

onMounted(async () => {
  applyRouteQuery()
  tabReady.value = true
  await loadPipeline()
})
</script>
