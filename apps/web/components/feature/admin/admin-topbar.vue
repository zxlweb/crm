<template>
  <header
    class="sticky top-0 z-20 flex h-16 items-center justify-between gap-4 border-b border-ds-border bg-ds-bg-topbar px-6 backdrop-blur-md lg:px-8"
  >
    <div class="min-w-0 flex-1">
      <p class="text-[10px] font-medium uppercase tracking-wider text-ds-fg-subtle">
        {{ $t('adminBreadcrumb') }}
      </p>
      <h1 class="truncate text-lg font-semibold text-ds-fg-heading leading-tight">
        {{ $t('superAdminTitle') }}
      </h1>
    </div>

    <div class="flex items-center gap-2 sm:gap-3">
      <AppUserMenu v-if="auth.user" align-right />
    </div>
  </header>
</template>

<script setup lang="ts">
const auth = useAuth()
</script>
