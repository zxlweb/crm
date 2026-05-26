<template>
  <div class="mx-auto w-full max-w-7xl space-y-6" data-testid="settings-page">
    <!-- Hero header -->
    <header
      class="ds-card relative overflow-hidden rounded-2xl px-6 py-6 sm:px-8"
      data-testid="settings-header"
    >
      <div
        class="pointer-events-none absolute inset-y-0 right-0 hidden w-1/2 opacity-60 sm:block"
        aria-hidden="true"
        style="background: radial-gradient(circle at top right, var(--ds-blur-brand) 0%, transparent 70%);"
      />
      <div class="relative flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div class="min-w-0">
          <div class="flex items-center gap-2 text-xs font-medium text-ds-fg-muted">
            <span>{{ $t('crmNavSettings') }}</span>
            <span aria-hidden="true">/</span>
            <span class="text-ds-fg-heading">{{ activeNavLabel }}</span>
          </div>
          <h1 class="mt-2 text-ds-2xl font-bold tracking-tight text-ds-fg-heading sm:text-3xl">
            {{ $t('settingsPageTitle') }}
          </h1>
          <p class="mt-1 max-w-2xl text-sm text-ds-fg-muted">
            {{ $t('settingsPageDesc') }}
          </p>
          <div v-if="settings" class="mt-4 flex flex-wrap items-center gap-2 text-xs">
            <span
              class="inline-flex items-center gap-1.5 rounded-full border border-ds-border-muted bg-ds-bg-muted px-2.5 py-1 text-ds-fg-muted"
            >
              <UIcon name="i-heroicons-clock" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ $t('settingsLastUpdated', { time: formatUpdatedAt(settings.updated_at) }) }}
            </span>
            <span
              v-if="settings.updated_by"
              class="inline-flex items-center gap-1.5 rounded-full border border-ds-border-muted bg-ds-bg-muted px-2.5 py-1 text-ds-fg-muted"
            >
              <UIcon name="i-heroicons-user-circle" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ $t('settingsUpdatedBy', { user: settings.updated_by }) }}
            </span>
            <span
              class="inline-flex items-center gap-1.5 rounded-full border border-ds-success/30 bg-ds-success-subtle px-2.5 py-1 font-medium text-ds-success"
            >
              <span class="inline-block h-1.5 w-1.5 rounded-full bg-ds-success" aria-hidden="true" />
              {{ $t('settingsHealthy') }}
            </span>
          </div>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <NuxtLink
            to="/settings/audit"
            class="inline-flex cursor-pointer items-center gap-2 rounded-xl border border-ds-border bg-ds-bg-elevated px-4 py-2 text-sm font-medium text-ds-fg-heading shadow-ds-sm transition-colors duration-200 hover:border-ds-brand hover:text-ds-fg-brand"
            data-testid="settings-audit-link"
          >
            <UIcon name="i-heroicons-document-magnifying-glass" class="h-4 w-4" aria-hidden="true" />
            <span>{{ $t('auditPageTitle') }}</span>
          </NuxtLink>
        </div>
      </div>
    </header>

    <!-- Loading -->
    <div v-if="pending" class="flex items-center justify-center py-24">
      <div class="h-8 w-8 animate-spin rounded-full border-2 border-ds-brand-muted border-t-ds-brand" />
    </div>

    <!-- Error -->
    <p
      v-else-if="loadError"
      class="flex items-start gap-3 rounded-xl border border-ds-danger/20 bg-ds-danger-subtle px-4 py-3 text-sm text-ds-danger"
    >
      <UIcon name="i-heroicons-exclamation-triangle" class="h-5 w-5 shrink-0" aria-hidden="true" />
      <span>{{ loadError }}</span>
    </p>

    <!-- Master-detail layout -->
    <div v-else class="grid gap-6 lg:grid-cols-[260px_minmax(0,1fr)]">
      <aside class="lg:sticky lg:top-6 lg:self-start">
        <div class="ds-card rounded-2xl p-3">
          <SettingsNav
            v-model="activeTab"
            :items="navItems"
          />
        </div>

        <div class="ds-card mt-4 hidden rounded-2xl p-4 lg:block">
          <p class="text-xs font-semibold uppercase tracking-wider text-ds-fg-muted">
            {{ $t('settingsTipsTitle') }}
          </p>
          <p class="mt-2 text-xs leading-relaxed text-ds-fg-muted">
            {{ activeTabTip }}
          </p>
          <NuxtLink
            to="/settings/audit"
            class="mt-3 inline-flex cursor-pointer items-center gap-1 text-xs font-medium text-ds-fg-brand hover:underline"
          >
            <UIcon name="i-heroicons-arrow-top-right-on-square" class="h-3.5 w-3.5" aria-hidden="true" />
            {{ $t('settingsTipsLink') }}
          </NuxtLink>
        </div>
      </aside>

      <main class="min-w-0">
        <section
          v-if="activeTab === 'general'"
          class="ds-card overflow-hidden rounded-2xl"
          data-testid="settings-section-general"
        >
          <header class="flex items-start gap-3 border-b border-ds-border-muted px-6 py-5">
            <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-ds-border-muted bg-ds-brand-subtle text-ds-fg-brand">
              <UIcon name="i-heroicons-building-office-2" class="h-5 w-5" aria-hidden="true" />
            </span>
            <div class="min-w-0">
              <h2 class="text-lg font-semibold text-ds-fg-heading">{{ $t('settingsGeneralTitle') }}</h2>
              <p class="mt-0.5 text-sm text-ds-fg-muted">{{ $t('settingsGeneralDesc') }}</p>
            </div>
          </header>
          <div class="p-6 sm:p-8">
            <TenantSettingsForm
              v-if="settings"
              :settings="settings"
              :readonly="!canEditSettings"
              @save="handleSave"
            />
          </div>
        </section>

        <section
          v-if="activeTab === 'fields'"
          class="ds-card overflow-hidden rounded-2xl"
          data-testid="settings-section-fields"
        >
          <header class="flex items-start gap-3 border-b border-ds-border-muted px-6 py-5">
            <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-ds-border-muted bg-ds-brand-subtle text-ds-fg-brand">
              <UIcon name="i-heroicons-rectangle-stack" class="h-5 w-5" aria-hidden="true" />
            </span>
            <div class="min-w-0">
              <h2 class="text-lg font-semibold text-ds-fg-heading">{{ $t('cfSectionTitle') }}</h2>
              <p class="mt-0.5 text-sm text-ds-fg-muted">{{ $t('cfSectionDesc') }}</p>
            </div>
          </header>
          <div class="p-6 sm:p-8">
            <CustomFieldManager :can-edit="canEditFields" />
          </div>
        </section>

        <section
          v-if="activeTab === 'roles'"
          class="ds-card overflow-hidden rounded-2xl"
          data-testid="settings-section-roles"
        >
          <header class="flex items-start gap-3 border-b border-ds-border-muted px-6 py-5">
            <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-ds-border-muted bg-ds-brand-subtle text-ds-fg-brand">
              <UIcon name="i-heroicons-shield-check" class="h-5 w-5" aria-hidden="true" />
            </span>
            <div class="min-w-0">
              <h2 class="text-lg font-semibold text-ds-fg-heading">{{ $t('roleManagement') }}</h2>
              <p class="mt-0.5 text-sm text-ds-fg-muted">{{ $t('roleManagementHint') }}</p>
            </div>
          </header>
          <div class="p-4 sm:p-6">
            <RolePermissionManager />
          </div>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TenantSettings, TenantSettingsUpdateInput } from '~/types/settings'

