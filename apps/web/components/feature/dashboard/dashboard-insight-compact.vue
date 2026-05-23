<template>
  <aside
    class="ds-panel-ai flex h-full flex-col overflow-hidden rounded-2xl shadow-ds-sm"
    data-testid="dashboard-insight-compact"
  >
    <div
      class="relative overflow-hidden border-b border-ds-border-muted bg-gradient-to-br from-ds-bg-elevated via-ds-bg-elevated to-violet-500/5 px-5 py-4"
    >
      <div
        class="pointer-events-none absolute -right-6 -top-8 h-24 w-24 rounded-full blur-2xl"
        :style="{ background: 'var(--ds-blur-brand)' }"
        aria-hidden="true"
      />

      <div class="relative flex items-center gap-3">
        <AiAssistantAvatar size="lg" />
        <div class="min-w-0 flex-1">
          <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('dashboardInsightPanelTitle') }}</h3>
          <p class="mt-0.5 text-xs text-ds-fg-muted">{{ $t('dashboardAiAvatarTagline') }}</p>
        </div>
      </div>
    </div>

    <div class="flex flex-1 flex-col gap-4 p-4">
      <ul class="space-y-3">
        <DashboardInsightCard
          v-for="item in items"
          :key="item.id"
          :item="item"
        />
      </ul>

      <NuxtLink
        v-if="detailHref"
        :to="detailHref"
        class="mt-auto inline-flex cursor-pointer items-center justify-center gap-1.5 rounded-xl border border-ds-border bg-ds-bg px-4 py-2.5 text-sm font-semibold text-ds-fg-brand transition-colors hover:border-ds-brand/30 hover:bg-ds-brand-subtle/40"
      >
        {{ $t('dashboardInsightPanelCta') }}
        <UIcon name="i-heroicons-arrow-right-20-solid" class="h-4 w-4" />
      </NuxtLink>
    </div>
  </aside>
</template>

<script setup lang="ts">
import type { DashboardInsightItem } from '~/types/dashboard'

defineProps<{
  items: DashboardInsightItem[]
  detailHref?: string
}>()
</script>
