<template>
  <aside
    data-testid="ai-relation-panel"
    class="flex w-full shrink-0 flex-col gap-4 rounded-2xl border border-ds-border bg-ds-bg-muted p-4 lg:w-80 xl:w-96"
  >
    <div class="flex items-center justify-between gap-2">
      <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('aiRelationTitle') }}</h3>
      <AiPreviewBadge v-if="showPreview" />
    </div>

    <section class="space-y-3 rounded-xl border border-ds-border bg-ds-bg p-4">
      <p class="text-xs font-medium uppercase tracking-wide text-ds-fg-muted">{{ $t('aiInsightsSection') }}</p>
      <ul v-if="insights.length" class="space-y-2">
        <li
          v-for="item in insights"
          :key="item.id"
          class="rounded-lg bg-ds-bg-muted px-3 py-2 text-sm text-ds-fg"
        >
          <p class="font-medium text-ds-fg-heading">{{ item.title }}</p>
          <p class="mt-1 text-xs text-ds-fg-muted">{{ item.body }}</p>
        </li>
      </ul>
      <p v-else class="text-sm text-ds-fg-muted">{{ $t('aiInsightsEmpty') }}</p>
      <p v-if="engagementScore != null" class="text-xs text-ds-fg-muted">
        {{ $t('aiEngagementScore', { score: engagementScore }) }}
      </p>
    </section>

    <section class="space-y-3 rounded-xl border border-ds-border bg-ds-bg p-4">
      <p class="text-xs font-medium uppercase tracking-wide text-ds-fg-muted">{{ $t('aiCopilotSection') }}</p>
      <div class="flex flex-wrap gap-2">
        <button
          type="button"
          class="ds-btn-secondary rounded-lg px-3 py-1.5 text-xs font-medium transition-transform duration-150 active:scale-95"
          disabled
        >
          {{ $t('aiCopilotFollowup') }}
        </button>
        <button
          type="button"
          class="ds-btn-secondary rounded-lg px-3 py-1.5 text-xs font-medium transition-transform duration-150 active:scale-95"
          disabled
        >
          {{ $t('aiCopilotEmail') }}
        </button>
      </div>
      <p class="text-xs text-ds-fg-subtle">{{ $t('aiCopilotSkeletonHint') }}</p>
    </section>

    <button
      type="button"
      data-testid="insight-adopt-btn"
      class="ds-btn-primary w-full rounded-xl py-2.5 text-sm font-medium transition-transform duration-150 active:scale-[0.98] disabled:opacity-50"
      disabled
    >
      {{ $t('aiAdoptSuggestion') }}
    </button>
  </aside>
</template>

<script setup lang="ts">
export type AiInsightStub = {
  id: string
  title: string
  body: string
}

withDefaults(
  defineProps<{
    showPreview?: boolean
    insights?: AiInsightStub[]
    engagementScore?: number | null
  }>(),
  {
    showPreview: true,
    insights: () => [],
    engagementScore: null,
  },
)
</script>
