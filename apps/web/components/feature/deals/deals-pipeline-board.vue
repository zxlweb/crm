<template>
  <div class="space-y-4" data-testid="deals-pipeline">
    <!-- Toolbar: search + inline pipeline ribbon + funnel toggle -->
    <div
      class="ds-pipeline-toolbar flex flex-col gap-3 rounded-2xl border border-ds-border-muted bg-ds-bg-elevated/85 p-3 shadow-ds-sm backdrop-blur-sm lg:flex-row lg:items-center"
    >
      <UiInput
        v-model="searchQuery"
        search
        type="search"
        class="max-w-md flex-1"
        :placeholder="$t('dealsSearchPlaceholder')"
        data-testid="deals-pipeline-search"
      />

      <!-- Inline pipeline ribbon (always visible, low-noise flow visualization) -->
      <div
        v-if="ribbonSegments.length > 0"
        class="ds-pipeline-ribbon flex min-w-0 flex-1 items-center gap-3"
        data-testid="deals-pipeline-ribbon"
      >
        <div
          class="relative h-2.5 min-w-0 flex-1 overflow-hidden rounded-full bg-ds-bg-muted ring-1 ring-inset ring-ds-border-muted"
          :aria-label="$t('dealsPipelineRibbonAria')"
        >
          <span
            v-for="seg in ribbonSegments"
            :key="seg.stage"
            class="absolute inset-y-0 transition-[left,width,background-color] duration-500 ease-out"
            :class="seg.barClass"
            :style="{ left: `${seg.offset}%`, width: `${seg.width}%` }"
            :title="seg.tooltip"
          />
        </div>
        <p class="hidden whitespace-nowrap text-[11px] font-medium tabular-nums text-ds-fg-muted xl:block">
          {{ $t('dealsPipelineRibbonTotal', { value: formatDealAmount(ribbonTotal) }) }}
        </p>
      </div>

      <div class="flex shrink-0 flex-wrap items-center gap-2">
        <p
          v-if="canUpdate && !searchActive"
          class="hidden text-xs text-ds-fg-subtle xl:block"
        >
          {{ $t('dealsDragHint') }}
        </p>
        <button
          type="button"
          class="ds-pipeline-funnel-toggle inline-flex cursor-pointer items-center gap-1.5 rounded-xl border px-3 py-2 text-xs font-medium transition-colors duration-200"
          :class="funnelOpen
            ? 'border-ds-brand bg-ds-brand-subtle text-ds-fg-brand shadow-ds-sm'
            : 'border-ds-border bg-ds-bg-elevated text-ds-fg-muted hover:border-ds-brand-muted hover:bg-ds-bg-muted hover:text-ds-fg'"
          :aria-expanded="funnelOpen"
          aria-controls="deals-funnel-panel"
          data-testid="deals-funnel-toggle"
          @click="funnelOpen = !funnelOpen"
        >
          <UIcon name="i-heroicons-funnel" class="h-3.5 w-3.5" aria-hidden="true" />
          <span>{{ funnelOpen ? $t('dealsHideFunnel') : $t('dealsShowFunnel') }}</span>
          <UIcon
            name="i-heroicons-chevron-down"
            class="h-3.5 w-3.5 transition-transform duration-200"
            :class="{ 'rotate-180': funnelOpen }"
            aria-hidden="true"
          />
        </button>
      </div>
    </div>

    <!-- Funnel detail panel — styled to match new card language -->
    <Transition
      enter-active-class="overflow-hidden transition-[max-height,opacity] duration-300 ease-out"
      enter-from-class="max-h-0 opacity-0"
      enter-to-class="max-h-[420px] opacity-100"
      leave-active-class="overflow-hidden transition-[max-height,opacity] duration-200 ease-in"
      leave-from-class="max-h-[420px] opacity-100"
      leave-to-class="max-h-0 opacity-0"
    >
      <section
        v-if="funnelOpen"
        id="deals-funnel-panel"
        class="ds-pipeline-funnel relative overflow-hidden rounded-2xl border border-ds-border-muted bg-ds-bg-elevated shadow-ds-sm"
        data-testid="deals-pipeline-funnel"
      >
        <span
          class="pointer-events-none absolute inset-x-0 top-0 h-0.5"
          :style="{ background: 'var(--ds-brand-gradient)' }"
          aria-hidden="true"
        />
        <span
          class="pointer-events-none absolute -right-10 -top-10 h-40 w-40 rounded-full opacity-15 blur-3xl"
          :style="{ background: 'var(--ds-brand-gradient)' }"
          aria-hidden="true"
        />
        <header class="flex items-start justify-between gap-3 px-5 pb-3 pt-4">
          <div class="min-w-0">
            <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('dealsPipelineFunnelTitle') }}</h3>
            <p class="mt-0.5 text-xs text-ds-fg-muted">{{ $t('dealsPipelineFunnelHint') }}</p>
          </div>
          <button
            type="button"
            class="inline-flex h-7 w-7 cursor-pointer items-center justify-center rounded-lg text-ds-fg-muted transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-heading"
            :aria-label="$t('dealsHideFunnel')"
            @click="funnelOpen = false"
          >
            <UIcon name="i-heroicons-x-mark" class="h-4 w-4" aria-hidden="true" />
          </button>
        </header>
        <div class="px-3 pb-3">
          <LeadsReportChartSlot
            :pending="pending"
            :empty="!funnelItems.length"
            :height="220"
            :empty-text="$t('dealsPipelineEmpty')"
          >
            <ChartFunnel :key="chartLocale" :items="funnelItems" :height="220" />
          </LeadsReportChartSlot>
        </div>
      </section>
    </Transition>

    <div
      v-if="isPipelineEmpty && !pending"
      class="rounded-xl border border-dashed border-ds-border bg-ds-bg-muted/30 px-6 py-16 text-center"
      data-testid="deals-pipeline-empty"
    >
      <svg class="mx-auto h-10 w-10 text-ds-fg-subtle" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="mt-3 text-sm font-medium text-ds-fg-heading">{{ $t('dealsPipelineEmpty') }}</p>
      <p class="mt-1 text-xs text-ds-fg-muted">{{ $t('dealsPipelineEmptyHint') }}</p>
      <slot name="empty-action" />
    </div>

    <div
      v-else
      class="ds-pipeline-track relative flex items-stretch gap-3 overflow-x-auto pb-3 pt-1 snap-x snap-mandatory scroll-px-4"
      data-testid="deals-pipeline-columns"
    >
      <template v-for="(column, idx) in filteredStages" :key="column.stage">
        <section
          class="ds-pipeline-col group/col relative flex w-[min(100%,18.5rem)] shrink-0 snap-start flex-col overflow-hidden rounded-2xl border bg-ds-bg-elevated shadow-ds-sm transition-[border-color,background-color,box-shadow,transform] duration-200"
          :class="columnDropClass(column.stage)"
          :data-testid="`deals-pipeline-stage-${column.stage}`"
          @dragover="onColumnDragOver($event, column.stage)"
          @dragleave="onColumnDragLeave(column.stage)"
          @drop="onColumnDrop($event, column.stage)"
        >
          <!-- Stage accent bar -->
          <span
            class="pointer-events-none absolute inset-x-0 top-0 h-1"
            :class="stageBarClass(column.stage)"
            aria-hidden="true"
          />

          <header
            class="border-b border-ds-border-muted px-3 pb-3 pt-3.5"
            :class="stageHeaderClass(column.stage)"
          >
            <div class="flex items-center justify-between gap-2">
              <DealsDealStageBadge :stage="column.stage" variant="plain" />
              <span
                class="inline-flex items-center justify-center rounded-full bg-ds-bg-elevated/95 px-2 py-0.5 text-xs font-semibold tabular-nums text-ds-fg-heading ring-1 ring-inset ring-ds-border-muted"
              >
                {{ column.count }}
              </span>
            </div>
            <div class="mt-2.5 flex items-end justify-between gap-2">
              <p class="text-sm font-bold tabular-nums tracking-tight text-ds-fg-heading">
                {{ formatDealAmount(column.amount_total) }}
              </p>
              <p
                v-if="stageShare(column) !== null"
                class="text-[10px] font-semibold uppercase tracking-wider text-ds-fg-muted"
              >
                {{ $t('dealsStageShare', { value: stageShare(column) }) }}
              </p>
            </div>
            <div
              v-if="openPipelineTotal > 0 && column.stage !== 'won' && column.stage !== 'lost'"
              class="mt-2 h-1 overflow-hidden rounded-full bg-ds-bg-muted"
              role="presentation"
            >
              <div
                class="h-full rounded-full transition-all duration-500"
                :class="stageBarClass(column.stage)"
                :style="{ width: `${Math.round((column.amount_total / openPipelineTotal) * 100)}%` }"
              />
            </div>
          </header>

          <div
            class="flex min-h-[8rem] flex-1 flex-col gap-2 p-2 transition-colors duration-200"
            :class="columnBodyClass(column.stage)"
          >
            <DealsDealPipelineCard
              v-for="item in column.items"
              :key="item.id"
              :item="item"
              :stage="column.stage"
              :can-update="canUpdate"
              :draggable="dragEnabled"
              :dragging="draggingDealId === item.id"
              :moving="movingDealId === item.id"
              @move-stage="(next) => $emit('move-stage', item.id, next)"
              @drag-start="onCardDragStart"
              @drag-end="onCardDragEnd"
            />
            <p
              v-if="!column.items.length"
              class="rounded-lg border border-dashed px-3 py-8 text-center text-xs transition-colors duration-200"
              :class="emptySlotClass(column.stage)"
            >
              {{ emptyColumnText(column) }}
            </p>
          </div>
        </section>

        <!-- Flow connector between open stages -->
        <div
          v-if="idx < filteredStages.length - 1 && isFlowConnector(column.stage, filteredStages[idx + 1].stage)"
          class="ds-pipeline-flow pointer-events-none hidden shrink-0 self-stretch items-center xl:flex"
          aria-hidden="true"
        >
          <div class="flex flex-col items-center gap-1.5">
            <UIcon
              name="i-heroicons-chevron-double-right"
              class="h-4 w-4 text-ds-fg-subtle"
            />
            <span
              v-if="conversionBetween(column.stage, filteredStages[idx + 1].stage) !== null"
              class="rounded-full border border-ds-border-muted bg-ds-bg-elevated px-1.5 py-0.5 text-[10px] font-semibold tabular-nums text-ds-fg-muted shadow-ds-sm"
            >
              {{ conversionBetween(column.stage, filteredStages[idx + 1].stage) }}%
            </span>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ChartFunnelItem } from '@crm/ui-kit'
