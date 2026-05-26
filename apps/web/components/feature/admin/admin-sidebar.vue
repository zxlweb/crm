<template>
  <aside
    class="hidden w-60 shrink-0 flex-col border-r border-ds-border bg-ds-bg-sidebar lg:flex"
  >
    <NuxtLink
      to="/admin"
      class="flex h-16 items-center gap-2.5 border-b border-ds-border-muted px-5 transition-colors hover:bg-ds-bg-muted/40"
    >
      <div
        class="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-ds-brand to-ds-brand-strong shadow-ds-sm"
      >
        <UIcon
          name="i-heroicons-shield-check-20-solid"
          class="h-5 w-5 text-ds-on-brand"
          aria-hidden="true"
        />
      </div>
      <div class="min-w-0">
        <p class="text-sm font-bold text-ds-fg-heading leading-tight">CRM</p>
        <p class="truncate text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle">
          {{ $t('adminNavSubtitle') }}
        </p>
      </div>
    </NuxtLink>

    <nav class="flex flex-1 flex-col p-3" :aria-label="$t('adminBreadcrumb')">
      <div class="flex-1 space-y-0.5">
        <p
          class="px-2 pb-1.5 pt-1 text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle"
        >
          {{ $t('adminBreadcrumb') }}
        </p>

        <NuxtLink
          to="/admin"
          class="flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium transition-colors duration-200"
          :class="
            isActive('/admin')
              ? 'ds-nav-active'
              : 'text-ds-fg-nav hover:bg-ds-bg-muted hover:text-ds-fg-nav-active'
          "
        >
          <UIcon name="i-heroicons-squares-2x2-20-solid" class="h-5 w-5 shrink-0" aria-hidden="true" />
          <span class="truncate">{{ $t('adminNavDashboard') }}</span>
        </NuxtLink>

        <button
          type="button"
          data-testid="admin-nav-tenants"
          class="flex w-full cursor-pointer items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium text-ds-fg-nav transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
          @click="scrollToTenants"
        >
          <UIcon
            name="i-heroicons-building-office-2-20-solid"
            class="h-5 w-5 shrink-0"
            aria-hidden="true"
          />
          <span class="truncate">{{ $t('adminNavTenants') }}</span>
        </button>
      </div>

      <div class="space-y-0.5 border-t border-ds-border-muted pt-3 mt-3">
        <p
          class="px-2 pb-1.5 text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle"
        >
          {{ $t('shellDevToolsLabel') }}
        </p>
        <NuxtLink
          to="/cards"
          class="flex items-center gap-3 rounded-xl px-3 py-2 text-sm text-ds-fg-nav transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
        >
          <UIcon
            name="i-heroicons-square-3-stack-3d-20-solid"
            class="h-4 w-4 shrink-0"
            aria-hidden="true"
          />
          <span class="truncate">{{ $t('cardsPageTitle') }}</span>
        </NuxtLink>
        <NuxtLink
          to="/charts"
          class="flex items-center gap-3 rounded-xl px-3 py-2 text-sm text-ds-fg-nav transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
        >
          <UIcon name="i-heroicons-chart-bar-20-solid" class="h-4 w-4 shrink-0" aria-hidden="true" />
          <span class="truncate">{{ $t('chartsPageTitle') }}</span>
        </NuxtLink>
        <NuxtLink
          to="/design"
          class="flex items-center gap-3 rounded-xl px-3 py-2 text-sm text-ds-fg-nav transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
        >
          <UIcon name="i-heroicons-swatch-20-solid" class="h-4 w-4 shrink-0" aria-hidden="true" />
          <span class="truncate">{{ $t('themePreviewTitle') }}</span>
        </NuxtLink>
      </div>

      <div
        v-if="auth.isSuperAdmin.value"
        class="space-y-0.5 border-t border-ds-border-muted pt-3 mt-3"
      >
        <NuxtLink
          to="/"
          class="group flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium text-ds-fg-nav transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
        >
          <UIcon name="i-heroicons-briefcase-20-solid" class="h-5 w-5 shrink-0" aria-hidden="true" />
          <span class="truncate">{{ $t('userMenuOpenCrm') }}</span>
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
        {{ $t('adminNavSubtitle') }}
      </p>
      <p class="px-1.5 pt-1 text-[10px] text-ds-fg-subtle">v1.0 · Demo</p>
    </div>
  </aside>
</template>

<script setup lang="ts">
const route = useRoute()
const auth = useAuth()

function isActive(path: string) {
  return route.path === path
}

async function scrollToTenants() {
  if (route.path !== '/admin') {
    await navigateTo('/admin')
    await nextTick()
  }
  document.getElementById('tenants')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}
</script>
