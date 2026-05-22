<template>
  <header class="sticky top-0 z-20 flex h-16 items-center justify-between gap-4 border-b border-ds-border bg-ds-bg-topbar px-6 backdrop-blur-md lg:px-8">
    <div class="min-w-0 flex-1">
      <p class="text-xs font-medium text-ds-fg-brand">{{ $t('adminBreadcrumb') }}</p>
      <h1 class="truncate text-lg font-semibold text-ds-fg-heading">{{ $t('superAdminTitle') }}</h1>
    </div>

    <div class="flex items-center gap-3">
      <UiThemeToggle class="hidden sm:inline-flex" />

      <div
        v-if="auth.user"
        class="flex items-center gap-2 rounded-xl border border-ds-border bg-ds-bg-elevated py-1.5 pl-1.5 pr-3 shadow-sm"
      >
        <div class="ds-brand-gradient flex h-8 w-8 items-center justify-center rounded-lg text-xs font-bold text-ds-on-brand">
          {{ initials }}
        </div>
        <div class="hidden text-left sm:block">
          <p class="max-w-[120px] truncate text-xs font-semibold text-ds-fg-heading">{{ auth.user.email }}</p>
          <p class="text-[10px] text-ds-fg-subtle">Super Admin</p>
        </div>
      </div>

      <button
        type="button"
        class="cursor-pointer rounded-xl border border-ds-border px-3 py-2 text-sm text-ds-fg-muted transition-colors hover:border-ds-brand hover:bg-ds-bg-muted hover:text-ds-fg-brand"
        @click="logout"
      >
        {{ $t('logout') }}
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
const auth = useAuth()
const tenant = useTenant()

const initials = computed(() => {
  const email = auth.user.value?.email ?? 'A'
  return email.charAt(0).toUpperCase()
})

function logout() {
  auth.clearTokens()
  tenant.setTenant('')
  tenant.setTenantList([])
  navigateTo('/login')
}
</script>
