<template>
  <div data-testid="lead-decision-panel">
    <LeadsLeadDetailMetrics :lead="lead" />

    <CardShell :title="$t('leadsEmotionTrendTitle')" class="mt-4">
      <template #header-extra>
        <div class="mt-2 flex justify-end sm:mt-0">
          <UiTabs
            v-model="emotionRange"
            :items="rangeTabItems"
            class="max-w-[280px]"
            data-testid="emotion-journey-range"
          />
        </div>
      </template>
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
    </CardShell>
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
  { id: '30d', label: t('leadsEmotionRange30d') },
  { id: '90d', label: t('leadsEmotionRange90d') },
  { id: 'all', label: t('leadsEmotionRangeAll') },
])

async function reloadEmotionJourney() {
  journeyRefreshKey.value += 1
  await nextTick()
  await mapRef.value?.reload?.()
}

defineExpose({ reloadEmotionJourney })
</script>
