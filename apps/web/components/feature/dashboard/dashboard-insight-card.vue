<template>
  <li
    class="rounded-xl border p-3.5"
    :class="surfaceClass"
    data-testid="dashboard-insight-card"
  >
    <div class="flex items-start gap-2.5">
      <div
        class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg"
        :class="iconWrapClass"
        aria-hidden="true"
      >
        <UIcon :name="iconName" class="h-4 w-4" />
      </div>
      <div class="min-w-0 flex-1">
        <div class="flex flex-wrap items-center gap-1.5">
          <span class="rounded-md px-1.5 py-0.5 text-[10px] font-bold" :class="titlePillClass">
            {{ item.title }}
          </span>
          <span
            v-if="item.urgent"
            class="rounded-md bg-red-100 px-1.5 py-0.5 text-[10px] font-bold text-red-700 dark:bg-red-950/50 dark:text-red-300"
          >
            {{ $t('dashboardInsightTagUrgent') }}
          </span>
        </div>
        <p class="mt-2 text-xs leading-relaxed" :class="bodyClass">{{ item.body }}</p>
      </div>
    </div>
  </li>
</template>

<script setup lang="ts">
import type { DashboardInsightItem } from '~/types/dashboard'

const props = defineProps<{
  item: DashboardInsightItem
}>()

const variantStyles = computed(() => {
  switch (props.item.variant) {
    case 'churn':
      return {
        surface: 'border-orange-200/80 bg-orange-50/80 dark:border-orange-900/40 dark:bg-orange-950/25',
        iconWrap: 'bg-orange-100 text-orange-600 dark:bg-orange-900/50 dark:text-orange-400',
        icon: 'i-heroicons-exclamation-triangle-20-solid',
        titlePill: 'bg-orange-100 text-orange-800 dark:bg-orange-900/50 dark:text-orange-300',
        body: 'text-orange-900/85 dark:text-orange-200/90',
      }
    case 'opportunity':
      return {
        surface: 'border-emerald-200/80 bg-emerald-50/80 dark:border-emerald-900/40 dark:bg-emerald-950/25',
        iconWrap: 'bg-emerald-100 text-emerald-600 dark:bg-emerald-900/50 dark:text-emerald-400',
        icon: 'i-heroicons-arrow-trending-up-20-solid',
        titlePill: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/50 dark:text-emerald-300',
        body: 'text-emerald-900/85 dark:text-emerald-200/90',
      }
    default:
      return {
        surface: 'border-sky-200/80 bg-sky-50/80 dark:border-sky-900/40 dark:bg-sky-950/25',
        iconWrap: 'bg-sky-100 text-sky-600 dark:bg-sky-900/50 dark:text-sky-400',
        icon: 'i-heroicons-light-bulb-20-solid',
        titlePill: 'bg-sky-100 text-sky-800 dark:bg-sky-900/50 dark:text-sky-300',
        body: 'text-sky-900/85 dark:text-sky-200/90',
      }
  }
})

const surfaceClass = computed(() => variantStyles.value.surface)
const iconWrapClass = computed(() => variantStyles.value.iconWrap)
const iconName = computed(() => variantStyles.value.icon)
const titlePillClass = computed(() => variantStyles.value.titlePill)
const bodyClass = computed(() => variantStyles.value.body)
</script>
