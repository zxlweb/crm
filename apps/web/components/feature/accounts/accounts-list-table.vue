<template>
  <UiTable
    :rows="items"
    :columns="columns"
    :empty-state="{ label: $t('accountsEmpty') }"
    data-testid="accounts-list-table"
  >
    <template v-if="$slots.toolbar" #toolbar>
      <slot name="toolbar" />
    </template>

    <template #name-data="{ row }">
      <NuxtLink
        :to="`/accounts/${row.id}`"
        class="cursor-pointer font-medium text-ds-fg-heading underline-offset-2 transition-colors duration-200 hover:text-ds-fg hover:underline"
      >
        {{ row.name }}
      </NuxtLink>
    </template>

    <template #industry-data="{ row }">
      <span class="text-ds-fg-muted">{{ row.industry || '—' }}</span>
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

    <template #actions-data="{ row }">
      <div class="flex items-center justify-end gap-0.5">
        <CrmTableIconAction
          v-if="canEdit"
          icon="i-heroicons-pencil-square-20-solid"
          :label="$t('edit')"
          @click="$emit('edit', row)"
        />
        <CrmTableIconAction
          :to="`/accounts/${row.id}`"
          icon="i-heroicons-eye-20-solid"
          :label="$t('accountsViewDetail')"
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
import type { Account } from '~/types/account'

defineProps<{
  items: Account[]
  canEdit?: boolean
}>()

defineEmits<{
  edit: [row: Account]
}>()

const { t } = useI18n()

const columns = computed<UiTableColumn[]>(() => [
  { key: 'name', label: t('accountsColName'), sortable: true },
  { key: 'industry', label: t('accountsColIndustry'), sortable: true },
  { key: 'lifecycle_stage', label: t('leadsColLifecycle'), sortable: true },
  { key: 'relationship_health', label: t('leadsColHealth'), sortable: true },
  { key: 'engagement_score', label: t('leadsColEngagement'), sortable: true, class: 'text-right' },
  { key: 'actions', label: t('actions'), class: 'text-right w-24' },
])
</script>
