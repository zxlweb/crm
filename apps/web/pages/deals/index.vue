<template>
  <PermissionGuard resource="deals" action="view">
    <div class="deals-page relative" data-testid="deals-page">
      <!-- Ambient brand glow -->
      <div
        class="pointer-events-none absolute inset-x-0 top-0 h-72 overflow-hidden"
        aria-hidden="true"
      >
        <div
          class="absolute -left-16 top-0 h-60 w-60 rounded-full blur-3xl opacity-70"
          :style="{ background: 'var(--ds-blur-brand)' }"
        />
        <div
          class="absolute right-0 top-8 h-44 w-44 rounded-full blur-3xl opacity-60"
          :style="{ background: 'var(--ds-blur-accent)' }"
        />
      </div>

      <div class="relative mx-auto max-w-[1400px] space-y-6">
        <!-- HERO COMMAND BAR -->
        <header
          class="ds-deals-hero relative overflow-hidden rounded-3xl border border-ds-border bg-ds-bg-elevated/85 px-5 py-5 shadow-ds-md backdrop-blur-md sm:px-7 sm:py-7"
          data-testid="deals-hero"
        >
          <div
            class="pointer-events-none absolute -right-12 -top-12 h-56 w-56 rounded-full opacity-30 blur-3xl"
            :style="{ background: 'var(--ds-brand-gradient)' }"
            aria-hidden="true"
          />
          <div
            class="pointer-events-none absolute -bottom-16 left-1/3 h-40 w-40 rounded-full opacity-20 blur-3xl"
            :style="{ background: 'var(--ds-blur-brand)' }"
            aria-hidden="true"
          />

          <div class="relative flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-[0.18em] text-ds-fg-brand">
                <span class="relative flex h-2 w-2" aria-hidden="true">
                  <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-ds-brand opacity-50" />
                  <span class="relative inline-flex h-2 w-2 rounded-full bg-ds-brand" />
                </span>
                <span>{{ $t('dealsPageLabel') }}</span>
                <span aria-hidden="true" class="text-ds-fg-subtle">·</span>
                <span class="text-ds-fg-muted">{{ isPreviewMode ? $t('dealsHeroDemoBadge') : $t('dealsHeroLiveBadge') }}</span>
              </div>
              <h1 class="mt-1.5 text-2xl font-bold tracking-tight text-ds-fg-heading sm:text-3xl">
                {{ $t('dealsPageTitle') }}
              </h1>

              <!-- Headline pipeline value -->
              <div class="mt-4 flex flex-wrap items-end gap-x-4 gap-y-2">
                <div>
                  <p class="text-[11px] font-medium uppercase tracking-wider text-ds-fg-muted">
                    {{ $t('dealsHeroTotalLabel') }}
                  </p>
                  <p class="mt-0.5 text-3xl font-extrabold tabular-nums tracking-tight text-ds-fg-heading sm:text-4xl">
                    <span class="bg-clip-text" :style="brandGradientText">{{ heroPipelineValue }}</span>
                  </p>
                </div>
                <div class="flex flex-wrap items-center gap-1.5">
                  <span
                    class="inline-flex items-center gap-1.5 rounded-full border border-ds-info/25 bg-ds-info-subtle px-2.5 py-1 text-xs font-medium text-ds-info"
                  >
                    <UIcon name="i-heroicons-briefcase" class="h-3.5 w-3.5" aria-hidden="true" />
                    {{ $t('dealsHeroOpenDeals', { n: openDealsCount }) }}
                  </span>
                  <span
                    v-if="closingThisWeekCount > 0"
                    class="inline-flex items-center gap-1.5 rounded-full border border-ds-warning/25 bg-ds-warning-subtle px-2.5 py-1 text-xs font-medium text-ds-warning"
                  >
                    <UIcon name="i-heroicons-bell-alert" class="h-3.5 w-3.5" aria-hidden="true" />
                    {{ $t('dealsHeroClosingWeek', { n: closingThisWeekCount }) }}
                  </span>
                  <span
                    v-if="overdueCount > 0"
                    class="inline-flex items-center gap-1.5 rounded-full border border-ds-danger/25 bg-ds-danger-subtle px-2.5 py-1 text-xs font-medium text-ds-danger"
                  >
                    <UIcon name="i-heroicons-exclamation-triangle" class="h-3.5 w-3.5" aria-hidden="true" />
                    {{ $t('dealsHeroOverdue', { n: overdueCount }) }}
                  </span>
                  <span
                    v-if="wonThisMonthCount > 0"
                    class="inline-flex items-center gap-1.5 rounded-full border border-ds-success/25 bg-ds-success-subtle px-2.5 py-1 text-xs font-medium text-ds-success"
                  >
                    <UIcon name="i-heroicons-trophy" class="h-3.5 w-3.5" aria-hidden="true" />
                    {{ $t('dealsHeroWonMonth', { n: wonThisMonthCount }) }}
                  </span>
                </div>
              </div>

              <p class="mt-3 max-w-xl text-sm text-ds-fg-muted">
                {{ isPreviewMode ? $t('dealsPageDescMock') : $t('dealsPageDesc') }}
              </p>
            </div>

            <div class="flex shrink-0 flex-wrap items-center gap-2 lg:flex-col lg:items-end">
              <button
                v-if="canCreate"
                type="button"
                class="ds-deals-hero__cta group relative inline-flex cursor-pointer items-center gap-2 overflow-hidden rounded-2xl px-5 py-3 text-sm font-semibold text-ds-on-brand shadow-ds-brand transition-[transform,box-shadow] duration-200 hover:shadow-ds-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ds-brand focus-visible:ring-offset-2 focus-visible:ring-offset-ds-bg disabled:cursor-not-allowed disabled:opacity-60"
                :style="{ background: 'var(--ds-brand-gradient)' }"
                data-testid="deal-create-btn"
                :disabled="creating"
                @click="createOpen = true"
              >
                <span
                  class="pointer-events-none absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/30 to-transparent opacity-0 transition-[transform,opacity] duration-500 group-hover:translate-x-full group-hover:opacity-100"
                  aria-hidden="true"
                />
                <UIcon name="i-heroicons-plus" class="h-4 w-4" aria-hidden="true" />
                <span>{{ $t('dealsCreate') }}</span>
              </button>
              <div class="hidden items-center gap-2 text-xs text-ds-fg-muted lg:flex">
                <kbd class="rounded border border-ds-border bg-ds-bg-muted px-1.5 py-0.5 font-mono text-[10px] text-ds-fg-muted">N</kbd>
                <span>{{ $t('dealsHeroShortcut') }}</span>
              </div>
            </div>
          </div>
        </header>

        <DealsPipelineSkeleton v-if="pending && !pipeline" />

        <template v-else-if="pipeline">
          <DealsKpiRow :summary="pipeline.summary" :stages="pipeline.stages" />

          <DealsFocusStream :stages="pipeline.stages" />

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
const { formatDealAmount } = useDealLabels()

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
  { id: 'pipeline', label: t('dealsTabPipeline'), icon: 'i-heroicons-rectangle-stack' },
  { id: 'reports', label: t('dealsTabReports'), icon: 'i-heroicons-chart-bar' },
])

