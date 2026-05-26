<template>
  <form
    class="space-y-8"
    data-testid="tenant-settings-form"
    @submit.prevent="handleSave"
  >
    <!-- Group · Workspace identity -->
    <fieldset class="space-y-4">
      <div class="flex items-center gap-2">
        <UIcon name="i-heroicons-identification" class="h-4 w-4 text-ds-fg-brand" aria-hidden="true" />
        <legend class="text-sm font-semibold uppercase tracking-wider text-ds-fg-brand">
          {{ $t('settingsGroupIdentity') }}
        </legend>
      </div>
      <div class="grid gap-5 sm:grid-cols-2">
        <div class="space-y-1.5">
          <label class="flex items-center justify-between text-sm font-medium text-ds-fg-heading">
            <span>{{ $t('settingsTenantName') }}</span>
            <span v-if="readonly" class="text-xs font-normal text-ds-fg-muted">{{ $t('settingsReadonly') }}</span>
          </label>
          <input
            v-model="form.tenant_name"
            type="text"
            class="ds-input w-full rounded-xl px-4 py-2.5 text-sm shadow-ds-sm disabled:cursor-not-allowed disabled:opacity-60"
            data-testid="settings-tenant-name"
            :disabled="readonly"
            :placeholder="$t('settingsTenantNamePh')"
          >
          <p class="text-xs text-ds-fg-muted">{{ $t('settingsTenantNameHelp') }}</p>
        </div>
      </div>
    </fieldset>

    <!-- Group · Locale & time -->
    <fieldset class="space-y-4 border-t border-ds-border-muted pt-8">
      <div class="flex items-center gap-2">
        <UIcon name="i-heroicons-globe-alt" class="h-4 w-4 text-ds-fg-brand" aria-hidden="true" />
        <legend class="text-sm font-semibold uppercase tracking-wider text-ds-fg-brand">
          {{ $t('settingsGroupLocale') }}
        </legend>
      </div>
      <div class="grid gap-5 sm:grid-cols-2">
        <div class="space-y-1.5">
          <label class="block text-sm font-medium text-ds-fg-heading">{{ $t('settingsTimezone') }}</label>
          <div class="relative">
            <select
              v-model="form.timezone"
              class="ds-input w-full appearance-none rounded-xl px-4 py-2.5 pr-10 text-sm shadow-ds-sm disabled:cursor-not-allowed disabled:opacity-60"
              data-testid="settings-timezone"
              :disabled="readonly"
            >
              <option v-for="tz in timezones" :key="tz" :value="tz">{{ tz }}</option>
            </select>
            <UIcon
              name="i-heroicons-chevron-down"
              class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ds-fg-muted"
              aria-hidden="true"
            />
          </div>
          <p class="text-xs text-ds-fg-muted">{{ $t('settingsTimezoneHelp') }}</p>
        </div>

        <div class="space-y-1.5">
          <label class="block text-sm font-medium text-ds-fg-heading">{{ $t('settingsDefaultLocale') }}</label>
          <div class="relative">
            <select
              v-model="form.default_locale"
              class="ds-input w-full appearance-none rounded-xl px-4 py-2.5 pr-10 text-sm shadow-ds-sm disabled:cursor-not-allowed disabled:opacity-60"
              data-testid="settings-locale"
              :disabled="readonly"
            >
              <option value="zh-CN">中文 (zh-CN)</option>
              <option value="en-US">English (en-US)</option>
            </select>
            <UIcon
              name="i-heroicons-chevron-down"
              class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ds-fg-muted"
              aria-hidden="true"
            />
          </div>
          <p class="text-xs text-ds-fg-muted">{{ $t('settingsDefaultLocaleHelp') }}</p>
        </div>
      </div>
    </fieldset>

    <!-- Group · Business switches -->
    <fieldset class="space-y-4 border-t border-ds-border-muted pt-8">
      <div class="flex items-center gap-2">
        <UIcon name="i-heroicons-adjustments-horizontal" class="h-4 w-4 text-ds-fg-brand" aria-hidden="true" />
        <legend class="text-sm font-semibold uppercase tracking-wider text-ds-fg-brand">
          {{ $t('settingsGroupBusiness') }}
        </legend>
      </div>

      <!-- Lead import mode as segmented control -->
      <div class="space-y-2">
        <label class="block text-sm font-medium text-ds-fg-heading">{{ $t('settingsLeadImportMode') }}</label>
        <div
          class="inline-flex rounded-xl border border-ds-border bg-ds-bg-muted p-1"
          role="radiogroup"
          :aria-label="$t('settingsLeadImportMode')"
        >
          <button
            v-for="opt in importModeOptions"
            :key="opt.value"
            type="button"
            role="radio"
            :aria-checked="form.lead_import_mode === opt.value"
            class="inline-flex cursor-pointer items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-colors duration-200 disabled:cursor-not-allowed disabled:opacity-60"
            :class="form.lead_import_mode === opt.value
              ? 'bg-ds-bg-elevated text-ds-fg-brand shadow-ds-sm ring-1 ring-inset ring-ds-border'
              : 'text-ds-fg-muted hover:text-ds-fg-heading'"
            :data-testid="opt.value === 'manual_review' ? 'settings-import-manual' : 'settings-import-auto'"
            :disabled="readonly"
            @click="form.lead_import_mode = opt.value"
          >
            <UIcon :name="opt.icon" class="h-4 w-4" aria-hidden="true" />
            <span>{{ opt.label }}</span>
          </button>
        </div>
        <!-- Hidden mirror to preserve original test selector contract -->
        <select
          v-model="form.lead_import_mode"
          data-testid="settings-import-mode"
          class="sr-only"
          tabindex="-1"
          aria-hidden="true"
          :disabled="readonly"
        >
          <option value="manual_review">{{ $t('settingsImportManual') }}</option>
          <option value="auto_merge">{{ $t('settingsImportAuto') }}</option>
        </select>
        <p class="text-xs text-ds-fg-muted">{{ activeImportHelp }}</p>
      </div>

      <!-- AI Preview as toggle -->
      <div class="flex items-start justify-between gap-4 rounded-xl border border-ds-border-muted bg-ds-bg-muted/40 p-4">
        <div class="flex min-w-0 items-start gap-3">
          <span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-ds-ai-subtle text-ds-ai" aria-hidden="true">
            <UIcon name="i-heroicons-sparkles" class="h-5 w-5" />
          </span>
          <div class="min-w-0">
            <p class="text-sm font-medium text-ds-fg-heading">{{ $t('settingsAiPreview') }}</p>
            <p class="mt-0.5 text-xs leading-relaxed text-ds-fg-muted">{{ $t('settingsAiPreviewHelp') }}</p>
          </div>
        </div>
        <label class="relative inline-flex cursor-pointer items-center" :class="{ 'cursor-not-allowed opacity-60': readonly }">
          <input
            v-model="form.ai_preview_enabled"
            type="checkbox"
            class="peer sr-only"
            data-testid="settings-ai-preview"
            :disabled="readonly"
          >
          <span
            class="h-6 w-11 rounded-full bg-ds-border transition-colors duration-200 peer-checked:bg-ds-brand peer-focus-visible:ring-2 peer-focus-visible:ring-ds-brand peer-focus-visible:ring-offset-2 peer-focus-visible:ring-offset-ds-bg"
            aria-hidden="true"
          />
          <span
            class="pointer-events-none absolute left-0.5 top-0.5 h-5 w-5 rounded-full bg-white shadow-ds-sm transition-transform duration-200 peer-checked:translate-x-5"
            aria-hidden="true"
          />
        </label>
      </div>
    </fieldset>

    <!-- Group · Sales quota -->
    <fieldset class="space-y-4 border-t border-ds-border-muted pt-8">
      <div class="flex items-center gap-2">
        <UIcon name="i-heroicons-currency-dollar" class="h-4 w-4 text-ds-fg-brand" aria-hidden="true" />
        <legend class="text-sm font-semibold uppercase tracking-wider text-ds-fg-brand">
          {{ $t('settingsGroupQuota') }}
        </legend>
      </div>
      <p class="text-xs text-ds-fg-muted">{{ $t('settingsSalesQuotaHelp') }}</p>
      <div class="grid gap-5 sm:grid-cols-3">
        <div class="space-y-1.5 sm:col-span-1">
          <label class="block text-sm font-medium text-ds-fg-heading">{{ $t('settingsQuotaAmount') }}</label>
          <div class="relative">
            <span class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-sm font-medium text-ds-fg-muted">
              {{ currencySymbol }}
            </span>
            <input
              v-model.number="form.quota_amount"
              type="number"
              min="0"
              class="ds-input w-full rounded-xl pl-8 pr-3 py-2.5 text-sm tabular-nums shadow-ds-sm disabled:cursor-not-allowed disabled:opacity-60"
              data-testid="settings-quota-amount"
              :disabled="readonly"
            >
          </div>
        </div>
        <div class="space-y-1.5">
          <label class="block text-sm font-medium text-ds-fg-heading">{{ $t('settingsQuotaCurrency') }}</label>
          <div class="relative">
            <select
              v-model="form.quota_currency"
              class="ds-input w-full appearance-none rounded-xl px-4 py-2.5 pr-10 text-sm shadow-ds-sm disabled:cursor-not-allowed disabled:opacity-60"
              data-testid="settings-quota-currency"
              :disabled="readonly"
            >
              <option value="CNY">CNY · ¥</option>
              <option value="USD">USD · $</option>
            </select>
            <UIcon
              name="i-heroicons-chevron-down"
              class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ds-fg-muted"
              aria-hidden="true"
            />
          </div>
        </div>
        <div class="space-y-1.5">
          <label class="block text-sm font-medium text-ds-fg-heading">{{ $t('settingsQuotaPeriod') }}</label>
          <input
            v-model="form.quota_period"
            type="month"
            class="ds-input w-full rounded-xl px-4 py-2.5 text-sm shadow-ds-sm disabled:cursor-not-allowed disabled:opacity-60"
            data-testid="settings-quota-period"
            :disabled="readonly"
          >
        </div>
      </div>
    </fieldset>

    <!-- Sticky action footer -->
    <div
      v-if="!readonly"
      class="sticky bottom-0 -mx-6 -mb-6 flex items-center justify-between gap-3 border-t border-ds-border-muted bg-ds-bg-elevated/95 px-6 py-4 backdrop-blur sm:-mx-8 sm:-mb-8 sm:px-8"
    >
      <div class="flex items-center gap-2 text-xs">
        <template v-if="savedOk">
          <UIcon name="i-heroicons-check-circle" class="h-4 w-4 text-ds-success" aria-hidden="true" />
          <span class="font-medium text-ds-success">{{ $t('settingsSaved') }}</span>
        </template>
        <template v-else-if="dirty">
          <span class="inline-block h-2 w-2 rounded-full bg-ds-warning" aria-hidden="true" />
          <span class="font-medium text-ds-warning">{{ $t('settingsUnsavedChanges') }}</span>
        </template>
        <template v-else>
          <UIcon name="i-heroicons-check" class="h-4 w-4 text-ds-fg-muted" aria-hidden="true" />
          <span class="text-ds-fg-muted">{{ $t('settingsUpToDate') }}</span>
        </template>
      </div>
      <div class="flex items-center gap-2">
        <button
          type="button"
          class="cursor-pointer rounded-xl border border-ds-border bg-ds-bg-elevated px-4 py-2 text-sm font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-border hover:text-ds-fg-heading disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="!dirty || saving"
          data-testid="settings-reset-btn"
          @click="resetForm"
        >
          {{ $t('settingsDiscard') }}
        </button>
        <button
          type="submit"
          class="ds-btn-primary inline-flex cursor-pointer items-center gap-2 rounded-xl px-5 py-2 text-sm font-semibold transition-colors duration-200 disabled:cursor-not-allowed disabled:opacity-60"
          data-testid="settings-save-btn"
          :disabled="saving || !dirty"
        >
          <UIcon
            v-if="saving"
            name="i-heroicons-arrow-path"
            class="h-4 w-4 animate-spin"
            aria-hidden="true"
          />
          <UIcon
            v-else
            name="i-heroicons-check-circle"
            class="h-4 w-4"
            aria-hidden="true"
          />
          <span>{{ saving ? $t('loading') : $t('save') }}</span>
        </button>
      </div>
    </div>
  </form>
