<template>
  <div
    class="group rounded-xl border border-ds-border bg-ds-bg-elevated p-4 transition-all duration-200 hover:border-ds-brand-muted hover:shadow-ds-sm motion-reduce:transition-none"
    data-testid="lead-emotion-preview"
  >
    <div class="mb-3 flex flex-wrap items-center justify-between gap-2">
      <div>
        <h4 class="text-sm font-semibold text-ds-fg-heading">{{ $t('leadsEmotionPreviewTitle') }}</h4>
        <p class="text-xs text-ds-fg-muted">{{ $t('leadsEmotionPreviewHint') }}</p>
      </div>
      <span
        v-if="currentSentiment"
        class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium"
        :class="sentimentTone(currentSentiment)"
      >
        <UiSentimentEmoji :sentiment="currentSentiment" size="sm" />
        {{ $t(`sentiment.${currentSentiment}`) }}
      </span>
    </div>

    <div v-if="pending" class="h-[88px] animate-pulse rounded-lg bg-ds-bg-muted motion-reduce:animate-none" />

    <p v-else-if="loadError" class="text-xs text-ds-danger">{{ loadError }}</p>

    <p v-else-if="!hasPoints" class="py-6 text-center text-sm text-ds-fg-muted">
      {{ $t('emotionJourneyEmpty') }}
    </p>

    <ChartLine
      v-else
      :categories="chartCategories"
      :series="chartSeries"
      :height="88"
      :y-min="-2"
      :y-max="2"
      :y-interval="2"
      :show-area="true"
      :loading-text="$t('loading')"
    />

    <button
      type="button"
      class="mt-3 flex w-full cursor-pointer items-center justify-center gap-1 rounded-lg border border-ds-border bg-ds-bg-muted px-3 py-2 text-xs font-medium text-ds-fg-brand transition-colors duration-200 hover:bg-ds-brand-subtle hover:text-ds-brand-hover"
      @click="$emit('expand')"
    >
      {{ $t('leadsEmotionPreviewExpand') }}
      <svg class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import type { ActivitySentiment } from '~/types/emotion-journey'

const props = defineProps<{
  leadId: string
}>()

defineEmits<{
  expand: []
}>()

const { t, locale } = useI18n()
const emotionApi = useEmotionJourney()

const journey = ref<Awaited<ReturnType<typeof emotionApi.fetchJourney>>['journey'] | null>(null)
const pending = ref(true)
const loadError = ref('')

const sortedPoints = computed(() => {
  if (!journey.value) return []
  return [...journey.value.points].sort(
    (a, b) => new Date(a.at).getTime() - new Date(b.at).getTime(),
  )
})

const hasPoints = computed(() => sortedPoints.value.length > 0)

const currentSentiment = computed(() => journey.value?.summary.current_sentiment ?? null)

const chartCategories = computed(() =>
  sortedPoints.value.map((p) => formatAt(p.at, true)),
)

const chartSeries = computed(() => [
  {
    name: t('emotionChartSeries'),
    data: sortedPoints.value.map((p) => p.sentiment_score ?? 0),
    primary: true,
  },
])

function sentimentTone(sentiment: ActivitySentiment) {
  const map: Record<ActivitySentiment, string> = {
    positive: 'bg-ds-success-subtle text-ds-success',
    neutral: 'bg-ds-bg-muted text-ds-fg-muted',
    hesitant: 'bg-ds-warning-subtle text-ds-warning',
    negative: 'bg-ds-danger-subtle text-ds-danger',
    unknown: 'bg-ds-bg-muted text-ds-fg-subtle',
  }
  return map[sentiment]
}

function formatAt(iso: string, short = false) {
  const d = new Date(iso)
  if (short) {
    return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
      month: 'short',
      day: 'numeric',
    }).format(d)
  }
  return d.toLocaleString()
}

async function load() {
  pending.value = true
  loadError.value = ''
  try {
    const { journey: data } = await emotionApi.fetchJourney('lead', props.leadId, { range: '90d' })
    journey.value = data
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

watch(() => props.leadId, load, { immediate: true })
</script>
