<template>
  <aside class="hidden w-64 shrink-0 flex-col border-r border-ds-border bg-ds-bg-sidebar lg:flex">
    <div class="flex h-16 items-center gap-2.5 border-b border-ds-border-muted px-6">
      <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-ds-brand shadow-ds-sm">
        <svg class="h-5 w-5 text-ds-on-brand" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
      </div>
      <div>
        <p class="text-sm font-bold text-ds-fg-heading">CRM</p>
        <p class="text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle">{{ $t('adminNavSubtitle') }}</p>
      </div>
    </div>

    <nav class="flex flex-1 flex-col p-4">
      <AppShellContextSwitch class="mb-4" />

      <div class="flex-1 space-y-1">
        <NuxtLink
          to="/admin"
          class="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors duration-200"
          :class="isActive('/admin') ? 'ds-nav-active' : 'text-ds-fg-nav hover:bg-ds-bg-muted hover:text-ds-fg-nav-active'"
        >
          <svg class="h-5 w-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
          {{ $t('adminNavDashboard') }}
        </NuxtLink>

        <button
          type="button"
          data-testid="admin-nav-tenants"
          class="flex w-full cursor-pointer items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-ds-fg-nav transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
          @click="scrollToTenants"
        >
          <svg class="h-5 w-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
          </svg>
          {{ $t('adminNavTenants') }}
        </button>
      </div>
    </nav>

    <div class="space-y-1 border-t border-ds-border-muted p-4">
      <p class="px-0.5 pb-1 text-[10px] font-medium uppercase tracking-wide text-ds-fg-subtle">
        {{ $t('shellDevToolsLabel') }}
      </p>
      <NuxtLink
        to="/cards"
        class="flex items-center gap-2 rounded-xl px-3 py-2 text-sm text-ds-fg-nav transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
      >
        <svg class="h-4 w-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zm10 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zm10 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
        </svg>
        {{ $t('cardsPageTitle') }}
      </NuxtLink>
      <NuxtLink
        to="/charts"
        class="flex items-center gap-2 rounded-xl px-3 py-2 text-sm text-ds-fg-nav transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
      >
        <svg class="h-4 w-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        {{ $t('chartsPageTitle') }}
      </NuxtLink>
      <NuxtLink
        to="/design"
        class="flex items-center gap-2 rounded-xl px-3 py-2 text-sm text-ds-fg-nav transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
      >
        <svg class="h-4 w-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01" />
        </svg>
        {{ $t('themePreviewTitle') }}
      </NuxtLink>
    </div>
  </aside>
</template>

<script setup lang="ts">
const route = useRoute()

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