interface SettingsNavItem {
  key: string
  label: string
  description: string
  icon: string
  badge?: string | number
}

definePageMeta({ layout: 'app', middleware: 'auth' })

const { t, locale } = useI18n()
const { can } = usePermission()
const settingsApi = useSettings()

const pending = ref(true)
const loadError = ref('')
const settings = ref<TenantSettings | null>(null)
type SettingsTab = 'general' | 'fields' | 'roles'

const route = useRoute()
const router = useRouter()
const activeTab = ref<SettingsTab>('general')

const canEditSettings = computed(() => can('settings', 'update'))
const canEditFields = computed(() => can('custom_fields', 'update'))
const canViewRoles = computed(() => can('rbac', 'view') || can('rbac', 'manage'))

const navItems = computed<SettingsNavItem[]>(() => {
  const items: SettingsNavItem[] = [
    {
      key: 'general',
      label: t('settingsTabGeneral'),
      description: t('settingsNavGeneralDesc'),
      icon: 'i-heroicons-building-office-2',
    },
    {
      key: 'fields',
      label: t('settingsTabFields'),
      description: t('settingsNavFieldsDesc'),
      icon: 'i-heroicons-rectangle-stack',
    },
  ]
  if (canViewRoles.value) {
    items.push({
      key: 'roles',
      label: t('settingsTabRoles'),
      description: t('settingsNavRolesDesc'),
      icon: 'i-heroicons-shield-check',
    })
  }
  return items
})

const activeNavLabel = computed(
  () => navItems.value.find((i) => i.key === activeTab.value)?.label ?? t('settingsTabGeneral'),
)

const activeTabTip = computed(() => {
  switch (activeTab.value) {
    case 'fields':
      return t('settingsTipFields')
    case 'roles':
      return t('settingsTipRoles')
    default:
      return t('settingsTipGeneral')
  }
})

function syncTabFromRoute() {
  const q = route.query.tab
  if (q === 'roles' && canViewRoles.value) activeTab.value = 'roles'
  else if (q === 'fields') activeTab.value = 'fields'
  else if (q === 'general') activeTab.value = 'general'
}

watch(activeTab, (key) => {
  if (route.query.tab !== key) {
    router.replace({ query: { ...route.query, tab: key } })
  }
})

function formatUpdatedAt(iso: string): string {
  if (!iso) return '—'
  try {
    const d = new Date(iso)
    if (Number.isNaN(d.getTime())) return iso
    const diffMs = Date.now() - d.getTime()
    const minutes = Math.round(diffMs / 60_000)
    if (minutes < 1) return t('justNow')
    if (minutes < 60) return t('minutesAgo', { n: minutes })
    const hours = Math.round(minutes / 60)
    if (hours < 24) return t('hoursAgo', { n: hours })
    const days = Math.round(hours / 24)
    if (days < 30) return t('daysAgo', { n: days })
    return new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium' }).format(d)
  } catch {
    return iso
  }
}

async function load() {
  pending.value = true
  loadError.value = ''
  try {
    settings.value = await settingsApi.fetchTenantSettings()
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

async function handleSave(payload: TenantSettingsUpdateInput) {
  try {
    settings.value = await settingsApi.updateTenantSettings(payload)
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  }
}

onMounted(() => {
  syncTabFromRoute()
  load()
})

watch(() => route.query.tab, syncTabFromRoute)
</script>
