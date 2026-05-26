<template>
  <header
    class="ds-lead-hero relative scroll-mt-24 overflow-hidden rounded-3xl border border-ds-border bg-ds-bg-elevated/85 px-5 py-5 shadow-ds-md backdrop-blur-md sm:px-7 sm:py-6"
    data-testid="lead-detail-header"
  >
    <div
      class="pointer-events-none absolute -right-12 -top-12 h-56 w-56 rounded-full opacity-30 blur-3xl"
      :style="{ background: heroGlowBackground }"
      aria-hidden="true"
    />
    <div
      class="pointer-events-none absolute -bottom-20 left-1/4 h-44 w-44 rounded-full opacity-20 blur-3xl"
      :style="{ background: 'var(--ds-blur-brand)' }"
      aria-hidden="true"
    />

    <!-- Breadcrumb + actions row -->
    <div class="relative flex flex-wrap items-center justify-between gap-3">
      <NuxtLink
        to="/leads"
        class="inline-flex shrink-0 cursor-pointer items-center gap-1.5 rounded-lg px-2 py-1 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-brand"
      >
        <UIcon name="i-heroicons-arrow-left-20-solid" class="h-3.5 w-3.5" aria-hidden="true" />
        {{ $t('leadsBackToList') }}
      </NuxtLink>

      <div class="flex shrink-0 flex-wrap items-center gap-2">
        <button
          v-if="canUpdate && lead.status !== 'converted'"
          type="button"
          class="inline-flex h-9 w-9 cursor-pointer items-center justify-center rounded-xl border border-ds-border bg-ds-bg-elevated text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand-muted hover:bg-ds-brand-subtle/40 hover:text-ds-fg-brand"
          data-testid="lead-edit-btn"
          :aria-label="$t('edit')"
          @click="$emit('edit')"
        >
          <UIcon name="i-heroicons-pencil-square" class="h-4 w-4" aria-hidden="true" />
        </button>
        <UiSelect
          v-if="canUpdate && lead.status !== 'converted'"
          :model-value="lead.status"
          :items="statusItems"
          class="w-[8rem]"
          :disabled="statusItems.length <= 1 || statusSaving"
          data-testid="lead-status-select"
          @update:model-value="$emit('status-change', $event as LeadStatus)"
        />
        <NuxtLink
          v-if="lead.status === 'converted' && lead.converted_account_id"
          :to="`/accounts/${lead.converted_account_id}`"
          class="inline-flex cursor-pointer items-center gap-1.5 rounded-xl border border-ds-success/30 bg-ds-success-subtle px-3 py-2 text-xs font-semibold text-ds-success transition-colors duration-200 hover:border-ds-success/50"
          data-testid="lead-view-account-btn"
        >
          <UIcon name="i-heroicons-arrow-top-right-on-square" class="h-3.5 w-3.5" aria-hidden="true" />
          {{ $t('leadsViewAccount') }}
        </NuxtLink>
        <button
          v-else-if="canConvert"
          type="button"
          class="ds-lead-hero__cta group relative inline-flex cursor-pointer items-center gap-1.5 overflow-hidden rounded-xl px-4 py-2 text-xs font-semibold text-ds-on-brand shadow-ds-brand transition-[transform,box-shadow] duration-200 hover:shadow-ds-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ds-brand focus-visible:ring-offset-2 focus-visible:ring-offset-ds-bg disabled:cursor-not-allowed disabled:opacity-60"
          :style="{ background: 'var(--ds-brand-gradient)' }"
          data-testid="lead-convert-btn"
          :disabled="convertSaving"
          @click="$emit('convert')"
        >
          <span
            class="pointer-events-none absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/30 to-transparent opacity-0 transition-[transform,opacity] duration-500 group-hover:translate-x-full group-hover:opacity-100"
            aria-hidden="true"
          />
          <UIcon name="i-heroicons-arrow-right-circle" class="h-4 w-4" aria-hidden="true" />
          <span>{{ $t('leadsConvert') }}</span>
        </button>
      </div>
    </div>

    <!-- Title + identity row -->
    <div class="relative mt-4 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
      <div class="min-w-0 flex-1">
        <div
          class="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-[0.18em]"
          :class="healthAccentText"
        >
          <span class="relative flex h-2 w-2" aria-hidden="true">
            <span
              v-if="lead.relationship_health === 'low'"
              class="absolute inline-flex h-full w-full animate-ping rounded-full opacity-60"
              :class="healthDotBg"
            />
            <span class="relative inline-flex h-2 w-2 rounded-full" :class="healthDotBg" />
          </span>
          <span>{{ healthEyebrowLabel }}</span>
          <span aria-hidden="true" class="text-ds-fg-subtle">·</span>
          <span class="text-ds-fg-muted">
            {{ $t(`lifecycle.${lead.lifecycle_stage}`) }}
          </span>
          <AiPreviewBadge v-if="showPreview" />
        </div>

        <div class="mt-1.5 flex items-center gap-2">
          <h1
            class="truncate text-2xl font-bold tracking-tight text-ds-fg-heading sm:text-3xl"
          >
            {{ lead.title }}
          </h1>
        </div>

        <!-- Headline amount -->
        <div class="mt-4 flex flex-wrap items-end gap-x-4 gap-y-2">
          <div class="min-w-0">
            <p class="text-[11px] font-medium uppercase tracking-wider text-ds-fg-muted">
              {{ $t('leadsHeroAmountLabel') }}
            </p>
            <p
              class="mt-0.5 text-3xl font-extrabold tabular-nums tracking-tight text-ds-fg-heading sm:text-4xl"
            >
              <span v-if="lead.amount > 0" class="bg-clip-text" :style="brandGradientText">
                {{ amountDisplay }}
              </span>
              <span v-else class="text-ds-fg-muted">{{ $t('leadsAmountUnset') }}</span>
            </p>
          </div>

          <div class="flex flex-wrap items-center gap-1.5">
            <span
              class="inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-medium"
              :class="engagementChipClass"
            >
              <UIcon name="i-heroicons-bolt" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ $t('leadsHeroEngagement', { score: lead.engagement_score }) }}
            </span>
            <span
              class="inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-medium"
              :class="idleChipClass"
            >
              <UIcon name="i-heroicons-clock" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ idleChipLabel }}
            </span>
            <span
              v-if="closeChipLabel"
              class="inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-medium"
              :class="closeChipClass"
            >
              <UIcon name="i-heroicons-calendar-days" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ closeChipLabel }}
            </span>
            <span
              v-if="lead.source"
              class="inline-flex items-center gap-1.5 rounded-full border border-ds-border bg-ds-bg-muted px-2.5 py-1 text-xs font-medium text-ds-fg-muted"
            >
              <UIcon name="i-heroicons-globe-alt" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ sourceLabel }}
            </span>
          </div>
        </div>

        <ul
          v-if="extraTags.length"
          class="mt-3 flex flex-wrap gap-1.5"
          data-testid="lead-header-tags"
        >
          <li
            v-for="tag in extraTags"
            :key="tag"
            class="list-none rounded-md border border-ds-border bg-ds-bg-muted/60 px-2 py-0.5 text-[11px] font-medium text-ds-fg-muted"
          >
            {{ tag }}
          </li>
        </ul>
      </div>

      <div class="flex shrink-0 flex-col items-start gap-2 lg:items-end">
        <div
          v-if="ownerProfile"
          data-testid="lead-header-owner"
        >
          <LeadsLeadOwnerChip :profile="ownerProfile" variant="inline" />
        </div>
        <p
          class="text-[11px] text-ds-fg-muted"
          data-testid="lead-header-meta"
        >
          <span>{{ $t('leadsHeaderUpdated', { date: updatedLabel }) }}</span>
        </p>
        <p
          v-if="statusLockedHint && !canUpdate"
          class="text-[11px] text-ds-danger"
        >
          {{ statusLockedHint }}
        </p>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import type { OwnerProfile } from '~/composables/use-owner-profile'