import type { DealPipelineStage, DealStage } from '~/types/deal'
import { canTransitionDealStage, readDealDragData } from '~/utils/deal-stage-transition'

const props = defineProps<{
  stages: DealPipelineStage[]
  pending?: boolean
  canUpdate?: boolean
  movingDealId?: string | null
}>()

const emit = defineEmits<{
  'move-stage': [dealId: string, stage: DealStage]
}>()

const { t, locale } = useI18n()
const { dealStageLabel, formatDealAmount } = useDealLabels()

const searchQuery = ref('')
const funnelOpen = ref(false)
const draggingDealId = ref<string | null>(null)
const dragFromStage = ref<DealStage | null>(null)
const dragOverStage = ref<DealStage | null>(null)

const chartLocale = computed(() => locale.value)
const searchActive = computed(() => searchQuery.value.trim().length > 0)
const dragEnabled = computed(() => Boolean(props.canUpdate) && !searchActive.value)

const openPipelineTotal = computed(() =>
  props.stages
    .filter((s) => s.stage !== 'won' && s.stage !== 'lost')
    .reduce((sum, s) => sum + s.amount_total, 0),
)

const normalizedSearch = computed(() => searchQuery.value.trim().toLowerCase())

const filteredStages = computed(() => {
  const q = normalizedSearch.value
  if (!q) return props.stages
  return props.stages.map((column) => ({
    ...column,
    items: column.items.filter((item) => item.title.toLowerCase().includes(q)),
    count: column.items.filter((item) => item.title.toLowerCase().includes(q)).length,
    amount_total: column.items
      .filter((item) => item.title.toLowerCase().includes(q))
      .reduce((sum, item) => sum + item.amount, 0),
  }))
})

