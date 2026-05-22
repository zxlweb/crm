<template>
  <slot v-if="allowed" />
  <slot v-else name="fallback" />
</template>

<script setup lang="ts">
const props = defineProps<{
  resource: string
  action: string
}>()

const { can } = usePermission()
const allowed = computed(() => can(props.resource, props.action))
</script>
