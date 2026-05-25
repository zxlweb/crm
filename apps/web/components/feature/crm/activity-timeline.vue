<template>
  <div class="space-y-4" data-testid="activity-timeline">
    <div v-if="pending" class="flex justify-center py-12">
      <UIcon name="i-heroicons-arrow-path" class="h-6 w-6 animate-spin text-primary" />
    </div>

    <UAlert v-else-if="loadError" color="red" variant="soft" :title="loadError" />

    <div
      v-else-if="items.length === 0"
      class="rounded-xl border border-dashed border-ds-border bg-ds-bg-muted px-6 py-10 text-center"
    >
      <p class="text-sm font-medium text-ds-fg-heading">{{ $t('timelineEmpty') }}</p>
      <p class="mt-1 text-sm text-ds-fg-muted">{{ $t('timelineEmptyHint') }}</p>
    </div>

    <ol v-else class="relative">
      <span
        class="pointer-events-none absolute bottom-0 left-3 top-0 border-l border-dashed border-ds-border"
        aria-hidden="true"
      />
      <li
        v-for="item in items"
        :key="item.id"
        class="relative flex gap-4 pb-6 last:pb-0"
      >
        <div
          class="relative z-10 flex w-6 shrink-0 items-center justify-center self-stretch"
          aria-hidden="true"
        >
          <span
            class="h-3 w-3 shrink-0 rounded-full border-2 border-ds-bg-elevated ring-2 ring-ds-bg-elevated"
            :class="dotClass(item.sentiment)"
          />
        </div>
        <div class="min-w-0 flex-1 rounded-xl border border-ds-border bg-ds-bg-elevated px-4 py-3">
          <div class="flex flex-wrap items-start justify-between gap-2">
            <div class="flex min-w-0 items-center gap-2">
              <span
                class="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-ds-brand-subtle text-ds-fg-brand"
              >
                <UiTagIcon :name="activityIcon(item.event_type)" size="sm" />
              </span>
              <p class="text-sm font-medium text-ds-fg-heading">{{ displayLabel(item) }}</p>
            </div>
            <time class="shrink-0 text-xs text-ds-fg-muted">{{ formatAt(item.occurred_at) }}</time>
          </div>
          <p v-if="item.body" class="mt-1 pl-10 text-sm text-ds-fg-muted">{{ item.body }}</p>
          <div class="mt-2 flex flex-wrap gap-2 pl-10 text-xs text-ds-fg-subtle">
            <span class="inline-flex items-center gap-1 rounded bg-ds-bg-muted px-2 py-0.5">
              <UiTagIcon :name="activityIcon(item.event_type)" size="xs" />
              {{ $t(`activityType.${item.event_type}`) }}
            </span>
            <span
              v-if="item.sentiment"
              class="inline-flex items-center gap-1 rounded px-2 py-0.5"
              :class="sentimentClass(item.sentiment)"
              :data-testid="`activity-sentiment-${item.sentiment}`"
            >
              <UiSentimentEmoji :sentiment="item.sentiment" size="xs" />
              {{ $t(`sentiment.${item.sentiment}`) }}
            </span>
          </div>
        </div>
      </li>
    </ol>
  </div>
</template>

<script setup lang="ts">
import { resolveActivityTypeIcon } from '@crm/ui-kit'
import type { Activity } from '~/types/activity'

function activityIcon(eventType: string) {
  return resolveActivityTypeIcon(eventType)
}

const props = defineProps<{
  subjectType: 'lead' | 'account' | 'contact'
  subjectId: string
}>()

const { t, locale } = useI18n()
const activitiesApi = useActivities()

const items = ref<Activity[]>([])
const pending = ref(true)
const loadError = ref('')

function displayLabel(item: Activity) {
  if (item.label) return item.label
  return t(`activityType.${item.event_type}`)
}

function formatAt(iso: string) {
  return new Intl.DateTimeFormat(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(iso))
}

function dotClass(sentiment?: string | null) {
  const map: Record<string, string> = {
    positive: 'bg-ds-success',
    neutral: 'bg-ds-fg-subtle',
    hesitant: 'bg-ds-warning',
    negative: 'bg-ds-danger',
  }
  return map[sentiment ?? ''] ?? 'bg-ds-brand'
}

function sentimentClass(sentiment: string) {
  const map: Record<string, string> = {
    positive: 'bg-ds-success-subtle text-ds-success',
    neutral: 'bg-ds-bg-muted text-ds-fg-muted',
    hesitant: 'bg-ds-warning-subtle text-ds-warning',
    negative: 'bg-ds-danger-subtle text-ds-danger',
  }
  return map[sentiment] ?? 'bg-ds-bg-muted text-ds-fg-muted'
}

async function load() {
  pending.value = true
  loadError.value = ''
  try {
    const { items: rows } = await activitiesApi.fetchList({
      subjectType: props.subjectType,
      subjectId: props.subjectId,
      pageSize: 50,
    })
    items.value = rows
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
    items.value = []
  } finally {
    pending.value = false
  }
}

watch(() => [props.subjectType, props.subjectId], load, { immediate: true })
defineExpose({ reload: load })
</script>
