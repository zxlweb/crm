<template>
  <UButton
    :color="mappedColor"
    :variant="mappedVariant"
    :size="size"
    :loading="loading"
    :disabled="disabled"
    class="cursor-pointer transition-transform duration-150 active:scale-[0.98] motion-reduce:transform-none"
    v-bind="$attrs"
  >
    <slot />
  </UButton>
</template>

<script setup lang="ts">
type Variant = 'primary' | 'secondary' | 'ghost' | 'danger'

const props = withDefaults(
  defineProps<{
    variant?: Variant
    size?: 'xs' | 'sm' | 'md' | 'lg'
    loading?: boolean
    disabled?: boolean
  }>(),
  {
    variant: 'primary',
    size: 'md',
    loading: false,
    disabled: false,
  },
)

const mappedColor = computed(() => {
  const map: Record<Variant, string> = {
    primary: 'primary',
    secondary: 'gray',
    ghost: 'gray',
    danger: 'red',
  }
  return map[props.variant]
})

const mappedVariant = computed(() => {
  if (props.variant === 'ghost') return 'ghost'
  if (props.variant === 'secondary') return 'outline'
  return 'solid'
})
</script>
