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
