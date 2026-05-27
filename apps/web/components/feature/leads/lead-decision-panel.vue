<template>
  <div class="space-y-4" data-testid="lead-decision-panel">
    <LeadsLeadDetailMetrics :lead="lead" />

    <section
      class="ds-lead-emotion relative overflow-hidden rounded-2xl border border-ds-border-muted bg-ds-bg-elevated shadow-ds-sm"
    >
      <span
        class="pointer-events-none absolute inset-x-0 top-0 h-0.5 bg-ds-info opacity-80"
        aria-hidden="true"
      />
      <span
        class="pointer-events-none absolute -right-10 -top-10 h-32 w-32 rounded-full opacity-15 blur-3xl"
        :style="{ background: 'var(--ds-blur-brand)' }"
        aria-hidden="true"
      />
      <header
        class="relative flex flex-col gap-3 border-b border-ds-border-muted px-4 py-3 sm:flex-row sm:items-center sm:justify-between sm:px-5"
      >
        <div class="flex min-w-0 items-start gap-2.5">
          <span
            class="mt-0.5 inline-flex h-6 w-6 shrink-0 items-center justify-center rounded-lg bg-ds-info-subtle text-ds-info ring-1 ring-inset ring-ds-info/20"
            aria-hidden="true"
          >
            <UIcon name="i-heroicons-heart" class="h-3.5 w-3.5" />
          </span>
          <div class="min-w-0">
            <h3 class="text-sm font-semibold text-ds-fg-heading">
              {{ $t('leadsEmotionTrendTitle') }}
            </h3>
            <p class="mt-0.5 text-xs text-ds-fg-muted">
              {{ $t('leadsEmotionTrendHint') }}
            </p>
          </div>
        </div>
        <UiTabs
          v-model="emotionRange"
          :items="rangeTabItems"
          class="max-w-[280px] shrink-0"
          data-testid="emotion-journey-range"
        />
      </header>
      <div class="relative p-3 sm:p-4">
        <div
          class="mb-4 rounded-xl border border-ds-border-muted bg-ds-bg-muted/45 p-3 sm:p-4"
          data-testid="emotion-confidence-demo-v2"
        >
          <div class="grid gap-3 lg:grid-cols-[240px_minmax(0,1fr)]">
            <div class="rounded-lg border border-ds-border-muted bg-ds-bg-elevated p-3 shadow-ds-sm">
              <div class="flex items-center gap-3">
                <div
                  class="flex h-14 w-14 shrink-0 items-center justify-center rounded-full bg-ds-info-subtle text-ds-info ring-1 ring-inset ring-ds-info/25"
                  aria-label="情绪可信分 78%"
                >
                  <div class="text-center leading-none">
                    <p class="text-xl font-bold">78</p>
                    <p class="mt-0.5 text-[10px] font-medium">可信分</p>
                  </div>
                </div>
                <div class="min-w-0">
                  <div class="flex flex-wrap items-center gap-1.5">
                    <h4 class="text-sm font-semibold text-ds-fg-heading">情绪可信度</h4>
                    <span class="rounded-full bg-ds-success-subtle px-2 py-0.5 text-[11px] font-medium text-ds-success">
                      中高可信
                    </span>
                  </div>
                  <p class="mt-1 text-xs leading-5 text-ds-fg-muted">
                    四类信号加权校准，不只看人工标注。
                  </p>
                </div>
              </div>
              <div class="mt-3 h-1.5 overflow-hidden rounded-full bg-ds-bg-muted">
                <div class="h-full w-[78%] rounded-full bg-ds-info" />
              </div>
            </div>

            <div class="grid gap-2 sm:grid-cols-2 xl:grid-cols-4">
              <div
                v-for="item in confidenceWeights"
                :key="item.label"
                class="rounded-lg border border-ds-border-muted bg-ds-bg-elevated p-3"
              >
                <div class="flex items-center justify-between gap-2">
                  <div class="flex min-w-0 items-center gap-1.5">
                    <UIcon :name="item.icon" class="h-3.5 w-3.5 shrink-0 text-ds-info" />
                    <span class="truncate text-xs font-semibold text-ds-fg-heading">{{ item.label }}</span>
                  </div>
                  <span class="shrink-0 text-xs font-bold text-ds-info">{{ item.weight }}</span>
                </div>
                <p class="mt-1.5 min-h-[34px] text-xs leading-[17px] text-ds-fg-muted">{{ item.desc }}</p>
                <div class="mt-2 h-1 overflow-hidden rounded-full bg-ds-bg-muted">
                  <div class="h-full rounded-full bg-ds-info" :style="{ width: item.weight }" />
                </div>
              </div>
            </div>
          </div>

          <div class="mt-3 flex flex-wrap items-center gap-2 text-xs">
            <span class="font-medium text-ds-fg-muted">本次证据</span>
            <span
              v-for="signal in confidenceSignalsClean"
              :key="signal"
              class="inline-flex items-center rounded-full bg-ds-bg-elevated px-2.5 py-1 text-ds-fg-muted ring-1 ring-inset ring-ds-border-muted"
            >
              {{ signal }}
            </span>
          </div>
        </div>

        <CrmEmotionJourneyMap
          ref="mapRef"
          subject-type="lead"
          :subject-id="lead.id"
          embedded
          hide-touchpoints
          :chart-height="280"
          :range="emotionRange"
          :refresh-key="journeyRefreshKey"
          :demo-badge-only-when-preview="demoBadgeOnlyWhenPreview"
        />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { EmotionJourneyQuery } from '~/composables/use-emotion-journey'
