<template>
  <PermissionGuard resource="leads" action="view">
    <LeadsLeadDetailSkeleton v-if="pending" />

    <p v-else-if="loadError" class="rounded-xl bg-ds-danger-subtle px-4 py-3 text-sm text-ds-danger" role="alert">
      {{ loadError }}
    </p>

    <p v-else-if="!lead" class="text-sm text-ds-fg-muted">{{ $t('leadsNotFound') }}</p>

    <div
      v-else
      class="lead-detail-page relative"
      data-testid="lead-detail-page"
    >
      <!-- Ambient backdrop glow -->
      <div
        class="pointer-events-none absolute inset-x-0 top-0 h-72 overflow-hidden"
        aria-hidden="true"
      >
        <div
          class="absolute -left-16 top-0 h-60 w-60 rounded-full blur-3xl opacity-70"
          :style="{ background: 'var(--ds-blur-brand)' }"
        />
        <div
          class="absolute right-0 top-8 h-44 w-44 rounded-full blur-3xl opacity-60"
          :style="{ background: 'var(--ds-blur-accent)' }"
        />
      </div>

      <div class="relative flex flex-col gap-5 xl:flex-row">
        <div class="min-w-0 flex-1 space-y-5">
          <!-- 1. Hero command bar -->
          <LeadsLeadHeroCommand
            :lead="lead"
            :owner-profile="ownerProfile"
            :can-update="canUpdate"
            :show-preview="isPreviewMode"
            :status-saving="statusSaving"
            :convert-saving="convertSaving"
            @edit="openEdit"
            @convert="openConvert"
            @status-change="onStatusChange"
          />

          <!-- 2. Pulse bar (composite health score) -->
          <LeadsLeadPulseBar :lead="lead" />

          <!-- 3. Focus signals (what to act on) -->
          <LeadsLeadFocusSignals
            :lead="lead"
            :churn-score="aiPreview.churnRiskScore"
          />

          <!-- 4. KPI metrics + emotion trend -->
          <LeadsLeadDecisionPanel
            ref="decisionPanelRef"
            :lead="lead"
            :emotion-refresh-key="emotionRefreshKey"
            :demo-badge-only-when-preview="true"
          />

          <!-- 5. Activity story (timeline) -->
          <section
            class="ds-lead-timeline relative overflow-hidden rounded-2xl border border-ds-border-muted bg-ds-bg-elevated shadow-ds-sm"
          >
            <span
              class="pointer-events-none absolute inset-x-0 top-0 h-0.5 bg-ds-brand opacity-70"
              aria-hidden="true"
            />
            <header
              class="flex flex-col gap-3 border-b border-ds-border-muted px-4 py-3 sm:flex-row sm:items-center sm:justify-between sm:px-5"
            >
              <div class="flex min-w-0 items-start gap-2.5">
                <span
                  class="mt-0.5 inline-flex h-6 w-6 shrink-0 items-center justify-center rounded-lg bg-ds-brand-subtle text-ds-fg-brand ring-1 ring-inset ring-ds-brand/20"
                  aria-hidden="true"
                >
                  <UIcon name="i-heroicons-bolt" class="h-3.5 w-3.5" />
                </span>
                <div class="min-w-0">
                  <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('leadsTabTimeline') }}</h3>
                  <p class="mt-0.5 text-xs text-ds-fg-muted">{{ $t('leadsSectionTimelineHint') }}</p>
                </div>
              </div>
              <UiButton
                v-if="canCreateActivity"
                size="sm"
                variant="primary"
                icon="i-heroicons-plus-20-solid"
                data-testid="activity-create-btn"
                @click="activityOpen = true"
              >
                {{ $t('activityCreateBtn') }}
              </UiButton>
            </header>

            <div id="timeline" class="space-y-5 p-4 sm:p-5" data-testid="lead-activity-timeline-section">
              <div>
                <h4 class="mb-1.5 text-xs font-semibold uppercase tracking-wider text-ds-fg-muted">
                  {{ $t('activitySummaryTitle') }}
                </h4>
                <p class="mb-2.5 text-xs text-ds-fg-subtle">{{ $t('activitySummaryHint') }}</p>
                <CrmActivitySummaryChart ref="summaryRef" subject-type="lead" :subject-id="lead.id" />
              </div>
              <div>
                <h4 class="mb-2.5 text-xs font-semibold uppercase tracking-wider text-ds-fg-muted">
                  {{ $t('activityTimelineTitle') }}
                </h4>
                <CrmActivityTimeline ref="timelineRef" subject-type="lead" :subject-id="lead.id" />
              </div>
            </div>
          </section>
        </div>

        <!-- AI Copilot rail -->
        <AiRelationPanel
          class="xl:sticky xl:top-4 xl:self-start"
          :show-preview="isPreviewMode"
          :insights="panelInsights"
          :engagement-score="panelEngagementScore"
          :churn-score="aiPreview.churnRiskScore"
          :disclaimer="aiPreview.disclaimer"
          :generate-copilot="aiPreview.generateCopilot"
        />
      </div>
    </div>

    <UiModal v-model:open="editOpen" :title="$t('leadsEditTitle')">
      <form class="space-y-4" @submit.prevent="submitEdit">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('leadsColTitle') }}</label>
          <UiInput v-model="formTitle" required data-testid="lead-form-title" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('leadsColSource') }}</label>
          <UiInput v-model="formSource" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('leadsFieldAmount') }}</label>
          <UiInput v-model="formAmount" type="number" min="0" step="1000" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('leadsFieldTags') }}</label>
          <UiInput v-model="formTagsText" :placeholder="$t('accountsTagsHint')" />
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="secondary" @click="editOpen = false">{{ $t('cancel') }}</UiButton>
          <UiButton :loading="saving" data-testid="lead-form-submit" @click="submitEdit">{{ $t('save') }}</UiButton>
        </div>
      </template>
    </UiModal>

    <UiModal v-model:open="activityOpen" :title="$t('activityCreateTitle')">
      <CrmActivityForm
        subject-type="lead"
        :subject-id="lead?.id ?? ''"
        :loading="activitySaving"
        @submit="submitActivity"
        @cancel="activityOpen = false"
      />
    </UiModal>

    <UiModal v-model:open="convertOpen" :title="$t('leadsConvertTitle')">
      <p class="mb-4 text-sm text-ds-fg-muted">{{ $t('leadsConvertDesc') }}</p>
      <form class="space-y-4" @submit.prevent="submitConvert">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('leadsConvertAccountName') }}</label>
          <UiInput v-model="convertAccountName" required data-testid="lead-convert-account-name" />
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="secondary" @click="convertOpen = false">{{ $t('cancel') }}</UiButton>
          <UiButton :loading="convertSaving" data-testid="lead-convert-submit" @click="submitConvert">
            {{ $t('leadsConvertConfirm') }}
          </UiButton>
        </div>
      </template>
    </UiModal>
  </PermissionGuard>
