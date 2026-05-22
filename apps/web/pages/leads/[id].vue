<template>
  <PermissionGuard resource="leads" action="view">
    <div v-if="pending" class="flex justify-center py-24">
      <div class="h-8 w-8 animate-spin rounded-full border-2 border-ds-brand-muted border-t-ds-brand" />
    </div>

    <p v-else-if="loadError" class="rounded-xl bg-ds-danger-subtle px-4 py-3 text-sm text-ds-danger">{{ loadError }}</p>

    <p v-else-if="!lead" class="text-sm text-ds-fg-muted">{{ $t('leadsNotFound') }}</p>

    <div v-else class="flex flex-col gap-6 xl:flex-row" data-testid="lead-detail-page">
      <div class="min-w-0 flex-1 space-y-6">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div class="space-y-2">
            <NuxtLink
              to="/leads"
              class="inline-flex items-center gap-1 text-xs font-medium text-ds-fg-brand transition-colors hover:text-ds-brand-hover"
            >
              <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              {{ $t('leadsBackToList') }}
            </NuxtLink>
            <h2 class="text-2xl font-bold text-ds-fg-heading">{{ lead.title }}</h2>
            <div class="flex flex-wrap items-center gap-2">
              <CrmLeadStatusBadge :status="lead.status" />
              <CrmLifecycleBadge :stage="lead.lifecycle_stage" />
              <CrmRelationshipHealthBadge :health="lead.relationship_health" />
              <span class="text-xs text-ds-fg-muted">
                {{ $t('leadsEngagementLabel', { score: lead.engagement_score }) }}
              </span>
            </div>
          </div>
        </div>

        <LeadsLeadRelationshipHub :lead="lead" />
      </div>

      <AiRelationPanel
        v-if="showAiPanel"
        :show-preview="isPreviewMode"
        :insights="stubInsights"
        :engagement-score="lead.engagement_score"
      />
    </div>
  </PermissionGuard>
</template>

<script setup lang="ts">
import { DEMO_TENANT_ID } from '~/constants/demo'
import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'
import type { Lead } from '~/types/lead'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const { t } = useI18n()
const leadsApi = useLeads()

const lead = ref<Lead | null>(null)
const pending = ref(true)
const loadError = ref('')
const tenantCookie = useCookie('crm.tenant_id')

const isPreviewMode = computed(
  () =>
    route.query.preview === '1' ||
    lead.value?.id === DEMO_LEAD_ID ||
    tenantCookie.value === DEMO_TENANT_ID,
)

/** 2.3 骨架常显侧栏；2.13 接 tenant.config.ai_enabled */
const showAiPanel = true

const stubInsights = computed(() => {
  if (!isPreviewMode.value) return []
  return [
    {
      id: 'INS-001',
      title: t('aiStubInsight1Title'),
      body: t('aiStubInsight1Body'),
    },
    {
      id: 'INS-002',
      title: t('aiStubInsight2Title'),
      body: t('aiStubInsight2Body'),
    },
  ]
})

async function load() {
  const id = String(route.params.id)
  pending.value = true
  loadError.value = ''
  try {
    lead.value = await leadsApi.fetchById(id)
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

watch(() => route.params.id, load, { immediate: true })
</script>
