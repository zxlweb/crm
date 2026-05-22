<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-lg font-medium">{{ $t('roleManagement') }}</h2>
        <p class="text-sm text-gray-500">{{ $t('roleManagementHint') }}</p>
      </div>
      <PermissionGuard resource="rbac" action="manage">
        <button
          type="button"
          class="rounded bg-blue-600 px-3 py-1.5 text-sm text-white hover:bg-blue-700"
          @click="openCreate"
        >
          {{ $t('createRole') }}
        </button>
        <template #fallback />
      </PermissionGuard>
    </div>

    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
    <p v-if="pending" class="text-sm text-gray-500">{{ $t('loading') }}</p>

    <section v-else class="rounded-lg bg-white p-4 shadow">
      <table class="w-full text-left text-sm">
        <thead class="border-b text-gray-500">
          <tr>
            <th class="py-2">{{ $t('name') }}</th>
            <th class="py-2">{{ $t('description') }}</th>
            <th class="py-2">{{ $t('permissions') }}</th>
            <th class="py-2">{{ $t('members') }}</th>
            <th class="py-2">{{ $t('actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="role in roles" :key="role.id" class="border-b last:border-0">
            <td class="py-2">
              {{ role.name }}
              <span v-if="role.is_system" class="ml-1 text-xs text-amber-600">({{ $t('system') }})</span>
            </td>
            <td class="py-2 text-gray-600">{{ role.description || '—' }}</td>
            <td class="py-2">{{ role.permission_ids.length }}</td>
            <td class="py-2">{{ role.user_count }}</td>
            <td class="py-2">
              <PermissionGuard resource="rbac" action="manage">
                <button type="button" class="text-blue-600 hover:underline" @click="openEdit(role)">
                  {{ $t('editPermissions') }}
                </button>
                <template #fallback>—</template>
              </PermissionGuard>
            </td>
          </tr>
        </tbody>
      </table>
    </section>

    <div
      v-if="modalOpen"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/30 p-4"
      @click.self="closeModal"
    >
      <div class="w-full max-w-lg rounded-lg bg-white p-6 shadow-lg">
        <h3 class="mb-4 font-medium">
          {{ editingRole ? $t('editPermissions') : $t('createRole') }}
        </h3>
        <form class="space-y-3" @submit.prevent="save">
          <div v-if="!editingRole">
            <label class="mb-1 block text-sm text-gray-600">{{ $t('name') }}</label>
            <input v-model="formName" required class="w-full rounded border px-3 py-2 text-sm">
          </div>
          <div v-if="!editingRole">
            <label class="mb-1 block text-sm text-gray-600">{{ $t('description') }}</label>
            <input v-model="formDesc" class="w-full rounded border px-3 py-2 text-sm">
          </div>
          <div>
            <p class="mb-2 text-sm text-gray-600">{{ $t('permissions') }}</p>
            <div class="max-h-48 space-y-2 overflow-y-auto rounded border p-3">
              <label
                v-for="perm in permissionItems"
                :key="perm.id"
                class="flex items-center gap-2 text-sm"
              >
                <input v-model="selectedPermIds" type="checkbox" :value="perm.id">
                <span>{{ perm.resource }}:{{ perm.action }}</span>
              </label>
            </div>
          </div>
          <p v-if="modalError" class="text-sm text-red-600">{{ modalError }}</p>
          <div class="flex justify-end gap-2 pt-2">
            <button type="button" class="rounded border px-3 py-1.5 text-sm" @click="closeModal">
              {{ $t('cancel') }}
            </button>
            <button type="submit" class="rounded bg-blue-600 px-3 py-1.5 text-sm text-white" :disabled="saving">
              {{ saving ? $t('loading') : $t('save') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { PermissionDictItem, RoleItem } from '~/composables/use-rbac'

definePageMeta({ middleware: 'auth' })

const { t } = useI18n()
const rbac = useRbac()
const { can } = usePermission()

const roles = ref<RoleItem[]>([])
const permissionItems = ref<PermissionDictItem[]>([])
const pending = ref(true)
const error = ref('')
const modalOpen = ref(false)
const editingRole = ref<RoleItem | null>(null)
const formName = ref('')
const formDesc = ref('')
const selectedPermIds = ref<string[]>([])
const saving = ref(false)
const modalError = ref('')

onMounted(async () => {
  if (!can('rbac', 'view') && !can('rbac', 'manage')) {
    error.value = t('noPermission')
    pending.value = false
    return
  }
  await load()
})

async function load() {
  pending.value = true
  error.value = ''
  try {
    ;[roles.value, permissionItems.value] = await Promise.all([
      rbac.fetchRoles(),
      rbac.fetchPermissionItems(),
    ])
  } catch (e) {
    error.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

function openCreate() {
  editingRole.value = null
  formName.value = ''
  formDesc.value = ''
  selectedPermIds.value = []
  modalOpen.value = true
}

function openEdit(role: RoleItem) {
  editingRole.value = role
  selectedPermIds.value = [...role.permission_ids]
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
  editingRole.value = null
  modalError.value = ''
}

async function save() {
  saving.value = true
  modalError.value = ''
  try {
    if (editingRole.value) {
      await rbac.assignRolePermissions(editingRole.value.id, selectedPermIds.value)
    } else {
      const created = await rbac.createRole({ name: formName.value, description: formDesc.value })
      if (selectedPermIds.value.length > 0) {
        await rbac.assignRolePermissions(created.id, selectedPermIds.value)
      }
    }
    closeModal()
    await load()
    await rbac.loadMyPermissions()
  } catch (e) {
    modalError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    saving.value = false
  }
}
</script>
