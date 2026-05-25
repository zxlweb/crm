<template>
  <div data-testid="contacts-list-table">
  <UiTable
    :rows="items"
    :columns="columns"
    :empty-state="{ label: $t('contactsEmpty') }"
  >
    <template v-if="$slots.toolbar" #toolbar>
      <slot name="toolbar" />
    </template>

    <template #display_name-data="{ row }">
      <div class="flex items-center gap-2">
        <NuxtLink
          :to="`/contacts/${row.id}`"
          class="cursor-pointer font-medium text-ds-fg-heading underline-offset-2 transition-colors hover:text-ds-fg hover:underline"
        >
          {{ row.display_name }}
        </NuxtLink>
        <span
          v-if="row.is_primary"
          class="rounded-full bg-ds-brand-subtle px-2 py-0.5 text-[10px] font-medium uppercase text-ds-fg-brand"
        >
          {{ $t('contactsPrimaryBadge') }}
        </span>
      </div>
    </template>

    <template #email-data="{ row }">
      <span class="text-ds-fg-muted">{{ row.email || '—' }}</span>
    </template>

    <template #account_id-data="{ row }">
      <NuxtLink
        v-if="row.account_id"
        :to="`/accounts/${row.account_id}`"
        class="text-sm text-ds-fg-brand hover:underline"
      >
        {{ accountName(row.account_id) }}
      </NuxtLink>
      <span v-else class="text-ds-fg-muted">—</span>
    </template>

    <template #lifecycle_stage-data="{ row }">
      <CrmLifecycleBadge variant="plain" muted-dot :stage="row.lifecycle_stage" />
    </template>

    <template #relationship_health-data="{ row }">
      <CrmRelationshipHealthBadge variant="plain" muted-dot :health="row.relationship_health" />
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
          :to="`/contacts/${row.id}`"
          icon="i-heroicons-eye-20-solid"
          :label="$t('contactsViewDetail')"
        />
      </div>
    </template>

    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </UiTable>
  </div>
</template>

<script setup lang="ts">
import type { UiTableColumn } from '@crm/ui-kit'
import type { Contact } from '~/types/contact'

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
  { key: 'actions', label: t('actions'), class: 'text-right w-24' },
])
</script>
