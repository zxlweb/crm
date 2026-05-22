<template>
  <div class="ds-card overflow-hidden rounded-2xl shadow-sm">
    <div class="overflow-x-auto">
      <table class="w-full min-w-[720px] text-left text-sm">
        <thead>
          <tr class="border-b border-ds-border-muted bg-ds-bg-muted text-xs font-medium uppercase tracking-wide text-ds-fg-brand">
            <th class="px-5 py-3">{{ $t('accountsColName') }}</th>
            <th class="px-5 py-3">{{ $t('accountsColIndustry') }}</th>
            <th class="px-5 py-3">{{ $t('leadsColLifecycle') }}</th>
            <th class="px-5 py-3">{{ $t('leadsColHealth') }}</th>
            <th class="px-5 py-3">{{ $t('leadsColEngagement') }}</th>
            <th class="px-5 py-3 text-right">{{ $t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-ds-border-muted">
          <tr
            v-for="row in items"
            :key="row.id"
            class="transition-colors duration-200 hover:bg-ds-bg-muted"
          >
            <td class="px-5 py-4">
              <NuxtLink
                :to="`/accounts/${row.id}`"
                class="font-medium text-ds-fg-heading transition-colors hover:text-ds-fg-brand"
              >
                {{ row.name }}
              </NuxtLink>
            </td>
            <td class="px-5 py-4 text-ds-fg-muted">{{ row.industry || '—' }}</td>
            <td class="px-5 py-4">
              <CrmLifecycleBadge :stage="row.lifecycle_stage" />
            </td>
            <td class="px-5 py-4">
              <CrmRelationshipHealthBadge :health="row.relationship_health" />
            </td>
            <td class="px-5 py-4 text-ds-fg-muted">{{ row.engagement_score }}</td>
            <td class="px-5 py-4 text-right">
              <div class="flex items-center justify-end gap-3">
                <button
                  v-if="canEdit"
                  type="button"
                  class="text-xs font-medium text-ds-fg-muted transition-colors hover:text-ds-fg-brand"
                  @click="$emit('edit', row)"
                >
                  {{ $t('edit') }}
                </button>
                <NuxtLink
                  :to="`/accounts/${row.id}`"
                  class="text-xs font-medium text-ds-fg-brand transition-colors hover:text-ds-brand-hover"
                >
                  {{ $t('accountsViewDetail') }}
                </NuxtLink>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p v-if="items.length === 0" class="py-12 text-center text-sm text-ds-fg-muted">{{ $t('accountsEmpty') }}</p>
  </div>
</template>

<script setup lang="ts">
import type { Account } from '~/types/account'

defineProps<{
  items: Account[]
  canEdit?: boolean
}>()

defineEmits<{
  edit: [row: Account]
}>()
</script>
