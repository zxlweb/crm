<!-- @deprecated Use lead-detail-header + lead-detail-metrics + lead-detail-tabs on pages/leads/[id].vue -->
<template>
  <CardShell
    :title="$t('leadsDecisionHubTitle')"
    :subtitle="$t('leadsDecisionHubSubtitle')"
    class="rounded-2xl"
    data-testid="lead-relationship-hub"
  >
    <dl
      class="mb-8 grid gap-4 border-b border-ds-border pb-6 sm:grid-cols-2 lg:grid-cols-4"
      data-testid="lead-overview-facts"
    >
      <div>
        <dt class="text-xs text-ds-fg-muted">{{ $t('leadsColSource') }}</dt>
        <dd class="mt-1 text-sm font-medium text-ds-fg-heading">{{ lead.source || '—' }}</dd>
      </div>
      <div>
        <dt class="text-xs text-ds-fg-muted">{{ $t('leadsFieldAmount') }}</dt>
        <dd class="mt-1 text-sm font-medium text-ds-fg-heading">{{ formatAmount(lead.amount) }}</dd>
      </div>
      <div>
        <dt class="text-xs text-ds-fg-muted">{{ $t('leadsFieldExpectedClose') }}</dt>
        <dd class="mt-1 text-sm font-medium text-ds-fg-heading">{{ lead.expected_close_date || '—' }}</dd>
      </div>
      <div>
        <dt class="text-xs text-ds-fg-muted">{{ $t('leadsFieldLastActivity') }}</dt>
        <dd class="mt-1 text-sm font-medium text-ds-fg-heading">{{ formatDate(lead.last_activity_at) }}</dd>
      </div>
      <div v-if="lead.tags.length" class="sm:col-span-2 lg:col-span-4">
        <dt class="text-xs text-ds-fg-muted">{{ $t('leadsFieldTags') }}</dt>
        <dd class="mt-2 flex flex-wrap gap-2">
          <span
            v-for="tag in lead.tags"
            :key="tag"
            class="rounded-full bg-ds-bg-muted px-2.5 py-0.5 text-xs text-ds-fg-muted"
          >
            {{ tag }}
          </span>
        </dd>
      </div>
    </dl>

    <section class="space-y-4" data-testid="tab-emotion-journey">
      <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('leadsSectionEmotion') }}</h3>
      <CrmEmotionJourneyMap
        subject-type="lead"
        :subject-id="lead.id"
        hide-touchpoints
      />
    </section>

    <hr class="my-8 border-0 border-t border-ds-border" />

    <section class="space-y-4" data-testid="lead-activity-timeline-section">
      <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('leadsSectionTimeline') }}</h3>
      <p class="text-xs text-ds-fg-muted">{{ $t('leadsSectionTimelineHint') }}</p>
      <CrmActivityTimeline :lead-id="lead.id" />
    </section>
  </CardShell>
</template>

<script setup lang="ts">
import type { Lead } from '~/types/lead'

defineProps<{
  lead: Lead
}>()

const { locale } = useI18n()

function formatAmount(value: number) {
  return new Intl.NumberFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    style: 'currency',
    currency: 'CNY',
    maximumFractionDigits: 0,
  }).format(value)
}

function formatDate(iso: string | null) {
  if (!iso) return '—'
  return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(iso))
}
</script>
