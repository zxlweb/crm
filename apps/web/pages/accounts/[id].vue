<template>
  <PermissionGuard resource="accounts" action="view">
    <div v-if="pending" class="flex justify-center py-24">
      <div class="h-8 w-8 animate-spin rounded-full border-2 border-ds-brand-muted border-t-ds-brand" />
    </div>

    <p v-else-if="loadError" class="rounded-xl bg-ds-danger-subtle px-4 py-3 text-sm text-ds-danger">{{ loadError }}</p>

    <p v-else-if="!account" class="text-sm text-ds-fg-muted">{{ $t('accountsNotFound') }}</p>

    <div v-else class="flex flex-col gap-6 xl:flex-row" data-testid="account-detail-page">
      <div class="min-w-0 flex-1 space-y-6">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div class="space-y-2">
            <NuxtLink
              to="/accounts"
              class="inline-flex items-center gap-1 text-xs font-medium text-ds-fg-brand transition-colors hover:text-ds-brand-hover"
            >
              <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              {{ $t('accountsBackToList') }}
            </NuxtLink>
            <h2 class="text-2xl font-bold text-ds-fg-heading">{{ account.name }}</h2>
            <div class="flex flex-wrap items-center gap-2">
              <CrmLifecycleBadge :stage="account.lifecycle_stage" />
              <CrmRelationshipHealthBadge :health="account.relationship_health" />
              <span class="text-xs text-ds-fg-muted">
                {{ $t('leadsEngagementLabel', { score: account.engagement_score }) }}
              </span>
            </div>
          </div>
          <UiButton
            v-if="canUpdate"
            variant="secondary"
            data-testid="account-edit-btn"
            @click="editOpen = true"
          >
            {{ $t('edit') }}
          </UiButton>
        </div>

        <AccountsAccountDetailTabs v-model="activeTab">
          <template #overview>
            <CardShell :title="$t('leadsTabOverview')" class="rounded-2xl">
              <dl class="grid gap-4 sm:grid-cols-2">
                <div>
                  <dt class="text-xs text-ds-fg-muted">{{ $t('accountsColIndustry') }}</dt>
                  <dd class="mt-1 text-sm font-medium text-ds-fg-heading">{{ account.industry || '—' }}</dd>
                </div>
                <div>
                  <dt class="text-xs text-ds-fg-muted">{{ $t('accountsFieldWebsite') }}</dt>
                  <dd class="mt-1 text-sm font-medium text-ds-fg-heading">
                    <a
                      v-if="account.website"
                      :href="websiteHref"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-ds-fg-brand hover:underline"
                    >
                      {{ account.website }}
                    </a>
                    <span v-else>—</span>
                  </dd>
                </div>
                <div>
                  <dt class="text-xs text-ds-fg-muted">{{ $t('leadsFieldLastActivity') }}</dt>
                  <dd class="mt-1 text-sm font-medium text-ds-fg-heading">{{ formatDate(account.last_activity_at) }}</dd>
                </div>
                <div>
                  <dt class="text-xs text-ds-fg-muted">{{ $t('accountsFieldUpdated') }}</dt>
                  <dd class="mt-1 text-sm font-medium text-ds-fg-heading">{{ formatDate(account.updated_at) }}</dd>
                </div>
                <div v-if="account.tags.length" class="sm:col-span-2">
                  <dt class="text-xs text-ds-fg-muted">{{ $t('leadsFieldTags') }}</dt>
                  <dd class="mt-2 flex flex-wrap gap-2">
                    <span
                      v-for="tag in account.tags"
                      :key="tag"
                      class="rounded-full bg-ds-bg-muted px-2.5 py-0.5 text-xs text-ds-fg-muted"
                    >
                      {{ tag }}
                    </span>
                  </dd>
                </div>
              </dl>
            </CardShell>

            <ContactsAccountContactsPanel
              :items="linkedContacts"
              :pending="contactsPending"
              :load-error="contactsError"
              :can-create="canCreateContact"
              @create="goCreateContact"
            />
          </template>

          <template #timeline>
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
              <div class="space-y-6" data-testid="account-activity-timeline-section">
                <div>
                  <h3 class="mb-2 text-sm font-semibold text-ds-fg-heading">{{ $t('activitySummaryTitle') }}</h3>
                  <p class="mb-3 text-xs text-ds-fg-muted">{{ $t('activitySummaryHint') }}</p>
                  <CrmActivitySummaryChart ref="summaryRef" subject-type="account" :subject-id="account.id" />
                </div>
                <div>
                  <h3 class="mb-3 text-sm font-semibold text-ds-fg-heading">{{ $t('activityTimelineTitle') }}</h3>
                  <CrmActivityTimeline ref="timelineRef" subject-type="account" :subject-id="account.id" />
                </div>
              </div>
            </CardShell>
          </template>

          <template #emotion>
            <div data-testid="tab-emotion-journey">
              <CardShell :title="$t('leadsTabEmotion')" class="rounded-2xl">
                <CrmEmotionJourneyMap ref="emotionMapRef" subject-type="account" :subject-id="account.id" />
              </CardShell>
            </div>
          </template>
        </AccountsAccountDetailTabs>
      </div>

      <AiRelationPanel
        v-if="showAiPanel"
        class="xl:sticky xl:top-4 xl:self-start"
        :show-preview="aiPreview.isPreviewMode"
        :insights="panelInsights"
        :engagement-score="panelEngagementScore"
        :churn-score="aiPreview.churnRiskScore"
        :disclaimer="aiPreview.disclaimer"
        :generate-copilot="aiPreview.generateCopilot"
      />
    </div>

    <UiModal v-model:open="activityOpen" :title="$t('activityCreateTitle')">
      <CrmActivityForm
        subject-type="account"
        :subject-id="account?.id ?? ''"
        :loading="activitySaving"
        @submit="submitActivity"
        @cancel="activityOpen = false"
      />
    </UiModal>

    <UiModal v-model:open="editOpen" :title="$t('accountsEditTitle')">
      <form class="space-y-4" @submit.prevent="submitEdit">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('accountsColName') }}</label>
          <UiInput v-model="formName" required />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('accountsColIndustry') }}</label>
          <UiInput v-model="formIndustry" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('accountsFieldWebsite') }}</label>
          <UiInput v-model="formWebsite" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('leadsColLifecycle') }}</label>
          <UiSelect v-model="formLifecycle" :items="lifecycleFormItems" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('leadsFieldTags') }}</label>
          <UiInput v-model="formTagsText" :placeholder="$t('accountsTagsHint')" />
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="secondary" @click="editOpen = false">{{ $t('cancel') }}</UiButton>
          <UiButton :loading="saving" @click="submitEdit">{{ $t('save') }}</UiButton>
        </div>
      </template>
    </UiModal>
  </PermissionGuard>