import type { Lead } from '~/types/lead'

const props = defineProps<{
  lead: Lead
  demoBadgeOnlyWhenPreview?: boolean
  /** 父页递增以刷新情绪旅程（如新建 Activity） */
  emotionRefreshKey?: number
}>()

const { t } = useI18n()

const mapRef = useTemplateRef<{ reload: () => Promise<void> }>('mapRef')

const emotionRange = ref<NonNullable<EmotionJourneyQuery['range']>>('90d')
const journeyRefreshKey = ref(0)

const confidenceWeights = [
  {
    label: '人工标注',
    weight: '30%',
    desc: '销售选择的基础情绪。',
    icon: 'i-heroicons-pencil-square',
  },
  {
    label: '文本识别',
    weight: '30%',
    desc: '会议、邮件、微信内容识别。',
    icon: 'i-heroicons-chat-bubble-left-right',
  },
  {
    label: '行为数据',
    weight: '25%',
    desc: '回复、参会、沉默等行为校准。',
    icon: 'i-heroicons-chart-bar',
  },
  {
    label: '业务阶段',
    weight: '15%',
    desc: '阶段推进验证关系真实度。',
    icon: 'i-heroicons-arrow-trending-up',
  },
]

const confidenceSignalsClean = [
  '人工标注：犹豫',
  '会议纪要命中预算顾虑',
  '近 7 天有触达但未推进',
  '生命周期停留在激活阶段',
]

watch(
  () => props.emotionRefreshKey,
  (key, prev) => {
    if (key != null && key > 0 && key !== prev) {
      journeyRefreshKey.value = key
    }
  },
)

const rangeTabItems = computed(() => [
  { id: '30d', label: t('leadsEmotionRange30d'), icon: 'i-heroicons-calendar-days' },
  { id: '90d', label: t('leadsEmotionRange90d'), icon: 'i-heroicons-calendar' },
  { id: 'all', label: t('leadsEmotionRangeAll'), icon: 'i-heroicons-clock' },
])

async function reloadEmotionJourney() {
  journeyRefreshKey.value += 1
  await nextTick()
  await mapRef.value?.reload?.()
}

defineExpose({ reloadEmotionJourney })
</script>
