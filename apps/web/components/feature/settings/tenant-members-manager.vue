<template>
  <div class="space-y-4" data-testid="tenant-members-manager">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="relative min-w-0 flex-1 sm:max-w-sm">
        <UIcon
          name="i-heroicons-magnifying-glass"
          class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ds-fg-muted"
          aria-hidden="true"
        />
        <input
          v-model="search"
          type="search"
          class="ds-input w-full rounded-xl py-2 pl-9 pr-3 text-sm"
          :placeholder="$t('membersSearchPh')"
          data-testid="members-search"
        >
      </div>
      <p class="text-xs text-ds-fg-muted">
        {{ $t('membersTotal', { n: filteredMembers.length }) }}
      </p>
    </div>

    <p v-if="loading" class="py-12 text-center text-sm text-ds-fg-muted">{{ $t('loading') }}</p>
    <p v-else-if="loadError" class="rounded-xl border border-ds-danger/20 bg-ds-danger-subtle px-4 py-3 text-sm text-ds-danger">
      {{ loadError }}
    </p>
    <p v-else-if="filteredMembers.length === 0" class="py-12 text-center text-sm text-ds-fg-muted">
      {{ search ? $t('membersNoMatch') : $t('membersEmpty') }}
    </p>

    <div v-else class="overflow-x-auto rounded-xl border border-ds-border-muted">
      <table class="w-full min-w-[560px] text-left text-sm">
        <thead>
          <tr class="border-b border-ds-border-muted bg-ds-bg-muted text-xs font-medium uppercase tracking-wide text-ds-fg-muted">
            <th class="px-4 py-3">{{ $t('membersColUser') }}</th>
            <th class="px-4 py-3">{{ $t('membersColDepartment') }}</th>
            <th class="px-4 py-3">{{ $t('membersColRoles') }}</th>
            <th class="px-4 py-3">{{ $t('membersColJoined') }}</th>
            <th v-if="canManage" class="px-4 py-3 text-right">{{ $t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-ds-border-muted">
          <tr
            v-for="member in filteredMembers"
            :key="member.id"
            class="transition-colors hover:bg-ds-bg-muted/60"
            :data-testid="`member-row-${member.id}`"
          >
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div
                  class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-ds-brand-subtle text-xs font-bold text-ds-fg-brand"
                  aria-hidden="true"
                >
                  {{ memberInitial(member) }}
                </div>
                <div class="min-w-0">
                  <p class="truncate font-medium text-ds-fg-heading">{{ member.name || '—' }}</p>
                  <p class="truncate text-xs text-ds-fg-muted">{{ member.email }}</p>
                </div>
              </div>
            </td>
            <td class="px-4 py-3 text-sm text-ds-fg-heading">
              {{ member.department || '—' }}
            </td>
            <td class="px-4 py-3">
              <div v-if="member.roles.length" class="flex flex-wrap gap-1.5">
                <span
                  v-for="role in member.roles"
                  :key="role.id"
                  class="inline-flex items-center rounded-full border border-ds-border-muted bg-ds-bg-muted px-2 py-0.5 text-xs font-medium text-ds-fg-heading"
                >
                  {{ role.name }}
                </span>
              </div>
              <span v-else class="text-xs text-ds-fg-muted">{{ $t('membersNoRoles') }}</span>
            </td>
            <td class="px-4 py-3 text-xs text-ds-fg-muted whitespace-nowrap">
              {{ formatJoined(member.joined_at) }}
            </td>
            <td v-if="canManage" class="px-4 py-3 text-right">
              <button
                type="button"
                class="cursor-pointer rounded-lg border border-ds-border px-3 py-1.5 text-xs font-medium text-ds-fg-brand transition-colors hover:border-ds-brand hover:bg-ds-brand-subtle"
                :data-testid="`member-edit-roles-${member.id}`"
                @click="openEditRoles(member)"
              >
                {{ $t('membersEditRoles') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <UiModal
      v-model:open="rolesDialogOpen"
      :title="$t('membersEditRolesTitle')"
      :subtitle="editRolesSubtitle"
      data-testid="members-roles-dialog"
    >
      <ul class="max-h-72 space-y-2 overflow-y-auto pr-0.5">
        <li v-for="role in allRoles" :key="role.id">
          <label
            class="flex cursor-pointer items-start gap-3 rounded-xl border px-3 py-2.5 transition-colors duration-200"
            :class="draftRoleIds.includes(role.id)
              ? 'border-ds-brand/40 bg-ds-brand-subtle/50 shadow-ds-sm'
              : 'border-ds-border-muted bg-ds-bg-muted/30 hover:border-ds-border hover:bg-ds-bg-muted'"
            :data-testid="`member-role-option-${role.id}`"
          >
            <input
              v-model="draftRoleIds"
              type="checkbox"
              class="ds-checkbox mt-0.5 h-4 w-4 shrink-0 cursor-pointer rounded"
              :value="role.id"
            >
            <span class="min-w-0 flex-1">
              <span class="flex flex-wrap items-center gap-2">
                <span class="text-sm font-medium text-ds-fg-heading">{{ role.name }}</span>
                <span
                  v-if="role.is_system"
                  class="rounded-full bg-amber-100 px-1.5 py-0.5 text-[10px] font-semibold text-amber-700 ring-1 ring-inset ring-amber-200 dark:bg-amber-500/15 dark:text-amber-300 dark:ring-amber-500/30"
                >
                  {{ $t('system') }}
                </span>
              </span>
              <span v-if="role.description" class="mt-0.5 block text-xs leading-relaxed text-ds-fg-muted">
                {{ role.description }}
              </span>
            </span>
          </label>
        </li>
      </ul>
      <p
        v-if="saveError"
        class="mt-4 flex items-center gap-2 rounded-lg border border-ds-danger/20 bg-ds-danger-subtle px-3 py-2 text-sm text-ds-danger"
      >
        <UIcon name="i-heroicons-exclamation-circle" class="h-4 w-4 shrink-0" aria-hidden="true" />
        {{ saveError }}
      </p>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="secondary" :disabled="saving" @click="rolesDialogOpen = false">
            {{ $t('cancel') }}
          </UiButton>
          <UiButton
            :loading="saving"
            :disabled="draftRoleIds.length === 0"
            data-testid="members-roles-save"
            @click="saveRoles"
          >
            {{ $t('save') }}
          </UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>

<script setup lang="ts">
import type { MemberItem, RoleItem } from '~/composables/use-rbac'

const { t, locale } = useI18n()
const { can } = usePermission()
const rbac = useRbac()

const canManage = computed(() => can('rbac', 'manage'))

const loading = ref(true)
const loadError = ref('')
const members = ref<MemberItem[]>([])
const allRoles = ref<RoleItem[]>([])
const search = ref('')

const rolesDialogOpen = ref(false)
const editingMember = ref<MemberItem | null>(null)
const draftRoleIds = ref<string[]>([])
const saving = ref(false)
const saveError = ref('')

const editRolesSubtitle = computed(() => {
  if (!editingMember.value) return ''
  const m = editingMember.value
  return m.name ? `${m.name} · ${m.email}` : m.email
})

const filteredMembers = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return members.value
  return members.value.filter((m) => {
    const hay = `${m.name} ${m.email} ${m.department ?? ''} ${m.roles.map((r) => r.name).join(' ')}`.toLowerCase()
    return hay.includes(q)
  })
})

function memberInitial(member: MemberItem): string {
  const s = (member.name || member.email || '?').trim()
  return s.charAt(0).toUpperCase()
}

function formatJoined(iso: string): string {
  if (!iso) return '—'
  try {
    const d = new Date(iso)
    if (Number.isNaN(d.getTime())) return iso
    return new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium' }).format(d)
  } catch {
    return iso
  }
}

async function load() {
  loading.value = true
  loadError.value = ''
  try {
    const memberRows = await rbac.fetchMembers()
    members.value = memberRows
    if (canManage.value) {
      allRoles.value = await rbac.fetchRoles()
    }
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    loading.value = false
  }
}

function openEditRoles(member: MemberItem) {
  editingMember.value = member
  draftRoleIds.value = member.roles.map((r) => r.id)
  saveError.value = ''
  rolesDialogOpen.value = true
}

function resetRolesDialog() {
  editingMember.value = null
  draftRoleIds.value = []
  saveError.value = ''
}

watch(rolesDialogOpen, (open) => {
  if (!open) resetRolesDialog()
})

async function saveRoles() {
  if (!editingMember.value) return
  saving.value = true
  saveError.value = ''
  try {
    await rbac.assignMemberRoles(editingMember.value.id, draftRoleIds.value)
    await load()
    rolesDialogOpen.value = false
  } catch (e) {
    saveError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
