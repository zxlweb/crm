<template>
  <UTabs
    v-model="activeIndex"
    :items="tabItems"
    class="w-full"
  />
</template>

<script setup lang="ts">
export type TabItem = { id: string; label: string }

const props = defineProps<{
  modelValue: string
  items: TabItem[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

/** Nuxt UI UTabs v-model 为数字索引；对外仍用 string id */
const activeIndex = computed({
  get() {
    const idx = props.items.findIndex((i) => i.id === props.modelValue)
    return idx >= 0 ? idx : 0
  },
  set(idx: number) {
    const item = props.items[idx]
    if (item) emit('update:modelValue', item.id)
  },
})

const tabItems = computed(() =>
  props.items.map((item) => ({
    label: item.label,
  })),
)
</script>
