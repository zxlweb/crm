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

        <Transition name="auth-hero" mode="out-in">
          <div :key="mode">
            <h1 class="text-4xl font-bold leading-[1.15] tracking-tight text-ds-fg-heading">
              {{ heroEmpower }}
            </h1>
            <p class="mt-2 text-lg font-medium text-ds-fg-brand">
              {{ heroEmpowerSub }}
            </p>
          </div>
        </Transition>
      </div>

      <div class="relative z-10 flex flex-1 items-center justify-center py-6">
        <LoginBrandIllustration v-if="!isDark" />
        <LoginBrandIllustrationDark v-else />
      </div>

      <div class="relative z-10 flex items-center justify-between text-xs text-ds-fg-subtle">
        <span>© {{ year }} EnterpriseFlow CRM</span>
        <UiThemeToggle variant="icon" :format-label="skinLabel" :aria-label="$t('themeSkinLabel')" />
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
        <UiThemeToggle variant="icon" :format-label="skinLabel" :aria-label="$t('themeSkinLabel')" />
      </div>

      <div
        class="auth-panel relative z-10 mt-14 w-full max-w-sm lg:mt-0"
        :class="{ 'auth-panel--visible': panelVisible }"
      >
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

        <div class="relative mt-8 min-h-[22rem]">
          <Transition name="auth-fade" mode="out-in">
            <div :key="mode" class="auth-scene">
              <h2 class="text-3xl font-bold tracking-tight text-ds-fg-heading" data-testid="auth-form-title">
                {{ isRegister ? $t('registerTitle') : $t('loginTitle') }}
              </h2>
              <p class="mt-2 text-sm text-ds-fg-muted">
                {{ isRegister ? $t('registerSubtitle') : $t('loginFormHint') }}
              </p>

              <form class="mt-8 space-y-5" data-testid="auth-form" @submit.prevent="onSubmit">
                <template v-if="isRegister">
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="name">{{ $t('name') }}</label>
                    <input
                      id="name"
                      v-model="name"
                      type="text"
                      required
                      autocomplete="name"
                      data-testid="auth-input-name"
                      class="ds-input ds-transition w-full rounded-xl px-4 py-3 text-sm motion-reduce:transition-none"
                    >
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="company">{{ $t('companyName') }}</label>
                    <input
                      id="company"
                      v-model="companyName"
                      type="text"
                      required
                      autocomplete="organization"
                      data-testid="auth-input-company"
                      class="ds-input ds-transition w-full rounded-xl px-4 py-3 text-sm motion-reduce:transition-none"
                    >
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="domain">{{ $t('tenantDomain') }}</label>
                    <input
                      id="domain"
                      v-model="domain"
                      type="text"
                      data-testid="auth-input-domain"
                      class="ds-input ds-transition w-full rounded-xl px-4 py-3 text-sm motion-reduce:transition-none"
                      :placeholder="$t('tenantDomainHint')"
                    >
                  </div>
                </template>

                <div>
                  <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="email">{{ $t('email') }}</label>
                  <input
                    id="email"
                    v-model="email"
                    type="email"
                    required
                    autocomplete="email"
                    data-testid="auth-input-email"
                    class="ds-input ds-transition w-full rounded-xl px-4 py-3 text-sm motion-reduce:transition-none"
                    :placeholder="isRegister ? undefined : 'admin@demo.com'"
                  >
                </div>
                <div>
                  <label class="mb-1.5 block text-sm font-medium text-ds-fg" for="password">{{ $t('password') }}</label>
                  <input
                    id="password"
                    v-model="password"
                    type="password"
                    required
                    minlength="6"
                    data-testid="auth-input-password"
                    :autocomplete="isRegister ? 'new-password' : 'current-password'"
                    class="ds-input ds-transition w-full rounded-xl px-4 py-3 text-sm motion-reduce:transition-none"
                  >
                </div>

                <p v-if="error" class="rounded-xl bg-ds-danger-subtle px-3 py-2 text-sm text-ds-danger" role="alert">{{ error }}</p>

                <button
                  type="submit"
                  data-testid="auth-submit"
                  class="ds-btn-primary w-full cursor-pointer rounded-xl px-4 py-3 text-sm font-semibold transition-all focus:outline-none focus:ring-2 focus:ring-ds-brand/40 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60 motion-reduce:transition-none"
                  :class="isDark ? 'focus:ring-offset-ds-bg-elevated' : 'focus:ring-offset-white'"
                  :disabled="loading"
                >
                  {{ submitLabel }}
                </button>
              </form>

              <p class="mt-6 text-center text-sm text-ds-fg-muted">
                <template v-if="isRegister">
                  {{ $t('hasAccount') }}
                  <button
                    type="button"
                    class="font-medium text-ds-fg-brand hover:underline"
                    data-testid="auth-switch-to-login"
                    @click="setMode('login')"
                  >
                    {{ $t('login') }}
                  </button>
                </template>
                <template v-else>
                  {{ $t('noAccount') }}
                  <button
                    type="button"
                    class="font-medium text-ds-fg-brand hover:underline"
                    data-testid="auth-switch-to-register"
                    @click="setMode('register')"
                  >
                    {{ $t('register') }}
                  </button>
                </template>
              </p>

              <p v-if="!isRegister" class="mt-3 text-center text-xs text-ds-fg-muted">
                {{ $t('loginFooterNote') }}
              </p>
            </div>
          </Transition>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { registerSchema } from '~/schemas/register'
