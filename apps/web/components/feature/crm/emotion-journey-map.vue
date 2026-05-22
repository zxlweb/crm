<template>
  <div class="space-y-6" data-testid="emotion-journey-map">
    <div v-if="pending" class="flex justify-center py-16">
      <UIcon name="i-heroicons-arrow-path" class="h-8 w-8 animate-spin text-primary" />
    </div>

    <UAlert v-else-if="loadError" color="red" variant="soft" :title="loadError" />

    <template v-else-if="journey">
      <div class="flex flex-wrap items-center gap-3">
        <span
          v-if="journey.summary.current_sentiment"
          class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium"
          :class="sentimentTone(journey.summary.current_sentiment)"
        >
          <UiSentimentEmoji
            :sentiment="journey.summary.current_sentiment"
            size="md"
            :label="$t(`sentiment.${journey.summary.current_sentiment}`)"
          />
          {{ $t('emotionCurrentSentiment') }}：{{ $t(`sentiment.${journey.summary.current_sentiment}`) }}
        </span>
        <span class="inline-flex items-center rounded-full bg-ds-bg-muted px-3 py-1 text-xs font-medium text-ds-fg-muted">
          {{ $t('emotionTrendLabel') }}：{{ $t(`emotionTrend.${journey.summary.trend}`) }}
        </span>
        <span
          v-if="journey.summary.days_since_positive != null"
          class="text-xs text-ds-fg-muted"
        >
          {{ $t('emotionDaysSincePositive', { days: journey.summary.days_since_positive }) }}
        </span>
        <UBadge v-if="showDemoBadge" color="amber" variant="subtle" size="sm">
          {{ $t('emotionDemoData') }}
        </UBadge>
      </div>

      <div v-if="journey.lifecycle_bands.length" class="space-y-2">
        <p class="text-xs font-medium uppercase tracking-wide text-ds-fg-muted">{{ $t('emotionLifecycleBands') }}</p>
        <div class="flex h-2 overflow-hidden rounded-full bg-ds-bg-muted">
          <div
            v-for="(band, idx) in journey.lifecycle_bands"
            :key="`${band.stage}-${idx}`"
            class="h-full min-w-[2rem] flex-1"
            :class="lifecycleBandClass(band.stage)"
            :title="$t(`lifecycle.${band.stage}`)"
          />
        </div>
        <div class="flex flex-wrap gap-2 text-xs text-ds-fg-subtle">
          <span v-for="band in journey.lifecycle_bands" :key="band.stage">
            {{ $t(`lifecycle.${band.stage}`) }}
          </span>
        </div>
      </div>

      <div v-if="journey.points.length === 0" class="rounded-2xl border border-dashed border-ds-border bg-ds-bg-muted px-6 py-12 text-center">
        <p class="text-sm font-medium text-ds-fg-heading">{{ $t('emotionJourneyEmpty') }}</p>
        <p class="mt-2 text-sm text-ds-fg-muted">{{ $t('emotionJourneyEmptyHint') }}</p>
      </div>

      <template v-else>
        <CardShell :title="$t('emotionJourneyChartTitle')" class="rounded-2xl">
          <ChartLine
            :categories="chartCategories"
            :series="chartSeries"
            :height="240"
            :y-formatter="yAxisLabelFormatter"
            :y-min="-2"
            :y-max="2"
            :y-interval="1"
            :point-emojis="chartPointEmojis"
            :point-labels="chartPointLabels"
          />
        </CardShell>

        <section v-if="!hideTouchpoints">
          <h4 class="mb-3 text-sm font-semibold text-ds-fg-heading">{{ $t('emotionTouchpoints') }}</h4>
          <ul class="space-y-2">
            <li
              v-for="point in sortedPoints"
              :key="point.activity_id"
              class="flex flex-col gap-1 rounded-xl border border-ds-border bg-ds-bg-elevated px-4 py-3 text-sm sm:flex-row sm:items-center sm:justify-between"
            >
              <div>
                <p class="font-medium text-ds-fg-heading">{{ point.label || point.event_type }}</p>
                <p class="text-xs text-ds-fg-muted">{{ formatAt(point.at) }}</p>
              </div>
              <span
                v-if="point.sentiment"
                class="inline-flex shrink-0 items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium"
                :class="sentimentTone(point.sentiment)"
              >
                <UiSentimentEmoji
                  :sentiment="point.sentiment"
                  size="sm"
                  :label="$t(`sentiment.${point.sentiment}`)"
                />
                {{ $t(`sentiment.${point.sentiment}`) }}
              </span>
            </li>
          </ul>
        </section>

        <section v-if="journey.milestones.length">
          <h4 class="mb-3 text-sm font-semibold text-ds-fg-heading">{{ $t('emotionMilestones') }}</h4>
          <ul class="space-y-2">
            <li
              v-for="m in journey.milestones"
              :key="`${m.type}-${m.at}`"
              class="rounded-xl border border-ds-border-muted bg-ds-bg-muted px-4 py-2 text-sm text-ds-fg-muted"
            >
              <span class="font-medium text-ds-fg-heading">{{ m.label }}</span>
              <span class="mx-2">·</span>
              <span>{{ formatAt(m.at) }}</span>
            </li>
          </ul>
        </section>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { sentimentEmoji } from '@crm/ui-kit'
