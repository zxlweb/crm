<template>
  <aside
    class="flex min-h-[20rem] flex-col overflow-hidden rounded-2xl border border-ds-border/60 bg-ds-bg-elevated shadow-ds-sm xl:min-h-[22rem]"
    data-testid="dashboard-calendar"
  >
    <div class="flex items-center justify-between border-b border-ds-border-muted px-4 py-3 sm:px-5">
      <div class="min-w-0">
        <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('dashboardCalendarTitle') }}</h3>
        <p class="mt-0.5 text-xs text-ds-fg-muted">{{ todayLabel }}</p>
      </div>
      <div class="flex items-center gap-1">
        <button
          type="button"
          class="inline-flex cursor-pointer items-center gap-1 rounded-lg border border-ds-border px-2 py-1 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand/30 hover:bg-ds-bg-muted/50 hover:text-ds-fg-brand"
          @click="shiftWeek(-1)"
        >
          <UIcon name="i-heroicons-chevron-left-20-solid" class="h-3.5 w-3.5" aria-hidden="true" />
        </button>
        <button
          type="button"
          class="inline-flex cursor-pointer items-center gap-1 rounded-lg border border-ds-border px-2 py-1 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand/30 hover:bg-ds-bg-muted/50 hover:text-ds-fg-brand"
          @click="shiftWeek(1)"
        >
          <UIcon name="i-heroicons-chevron-right-20-solid" class="h-3.5 w-3.5" aria-hidden="true" />
        </button>
      </div>
    </div>

    <div class="border-b border-ds-border-muted px-3 py-3 sm:px-4">
      <div class="grid grid-cols-7 gap-1">
        <span
          v-for="weekday in weekdayLabels"
          :key="weekday"
          class="text-center text-[10px] font-semibold uppercase tracking-wide text-ds-fg-subtle"
        >
          {{ weekday }}
        </span>
        <button
          v-for="day in weekDays"
          :key="day.iso"
          type="button"
          class="flex cursor-pointer flex-col items-center rounded-lg px-1 py-1.5 text-xs transition-colors duration-200"
          :class="dayButtonClass(day)"
          @click="selectedDate = day.iso"
        >
          <span class="font-semibold tabular-nums">{{ day.date }}</span>
        </button>
      </div>
    </div>

    <div class="flex min-h-0 flex-1 flex-col overflow-hidden">
      <div class="border-b border-ds-border-muted px-4 py-2 sm:px-5">
        <p class="text-xs font-medium text-ds-fg-muted">
          {{ selectedDateLabel }}
        </p>
      </div>

      <ul class="min-h-0 flex-1 space-y-2 overflow-y-auto p-3 sm:p-4">
        <li v-if="visibleEvents.length === 0" class="rounded-xl border border-dashed border-ds-border px-3 py-6 text-center">
          <UIcon name="i-heroicons-calendar-days-20-solid" class="mx-auto h-5 w-5 text-ds-fg-subtle" aria-hidden="true" />
          <p class="mt-2 text-xs text-ds-fg-muted">{{ $t('dashboardCalendarEmpty') }}</p>
        </li>
        <li v-for="event in visibleEvents" :key="event.id">
          <NuxtLink
            v-if="event.href"
            :to="event.href"
            class="group flex cursor-pointer gap-3 rounded-xl border border-ds-border-muted bg-ds-bg-muted/20 px-3 py-2.5 transition-colors duration-200 hover:border-ds-brand/25 hover:bg-ds-brand-subtle/30"
          >
            <span class="w-10 shrink-0 pt-0.5 text-xs font-semibold tabular-nums text-ds-fg-brand">{{ event.time }}</span>
            <span class="min-w-0 flex-1">
              <span class="block truncate text-xs font-semibold text-ds-fg-heading group-hover:text-ds-fg-brand">{{ event.title }}</span>
              <span class="mt-0.5 block text-[11px] text-ds-fg-muted">{{ event.subtitle }}</span>
            </span>
          </NuxtLink>
          <div
            v-else
            class="flex gap-3 rounded-xl border border-dashed border-ds-border-muted px-3 py-2.5"
          >
            <span class="w-10 shrink-0 pt-0.5 text-xs font-semibold tabular-nums text-ds-fg-muted">{{ event.time }}</span>
            <span class="min-w-0 flex-1">
              <span class="block text-xs font-medium text-ds-fg-muted">{{ event.title }}</span>
              <span class="mt-0.5 block text-[11px] text-ds-fg-subtle">{{ event.subtitle }}</span>
            </span>
          </div>
        </li>
      </ul>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { DASHBOARD_CALENDAR_FIXTURE } from '~/fixtures/dashboard-preview'
