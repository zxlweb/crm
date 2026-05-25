<template>
  <PermissionGuard resource="contacts" action="view">
    <div v-if="pending" class="flex justify-center py-24">
      <UIcon name="i-heroicons-arrow-path" class="h-8 w-8 animate-spin text-primary" />
    </div>

    <p v-else-if="loadError" class="rounded-xl bg-ds-danger-subtle px-4 py-3 text-sm text-ds-danger">{{ loadError }}</p>

    <p v-else-if="!contact" class="text-sm text-ds-fg-muted">{{ $t('contactsNotFound') }}</p>

    <div v-else class="flex flex-col gap-4 xl:flex-row" data-testid="contact-detail-page">
      <div class="min-w-0 flex-1 space-y-4">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div class="space-y-2">
            <NuxtLink
              to="/contacts"
              class="inline-flex items-center gap-1 text-xs font-medium text-ds-fg-brand hover:underline"
            >
              {{ $t('contactsBackToList') }}
            </NuxtLink>
            <h2 class="text-2xl font-bold text-ds-fg-heading">{{ contact.display_name }}</h2>
            <div class="flex flex-wrap items-center gap-2">
              <CrmLifecycleBadge :stage="contact.lifecycle_stage" />
              <CrmRelationshipHealthBadge :health="contact.relationship_health" />
              <span v-if="contact.is_primary" class="text-xs text-ds-fg-brand">{{ $t('contactsPrimaryBadge') }}</span>
            </div>
          </div>
          <UiButton v-if="canUpdate" variant="secondary" data-testid="contact-edit-btn" @click="editOpen = true">
            {{ $t('edit') }}
          </UiButton>
        </div>

        <CardShell :title="$t('leadsTabOverview')" class="rounded-2xl">
          <dl class="grid gap-4 sm:grid-cols-2">
            <div>
              <dt class="text-xs text-ds-fg-muted">{{ $t('contactsColEmail') }}</dt>
              <dd class="mt-1 text-sm font-medium text-ds-fg-heading">{{ contact.email || '—' }}</dd>
            </div>
            <div>
              <dt class="text-xs text-ds-fg-muted">{{ $t('contactsFieldPhone') }}</dt>
              <dd class="mt-1 text-sm font-medium text-ds-fg-heading">{{ contact.phone || '—' }}</dd>
            </div>
            <div>
              <dt class="text-xs text-ds-fg-muted">{{ $t('contactsColAccount') }}</dt>
              <dd class="mt-1 text-sm font-medium text-ds-fg-heading">
                <NuxtLink
                  v-if="contact.account_id"
                  :to="`/accounts/${contact.account_id}`"
                  class="text-ds-fg-brand hover:underline"
                >
                  {{ accountLabel }}
                </NuxtLink>
                <span v-else>—</span>
              </dd>
            </div>
            <div>
              <dt class="text-xs text-ds-fg-muted">{{ $t('leadsFieldLastActivity') }}</dt>
              <dd class="mt-1 text-sm font-medium text-ds-fg-heading">{{ formatDate(contact.last_activity_at) }}</dd>
            </div>
          </dl>
        </CardShell>

        <CardShell :title="$t('leadsTabEmotion')" class="rounded-2xl">
          <div data-testid="tab-emotion-journey">
            <CrmEmotionJourneyMap ref="emotionMapRef" subject-type="contact" :subject-id="contact.id" />
          </div>
        </CardShell>

        <CardShell :title="$t('leadsTabTimeline')" :subtitle="$t('leadsSectionTimelineHint')" class="rounded-2xl">
          <template v-if="canCreateActivity" #header-extra>
            <div class="mt-3 flex justify-end">
              <UiButton
                size="sm"
                variant="secondary"
                icon="i-heroicons-plus-20-solid"
                data-testid="activity-create-btn"
                @click="activityOpen = true"
              >
                {{ $t('activityCreateBtn') }}
              </UiButton>
            </div>
          </template>
          <div class="space-y-6">
            <div>
              <h3 class="mb-2 text-sm font-semibold text-ds-fg-heading">{{ $t('activitySummaryTitle') }}</h3>
              <CrmActivitySummaryChart ref="summaryRef" subject-type="contact" :subject-id="contact.id" />
            </div>
            <div>
              <h3 class="mb-3 text-sm font-semibold text-ds-fg-heading">{{ $t('activityTimelineTitle') }}</h3>
              <CrmActivityTimeline ref="timelineRef" subject-type="contact" :subject-id="contact.id" />
            </div>
          </div>
        </CardShell>
      </div>

      <AiRelationPanel
        class="xl:sticky xl:top-4 xl:self-start"
        :show-preview="aiPreview.isPreviewMode"
        :insights="panelInsights"
        :engagement-score="panelEngagementScore"
        :churn-score="aiPreview.churnRiskScore"
        :disclaimer="aiPreview.disclaimer"
        :generate-copilot="aiPreview.generateCopilot"
      />
    </div>

    <UiModal v-model:open="editOpen" :title="$t('contactsEditTitle')">
      <form class="space-y-4" @submit.prevent="submitEdit">
        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('contactsFieldFirstName') }}</label>
            <UiInput v-model="formFirstName" />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('contactsFieldLastName') }}</label>
            <UiInput v-model="formLastName" />
          </div>
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('contactsColEmail') }}</label>
          <UiInput v-model="formEmail" type="email" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('contactsFieldPhone') }}</label>
          <UiInput v-model="formPhone" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('contactsColAccount') }}</label>
          <UiSelect v-model="formAccountId" :items="accountFormItems" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('leadsColLifecycle') }}</label>
          <UiSelect v-model="formLifecycle" :items="lifecycleFormItems" />
        </div>
        <label class="flex items-center gap-2 text-sm">
          <input v-model="formIsPrimary" type="checkbox" class="rounded border-ds-border" />
          {{ $t('contactsFieldPrimary') }}
        </label>
      </form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="secondary" @click="editOpen = false">{{ $t('cancel') }}</UiButton>
          <UiButton :loading="saving" @click="submitEdit">{{ $t('save') }}</UiButton>
        </div>
      </template>
    </UiModal>

    <UiModal v-model:open="activityOpen" :title="$t('activityCreateTitle')">
      <CrmActivityForm
        subject-type="contact"
        :subject-id="contact?.id ?? ''"
        :loading="activitySaving"
        @submit="submitActivity"
        @cancel="activityOpen = false"
      />
    </UiModal>
  </PermissionGuard>
