<template>
  <div class="flex min-h-[520px] flex-col gap-4 lg:flex-row" data-testid="role-permission-manager">
    <!-- 左侧：角色列表 -->
    <aside class="flex w-full shrink-0 flex-col rounded-2xl border border-ds-border-muted bg-ds-bg-muted/40 p-3 lg:w-72">
      <div class="mb-3 flex items-center justify-between gap-2 px-1">
        <div class="flex items-center gap-2">
          <UIcon name="i-heroicons-user-group" class="h-4 w-4 text-ds-fg-brand" aria-hidden="true" />
          <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('roleListTitle') }}</h3>
          <span class="rounded-full bg-ds-bg-elevated px-1.5 py-0.5 text-[10px] font-semibold text-ds-fg-muted ring-1 ring-inset ring-ds-border-muted">
            {{ roles.length }}
          </span>
        </div>
        <button
          v-if="canManage"
          type="button"
          class="ds-btn-primary inline-flex cursor-pointer items-center gap-1 rounded-lg px-2.5 py-1 text-xs font-semibold transition-colors duration-200"
          data-testid="role-create-btn"
          @click="openCreateRole"
        >
          <UIcon name="i-heroicons-plus" class="h-3.5 w-3.5" aria-hidden="true" />
          <span>{{ $t('createRole') }}</span>
        </button>
      </div>

      <div class="relative mb-2">
        <UIcon
          name="i-heroicons-magnifying-glass"
          class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ds-fg-muted"
          aria-hidden="true"
        />
        <input
          v-model="roleSearch"
          type="search"
          class="ds-input w-full rounded-lg py-1.5 pl-9 pr-3 text-sm"
          :placeholder="$t('roleSearchPh')"
          data-testid="role-search"
        >
      </div>

      <p v-if="loading" class="py-8 text-center text-sm text-ds-fg-muted">{{ $t('loading') }}</p>
      <p v-else-if="filteredRoles.length === 0" class="py-8 text-center text-xs text-ds-fg-muted">{{ $t('roleNoMatch') }}</p>
      <ul v-else class="space-y-1 overflow-y-auto pr-0.5" data-testid="role-list">
        <li v-for="role in filteredRoles" :key="role.id">
          <button
            type="button"
            class="group flex w-full cursor-pointer flex-col gap-1 rounded-xl border border-transparent px-3 py-2.5 text-left text-sm transition-colors duration-200"
            :class="selectedRole?.id === role.id
              ? 'border-ds-border bg-ds-bg-elevated text-ds-fg-brand shadow-ds-sm'
              : 'text-ds-fg-heading hover:bg-ds-bg-muted'"
            :data-testid="`role-item-${role.id}`"
            @click="selectRole(role)"
          >
            <span class="flex items-center justify-between gap-2">
              <span class="truncate font-semibold">{{ role.name }}</span>
              <span
                v-if="role.is_system"
                class="rounded-full bg-amber-100 px-1.5 py-0.5 text-[10px] font-semibold text-amber-700 ring-1 ring-inset ring-amber-200 dark:bg-amber-500/15 dark:text-amber-300 dark:ring-amber-500/30"
              >
                {{ $t('system') }}
              </span>
            </span>
            <span class="flex flex-wrap items-center gap-1.5 text-[11px] text-ds-fg-muted">
              <span class="inline-flex items-center gap-1">
                <UIcon name="i-heroicons-users" class="h-3 w-3" aria-hidden="true" />
                {{ $t('roleMemberCount', { n: role.user_count }) }}
              </span>
              <span aria-hidden="true">·</span>
              <span class="inline-flex items-center gap-1">
                <UIcon name="i-heroicons-key" class="h-3 w-3" aria-hidden="true" />
                {{ $t('rolePermCount', { n: role.permission_ids.length }) }}
              </span>
            </span>
          </button>
        </li>
      </ul>
    </aside>

    <!-- 右侧：权限树 -->
    <section class="min-w-0 flex-1 rounded-2xl border border-ds-border-muted bg-ds-bg-elevated p-5 shadow-ds-sm">
      <template v-if="!selectedRole">
        <div class="flex flex-col items-center justify-center gap-3 py-16 text-center">
          <span class="flex h-12 w-12 items-center justify-center rounded-2xl bg-ds-brand-subtle text-ds-fg-brand" aria-hidden="true">
            <UIcon name="i-heroicons-shield-check" class="h-6 w-6" />
          </span>
          <p class="text-sm text-ds-fg-muted">{{ $t('roleSelectHint') }}</p>
        </div>
      </template>

      <template v-else>
        <div class="mb-5 flex flex-wrap items-start justify-between gap-3 border-b border-ds-border-muted pb-4">
          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-center gap-2">
              <h3 class="text-lg font-semibold text-ds-fg-heading" data-testid="role-perm-title">{{ selectedRole.name }}</h3>
              <span
                v-if="selectedRole.is_system"
                class="rounded-full bg-amber-100 px-2 py-0.5 text-[10px] font-semibold text-amber-700 ring-1 ring-inset ring-amber-200 dark:bg-amber-500/15 dark:text-amber-300 dark:ring-amber-500/30"
              >
                {{ $t('system') }}
              </span>
            </div>
            <p class="mt-1 text-sm text-ds-fg-muted">{{ selectedRole.description || $t('roleManagementHint') }}</p>
            <div class="mt-3 flex flex-wrap items-center gap-2">
              <span class="inline-flex items-center gap-1.5 rounded-full border border-ds-border-muted bg-ds-bg-muted px-2.5 py-1 text-xs text-ds-fg-muted">
                <UIcon name="i-heroicons-users" class="h-3.5 w-3.5" aria-hidden="true" />
                {{ $t('roleMemberCount', { n: selectedRole.user_count }) }}
              </span>
              <span class="inline-flex items-center gap-1.5 rounded-full border border-ds-border-muted bg-ds-bg-muted px-2.5 py-1 text-xs text-ds-fg-muted">
                <UIcon name="i-heroicons-key" class="h-3.5 w-3.5" aria-hidden="true" />
                {{ $t('rolePermSelected', { selected: draftPermIds.length, total: totalPermissionCount }) }}
              </span>
              <span
                v-if="dirty"
                class="inline-flex items-center gap-1.5 rounded-full border border-ds-warning/30 bg-ds-warning-subtle px-2.5 py-1 text-xs font-medium text-ds-warning"
              >
                <span class="inline-block h-1.5 w-1.5 rounded-full bg-ds-warning" aria-hidden="true" />
                {{ $t('settingsUnsavedChanges') }}
              </span>
            </div>
          </div>
          <div v-if="canManage" class="flex shrink-0 items-center gap-2">
            <button
              type="button"
              class="cursor-pointer rounded-xl border border-ds-border bg-ds-bg-elevated px-4 py-2 text-sm font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-border hover:text-ds-fg-heading disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="!dirty || saving"
              data-testid="role-perm-reset"
              @click="resetDraft"
            >
              {{ $t('settingsDiscard') }}
            </button>
            <button
              type="button"
              class="ds-btn-primary inline-flex cursor-pointer items-center gap-2 rounded-xl px-4 py-2 text-sm font-semibold disabled:cursor-not-allowed disabled:opacity-60"
              data-testid="role-perm-save"
              :disabled="saving || !dirty"
              @click="savePermissions"
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
              <span>{{ saving ? $t('loading') : $t('savePermissions') }}</span>
            </button>
          </div>
        </div>

        <p
          v-if="!canManage"
          class="mb-4 flex items-center gap-2 rounded-xl border border-ds-border-muted bg-ds-bg-muted px-3 py-2 text-sm text-ds-fg-muted"
        >
          <UIcon name="i-heroicons-eye" class="h-4 w-4 shrink-0" aria-hidden="true" />
          {{ $t('roleReadonlyHint') }}
        </p>

        <p
          v-if="saveError"
          class="mb-4 flex items-center gap-2 rounded-xl border border-ds-danger/20 bg-ds-danger-subtle px-3 py-2 text-sm text-ds-danger"
        >
          <UIcon name="i-heroicons-exclamation-circle" class="h-4 w-4 shrink-0" aria-hidden="true" />
          {{ saveError }}
        </p>

        <div class="max-h-[56vh] space-y-4 overflow-y-auto pr-1" data-testid="permission-tree">
          <div
            v-for="mod in permissionModules"
            :key="mod.key"
            class="rounded-xl border border-ds-border-muted"
          >
            <div class="flex items-center gap-2 border-b border-ds-border-muted bg-ds-bg-muted/50 px-4 py-2.5">
              <input
                :id="`mod-${mod.key}`"
                type="checkbox"
                class="h-4 w-4 cursor-pointer rounded border-ds-border text-ds-brand"
                :disabled="!canManage"
                :checked="isModuleChecked(mod)"
                :indeterminate="isModuleIndeterminate(mod)"
                @change="toggleModule(mod, ($event.target as HTMLInputElement).checked)"
              >
              <label :for="`mod-${mod.key}`" class="cursor-pointer text-sm font-semibold text-ds-fg-heading">
                {{ $t(`permModule.${mod.key}`) }}
              </label>
            </div>

            <div class="divide-y divide-ds-border-muted">
              <div
                v-for="res in mod.resources"
                :key="res.resource"
                class="px-4 py-3"
              >
                <div class="flex items-center gap-2">
                  <input
                    :id="`res-${res.resource}`"
                    type="checkbox"
                    class="h-4 w-4 cursor-pointer rounded border-ds-border text-ds-brand"
                    :disabled="!canManage"
                    :checked="isResourceChecked(res)"
                    :indeterminate="isResourceIndeterminate(res)"
                    @change="toggleResource(res, ($event.target as HTMLInputElement).checked)"
                  >
                  <label :for="`res-${res.resource}`" class="cursor-pointer text-sm font-medium text-ds-fg-brand">
                    {{ $t(`permResource.${res.resource}`) }}
                  </label>
                </div>
                <ul class="ml-6 mt-2 space-y-1.5 border-l border-ds-border-muted pl-4">
                  <li v-for="perm in res.items" :key="perm.id" class="flex items-start gap-2">
                    <input
                      :id="`perm-${perm.id}`"
                      v-model="draftPermIds"
                      type="checkbox"
                      class="mt-0.5 h-4 w-4 shrink-0 cursor-pointer rounded border-ds-border text-ds-brand"
                      :value="perm.id"
                      :disabled="!canManage"
                      :data-testid="`perm-${perm.resource}-${perm.action}`"
                    >
                    <label :for="`perm-${perm.id}`" class="cursor-pointer text-sm text-ds-fg-heading">
                      <span class="font-mono text-xs text-ds-fg-muted">{{ perm.resource }}:{{ perm.action }}</span>
                      <span v-if="perm.description" class="ml-2 text-ds-fg-muted">— {{ perm.description }}</span>
                    </label>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </template>
    </section>

    <!-- 新建角色 -->
    <div
      v-if="createOpen"
      class="fixed inset-0 z-ds-modal flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
      data-testid="role-create-modal"
      @click.self="createOpen = false"
    >
      <form
        class="ds-card flex w-full max-w-md flex-col overflow-hidden rounded-2xl shadow-ds-xl"
        @submit.prevent="submitCreateRole"
      >
        <header class="flex items-start gap-3 border-b border-ds-border-muted bg-ds-bg-muted/40 px-6 py-4">
          <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-ds-brand-subtle text-ds-fg-brand" aria-hidden="true">
            <UIcon name="i-heroicons-user-plus" class="h-5 w-5" />
          </span>
          <div class="min-w-0 flex-1">
            <h3 class="text-base font-semibold text-ds-fg-heading">{{ $t('createRole') }}</h3>
            <p class="mt-0.5 text-xs text-ds-fg-muted">{{ $t('createRoleHint') }}</p>
          </div>
          <button
            type="button"
            class="inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-ds-fg-muted transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-heading"
            :aria-label="$t('cancel')"
            @click="createOpen = false"
          >
            <UIcon name="i-heroicons-x-mark" class="h-5 w-5" aria-hidden="true" />
          </button>
        </header>
        <div class="space-y-4 px-6 py-5">
          <div class="space-y-1.5">
            <label class="block text-xs font-semibold uppercase tracking-wide text-ds-fg-muted">{{ $t('name') }}</label>
            <input v-model="createName" required class="ds-input w-full rounded-xl px-3 py-2 text-sm" placeholder="Sales Manager">
          </div>
          <div class="space-y-1.5">
            <label class="block text-xs font-semibold uppercase tracking-wide text-ds-fg-muted">{{ $t('description') }}</label>
            <input v-model="createDesc" class="ds-input w-full rounded-xl px-3 py-2 text-sm" :placeholder="$t('roleDescPh')">
          </div>
          <p
            v-if="createError"
            class="flex items-center gap-2 rounded-lg border border-ds-danger/20 bg-ds-danger-subtle px-3 py-2 text-sm text-ds-danger"
          >
            <UIcon name="i-heroicons-exclamation-circle" class="h-4 w-4 shrink-0" aria-hidden="true" />
            {{ createError }}
          </p>
        </div>
        <footer class="flex items-center justify-end gap-2 border-t border-ds-border-muted bg-ds-bg-muted/30 px-6 py-4">
          <button type="button" class="cursor-pointer rounded-xl border border-ds-border bg-ds-bg-elevated px-4 py-2 text-sm font-medium text-ds-fg-muted transition-colors duration-200 hover:text-ds-fg-heading" @click="createOpen = false">
            {{ $t('cancel') }}
          </button>
          <button
            type="submit"
            class="ds-btn-primary inline-flex cursor-pointer items-center gap-2 rounded-xl px-4 py-2 text-sm font-semibold disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="creating"
          >
            <UIcon
              v-if="creating"
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
            <span>{{ creating ? $t('loading') : $t('save') }}</span>
          </button>
        </footer>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { PermissionDictItem, RoleItem } from '~/composables/use-rbac'
