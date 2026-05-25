<template>
  <CardShell :title="$t('contactsLinkedTitle')" class="rounded-2xl" data-testid="account-contacts-panel">
    <div v-if="pending" class="flex justify-center py-8">
      <UIcon name="i-heroicons-arrow-path" class="h-6 w-6 animate-spin text-primary" />
    </div>
    <p v-else-if="loadError" class="text-sm text-ds-danger">{{ loadError }}</p>
    <p v-else-if="items.length === 0" class="text-sm text-ds-fg-muted">{{ $t('contactsLinkedEmpty') }}</p>
    <ul v-else class="divide-y divide-ds-border">
      <li v-for="row in items" :key="row.id" class="flex items-center justify-between gap-3 py-3">
        <div class="min-w-0">
          <NuxtLink :to="`/contacts/${row.id}`" class="font-medium text-ds-fg-heading hover:underline">
            {{ row.display_name }}
          </NuxtLink>
          <p class="truncate text-xs text-ds-fg-muted">{{ row.email || row.phone || '—' }}</p>
        </div>
        <span
          v-if="row.is_primary"
          class="shrink-0 rounded-full bg-ds-brand-subtle px-2 py-0.5 text-[10px] font-medium text-ds-fg-brand"
        >
          {{ $t('contactsPrimaryBadge') }}
        </span>
      </li>
    </ul>
    <template v-if="canCreate" #header-extra>
      <div class="mt-3 flex justify-end">
        <UiButton
          size="sm"
          variant="secondary"
          icon="i-heroicons-plus-20-solid"
          data-testid="account-contact-create-btn"
          @click="$emit('create')"
        >
          {{ $t('contactsCreate') }}
        </UiButton>
      </div>
    </template>
  </CardShell>
</template>

<script setup lang="ts">
import type { Contact } from '~/types/contact'

defineProps<{
  items: Contact[]
  pending?: boolean
  loadError?: string
  canCreate?: boolean
}>()

defineEmits<{
  create: []
}>()
</script>