</template>

<script setup lang="ts">
import { canTransitionLeadStatus } from '~/utils/lead-status-transition'
import type { ActivityCreateInput } from '~/types/activity'
import type { Lead, LeadStatus } from '~/types/lead'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const { t } = useI18n()
const permission = usePermission()
const leadsApi = useLeads()
const activitiesApi = useActivities()
const ownerProfiles = useOwnerProfile()

const lead = ref<Lead | null>(null)
const pending = ref(true)
const loadError = ref('')
const editOpen = ref(false)
const convertOpen = ref(false)
const activityOpen = ref(false)
const activitySaving = ref(false)
const saving = ref(false)
const timelineRef = useTemplateRef<{ reload: () => Promise<void> }>('timelineRef')
const summaryRef = useTemplateRef<{ reload: () => Promise<void> }>('summaryRef')
const decisionPanelRef = useTemplateRef<{ reloadEmotionJourney: () => Promise<void> }>('decisionPanelRef')
const emotionRefreshKey = ref(0)
const statusSaving = ref(false)
const convertSaving = ref(false)

const formTitle = ref('')
const formSource = ref('')
const formAmount = ref('0')
const formTagsText = ref('')
const convertAccountName = ref('')

const leadId = computed(() => lead.value?.id ?? null)
const aiPreview = useAiPreview(leadId)