import type { LoginResponse } from '~/composables/use-auth'
import type { ThemeId } from '@crm/ui-kit'

definePageMeta({ layout: 'auth' })

type AuthMode = 'login' | 'register'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { isDark } = useTheme()
const auth = useAuth()
const tenant = useTenant()
const config = useRuntimeConfig()

const mode = ref<AuthMode>('login')
const panelVisible = ref(false)
const email = ref('admin@demo.com')
const password = ref('password123')
const name = ref('')
const companyName = ref('')
const domain = ref('')
const loading = ref(false)
const error = ref('')
const year = new Date().getFullYear()

const isRegister = computed(() => mode.value === 'register')
const brandName = computed(() => (isDark.value ? 'Optrixx' : 'EnterpriseFlow'))

function skinLabel(id: ThemeId) {
  return id === 'v1' ? t('themeSkinLight') : t('themeSkinDark')
}
const heroEmpower = computed(() => (isRegister.value ? t('registerHeroEmpower') : t('loginHeroEmpower')))
const heroEmpowerSub = computed(() => (isRegister.value ? t('registerHeroEmpowerSub') : t('loginHeroEmpowerSub')))
const submitLabel = computed(() => {
  if (loading.value) return isRegister.value ? t('registering') : t('loggingIn')
  return isRegister.value ? t('register') : t('login')
})

function setMode(next: AuthMode) {
  if (mode.value === next) return
  mode.value = next
  error.value = ''
  if (next === 'login') {
    name.value = ''
    companyName.value = ''
    domain.value = ''
    router.replace({ query: {} })
  } else {
    if (email.value === 'admin@demo.com') email.value = ''
    if (password.value === 'password123') password.value = ''
    router.replace({ query: { mode: 'register' } })
  }
  nextTick(() => {
    const focusId = next === 'register' ? 'name' : 'email'
    document.getElementById(focusId)?.focus()
  })
}

onMounted(async () => {
  if (auth.isAuthenticated.value) {
    const tenant = useTenant()
    try {
      const list = await tenant.fetchTenants()
      const tid = tenant.currentTenantId.value
      const valid = !!tid && list.some((t) => t.id === tid)
      if (list.length > 0 && !auth.isSuperAdmin.value) {
        await tenant.switchTenant(valid ? tid! : list[0].id)
      }
    } catch {
      // 仍进入首页，由 tenant 插件尝试修复
    }
    await navigateTo(auth.isSuperAdmin.value ? '/admin' : '/')
    return
  }
  if (route.query.mode === 'register') {
    mode.value = 'register'
    email.value = ''
    password.value = ''
  }
  requestAnimationFrame(() => {
    panelVisible.value = true
  })
})

