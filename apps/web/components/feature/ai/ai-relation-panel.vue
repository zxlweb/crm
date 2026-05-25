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
      <AiPreviewBadge v-if="showPreview && (hasInsights || churnScore != null)" />
    </div>

    <p
      v-if="showPreview && disclaimer"
      class="rounded-lg border border-ds-border-muted bg-ds-bg-muted/40 px-3 py-2 text-xs leading-relaxed text-ds-fg-subtle"
      data-testid="ai-preview-disclaimer"
    >
      {{ disclaimer }}
    </p>

    <p
      v-if="!hasInsights && churnScore == null"
      class="rounded-xl border border-dashed border-ds-border bg-ds-bg px-4 py-6 text-center text-sm text-ds-fg-muted"
    >
      {{ $t('aiInsightsCompactEmpty') }}
    </p>

    <template v-else>
      <section
        v-if="showPreview && churnScore != null"
        class="space-y-2 rounded-xl border border-ds-border/80 bg-ds-bg p-4"
        data-testid="ai-churn-risk-bar"
      >
        <div class="flex items-center justify-between gap-2">
          <p class="text-xs font-medium uppercase tracking-wide text-ds-fg-muted">
            {{ $t('aiChurnRiskLabel') }}
          </p>
          <span class="text-sm font-semibold tabular-nums text-ds-warning">{{ churnScore }}%</span>
        </div>
        <div class="h-2 overflow-hidden rounded-full bg-ds-bg-muted" role="presentation">
          <div
            class="h-full rounded-full bg-gradient-to-r from-amber-400 to-orange-500 transition-all duration-500"
            :style="{ width: `${Math.min(100, Math.max(0, churnScore))}%` }"
          />
        </div>
        <p class="text-xs text-ds-fg-subtle">{{ $t('aiPreviewScoreHint') }}</p>
      </section>

      <section v-if="hasInsights" class="space-y-3 rounded-xl border border-ds-border/80 bg-ds-bg p-4">
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
            data-testid="ai-copilot-followup-btn"
            @click="onCopilotAction('followup')"
          >
            {{ $t('aiCopilotFollowup') }}
          </button>
          <button
            type="button"
            class="ds-btn-secondary cursor-pointer rounded-lg px-3 py-1.5 text-xs font-medium transition-colors duration-200 active:scale-95"
            data-testid="ai-copilot-email-btn"
            @click="onCopilotAction('email')"
          >
            {{ $t('aiCopilotEmail') }}
          </button>
        </div>
        <div
          v-if="copilotOutput"
          class="rounded-lg border border-ds-border bg-ds-bg-muted/50 px-3 py-2.5 text-xs leading-relaxed text-ds-fg"
          data-testid="ai-copilot-output"
          role="status"
        >
          {{ copilotOutput }}
        </div>
        <p v-if="actionHint" class="text-xs text-ds-fg-muted" role="status">{{ actionHint }}</p>
      </section>

      <button
        v-if="hasInsights || copilotOutput"
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
    churnScore?: number | null
    disclaimer?: string
    generateCopilot?: (scene: 'followup_script' | 'email_draft') => string
  }>(),
  {
    showPreview: false,
    insights: () => [],
    engagementScore: null,
    churnScore: null,
    disclaimer: '',
    generateCopilot: undefined,
  },
)

const { t } = useI18n()

const hasInsights = computed(() => props.insights.length > 0)
const copilotOutput = ref('')
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
  actionHint.value = ''
  if (props.showPreview && props.generateCopilot) {
    const scene = kind === 'followup' ? 'followup_script' : 'email_draft'
    copilotOutput.value = props.generateCopilot(scene)
    return
  }
  const key = kind === 'followup' ? 'aiCopilotFollowup' : 'aiCopilotEmail'
  flashHint(t('aiCopilotComingSoon', { action: t(key) }))
}

function onAdopt() {
  if (props.showPreview && copilotOutput.value) {
    flashHint(t('aiAdoptPreviewSuccess'))
    return
  }
  flashHint(t('aiAdoptComingSoon'))
}

onUnmounted(() => {
  if (hintTimer) clearTimeout(hintTimer)
})
</script>
