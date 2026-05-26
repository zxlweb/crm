<template>
  <div
    class="dashboard-page relative"
    data-testid="dashboard-page"
    :data-dashboard-layout="dashboardLayoutTag"
  >
    <div
      class="pointer-events-none absolute inset-x-0 top-0 h-72 overflow-hidden"
      aria-hidden="true"
    >
      <div
        class="absolute -left-16 top-0 h-60 w-60 rounded-full blur-3xl opacity-70"
        :style="{ background: 'var(--ds-blur-brand)' }"
      />
      <div
        class="absolute right-0 top-6 h-44 w-44 rounded-full blur-3xl opacity-60"
        :style="{ background: 'var(--ds-blur-accent)' }"
      />
    </div>

    <div class="relative mx-auto max-w-[1400px] space-y-6">
      <!-- Hero Command Bar -->
      <DashboardHeroCommand
        :snapshot="snapshot"
        :greeting="greeting"
        :headline="headline"
        :weekly-follow-up-note="weeklyFollowUpNote"
        :is-preview-mode="isPreviewMode"
      />

      <!-- Focus Stream — merged attention queue for managers -->
      <DashboardFocusStream
        v-if="showFocusStream"
        :priorities="snapshot.priorities"
        :insights="insightItems"
        :insight-detail-href="insightDetailHref"
        :max-priorities="focusMaxPriorities"
        :max-cards="focusMaxCards"
      />

      <!-- KPI Row · at-a-glance numbers -->
      <section data-testid="dashboard-zone-metrics" aria-labelledby="dashboard-metrics-heading">
        <h2 id="dashboard-metrics-heading" class="sr-only">{{ headline }}</h2>
        <DashboardKpiRow
          variant="hero"
          :mode="kpiMode"
          :leads-total="snapshot.leadsTotal"
          :accounts-total="snapshot.accountsTotal"
          :deals-open-count="snapshot.dealsOpenCount"
          :deals-open-amount="snapshot.dealsOpenAmount"
          :avg-engagement="snapshot.avgEngagement"
          :at-risk-total="snapshot.atRiskTotal"
          :weekly-follow-up-count="weeklyFollowUpCount"
          :kpi-trends="snapshot.kpiTrends"
          :sparklines="snapshot.sparklines"
        />
      </section>

      <!-- Analytics — trends/charts -->
      <DashboardAnalyticsPanel
        :show-team-ranking="showTeamRanking"
        :layout-mode="analyticsLayoutMode"
        data-testid="dashboard-zone-analytics"
      />

      <!-- Priority + Calendar — today's lists (reps; managers use focus stream above) -->
      <div
        v-if="showPrioritySection"
        id="dashboard-priority-section"
        class="grid scroll-mt-24 gap-4 xl:grid-cols-[minmax(0,1fr)_min(340px,100%)] xl:items-start"
        data-testid="dashboard-zone-priority-calendar"
      >
        <DashboardPrioritySection
          :items="snapshot.priorities"
          :greeting="greeting"
          :headline="headline"
          :priority-summary="prioritySummary"
          :weekly-follow-up-note="weeklyFollowUpNote"
          at-risk-href="/leads?health=low"
          :can-create-lead="canCreateLead && !readOnly"
          :can-create-account="canCreateAccount && !readOnly"
          :read-only="readOnly"
        />
        <DashboardCalendar
          :priorities="snapshot.priorities"
          :is-preview-mode="isPreviewMode"
        />
      </div>

      <!-- Pipeline activity — reps only (managers/admins use analytics above) -->
      <DashboardPipelineTabs
        v-if="showOperationalFeed"
        :leads="snapshot.pipelineLeads"
        :accounts="snapshot.pipelineAccounts"
      />

      <!-- Zone E · Preview 展望 -->
      <section v-if="showZoneE" data-testid="dashboard-zone-e">
        <details ref="zoneEDetailsRef" class="group">
          <summary
            class="flex cursor-pointer list-none items-center justify-between gap-3 border-b border-ds-border-muted pb-3 marker:content-none"
          >
            <div class="min-w-0">
              <p class="text-xs font-semibold uppercase tracking-[0.12em] text-ds-fg-subtle">
                {{ $t('dashboardZoneELabel') }}
              </p>
              <h2 class="mt-1 text-lg font-semibold tracking-tight text-ds-fg-heading">
                {{ $t('dashboardZoneETitle') }}
              </h2>
            </div>
            <span class="inline-flex shrink-0 items-center gap-1.5 text-xs font-medium text-ds-fg-muted">
              <span class="hidden sm:inline">{{ zoneEOpen ? $t('dashboardZoneECollapse') : $t('dashboardZoneEExpand') }}</span>
              <UIcon
                name="i-heroicons-chevron-down"
                class="h-4 w-4 transition-transform duration-200 group-open:rotate-180"
                aria-hidden="true"
              />
            </span>
          </summary>

          <div class="mt-4 space-y-4">
            <div
              v-if="isPreviewMode"
              class="flex flex-wrap items-center gap-x-3 gap-y-2 rounded-xl border border-ds-border-muted bg-ds-bg-muted/30 px-4 py-2.5"
            >
              <AiPreviewBadge />
              <p class="text-xs leading-relaxed text-ds-fg-subtle">
                {{ $t('dashboardPreviewDisclaimer') }}
              </p>
            </div>

            <div
              class="-mx-1 flex gap-4 overflow-x-auto px-1 pb-1 snap-x snap-mandatory scroll-smooth lg:mx-0 lg:grid lg:gap-5 lg:overflow-visible lg:px-0 lg:pb-0"
              :class="showTeamHeatmap ? 'lg:grid-cols-2' : 'lg:max-w-2xl'"
            >
              <div class="w-[min(calc(100vw-3rem),22rem)] shrink-0 snap-center lg:w-auto lg:min-w-0 lg:shrink">
                <DashboardPreviewFunnel />
              </div>
              <div
                v-if="showTeamHeatmap"
                class="w-[min(calc(100vw-3rem),22rem)] shrink-0 snap-center lg:w-auto lg:min-w-0 lg:shrink"
              >
                <DashboardTeamHeatmap />
              </div>
            </div>
          </div>
        </details>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DashboardDataScope } from '~/types/dashboard-stats'
