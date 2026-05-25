<template>
  <div class="space-y-5" data-testid="deals-pipeline">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <UiInput
        v-model="searchQuery"
        search
        type="search"
        class="max-w-md flex-1"
        :placeholder="$t('dealsSearchPlaceholder')"
        data-testid="deals-pipeline-search"
      />
      <div class="flex flex-wrap items-center gap-2">
        <p
          v-if="canUpdate && !searchActive"
          class="hidden text-xs text-ds-fg-subtle sm:block"
        >
          {{ $t('dealsDragHint') }}
        </p>
        <button
          type="button"
          class="inline-flex cursor-pointer items-center gap-1.5 rounded-lg border border-ds-border px-3 py-2 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand-muted hover:bg-ds-bg-muted hover:text-ds-fg"
          data-testid="deals-funnel-toggle"
          @click="funnelOpen = !funnelOpen"
        >
          <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
          </svg>
          {{ funnelOpen ? $t('dealsHideFunnel') : $t('dealsShowFunnel') }}
        </button>
      </div>
    </div>

    <Transition
      enter-active-class="transition-all duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-1"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition-all duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-1"
    >
      <CardShell
        v-if="funnelOpen"
        :title="$t('dealsPipelineFunnelTitle')"
        :subtitle="$t('dealsPipelineFunnelHint')"
        :height="320"
        data-testid="deals-pipeline-funnel"
      >
        <LeadsReportChartSlot
          :pending="pending"
          :empty="!funnelItems.length"
          :height="280"
          :empty-text="$t('dealsPipelineEmpty')"
        >
          <ChartFunnel :key="chartLocale" :items="funnelItems" :height="260" />
        </LeadsReportChartSlot>
      </CardShell>
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
      class="flex gap-4 overflow-x-auto pb-2 snap-x snap-mandatory scroll-px-4"
      data-testid="deals-pipeline-columns"
    >
      <section
        v-for="column in filteredStages"
        :key="column.stage"
        class="flex w-[min(100%,18rem)] shrink-0 snap-start flex-col rounded-xl border shadow-ds-sm transition-[border-color,background-color,box-shadow] duration-200"
        :class="columnDropClass(column.stage)"
        :data-testid="`deals-pipeline-stage-${column.stage}`"
        @dragover="onColumnDragOver($event, column.stage)"
        @dragleave="onColumnDragLeave(column.stage)"
        @drop="onColumnDrop($event, column.stage)"
      >
        <header
          class="rounded-t-xl border-b border-ds-border/60 px-3 py-3"
          :class="stageHeaderClass(column.stage)"
        >
          <div class="flex items-center justify-between gap-2">
            <DealsDealStageBadge :stage="column.stage" variant="plain" />
            <span
              class="rounded-full bg-ds-bg/80 px-2 py-0.5 text-xs font-semibold tabular-nums text-ds-fg-heading"
            >
              {{ column.count }}
            </span>
          </div>
          <p class="mt-2 text-sm font-semibold tabular-nums text-ds-fg-heading">
            {{ formatDealAmount(column.amount_total) }}
          </p>
          <div
            v-if="openPipelineTotal > 0 && column.stage !== 'won' && column.stage !== 'lost'"
            class="mt-2 h-1 overflow-hidden rounded-full bg-ds-bg/60"
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
const funnelOpen = ref(true)
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

function isValidDropTarget(toStage: DealStage): boolean {
  if (!dragFromStage.value) return false
  return canTransitionDealStage(dragFromStage.value, toStage)
}

function columnDropClass(stage: DealStage): string {
  if (dragOverStage.value !== stage || !draggingDealId.value) {
    return 'border-ds-border/80 bg-ds-bg-muted/20'
  }
  if (isValidDropTarget(stage)) {
    return 'border-ds-brand bg-ds-brand-subtle/25 shadow-ds-md ring-2 ring-ds-brand/20'
  }
  return 'border-ds-danger/40 bg-ds-danger-subtle/20 ring-2 ring-ds-danger/15'
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
    qualification: 'bg-ds-info-subtle/30',
    proposal: 'bg-ds-brand-subtle/30',
    negotiation: 'bg-ds-warning-subtle/30',
    won: 'bg-ds-success-subtle/30',
    lost: 'bg-ds-bg-muted/50',
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
</script>