const isPipelineEmpty = computed(() =>
  props.stages.every((s) => s.count === 0),
)

const funnelItems = computed<ChartFunnelItem[]>(() => {
  void locale.value
  return props.stages
    .filter((s) => s.count > 0)
    .map((s) => ({
      name: dealStageLabel(s.stage),
      value: s.count,
    }))
})

// Inline pipeline ribbon — always-visible, low-noise flow visualization
const ribbonStageOrder: DealStage[] = ['qualification', 'proposal', 'negotiation', 'won']

const ribbonTotal = computed(() =>
  props.stages
    .filter((s) => ribbonStageOrder.includes(s.stage))
    .reduce((sum, s) => sum + s.amount_total, 0),
)

interface RibbonSegment {
  stage: DealStage
  offset: number
  width: number
  barClass: string
  tooltip: string
}

const ribbonSegments = computed<RibbonSegment[]>(() => {
  if (ribbonTotal.value === 0) return []
  const segments: RibbonSegment[] = []
  let cursor = 0
  for (const stage of ribbonStageOrder) {
    const col = props.stages.find((s) => s.stage === stage)
    if (!col || col.amount_total === 0) continue
    const width = (col.amount_total / ribbonTotal.value) * 100
    segments.push({
      stage,
      offset: cursor,
      width,
      barClass: stageBarClass(stage),
      tooltip: `${dealStageLabel(stage)} · ${formatDealAmount(col.amount_total)} · ${col.count}`,
    })
    cursor += width
  }
  return segments
})

