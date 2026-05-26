<template>
  <div class="space-y-6" data-testid="audit-page">
    <div>
      <h1 class="text-2xl font-bold text-ds-fg-heading">{{ $t('auditPageTitle') }}</h1>
      <p class="mt-1 text-sm text-ds-fg-muted">{{ $t('auditPageDesc') }}</p>
    </div>

    <p v-if="!canView" class="rounded-xl bg-ds-danger-subtle px-4 py-3 text-sm text-ds-danger">
      {{ $t('noPermission') }}
    </p>

    <AuditStatsPanel v-else />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'app', middleware: 'auth' })

const { can } = usePermission()
const canView = computed(() => can('audit', 'view'))
</script>