import { allowedLeadStatusTargets, canConvertLead } from '~/utils/lead-status-transition'
import type { Lead, LeadStatus } from '~/types/lead'

const LOW_INTENT_TAG = '低意向'

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

const brandGradientText = {
  background: 'var(--ds-brand-gradient)',
  '-webkit-background-clip': 'text',
  'background-clip': 'text',
  color: 'transparent',
}

const canConvert = computed(() => Boolean(props.canUpdate) && canConvertLead(props.lead.status))

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

const sourceLabel = computed(() => {
  const key = `leadSource.${props.lead.source}`
  const translated = t(key)
  return translated === key ? props.lead.source : translated
})

const extraTags = computed(() =>
  (props.lead.tags ?? []).filter((tag) => {
    const normalized = tag.toLowerCase()
    return tag !== LOW_INTENT_TAG && normalized !== 'low intent'
  }),
)

const amountDisplay = computed(() => {
  if (props.lead.amount <= 0) return t('leadsAmountUnset')
  return new Intl.NumberFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    style: 'currency',
    currency: 'CNY',
    maximumFractionDigits: 0,
  }).format(props.lead.amount)
})

const daysIdle = computed<number | null>(() => {
  if (!props.lead.last_activity_at) return null
  return Math.max(
    0,
    Math.floor((Date.now() - new Date(props.lead.last_activity_at).getTime()) / 86_400_000),
  )
})

