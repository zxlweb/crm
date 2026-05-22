<template>
  <section class="rounded-lg bg-white p-6 shadow">
    <h2 class="mb-2 text-lg font-medium">{{ $t('welcome') }}</h2>
    <p class="text-sm text-gray-600">{{ $t('phase0Ready') }}</p>
    <p class="mt-4 text-xs text-gray-400">API: {{ config.public.apiBase }}</p>
    <p v-if="auth.user?.email" class="mt-2 text-sm text-gray-600">
      {{ auth.user?.email }}
      <span v-if="auth.isSuperAdmin" class="ml-2 rounded bg-amber-100 px-2 py-0.5 text-xs text-amber-800">
        Super Admin
      </span>
    </p>
    <NuxtLink
      v-if="auth.isAuthenticated"
      to="/leads"
      class="mt-4 inline-flex rounded-lg bg-violet-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-violet-700"
    >
      {{ $t('leadsNavEntry') }}
    </NuxtLink>
  </section>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()
const auth = useAuth()
</script>
