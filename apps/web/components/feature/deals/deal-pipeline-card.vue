<template>
  <article
    class="ds-deal-card group relative flex min-h-[124px] flex-col overflow-hidden rounded-xl border border-ds-border-muted bg-ds-bg-elevated p-3 shadow-ds-sm transition-[box-shadow,border-color,opacity,transform] duration-200 hover:border-ds-brand-muted hover:shadow-ds-md"
    :class="[
      draggable ? 'cursor-grab touch-none active:cursor-grabbing' : '',
      dragging ? 'scale-[0.98] border-ds-brand-muted opacity-60 shadow-none' : '',
      moving ? 'pointer-events-none opacity-50' : '',
      isOverdue ? 'ring-1 ring-inset ring-ds-danger/20' : '',
    ]"
    :draggable="draggable"
    :data-testid="`deal-pipeline-card-${item.id}`"
    @dragstart="onDragStart"
    @dragend="onDragEnd"
  >
    <!-- Overdue accent strip -->
    <span
      v-if="isOverdue"
      class="pointer-events-none absolute inset-y-0 left-0 w-0.5 bg-ds-danger"
      aria-hidden="true"
    />

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

    <div class="flex items-start gap-2 pr-5">
      <span
        class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-[10px] font-bold uppercase ring-1 ring-inset ring-ds-border-muted"
        :class="avatarClass"
        :title="accountLabel"
        aria-hidden="true"
      >
        {{ avatarInitials }}
      </span>
      <NuxtLink
        :to="`/deals/${item.id}`"
        class="min-w-0 flex-1 cursor-pointer text-sm font-semibold leading-snug text-ds-fg-heading underline-offset-2 transition-colors duration-200 group-hover:text-ds-fg-brand hover:underline"
        @click="onLinkClick"
      >
        <span class="line-clamp-2">{{ item.title }}</span>
      </NuxtLink>
    </div>

    <div class="mt-2 flex items-baseline justify-between gap-2">
      <p class="text-base font-bold tabular-nums tracking-tight text-ds-fg-heading">
        {{ formatDealAmount(item.amount, item.currency) }}
      </p>
      <span
        v-if="isHotDeal"
        class="inline-flex items-center gap-1 rounded-full bg-ds-brand-subtle px-1.5 py-0.5 text-[10px] font-semibold text-ds-fg-brand ring-1 ring-inset ring-ds-brand/25"
        :title="$t('dealsCardHotHint')"
      >
        <UIcon name="i-heroicons-fire" class="h-3 w-3" aria-hidden="true" />
        {{ $t('dealsCardHot') }}
      </span>
    </div>

    <div class="mt-3 space-y-1.5">
      <div class="flex items-center justify-between gap-2 text-xs">
        <span class="text-ds-fg-muted">{{ $t('dealsCardProbability', { value: item.probability }) }}</span>
        <span
          v-if="item.expected_close_date"
          class="inline-flex items-center gap-1 tabular-nums"
          :class="closeDateClass"
        >
          <span
            v-if="isOverdue"
            class="relative flex h-1.5 w-1.5"
            aria-hidden="true"
          >
            <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-ds-danger opacity-60" />
            <span class="relative inline-flex h-1.5 w-1.5 rounded-full bg-ds-danger" />
          </span>
          <span>{{ closeDateLabel }}</span>
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

const { t } = useI18n()
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

function startOfDay(d: Date): number {
  const x = new Date(d)
  x.setHours(0, 0, 0, 0)
  return x.getTime()
}

const daysToClose = computed<number | null>(() => {
  if (!props.item.expected_close_date) return null
  const close = new Date(`${props.item.expected_close_date}T00:00:00`)
  if (Number.isNaN(close.getTime())) return null
  const todayStart = startOfDay(new Date())
  const closeStart = startOfDay(close)
  return Math.round((closeStart - todayStart) / 86_400_000)
})

const isOverdue = computed(() => {
  if (daysToClose.value === null) return false
  if (props.stage === 'won' || props.stage === 'lost') return false
  return daysToClose.value < 0
})

const closeDateClass = computed(() => {
  if (isOverdue.value) return 'font-semibold text-ds-danger'
  if (daysToClose.value !== null && daysToClose.value <= 7) return 'text-ds-warning'
  return 'text-ds-fg-subtle'
})

const closeDateLabel = computed(() => {
  if (!props.item.expected_close_date) return '—'
  if (props.stage === 'won' || props.stage === 'lost') return formatCloseDate(props.item.expected_close_date)
  const days = daysToClose.value
  if (days === null) return formatCloseDate(props.item.expected_close_date)
  if (days < 0) return t('dealsCardOverdueBy', { n: Math.abs(days) })
  if (days === 0) return t('dealsCardDueToday')
  if (days <= 7) return t('dealsCardInDays', { n: days })
  return formatCloseDate(props.item.expected_close_date)
})

const isHotDeal = computed(() => {
  if (props.stage === 'won' || props.stage === 'lost') return false
  return props.item.probability >= 70 && props.item.amount >= 100_000
})

const probabilityBarClass = computed(() => {
  if (props.item.probability >= 70) return 'bg-gradient-to-r from-ds-success to-emerald-500'
  if (props.item.probability >= 40) return 'bg-gradient-to-r from-ds-brand to-sky-500'
  if (props.item.probability >= 20) return 'bg-gradient-to-r from-ds-warning to-amber-500'
  return 'bg-ds-fg-subtle'
})

const accountLabel = computed(() => props.item.account_id ?? props.item.title)

const avatarInitials = computed(() => {
  const source = accountLabel.value || props.item.title || '?'
  const cleaned = source.replace(/^acc[-_]?/i, '').replace(/[-_]/g, ' ')
  const parts = cleaned.trim().split(/\s+/).filter(Boolean)
  if (parts.length === 0) return '?'
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase()
  return (parts[0][0] + parts[1][0]).toUpperCase()
})

const avatarPalette = [
  'bg-ds-info-subtle text-ds-info',
  'bg-ds-brand-subtle text-ds-fg-brand',
  'bg-ds-warning-subtle text-ds-warning',
  'bg-ds-success-subtle text-ds-success',
  'bg-ds-danger-subtle text-ds-danger',
] as const

const avatarClass = computed(() => {
  const key = accountLabel.value || props.item.id
  let hash = 0
  for (let i = 0; i < key.length; i++) {
    hash = (hash * 31 + (key.codePointAt(i) ?? 0)) >>> 0
  }
  return avatarPalette[hash % avatarPalette.length]
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
