<template>
  <div v-if="auth.isSuperAdmin" class="space-y-1.5" data-testid="app-context-switch">
    <p class="px-0.5 text-[10px] font-medium uppercase tracking-wide text-ds-fg-subtle">
      {{ $t('appContextSwitchLabel') }}
    </p>
    <div
      class="inline-flex w-full items-center gap-1 rounded-xl border border-ds-border bg-ds-bg-elevated p-1"
      role="group"
      :aria-label="$t('appContextSwitchLabel')"
    >
      <NuxtLink
        to="/"
        class="flex flex-1 cursor-pointer items-center justify-center rounded-lg px-2 py-1.5 text-xs font-medium transition-colors duration-200"
        :class="isCrmContext ? 'bg-ds-brand text-ds-on-brand shadow-sm' : 'text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg'"
        :aria-current="isCrmContext ? 'page' : undefined"
      >
        {{ $t('appContextCrm') }}
      </NuxtLink>
      <NuxtLink
        to="/admin"
        class="flex flex-1 cursor-pointer items-center justify-center rounded-lg px-2 py-1.5 text-xs font-medium transition-colors duration-200"
        :class="isAdminContext ? 'bg-ds-brand text-ds-on-brand shadow-sm' : 'text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg'"
        :aria-current="isAdminContext ? 'page' : undefined"
      >
        {{ $t('appContextAdmin') }}
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const auth = useAuth()

const isAdminContext = computed(() => route.path.startsWith('/admin'))
const isCrmContext = computed(() => !isAdminContext.value)
</script>
