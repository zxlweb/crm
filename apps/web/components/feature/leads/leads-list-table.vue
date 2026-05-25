<template>
  <UiTable
    :rows="items"
    :columns="columns"
    :empty-state="{ label: $t('leadsEmpty') }"
    data-testid="leads-list-table"
  >
    <template v-if="$slots.toolbar" #toolbar>
      <slot name="toolbar" />
    </template>

    <template #title-data="{ row }">
      <NuxtLink
        :to="`/leads/${row.id}`"
        class="cursor-pointer font-medium text-ds-fg-heading underline-offset-2 transition-colors duration-200 hover:text-ds-fg hover:underline"
      >
        {{ row.title }}
      </NuxtLink>
    </template>

    <template #status-data="{ row }">
      <CrmLeadStatusBadge variant="plain" muted-dot :status="row.status" />
    </template>

    <template #lifecycle_stage-data="{ row }">
      <CrmLifecycleBadge variant="plain" muted-dot :stage="row.lifecycle_stage" />
    </template>

    <template #relationship_health-data="{ row }">
      <CrmRelationshipHealthBadge variant="plain" muted-dot :health="row.relationship_health" />
    </template>

    <template #engagement_score-data="{ row }">
      <span class="tabular-nums text-ds-fg">{{ row.engagement_score }}</span>
    </template>

    <template #source-data="{ row }">
      <span class="text-ds-fg-muted">{{ formatSource(row.source) }}</span>
    </template>

    <template #actions-data="{ row }">
      <div class="flex items-center justify-end gap-0.5">
        <CrmTableIconAction
          :to="`/leads/${row.id}`"
          icon="i-heroicons-eye-20-solid"
          :label="$t('leadsViewDetail')"
        />
      </div>
    </template>

    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </UiTable>
</template>

<script setup lang="ts">
import type { UiTableColumn } from '@crm/ui-kit'
import type { Lead } from '~/types/lead'

defineProps<{
  items: Lead[]
}>()

const { t } = useI18n()

const columns = computed<UiTableColumn[]>(() => [
  { key: 'title', label: t('leadsColTitle'), sortable: true },
  { key: 'status', label: t('status'), sortable: true },
  { key: 'lifecycle_stage', label: t('leadsColLifecycle'), sortable: true },
  { key: 'relationship_health', label: t('leadsColHealth'), sortable: true },
  { key: 'engagement_score', label: t('leadsColEngagement'), sortable: true, class: 'text-right' },
  { key: 'source', label: t('leadsColSource'), sortable: true },
  { key: 'actions', label: t('actions'), class: 'text-right w-16' },
])

const { leadSourceLabel } = useLeadLabels()

function formatSource(source: string) {
  return leadSourceLabel(source) || '—'
}
</script>
