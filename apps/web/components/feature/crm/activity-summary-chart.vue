<template>
  <div class="space-y-2" data-testid="activity-summary-chart">
    <UAlert v-if="loadError" color="red" variant="soft" :title="loadError" />
    <LeadsReportChartSlot
      :pending="pending"
      :empty="!barItems.length && !loadError"
      :empty-text="$t('activitySummaryEmpty')"
      :height="height"
    >
      <ChartBar
        :key="`activity-summary-${chartLocale}`"
        :items="barItems"
        :height="height"
        :horizontal="false"
      />
    </LeadsReportChartSlot>
  </div>
</template>

<script setup lang="ts">
import type { ChartBarItem } from '@crm/ui-kit'
import type { ActivitySubjectType, ActivitySummary } from '~/types/activity'

const props = withDefaults(
  defineProps<{
    subjectType: ActivitySubjectType
    subjectId: string
    height?: number
  }>(),
  { height: 220 },
)

const { t, locale } = useI18n()
const { activityTypeLabel } = useActivityLabels()
const activitiesApi = useActivities()

const chartLocale = computed(() => locale.value)
const pending = ref(true)
const loadError = ref('')
const summary = ref<ActivitySummary | null>(null)

const barItems = computed<ChartBarItem[]>(() => {
  void locale.value
  return (summary.value?.items ?? []).map((item) => ({
    name: activityTypeLabel(item.label),
    value: item.value,
  }))
})

async function load() {
  if (!props.subjectId) return
  pending.value = true
  loadError.value = ''
  try {
    summary.value = await activitiesApi.fetchSummary({
      subjectType: props.subjectType,
      subjectId: props.subjectId,
    })
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
    summary.value = null
  } finally {
    pending.value = false
  }
}

watch(() => [props.subjectType, props.subjectId], load, { immediate: true })
defineExpose({ reload: load })
</script>
