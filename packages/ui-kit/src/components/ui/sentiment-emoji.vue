<template>
  <UiTagIcon
    :name="iconName"
    :size="tagSize"
    :title="ariaLabel"
    class="inline-flex shrink-0"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { resolveSentimentIcon } from '../../icons/tag-icons'
import UiTagIcon from './tag-icon.vue'

const props = withDefaults(
  defineProps<{
    sentiment: string
    size?: 'xs' | 'sm' | 'md' | 'lg'
    /** 无障碍朗读 */
    label?: string
  }>(),
  { size: 'sm' },
)

const iconName = computed(() => resolveSentimentIcon(props.sentiment))

const ariaLabel = computed(() => props.label ?? props.sentiment)

const tagSize = computed((): 'xs' | 'sm' | 'md' => {
  if (props.size === 'xs') return 'xs'
  if (props.size === 'lg') return 'md'
  return props.size === 'md' ? 'md' : 'sm'
})
</script>
