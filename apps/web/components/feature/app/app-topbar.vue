<template>
  <header
    class="sticky top-0 z-10 flex h-16 shrink-0 items-center justify-between gap-4 border-b border-ds-border bg-ds-bg-topbar px-6 backdrop-blur-md lg:px-8"
  >
    <div class="min-w-0 flex-1">
      <h1 class="truncate text-lg font-semibold leading-tight text-ds-fg-heading">
        <template v-if="breadcrumbPrefix">
          <span class="font-medium text-ds-fg-muted">{{ breadcrumbPrefix }}</span>
          <span class="mx-1.5 font-normal text-ds-fg-subtle" aria-hidden="true">/</span>
        </template>
        <span>{{ pageTitle }}</span>
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

/** 子页面在顶栏展示「销售工作台 / 当前页」，页面内不再重复大标题 */
const breadcrumbPrefix = computed(() => {
  if (route.path === '/') return null
  return t('crmBreadcrumb')
})

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
    return t('leadsTitle')
  }
  if (route.path.startsWith('/contacts')) {
    return route.path === '/contacts' ? t('contactsPageTitle') : t('crmNavContacts')
  }
  if (route.path.startsWith('/deals')) {
    return t('dealsPageTitle')
  }
  if (route.path.startsWith('/settings')) {
    return t('settingsPageTitle')
  }
  if (route.path.startsWith('/charts')) {
    return t('chartsPageTitle')
  }
  if (route.path.startsWith('/admin')) {
    return t('adminBreadcrumb')
  }
  return t('crmBreadcrumb')
})
</script>
