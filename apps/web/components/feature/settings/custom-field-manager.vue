<template>
  <div data-testid="custom-field-manager" class="space-y-5">
    <!-- Toolbar -->
    <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
      <!-- Entity chip filter -->
      <div class="flex flex-wrap items-center gap-1.5" role="tablist" :aria-label="$t('cfFilterAria')">
        <button
          v-for="chip in entityChips"
          :key="chip.value || 'all'"
          type="button"
          role="tab"
          :aria-selected="selectedEntity === chip.value"
          class="inline-flex cursor-pointer items-center gap-2 rounded-full border px-3 py-1.5 text-xs font-medium transition-colors duration-200"
          :class="selectedEntity === chip.value
            ? 'border-ds-brand bg-ds-brand-subtle text-ds-fg-brand shadow-ds-sm'
            : 'border-ds-border-muted bg-ds-bg-muted text-ds-fg-muted hover:border-ds-border hover:text-ds-fg-heading'"
          :data-testid="`cf-chip-${chip.value || 'all'}`"
          @click="selectedEntity = chip.value"
        >
          <UIcon :name="chip.icon" class="h-3.5 w-3.5" aria-hidden="true" />
          <span>{{ chip.label }}</span>
          <span
            class="rounded-full bg-ds-bg-elevated px-1.5 py-0.5 text-[10px] font-semibold leading-none ring-1 ring-inset ring-ds-border-muted"
            :class="selectedEntity === chip.value ? 'text-ds-fg-brand' : 'text-ds-fg-muted'"
          >
            {{ chip.count }}
          </span>
        </button>
      </div>

      <div class="flex items-center gap-2">
        <!-- Hidden mirror to preserve original test selector contract -->
        <select
          v-model="selectedEntity"
          data-testid="cf-entity-filter"
          class="sr-only"
          tabindex="-1"
          aria-hidden="true"
        >
          <option value="">{{ $t('cfAllEntities') }}</option>
          <option v-for="e in entityTypes" :key="e" :value="e">{{ $t(`cfEntity.${e}`) }}</option>
        </select>

        <div class="relative">
          <UIcon
            name="i-heroicons-magnifying-glass"
            class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ds-fg-muted"
            aria-hidden="true"
          />
          <input
            v-model="searchQuery"
            type="search"
            class="ds-input w-full rounded-xl py-2 pl-9 pr-3 text-sm shadow-ds-sm lg:w-64"
            :placeholder="$t('cfSearchPh')"
            data-testid="cf-search"
          >
        </div>

        <button
          v-if="canEdit"
          type="button"
          class="ds-btn-primary inline-flex cursor-pointer items-center gap-2 rounded-xl px-4 py-2 text-sm font-semibold transition-colors duration-200"
          data-testid="cf-create-btn"
          @click="openCreate"
        >
          <UIcon name="i-heroicons-plus" class="h-4 w-4" aria-hidden="true" />
          <span>{{ $t('cfCreate') }}</span>
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <div class="h-6 w-6 animate-spin rounded-full border-2 border-ds-brand-muted border-t-ds-brand" />
    </div>

    <!-- Empty state -->
    <div
      v-else-if="filteredFields.length === 0"
      class="flex flex-col items-center justify-center gap-3 rounded-2xl border border-dashed border-ds-border-muted bg-ds-bg-muted/40 px-6 py-14 text-center"
    >
      <span class="flex h-12 w-12 items-center justify-center rounded-2xl bg-ds-brand-subtle text-ds-fg-brand" aria-hidden="true">
        <UIcon name="i-heroicons-square-3-stack-3d" class="h-6 w-6" />
      </span>
      <div class="space-y-1">
        <p class="text-sm font-semibold text-ds-fg-heading">
          {{ searchQuery || selectedEntity ? $t('cfNoMatch') : $t('cfEmpty') }}
        </p>
        <p class="text-xs text-ds-fg-muted">
          {{ searchQuery || selectedEntity ? $t('cfNoMatchHint') : $t('cfEmptyHint') }}
        </p>
      </div>
      <button
        v-if="canEdit && !searchQuery && !selectedEntity"
        type="button"
        class="ds-btn-primary mt-2 inline-flex cursor-pointer items-center gap-2 rounded-xl px-4 py-2 text-sm font-medium"
        @click="openCreate"
      >
        <UIcon name="i-heroicons-plus" class="h-4 w-4" aria-hidden="true" />
        <span>{{ $t('cfCreate') }}</span>
      </button>
    </div>

    <!-- Table -->
    <div v-else class="overflow-hidden rounded-xl border border-ds-border-muted">
      <div class="overflow-x-auto">
        <table class="w-full min-w-[720px] text-left text-sm">
          <thead class="bg-ds-bg-muted/60">
            <tr class="text-xs font-semibold uppercase tracking-wide text-ds-fg-muted">
              <th class="px-4 py-3 font-semibold">{{ $t('cfColKey') }}</th>
              <th class="px-4 py-3 font-semibold">{{ $t('cfColLabel') }}</th>
              <th class="px-4 py-3 font-semibold">{{ $t('cfColEntity') }}</th>
              <th class="px-4 py-3 font-semibold">{{ $t('cfColType') }}</th>
              <th class="px-4 py-3 font-semibold">{{ $t('cfColRequired') }}</th>
              <th v-if="canEdit" class="px-4 py-3 text-right font-semibold">{{ $t('actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-ds-border-muted">
            <tr
              v-for="field in filteredFields"
              :key="field.id"
              class="transition-colors duration-200 hover:bg-ds-bg-muted/60"
            >
              <td class="px-4 py-3">
                <code class="rounded-md bg-ds-bg-muted px-1.5 py-0.5 font-mono text-xs text-ds-fg-heading ring-1 ring-inset ring-ds-border-muted">
                  {{ field.field_key }}
                </code>
              </td>
              <td class="px-4 py-3 font-medium text-ds-fg-heading">{{ localizedLabel(field) }}</td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center gap-1.5 rounded-full bg-ds-brand-subtle px-2 py-0.5 text-xs font-medium text-ds-fg-brand ring-1 ring-inset ring-ds-brand/20"
                >
                  <UIcon :name="entityIcon(field.entity_type)" class="h-3 w-3" aria-hidden="true" />
                  {{ $t(`cfEntity.${field.entity_type}`) }}
                </span>
              </td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-xs font-medium ring-1 ring-inset"
                  :class="typeBadgeClass(field.field_type)"
                >
                  <UIcon :name="typeIcon(field.field_type)" class="h-3 w-3" aria-hidden="true" />
                  {{ $t(`cfType.${field.field_type}`) }}
                </span>
              </td>
              <td class="px-4 py-3">
                <span v-if="field.required" class="inline-flex items-center gap-1.5 text-xs font-medium text-ds-danger">
                  <span class="inline-block h-1.5 w-1.5 rounded-full bg-ds-danger" aria-hidden="true" />
                  {{ $t('cfYes') }}
                </span>
                <span v-else class="text-xs text-ds-fg-muted">{{ $t('cfNo') }}</span>
              </td>
              <td v-if="canEdit" class="px-4 py-3 text-right">
                <div class="inline-flex items-center gap-1">
                  <button
                    type="button"
                    class="inline-flex cursor-pointer items-center gap-1 rounded-lg border border-ds-border-muted bg-ds-bg-elevated px-2.5 py-1.5 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-brand hover:text-ds-fg-brand"
                    :title="$t('edit')"
                    @click="openEdit(field)"
                  >
                    <UIcon name="i-heroicons-pencil-square" class="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{{ $t('edit') }}</span>
                  </button>
                  <button
                    type="button"
                    class="inline-flex cursor-pointer items-center gap-1 rounded-lg border border-ds-border-muted bg-ds-bg-elevated px-2.5 py-1.5 text-xs font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-danger hover:text-ds-danger"
                    :title="$t('cfDelete')"
                    @click="handleRemove(field)"
                  >
                    <UIcon name="i-heroicons-trash" class="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{{ $t('cfDelete') }}</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <div
      v-if="modalOpen"
      class="fixed inset-0 z-ds-modal flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
      @click.self="closeModal"
    >
      <div
        class="ds-card flex w-full max-w-xl flex-col overflow-hidden rounded-2xl shadow-ds-xl"
        data-testid="cf-modal"
      >
        <header class="flex items-start gap-3 border-b border-ds-border-muted bg-ds-bg-muted/40 px-6 py-4">
          <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-ds-brand-subtle text-ds-fg-brand" aria-hidden="true">
            <UIcon
              :name="editing ? 'i-heroicons-pencil-square' : 'i-heroicons-plus-circle'"
              class="h-5 w-5"
            />
          </span>
          <div class="min-w-0 flex-1">
            <h3 class="text-base font-semibold text-ds-fg-heading">
              {{ editing ? $t('cfEditTitle') : $t('cfCreateTitle') }}
            </h3>
            <p class="mt-0.5 text-xs text-ds-fg-muted">{{ $t('cfModalHint') }}</p>
          </div>
          <button
            type="button"
            class="inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-ds-fg-muted transition-colors duration-200 hover:bg-ds-bg-muted hover:text-ds-fg-heading"
            :aria-label="$t('cancel')"
            @click="closeModal"
          >
            <UIcon name="i-heroicons-x-mark" class="h-5 w-5" aria-hidden="true" />
          </button>
        </header>

        <form class="space-y-5 px-6 py-5" @submit.prevent="handleSubmit">
          <div class="grid gap-4 sm:grid-cols-2">
            <div class="space-y-1.5">
              <label class="block text-xs font-semibold uppercase tracking-wide text-ds-fg-muted">{{ $t('cfColEntity') }}</label>
              <div class="relative">
                <select
                  v-model="formEntity"
                  class="ds-input w-full appearance-none rounded-xl px-3 py-2 pr-9 text-sm"
                  :disabled="!!editing"
                  data-testid="cf-form-entity"
                >
                  <option v-for="e in entityTypes" :key="e" :value="e">{{ $t(`cfEntity.${e}`) }}</option>
                </select>
                <UIcon
                  name="i-heroicons-chevron-down"
                  class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ds-fg-muted"
                  aria-hidden="true"
                />
              </div>
            </div>
            <div class="space-y-1.5">
              <label class="block text-xs font-semibold uppercase tracking-wide text-ds-fg-muted">{{ $t('cfColKey') }}</label>
              <input
                v-model="formKey"
                class="ds-input w-full rounded-xl px-3 py-2 font-mono text-sm"
                :disabled="!!editing"
                placeholder="custom_field_key"
                data-testid="cf-form-key"
              >
              <p v-if="!editing" class="text-[11px] text-ds-fg-muted">{{ $t('cfKeyHelp') }}</p>
            </div>
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <div class="space-y-1.5">
              <label class="block text-xs font-semibold uppercase tracking-wide text-ds-fg-muted">{{ $t('cfLabelZh') }}</label>
              <input v-model="formLabelZh" class="ds-input w-full rounded-xl px-3 py-2 text-sm" placeholder="例如：行业" data-testid="cf-form-label-zh">
            </div>
            <div class="space-y-1.5">
              <label class="block text-xs font-semibold uppercase tracking-wide text-ds-fg-muted">{{ $t('cfLabelEn') }}</label>
              <input v-model="formLabelEn" class="ds-input w-full rounded-xl px-3 py-2 text-sm" placeholder="e.g. Industry" data-testid="cf-form-label-en">
            </div>
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <div class="space-y-1.5">
              <label class="block text-xs font-semibold uppercase tracking-wide text-ds-fg-muted">{{ $t('cfColType') }}</label>
              <div class="relative">
                <select
                  v-model="formType"
                  class="ds-input w-full appearance-none rounded-xl px-3 py-2 pr-9 text-sm"
                  :disabled="!!editing"
                  data-testid="cf-form-type"
                >
                  <option value="text">{{ $t('cfType.text') }}</option>
                  <option value="select">{{ $t('cfType.select') }}</option>
                  <option value="date">{{ $t('cfType.date') }}</option>
                </select>
                <UIcon
                  name="i-heroicons-chevron-down"
                  class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ds-fg-muted"
                  aria-hidden="true"
                />
              </div>
            </div>
            <div class="flex items-end">
              <label class="inline-flex cursor-pointer items-center gap-2 rounded-xl border border-ds-border-muted bg-ds-bg-muted/40 px-3 py-2 text-sm">
                <input v-model="formRequired" type="checkbox" class="ds-checkbox h-4 w-4 cursor-pointer rounded">
                <span class="text-ds-fg-heading">{{ $t('cfColRequired') }}</span>
              </label>
            </div>
          </div>

          <div
            v-if="formType === 'select'"
            class="space-y-3 rounded-xl border border-ds-border-muted bg-ds-bg-muted/30 p-4"
          >
            <div class="flex items-center justify-between">
              <label class="block text-xs font-semibold uppercase tracking-wide text-ds-fg-muted">{{ $t('cfOptions') }}</label>
              <span class="text-[11px] text-ds-fg-muted">{{ formOptions.length }} {{ $t('cfOptionsCountSuffix') }}</span>
            </div>
            <div v-if="formOptions.length === 0" class="rounded-lg border border-dashed border-ds-border-muted bg-ds-bg px-4 py-3 text-center text-xs text-ds-fg-muted">
              {{ $t('cfOptionsEmpty') }}
            </div>
            <div v-for="(opt, idx) in formOptions" :key="idx" class="grid grid-cols-[1fr_1fr_1fr_auto] items-center gap-2">
              <input v-model="opt.value" class="ds-input rounded-lg px-2.5 py-1.5 font-mono text-xs" placeholder="value">
              <input v-model="opt.label['zh-CN']" class="ds-input rounded-lg px-2.5 py-1.5 text-xs" placeholder="中文">
              <input v-model="opt.label['en-US']" class="ds-input rounded-lg px-2.5 py-1.5 text-xs" placeholder="English">
              <button
                type="button"
                class="inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-ds-fg-muted transition-colors duration-200 hover:bg-ds-danger-subtle hover:text-ds-danger"
                :aria-label="$t('cfDelete')"
                @click="formOptions.splice(idx, 1)"
              >
                <UIcon name="i-heroicons-trash" class="h-4 w-4" aria-hidden="true" />
              </button>
            </div>
            <button
              type="button"
              class="inline-flex cursor-pointer items-center gap-1 text-xs font-medium text-ds-fg-brand hover:underline"
              @click="addOption"
            >
              <UIcon name="i-heroicons-plus" class="h-3.5 w-3.5" aria-hidden="true" />
              {{ $t('cfAddOption') }}
            </button>
          </div>

          <p
            v-if="modalError"
            class="flex items-center gap-2 rounded-lg border border-ds-danger/20 bg-ds-danger-subtle px-3 py-2 text-sm text-ds-danger"
          >
            <UIcon name="i-heroicons-exclamation-circle" class="h-4 w-4 shrink-0" aria-hidden="true" />
            {{ modalError }}
          </p>

          <footer class="-mx-6 flex items-center justify-end gap-2 border-t border-ds-border-muted bg-ds-bg-muted/30 px-6 pt-4">
            <button
              type="button"
              class="cursor-pointer rounded-xl border border-ds-border bg-ds-bg-elevated px-4 py-2 text-sm font-medium text-ds-fg-muted transition-colors duration-200 hover:border-ds-border hover:text-ds-fg-heading"
              @click="closeModal"
            >
              {{ $t('cancel') }}
            </button>
            <button
              type="submit"
              class="ds-btn-primary inline-flex cursor-pointer items-center gap-2 rounded-xl px-5 py-2 text-sm font-semibold disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="submitting"
              data-testid="cf-form-submit"
            >
              <UIcon
                v-if="submitting"
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
              <span>{{ submitting ? $t('loading') : $t('save') }}</span>
            </button>
          </footer>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { CustomField, CustomFieldOption, EntityType, FieldType } from '~/types/settings'

defineProps<{ canEdit: boolean }>()
const emit = defineEmits<{ refresh: [] }>()

const { t, locale } = useI18n()
const customFields = useCustomFields()

const fields = ref<CustomField[]>([])
const loading = ref(true)
const selectedEntity = ref<EntityType | ''>('')
const searchQuery = ref('')
const entityTypes: EntityType[] = ['account', 'contact', 'lead', 'deal']

const modalOpen = ref(false)
const editing = ref<CustomField | null>(null)
const formEntity = ref<EntityType>('lead')
const formKey = ref('')
const formLabelZh = ref('')
const formLabelEn = ref('')
const formType = ref<FieldType>('text')
const formRequired = ref(false)
const formOptions = ref<CustomFieldOption[]>([])
const submitting = ref(false)
const modalError = ref('')

const entityIconMap: Record<EntityType, string> = {
  account: 'i-heroicons-building-office',
  contact: 'i-heroicons-user-circle',
  lead: 'i-heroicons-sparkles',
  deal: 'i-heroicons-banknotes',
}

function entityIcon(e: EntityType) {
  return entityIconMap[e]
}

const typeIconMap: Record<FieldType, string> = {
  text: 'i-heroicons-bars-3-bottom-left',
  select: 'i-heroicons-list-bullet',
  date: 'i-heroicons-calendar-days',
}

function typeIcon(t: FieldType) {
  return typeIconMap[t]
}

function typeBadgeClass(ty: FieldType) {
  switch (ty) {
    case 'select':
      return 'bg-ds-brand-subtle text-ds-fg-brand ring-ds-brand/20'
    case 'date':
      return 'bg-ds-info-subtle text-ds-info ring-ds-info/20'
    case 'text':
    default:
      return 'bg-ds-bg-muted text-ds-fg-muted ring-ds-border-muted'
  }
}

const entityChips = computed(() => {
  const counts = new Map<EntityType, number>()
  for (const f of fields.value) counts.set(f.entity_type, (counts.get(f.entity_type) ?? 0) + 1)
  const chips: { value: EntityType | ''; label: string; icon: string; count: number }[] = [
    {
      value: '',
      label: t('cfAllEntities'),
      icon: 'i-heroicons-squares-2x2',
      count: fields.value.length,
    },
  ]
  for (const e of entityTypes) {
    chips.push({
      value: e,
      label: t(`cfEntity.${e}`),
      icon: entityIconMap[e],
      count: counts.get(e) ?? 0,
    })
  }
  return chips
})

const filteredFields = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  return fields.value.filter((f) => {
    if (selectedEntity.value && f.entity_type !== selectedEntity.value) return false
    if (!q) return true
    const label = `${f.field_label['zh-CN']} ${f.field_label['en-US']}`.toLowerCase()
    return f.field_key.toLowerCase().includes(q) || label.includes(q)
  })
})

function localizedLabel(field: CustomField): string {
  const loc = locale.value as 'zh-CN' | 'en-US'
  return field.field_label[loc] || field.field_label['zh-CN'] || field.field_key
}

async function loadFields() {
  loading.value = true
  try {
    fields.value = await customFields.fetchList()
  } finally {
    loading.value = false
  }
}

onMounted(loadFields)

function openCreate() {
  editing.value = null
  formEntity.value = (selectedEntity.value as EntityType) || 'lead'
  formKey.value = ''
  formLabelZh.value = ''
  formLabelEn.value = ''
  formType.value = 'text'
  formRequired.value = false
  formOptions.value = []
  modalError.value = ''
  modalOpen.value = true
}

function openEdit(field: CustomField) {
  editing.value = field
  formEntity.value = field.entity_type
  formKey.value = field.field_key
  formLabelZh.value = field.field_label['zh-CN']
  formLabelEn.value = field.field_label['en-US']
  formType.value = field.field_type
  formRequired.value = field.required
  formOptions.value = field.options.map((o) => ({ ...o, label: { ...o.label } }))
  modalError.value = ''
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
  editing.value = null
}

function addOption() {
  formOptions.value.push({ value: '', label: { 'zh-CN': '', 'en-US': '' } })
}

async function handleSubmit() {
  submitting.value = true
  modalError.value = ''
  try {
    if (editing.value) {
      await customFields.update(editing.value.id, {
        field_label: { 'zh-CN': formLabelZh.value, 'en-US': formLabelEn.value },
        required: formRequired.value,
        options: formType.value === 'select' ? formOptions.value : undefined,
      })
    } else {
      await customFields.create({
        entity_type: formEntity.value,
        field_key: formKey.value,
        field_label: { 'zh-CN': formLabelZh.value, 'en-US': formLabelEn.value },
        field_type: formType.value,
        required: formRequired.value,
        options: formType.value === 'select' ? formOptions.value : undefined,
      })
    }
    closeModal()
    await loadFields()
    emit('refresh')
  } catch (e) {
    modalError.value = e instanceof Error ? e.message : 'Error'
  } finally {
    submitting.value = false
  }
}

async function handleRemove(field: CustomField) {
  try {
    await customFields.remove(field.id)
    await loadFields()
    emit('refresh')
  } catch (e) {
    modalError.value = e instanceof Error ? e.message : 'Error'
  }
}
</script>
