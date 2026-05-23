<template>
  <header
    class="relative scroll-mt-24 overflow-hidden rounded-2xl border border-ds-border bg-ds-bg-elevated shadow-ds-sm"
    data-testid="lead-detail-header"
  >
    <button
      v-if="canUpdate && lead.status !== 'converted'"
      type="button"
      class="absolute right-4 top-4 z-10 flex h-9 w-9 cursor-pointer items-center justify-center rounded-lg border border-ds-border bg-ds-bg-elevated text-ds-fg-brand transition-colors duration-200 hover:border-ds-brand-muted hover:bg-ds-brand-subtle sm:right-5 sm:top-3"
      data-testid="lead-edit-btn"
      :aria-label="$t('edit')"
      @click="$emit('edit')"
    >
      <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
      </svg>
    </button>

    <div
      class="relative flex flex-wrap items-center justify-between gap-3 border-b border-ds-border px-4 py-3 pr-14 sm:px-5 sm:pr-16"
    >
      <NuxtLink
        to="/leads"
        class="inline-flex shrink-0 cursor-pointer items-center gap-1.5 text-xs font-medium text-ds-fg-brand transition-colors duration-200 hover:text-ds-brand-hover"
      >
        <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        {{ $t('leadsBackToList') }}
      </NuxtLink>

      <div class="flex shrink-0 flex-wrap items-center justify-end gap-2">
        <NuxtLink
          v-if="lead.status === 'converted' && lead.converted_account_id"
          :to="`/accounts/${lead.converted_account_id}`"
          class="ds-btn-secondary inline-flex cursor-pointer items-center rounded-lg px-3 py-1.5 text-xs font-medium sm:text-sm"
          data-testid="lead-view-account-btn"
        >
          {{ $t('leadsViewAccount') }}
        </NuxtLink>
        <UiButton
          v-else-if="canConvert"
          data-testid="lead-convert-btn"
          :loading="convertSaving"
          @click="$emit('convert')"
        >
          {{ $t('leadsConvert') }}
        </UiButton>
        <UiSelect
          v-if="canUpdate && lead.status !== 'converted'"
          :model-value="lead.status"
          :items="statusItems"
          class="w-[7.5rem]"
          :disabled="statusItems.length <= 1 || statusSaving"
          data-testid="lead-status-select"
          @update:model-value="$emit('status-change', $event as LeadStatus)"
        />
      </div>
    </div>

    <div class="relative px-4 py-4 sm:px-5">
      <div class="flex flex-wrap items-center gap-2 pr-10 sm:pr-12">
        <h1 class="text-xl font-bold tracking-tight text-ds-fg-heading sm:text-2xl">
          {{ lead.title }}
        </h1>
        <AiPreviewBadge v-if="showPreview" />
      </div>

      <ul
        v-if="headerTags.length"
        class="mt-2.5 flex flex-wrap gap-1.5"
        data-testid="lead-header-tags"
      >
        <li
          v-for="tag in headerTags"
          :key="tag.id"
          class="list-none rounded-md border px-2.5 py-0.5 text-xs font-medium"
          :class="tagToneClass(tag.kind)"
        >
          {{ tag.label }}
        </li>
      </ul>

      <p
        class="mt-3 text-xs text-ds-fg-muted sm:text-sm"
        data-testid="lead-header-meta"
      >
        <span>{{ $t('leadsHeaderUpdated', { date: formatDate(lead.updated_at) }) }}</span>
        <template v-if="statusLockedHint && !canUpdate">
          <span class="mx-2 hidden text-ds-fg-subtle sm:inline" aria-hidden="true">·</span>
          <span class="mt-1 block text-ds-danger sm:mt-0 sm:inline">{{ statusLockedHint }}</span>
        </template>
      </p>

      <div
        v-if="ownerProfile"
        class="mt-2.5 border-t border-ds-border/60 pt-2.5"
        data-testid="lead-header-owner"
      >
        <LeadsLeadOwnerChip :profile="ownerProfile" variant="inline" />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import type { OwnerProfile } from '~/composables/use-owner-profile'
import { allowedLeadStatusTargets, canConvertLead } from '~/utils/lead-status-transition'
import type { Lead, LeadStatus } from '~/types/lead'

const LOW_INTENT_TAG = '低意向'

type HeaderTagKind = 'risk' | 'intent' | 'source'

type HeaderTag = {
  id: string
  label: string
  kind: HeaderTagKind
}

const props = defineProps<{
  lead: Lead
  ownerProfile?: OwnerProfile | null
  canUpdate?: boolean
  showPreview?: boolean
  statusSaving?: boolean
  convertSaving?: boolean
}>()

defineEmits<{
  edit: []
  convert: []
  'status-change': [status: LeadStatus]
}>()

const { t, locale } = useI18n()

const canConvert = computed(() => !!props.canUpdate && canConvertLead(props.lead.status))

const sourceLabel = computed(() => {
  const key = `leadSource.${props.lead.source}`
  const translated = t(key)
  return translated === key ? props.lead.source || '—' : translated
})

/** 头部仅展示：风险 · 低意向 · 来源 */
const headerTags = computed<HeaderTag[]>(() => {
  const tags: HeaderTag[] = []

  if (props.lead.relationship_health === 'low') {
    tags.push({
      id: 'risk',
      kind: 'risk',
      label: t('relationshipHealth.low'),
    })
  }

  const lowIntent = props.lead.tags.find(
    (tag) => tag === LOW_INTENT_TAG || tag.toLowerCase() === 'low intent',
  )
  if (lowIntent) {
    tags.push({
      id: 'intent',
      kind: 'intent',
      label: lowIntent === LOW_INTENT_TAG ? t('leadsTagLowIntent') : lowIntent,
    })
  }

  if (props.lead.source) {
    tags.push({
      id: 'source',
      kind: 'source',
      label: t('leadsHeaderTagSource', { source: sourceLabel.value }),
    })
  }

  return tags
})

function tagToneClass(kind: HeaderTagKind) {
  const map: Record<HeaderTagKind, string> = {
    risk: 'border-ds-danger/25 bg-ds-danger-subtle/60 text-ds-danger',
    intent: 'border-ds-border bg-ds-bg-muted text-ds-fg-muted',
    source: 'border-ds-border bg-ds-bg-muted text-ds-fg-muted',
  }
  return map[kind]
}

const statusItems = computed(() =>
  allowedLeadStatusTargets(props.lead.status).map((s) => ({
    label: t(`leadStatus.${s}`),
    value: s,
  })),
)

const statusLockedHint = computed(() => {
  if (props.canUpdate) return ''
  if (props.lead.status === 'converted') return t('leadsStatusLockedConverted')
  if (props.lead.status === 'unqualified') return t('leadsStatusLockedTerminal')
  return ''
})

function formatDate(iso: string | null | undefined) {
  if (!iso) return '—'
  return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    dateStyle: 'medium',
  }).format(new Date(iso))
}
</script>