import type { LifecycleStage } from '~/types/lead'
import type { ActivitySentiment, EmotionJourney, EmotionSubjectType } from '~/types/emotion-journey'
const props = withDefaults(
  defineProps<{
    subjectType: EmotionSubjectType
    subjectId: string
    /** 与时间线合并展示时隐藏触点列表，避免重复 */
    hideTouchpoints?: boolean
  }>(),
  { hideTouchpoints: false },
)

const { t, locale } = useI18n()
const emotionApi = useEmotionJourney()

const journey = ref<EmotionJourney | null>(null)
const pending = ref(true)
const loadError = ref('')
const usingDemoFixture = ref(false)

const showDemoBadge = computed(
  () => usingDemoFixture.value || emotionApi.isPreview.value,
)

const sortedPoints = computed(() => {
  if (!journey.value) return []
  return [...journey.value.points].sort(
    (a, b) => new Date(a.at).getTime() - new Date(b.at).getTime(),
  )
})

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

const chartPointEmojis = computed(() =>
  sortedPoints.value.map((p) =>
    p.sentiment ? sentimentEmoji(p.sentiment) : sentimentEmoji('unknown'),
  ),
)

const chartPointLabels = computed(() =>
  sortedPoints.value.map((p) => {
    const key = p.sentiment ?? 'unknown'
    const em = sentimentEmoji(key)
    const label = t(`sentiment.${key}`)
    return `${em} ${label}`
  }),
)

/** Y 轴仅文字刻度；emoji 由 pointEmojis 标在数据点上方 */
function yAxisLabelFormatter(value: number) {
  const map: Record<number, string> = {
    2: t('sentiment.positive'),
    0: t('sentiment.neutral'),
    [-1]: t('sentiment.hesitant'),
    [-2]: t('sentiment.negative'),
  }
  return map[value] ?? ''
}

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

function lifecycleBandClass(stage: string) {
  const map: Record<LifecycleStage, string> = {
    acquire: 'bg-ds-bg-muted',
    activate: 'bg-blue-200 dark:bg-blue-900/50',
    grow: 'bg-violet-300 dark:bg-violet-800/60',
    retain: 'bg-emerald-300 dark:bg-emerald-800/50',
    revive: 'bg-amber-300 dark:bg-amber-800/50',
  }
  return map[stage as LifecycleStage] ?? 'bg-ds-brand-subtle'
}

function formatAt(iso: string, short = false) {
  const d = new Date(iso)
  if (short) {
    return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
      month: 'short',
      day: 'numeric',
    }).format(d)
  }
  return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(d)
}

async function load() {
  pending.value = true
  loadError.value = ''
  usingDemoFixture.value = false
  try {
    const { journey: data, fromFixture } = await emotionApi.fetchJourney(
      props.subjectType,
      props.subjectId,
    )
    journey.value = data
    usingDemoFixture.value = fromFixture
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

watch(() => [props.subjectType, props.subjectId], load, { immediate: true })
</script>