</template>

<script setup lang="ts">
import type { ActivityCreateInput } from '~/types/activity'
import type { Account, AccountDetailTab, LifecycleStage } from '~/types/account'
import type { Contact } from '~/types/contact'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const { t, locale } = useI18n()
const permission = usePermission()
const accountsApi = useAccounts()
const contactsApi = useContacts()
const activitiesApi = useActivities()
const router = useRouter()

const lifecycleOptions: LifecycleStage[] = ['acquire', 'activate', 'grow', 'retain', 'revive']

const account = ref<Account | null>(null)
const pending = ref(true)
const loadError = ref('')
const activeTab = ref<AccountDetailTab>(
  route.query.tab === 'timeline' || route.query.tab === 'emotion' ? route.query.tab : 'overview',
)
const editOpen = ref(false)
const activityOpen = ref(false)
const saving = ref(false)
const activitySaving = ref(false)
const timelineRef = ref<{ reload: () => Promise<void> } | null>(null)
const summaryRef = ref<{ reload: () => Promise<void> } | null>(null)
const emotionMapRef = ref<{ reload: () => Promise<void> } | null>(null)

const formName = ref('')
const formIndustry = ref('')
const formWebsite = ref('')
const formLifecycle = ref<LifecycleStage>('acquire')
const formTagsText = ref('')

