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
          </template>

          <template #timeline>
            <CardShell :title="$t('leadsTabTimeline')" class="rounded-2xl">
              <p class="text-sm text-ds-fg-muted">{{ $t('accountsTimelinePlaceholder') }}</p>
            </CardShell>
          </template>

          <template #emotion>
            <div data-testid="tab-emotion-journey">
              <CardShell :title="$t('leadsTabEmotion')" class="rounded-2xl">
                <CrmEmotionJourneyMap subject-type="account" :subject-id="account.id" />
              </CardShell>
            </div>
          </template>
        </AccountsAccountDetailTabs>
      </div>

      <AiRelationPanel
        v-if="showAiPanel"
        :show-preview="isPreviewMode"
        :insights="stubInsights"
        :engagement-score="account.engagement_score"
      />
    </div>

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
import type { Account, AccountDetailTab, LifecycleStage } from '~/types/account'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const { t, locale } = useI18n()
const permission = usePermission()
const accountsApi = useAccounts()

const lifecycleOptions: LifecycleStage[] = ['acquire', 'activate', 'grow', 'retain', 'revive']

const account = ref<Account | null>(null)
const pending = ref(true)
const loadError = ref('')
const activeTab = ref<AccountDetailTab>('overview')
const editOpen = ref(false)
const saving = ref(false)

const formName = ref('')
const formIndustry = ref('')
const formWebsite = ref('')
const formLifecycle = ref<LifecycleStage>('acquire')
const formTagsText = ref('')

const canUpdate = computed(() => permission.can('accounts', 'update'))
const isPreviewMode = computed(() => route.query.preview === '1')
const showAiPanel = true

const websiteHref = computed(() => {
  const w = account.value?.website?.trim()
  if (!w) return ''
  return w.startsWith('http') ? w : `https://${w}`
})

const lifecycleFormItems = computed(() =>
  lifecycleOptions.map((s) => ({ label: t(`lifecycle.${s}`), value: s })),
)

const stubInsights = computed(() => {
  if (!isPreviewMode.value) return []
  return [
    { id: 'INS-A1', title: t('aiStubInsight1Title'), body: t('aiStubInsight1Body') },
    { id: 'INS-A2', title: t('aiStubInsight2Title'), body: t('aiStubInsight2Body') },
  ]
})

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

async function load() {
  pending.value = true
  loadError.value = ''
  const id = route.params.id as string
  try {
    account.value = await accountsApi.fetchById(id)
    if (account.value) fillForm(account.value)
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
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
