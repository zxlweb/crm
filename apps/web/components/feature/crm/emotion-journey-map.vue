<template>
  <div class="space-y-6" :class="embedded ? 'space-y-4' : 'space-y-6'" data-testid="emotion-journey-map">
    <div v-if="pending" class="flex justify-center py-16">
      <UIcon name="i-heroicons-arrow-path" class="h-8 w-8 animate-spin text-primary" />
    </div>

    <UAlert v-else-if="loadError" color="red" variant="soft" :title="loadError" />

    <template v-else-if="journey">
      <div
        class="gap-3"
        :class="
          embedded
            ? 'flex flex-col sm:flex-row sm:items-center sm:gap-4'
            : 'flex flex-col gap-4'
        "
      >
        <div class="flex min-w-0 flex-wrap items-center gap-2 sm:shrink-0">
          <span
            v-if="journey.summary.current_sentiment"
            class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium"
            :class="sentimentTone(journey.summary.current_sentiment)"
          >
            <UiSentimentEmoji
              :sentiment="journey.summary.current_sentiment"
              size="xs"
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
     
        </div>

        <div
          v-if="journey.lifecycle_bands.length"
          class="min-w-0 space-y-1.5"
          :class="embedded ? 'flex-1 sm:border-l sm:border-ds-border sm:pl-4' : 'w-full'"
        >
          <div class="flex items-center gap-2">
            <p
              class="shrink-0 text-xs font-medium text-ds-fg-muted"
              :class="embedded ? '' : 'uppercase tracking-wide'"
            >
              {{ $t('emotionLifecycleBands') }}
            </p>
            <div
              class="flex h-3 min-w-0 flex-1 overflow-hidden rounded-full bg-ds-bg-muted"
              role="img"
              :aria-label="$t('emotionLifecycleBands')"
            >
              <div
                v-for="(band, idx) in proportionalBands"
                :key="`${band.stage}-${idx}`"
                class="h-full shrink-0 transition-all duration-300"
                :class="lifecycleBandClass(band.stage)"
                :style="{ width: `${band.percent}%` }"
                :title="bandTooltip(band)"
              />
            </div>
            <span
              v-if="embedded && proportionalBands.length === 1"
              class="inline-flex shrink-0 items-center gap-1 text-xs text-ds-fg-subtle"
            >
              <span class="h-2 w-2 rounded-full" :class="lifecycleBandClass(proportionalBands[0].stage)" />
              {{ $t(`lifecycle.${proportionalBands[0].stage}`) }}
              <span class="text-ds-fg-muted">({{ proportionalBands[0].percent }}%)</span>
            </span>
          </div>
          <div
            v-if="!embedded || proportionalBands.length > 1"
            class="flex flex-wrap gap-x-3 gap-y-0.5 text-xs text-ds-fg-subtle"
          >
            <span
              v-for="band in proportionalBands"
              :key="`legend-${band.stage}-${band.from}`"
              class="inline-flex items-center gap-1.5"
            >
              <span class="h-2 w-2 rounded-full" :class="lifecycleBandClass(band.stage)" />
              {{ $t(`lifecycle.${band.stage}`) }}
              <span class="text-ds-fg-muted">({{ band.percent }}%)</span>
            </span>
          </div>
        </div>
      </div>

      <div v-if="journey.points.length === 0" class="rounded-2xl border border-dashed border-ds-border bg-ds-bg-muted px-6 py-12 text-center">
        <p class="text-sm font-medium text-ds-fg-heading">{{ $t('emotionJourneyEmpty') }}</p>
        <p class="mt-2 text-sm text-ds-fg-muted">{{ $t('emotionJourneyEmptyHint') }}</p>
      </div>

      <template v-else>
        <div :class="embedded ? '' : ''">
          <p v-if="!embedded" class="mb-3 text-sm font-semibold text-ds-fg-heading">
            {{ $t('emotionJourneyChartTitle') }}
          </p>
          <ChartLine
            :categories="chartCategories"
            :series="chartSeries"
            :height="chartHeight"
            :y-formatter="yAxisLabelFormatter"
            :y-min="-2"
            :y-max="2"
            :y-interval="1"
            :point-symbols="chartPointSymbols"
            :point-item-styles="chartPointItemStyles"
            :point-symbol-size="26"
            :point-labels="chartPointLabels"
            :show-area="true"
          />
        </div>

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
import { resolveSentimentChartColor, sentimentEchartSymbol } from '@crm/ui-kit'
import type { LifecycleStage } from '~/types/lead'
import type { ActivitySentiment, EmotionJourney, EmotionSubjectType } from '~/types/emotion-journey'
const props = withDefaults(
  defineProps<{
    subjectType: EmotionSubjectType
    subjectId: string
    /** 与时间线合并展示时隐藏触点列表，避免重复 */
    hideTouchpoints?: boolean
    chartHeight?: number
    /** 仅 Preview / 演示租户显示「演示数据」标签 */
    demoBadgeOnlyWhenPreview?: boolean
    /** 嵌入决策面板：不包 Chart CardShell */
    embedded?: boolean
    /** 时间范围，变更时重新拉取 */
    range?: '30d' | '90d' | 'all'
    /** 父级递增以触发重新拉取（如新建 Activity 后） */
    refreshKey?: number
  }>(),
  { hideTouchpoints: false, chartHeight: 280, demoBadgeOnlyWhenPreview: true, embedded: false },
)