import {
  buildPermissionModules,
  modulePermissionIds,
  resourcePermissionIds,
  type PermissionModuleNode,
  type PermissionResourceNode,
} from '~/utils/rbac-permission-groups'

const { t } = useI18n()
const rbac = useRbac()
const { can } = usePermission()

const canManage = computed(() => can('rbac', 'manage'))
const canView = computed(() => can('rbac', 'view') || canManage.value)

const loading = ref(true)
const roles = ref<RoleItem[]>([])
const permissionItems = ref<PermissionDictItem[]>([])
const selectedRole = ref<RoleItem | null>(null)
const draftPermIds = ref<string[]>([])
const baselinePermIds = ref<string[]>([])
const saving = ref(false)
const saveError = ref('')
const roleSearch = ref('')

const filteredRoles = computed(() => {
  const q = roleSearch.value.trim().toLowerCase()
  if (!q) return roles.value
  return roles.value.filter((r) => {
    const desc = (r.description || '').toLowerCase()
    return r.name.toLowerCase().includes(q) || desc.includes(q)
  })
})

const totalPermissionCount = computed(() => permissionItems.value.length)

const createOpen = ref(false)
const createName = ref('')
const createDesc = ref('')
const createError = ref('')
const creating = ref(false)

const permissionModules = computed(() => buildPermissionModules(permissionItems.value))