</template>

<script setup lang="ts">
import type { TenantSettings } from '~/types/settings'

const props = defineProps<{
  settings: TenantSettings
  readonly?: boolean
}>()

const emit = defineEmits<{ save: [payload: ReturnType<typeof buildPayload>] }>()

const { t } = useI18n()

const saving = ref(false)
const savedOk = ref(false)

function snapshot(s: TenantSettings) {
  return {
    tenant_name: s.tenant_name,
    timezone: s.timezone,
    default_locale: s.default_locale,
    lead_import_mode: s.business_switches.lead_import_mode,
    ai_preview_enabled: s.business_switches.ai_preview_enabled,
    quota_amount: s.sales_quota.amount,
    quota_currency: s.sales_quota.currency,
    quota_period: s.sales_quota.period,
  }
}

const form = reactive(snapshot(props.settings))
const initial = ref(snapshot(props.settings))

const timezones = [
  'Asia/Shanghai',
  'Asia/Tokyo',
  'America/New_York',
  'America/Los_Angeles',
  'Europe/London',
  'UTC',
]

const importModeOptions = computed(() => [
  {
    value: 'manual_review' as const,
    label: t('settingsImportManual'),
    icon: 'i-heroicons-clipboard-document-check',
    help: t('settingsImportManualHelp'),
  },
  {
    value: 'auto_merge' as const,
    label: t('settingsImportAuto'),
    icon: 'i-heroicons-bolt',
    help: t('settingsImportAutoHelp'),
  },
])

