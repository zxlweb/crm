<template>
  <div>
    <div
      class="sticky top-0 z-10 -mx-1 mb-4 rounded-xl border border-ds-border/80 bg-ds-bg/90 px-1 py-2 backdrop-blur-md supports-[backdrop-filter]:bg-ds-bg/75"
    >
      <UiTabs
        :model-value="modelValue"
        :items="tabs"
        @update:model-value="$emit('update:modelValue', $event as LeadDetailTab)"
      />
    </div>

    <div class="py-2" role="tabpanel">
      <Transition
        mode="out-in"
        enter-active-class="transition-opacity duration-200 ease-out"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition-opacity duration-150 ease-in"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div v-if="modelValue === 'timeline'" key="timeline">
          <slot name="timeline" />
        </div>
        <div v-else-if="modelValue === 'emotion'" key="emotion">
          <slot name="emotion" />
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { LeadDetailTab } from '~/types/lead'

defineProps<{
  modelValue: LeadDetailTab
}>()

defineEmits<{
  'update:modelValue': [value: LeadDetailTab]
}>()

const { t } = useI18n()

const tabs = computed(() => [
  { id: 'timeline' as const, label: t('leadsTabTimeline') },
  { id: 'emotion' as const, label: t('leadsTabEmotion') },
])
</script>
