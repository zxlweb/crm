<template>
  <fieldset class="space-y-2" data-testid="activity-form-sentiment">
    <legend class="mb-1 block text-sm font-medium text-ds-fg">
      {{ $t('activityFormSentiment') }}
    </legend>
    <p class="mb-2 text-xs leading-relaxed text-ds-fg-muted">
      {{ $t('activityFormSentimentHint') }}
    </p>
    <div class="space-y-2" role="radiogroup" :aria-label="$t('activityFormSentiment')">
      <label
        v-for="code in ACTIVITY_SENTIMENTS"
        :key="code"
        class="flex cursor-pointer gap-3 rounded-xl border px-3 py-2.5 transition-colors duration-200"
        :class="
          model === code
            ? 'border-ds-brand bg-ds-brand-subtle/50 ring-1 ring-ds-brand/30'
            : 'border-ds-border-muted bg-ds-bg-elevated hover:border-ds-brand-muted'
        "
        :data-testid="`activity-sentiment-option-${code}`"
      >
        <input
          v-model="model"
          type="radio"
          class="mt-1 h-4 w-4 shrink-0 cursor-pointer accent-ds-brand"
          name="activity-sentiment"
          :value="code"
        />
        <span class="min-w-0 flex-1">
          <span class="flex items-center gap-2 text-sm font-semibold text-ds-fg-heading">
            <UiSentimentEmoji :sentiment="code" size="sm" :label="$t(`sentiment.${code}`)" />
            {{ $t(`sentiment.${code}`) }}
          </span>
          <span class="mt-1 block text-xs leading-relaxed text-ds-fg-muted">
            {{ $t(`sentiment.guide.${code}`) }}
          </span>
        </span>
      </label>
    </div>
  </fieldset>
</template>

<script setup lang="ts">
import { ACTIVITY_SENTIMENTS, type ActivitySentiment } from '~/types/activity'

const model = defineModel<ActivitySentiment>({ required: true })
</script>
