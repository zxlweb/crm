<template>
  <aside
    data-testid="ai-relation-panel"
    class="ds-panel-ai flex w-full shrink-0 flex-col gap-4 rounded-2xl p-4 lg:w-80 xl:w-96"
  >
    <div class="flex items-center justify-between gap-2">
      <div class="flex min-w-0 items-center gap-2.5">
        <AiAssistantAvatar size="md" />
        <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('aiRelationTitle') }}</h3>
      </div>
      <AiPreviewBadge v-if="showPreview && hasInsights" />
    </div>

    <p
      v-if="!hasInsights"
      class="rounded-xl border border-dashed border-ds-border bg-ds-bg px-4 py-6 text-center text-sm text-ds-fg-muted"
    >
      {{ $t('aiInsightsCompactEmpty') }}
    </p>

    <template v-else>
      <section class="space-y-3 rounded-xl border border-ds-border/80 bg-ds-bg p-4">
        <p class="text-xs font-medium uppercase tracking-wide text-ds-fg-muted">{{ $t('aiInsightsSection') }}</p>
        <ul class="space-y-2">
          <li
            v-for="item in insights"
            :key="item.id"
            class="rounded-lg border border-transparent px-3 py-2 text-sm text-ds-fg transition-colors duration-200 hover:border-ds-border hover:bg-ds-bg-muted"
          >
            <p class="font-medium text-ds-fg-heading">{{ item.title }}</p>
            <p class="mt-1 text-xs text-ds-fg-muted">{{ item.body }}</p>
          </li>
        </ul>
        <p v-if="engagementScore != null" class="text-xs text-ds-fg-muted">
          {{ $t('aiEngagementScore', { score: engagementScore }) }}
        </p>
      </section>

      <section class="space-y-3 rounded-xl border border-ds-border/80 bg-ds-bg p-4">
        <p class="text-xs font-medium uppercase tracking-wide text-ds-fg-muted">{{ $t('aiCopilotSection') }}</p>
        <div class="flex flex-wrap gap-2">
          <button
            type="button"
            class="ds-btn-secondary cursor-pointer rounded-lg px-3 py-1.5 text-xs font-medium transition-colors duration-200 active:scale-95"
            @click="onCopilotAction('followup')"
          >
            {{ $t('aiCopilotFollowup') }}
          </button>
          <button
            type="button"
            class="ds-btn-secondary cursor-pointer rounded-lg px-3 py-1.5 text-xs font-medium transition-colors duration-200 active:scale-95"
            @click="onCopilotAction('email')"
          >
            {{ $t('aiCopilotEmail') }}
          </button>
        </div>
        <p v-if="actionHint" class="text-xs text-ds-fg-muted" role="status">{{ actionHint }}</p>
      </section>

      <button
        type="button"
        data-testid="insight-adopt-btn"
        class="ds-btn-secondary w-full cursor-pointer rounded-xl py-2.5 text-sm font-medium transition-colors duration-200 active:scale-[0.98]"
        @click="onAdopt"
      >
        {{ $t('aiAdoptSuggestion') }}
      </button>
    </template>
  </aside>
</template>

<script setup lang="ts">
export type AiInsightStub = {
  id: string
  title: string
  body: string
}

const props = withDefaults(
  defineProps<{
    showPreview?: boolean
    insights?: AiInsightStub[]
    engagementScore?: number | null
  }>(),
  {
    showPreview: false,
    insights: () => [],
    engagementScore: null,
  },
)

const { t } = useI18n()

const hasInsights = computed(() => props.insights.length > 0)
const actionHint = ref('')
let hintTimer: ReturnType<typeof setTimeout> | undefined

function flashHint(message: string) {
  actionHint.value = message
  if (hintTimer) clearTimeout(hintTimer)
  hintTimer = setTimeout(() => {
    actionHint.value = ''
  }, 2800)
}

function onCopilotAction(kind: 'followup' | 'email') {
  const key = kind === 'followup' ? 'aiCopilotFollowup' : 'aiCopilotEmail'
  flashHint(t('aiCopilotComingSoon', { action: t(key) }))
}

function onAdopt() {
  flashHint(t('aiAdoptComingSoon'))
}

onUnmounted(() => {
  if (hintTimer) clearTimeout(hintTimer)
})
</script>