</template>

<script setup lang="ts">
import type { ActivityCreateInput } from '~/types/activity'
import type { Contact, LifecycleStage } from '~/types/contact'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const { t, locale } = useI18n()
const permission = usePermission()
const contactsApi = useContacts()
const accountsApi = useAccounts()
const activitiesApi = useActivities()

const lifecycleOptions: LifecycleStage[] = ['acquire', 'activate', 'grow', 'retain', 'revive']

const contact = ref<Contact | null>(null)
const pending = ref(true)
const loadError = ref('')
const editOpen = ref(false)
const activityOpen = ref(false)
const saving = ref(false)
const activitySaving = ref(false)
const accountLabel = ref('—')
const timelineRef = ref<{ reload: () => Promise<void> } | null>(null)
const summaryRef = ref<{ reload: () => Promise<void> } | null>(null)
const emotionMapRef = ref<{ reload: () => Promise<void> } | null>(null)

const formFirstName = ref('')
const formLastName = ref('')
const formEmail = ref('')
const formPhone = ref('')
const formAccountId = ref('')
const formLifecycle = ref<LifecycleStage>('acquire')
const formIsPrimary = ref(false)
const accountFormItems = ref<Array<{ label: string; value: string }>>([])

const canUpdate = computed(() => permission.can('contacts', 'update'))
const canCreateActivity = computed(() => permission.can('activities', 'create'))

const contactId = computed(() => contact.value?.id ?? null)
const aiPreview = useAiPreview(contactId)
const fallbackEngagement = computed(() => contact.value?.engagement_score ?? null)

const previewInsights = aiPreview.previewInsights

const {
  insights: panelInsights,
  engagementScore: panelEngagementScore,
  reload: reloadInsights,
} = useDetailInsights({
  subjectType: 'contact',
  subjectId: contactId,
  previewInsights,
  fallbackEngagement,
})

const lifecycleFormItems = computed(() =>
  lifecycleOptions.map((s) => ({ label: t(`lifecycle.${s}`), value: s })),
)

function formatDate(iso: string | null) {
  if (!iso) return '—'
  return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(iso))
}

function fillForm(row: Contact) {
  formFirstName.value = row.first_name
  formLastName.value = row.last_name
  formEmail.value = row.email
  formPhone.value = row.phone
  formAccountId.value = row.account_id ?? ''
  formLifecycle.value = row.lifecycle_stage
  formIsPrimary.value = row.is_primary
}

async function resolveAccountLabel(accountId: string | null) {
  if (!accountId) {
    accountLabel.value = '—'
    return
  }
  const acc = await accountsApi.fetchById(accountId)
  accountLabel.value = acc?.name ?? accountId
}

async function loadAccountOptions() {
  const { data } = await accountsApi.fetchList({ page_size: 100 })
  accountFormItems.value = [
    { label: t('contactsNoAccount'), value: '' },
    ...data.items.map((a) => ({ label: a.name, value: a.id })),
  ]
}

async function load() {
  const id = String(route.params.id)
  pending.value = true
  loadError.value = ''
  try {
    contact.value = await contactsApi.fetchById(id)
    if (contact.value) {
      fillForm(contact.value)
      await resolveAccountLabel(contact.value.account_id)
    }
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

async function submitEdit() {
  if (!contact.value) return
  saving.value = true
  try {
    contact.value = await contactsApi.update(contact.value.id, {
      first_name: formFirstName.value.trim(),
      last_name: formLastName.value.trim(),
      email: formEmail.value.trim(),
      phone: formPhone.value.trim(),
      account_id: formAccountId.value || null,
      lifecycle_stage: formLifecycle.value,
      is_primary: formIsPrimary.value,
    })
    await resolveAccountLabel(contact.value.account_id)
    editOpen.value = false
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    saving.value = false
  }
}

async function submitActivity(payload: ActivityCreateInput) {
  if (!contact.value) return
  activitySaving.value = true
  try {
    await activitiesApi.create(payload)
    activityOpen.value = false
    await Promise.all([
      timelineRef.value?.reload(),
      summaryRef.value?.reload(),
      reloadInsights(),
      emotionMapRef.value?.reload(),
    ])
    contact.value = await contactsApi.fetchById(contact.value.id)
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    activitySaving.value = false
  }
}

onMounted(async () => {
  await loadAccountOptions()
  await load()
})

watch(() => route.params.id, load)
</script>
