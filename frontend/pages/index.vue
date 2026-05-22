<template>
  <section class="rounded-lg bg-white p-6 shadow">
    <h2 class="mb-2 text-lg font-medium">{{ $t('welcome') }}</h2>
    <p class="text-sm text-gray-600">{{ $t('phase0Ready') }}</p>
    <p class="mt-4 text-xs text-gray-400">API: {{ config.public.apiBase }}</p>
    <p v-if="auth.user" class="mt-2 text-sm text-gray-600">
      {{ auth.user.email }}
      <span v-if="auth.isSuperAdmin" class="ml-2 rounded bg-amber-100 px-2 py-0.5 text-xs text-amber-800">
        Super Admin
      </span>
    </p>
  </section>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()
const auth = useAuth()
</script>
