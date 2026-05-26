<template>
  <section
    class="ds-lead-pulse relative overflow-hidden rounded-2xl border border-ds-border-muted bg-ds-bg-elevated p-4 shadow-ds-sm sm:p-5"
    data-testid="lead-pulse-bar"
    :aria-label="$t('leadsPulseAria')"
  >
    <span
      class="pointer-events-none absolute inset-x-0 top-0 h-0.5"
      :style="{ background: 'var(--ds-brand-gradient)' }"
      aria-hidden="true"
    />
    <span
      class="pointer-events-none absolute -right-12 -top-12 h-36 w-36 rounded-full opacity-15 blur-3xl"
      :style="{ background: 'var(--ds-brand-gradient)' }"
      aria-hidden="true"
    />

    <div class="relative flex flex-col gap-4 lg:flex-row lg:items-center">
      <!-- Composite score gauge -->
      <div class="flex shrink-0 items-center gap-3">
        <div class="relative h-16 w-16 shrink-0">
          <svg
            class="h-full w-full -rotate-90"
            viewBox="0 0 36 36"
            aria-hidden="true"
          >
            <circle
              cx="18"
              cy="18"
              r="15.915"
              fill="none"
              stroke="var(--ds-border-muted)"
              stroke-width="3"
            />
            <circle
              cx="18"
              cy="18"
              r="15.915"
              fill="none"
              :stroke="ringStroke"
              stroke-width="3"
              stroke-linecap="round"
              :stroke-dasharray="`${pulseScore} 100`"
              class="transition-[stroke-dasharray] duration-700 ease-out"
            />
          </svg>
          <div
            class="absolute inset-0 flex flex-col items-center justify-center"
          >
            <span class="text-base font-bold tabular-nums" :class="bandText">
              {{ pulseScore }}
            </span>
            <span class="text-[9px] font-medium uppercase tracking-wide text-ds-fg-subtle">
              {{ $t('leadsPulseScoreSuffix') }}
            </span>
          </div>
        </div>
        <div class="min-w-0">
          <p class="text-[11px] font-semibold uppercase tracking-[0.14em] text-ds-fg-muted">
            {{ $t('leadsPulseTitle') }}
          </p>
          <p class="mt-0.5 text-sm font-semibold leading-snug" :class="bandText">
            {{ bandLabel }}
          </p>
          <p class="mt-0.5 max-w-[18rem] text-[11px] leading-relaxed text-ds-fg-muted">
            {{ bandHint }}
          </p>
        </div>
      </div>

      <!-- Composite ribbon -->
      <div class="min-w-0 flex-1">
        <div
          class="relative h-2.5 overflow-hidden rounded-full bg-ds-bg-muted ring-1 ring-inset ring-ds-border-muted"
          :aria-label="$t('leadsPulseRibbonAria')"
        >
          <span
            v-for="seg in segments"
            :key="seg.key"
            class="absolute inset-y-0 transition-[left,width] duration-700 ease-out"
            :class="seg.barClass"
            :style="{ left: `${seg.offset}%`, width: `${seg.width}%` }"
            :title="seg.tooltip"
          />
        </div>

        <ul class="mt-3 grid grid-cols-2 gap-2 sm:grid-cols-4">
          <li
            v-for="seg in segments"
            :key="`legend-${seg.key}`"
            class="flex min-w-0 items-start gap-2 rounded-xl border border-ds-border-muted bg-ds-bg-elevated/60 px-2.5 py-2"
          >
            <span
              class="mt-1 h-2 w-2 shrink-0 rounded-full"
              :class="seg.dotClass"
              aria-hidden="true"
            />
            <div class="min-w-0 flex-1">
              <p class="truncate text-[10px] font-semibold uppercase tracking-wide text-ds-fg-subtle">
                {{ seg.label }}
              </p>
              <p class="mt-0.5 truncate text-xs font-semibold tabular-nums" :class="seg.valueClass">
                {{ seg.valueLabel }}
              </p>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { Lead } from '~/types/lead'

interface PulseSegment {
  key: 'engagement' | 'recency' | 'lifecycle' | 'intent'
  label: string
  score: number
  weight: number
  offset: number
  width: number
  valueLabel: string
  valueClass: string
  barClass: string
  dotClass: string
  tooltip: string
}

const props = defineProps<{
  lead: Lead
}>()

const { t } = useI18n()

const LIFECYCLE_SCORE: Record<Lead['lifecycle_stage'], number> = {
  acquire: 30,
  activate: 55,
  grow: 80,
  retain: 95,
  revive: 45,
}

const STATUS_INTENT_SCORE: Record<Lead['status'], number> = {
  new: 35,
  contacted: 55,
  qualified: 85,
  unqualified: 10,
  converted: 100,
}

const LOW_INTENT_TAG = '低意向'

const daysIdle = computed<number | null>(() => {
  if (!props.lead.last_activity_at) return null
  return Math.max(
    0,
    Math.floor((Date.now() - new Date(props.lead.last_activity_at).getTime()) / 86_400_000),
  )
})