const { t, locale } = useI18n()
const emotionApi = useEmotionJourney()

const journey = ref<EmotionJourney | null>(null)
const pending = ref(true)
const loadError = ref('')
const usingDemoFixture = ref(false)

const showDemoBadge = computed(() => {
  if (!props.demoBadgeOnlyWhenPreview) {
    return usingDemoFixture.value
  }
  return usingDemoFixture.value && emotionApi.isPreview.value
})

type BandWithPercent = EmotionJourney['lifecycle_bands'][number] & { percent: number }

const proportionalBands = computed<BandWithPercent[]>(() => {
  if (!journey.value?.lifecycle_bands.length) return []
  const bands = journey.value.lifecycle_bands
  const spans = bands.map((b) => {
    const ms = Math.max(new Date(b.to).getTime() - new Date(b.from).getTime(), 86400000)
    return ms
  })
  const total = spans.reduce((a, b) => a + b, 0) || 1
  const raw = bands.map((band, i) => ({
    ...band,
    percent: Math.round((spans[i] / total) * 100),
  }))
  const sum = raw.reduce((a, b) => a + b.percent, 0) || 1
  return raw.map((band) => ({
    ...band,
    percent: Math.max(6, Math.round((band.percent / sum) * 100)),
  }))
})

function bandTooltip(band: BandWithPercent) {
  return `${t(`lifecycle.${band.stage}`)} · ${formatAt(band.from, true)} – ${formatAt(band.to, true)}`
}

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

const chartPointSymbols = computed(() =>
  sortedPoints.value.map((p) => sentimentEchartSymbol(resolvePointSentiment(p))),
)

const chartPointItemStyles = computed(() =>
  sortedPoints.value.map((p) => ({
    color: resolveSentimentChartColor(resolvePointSentiment(p)),
  })),
)

const chartPointLabels = computed(() =>
  sortedPoints.value.map((p) => t(`sentiment.${resolvePointSentiment(p)}`)),
)

function resolvePointSentiment(point: { sentiment?: ActivitySentiment | null; sentiment_score?: number | null }) {
  if (point.sentiment) return point.sentiment
  const byScore: Record<number, ActivitySentiment> = {
    2: 'positive',
    0: 'neutral',
    [-1]: 'hesitant',
    [-2]: 'negative',
  }
  if (point.sentiment_score != null && byScore[point.sentiment_score]) {
    return byScore[point.sentiment_score]
  }
  return 'unknown'
}

/** Y 轴文字刻度；数据点上方为 SVG 情绪图标（pointSymbols） */
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
      { range: props.range ?? '90d' },
    )
    journey.value = data
    usingDemoFixture.value = fromFixture
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

watch(() => [props.subjectType, props.subjectId, props.range], load, { immediate: true })
watch(
  () => props.refreshKey,
  (key, prev) => {
    if (key != null && key > 0 && key !== prev) void load()
  },
)
defineExpose({ reload: load })
</script>
