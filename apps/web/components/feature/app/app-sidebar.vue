<template>
  <aside
    class="hidden w-60 shrink-0 flex-col border-r border-ds-border bg-ds-bg-sidebar lg:flex"
  >
    <NuxtLink
      to="/"
      class="flex h-16 items-center gap-2.5 border-b border-ds-border-muted px-5 transition-colors hover:bg-ds-bg-muted/40"
    >
      <div
        class="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-ds-brand to-ds-brand-strong shadow-ds-sm"
      >
        <svg
          class="h-5 w-5 text-ds-on-brand"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          aria-hidden="true"
        >
          <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
      </div>
      <div class="min-w-0">
        <p class="text-sm font-bold text-ds-fg-heading leading-tight">CRM</p>
        <p class="truncate text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle">
          {{ $t('crmNavSubtitle') }}
        </p>
      </div>
    </NuxtLink>

    <nav class="flex flex-1 flex-col p-3" :aria-label="$t('crmBreadcrumb')">
      <div class="flex-1 space-y-0.5">
        <p
          class="px-2 pb-1.5 pt-1 text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle"
        >
          {{ $t('crmBreadcrumb') }}
        </p>

        <NuxtLink
          v-for="item in primaryNav"
          :key="item.to"
          :to="item.to"
          :data-testid="item.testId"
          class="group flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium transition-colors duration-200"
          :class="
            item.active
              ? 'ds-nav-active'
              : 'text-ds-fg-nav hover:bg-ds-bg-muted hover:text-ds-fg-nav-active'
          "
        >
          <UIcon :name="item.icon" class="h-5 w-5 shrink-0" aria-hidden="true" />
          <span class="truncate">{{ item.label }}</span>
        </NuxtLink>
      </div>

      <div
        v-if="auth.isSuperAdmin.value"
        class="space-y-0.5 border-t border-ds-border-muted pt-3 mt-3"
      >
        <p
          class="px-2 pb-1.5 text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle"
        >
          {{ $t('userMenuSuperAdmin') }}
        </p>
        <NuxtLink
          to="/admin"
          class="group flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium text-ds-fg-nav transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
        >
          <UIcon name="i-heroicons-shield-check-20-solid" class="h-5 w-5 shrink-0" aria-hidden="true" />
          <span class="truncate">{{ $t('userMenuOpenAdmin') }}</span>
          <UIcon
            name="i-heroicons-arrow-up-right-20-solid"
            class="ml-auto h-3.5 w-3.5 shrink-0 text-ds-fg-subtle opacity-0 transition-opacity group-hover:opacity-100"
            aria-hidden="true"
          />
        </NuxtLink>
      </div>
    </nav>

    <div class="border-t border-ds-border-muted px-3 py-2.5">
      <p
        class="px-1.5 text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle leading-none"
      >
        {{ $t('crmNavSubtitle') }}
      </p>
      <p class="px-1.5 pt-1 text-[10px] text-ds-fg-subtle">v1.0 · Demo</p>
    </div>
  </aside>
</template>

<script setup lang="ts">
const route = useRoute()
const { t } = useI18n()
const auth = useAuth()
const permission = usePermission()

type NavItem = {
  to: string
  label: string
  icon: string
  active: boolean
  testId?: string
}

const primaryNav = computed<NavItem[]>(() => {
  const items: NavItem[] = [
    {
      to: '/',
      label: t('crmNavDashboard'),
      icon: 'i-heroicons-home-20-solid',
      active: route.path === '/',
    },
    {
      to: '/leads',
      label: t('crmNavLeads'),
      icon: 'i-heroicons-light-bulb-20-solid',
      active: route.path.startsWith('/leads'),
    },
    {
      to: '/accounts',
      label: t('crmNavAccounts'),
      icon: 'i-heroicons-building-office-2-20-solid',
      active: route.path.startsWith('/accounts'),
    },
    {
      to: '/contacts',
      label: t('crmNavContacts'),
      icon: 'i-heroicons-user-group-20-solid',
      active: route.path.startsWith('/contacts'),
    },
    {
      to: '/deals',
      label: t('crmNavDeals'),
      icon: 'i-heroicons-currency-dollar-20-solid',
      active: route.path.startsWith('/deals'),
    },
  ]

  // 设置入口：仅当用户有任意配置类权限才出现，避免普通销售看到「平台管理」类入口
  const canViewSettings =
    permission.can('settings', 'view') ||
    permission.can('settings', 'update') ||
    permission.can('rbac', 'view') ||
    permission.can('rbac', 'manage')

  if (canViewSettings) {
    items.push({
      to: '/settings',
      label: t('crmNavSettings'),
      icon: 'i-heroicons-cog-6-tooth-20-solid',
      active: route.path.startsWith('/settings'),
      testId: 'nav-settings',
    })
  }

  return items
})
</script>
