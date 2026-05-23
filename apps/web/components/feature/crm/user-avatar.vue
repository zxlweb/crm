<template>
  <span
    class="inline-flex shrink-0 items-center justify-center rounded-full font-semibold"
    :class="[sizeClass, toneClass]"
    :style="tone === 'accent' ? { background: bgColor } : undefined"
    :title="name"
    role="img"
    :aria-label="name"
  >
    {{ initials }}
  </span>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    name: string
    /** 用于稳定配色 */
    seed?: string
    size?: 'sm' | 'md' | 'lg' | 'xl'
    /** accent：彩色头像；neutral：低对比，用于详情页销售行 */
    tone?: 'accent' | 'neutral'
  }>(),
  { size: 'md', tone: 'accent' },
)

const initials = computed(() => {
  const parts = props.name.trim().split(/\s+/)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  const s = props.name.trim()
  if (s.length <= 2) return s.toUpperCase()
  return s.slice(0, 2).toUpperCase()
})

const sizeClass = computed(() => {
  const map = {
    sm: 'h-8 w-8 text-xs',
    md: 'h-10 w-10 text-sm',
    lg: 'h-12 w-12 text-base',
    xl: 'h-16 w-16 text-lg',
  }
  return map[props.size]
})

const toneClass = computed(() =>
  props.tone === 'neutral'
    ? 'bg-ds-bg-muted text-ds-fg-muted ring-1 ring-ds-border'
    : 'text-ds-on-brand shadow-sm ring-2 ring-ds-bg-elevated',
)

const bgColor = computed(() => {
  const key = props.seed ?? props.name
  let hash = 0
  for (let i = 0; i < key.length; i++) hash = key.charCodeAt(i) + ((hash << 5) - hash)
  const hue = Math.abs(hash) % 360
  return `hsl(${hue} 52% 42%)`
})
</script>
