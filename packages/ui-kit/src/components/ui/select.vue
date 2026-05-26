<template>
  <USelectMenu
    :model-value="modelValue"
    :options="options"
    :placeholder="placeholder"
    :disabled="disabled"
    :searchable="searchable"
    :searchable-placeholder="searchablePlaceholder"
    value-attribute="value"
    option-attribute="label"
    color="gray"
    size="md"
    variant="outline"
    :class="widthClass"
    v-bind="$attrs"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #label>
      <span v-if="selectedLabel" class="truncate text-ds-fg">{{ selectedLabel }}</span>
      <span v-else class="truncate text-ds-fg-subtle">{{ placeholder }}</span>
    </template>
  </USelectMenu>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export type SelectItem = { label: string; value: string }

const props = withDefaults(
  defineProps<{
    modelValue?: string
    items: SelectItem[]
    placeholder?: string
    disabled?: boolean
    searchable?: boolean
    searchablePlaceholder?: string
    /**
     * 默认 `auto`：跟随父容器宽度但有最大宽度上限；`full` 强制铺满；`inline` 仅按内容宽度。
     */
    width?: 'auto' | 'full' | 'inline'
  }>(),
  {
    modelValue: '',
    placeholder: '',
    disabled: false,
    searchable: false,
    searchablePlaceholder: '',
    width: 'auto',
  },
)

defineEmits<{
  'update:modelValue': [value: string]
}>()

const options = computed(() =>
  props.items.map((item) => ({
    label: item.label,
    value: item.value,
  })),
)

const selectedLabel = computed(
  () => props.items.find((item) => item.value === props.modelValue)?.label ?? '',
)

const widthClass = computed(() => {
  if (props.width === 'full') return 'w-full'
  if (props.width === 'inline') return ''
  return 'w-full max-w-md'
})
</script>