const activeImportHelp = computed(
  () => importModeOptions.value.find((o) => o.value === form.lead_import_mode)?.help ?? '',
)

const currencySymbol = computed(() => {
  switch (form.quota_currency) {
    case 'USD':
      return '$'
    case 'CNY':
    default:
      return '¥'
  }
})

const dirty = computed(() => {
  const a = JSON.stringify(initial.value)
  const b = JSON.stringify({ ...form })
  return a !== b
})

watch(
  () => props.settings,
  (s) => {
    const snap = snapshot(s)
    Object.assign(form, snap)
    initial.value = snap
  },
)

watch(
  () => ({ ...form }),
  () => {
    if (savedOk.value) savedOk.value = false
  },
  { deep: true },
)

function resetForm() {
  Object.assign(form, initial.value)
  savedOk.value = false
}

function buildPayload() {
  return {
    tenant_name: form.tenant_name,
    default_locale: form.default_locale,
    timezone: form.timezone,
    business_switches: {
      ai_preview_enabled: form.ai_preview_enabled,
      lead_import_mode: form.lead_import_mode as 'manual_review' | 'auto_merge',
    },
    sales_quota: {
      amount: form.quota_amount,
      currency: form.quota_currency,
      period: form.quota_period,
    },
  }
}

async function handleSave() {
  if (saving.value || !dirty.value) return
  saving.value = true
  savedOk.value = false
  try {
    emit('save', buildPayload())
    initial.value = { ...form }
    savedOk.value = true
    setTimeout(() => {
      savedOk.value = false
    }, 2400)
  } finally {
    saving.value = false
  }
}

defineExpose({ saving, savedOk, dirty })
</script>
