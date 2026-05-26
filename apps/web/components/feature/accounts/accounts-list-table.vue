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
        class="ds-accounts-row__name group/title flex min-w-0 cursor-pointer items-center gap-2.5"
      >
        <span
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl text-xs font-bold tracking-tight text-ds-fg-brand ring-1 ring-inset ring-ds-brand/15 transition-[background-color,box-shadow,transform] duration-200 group-hover/title:scale-[1.04] group-hover/title:shadow-ds-sm"
          :style="listRowAvatarStyle(row.id)"
        >
          {{ initialsOf(row.name) }}
        </span>
        <span class="flex min-w-0 flex-col">
          <span
            class="truncate text-sm font-semibold text-ds-fg-heading transition-colors duration-200 group-hover/title:text-ds-fg-brand"
          >
            {{ row.name }}
          </span>
          <span
            v-if="row.industry"
            class="mt-0.5 truncate text-[11px] text-ds-fg-subtle"
          >
            {{ row.industry }}
          </span>
        </span>
      </NuxtLink>
    </template>

    <template #industry-data="{ row }">
      <span
        class="inline-flex items-center gap-1 rounded-md border border-ds-border bg-ds-bg-muted px-2 py-0.5 text-[11px] font-medium text-ds-fg-muted"
      >
        <UIcon name="i-heroicons-building-office-2" class="h-3 w-3" aria-hidden="true" />
        {{ row.industry || '—' }}
      </span>
    </template>

    <template #lifecycle_stage-data="{ row }">
      <CrmLifecycleBadge variant="filled" :stage="row.lifecycle_stage" />
    </template>

    <template #relationship_health-data="{ row }">
      <CrmRelationshipHealthBadge variant="filled" :health="row.relationship_health" />
    </template>

    <template #engagement_score-data="{ row }">
      <div class="ds-accounts-row__engagement flex items-center justify-end gap-2">
        <div
          class="relative h-1.5 w-16 overflow-hidden rounded-full bg-ds-bg-muted"
          aria-hidden="true"
        >
          <span
            class="absolute inset-y-0 left-0 rounded-full transition-[width,background-color] duration-500"
            :class="engagementBarClass(row.engagement_score)"
            :style="{ width: `${Math.max(4, Math.min(100, row.engagement_score))}%` }"
          />
        </div>
        <span
          class="min-w-[2.25rem] text-right text-sm font-semibold tabular-nums"
          :class="engagementTextClass(row.engagement_score)"
        >
          {{ row.engagement_score }}
        </span>
      </div>
    </template>

    <template #actions-data="{ row }">
      <div class="flex items-center justify-end gap-0.5">
        <CrmTableIconAction
          v-if="canEdit"
          icon="i-heroicons-pencil-square-20-solid"
          :label="$t('edit')"
          data-testid="accounts-row-edit"
          @click="$emit('edit', row)"
        />
        <CrmTableIconAction
          :to="`/accounts/${row.id}`"
          icon="i-heroicons-eye-20-solid"
          :label="$t('accountsViewDetail')"
          data-testid="accounts-row-view"
        />
      </div>
    </template>

    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </UiTable>
</template>

<script setup lang="ts">
import type { Account } from '~/types/account'
import {
  engagementBarClass,
  engagementTextClass,
  initialsOf,
  listRowAvatarStyle,
} from '~/utils/list-table-row'

type UiTableColumn = { key: string; label: string; sortable?: boolean; class?: string }

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