import type { DashboardInsightItem, DashboardSnapshot } from '~/types/dashboard'

const props = withDefaults(
  defineProps<{
    snapshot: DashboardSnapshot
    dataScope?: DashboardDataScope
    greeting: string
    headline: string
    prioritySummary: string
    weeklyFollowUpNote?: string
    weeklyFollowUpCount?: number
    insightItems: DashboardInsightItem[]
    insightDetailHref?: string
    canCreateLead: boolean
    canCreateAccount: boolean
    readOnly?: boolean
    showZoneE?: boolean
    showTeamHeatmap?: boolean
    showTeamRanking?: boolean
    isPreviewMode?: boolean
    zoneEDefaultOpen?: boolean
  }>(),
  {
    showZoneE: true,
    showTeamHeatmap: false,
    showTeamRanking: false,
    isPreviewMode: false,
    zoneEDefaultOpen: false,
  },
)

const showZoneE = computed(() => props.showZoneE)
const showTeamHeatmap = computed(() => props.showTeamHeatmap)
const showTeamRanking = computed(() => props.showTeamRanking)

const resolvedScope = computed(
  () => props.dataScope ?? props.snapshot.dataScope ?? 'self',
)

const isManagerDashboard = computed(
  () => resolvedScope.value === 'department' && props.showTeamRanking,
)

const isTenantAdminDashboard = computed(() => resolvedScope.value === 'all')

const showFocusStream = computed(() => !isTenantAdminDashboard.value)

const focusMaxPriorities = computed(() => (isManagerDashboard.value ? 5 : 2))
const focusMaxCards = computed(() => (isManagerDashboard.value ? 8 : 6))

const showOperationalFeed = computed(
  () => !isManagerDashboard.value && !isTenantAdminDashboard.value,
)

/** 今日重点 + 今日日程 — 仅销售代表 */
const showPrioritySection = computed(() => showOperationalFeed.value)

const kpiMode = computed(() => (isManagerDashboard.value ? 'manager' : 'full'))

const analyticsLayoutMode = computed(() =>
  isManagerDashboard.value ? 'manager' : 'default',
)

const weeklyFollowUpCount = computed(
  () => props.weeklyFollowUpCount ?? props.snapshot.weeklyFollowUpCount ?? 0,
)

const dashboardLayoutTag = computed(() => {
  if (isManagerDashboard.value) return 'manager'
  if (isTenantAdminDashboard.value) return 'admin'
  return 'rep'
})

const zoneEDetailsRef = ref<HTMLDetailsElement | null>(null)
const zoneEOpen = ref(props.zoneEDefaultOpen)

function syncZoneEOpen() {
  zoneEOpen.value = zoneEDetailsRef.value?.open ?? zoneEOpen.value
}

onMounted(() => {
  if (zoneEDetailsRef.value) {
    zoneEDetailsRef.value.open = props.zoneEDefaultOpen
    zoneEOpen.value = props.zoneEDefaultOpen
    zoneEDetailsRef.value.addEventListener('toggle', syncZoneEOpen)
  }
})

onBeforeUnmount(() => {
  zoneEDetailsRef.value?.removeEventListener('toggle', syncZoneEOpen)
})

watch(
  () => props.zoneEDefaultOpen,
  (open) => {
    if (zoneEDetailsRef.value) {
      zoneEDetailsRef.value.open = open
    }
    zoneEOpen.value = open
  },
)
</script>
