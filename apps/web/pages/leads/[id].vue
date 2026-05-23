<template>
  <PermissionGuard resource="leads" action="view">
    <LeadsLeadDetailSkeleton v-if="pending" />

    <p v-else-if="loadError" class="rounded-xl bg-ds-danger-subtle px-4 py-3 text-sm text-ds-danger" role="alert">
      {{ loadError }}
    </p>

    <p v-else-if="!lead" class="text-sm text-ds-fg-muted">{{ $t('leadsNotFound') }}</p>

    <div v-else class="flex flex-col gap-4 xl:flex-row" data-testid="lead-detail-page">
      <div class="min-w-0 flex-1 space-y-3">
        <LeadsLeadDetailHeader
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



        <LeadsLeadDecisionPanel :lead="lead" :demo-badge-only-when-preview="true" />

        <CardShell :title="$t('leadsTabTimeline')" :subtitle="$t('leadsSectionTimelineHint')" class="rounded-2xl">
          <div id="timeline" data-testid="lead-activity-timeline-section">
            <CrmActivityTimeline :lead-id="lead.id" />
          </div>
        </CardShell>
      </div>

      <AiRelationPanel
        class="xl:sticky xl:top-4 xl:self-start"
        :show-preview="isPreviewMode && stubInsights.length > 0"
        :insights="stubInsights"
        :engagement-score="lead.engagement_score"
      />
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
import { DEMO_TENANT_ID } from '~/constants/demo'
import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'
import { canTransitionLeadStatus } from '~/utils/lead-status-transition'
import type { Lead, LeadStatus } from '~/types/lead'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const { t } = useI18n()
const permission = usePermission()
const leadsApi = useLeads()
const ownerProfiles = useOwnerProfile()

const lead = ref<Lead | null>(null)
const pending = ref(true)
const loadError = ref('')
const editOpen = ref(false)
const convertOpen = ref(false)
const saving = ref(false)
const statusSaving = ref(false)
const convertSaving = ref(false)

const formTitle = ref('')
const formSource = ref('')
const formAmount = ref('0')
const formTagsText = ref('')
const convertAccountName = ref('')

const tenantCookie = useCookie('crm.tenant_id')

const canUpdate = computed(() => permission.can('leads', 'update'))

const ownerProfile = computed(() => ownerProfiles.resolve(lead.value?.owner_id))

const isPreviewMode = computed(
  () =>
    route.query.preview === '1' ||
    lead.value?.id === DEMO_LEAD_ID ||
    tenantCookie.value === DEMO_TENANT_ID,
)



const stubInsights = computed(() => {
  if (!isPreviewMode.value) return []
  return [
    { id: 'INS-001', title: t('aiStubInsight1Title'), body: t('aiStubInsight1Body') },
    { id: 'INS-002', title: t('aiStubInsight2Title'), body: t('aiStubInsight2Body') },
  ]
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