const canCreate = computed(() => permission.can('deals', 'create'))
const canUpdate = computed(() => permission.can('deals', 'update'))
const isPreviewMode = computed(() => dealsApi.useMock.value)

const brandGradientText = {
  backgroundImage: 'var(--ds-brand-gradient)',
  WebkitBackgroundClip: 'text',
  WebkitTextFillColor: 'transparent',
  backgroundClip: 'text',
  color: 'transparent',
}

const heroPipelineValue = computed(() => {
  const summary = pipeline.value?.summary
  if (!summary) return formatDealAmount(0)
  return formatDealAmount(summary.open_amount)
})

const openDealsCount = computed(() => pipeline.value?.summary.open_count ?? 0)
const wonThisMonthCount = computed(() => pipeline.value?.summary.won_count_mtd ?? 0)

function startOfDay(d: Date): number {
  const x = new Date(d)
  x.setHours(0, 0, 0, 0)
  return x.getTime()
}

const closingThisWeekCount = computed(() => {
  if (!pipeline.value) return 0
  const today = startOfDay(new Date())
  let n = 0
  for (const col of pipeline.value.stages) {
    if (col.stage === 'won' || col.stage === 'lost') continue
    for (const item of col.items) {
      if (!item.expected_close_date) continue
      const close = startOfDay(new Date(`${item.expected_close_date}T00:00:00`))
      const days = Math.round((close - today) / 86_400_000)
      if (days >= 0 && days <= 7) n++
    }
  }
  return n
})

const overdueCount = computed(() => {
  if (!pipeline.value) return 0
  const today = startOfDay(new Date())
  let n = 0
  for (const col of pipeline.value.stages) {
    if (col.stage === 'won' || col.stage === 'lost') continue
    for (const item of col.items) {
      if (!item.expected_close_date) continue
      const close = startOfDay(new Date(`${item.expected_close_date}T00:00:00`))
      if (close < today) n++
    }
  }
  return n
})

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