const canUpdate = computed(() => permission.can('leads', 'update'))
const canCreateActivity = computed(() => permission.can('activities', 'create'))

const ownerProfile = computed(() => ownerProfiles.resolve(lead.value?.owner_id))

const isPreviewMode = aiPreview.isPreviewMode
const previewInsights = aiPreview.previewInsights

const fallbackEngagement = computed(() => lead.value?.engagement_score ?? null)

const {
  insights: panelInsights,
  engagementScore: panelEngagementScore,
  reload: reloadInsights,
} = useDetailInsights({
  subjectType: 'lead',
  subjectId: leadId,
  previewInsights,
  fallbackEngagement,
})

function parseTags(text: string): string[] {
  return text
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
}

function fillForm(row: Lead) {
  formTitle.value = row.title
  formSource.value = row.source
  formAmount.value = String(row.amount)
  formTagsText.value = (row.tags ?? []).join(', ')
}

function openEdit() {
  if (lead.value) fillForm(lead.value)
  editOpen.value = true
}

function openConvert() {
  if (!lead.value) return
  convertAccountName.value = lead.value.title
  convertOpen.value = true
}

async function load() {
  const id = String(route.params.id)
  pending.value = true
  loadError.value = ''
  try {
    lead.value = await leadsApi.fetchById(id)
    if (lead.value) fillForm(lead.value)
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

async function onStatusChange(next: LeadStatus) {
  if (!lead.value || next === lead.value.status) return
  if (!canTransitionLeadStatus(lead.value.status, next)) {
    loadError.value = t('leadsStatusTransitionDenied')
    return
  }
  statusSaving.value = true
  loadError.value = ''
  try {
    lead.value = await leadsApi.update(lead.value.id, { status: next })
    fillForm(lead.value)
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    statusSaving.value = false
  }
}

async function submitEdit() {
  if (!lead.value || !formTitle.value.trim()) return
  saving.value = true
  loadError.value = ''
  try {
    lead.value = await leadsApi.update(lead.value.id, {
      title: formTitle.value.trim(),
      source: formSource.value.trim() || undefined,
      amount: Number(formAmount.value) || 0,
      tags: parseTags(formTagsText.value),
    })
    editOpen.value = false
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    saving.value = false
  }
}

async function reloadAfterActivityChange() {
  emotionRefreshKey.value += 1
  await nextTick()
  await Promise.all([
    timelineRef.value?.reload?.(),
    summaryRef.value?.reload?.(),
    reloadInsights(),
    decisionPanelRef.value?.reloadEmotionJourney?.(),
  ])
}

async function submitActivity(payload: ActivityCreateInput) {
  if (!lead.value) return
  activitySaving.value = true
  loadError.value = ''
  try {
    await activitiesApi.create(payload)
    activityOpen.value = false
    await nextTick()
    await reloadAfterActivityChange()
    lead.value = await leadsApi.fetchById(lead.value.id)
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    activitySaving.value = false
  }
}

async function submitConvert() {
  if (!lead.value || !convertAccountName.value.trim()) return
  convertSaving.value = true
  loadError.value = ''
  try {
    lead.value = await leadsApi.convert(lead.value.id, {
      create_account: { name: convertAccountName.value.trim() },
    })
    convertOpen.value = false
    fillForm(lead.value)
    emotionRefreshKey.value += 1
    await nextTick()
    await decisionPanelRef.value?.reloadEmotionJourney?.()
  } catch (e) {
    const msg = e instanceof Error ? e.message : t('loadFailed')
    loadError.value =
      msg === 'convert_not_allowed' ? t('leadsConvertNotAllowed') : msg
  } finally {
    convertSaving.value = false
  }
}

watch(() => route.params.id, load, { immediate: true })
</script>
