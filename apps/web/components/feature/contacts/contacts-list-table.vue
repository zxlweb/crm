<template>
  <UiTable
    :rows="items"
    :columns="columns"
    :empty-state="{ label: $t('contactsEmpty') }"
    data-testid="contacts-list-table"
  >
    <template v-if="$slots.toolbar" #toolbar>
      <slot name="toolbar" />
    </template>

    <template #display_name-data="{ row }">
      <NuxtLink
        :to="`/contacts/${row.id}`"
        class="ds-contacts-row__name group/title flex min-w-0 cursor-pointer items-center gap-2.5"
      >
        <span
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl text-xs font-bold tracking-tight text-ds-fg-brand ring-1 ring-inset ring-ds-brand/15 transition-[background-color,box-shadow,transform] duration-200 group-hover/title:scale-[1.04] group-hover/title:shadow-ds-sm"
          :style="listRowAvatarStyle(row.id)"
        >
          {{ initialsOf(row.display_name) }}
        </span>
        <span class="flex min-w-0 flex-col">
          <span class="flex min-w-0 items-center gap-1.5">
            <span
              class="truncate text-sm font-semibold text-ds-fg-heading transition-colors duration-200 group-hover/title:text-ds-fg-brand"
            >
              {{ row.display_name }}
            </span>
            <span
              v-if="row.is_primary"
              class="shrink-0 rounded border border-ds-brand/25 bg-ds-brand-subtle px-1.5 py-0 text-[10px] font-medium leading-4 text-ds-fg-brand"
            >
              {{ $t('contactsPrimaryBadge') }}
            </span>
          </span>
          <span
            v-if="row.email"
            class="mt-0.5 truncate text-[11px] text-ds-fg-subtle"
          >
            {{ row.email }}
          </span>
        </span>
      </NuxtLink>
    </template>

    <template #email-data="{ row }">
      <span class="text-sm text-ds-fg-muted">{{ row.email || '—' }}</span>
    </template>

    <template #account_id-data="{ row }">
      <NuxtLink
        v-if="row.account_id"
        :to="`/accounts/${row.account_id}`"
        class="inline-flex items-center gap-1 rounded-md border border-ds-border bg-ds-bg-muted px-2 py-0.5 text-[11px] font-medium text-ds-fg-brand transition-colors duration-200 hover:border-ds-brand-muted hover:bg-ds-brand-subtle"
      >
        <UIcon name="i-heroicons-building-office-2" class="h-3 w-3" aria-hidden="true" />
        {{ accountName(row.account_id) }}
      </NuxtLink>
      <span v-else class="text-xs text-ds-fg-subtle">—</span>
    </template>

    <template #lifecycle_stage-data="{ row }">
      <CrmLifecycleBadge variant="filled" :stage="row.lifecycle_stage" />
    </template>

    <template #relationship_health-data="{ row }">
      <CrmRelationshipHealthBadge variant="filled" :health="row.relationship_health" />
    </template>

    <template #engagement_score-data="{ row }">
      <div class="ds-contacts-row__engagement flex items-center justify-end gap-2">
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
          data-testid="contacts-row-edit"
          @click="$emit('edit', row)"
        />
        <CrmTableIconAction
          :to="`/contacts/${row.id}`"
          icon="i-heroicons-eye-20-solid"
          :label="$t('contactsViewDetail')"
          data-testid="contacts-row-view"
        />
      </div>
    </template>

    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </UiTable>
</template>

<script setup lang="ts">
import type { Contact } from '~/types/contact'
import {
  engagementBarClass,
  engagementTextClass,
  initialsOf,
  listRowAvatarStyle,
} from '~/utils/list-table-row'

type UiTableColumn = { key: string; label: string; sortable?: boolean; class?: string }

const props = defineProps<{
  items: Contact[]
  canEdit?: boolean
  accountNames?: Record<string, string>
}>()

defineEmits<{
  edit: [row: Contact]
}>()

const { t } = useI18n()

function accountName(id: string) {
  return props.accountNames?.[id] ?? id.slice(0, 8)
}

const columns = computed<UiTableColumn[]>(() => [
  { key: 'display_name', label: t('contactsColName'), sortable: true },
  { key: 'email', label: t('contactsColEmail'), sortable: true },
  { key: 'account_id', label: t('contactsColAccount') },
  { key: 'lifecycle_stage', label: t('leadsColLifecycle'), sortable: true },
  { key: 'relationship_health', label: t('leadsColHealth'), sortable: true },
  { key: 'engagement_score', label: t('leadsColEngagement'), sortable: true, class: 'text-right' },
  { key: 'actions', label: t('actions'), class: 'text-right w-24' },
])
</script>