const dirty = computed(() => {
  const a = [...draftPermIds.value].sort().join(',')
  const b = [...baselinePermIds.value].sort().join(',')
  return a !== b
})

function isChecked(ids: string[], pool: string[]) {
  return ids.length > 0 && ids.every((id) => pool.includes(id))
}

function isIndeterminate(ids: string[], pool: string[]) {
  const hit = ids.filter((id) => pool.includes(id)).length
  return hit > 0 && hit < ids.length
}

function isModuleChecked(mod: PermissionModuleNode) {
  return isChecked(modulePermissionIds(mod), draftPermIds.value)
}

function isModuleIndeterminate(mod: PermissionModuleNode) {
  return isIndeterminate(modulePermissionIds(mod), draftPermIds.value)
}

function isResourceChecked(res: PermissionResourceNode) {
  return isChecked(resourcePermissionIds(res), draftPermIds.value)
}

function isResourceIndeterminate(res: PermissionResourceNode) {
  return isIndeterminate(resourcePermissionIds(res), draftPermIds.value)
}

function setIds(add: string[], remove: string[], checked: boolean) {
  const set = new Set(draftPermIds.value)
  for (const id of remove) set.delete(id)
  if (checked) for (const id of add) set.add(id)
  draftPermIds.value = [...set]
}

