<template>
  <form class="space-y-4" data-testid="activity-form" @submit.prevent="onSubmit">
    <div>
      <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('activityFormType') }}</label>
      <UiSelect v-model="eventType" :items="eventTypeItems" data-testid="activity-form-event-type" />
    </div>
    <div>
      <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('activityFormDirection') }}</label>
      <UiSelect v-model="direction" :items="directionItems" />
    </div>
    <div>
      <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('activityFormBody') }}</label>
      <textarea
        v-model="body"
        rows="3"
        class="w-full rounded-lg border border-ds-border bg-ds-bg-input px-3 py-2 text-sm text-ds-fg focus:border-ds-brand focus:outline-none focus:ring-1 focus:ring-ds-brand"
        data-testid="activity-form-body"
      />
    </div>
    <div>
      <label class="mb-1.5 block text-sm font-medium text-ds-fg">{{ $t('activityFormSentiment') }}</label>
      <UiSelect
        v-model="sentiment"
        :items="sentimentItems"
        :placeholder="$t('activityFormSentimentPlaceholder')"
        data-testid="activity-form-sentiment"
      />
    </div>
    <div class="flex justify-end gap-2 border-t border-ds-border pt-4">
      <UiButton type="button" variant="secondary" data-testid="activity-form-cancel" @click="emit('cancel')">
        {{ $t('cancel') }}
      </UiButton>
      <UiButton type="submit" :loading="loading" data-testid="activity-form-submit">
        {{ $t('save') }}
      </UiButton>
    </div>
  </form>
</template>

<script setup lang="ts">
import type { ActivityCreateInput, ActivityEventType, ActivitySentiment } from '~/types/activity'

const props = defineProps<{
  subjectType: ActivityCreateInput['subject_type']
  subjectId: string
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: ActivityCreateInput]
  cancel: []
}>()

const { t } = useI18n()

const eventType = ref<ActivityEventType>('call')
const direction = ref('outbound')
const body = ref('')
const sentiment = ref('')

const eventTypes: ActivityEventType[] = ['call', 'email', 'meeting', 'wechat', 'note', 'visit']
const sentiments: ActivitySentiment[] = ['positive', 'neutral', 'hesitant', 'negative']

const eventTypeItems = computed(() =>
  eventTypes.map((e) => ({ label: t(`activityType.${e}`), value: e })),
)
const directionItems = computed(() => [
  { label: t('activityDirection.outbound'), value: 'outbound' },
  { label: t('activityDirection.inbound'), value: 'inbound' },
])
const sentimentItems = computed(() => [
  { label: t('activityFormSentimentNone'), value: '' },
  ...sentiments.map((s) => ({ label: t(`sentiment.${s}`), value: s })),
])

function resetForm() {
  eventType.value = 'call'
  direction.value = 'outbound'
  body.value = ''
  sentiment.value = ''
}

function onSubmit() {
  const payload: ActivityCreateInput = {
    subject_type: props.subjectType,
    subject_id: props.subjectId,
    event_type: eventType.value,
    direction: direction.value as ActivityCreateInput['direction'],
    body: body.value.trim() || undefined,
  }
  if (sentiment.value) {
    payload.sentiment = sentiment.value as ActivitySentiment
    payload.sentiment_source = 'manual'
  }
  emit('submit', payload)
  resetForm()
}

defineExpose({ submit: onSubmit, reset: resetForm })
</script>
