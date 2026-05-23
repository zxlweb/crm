<template>
  <USelect
    :model-value="modelValue"
    :options="options"
    :placeholder="placeholder"
    color="gray"
    size="md"
    variant="outline"
    class="w-full"
    v-bind="$attrs"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup lang="ts">
export type SelectItem = { label: string; value: string }

const props = withDefaults(
  defineProps<{
    modelValue?: string
    items: SelectItem[]
    placeholder?: string
  }>(),
  {
    modelValue: '',
    placeholder: '',
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
</script>
