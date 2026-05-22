<template>
  <div class="grid min-h-screen lg:grid-cols-2">
    <section
      class="ds-panel-brand relative hidden flex-col justify-between overflow-hidden px-10 py-12 lg:flex"
    >
      <div class="pointer-events-none absolute -right-16 top-20 h-64 w-64 rounded-full blur-3xl" :style="{ background: 'var(--ds-blur-brand)' }" />
      <div class="pointer-events-none absolute -left-10 bottom-10 h-48 w-48 rounded-full blur-2xl" :style="{ background: 'var(--ds-blur-accent)' }" />

      <div class="relative z-10">
        <div class="mb-6 flex items-center gap-3">
          <div class="ds-brand-gradient flex h-10 w-10 items-center justify-center rounded-xl ds-brand-shadow">
            <svg class="h-5 w-5 text-ds-on-brand" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <span class="text-lg font-semibold tracking-tight text-ds-fg-heading">{{ brandName }}</span>
        </div>

        <h1 class="text-4xl font-bold leading-[1.15] tracking-tight text-ds-fg-heading">
          {{ $t('loginHeroEmpower') }}
        </h1>
        <p class="mt-2 text-lg font-medium text-ds-fg-brand">
          {{ $t('loginHeroEmpowerSub') }}
        </p>
      </div>

      <div class="relative z-10 flex flex-1 items-center justify-center py-6">
        <LoginBrandIllustration v-if="!isDark" />
        <LoginBrandIllustrationDark v-else />
      </div>

      <div class="relative z-10 flex items-center justify-between text-xs text-ds-fg-subtle">
        <span>© {{ year }} EnterpriseFlow CRM</span>
        <UiThemeToggle />
      </div>
    </section>

    <section class="relative flex items-center justify-center bg-ds-bg-elevated px-6 py-10 lg:px-14">
      <div class="absolute inset-x-0 top-0 flex items-center justify-between border-b border-ds-border bg-ds-bg-muted px-6 py-4 lg:hidden">
        <div class="flex items-center gap-2">
          <div class="ds-brand-gradient flex h-8 w-8 items-center justify-center rounded-lg">
            <svg class="h-4 w-4 text-ds-on-brand" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <span class="font-semibold text-ds-fg-heading">{{ brandName }}</span>
        </div>
        <UiThemeToggle />
      </div>

      <div class="relative z-10 mt-14 w-full max-w-sm lg:mt-0">
        <div class="mb-8 flex items-center justify-between gap-2">
          <div class="flex items-center gap-2">
            <div class="ds-brand-gradient flex h-9 w-9 items-center justify-center rounded-lg">
              <svg class="h-5 w-5 text-ds-on-brand" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <span class="text-xl font-bold text-ds-fg-heading">{{ brandName }}</span>
          </div>
          <NuxtLink to="/design" class="text-xs text-ds-fg-subtle transition-colors hover:text-ds-fg-brand">
            {{ $t('themePreviewTitle') }}
          </NuxtLink>
        </div>

        <h2 class="text-3xl font-bold tracking-tight text-ds-fg-heading">{{ $t('loginTitle') }}</h2>
        <p class="mt-2 text-sm text-ds-fg-muted">{{ $t('loginFormHint') }}</p>

        <form class="mt-8 space-y-5" @submit.prevent="onSubmit">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="email">{{ $t('email') }}</label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              autocomplete="email"
              class="ds-input w-full rounded-xl px-4 py-3 text-sm transition-all"
              placeholder="admin@demo.com"
            >
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="password">{{ $t('password') }}</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              autocomplete="current-password"
              class="ds-input w-full rounded-xl px-4 py-3 text-sm transition-all"
            >
          </div>

          <p v-if="error" class="rounded-xl bg-ds-danger-subtle px-3 py-2 text-sm text-ds-danger" role="alert">{{ error }}</p>

          <button
            type="submit"
            class="ds-btn-primary w-full cursor-pointer rounded-xl px-4 py-3 text-sm font-semibold transition-all focus:outline-none focus:ring-2 focus:ring-ds-brand/40 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60"
            :class="isDark ? 'focus:ring-offset-ds-bg-elevated' : 'focus:ring-offset-white'"
            :disabled="loading"
          >
            {{ loading ? $t('loggingIn') : $t('login') }}
          </button>
        </form>

        <p class="mt-8 text-center text-xs text-ds-fg-muted">{{ $t('loginFooterNote') }}</p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { LoginResponse } from '~/composables/use-auth'

definePageMeta({ layout: 'auth' })

const { t } = useI18n()
const { isDark } = useTheme()
const auth = useAuth()
const tenant = useTenant()
const config = useRuntimeConfig()

const email = ref('admin@demo.com')
const password = ref('password123')
const loading = ref(false)
const error = ref('')
const year = new Date().getFullYear()

const brandName = computed(() => (isDark.value ? 'Optrixx' : 'EnterpriseFlow'))

onMounted(() => {
  if (auth.isAuthenticated.value) navigateTo('/')
})

async function onSubmit() {
  loading.value = true
  error.value = ''
  try {
    const res = await fetch(`${config.public.apiBase}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value, password: password.value }),
    })
    const body = await res.json()
    if (!res.ok) throw new Error(body.message || t('loginFailed'))
    const data = body.data as LoginResponse
    auth.setSession(data.access_token, data.refresh_token, data.user)
    tenant.setTenantList(data.tenants)
    if (data.current_tenant) tenant.setTenant(data.current_tenant.id)
    else if (data.tenants.length > 0) tenant.setTenant(data.tenants[0].id)
    if (tenant.currentTenantId.value) await useRbac().loadMyPermissions()
    await navigateTo(data.user.is_super_admin ? '/admin' : '/')
  } catch (e) {
    error.value = e instanceof Error ? e.message : t('loginFailed')
  } finally {
    loading.value = false
  }
}
</script>
