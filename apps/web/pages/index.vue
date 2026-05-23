<template>
  <DashboardSkeleton v-if="pending" />

  <UAlert v-else-if="loadError" color="red" variant="soft" :title="loadError" />

  <DashboardOverview
    v-else-if="snapshot"
    :snapshot="snapshot"
    :greeting="greeting"
    :headline="headline"
    :priority-summary="prioritySummary"
    :weekly-follow-up-note="weeklyFollowUpNote"
    :insight-items="insightItems"
    :insight-detail-href="insightDetailHref"
    :can-create-lead="canCreateLead"
    :can-create-account="canCreateAccount"
    :read-only="readOnly"
    :show-zone-e="showZoneE"
    :show-team-heatmap="showTeamHeatmap"
    :is-preview-mode="isPreviewMode"
    :zone-e-default-open="zoneEDefaultOpen"
  />
</template>

<script setup lang="ts">
import { DEMO_TENANT_ID } from '~/constants/demo'
import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'
import type { DashboardInsightItem, DashboardSnapshot } from '~/types/dashboard'

definePageMeta({ layout: 'app', middleware: 'auth' })

const { t } = useI18n()
const auth = useAuth()
const tenant = useTenant()
const permission = usePermission()
const dashboard = useDashboard()
const leadsApi = useLeads()
const tenantCookie = useCookie<string | null>('crm.tenant_id')

const snapshot = ref<DashboardSnapshot | null>(null)
const insightItems = ref<DashboardInsightItem[]>([])
const pending = ref(true)
const loadError = ref('')

const canCreateLead = computed(() => permission.can('leads', 'create'))
const canCreateAccount = computed(() => permission.can('accounts', 'create'))
const canUpdateLead = computed(() => permission.can('leads', 'update'))

const readOnly = computed(
  () => !canCreateLead.value && !canCreateAccount.value && !canUpdateLead.value,
)

const isPreviewMode = computed(
  () => leadsApi.useMock.value || tenantCookie.value === DEMO_TENANT_ID,
)

const showZoneE = computed(() => true)
const showTeamHeatmap = computed(() => isPreviewMode.value || permission.can('rbac', 'view'))
const zoneEDefaultOpen = computed(
  () => isPreviewMode.value || permission.can('rbac', 'view'),
)

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return t('dashboardGreetingMorning')
  if (hour < 18) return t('dashboardGreetingAfternoon')
  return t('dashboardGreetingEvening')
})

const headline = computed(() => {
  const name = auth.user.value?.name?.trim() || auth.user.value?.email?.split('@')[0] || ''
  const workspace = tenant.currentTenant.value?.name
  if (name && workspace) {
    return t('dashboardHeadlineNamed', { name, workspace })
  }
  if (name) return t('dashboardHeadlineUser', { name })
  return t('dashboardHeadlineDefault')
})

const prioritySummary = computed(() => {
  const count = snapshot.value?.priorityCount ?? 0
  if (count === 0) return t('dashboardPrioritySummaryNone')
  return t('dashboardPrioritySummary', { count })
})

const weeklyFollowUpNote = computed(() => {
  const count = snapshot.value?.weeklyFollowUpCount ?? 0
  if (count <= 0) return ''
  return t('dashboardWeeklyFollowUps', { count })
})

const insightDetailHref = computed(() => {
  const firstLead = snapshot.value?.priorities.find((item) => item.entityType === 'lead')
  if (firstLead) return `/leads/${firstLead.entityId}`
  if (isPreviewMode.value) return `/leads/${DEMO_LEAD_ID}`
  return undefined
})

async function reload() {
  pending.value = true
  loadError.value = ''
  try {
    snapshot.value = await dashboard.loadSnapshot(isPreviewMode.value)
    insightItems.value = dashboard.loadInsightItems(isPreviewMode.value)
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

onMounted(() => {
  void reload()
})

watch(
  () => tenant.currentTenantId.value,
  () => {
    void reload()
  },
)
</script>
