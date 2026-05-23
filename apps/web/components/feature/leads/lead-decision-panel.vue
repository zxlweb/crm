<template>
  <div data-testid="lead-decision-panel">
    <LeadsLeadDetailMetrics :lead="lead" />

    <CardShell :title="$t('leadsEmotionTrendTitle')" class="mt-4">
      <CrmEmotionJourneyMap
        subject-type="lead"
        :subject-id="lead.id"
        embedded
        hide-touchpoints
        :chart-height="280"
        :range="emotionRange"
        :demo-badge-only-when-preview="demoBadgeOnlyWhenPreview"
      />
    </CardShell>
  </div>
</template>

<script setup lang="ts">
import type { EmotionJourneyQuery } from '~/composables/use-emotion-journey'
import type { Lead } from '~/types/lead'

defineProps<{
  lead: Lead
  demoBadgeOnlyWhenPreview?: boolean
}>()

const { t } = useI18n()

const emotionRange = ref<NonNullable<EmotionJourneyQuery['range']>>('90d')

const rangeItems = computed(() => [
  { id: '30d' as const, label: t('leadsEmotionRange30d') },
  { id: '90d' as const, label: t('leadsEmotionRange90d') },
  { id: 'all' as const, label: t('leadsEmotionRangeAll') },
])
</script>
