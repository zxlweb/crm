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
    class="w-full"
    v-bind="$attrs"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue?: string
    type?: string
    placeholder?: string
    disabled?: boolean
    /** 列表工具栏搜索：自动加放大镜图标 */
    search?: boolean
    leadingIcon?: string
  }>(),
  {
    modelValue: '',
    type: 'text',
    placeholder: '',
    disabled: false,
    search: false,
    leadingIcon: '',
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
</script>