const daysToClose = computed<number | null>(() => {
  if (!props.lead.expected_close_date) return null
  const close = new Date(`${props.lead.expected_close_date}T00:00:00`).getTime()
  if (Number.isNaN(close)) return null
  return Math.round((close - Date.now()) / 86_400_000)
})

const engagementChipClass = computed(() => {
  if (props.lead.engagement_score >= 60) {
    return 'border-ds-success/30 bg-ds-success-subtle text-ds-success'
  }
  if (props.lead.engagement_score < 35) {
    return 'border-ds-danger/30 bg-ds-danger-subtle text-ds-danger'
  }
  return 'border-ds-warning/30 bg-ds-warning-subtle text-ds-warning'
})

const idleChipLabel = computed(() => {
  if (daysIdle.value == null) return t('leadsMetricNoActivity')
  if (daysIdle.value === 0) return t('leadsHeroIdleToday')
  return t('leadsMetricDaysValue', { days: daysIdle.value })
})

const idleChipClass = computed(() => {
  if (daysIdle.value == null) {
    return 'border-ds-border bg-ds-bg-muted text-ds-fg-muted'
  }
  if (daysIdle.value > 14) return 'border-ds-danger/30 bg-ds-danger-subtle text-ds-danger'
  if (daysIdle.value > 7) return 'border-ds-warning/30 bg-ds-warning-subtle text-ds-warning'
  return 'border-ds-success/30 bg-ds-success-subtle text-ds-success'
})

const closeChipLabel = computed(() => {
  if (daysToClose.value == null) return ''
  if (daysToClose.value < 0) return t('leadsHeroOverdueBy', { n: Math.abs(daysToClose.value) })
  if (daysToClose.value === 0) return t('leadsHeroDueToday')
  return t('leadsHeroDueInDays', { n: daysToClose.value })
})

const closeChipClass = computed(() => {
  if (daysToClose.value == null) return 'border-ds-border bg-ds-bg-muted text-ds-fg-muted'
  if (daysToClose.value < 0) return 'border-ds-danger/30 bg-ds-danger-subtle text-ds-danger'
  if (daysToClose.value <= 7) return 'border-ds-warning/30 bg-ds-warning-subtle text-ds-warning'
  return 'border-ds-info/30 bg-ds-info-subtle text-ds-info'
})

const updatedLabel = computed(() => {
  const iso = props.lead.updated_at
  if (!iso) return '—'
  return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    dateStyle: 'medium',
  }).format(new Date(iso))
})

const healthAccentText = computed(() => {
  switch (props.lead.relationship_health) {
    case 'low':
      return 'text-ds-danger'
    case 'high':
      return 'text-ds-success'
    default:
      return 'text-ds-fg-brand'
  }
})

const healthDotBg = computed(() => {
  switch (props.lead.relationship_health) {
    case 'low':
      return 'bg-ds-danger'
    case 'high':
      return 'bg-ds-success'
    default:
      return 'bg-ds-brand'
  }
})

const healthEyebrowLabel = computed(() => {
  switch (props.lead.relationship_health) {
    case 'low':
      return t('leadsHeroEyebrowAtRisk')
    case 'high':
      return t('leadsHeroEyebrowHealthy')
    default:
      return t('leadsHeroEyebrowEngaging')
  }
})

const heroGlowBackground = computed(() => {
  if (props.lead.relationship_health === 'low') {
    return 'radial-gradient(ellipse at center, var(--ds-danger), transparent 70%)'
  }
  return 'var(--ds-brand-gradient)'
})
</script>
