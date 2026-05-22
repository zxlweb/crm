<template>
  <aside class="hidden w-64 shrink-0 flex-col border-r border-ds-border bg-ds-bg-sidebar lg:flex">
    <div class="flex h-16 items-center gap-2.5 border-b border-ds-border-muted px-6">
      <div class="ds-brand-gradient flex h-9 w-9 items-center justify-center rounded-xl ds-brand-shadow">
        <svg class="h-5 w-5 text-ds-on-brand" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
      </div>
      <div>
        <p class="text-sm font-bold text-ds-fg-heading">CRM</p>
        <p class="text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle">{{ $t('crmNavSubtitle') }}</p>
      </div>
    </div>

    <nav class="flex-1 space-y-1 p-4">
      <NuxtLink
        to="/leads"
        class="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors duration-200"
        :class="isLeadsActive ? 'ds-nav-active' : 'text-ds-fg-nav hover:bg-ds-bg-muted hover:text-ds-fg-nav-active'"
      >
        <svg class="h-5 w-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
        </svg>
        {{ $t('crmNavLeads') }}
      </NuxtLink>

      <NuxtLink
        to="/accounts"
        class="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors duration-200"
        :class="isAccountsActive ? 'ds-nav-active' : 'text-ds-fg-nav hover:bg-ds-bg-muted hover:text-ds-fg-nav-active'"
      >
        <svg class="h-5 w-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
        </svg>
        {{ $t('crmNavAccounts') }}
      </NuxtLink>

      <span
        class="flex cursor-not-allowed items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-ds-fg-subtle opacity-60"
        :title="$t('crmNavComingSoon')"
      >
        <svg class="h-5 w-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        {{ $t('crmNavContacts') }}
      </span>
    </nav>

    <div class="space-y-2 border-t border-ds-border-muted p-4">
      <div class="px-1">
        <UiThemeToggle />
      </div>
      <NuxtLink
        to="/"
        class="flex items-center gap-2 rounded-xl px-3 py-2 text-sm text-ds-fg-nav transition-colors hover:bg-ds-bg-muted hover:text-ds-fg-nav-active"
      >
        <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        {{ $t('backToApp') }}
      </NuxtLink>
    </div>
  </aside>
</template>

<script setup lang="ts">
const route = useRoute()

const isLeadsActive = computed(() => route.path.startsWith('/leads'))
const isAccountsActive = computed(() => route.path.startsWith('/accounts'))
</script>
