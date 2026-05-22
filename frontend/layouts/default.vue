<template>
  <div class="min-h-screen bg-gray-50 font-sans text-gray-900">
    <header class="flex items-center justify-between border-b bg-white px-6 py-4 shadow-sm">
      <NuxtLink to="/" class="text-xl font-semibold tracking-tight">CRM</NuxtLink>
      <nav class="flex items-center gap-4 text-sm">
        <NuxtLink
          v-if="auth.isSuperAdmin"
          to="/admin"
          class="text-amber-700 transition-colors hover:text-amber-800"
        >
          {{ $t('superAdminEntry') }}
        </NuxtLink>
        <PermissionGuard resource="rbac" action="view">
          <NuxtLink to="/settings/roles" class="text-gray-600 transition-colors hover:text-gray-900">
            {{ $t('roleManagement') }}
          </NuxtLink>
          <template #fallback />
        </PermissionGuard>
        <NuxtLink
          v-if="!auth.isAuthenticated"
          to="/login"
          class="text-blue-600 transition-colors hover:text-blue-700"
        >
          {{ $t('login') }}
        </NuxtLink>
        <button
          v-else
          type="button"
          class="cursor-pointer text-gray-500 transition-colors hover:text-gray-800"
          @click="logout"
        >
          {{ $t('logout') }}
        </button>
      </nav>
    </header>
    <main class="p-6">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
const auth = useAuth()
const tenant = useTenant()

function logout() {
  auth.clearTokens()
  tenant.setTenant('')
  tenant.setTenantList([])
  navigateTo('/login')
}
</script>
