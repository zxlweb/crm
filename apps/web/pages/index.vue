<template>
  <DashboardSkeleton v-if="pending" />

  <UAlert v-else-if="loadError" color="red" variant="soft" :title="loadError" />

  <DashboardOverview
    v-else-if="snapshot"
    :snapshot="snapshot"
    :data-scope="snapshot.dataScope"
    :greeting="greeting"
    :headline="headline"
    :priority-summary="prioritySummary"
    :weekly-follow-up-note="weeklyFollowUpNote"
    :weekly-follow-up-count="snapshot.weeklyFollowUpCount"
    :insight-items="insightItems"
    :insight-detail-href="insightDetailHref"
    :can-create-lead="canCreateLead"
    :can-create-account="canCreateAccount"
    :read-only="readOnly"
    :show-zone-e="showZoneE"
    :show-team-heatmap="showTeamHeatmap"
    :show-team-ranking="showTeamRanking"
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

/** 数据展望 · 演示样例：仅租户管理员（全租户视角），销售经理/代表不展示 */
const showZoneE = computed(() => snapshot.value?.dataScope === 'all')
const showTeamHeatmap = computed(() => isPreviewMode.value || permission.can('rbac', 'view'))
const showTeamRanking = computed(() => snapshot.value?.canViewTeamRanking === true)
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
  const department = tenant.currentDepartment.value?.trim()
  if (name && workspace && department) {
    return t('dashboardHeadlineNamedDept', { name, workspace, department })
  }
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
    const msg = e instanceof Error ? e.message : ''
    const forbidden = msg.includes('403') || msg.includes('无权限') || msg.includes('Forbidden')
    if (forbidden && !auth.isSuperAdmin.value) {
      try {
        const list = await tenant.fetchTenants()
        const tid = tenant.currentTenantId.value
        const valid = !!tid && list.some((t) => t.id === tid)
        if (list.length > 0) {
          await tenant.switchTenant(valid ? tid! : list[0].id)
          snapshot.value = await dashboard.loadSnapshot(isPreviewMode.value)
          insightItems.value = dashboard.loadInsightItems(isPreviewMode.value)
          return
        }
      } catch {
        // fall through
      }
    }
    loadError.value = msg || t('loadFailed')
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
