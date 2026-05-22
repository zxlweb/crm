<template>
  <slot v-if="allowed" />
  <slot v-else name="fallback">
    <p class="rounded-xl bg-ds-danger-subtle px-4 py-3 text-sm text-ds-danger" role="alert">
      {{ $t('noPermission') }}
    </p>
  </slot>
</template>

<script setup lang="ts">
const props = defineProps<{
  resource: string
  action: string
}>()

const { can } = usePermission()
const allowed = computed(() => can(props.resource, props.action))
</script>
