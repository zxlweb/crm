<template>
  <div
    class="overflow-hidden rounded-2xl border border-ds-border-table bg-ds-bg-elevated shadow-ds-sm"
    data-testid="data-table"
  >
    <div
      v-if="$slots.toolbar"
      class="border-b border-ds-divide px-5 py-3"
    >
      <slot name="toolbar" />
    </div>
    <UTable
      v-bind="tableAttrs"
      :rows="rows"
      :columns="columns"
      :ui="mergedUi"
    >
      <template v-for="(_, slotName) in tableBodySlots" #[slotName]="slotProps">
        <slot :name="slotName" v-bind="slotProps ?? {}" />
      </template>
    </UTable>

    <div
      v-if="$slots.footer"
      class="flex flex-wrap items-center justify-between gap-3 border-t border-ds-divide px-5 py-3.5 text-xs text-ds-fg-muted"
    >
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, useAttrs, useSlots } from 'vue'
import { crmTableUi } from '../../theme/nuxt-ui-table'
import type { UiTableColumn } from './table-types'

export type { UiTableColumn }

const props = defineProps<{
  rows: Record<string, unknown>[]
  columns: UiTableColumn[]
  /** 透传 UTable 的 ui，与 CRM 默认主题 merge */
  ui?: Record<string, unknown>
}>()

defineOptions({ inheritAttrs: false })

const slots = useSlots()

/** toolbar / footer 由 UiTable 外壳渲染，不传给 UTable */
const tableBodySlots = computed(() => {
  const reserved = new Set(['toolbar', 'footer'])
  return Object.fromEntries(
    Object.entries(slots).filter(([name]) => !reserved.has(name)),
  )
})

const attrs = useAttrs()

const tableAttrs = computed(() => {
  const { rows: _r, columns: _c, ui: _u, ...rest } = attrs
  return rest
})

const mergedUi = computed(() => deepMerge(crmTableUi as Record<string, unknown>, props.ui ?? {}))

function deepMerge(base: Record<string, unknown>, override: Record<string, unknown>): Record<string, unknown> {
  const out = { ...base }
  for (const key of Object.keys(override)) {
    const b = base[key]
    const o = override[key]
    if (b && o && typeof b === 'object' && typeof o === 'object' && !Array.isArray(b)) {
      out[key] = deepMerge(b as Record<string, unknown>, o as Record<string, unknown>)
    } else {
      out[key] = o
    }
  }
  return out
}
</script>
