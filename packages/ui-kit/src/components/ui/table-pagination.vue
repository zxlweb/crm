<template>
  <div
    class="flex w-full flex-wrap items-center justify-between gap-3"
    data-testid="table-pagination"
  >
    <span class="text-sm text-ds-fg-muted">{{ rangeText }}</span>
    <UPagination
      :model-value="page"
      :total="total"
      :page-count="pageSize"
      :max="max"
      :disabled="disabled || total === 0"
      :ui="crmPaginationUi"
      :prev-button="prevButtonConfig"
      :next-button="nextButtonConfig"
      @update:model-value="$emit('update:page', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { crmPaginationUi } from '../../theme/nuxt-ui-pagination'

const props = withDefaults(
  defineProps<{
    page: number
    pageSize: number
    total: number
    rangeText?: string
    prevLabel?: string
    nextLabel?: string
    max?: number
    disabled?: boolean
  }>(),
  {
    max: 7,
    disabled: false,
    prevLabel: 'Previous',
    nextLabel: 'Next',
  },
)

defineEmits<{
  'update:page': [page: number]
}>()

const rangeStart = computed(() =>
  props.total === 0 ? 0 : (props.page - 1) * props.pageSize + 1,
)
const rangeEnd = computed(() =>
  props.total === 0 ? 0 : Math.min(props.page * props.pageSize, props.total),
)

const rangeText = computed(
  () =>
    props.rangeText ??
    `${rangeStart.value}-${rangeEnd.value} / ${props.total}`,
)

const prevButtonConfig = computed(() => ({
  ...crmPaginationUi.default.prevButton,
  label: props.prevLabel,
}))

const nextButtonConfig = computed(() => ({
  ...crmPaginationUi.default.nextButton,
  label: props.nextLabel,
}))
</script>
