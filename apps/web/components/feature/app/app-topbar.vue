<template>
  <header
    class="sticky top-0 z-10 flex h-16 shrink-0 items-center justify-between gap-4 border-b border-ds-border bg-ds-bg-topbar px-6 backdrop-blur-md lg:px-8"
  >
    <div class="min-w-0 flex-1">
      <p class="text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle">
        {{ $t('crmBreadcrumb') }}
      </p>
      <h1 class="truncate text-lg font-semibold text-ds-fg-heading leading-tight">
        {{ pageTitle }}
      </h1>
    </div>

    <div class="flex items-center gap-2 sm:gap-3">
      <UBadge v-if="leads.useMock.value" color="amber" variant="subtle" size="sm">
        {{ $t('crmMockBadge') }}
      </UBadge>
      <AppUserMenu v-if="auth.user" align-right />
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
  if (route.path.startsWith('/contacts')) {
    return t('crmNavContacts')
  }
  if (route.path.startsWith('/deals')) {
    return t('dealsPageTitle')
  }
  if (route.path.startsWith('/settings')) {
    return t('settingsPageTitle')
  }
  return t('crmBreadcrumb')
})
</script>
