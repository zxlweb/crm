<template>
  <div class="space-y-6" data-testid="lead-overview-facts">
    <aside
      v-if="insightHint"
      class="flex gap-3 rounded-xl border border-ds-brand-muted bg-ds-brand-subtle/40 px-4 py-4 text-sm text-ds-fg"
    >
      <svg class="mt-0.5 h-5 w-5 shrink-0 text-ds-fg-brand" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <div>
        <p class="font-medium text-ds-fg-heading">{{ $t('leadsOverviewActionTitle') }}</p>
        <p class="mt-1 text-ds-fg-muted">{{ insightHint }}</p>
      </div>
    </aside>

    <CardShell :title="$t('leadsOverviewEmotionSection')" class="rounded-2xl">
      <LeadsLeadEmotionPreview :lead-id="lead.id" @expand="$emit('open-emotion')" />
    </CardShell>

    <CardShell :title="$t('leadsFieldTags')" class="rounded-2xl">
      <div v-if="lead.tags.length" class="flex flex-wrap gap-2">
        <span
          v-for="tag in lead.tags"
          :key="tag"
          class="rounded-full border border-ds-border bg-ds-bg-elevated px-3 py-1 text-xs font-medium text-ds-fg"
        >
          {{ tag }}
        </span>
      </div>
      <p v-else class="text-sm text-ds-fg-muted">{{ $t('leadsTagsEmpty') }}</p>
      <p class="mt-4 border-t border-ds-border pt-4 text-xs text-ds-fg-muted">
        {{ $t('leadsFieldCreated') }} {{ formatDate(lead.created_at) }}
        <span class="mx-2">·</span>
        {{ $t('leadsFieldUpdated') }} {{ formatDate(lead.updated_at) }}
      </p>
    </CardShell>
  </div>
</template>

<script setup lang="ts">
import type { Lead } from '~/types/lead'

defineProps<{
  lead: Lead
  insightHint?: string
}>()

defineEmits<{
  'open-emotion': []
}>()

const { locale } = useI18n()

function formatDate(iso: string | null) {
  if (!iso) return '—'
  return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    dateStyle: 'medium',
  }).format(new Date(iso))
}
</script>