const linkedContacts = ref<Contact[]>([])
const contactsPending = ref(false)
const contactsError = ref('')

const canUpdate = computed(() => permission.can('accounts', 'update'))
const canCreateContact = computed(() => permission.can('contacts', 'create'))
const canCreateActivity = computed(() => permission.can('activities', 'create'))
const showAiPanel = true

const accountId = computed(() => account.value?.id ?? null)
const aiPreview = useAiPreview(accountId)
const fallbackEngagement = computed(() => account.value?.engagement_score ?? null)

const previewInsights = aiPreview.previewInsights

const {
  insights: panelInsights,
  engagementScore: panelEngagementScore,
  reload: reloadInsights,
} = useDetailInsights({
  subjectType: 'account',
  subjectId: accountId,
  previewInsights,
  fallbackEngagement,
})

const websiteHref = computed(() => {
  const w = account.value?.website?.trim()
  if (!w) return ''
  return w.startsWith('http') ? w : `https://${w}`
})

const lifecycleFormItems = computed(() =>
  lifecycleOptions.map((s) => ({ label: t(`lifecycle.${s}`), value: s })),
)

function parseTags(text: string): string[] {
  return text
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
}

function formatDate(iso: string | null | undefined) {
  if (!iso) return '—'
  return new Date(iso).toLocaleString(locale.value === 'zh' ? 'zh-CN' : 'en-US')
}

function fillForm(row: Account) {
  formName.value = row.name
  formIndustry.value = row.industry ?? ''
  formWebsite.value = row.website ?? ''
  formLifecycle.value = row.lifecycle_stage
  formTagsText.value = (row.tags ?? []).join(', ')
}

function goCreateContact() {
  if (!account.value) return
  router.push({ path: '/contacts', query: { create: '1', account_id: account.value.id } })
}

async function loadLinkedContacts(accountId: string) {
  contactsPending.value = true
  contactsError.value = ''
  try {
    const { data } = await contactsApi.fetchByAccount(accountId, { page_size: 20 })
    linkedContacts.value = data.items
  } catch (e) {
    contactsError.value = e instanceof Error ? e.message : t('loadFailed')
    linkedContacts.value = []
  } finally {
    contactsPending.value = false
  }
}

async function load() {
  pending.value = true
  loadError.value = ''
  const id = route.params.id as string
  try {
    account.value = await accountsApi.fetchById(id)
    if (account.value) {
      fillForm(account.value)
      await loadLinkedContacts(account.value.id)
    }
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

async function submitActivity(payload: ActivityCreateInput) {
  if (!account.value) return
  activitySaving.value = true
  loadError.value = ''
  try {
    await activitiesApi.create(payload)
    activityOpen.value = false
    await Promise.all([
      timelineRef.value?.reload(),
      summaryRef.value?.reload(),
      reloadInsights(),
      emotionMapRef.value?.reload(),
    ])
    account.value = await accountsApi.fetchById(account.value.id)
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    activitySaving.value = false
  }
}

async function submitEdit() {
  if (!account.value || !formName.value.trim()) return
  saving.value = true
  try {
    account.value = await accountsApi.update(account.value.id, {
      name: formName.value.trim(),
      industry: formIndustry.value.trim() || undefined,
      website: formWebsite.value.trim() || undefined,
      lifecycle_stage: formLifecycle.value,
      tags: parseTags(formTagsText.value),
    })
    editOpen.value = false
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    saving.value = false
  }
}

watch(() => route.params.id, load)
onMounted(load)
</script>
