<template>
  <span
    class="inline-flex shrink-0 items-center justify-center leading-none"
    :class="sizeClass"
    role="img"
    :aria-label="ariaLabel"
  >{{ glyph }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { sentimentEmoji } from '../../icons/sentiment-emoji'

const props = withDefaults(
  defineProps<{
    sentiment: string
    size?: 'xs' | 'sm' | 'md' | 'lg'
    /** 无障碍朗读；默认与 sentiment 文案一致时可由父级传入 */
    label?: string
  }>(),
  { size: 'sm' },
)

const glyph = computed(() => sentimentEmoji(props.sentiment))

const ariaLabel = computed(() => props.label ?? props.sentiment)

const sizeClass = computed(() => {
  const map = {
    xs: 'text-sm',
    sm: 'text-base',
    md: 'text-lg',
    lg: 'text-xl',
  }
  return map[props.size]
})
</script>
