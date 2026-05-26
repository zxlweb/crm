<template>
  <UInput
    :model-value="modelValue"
    :type="type"
    :placeholder="placeholder"
    :disabled="disabled"
    :leading="hasLeadingIcon"
    :leading-icon="resolvedLeadingIcon"
    color="gray"
    size="md"
    variant="outline"
    :class="widthClass"
    v-bind="$attrs"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    type?: string
    placeholder?: string
    disabled?: boolean
    /** 列表工具栏搜索：自动加放大镜图标 */
    search?: boolean
    leadingIcon?: string
    /**
     * 默认 `auto`：在 block 容器（如 form field <div>）中铺满父级宽度，
     * 在 flex 容器中按内容/min-w 控制宽度，避免行内组件被强行拉满。
     * 设为 `full` 强制 w-full；设为 `inline` 仅按内容宽度。
     */
    width?: 'auto' | 'full' | 'inline'
  }>(),
  {
    modelValue: '',
    type: 'text',
    placeholder: '',
    disabled: false,
    search: false,
    leadingIcon: '',
    width: 'auto',
  },
)

defineEmits<{
  'update:modelValue': [value: string]
}>()

const resolvedLeadingIcon = computed(() => {
  if (props.leadingIcon) return props.leadingIcon
  if (props.search) return 'i-heroicons-magnifying-glass-20-solid'
  return undefined
})

const hasLeadingIcon = computed(() => !!resolvedLeadingIcon.value)

const widthClass = computed(() => {
  if (props.width === 'full') return 'w-full'
  if (props.width === 'inline') return ''
  return 'w-full max-w-md'
})
</script>