function isValidDropTarget(toStage: DealStage): boolean {
  if (!dragFromStage.value) return false
  return canTransitionDealStage(dragFromStage.value, toStage)
}

function columnDropClass(stage: DealStage): string {
  if (dragOverStage.value !== stage || !draggingDealId.value) {
    return 'border-ds-border-muted hover:border-ds-border'
  }
  if (isValidDropTarget(stage)) {
    return 'border-ds-brand shadow-ds-md ring-2 ring-ds-brand/25 -translate-y-0.5'
  }
  return 'border-ds-danger/50 ring-2 ring-ds-danger/20'
}

function columnBodyClass(stage: DealStage): string {
  if (dragOverStage.value === stage && draggingDealId.value && isValidDropTarget(stage)) {
    return 'bg-ds-brand-subtle/10'
  }
  return ''
}

function emptySlotClass(stage: DealStage): string {
  if (dragOverStage.value === stage && draggingDealId.value) {
    return isValidDropTarget(stage)
      ? 'border-ds-brand/50 bg-ds-brand-subtle/20 text-ds-fg-brand'
      : 'border-ds-danger/30 bg-ds-danger-subtle/10 text-ds-danger'
  }
  return 'border-ds-border/80 bg-ds-bg/50 text-ds-fg-muted'
}

function emptyColumnText(column: DealPipelineStage): string {
  if (dragOverStage.value === column.stage && draggingDealId.value) {
    return isValidDropTarget(column.stage)
      ? t('dealsDropHere')
      : t('dealsDropInvalid')
  }
  if (searchActive.value && !column.items.length) return t('dealsSearchNoResults')
  return t('dealsStageEmpty')
}

function onCardDragStart(payload: { dealId: string; fromStage: DealStage }) {
  draggingDealId.value = payload.dealId
  dragFromStage.value = payload.fromStage
}

function onCardDragEnd() {
  draggingDealId.value = null
  dragFromStage.value = null
  dragOverStage.value = null
}

function onColumnDragOver(event: DragEvent, stage: DealStage) {
  if (!dragEnabled.value || !draggingDealId.value) return
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = isValidDropTarget(stage) ? 'move' : 'none'
  }
  dragOverStage.value = stage
}

function onColumnDragLeave(stage: DealStage) {
  if (dragOverStage.value === stage) {
    dragOverStage.value = null
  }
}

function onColumnDrop(event: DragEvent, toStage: DealStage) {
  event.preventDefault()
  const payload = event.dataTransfer ? readDealDragData(event.dataTransfer) : null
  onCardDragEnd()
  if (!payload || !dragEnabled.value) return
  if (payload.fromStage === toStage) return
  if (!canTransitionDealStage(payload.fromStage, toStage)) return
  emit('move-stage', payload.dealId, toStage)
}

function stageHeaderClass(stage: DealStage): string {
  const map: Record<DealStage, string> = {
    qualification: 'bg-gradient-to-b from-ds-info-subtle/50 to-transparent',
    proposal: 'bg-gradient-to-b from-ds-brand-subtle/50 to-transparent',
    negotiation: 'bg-gradient-to-b from-ds-warning-subtle/50 to-transparent',
    won: 'bg-gradient-to-b from-ds-success-subtle/60 to-transparent',
    lost: 'bg-gradient-to-b from-ds-bg-muted/60 to-transparent',
  }
  return map[stage]
}

function stageBarClass(stage: DealStage): string {
  const map: Record<DealStage, string> = {
    qualification: 'bg-ds-info',
    proposal: 'bg-ds-brand',
    negotiation: 'bg-ds-warning',
    won: 'bg-ds-success',
    lost: 'bg-ds-fg-subtle',
  }
  return map[stage]
}

function stageShare(column: DealPipelineStage): number | null {
  if (column.stage === 'won' || column.stage === 'lost') return null
  if (openPipelineTotal.value === 0) return null
  return Math.round((column.amount_total / openPipelineTotal.value) * 100)
}

function isFlowConnector(from: DealStage, to: DealStage): boolean {
  const openOrder: DealStage[] = ['qualification', 'proposal', 'negotiation']
  const fromIdx = openOrder.indexOf(from)
  const toIdx = openOrder.indexOf(to)
  return fromIdx >= 0 && toIdx === fromIdx + 1
}

function conversionBetween(from: DealStage, to: DealStage): number | null {
  const fromCol = props.stages.find((s) => s.stage === from)
  const toCol = props.stages.find((s) => s.stage === to)
  if (!fromCol || !toCol) return null
  const total = fromCol.count + toCol.count
  if (total === 0) return null
  return Math.round((toCol.count / total) * 100)
}
</script>
