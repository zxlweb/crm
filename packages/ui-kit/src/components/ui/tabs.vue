<template>
  <div
    class="ds-ui-tabs inline-flex max-w-full items-center gap-0.5 rounded-xl border border-ds-border-muted bg-ds-bg-muted/40 p-0.5 backdrop-blur-sm"
    role="tablist"
    :aria-label="ariaLabel"
  >
    <button
      v-for="(item, index) in items"
      :key="item.id"
      ref="tabRefs"
      type="button"
      role="tab"
      :aria-selected="item.id === modelValue"
      :tabindex="item.id === modelValue ? 0 : -1"
      :data-active="item.id === modelValue ? 'true' : undefined"
      class="ds-ui-tabs__tab group relative inline-flex min-w-0 cursor-pointer items-center gap-1.5 rounded-lg px-2.5 py-1 text-xs font-medium transition-[color,background-color,box-shadow] duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ds-brand focus-visible:ring-offset-1 focus-visible:ring-offset-ds-bg-muted"
      :class="
        item.id === modelValue
          ? 'bg-ds-bg-elevated text-ds-fg-heading shadow-ds-sm ring-1 ring-inset ring-ds-border-muted'
          : 'text-ds-fg-muted hover:text-ds-fg-heading'
      "
      @click="select(item.id)"
      @keydown="onKeydown($event, index)"
    >
      <UIcon
        v-if="item.icon"
        :name="item.icon"
        class="h-3.5 w-3.5 shrink-0 transition-colors duration-200"
        :class="item.id === modelValue ? 'text-ds-fg-brand' : 'text-ds-fg-subtle group-hover:text-ds-fg-muted'"
        aria-hidden="true"
      />
      <span class="truncate">{{ item.label }}</span>
      <span
        v-if="item.count != null"
        class="ml-0.5 inline-flex min-w-[1.25rem] items-center justify-center rounded-full px-1.5 text-[10px] font-semibold tabular-nums leading-[1.1rem] ring-1 ring-inset transition-colors duration-200"
        :class="
          item.id === modelValue
            ? 'bg-ds-brand-subtle text-ds-fg-brand ring-ds-brand/20'
            : 'bg-ds-bg-muted text-ds-fg-muted ring-ds-border-muted'
        "
      >
        {{ item.count }}
      </span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { nextTick, ref } from 'vue'

export type TabItem = {
  id: string
  label: string
  icon?: string
  count?: number
}

const props = defineProps<{
  modelValue: string
  items: TabItem[]
  ariaLabel?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const tabRefs = ref<HTMLButtonElement[]>([])

function select(id: string) {
  if (id !== props.modelValue) emit('update:modelValue', id)
}

function onKeydown(event: KeyboardEvent, index: number) {
  const last = props.items.length - 1
  let next = -1
  switch (event.key) {
    case 'ArrowRight':
      next = index === last ? 0 : index + 1
      break
    case 'ArrowLeft':
      next = index === 0 ? last : index - 1
      break
    case 'Home':
      next = 0
      break
    case 'End':
      next = last
      break
    default:
      return
  }
  event.preventDefault()
  const target = props.items[next]
  if (target) {
    emit('update:modelValue', target.id)
    nextTick(() => {
      tabRefs.value[next]?.focus()
    })
  }
}
</script>
