<template>
  <svg
    :class="sizeClass"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.75"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
    :role="title ? 'img' : undefined"
    :aria-label="title"
  >
    <title v-if="title">{{ title }}</title>
    <component
      :is="shape.tag"
      v-for="(shape, idx) in shapes"
      :key="idx"
      v-bind="shapeAttrs(shape)"
    />
  </svg>
</template>

<script setup lang="ts">
import {
  TAG_ICONS,
  type TagIconName,
  type TagIconPath,
} from '../../icons/tag-icons'

const props = withDefaults(
  defineProps<{
    name: TagIconName
    size?: 'xs' | 'sm' | 'md'
    title?: string
  }>(),
  { size: 'sm' },
)

const shapes = computed(() => TAG_ICONS[props.name] ?? TAG_ICONS.note)

const sizeClass = computed(() => {
  const map = {
    xs: 'h-3.5 w-3.5 shrink-0',
    sm: 'h-4 w-4 shrink-0',
    md: 'h-5 w-5 shrink-0',
  }
  return map[props.size]
})

function shapeAttrs(shape: TagIconPath) {
  const { tag: _tag, fill, ...rest } = shape
  if (fill) return { ...rest, fill }
  if (shape.tag === 'rect' || shape.tag === 'circle') {
    return { ...rest, fill: 'none' }
  }
  return rest
}
</script>
