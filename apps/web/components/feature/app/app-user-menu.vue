<template>
  <div ref="rootRef" class="relative" data-testid="user-menu">
    <button
      type="button"
      class="group flex cursor-pointer items-center gap-2 rounded-xl border border-ds-border bg-ds-bg-elevated py-1.5 pl-1.5 pr-2 shadow-sm transition-all duration-200 hover:border-ds-brand/60 hover:bg-ds-bg-muted hover:shadow"
      :class="[
        compact ? 'pr-1.5' : 'sm:pr-2.5',
        block ? 'w-full' : '',
      ]"
      aria-haspopup="menu"
      :aria-expanded="open"
      :aria-label="$t('userMenuTriggerLabel')"
      data-testid="user-menu-trigger"
      @click.stop="toggle"
    >
      <div
        class="relative flex h-8 w-8 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-gradient-to-br from-ds-brand to-ds-brand-strong text-xs font-bold text-ds-on-brand shadow-inner"
        aria-hidden="true"
      >
        {{ initials }}
      </div>
      <div v-if="!compact" class="hidden min-w-0 flex-1 text-left sm:block">
        <p class="max-w-[160px] truncate text-xs font-semibold text-ds-fg-heading leading-tight">
          {{ displayName }}
        </p>
        <p
          v-if="contextLine"
          class="max-w-[160px] truncate text-[10px] text-ds-fg-subtle leading-tight mt-0.5"
        >
          {{ contextLine }}
        </p>
      </div>
      <UIcon
        name="i-heroicons-chevron-down-20-solid"
        class="hidden h-3.5 w-3.5 shrink-0 text-ds-fg-subtle transition-transform duration-200 sm:block"
        :class="open ? 'rotate-180' : ''"
        aria-hidden="true"
      />
    </button>

    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0 scale-95 -translate-y-1"
      enter-to-class="opacity-100 scale-100 translate-y-0"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100 scale-100 translate-y-0"
      leave-to-class="opacity-0 scale-95 -translate-y-1"
    >
      <div
        v-if="open"
        role="menu"
        class="absolute z-50 mt-2 w-[296px] origin-top-right overflow-hidden rounded-2xl border border-ds-border bg-ds-bg-elevated shadow-xl ring-1 ring-black/5"
        :class="alignRight ? 'right-0' : 'left-0'"
        data-testid="user-menu-panel"
        @click.stop
      >
        <!-- Identity header -->
        <div class="relative overflow-hidden border-b border-ds-border-muted px-4 pt-4 pb-3">
          <div
            class="absolute inset-0 -z-0 opacity-70"
            aria-hidden="true"
            :style="{
              background:
                'linear-gradient(135deg, color-mix(in srgb, var(--ds-brand) 12%, transparent), transparent 70%)',
            }"
          />
          <div class="relative flex items-start gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-gradient-to-br from-ds-brand to-ds-brand-strong text-sm font-bold text-ds-on-brand shadow-sm"
              aria-hidden="true"
            >
              {{ initials }}
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold text-ds-fg-heading">{{ displayName }}</p>
              <p v-if="email" class="truncate text-xs text-ds-fg-muted">
                {{ email }}
              </p>
              <div class="mt-1.5 flex flex-wrap items-center gap-1">
                <span
                  v-if="workspaceLabel"
                  class="inline-flex items-center gap-1 rounded-full bg-ds-brand/10 px-2 py-0.5 text-[10px] font-medium text-ds-brand"
                >
                  <UIcon :name="workspaceIcon" class="h-3 w-3" aria-hidden="true" />
                  {{ workspaceLabel }}
                </span>
                <span
                  v-if="roleLabel"
                  class="inline-flex items-center rounded-full border border-ds-border bg-ds-bg-elevated px-2 py-0.5 text-[10px] font-medium text-ds-fg-muted"
                >
                  {{ roleLabel }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="max-h-[60vh] overflow-y-auto">
          <!-- Workspace quick switch (super admin only) -->
          <section
            v-if="canSwitchWorkspace"
            class="border-b border-ds-border-muted px-3 py-3"
            data-testid="user-menu-workspace"
          >
            <p class="px-1 pb-2 text-[10px] font-medium uppercase tracking-wide text-ds-fg-subtle">
              {{ $t('userMenuSectionWorkspace') }}
            </p>
            <fieldset
              class="m-0 inline-flex w-full items-center gap-1 rounded-xl border border-ds-border bg-ds-bg p-1"
            >
              <legend class="sr-only">{{ $t('userMenuSectionWorkspace') }}</legend>
              <button
                type="button"
                class="flex flex-1 cursor-pointer items-center justify-center gap-1.5 rounded-lg px-2 py-1.5 text-xs font-medium transition-colors duration-200"
                :class="
                  isCrmContext
                    ? 'bg-ds-brand text-ds-on-brand shadow-sm'
                    : 'text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg'
                "
                :aria-pressed="isCrmContext"
                data-testid="user-menu-workspace-crm"
                @click="goWorkspace('/')"
              >
                <UIcon name="i-heroicons-briefcase-20-solid" class="h-3.5 w-3.5" aria-hidden="true" />
                {{ $t('appContextCrm') }}
              </button>
              <button
                type="button"
                class="flex flex-1 cursor-pointer items-center justify-center gap-1.5 rounded-lg px-2 py-1.5 text-xs font-medium transition-colors duration-200"
                :class="
                  isAdminContext
                    ? 'bg-ds-brand text-ds-on-brand shadow-sm'
                    : 'text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg'
                "
                :aria-pressed="isAdminContext"
                data-testid="user-menu-workspace-admin"
                @click="goWorkspace('/admin')"
              >
                <UIcon name="i-heroicons-shield-check-20-solid" class="h-3.5 w-3.5" aria-hidden="true" />
                {{ $t('appContextAdmin') }}
              </button>
            </fieldset>
          </section>

          <!-- Role switcher -->
          <section
            v-if="showRoleSection"
            class="border-b border-ds-border-muted px-3 py-3"
            data-testid="user-menu-roles"
          >
            <p class="px-1 pb-1 text-[10px] font-medium uppercase tracking-wide text-ds-fg-subtle">
              {{ $t('userMenuSectionRoles') }}
            </p>
            <ul class="space-y-0.5">
              <li v-for="r in activeRole.roles.value" :key="r.id">
                <button
                  type="button"
                  role="menuitemradio"
                  :aria-checked="r.id === activeRole.currentRoleId.value"
                  :disabled="rolePending"
                  class="flex w-full cursor-pointer items-center gap-2 rounded-lg px-2 py-1.5 text-left text-xs transition-colors hover:bg-ds-bg-muted disabled:cursor-not-allowed disabled:opacity-60"
                  :class="
                    r.id === activeRole.currentRoleId.value
                      ? 'bg-ds-brand/10 text-ds-fg-heading'
                      : 'text-ds-fg'
                  "
                  data-testid="user-menu-role-item"
                  @click="onSelectRole(r.id)"
                >
                  <span
                    class="flex h-5 w-5 shrink-0 items-center justify-center rounded-md"
                    :class="
                      r.id === activeRole.currentRoleId.value
                        ? 'bg-ds-brand text-ds-on-brand'
                        : 'bg-ds-bg text-ds-fg-subtle'
                    "
                  >
                    <UIcon
                      v-if="r.id === activeRole.currentRoleId.value"
                      name="i-heroicons-check-20-solid"
                      class="h-3 w-3"
                      aria-hidden="true"
                    />
                    <UIcon
                      v-else
                      name="i-heroicons-user-circle-20-solid"
                      class="h-3 w-3"
                      aria-hidden="true"
                    />
                  </span>
                  <span class="min-w-0 flex-1 truncate font-medium">{{ r.name }}</span>
                  <span
                    v-if="r.id === activeRole.currentRoleId.value"
                    class="text-[10px] font-medium text-ds-brand"
                  >
                    {{ $t('userMenuActive') }}
                  </span>
                </button>
              </li>
            </ul>
            <p v-if="roleError" class="mt-1.5 truncate px-1 text-[10px] text-ds-danger" role="alert">
              {{ roleError }}
            </p>
          </section>

          <!-- Tenant switcher -->
          <section
            v-if="showTenantSection"
            class="border-b border-ds-border-muted px-3 py-3"
            data-testid="user-menu-tenants"
          >
            <p class="px-1 pb-1 text-[10px] font-medium uppercase tracking-wide text-ds-fg-subtle">
              {{ $t('userMenuSectionTenants') }}
            </p>
            <ul class="space-y-0.5 max-h-48 overflow-y-auto pr-1">
              <li v-for="ti in tenant.tenants.value" :key="ti.id">
                <button
                  type="button"
                  role="menuitemradio"
                  :aria-checked="ti.id === tenant.currentTenantId.value"
                  :disabled="tenantPending"
                  class="flex w-full cursor-pointer items-center gap-2 rounded-lg px-2 py-1.5 text-left text-xs transition-colors hover:bg-ds-bg-muted disabled:cursor-not-allowed disabled:opacity-60"
                  :class="
                    ti.id === tenant.currentTenantId.value
                      ? 'bg-ds-brand/10 text-ds-fg-heading'
                      : 'text-ds-fg'
                  "
                  data-testid="user-menu-tenant-item"
                  @click="onSelectTenant(ti.id)"
                >
                  <span
                    class="flex h-5 w-5 shrink-0 items-center justify-center rounded-md"
                    :class="
                      ti.id === tenant.currentTenantId.value
                        ? 'bg-ds-brand text-ds-on-brand'
                        : 'bg-ds-bg text-ds-fg-subtle'
                    "
                  >
                    <UIcon
                      v-if="ti.id === tenant.currentTenantId.value"
                      name="i-heroicons-check-20-solid"
                      class="h-3 w-3"
                      aria-hidden="true"
                    />
                    <UIcon
                      v-else
                      name="i-heroicons-building-office-2-20-solid"
                      class="h-3 w-3"
                      aria-hidden="true"
                    />
                  </span>
                  <span class="min-w-0 flex-1 truncate font-medium">{{ ti.name }}</span>
                  <span
                    v-if="ti.id === tenant.currentTenantId.value"
                    class="text-[10px] font-medium text-ds-brand"
                  >
                    {{ $t('userMenuActive') }}
                  </span>
                </button>
              </li>
            </ul>
            <p v-if="tenantError" class="mt-1.5 truncate px-1 text-[10px] text-ds-danger" role="alert">
              {{ tenantError }}
            </p>
          </section>

          <!-- Appearance -->
          <section class="border-b border-ds-border-muted px-3 py-3" data-testid="user-menu-appearance">
            <p class="px-1 pb-2 text-[10px] font-medium uppercase tracking-wide text-ds-fg-subtle">
              {{ $t('userMenuSectionAppearance') }}
            </p>
            <fieldset
              class="m-0 inline-flex w-full items-center gap-1 rounded-xl border border-ds-border bg-ds-bg p-1"
            >
              <legend class="sr-only">{{ $t('themeSkinLabel') }}</legend>
              <button
                type="button"
                class="flex flex-1 cursor-pointer items-center justify-center gap-1.5 rounded-lg px-2 py-1.5 text-xs font-medium transition-colors duration-200"
                :class="
                  themeId === 'v1'
                    ? 'bg-ds-brand text-ds-on-brand shadow-sm'
                    : 'text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg'
                "
                :aria-pressed="themeId === 'v1'"
                data-testid="user-menu-theme-light"
                @click="setTheme('v1')"
              >
                <UIcon name="i-heroicons-sun-20-solid" class="h-3.5 w-3.5" aria-hidden="true" />
                {{ $t('themeSkinLight') }}
              </button>
              <button
                type="button"
                class="flex flex-1 cursor-pointer items-center justify-center gap-1.5 rounded-lg px-2 py-1.5 text-xs font-medium transition-colors duration-200"
                :class="
                  themeId === 'v2'
                    ? 'bg-ds-brand text-ds-on-brand shadow-sm'
                    : 'text-ds-fg-muted hover:bg-ds-bg-muted hover:text-ds-fg'
                "
                :aria-pressed="themeId === 'v2'"
                data-testid="user-menu-theme-dark"
                @click="setTheme('v2')"
              >
                <UIcon name="i-heroicons-moon-20-solid" class="h-3.5 w-3.5" aria-hidden="true" />
                {{ $t('themeSkinDark') }}
              </button>
            </fieldset>
          </section>
        </div>

        <!-- Account / logout -->
        <div class="px-1.5 py-1.5">
          <button
            type="button"
            role="menuitem"
            class="flex w-full cursor-pointer items-center gap-2 rounded-lg px-3 py-2 text-left text-sm font-medium text-ds-danger transition-colors hover:bg-ds-danger-subtle"
            data-testid="user-menu-logout"
            @click="handleLogout"
          >
            <UIcon name="i-heroicons-arrow-right-on-rectangle-20-solid" class="h-4 w-4" aria-hidden="true" />
            {{ $t('logout') }}
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { useUiKitTheme } from '@crm/ui-kit'

type ThemeId = 'v1' | 'v2'

withDefaults(
  defineProps<{
    /** 顶栏窄屏仅显示头像 */
    compact?: boolean
    alignRight?: boolean
    /** 侧栏底栏：按钮占满宽 */
    block?: boolean
  }>(),
  { compact: false, alignRight: true, block: false },
)

const { t } = useI18n()
const auth = useAuth()
const tenant = useTenant()
const activeRole = useActiveRole()
const route = useRoute()
const { logout } = useLogout()
const { id: themeId, setTheme: setUiKitTheme } = useUiKitTheme()

const open = ref(false)
const rootRef = ref<HTMLElement | null>(null)
const rolePending = ref(false)
const roleError = ref('')
const tenantPending = ref(false)
const tenantError = ref('')

const initials = computed(() => {
  const name = auth.user.value?.name?.trim()
  if (name) return name.charAt(0).toUpperCase()
  const email = auth.user.value?.email ?? '?'
  return email.charAt(0).toUpperCase()
})

const displayName = computed(() => auth.user.value?.name || auth.user.value?.email || '—')
const email = computed(() => auth.user.value?.email ?? '')

const isAdminContext = computed(() => route.path.startsWith('/admin'))
const isCrmContext = computed(() => !isAdminContext.value)

const canSwitchWorkspace = computed(() => auth.isSuperAdmin.value)

const workspaceLabel = computed(() => {
  if (!canSwitchWorkspace.value) return ''
  return isAdminContext.value ? t('appContextAdmin') : t('appContextCrm')
})

const workspaceIcon = computed(() =>
  isAdminContext.value ? 'i-heroicons-shield-check-20-solid' : 'i-heroicons-briefcase-20-solid',
)

const tenantLabel = computed(() => {
  if (auth.isSuperAdmin.value && !tenant.currentTenant.value) {
    return t('userMenuSuperAdmin')
  }
  return tenant.currentTenant.value?.name ?? ''
})

const roleLabel = computed(() => activeRole.currentRole.value?.name ?? '')

const contextLine = computed(() => {
  const parts = [tenantLabel.value, roleLabel.value].filter(Boolean)
  return parts.join(' · ')
})

const showRoleSection = computed(
  () => !auth.isSuperAdmin.value && activeRole.roles.value.length > 1,
)

const showTenantSection = computed(() => tenant.tenants.value.length > 1)

function toggle() {
  open.value = !open.value
}

function close() {
  open.value = false
}

function setTheme(t: ThemeId) {
  setUiKitTheme(t)
}

async function goWorkspace(path: '/' | '/admin') {
  close()
  if (route.path === path) return
  if (path === '/admin' && !auth.isSuperAdmin.value) return
  await navigateTo(path)
}

async function onSelectRole(id: string) {
  if (!id || id === activeRole.currentRoleId.value || rolePending.value) return
  rolePending.value = true
  roleError.value = ''
  try {
    await activeRole.switchRole(id)
    await refreshNuxtData()
    close()
  } catch (e) {
    roleError.value = e instanceof Error ? e.message : t('roleSwitchFailed')
  } finally {
    rolePending.value = false
  }
}

async function onSelectTenant(id: string) {
  if (!id || id === tenant.currentTenantId.value || tenantPending.value) return
  tenantPending.value = true
  tenantError.value = ''
  try {
    await tenant.switchTenant(id)
    await refreshNuxtData()
    close()
  } catch (e) {
    tenantError.value = e instanceof Error ? e.message : t('tenantSwitchFailed')
  } finally {
    tenantPending.value = false
  }
}

async function handleLogout() {
  close()
  await logout()
}

async function ensureRoles() {
  if (auth.isSuperAdmin.value || !tenant.currentTenantId.value) return
  if (activeRole.roles.value.length > 0) return
  try {
    await activeRole.fetchRoles()
  } catch {
    // 静默：菜单内会再次显示错误
  }
}

async function ensureTenants() {
  if (tenant.tenants.value.length > 0) return
  try {
    await tenant.fetchTenants()
  } catch {
    // 静默
  }
}

function onDocumentClick(e: MouseEvent) {
  if (!open.value || !rootRef.value) return
  if (!rootRef.value.contains(e.target as Node)) close()
}

function onEscape(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
}

onMounted(() => {
  void ensureRoles()
  void ensureTenants()
  document.addEventListener('click', onDocumentClick)
  document.addEventListener('keydown', onEscape)
})

onUnmounted(() => {
  document.removeEventListener('click', onDocumentClick)
  document.removeEventListener('keydown', onEscape)
})

watch(
  () => tenant.currentTenantId.value,
  () => {
    activeRole.clearRoles()
    void ensureRoles()
  },
)
</script>