async function finishSession(data: LoginResponse) {
  auth.setSession(data.access_token, data.refresh_token, data.user)
  tenant.setTenantList(data.tenants)

  const targetTenantId = data.current_tenant?.id ?? data.tenants[0]?.id
  if (targetTenantId && !data.user.is_super_admin) {
    // 必须 switch-tenant：把租户/角色写入 JWT，并清掉残留的 demo 租户 cookie（否则 Casbin 403、全站无数据）
    try {
      await tenant.switchTenant(targetTenantId)
    } catch {
      tenant.setTenant(targetTenantId)
      tenant.applyDepartmentFromLogin(data)
      useActiveRole().applyFromLogin(data)
      try {
        await useRbac().loadMyPermissions()
      } catch {
        // 登录已成功；权限稍后由 rbac 插件补拉
      }
    }
  } else {
    if (targetTenantId) tenant.setTenant(targetTenantId)
    tenant.applyDepartmentFromLogin(data)
    useActiveRole().applyFromLogin(data)
    if (targetTenantId && !data.user.is_super_admin) {
      try {
        await useRbac().loadMyPermissions()
      } catch {
        // ignore
      }
    }
  }

  await navigateTo(data.user.is_super_admin ? '/admin' : '/')
}

async function onSubmit() {
  loading.value = true
  error.value = ''
  try {
    if (isRegister.value) {
      const payload = {
        email: email.value,
        password: password.value,
        name: name.value,
        company_name: companyName.value,
        domain: domain.value.trim() || undefined,
      }
      const parsed = registerSchema.safeParse(payload)
      if (!parsed.success) {
        error.value = t('registerValidationFailed')
        loading.value = false
        return
      }
      const res = await fetch(`${config.public.apiBase}/api/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(parsed.data),
      })
      const body = await res.json()
      if (!res.ok) throw new Error(body.message || t('registerFailed'))
      await finishSession(body.data as LoginResponse)
    } else {
      const res = await fetch(`${config.public.apiBase}/api/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: email.value, password: password.value }),
      })
      const body = await res.json()
      if (!res.ok) throw new Error(body.message || t('loginFailed'))
      await finishSession(body.data as LoginResponse)
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : (isRegister.value ? t('registerFailed') : t('loginFailed'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-panel {
  opacity: 0;
  transform: translateX(2.5rem);
}

.auth-panel--visible {
  animation: auth-panel-in 0.7s cubic-bezier(0.22, 1, 0.36, 1) forwards;
}

.auth-scene {
  width: 100%;
}

/* 场景切换：从右往左滑入 */
.auth-fade-enter-active {
  transition:
    opacity 0.48s cubic-bezier(0.22, 1, 0.36, 1),
    transform 0.48s cubic-bezier(0.22, 1, 0.36, 1);
}

.auth-fade-leave-active {
  transition:
    opacity 0.26s cubic-bezier(0.4, 0, 0.6, 1),
    transform 0.26s cubic-bezier(0.4, 0, 0.6, 1);
  position: absolute;
  inset-inline: 0;
  top: 0;
}

.auth-fade-enter-from {
  opacity: 0;
  transform: translateX(2rem);
}

.auth-fade-leave-to {
  opacity: 0;
  transform: translateX(-1.25rem);
}

/* 左侧标语 */
.auth-hero-enter-active {
  transition:
    opacity 0.48s cubic-bezier(0.22, 1, 0.36, 1),
    transform 0.48s cubic-bezier(0.22, 1, 0.36, 1);
}

.auth-hero-leave-active {
  transition:
    opacity 0.28s cubic-bezier(0.4, 0, 0.6, 1),
    transform 0.28s cubic-bezier(0.4, 0, 0.6, 1);
}

.auth-hero-enter-from,
.auth-hero-leave-to {
  opacity: 0;
  transform: translateY(0.625rem);
}

@keyframes auth-panel-in {
  from {
    opacity: 0;
    transform: translateX(2.5rem);
  }

  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .auth-panel {
    opacity: 1;
    transform: none;
  }

  .auth-panel--visible {
    animation: none;
  }

  .auth-fade-enter-active,
  .auth-fade-leave-active,
  .auth-hero-enter-active,
  .auth-hero-leave-active {
    transition-duration: 0.15s;
  }

  .auth-fade-enter-from,
  .auth-fade-leave-to,
  .auth-hero-enter-from,
  .auth-hero-leave-to {
    transform: none;
  }
}
</style>