import type { DashboardCalendarEvent, PriorityActionItem } from '~/types/dashboard'

const props = defineProps<{
  priorities: PriorityActionItem[]
  isPreviewMode?: boolean
}>()

const { t, locale } = useI18n()

const weekOffset = ref(0)
const selectedDate = ref(toIsoDate(new Date()))

const weekdayLabels = computed(() => {
  const base = locale.value.startsWith('zh') ? ['一', '二', '三', '四', '五', '六', '日'] : ['Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa', 'Su']
  return base
})

const todayLabel = computed(() => {
  const formatter = new Intl.DateTimeFormat(locale.value, {
    month: 'long',
    day: 'numeric',
    weekday: 'short',
  })
  return formatter.format(new Date())
})

const weekDays = computed(() => {
  const anchor = new Date()
  anchor.setHours(12, 0, 0, 0)
  const day = anchor.getDay()
  const mondayOffset = day === 0 ? -6 : 1 - day
  anchor.setDate(anchor.getDate() + mondayOffset + weekOffset.value * 7)

  return Array.from({ length: 7 }, (_, index) => {
    const current = new Date(anchor)
    current.setDate(anchor.getDate() + index)
    return {
      iso: toIsoDate(current),
      date: current.getDate(),
      isToday: toIsoDate(current) === toIsoDate(new Date()),
      isSelected: toIsoDate(current) === selectedDate.value,
    }
  })
})

const selectedDateLabel = computed(() => {
  const formatter = new Intl.DateTimeFormat(locale.value, {
    month: 'long',
    day: 'numeric',
    weekday: 'long',
  })
  return formatter.format(parseIsoDate(selectedDate.value))
})

const events = computed<DashboardCalendarEvent[]>(() => {
  const todayIso = toIsoDate(new Date())
  const fromPriorities = props.priorities.slice(0, 4).map((item, index) => ({
    id: `priority-${item.id}`,
    date: todayIso,
    time: ['09:30', '11:00', '14:30', '16:00'][index] ?? '15:00',
    title: item.title,
    subtitle: t('dashboardCalendarFollowUp', {
      type: item.entityType === 'lead' ? t('dashboardTypeLead') : t('dashboardTypeAccount'),
    }),
    href: item.followHref,
  }))

  if (props.isPreviewMode) {
    return [...fromPriorities, ...DASHBOARD_CALENDAR_FIXTURE]
  }

  if (fromPriorities.length > 0) return fromPriorities

  return [
    {
      id: 'available-1',
      date: todayIso,
      time: '10:00',
      title: t('dashboardCalendarAvailable'),
      subtitle: t('dashboardCalendarAvailableHint'),
    },
  ]
})

const visibleEvents = computed(() =>
  events.value.filter((event) => event.date === selectedDate.value),
)

function shiftWeek(delta: number) {
  weekOffset.value += delta
}

function dayButtonClass(day: { isToday: boolean; isSelected: boolean }) {
  if (day.isSelected) return 'bg-ds-brand text-ds-on-brand shadow-ds-sm'
  if (day.isToday) return 'bg-ds-brand-subtle text-ds-fg-brand'
  return 'text-ds-fg-muted hover:bg-ds-bg-muted/60'
}

function toIsoDate(date: Date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function parseIsoDate(iso: string) {
  const [year, month, day] = iso.split('-').map(Number)
  return new Date(year, month - 1, day, 12)
}
</script>
