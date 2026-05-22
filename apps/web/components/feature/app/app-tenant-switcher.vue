<template>
  <div class="flex min-w-0 items-center gap-2" data-testid="tenant-switcher">
    <label v-if="showLabel" class="shrink-0 text-xs text-ds-fg-muted" for="tenant-select">
      {{ $t('tenantSwitcherLabel') }}
    </label>
    <UiSelect
      id="tenant-select"
      :model-value="tenant.currentTenantId.value ?? ''"
      class="min-w-[10rem] max-w-[14rem]"
      :items="selectItems"
      :placeholder="$t('tenantSelectPlaceholder')"
      :disabled="pending || selectItems.length === 0"
      @update:model-value="onSelect"
    />
    <p v-if="error" class="max-w-[12rem] truncate text-xs text-ds-danger" role="alert">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    showLabel?: boolean
  }>(),
  { showLabel: true },
)

const { t } = useI18n()
const tenant = useTenant()

const pending = ref(false)
const error = ref('')

const selectItems = computed(() =>
  tenant.tenants.value.map((item) => ({
    label: item.name,
    value: item.id,
  })),
)

async function ensureTenants() {
  if (tenant.tenants.value.length > 0) return
  pending.value = true
  error.value = ''
  try {
    await tenant.fetchTenants()
    if (!tenant.currentTenantId.value && tenant.tenants.value.length > 0) {
      await tenant.switchTenant(tenant.tenants.value[0].id)
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : t('tenantSwitchFailed')
  } finally {
    pending.value = false
  }
}

async function onSelect(id: string) {
  if (!id || id === tenant.currentTenantId.value) return
  pending.value = true
  error.value = ''
  try {
    await tenant.switchTenant(id)
    await refreshNuxtData()
  } catch (e) {
    error.value = e instanceof Error ? e.message : t('tenantSwitchFailed')
  } finally {
    pending.value = false
  }
}

onMounted(ensureTenants)
</script>