const recencyScore = computed(() => {
  if (daysIdle.value == null) return 25
  if (daysIdle.value === 0) return 100
  if (daysIdle.value <= 3) return 85
  if (daysIdle.value <= 7) return 65
  if (daysIdle.value <= 14) return 40
  if (daysIdle.value <= 30) return 20
  return 10
})

const lifecycleScore = computed(() => LIFECYCLE_SCORE[props.lead.lifecycle_stage] ?? 30)

const intentScore = computed(() => {
  let score = STATUS_INTENT_SCORE[props.lead.status] ?? 30
  const hasLowIntent = (props.lead.tags ?? []).some(
    (tag) => tag === LOW_INTENT_TAG || tag.toLowerCase() === 'low intent',
  )
  if (hasLowIntent) score = Math.max(15, score - 25)
  return score
})

const engagementScore = computed(() => Math.max(0, Math.min(100, props.lead.engagement_score)))

const WEIGHTS = {
  engagement: 0.35,
  recency: 0.3,
  lifecycle: 0.15,
  intent: 0.2,
}

const pulseScore = computed(() => {
  const weighted =
    engagementScore.value * WEIGHTS.engagement +
    recencyScore.value * WEIGHTS.recency +
    lifecycleScore.value * WEIGHTS.lifecycle +
    intentScore.value * WEIGHTS.intent
  return Math.round(weighted)
})

function bandFor(score: number): 'critical' | 'watch' | 'steady' | 'hot' {
  if (score < 30) return 'critical'
  if (score < 55) return 'watch'
  if (score < 75) return 'steady'
  return 'hot'
}

const band = computed(() => bandFor(pulseScore.value))

const bandText = computed(() => {
  switch (band.value) {
    case 'critical':
      return 'text-ds-danger'
    case 'watch':
      return 'text-ds-warning'
    case 'steady':
      return 'text-ds-info'
    default:
      return 'text-ds-success'
  }
})

const ringStroke = computed(() => {
  switch (band.value) {
    case 'critical':
      return 'var(--ds-danger)'
    case 'watch':
      return 'var(--ds-warning)'
    case 'steady':
      return 'var(--ds-info)'
    default:
      return 'var(--ds-success)'
  }
})

const bandLabel = computed(() => t(`leadsPulseBand.${band.value}`))
const bandHint = computed(() => t(`leadsPulseBandHint.${band.value}`))

function recencyValueLabel(): string {
  if (daysIdle.value == null) return t('leadsMetricNoActivity')
  if (daysIdle.value === 0) return t('leadsHeroIdleToday')
  return t('leadsMetricDaysValue', { days: daysIdle.value })
}

function toneFor(score: number) {
  if (score < 30) {
    return {
      barClass: 'bg-ds-danger',
      dotClass: 'bg-ds-danger',
      valueClass: 'text-ds-danger',
    }
  }
  if (score < 55) {
    return {
      barClass: 'bg-ds-warning',
      dotClass: 'bg-ds-warning',
      valueClass: 'text-ds-warning',
    }
  }
  if (score < 75) {
    return {
      barClass: 'bg-ds-info',
      dotClass: 'bg-ds-info',
      valueClass: 'text-ds-info',
    }
  }
  return {
    barClass: 'bg-ds-success',
    dotClass: 'bg-ds-success',
    valueClass: 'text-ds-success',
  }
}

const segments = computed<PulseSegment[]>(() => {
  const raw = [
    {
      key: 'engagement' as const,
      label: t('leadsPulseSegmentEngagement'),
      score: engagementScore.value,
      weight: WEIGHTS.engagement,
      valueLabel: t('leadsPulseEngagementValue', { score: engagementScore.value }),
    },
    {
      key: 'recency' as const,
      label: t('leadsPulseSegmentRecency'),
      score: recencyScore.value,
      weight: WEIGHTS.recency,
      valueLabel: recencyValueLabel(),
    },
    {
      key: 'lifecycle' as const,
      label: t('leadsPulseSegmentLifecycle'),
      score: lifecycleScore.value,
      weight: WEIGHTS.lifecycle,
      valueLabel: t(`lifecycle.${props.lead.lifecycle_stage}`),
    },
    {
      key: 'intent' as const,
      label: t('leadsPulseSegmentIntent'),
      score: intentScore.value,
      weight: WEIGHTS.intent,
      valueLabel: t(`leadStatus.${props.lead.status}`),
    },
  ]

  let cursor = 0
  return raw.map((row) => {
    const contribution = row.weight * row.score
    const tone = toneFor(row.score)
    const seg: PulseSegment = {
      ...row,
      offset: cursor,
      width: contribution,
      barClass: tone.barClass,
      dotClass: tone.dotClass,
      valueClass: tone.valueClass,
      tooltip: `${row.label} · ${row.valueLabel} · ${Math.round(row.score)}/100`,
    }
    cursor += contribution
    return seg
  })
})
</script>
