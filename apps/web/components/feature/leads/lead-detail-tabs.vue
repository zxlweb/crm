<template>
  <div>
    <UiTabs
      :model-value="modelValue"
      :items="tabs"
      class="mb-6"
      @update:model-value="$emit('update:modelValue', $event as LeadDetailTab)"
    />

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
        <div v-if="modelValue === 'overview'" key="overview">
          <slot name="overview" />
        </div>
        <div v-else-if="modelValue === 'timeline'" key="timeline">
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
  { id: 'overview' as const, label: t('leadsTabOverview') },
  { id: 'timeline' as const, label: t('leadsTabTimeline') },
  { id: 'emotion' as const, label: t('leadsTabEmotion') },
])

/** QA: emotion tab — pages/leads/[id].vue 内保留 data-testid="tab-emotion-journey" 于 slot */
</script>
