<template>
  <article
    class="group relative flex min-h-[108px] flex-col rounded-xl border border-ds-border/80 bg-ds-bg p-3 shadow-ds-sm transition-[box-shadow,border-color,opacity,transform] duration-200 hover:border-ds-brand-muted/40 hover:shadow-ds-md"
    :class="[
      draggable ? 'cursor-grab touch-none active:cursor-grabbing' : '',
      dragging ? 'scale-[0.98] border-ds-brand-muted opacity-60 shadow-none' : '',
      moving ? 'pointer-events-none opacity-50' : '',
    ]"
    :draggable="draggable"
    :data-testid="`deal-pipeline-card-${item.id}`"
    @dragstart="onDragStart"
    @dragend="onDragEnd"
  >
    <span
      v-if="draggable"
      class="pointer-events-none absolute right-2.5 top-2.5 text-ds-fg-subtle opacity-0 transition-opacity duration-200 group-hover:opacity-70"
      aria-hidden="true"
    >
      <svg class="h-4 w-4" viewBox="0 0 24 24" fill="currentColor">
        <circle cx="9" cy="7" r="1.25" />
        <circle cx="15" cy="7" r="1.25" />
        <circle cx="9" cy="12" r="1.25" />
        <circle cx="15" cy="12" r="1.25" />
        <circle cx="9" cy="17" r="1.25" />
        <circle cx="15" cy="17" r="1.25" />
      </svg>
    </span>

    <NuxtLink
      :to="`/deals/${item.id}`"
      class="cursor-pointer pr-5 text-sm font-semibold leading-snug text-ds-fg-heading underline-offset-2 transition-colors duration-200 group-hover:text-ds-fg-brand hover:underline"
      @click="onLinkClick"
    >
      {{ item.title }}
    </NuxtLink>

    <p class="mt-2 text-base font-semibold tabular-nums text-ds-fg-heading">
      {{ formatDealAmount(item.amount, item.currency) }}
    </p>

    <div class="mt-3 space-y-1.5">
      <div class="flex items-center justify-between gap-2 text-xs">
        <span class="text-ds-fg-muted">{{ $t('dealsCardProbability', { value: item.probability }) }}</span>
        <span
          v-if="item.expected_close_date"
          class="tabular-nums"
          :class="closeDateClass"
        >
          {{ formatCloseDate(item.expected_close_date) }}
        </span>
      </div>
      <div class="h-1.5 overflow-hidden rounded-full bg-ds-bg-muted" role="presentation">
        <div
          class="h-full rounded-full transition-all duration-500"
          :class="probabilityBarClass"
          :style="{ width: `${Math.min(100, Math.max(0, item.probability))}%` }"
        />
      </div>
    </div>

    <p
      v-if="draggable && canUpdate && nextStages.length"
      class="mt-2.5 truncate text-[10px] leading-relaxed text-ds-fg-subtle opacity-0 transition-opacity duration-200 group-hover:opacity-100"
      data-testid="deal-drag-hint"
    >
      {{ $t('dealsCardDragTargets', { stages: nextStageLabels }) }}
    </p>

    <div
      v-else-if="canUpdate && nextStages.length"
      class="relative mt-2.5"
      draggable="false"
      @mousedown.stop
      @click.stop
    >
      <button
        type="button"
        class="flex w-full cursor-pointer items-center justify-between gap-2 rounded-lg border border-ds-border/80 bg-ds-bg-muted/40 px-2.5 py-1.5 text-left text-[11px] font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand-muted hover:bg-ds-brand-subtle/40 hover:text-ds-fg-brand"
        data-testid="deal-move-menu-btn"
        :aria-expanded="menuOpen"
        @click="menuOpen = !menuOpen"
      >
        <span>{{ $t('dealsMoveStage') }}</span>
        <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
        </svg>
      </button>
      <div
        v-if="menuOpen"
        class="absolute inset-x-0 top-[calc(100%+4px)] z-[var(--ds-z-dropdown)] overflow-hidden rounded-lg border border-ds-border bg-ds-bg-elevated py-1 shadow-ds-md"
        data-testid="deal-move-menu"
      >
        <button
          v-for="next in nextStages"
          :key="next"
          type="button"
          class="flex w-full cursor-pointer items-center gap-2 px-3 py-2 text-left text-xs text-ds-fg transition-colors duration-200 hover:bg-ds-bg-muted"
          :data-testid="`deal-move-${next}`"
          @click="onMoveStage(next)"
        >
          <DealsDealStageBadge :stage="next" variant="plain" />
        </button>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { DealPipelineItem, DealStage } from '~/types/deal'
import { allowedNextStages, writeDealDragData } from '~/utils/deal-stage-transition'

const props = defineProps<{
  item: DealPipelineItem
  stage: DealStage
  canUpdate?: boolean
  draggable?: boolean
  dragging?: boolean
  moving?: boolean
}>()

const emit = defineEmits<{
  'move-stage': [stage: DealStage]
  'drag-start': [payload: { dealId: string; fromStage: DealStage }]
  'drag-end': []
}>()

const { dealStageLabel, formatDealAmount, formatCloseDate } = useDealLabels()

const didDrag = ref(false)
const menuOpen = ref(false)

const nextStages = computed(() => {
  if (!props.canUpdate) return []
  if (props.stage === 'won' || props.stage === 'lost') return []
  return allowedNextStages(props.stage)
})

const nextStageLabels = computed(() =>
  nextStages.value.map((stage) => dealStageLabel(stage)).join(' · '),
)

const isOverdue = computed(() => {
  if (!props.item.expected_close_date) return false
  if (props.stage === 'won' || props.stage === 'lost') return false
  const close = new Date(`${props.item.expected_close_date}T23:59:59`)
  return close.getTime() < Date.now()
})

const closeDateClass = computed(() =>
  isOverdue.value ? 'font-medium text-ds-danger' : 'text-ds-fg-subtle',
)

const probabilityBarClass = computed(() => {
  if (props.item.probability >= 70) return 'bg-gradient-to-r from-ds-success to-emerald-500'
  if (props.item.probability >= 40) return 'bg-gradient-to-r from-ds-brand to-sky-500'
  if (props.item.probability >= 20) return 'bg-gradient-to-r from-ds-warning to-amber-500'
  return 'bg-ds-fg-subtle'
})

function onDragStart(event: DragEvent) {
  if (!props.draggable || !event.dataTransfer) return
  menuOpen.value = false
  didDrag.value = true
  writeDealDragData(event.dataTransfer, { dealId: props.item.id, fromStage: props.stage })
  emit('drag-start', { dealId: props.item.id, fromStage: props.stage })
}

function onDragEnd() {
  emit('drag-end')
  window.setTimeout(() => {
    didDrag.value = false
  }, 0)
}

function onLinkClick(event: MouseEvent) {
  if (didDrag.value) {
    event.preventDefault()
  }
}

function onMoveStage(stage: DealStage) {
  menuOpen.value = false
  emit('move-stage', stage)
}

function closeMenu() {
  menuOpen.value = false
}

onMounted(() => document.addEventListener('click', closeMenu))
onUnmounted(() => document.removeEventListener('click', closeMenu))
</script>