function toggleModule(mod: PermissionModuleNode, checked: boolean) {
  setIds(modulePermissionIds(mod), modulePermissionIds(mod), checked)
}

function toggleResource(res: PermissionResourceNode, checked: boolean) {
  const ids = resourcePermissionIds(res)
  setIds(ids, ids, checked)
}

function selectRole(role: RoleItem) {
  selectedRole.value = role
  draftPermIds.value = [...role.permission_ids]
  baselinePermIds.value = [...role.permission_ids]
  saveError.value = ''
}

function resetDraft() {
  draftPermIds.value = [...baselinePermIds.value]
  saveError.value = ''
}

async function load() {
  loading.value = true
  try {
    ;[roles.value, permissionItems.value] = await Promise.all([
      rbac.fetchRoles(),
      rbac.fetchPermissionItems(),
    ])
    if (!selectedRole.value && roles.value.length > 0) {
      selectRole(roles.value[0])
    } else if (selectedRole.value) {
      const fresh = roles.value.find((r) => r.id === selectedRole.value!.id)
      if (fresh) selectRole(fresh)
    }
  } finally {
    loading.value = false
  }
}

async function savePermissions() {
  if (!selectedRole.value || !canManage.value) return
  saving.value = true
  saveError.value = ''
  try {
    const updated = await rbac.assignRolePermissions(selectedRole.value.id, draftPermIds.value)
    selectedRole.value = updated
    baselinePermIds.value = [...draftPermIds.value]
    await load()
    await rbac.loadMyPermissions()
  } catch (e) {
    saveError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    saving.value = false
  }
}

function openCreateRole() {
  createName.value = ''
  createDesc.value = ''
  createError.value = ''
  createOpen.value = true
}

async function submitCreateRole() {
  creating.value = true
  createError.value = ''
  try {
    const created = await rbac.createRole({ name: createName.value, description: createDesc.value })
    createOpen.value = false
    await load()
    const fresh = roles.value.find((r) => r.id === created.id)
    if (fresh) selectRole(fresh)
  } catch (e) {
    createError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    creating.value = false
  }
}

onMounted(async () => {
  if (!canView.value) return
  await load()
})

defineExpose({ load })
</script>
