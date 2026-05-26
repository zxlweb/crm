<template>
  <section
    class="flex flex-col rounded-2xl border border-ds-border/60 bg-ds-bg-elevated shadow-ds-sm"
    data-testid="dashboard-priority-section"
    data-zone="command-center"
  >
    <div class="flex items-start justify-between gap-3 border-b border-ds-border-muted px-4 py-3 sm:px-5">
      <div class="flex min-w-0 items-start gap-2">
        <span
          class="mt-0.5 inline-flex h-6 w-6 shrink-0 items-center justify-center rounded-lg bg-ds-brand-subtle text-ds-fg-brand ring-1 ring-inset ring-ds-brand/20"
          aria-hidden="true"
        >
          <UIcon name="i-heroicons-inbox-stack-20-solid" class="h-3.5 w-3.5" />
        </span>
        <div class="min-w-0">
          <h2 class="text-sm font-semibold text-ds-fg-heading">{{ $t('dashboardPriorityTitle') }}</h2>
          <p class="mt-0.5 text-xs text-ds-fg-muted">{{ prioritySummary }}</p>
        </div>
      </div>
      <span
        v-if="items.length > 0"
        class="inline-flex shrink-0 items-center rounded-full bg-ds-bg-muted px-2 py-0.5 text-[11px] font-semibold tabular-nums text-ds-fg ring-1 ring-inset ring-ds-border-muted"
      >
        {{ items.length }}
      </span>
    </div>

    <div class="p-3 sm:p-4">
      <div
        v-if="items.length === 0"
        class="rounded-xl border border-dashed border-ds-border bg-ds-bg-muted/20 px-5 py-6 text-center"
      >
        <div
          class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-ds-brand-subtle text-ds-fg-brand"
          aria-hidden="true"
        >
          <UIcon name="i-heroicons-check-circle-20-solid" class="h-5 w-5" />
        </div>
        <p class="text-sm font-semibold text-ds-fg-heading">{{ $t('dashboardPriorityEmptyTitle') }}</p>
        <p class="mx-auto mt-1 max-w-sm text-sm text-ds-fg-muted">{{ $t('dashboardPriorityEmptyDesc') }}</p>
        <div v-if="!readOnly" class="mt-4 flex flex-wrap justify-center gap-2">
          <UiButton
            v-if="canCreateLead"
            variant="primary"
            size="sm"
            icon="i-heroicons-plus-20-solid"
            to="/leads?create=1"
          >
            {{ $t('leadsCreate') }}
          </UiButton>
          <UiButton v-if="canCreateAccount" variant="secondary" size="sm" to="/accounts?create=1">
            {{ $t('accountsCreate') }}
          </UiButton>
        </div>
      </div>

      <div
        v-else
        class="grid gap-2 sm:gap-2.5"
        :class="[
          items.length > 2 ? 'xl:grid-cols-2' : 'grid-cols-1',
          items.length > 4 ? 'max-h-[min(24rem,60vh)] overflow-y-auto pr-0.5' : '',
        ]"
      >
        <DashboardPriorityActionCard
          v-for="item in items"
          :key="item.id"
          :item="item"
        />
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { PriorityActionItem } from '~/types/dashboard'

defineProps<{
  items: PriorityActionItem[]
  greeting: string
  headline: string
  prioritySummary: string
  weeklyFollowUpNote?: string
  atRiskHref: string
  canCreateLead: boolean
  canCreateAccount: boolean
  readOnly?: boolean
}>()
</script>
