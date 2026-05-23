<template>
  <header class="sticky top-0 z-10 flex h-16 shrink-0 items-center justify-between border-b border-ds-border bg-ds-bg-topbar px-6 backdrop-blur-md lg:px-8">
    <div class="min-w-0">
      <p class="text-xs font-medium uppercase tracking-wide text-ds-fg-muted">{{ $t('crmBreadcrumb') }}</p>
      <h1 class="truncate text-lg font-semibold text-ds-fg-heading">{{ pageTitle }}</h1>
    </div>
    <div class="flex items-center gap-3">
      <AppTenantSwitcher :show-label="false" />
      <UBadge v-if="leads.useMock.value" color="amber" variant="subtle" size="sm">
        {{ $t('crmMockBadge') }}
      </UBadge>
      <span v-if="auth.user" class="hidden text-sm text-ds-fg-muted sm:inline">{{ auth.user.email }}</span>
    </div>
  </header>
</template>

<script setup lang="ts">
const route = useRoute()
const { t } = useI18n()
const auth = useAuth()
const leads = useLeads()

const pageTitle = computed(() => {
  if (route.path === '/') {
    return t('dashboardPageTitle')
  }
  if (route.path.startsWith('/accounts/') && route.params.id) {
    return t('accountsDetailTitle')
  }
  if (route.path === '/accounts') {
    return t('accountsPageTitle')
  }
  if (route.path.startsWith('/leads/') && route.params.id) {
    return t('leadsDetailTitle')
  }
  if (route.path === '/leads') {
    return t('leadsPageTitle')
  }
  return t('crmBreadcrumb')
})
</script>
