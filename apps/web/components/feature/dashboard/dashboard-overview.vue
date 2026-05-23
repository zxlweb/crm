<template>
  <div class="dashboard-page relative" data-testid="dashboard-page">
    <div
      class="pointer-events-none absolute inset-x-0 top-0 h-64 overflow-hidden"
      aria-hidden="true"
    >
      <div
        class="absolute -left-16 top-0 h-48 w-48 rounded-full blur-3xl opacity-70"
        :style="{ background: 'var(--ds-blur-brand)' }"
      />
      <div
        class="absolute right-0 top-6 h-40 w-40 rounded-full blur-3xl opacity-60"
        :style="{ background: 'var(--ds-blur-accent)' }"
      />
    </div>

    <div class="relative mx-auto max-w-[1400px] space-y-6">
      <!-- Row 1 · KPI 全宽 -->
      <section data-testid="dashboard-zone-metrics" aria-labelledby="dashboard-metrics-heading">
        <div class="mb-4 flex flex-wrap items-end justify-between gap-3">
          <div class="min-w-0">
            <p class="text-xs font-medium text-ds-fg-muted">{{ greeting }}</p>
            <h1 id="dashboard-metrics-heading" class="mt-0.5 text-xl font-bold tracking-tight text-ds-fg-heading sm:text-2xl">
              {{ headline }}
            </h1>
          </div>
          <p v-if="weeklyFollowUpNote" class="text-xs text-ds-fg-subtle">{{ weeklyFollowUpNote }}</p>
        </div>

        <DashboardKpiRow
          variant="hero"
          :leads-total="snapshot.leadsTotal"
          :accounts-total="snapshot.accountsTotal"
          :avg-engagement="snapshot.avgEngagement"
          :at-risk-total="snapshot.atRiskTotal"
          :kpi-trends="snapshot.kpiTrends"
        />
      </section>

      <!-- Row 2 · 今日优先 + 今日日程 -->
      <div
        class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_min(340px,100%)] xl:items-stretch"
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

      <!-- Zone C + D · 管道 + 洞察 -->
      <div class="grid gap-6 xl:grid-cols-[minmax(0,1fr)_min(360px,100%)] xl:items-start">
        <DashboardPipelineTabs
          :leads="snapshot.pipelineLeads"
          :accounts="snapshot.pipelineAccounts"
        />

        <aside class="xl:sticky xl:top-6">
          <DashboardInsightCompact
            :items="insightItems"
            :detail-href="insightDetailHref"
          />
        </aside>
      </div>

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
import type { DashboardInsightItem, DashboardSnapshot } from '~/types/dashboard'

const props = withDefaults(
  defineProps<{
    snapshot: DashboardSnapshot
    greeting: string
    headline: string
    prioritySummary: string
    weeklyFollowUpNote?: string
    insightItems: DashboardInsightItem[]
    insightDetailHref?: string
    canCreateLead: boolean
    canCreateAccount: boolean
    readOnly?: boolean
    showZoneE?: boolean
    showTeamHeatmap?: boolean
    isPreviewMode?: boolean
    zoneEDefaultOpen?: boolean
  }>(),
  {
    showZoneE: true,
    showTeamHeatmap: false,
    isPreviewMode: false,
    zoneEDefaultOpen: false,
  },
)

const showZoneE = computed(() => props.showZoneE)
const showTeamHeatmap = computed(() => props.showTeamHeatmap)

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
